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
