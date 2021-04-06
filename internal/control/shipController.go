package control

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/stapi"
	"github.com/codemicro/spacetraders/internal/tool"
	"github.com/imdario/mergo"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const (
	ShipTypeTrader = iota
	ShipTypeProbe
)

type ShipController struct {
	ship *stapi.Ship
	core *Core

	logger zerolog.Logger

	shipType int
	data string
}

func NewShipController(ship *stapi.Ship, core *Core, shipType int, data string) *ShipController {
	s := new(ShipController)
	s.ship = ship
	s.core = core

	s.shipType = shipType
	s.data = data

	s.logger = log.With().Str("area", "ShipController").Str("shipID", s.ship.ID).Logger()

	go s.Start()

	return s
}

func (s *ShipController) log(format string, a ...interface{}) {
	prefix := s.ship.ID[:6] + ": "
	if s.shipType == ShipTypeProbe {
		format = "(PROBE) " + format
	}
	x := strings.ReplaceAll(fmt.Sprintf(format, a...), "\n", "\n"+strings.Repeat(" ", len(prefix)))
	s.core.Log("%s%s\n", aurora.Yellow(prefix), x)
}

func (s *ShipController) error(err error) {
	s.logger.Error().Err(err).Msg(tool.GetContext(2))
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
	if err = mergo.Merge(s.ship, newShip); err != nil { // TODO: check - is this actually doing the correct thing?
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

func (s *ShipController) updateShipInfo() error {
	newShip, err := s.core.user.GetShipInfo(s.ship.ID)
	if err != nil {
		return err
	}
	*s.ship = *newShip
	return nil
}

func (s *ShipController) Start() {
	s.log("online at %s (%d,%d)", s.ship.Location, s.ship.XCoordinate, s.ship.YCoordinate)

	if s.shipType == ShipTypeProbe {
		if s.ship.Location != s.data {

		}
	}

	fp, err := s.planCargoFlight()
	if err != nil {
		s.error(err)
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
		s.error(err)
		return
	}

	if fp.cargo != nil {

		s.log("selling cargo")

		order, err := s.sellGood(fp.cargo.Symbol, fp.unitsCargo)
		if err != nil {
			s.error(err)
			return
		}

		s.log("sold for %dcr - profit of %dcr", order.Total, order.Total-fp.flightCost)

	}

	s.log("done")
}
