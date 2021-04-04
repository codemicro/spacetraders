package analysis

import (
	"github.com/codemicro/spacetraders/internal/stapi"
	"sort"
)

type RoutingMethod uint8

const (
	RoutingMethodShortest = iota
	RoutingMethodShort
	RoutingMethodMedium
	RoutingMethodLong
	RoutingMethodLongest
)

func PickRoute(from *stapi.Location, options []*stapi.Location, method RoutingMethod) *stapi.Location {

	// remove from from options if exists
	for i, opt := range options {
		if opt == nil { // this will occur if anything has been deleted, since it is replaced with nil
			continue
		}
		if from.Symbol == opt.Symbol {
			options[i] = options[len(options)-1]
			options[len(options)-1] = nil
			options = options[:len(options)-1]
		}
	}

	distances := make(map[string]int)

	for _, routeOption := range options {
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
	case RoutingMethodShort:
		targetLocation = options[len(options) / 4]
	case RoutingMethodMedium:
		targetLocation = options[len(options) / 2]
	case RoutingMethodLong:
		targetLocation = options[(len(options) / 4) * 3]
	case RoutingMethodLongest:
		targetLocation = options[len(options) - 1]
	}

	return targetLocation
}