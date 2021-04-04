package analysis

import (
	"github.com/codemicro/spacetraders/internal/stapi"
	"strings"
)

func PickCargo(options []*stapi.MarketplaceGood) *stapi.MarketplaceGood {

	// TODO: this thing

	newOptions := make([]*stapi.MarketplaceGood, len(options))
	copy(newOptions, options)

	// remove fuel from cargo options
	for i, opt := range newOptions {
		if opt == nil { // this will occur if anything has been deleted, since it is replaced with nil
			continue
		}
		if strings.EqualFold(opt.Symbol, "FUEL") {
			newOptions[i] = newOptions[len(newOptions)-1]
			newOptions[len(newOptions)-1] = nil
			newOptions = newOptions[:len(newOptions)-1]
		}
	}

	return nil
}
