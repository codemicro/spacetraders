package analysis

import (
	"github.com/codemicro/spacetraders/internal/stapi"
	"sort"
)

type RoutingMethod uint8

const (
	RoutingMethodShortest = iota
	RoutingMethodMedium
	RoutingMethodLong
)

func PickRoute(from *stapi.Location, options []*stapi.Location, method RoutingMethod) *stapi.Location {
	distances := make(map[string]int64)

	for _, routeOption := range options {

		if from.Symbol == routeOption.Symbol {
			continue
		}

		distances[routeOption.Symbol] = FindDistance(from, routeOption)

	}

	sort.Slice(options, func(i, j int) bool {
		iSym := options[i].Symbol
		jSym := options[j].Symbol
		return distances[iSym] < distances[jSym]
	})

	var targetLocation *stapi.Location

	switch method {
	case RoutingMethodShortest:
		targetLocation = options[0]
	case RoutingMethodMedium:
		targetLocation = options[len(options) / 2] // will always tend towards the highest middle value
	case RoutingMethodLong:
		targetLocation = options[len(options) - 1]
	}

	return targetLocation
}