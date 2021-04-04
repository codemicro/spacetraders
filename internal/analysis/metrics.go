package analysis

import (
	"github.com/codemicro/spacetraders/internal/stapi"
	"math"
)

func FindDistance(from, to *stapi.Location) int {
	x := float64(to.XCoordinate - from.XCoordinate)
	y := float64(to.YCoordinate - from.YCoordinate)
	return int(math.Round(math.Sqrt(x*x + y*y)))
}

const fuelContingency = 2

func CalculateFuelForFlight(from, to *stapi.Location) int {

	dist := FindDistance(from, to)

	var rawFuelRequired int
	{
		d := int(math.Round(float64(dist) / 4))
		e := 1
		if from.Type == stapi.LocationTypePlanet {
			e += 2
		}
		rawFuelRequired = d + e
	}

	return rawFuelRequired + fuelContingency
}
