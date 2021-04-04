package analysis

import (
	"github.com/codemicro/spacetraders/internal/stapi"
	"sort"
	"strings"
)

type CargoMethod uint8

const (
	CargoMethodNone = iota
	CargoMethodCheapest
)

func PickCargo(options []*stapi.MarketplaceGood, method CargoMethod) *stapi.MarketplaceGood {

	// TODO: this thing

	// remove fuel from cargo options
	for i, opt := range options {
		if opt == nil { // this will occur if anything has been deleted, since it is replaced with nil
			continue
		}
		if strings.EqualFold(opt.Symbol, "FUEL") {
			options = append(options[:i], options[i+1:]...)
		}
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].PurchasePricePerUnit < options[j].PurchasePricePerUnit
	})

	var targetCargo *stapi.MarketplaceGood
	switch method {
	case CargoMethodCheapest:
		targetCargo = options[0]
	}

	return targetCargo
}
