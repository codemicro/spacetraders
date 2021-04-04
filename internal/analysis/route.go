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

	newOptions := make([]*stapi.Location, len(options))
	copy(newOptions, options)

	// remove from from options if exists
	for i, opt := range newOptions {
		if opt == nil { // this will occur if anything has been deleted, since it is replaced with nil
			continue
		}
		if from.Symbol == opt.Symbol {
			newOptions[i] = newOptions[len(newOptions)-1]
			newOptions[len(newOptions)-1] = nil
			newOptions = newOptions[:len(newOptions)-1]
		}
	}

	distances := make(map[string]int)

	for _, routeOption := range newOptions {
		distances[routeOption.Symbol] = FindDistance(from, routeOption)
	}

	sort.Slice(newOptions, func(i, j int) bool {
		iSym := newOptions[i].Symbol
		jSym := newOptions[j].Symbol
		return distances[iSym] < distances[jSym]
	})

	var targetLocation *stapi.Location

	switch method {
	case RoutingMethodShortest:
		targetLocation = newOptions[0]
	case RoutingMethodShort:
		targetLocation = newOptions[len(newOptions)/4]
	case RoutingMethodMedium:
		targetLocation = newOptions[len(newOptions)/2]
	case RoutingMethodLong:
		targetLocation = newOptions[(len(newOptions)/4)*3]
	case RoutingMethodLongest:
		targetLocation = newOptions[len(newOptions)-1]
	}

	return targetLocation
}
