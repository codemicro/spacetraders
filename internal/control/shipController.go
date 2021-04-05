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

	isScout bool
}

func NewShipController(ship *stapi.Ship, core *Core, scout bool) *ShipController {
	s := new(ShipController)
	s.ship = ship
	s.core = core
	s.isScout = scout

	go s.Start()

	return s
}

func (s *ShipController) log(format string, a ...interface{}) {
	prefix := s.ship.ID[:6] + ": "
	if s.isScout {
		format = "(SCOUT) " + format
	}
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

	for s.isScout {
		// runs while in scout mode
		if err := s.doScout(); err != nil {
			s.log("ERROR: %s", err.Error()) // TODO: nice error handling
			return
		}
	}

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

	if err = s.doFlight(fp); err != nil {
		s.log("ERROR: %s", err.Error()) // TODO: nice error handling
		return
	}

	if fp.cargo != nil {

		s.log("selling cargo")

		order, err := s.sellGood(fp.cargo.Symbol, fp.unitsCargo)
		if err != nil {
			s.log("ERROR: %s", err.Error()) // TODO: nice error handling
			return
		}

		s.log("sold for %dcr - profit of %dcr", order.Total, order.Total-fp.flightCost)

	}

	s.log("done")
}
