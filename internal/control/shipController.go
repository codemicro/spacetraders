package control

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/stapi"
	"github.com/imdario/mergo"
	"strings"
	"time"
)

type ShipController struct {
	ship *stapi.Ship
	core *Core
}

func NewShipController(ship *stapi.Ship, core *Core) *ShipController {
	s := new(ShipController)
	s.ship = ship
	s.core = core

	go s.Start()

	return s
}

func (s *ShipController) log(format string, a ...interface{}) {
	prefix := s.ship.ID[:6] + ": "
	x := strings.ReplaceAll(fmt.Sprintf(format, a...), "\n", "\n"+strings.Repeat(" ", len(prefix)))
	s.core.Log("%s%s\n", prefix, x)
}

func (s *ShipController) buyGood(good string, quantity int) error {
	newShip, err := s.core.user.SubmitPurchaseOrder(s.ship.ID, good, quantity)
	if err != nil {
		return err
	}
	return mergo.Merge(s.ship, newShip)
}

func (s *ShipController) sellGood(good string, quantity int) (*stapi.Order, error) {
	newShip, order, err := s.core.user.SubmitSellOrder(s.ship.ID, good, quantity)
	if err != nil {
		return nil, err
	}
	if err = mergo.Merge(s.ship, newShip); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *ShipController) refuel(amount int) error {
	return s.buyGood("FUEL", amount)
}

func (s *ShipController) fileFlightplan(fp *plannedFlight) (*stapi.Flightplan, error) {
	return s.core.user.SubmitFlightplan(s.ship.ID, fp.destination.Symbol)
}

func (s *ShipController) Start() {
	s.log("online at %s (%d,%d)", s.ship.Location, s.ship.XCoordinate, s.ship.YCoordinate)

	fp, err := s.planFlight()
	if err != nil {
		s.log("ERROR: %s", err.Error()) // TODO: nice error handling
		return
	}

	cargoString := "none"
	if fp.cargo != nil {
		cargoString = fmt.Sprintf("%s (%d units)", fp.cargo.Symbol, fp.unitsCargo)
	}

	s.log(
		"flightplan created\nCost: %dcr\nExtra fuel: %d units\nCargo: %s\nDestination: %s (%s)\nDistance: %d",
		fp.flightCost,
		fp.extraFuelRequired,
		cargoString,
		fp.destination.Name,
		fp.destination.Symbol,
		fp.distance,
	)

	s.log("waiting 5 seconds for cancellation...")
	time.Sleep(time.Second * 5)

	s.log("preparing for flight")

	for _, task := range fp.preflightTasks {
		if err = task(); err != nil {
			s.log("ERROR: %s", err.Error()) // TODO: nice error handling
			return
		}
	}

	flightplan, err := s.fileFlightplan(fp)
	if err != nil {
		s.log("ERROR: %s", err.Error()) // TODO: nice error handling
		return
	}

	s.log("departing...\nFlightplan ID: %s", flightplan.ID)

	sleepDuration := time.Minute
	totalFlightDuration := flightplan.ArrivesAt.Sub(*flightplan.CreatedAt)
	for {

		flightplan, err = s.core.user.GetFlightplan(flightplan.ID)
		if err != nil {
			s.log("ERROR: %s", err.Error()) // TODO: nice error handling
			return
		}

		if ut := time.Until(*flightplan.ArrivesAt); ut < sleepDuration {
			sleepDuration = ut + (time.Second * 2)
		}

		var percentageComplete float32
		{
			durationFlown := time.Since(*flightplan.CreatedAt)
			percentageComplete = float32(durationFlown) / float32(totalFlightDuration) * 100
		}

		if flightplan.TerminatedAt != nil {
			s.log("arrived at %s", flightplan.TerminatedAt.Format(time.Kitchen))
			break
		}

		s.log("en route - %.2f%% complete, %ds remaining", percentageComplete, flightplan.FlightTimeRemaining)
		time.Sleep(sleepDuration)
	}

	if fp.cargo != nil {

		s.log("selling cargo")

		order, err := s.sellGood(fp.cargo.Symbol, fp.unitsCargo)
		if err != nil {
			s.log("ERROR: %s", err.Error()) // TODO: nice error handling
			return
		}

		s.log("sold for %dcr - profit of %dcr", order.Total, order.Total - fp.flightCost)

	}

	s.log("done")
}
