package control

import (
	"github.com/codemicro/spacetraders/internal/analysis"
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
	s.core.Log("%s: "+format+"\n", append([]interface{}{s.ship.ID[:6]}, a...)...)
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

	locationsInThisSystem, err := stapi.GetSystemLocations(strings.Split(s.ship.Location, "-")[0])
	if err != nil {
		s.log(err.Error()) // TODO: nice error handling
		return
	}

	currentLocation, err := stapi.GetLocationInfo(s.ship.Location)
	if err != nil {
		s.log(err.Error()) // TODO: nice error handling
		return
	}

	flightDestination := analysis.PickRoute(currentLocation, locationsInThisSystem, analysis.RoutingMethodShort)
	flightDistance := analysis.FindDistance(currentLocation, flightDestination)

	journeyFuel := analysis.CalculateFuelForFlight(currentLocation, flightDestination)
	extraFuelRequired := journeyFuel - s.ship.GetCurrentFuel()

	s.log("found destination: dist %d, %#v", flightDistance, flightDestination)
	s.log("total fuel required for journey: %d", journeyFuel)

	if extraFuelRequired > 0 {
		s.log("fuelling with %d units of fuel", extraFuelRequired)
		err = s.refuel(extraFuelRequired)
		if err != nil {
			s.log(err.Error()) // TODO: nice error handling
			return
		}
	} else {
		s.log("no extra fuel required")
	}

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
