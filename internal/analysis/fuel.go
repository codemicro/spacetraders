package analysis

import (
	"github.com/codemicro/spacetraders/internal/stapi"
	"math"
)

func GetPlanetPenaltyByShip(ship string, cargoCapacity int) int {
	switch ship {
	case "GR-MK-II":
		return 3
	case "GR-MK-III":
		return 4
	}

	switch cargoCapacity {
	case 1000:
		return 3
	case 5000:
		return 4
	default:
		return 2
	}
}

const fuelContingency = 2

func CalculateFuelForFlight(from, to *stapi.Location, shipType string, cargoCapacity int) int {
	return CalculateFuelForDistance(FindDistance(from, to), from.Type, shipType, cargoCapacity)
}

func CalculateFuelForDistance(dist int, departureType stapi.LocationType, shipType string, cargoCapacity int) int {
	var rawFuelRequired int
	{
		d := int(math.Round(float64(dist) / 4))
		e := 1
		if departureType == stapi.LocationTypePlanet {
			e += GetPlanetPenaltyByShip(shipType, cargoCapacity)
		}
		rawFuelRequired = d + e
	}

	return rawFuelRequired + fuelContingency
}
