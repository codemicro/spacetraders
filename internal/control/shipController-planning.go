package control

import (
	"errors"
	"github.com/codemicro/spacetraders/internal/analysis"
	"github.com/codemicro/spacetraders/internal/stapi"
	"strings"
)

type plannedFlight struct {
	preflightTasks    []func() error
	destination       *stapi.Location
	cargo             *stapi.MarketplaceGood
	extraFuelRequired int
	unitsCargo        int
	flightCost        int
	distance          int
}

var (
	ErrorCannotPlanRoute = errors.New("shipController: could not plan route from specified destination")
	ErrorNotEnoughFuel   = errors.New("shipController: not enough fuel for journey and cannot buy any more")
	ErrorCannotPickCargo = errors.New("shipController: could not choose a cargo (this is probably a programming error")
)

const cargoSpendLimit = 8000

func (s *ShipController) planFlight(destinationString string) (*plannedFlight, error) {
	fp := new(plannedFlight)

	currentLocation, err := stapi.GetLocationInfo(s.ship.Location)
	if err != nil {
		return nil, err
	}

	marketplace, err := stapi.GetMarketplaceAtLocation(currentLocation.Symbol)
	if err != nil {
		return nil, err
	}

	destination, err := stapi.GetLocationInfo(destinationString)
	if err != nil {
		return nil, err
	}

	flightDistance := analysis.FindDistance(currentLocation, destination)

	fp.destination = destination
	fp.distance = flightDistance
	fp.cargo = nil

	if err = s.planFuel(fp, currentLocation, marketplace); err != nil {
		return nil, err
	}

	return fp, nil
}

func (s *ShipController) planCargoFlight() (*plannedFlight, error) {

	fp := new(plannedFlight)

	currentLocation, err := stapi.GetLocationInfo(s.ship.Location)
	if err != nil {
		return nil, err
	}

	marketplace, err := stapi.GetMarketplaceAtLocation(currentLocation.Symbol)
	if err != nil {
		return nil, err
	}

	destination, cargo, err := analysis.FindCombinedRouteAndCargo(s.ship.Location)
	if err != nil {
		return nil, err
	}

	flightDistance := analysis.FindDistance(currentLocation, destination)

	fp.destination = destination
	fp.distance = flightDistance
	fp.cargo = cargo

	if err = s.planFuel(fp, currentLocation, marketplace); err != nil {
		return nil, err
	}

	unitsToBuy := (s.ship.SpaceAvailable - fp.extraFuelRequired) / fp.cargo.VolumePerUnit
	for {
		cost := fp.cargo.PurchasePricePerUnit * unitsToBuy
		if cost > cargoSpendLimit {
			unitsToBuy -= 1
		} else {
			break
		}
	}

	fp.unitsCargo = unitsToBuy
	fp.flightCost += fp.cargo.PurchasePricePerUnit * fp.unitsCargo

	fp.preflightTasks = append(fp.preflightTasks, func() error {
		s.log("purchasing %d units of cargo %s", fp.unitsCargo, fp.cargo.Symbol)
		return s.buyGood(fp.cargo.Symbol, fp.unitsCargo)
	})

	return fp, nil
}

func (s *ShipController) planRoute(fp *plannedFlight, currentLocation *stapi.Location, system string, method analysis.RoutingMethod) error {
	locationsInThisSystem, err := stapi.GetSystemLocations(system)
	if err != nil {
		return err
	}

	flightDestination := analysis.PickRoute(currentLocation, locationsInThisSystem, method)
	if flightDestination == nil {
		return ErrorCannotPlanRoute
	}

	flightDistance := analysis.FindDistance(currentLocation, flightDestination)

	fp.destination = flightDestination
	fp.distance = flightDistance

	return nil
}

func (s *ShipController) planFuel(fp *plannedFlight, currentLocation *stapi.Location, marketplace []*stapi.MarketplaceGood) error {
	journeyFuel := analysis.CalculateFuelForFlight(currentLocation, fp.destination)
	extraFuelRequired := journeyFuel - s.ship.GetCurrentFuel()

	if extraFuelRequired > 0 {
		fp.extraFuelRequired = extraFuelRequired

		var fuelCost int
		for _, g := range marketplace {

			if strings.EqualFold(g.Symbol, "FUEL") {
				fuelCost = g.PurchasePricePerUnit
				break
			}
		}
		if fuelCost == 0 {
			return ErrorNotEnoughFuel
		}

		fp.flightCost += fuelCost * extraFuelRequired

		fp.preflightTasks = append(fp.preflightTasks, func() error {
			s.log("fuelling with %d units of fuel", extraFuelRequired)
			return s.refuel(extraFuelRequired)
		})
	}

	return nil
}
