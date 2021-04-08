package control

import (
	"errors"
	"fmt"
	"github.com/codemicro/spacetraders/internal/analysis"
	"github.com/codemicro/spacetraders/internal/stapi"
	"github.com/codemicro/spacetraders/internal/tool"
	"github.com/imdario/mergo"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"strings"
	"time"
)

const (
	doNotUse = iota // if there is a ship type with 0, it is treated as all ships by GORM
	ShipTypeTrader
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
	prefix := s.ship.ID + ": "
	if s.shipType == ShipTypeProbe {
		format = "(PROBE) " + format
	}
	x := strings.ReplaceAll(fmt.Sprintf(format, a...), "\n", "\n"+strings.Repeat(" ", len(prefix)))
	var highlighted string
	if s.shipType == ShipTypeProbe {
		highlighted = aurora.BrightBlue(prefix).String()
	} else {
		highlighted = aurora.Yellow(prefix).String()
	}
	s.core.WriteToStdout("%s%s\n", highlighted, x)
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

	for s.core.allowStartNewFlight {
		var fp *plannedFlight

		if s.shipType == ShipTypeProbe {
			if s.ship.Location != s.data {
				var err error
				s.log("planning flight to %s", s.data)
				fp, err = s.planFlight(s.data)
				if err != nil {
					s.error(err)
					s.core.stopNotifier <- s.ship.ID
					return
				}
			} else {
				s.log("already at target location %s (%d,%d), not moving", s.ship.Location, s.ship.XCoordinate, s.ship.YCoordinate)
				time.Sleep(time.Second * time.Duration(rand.Intn(59))) // so we don't get a huge barrage of requests all at once
				s.probeAction()
				s.core.stopNotifier <- s.ship.ID
				return // means that this thread can stop whenever it wants
			}

		} else {
			var err error
			s.log("planning cargo flight")

			var attempts int
			for {

				if !s.core.allowStartNewFlight {
					// If we've been sleeping and a shutdown has been requested, this needs to stop
					s.core.stopNotifier <- s.ship.ID
					return
				}

				fp, err = s.planCargoFlight()
				if err != nil {
					if errors.Is(err, analysis.ErrorNoSuitableRoutes) {
						attempts += 1

						if attempts == 2 {
							s.log("2 minutes have passed without being able to find a route. Moving to a new location")
							fp, err = s.planFlightWithMethod(analysis.RoutingMethod(rand.Intn(3))) // the random is in an effort to prevent any infinite looping
							if err != nil {
								s.error(err)
								s.core.stopNotifier <- s.ship.ID
								return
							}
							break
						}

						s.log("did not find any suitable routes at %s, waiting one minute and trying again", s.ship.Location)
						time.Sleep(time.Minute)
						continue
					}
					s.error(err)
					s.core.stopNotifier <- s.ship.ID
					return
				} else {
					break
				}
			}

		}

		cargoString := "none"
		if fp.cargo != nil {
			cargoString = fmt.Sprintf("%s (%d units)", fp.cargo.Symbol, fp.unitsCargo)
		}

		s.log(
			"flightplan created\nCost: %dcr\nPredicted profit: %dcr\nExtra fuel: %d units\nCargo: %s\nDeparture %s\nDestination: %s (%s)\nDistance: %d",
			fp.flightCost,
			fp.expectedProfit,
			fp.extraFuelRequired,
			cargoString,
			s.ship.Location,
			fp.destination.Name,
			fp.destination.Symbol,
			fp.distance,
		)

		if err := s.doFlight(fp); err != nil {
			s.error(err)
			s.core.stopNotifier <- s.ship.ID
			return
		}

		if fp.cargo != nil {

			s.log("selling cargo")

			order, err := s.sellGood(fp.cargo.Symbol, fp.unitsCargo)
			if err != nil {
				s.error(err)
				s.core.stopNotifier <- s.ship.ID
				return
			}

			s.log("sold for %dcr", order.Total)
			s.core.ReportProfit(order.Total-fp.flightCost)

		} else {
			s.core.ReportProfit(-fp.flightCost)
		}

		s.log("updating ship information")
		if err := s.updateShipInfo(); err != nil {
			s.error(err)
			s.core.stopNotifier <- s.ship.ID
			return
		}

		s.log("done")
	}

	s.core.stopNotifier <- s.ship.ID

}
