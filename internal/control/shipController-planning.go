package control

import (
	"errors"
	"github.com/codemicro/spacetraders/internal/analysis"
	"github.com/codemicro/spacetraders/internal/stapi"
	"strings"
)

type flightplan struct {
	preflightTasks []func() error
	destination    *stapi.Location
	flightCost     int
	distance       int
}

var (
	ErrorCannotPlanRoute = errors.New("shipController: could not plan route from specified destination")
	ErrorNotEnoughFuel   = errors.New("shipController: not enough fuel for journey and cannot buy any more")
)

func (s *ShipController) planFlight() (*flightplan, error) {

	fp := new(flightplan)

	locationsInThisSystem, err := stapi.GetSystemLocations(strings.Split(s.ship.Location, "-")[0])
	if err != nil {
		return nil, err
	}

	currentLocation, err := stapi.GetLocationInfo(s.ship.Location)
	if err != nil {
		return nil, err
	}

	flightDestination := analysis.PickRoute(currentLocation, locationsInThisSystem, analysis.RoutingMethodShort)
	if flightDestination == nil {
		return nil, ErrorCannotPlanRoute
	}

	flightDistance := analysis.FindDistance(currentLocation, flightDestination)

	fp.destination = flightDestination
	fp.distance = flightDistance

	journeyFuel := analysis.CalculateFuelForFlight(currentLocation, flightDestination)
	extraFuelRequired := journeyFuel - s.ship.GetCurrentFuel()

	marketplace, err := stapi.GetMarketplaceAtLocation(currentLocation.Symbol)
	if err != nil {
		return nil, err
	}

	if extraFuelRequired > 0 {

		var fuelCost int
		for _, g := range marketplace {
			if strings.EqualFold(g.Symbol, "FUEL") {
				fuelCost = g.PurchasePricePerUnit
				break
			}
		}
		if fuelCost == 0 {
			return nil, ErrorNotEnoughFuel
		}

		fp.preflightTasks = append(fp.preflightTasks, func() error {
			s.log("fuelling with %d units of fuel", extraFuelRequired)
			return s.refuel(extraFuelRequired)
		})

	}

	// TODO: y'know, uuuhh, cargo???

	return fp, nil
}
