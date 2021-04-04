package control

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/stapi"
	"github.com/imdario/mergo"
	"strings"
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

func (s *ShipController) refuel(amount int) error {
	return s.buyGood("FUEL", amount)
}

func (s *ShipController) Start() {
	s.log("online at %s (%d,%d)", s.ship.Location, s.ship.XCoordinate, s.ship.YCoordinate)

	fp, err := s.planFlight()
	if err != nil {
		s.log("ERROR: %s", err.Error()) // TODO: nice error handling
		return
	}

	s.log("flightplan created\nCost: %d\nDestination: %s (%s)\nDistance: %d", fp.flightCost, fp.destination.Name, fp.destination.Symbol, fp.distance)

	//goods, err := stapi.GetMarketplaceAtLocation(s.ship.Location)
	//if err != nil {
	//	s.log(err.Error()) // TODO: nice error handling
	//	return
	//}
	//
	//for _, good := range goods {
	//	s.log("good %#v", good)
	//}
}
