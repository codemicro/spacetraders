package control

import (
	"github.com/codemicro/spacetraders/internal/stapi"
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

func (s *ShipController) Start() {
	s.core.Log("%s: online at %s (%d,%d)\n", s.ship.ID[:6], s.ship.Location, s.ship.XCoordinate, s.ship.YCoordinate)
}
