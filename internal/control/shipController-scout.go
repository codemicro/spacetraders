package control

import (
	"github.com/codemicro/spacetraders/internal/analysis"
	"github.com/codemicro/spacetraders/internal/stapi"
	"strings"
	"time"
)

// This file contains ship controller functionality that scouts for locations without a known marketplace

func (s *ShipController) doScout() error {

	fp := new(plannedFlight)

	// TODO: integrate currentLocation into ShipController
	currentLocation, err := stapi.GetLocationInfo(s.ship.Location)
	if err != nil {
		return err
	}

	// find locations with an unknown market in the current system
	{
		currentSystem := strings.Split(s.ship.Location, "-")[0]

		possibleLocations, err := stapi.GetSystemLocations(currentSystem)
		if err != nil {
			return err
		}

		// filter possible locations

		{
			n := 0
			for _, loc := range possibleLocations {
				if _, found := analysis.GetMarketplaceAtLocation(loc.Symbol); !found && loc.Type != stapi.LocationTypeWormhole {
					possibleLocations[n] = loc
					n++
				}
			}
			possibleLocations = possibleLocations[:n]
		}

		if len(possibleLocations) == 0 {
			// if we've visited all locations, we no longer need to scout
			s.log("finished scouting")
			s.isScout = false
			return nil
		} else {
			s.log("%d locations left to scout", len(possibleLocations))
		}

		flightDestination := analysis.PickRoute(currentLocation, possibleLocations, analysis.RoutingMethodShortest)
		if flightDestination == nil {
			return ErrorCannotPlanRoute
		}

		flightDistance := analysis.FindDistance(currentLocation, flightDestination)

		fp.destination = flightDestination
		fp.distance = flightDistance
	}

	var currentMarketplace []*stapi.MarketplaceGood
	if mk, found := analysis.GetMarketplaceAtLocation(s.ship.Location); !found {

		// this shouldn't really happen, but just in case...
		var err error
		mk, err = stapi.GetMarketplaceAtLocation(s.ship.Location)
		if err != nil {
			return err
		}

		currentMarketplace = mk

	} else {
		currentMarketplace = mk
	}

	// fuel time!
	if err = s.planFuel(fp, currentLocation, currentMarketplace); err != nil {
		return err
	}

	s.log(
		"flightplan created\nCost: %dcr\nExtra fuel: %d units\nDestination: %s (%s)\nDistance: %d",
		fp.flightCost,
		fp.extraFuelRequired,
		fp.destination.Name,
		fp.destination.Symbol,
		fp.distance,
	)

	s.log("sleeping for 5 seconds...")
	time.Sleep(time.Second * 5)

	// let's fly!
	if err = s.doFlight(fp); err != nil {
		return err
	}

	return nil
}

func (s *ShipController) grabMarketplaceData() error {
	marketplace, err := stapi.GetMarketplaceAtLocation(s.ship.Location)
	if err != nil {
		return err
	}
	return analysis.RecordMarketplaceAtLocation(s.ship.Location, marketplace)
}