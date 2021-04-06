package analysis

import (
	"errors"
	"github.com/codemicro/spacetraders/internal/stapi"
	"github.com/codemicro/spacetraders/internal/tool"
	"sort"
	"strings"
)

type cargoDestination struct {
	Cargo string
	Destination string
	Value float64
}

func FindCombinedRouteAndCargo(currentLocationSymbol string) (*stapi.Location, *stapi.MarketplaceGood, error) {

	systemLocations, err := stapi.GetSystemLocations(tool.SystemFromSymbol(currentLocationSymbol))
	if err != nil {
		return nil, nil, err
	}

	var currentLocation *stapi.Location
	{
		// filter out current location and any wormholes
		var n int
		for _, location := range systemLocations {
			isCurrent := strings.EqualFold(location.Symbol, currentLocationSymbol)
			if !isCurrent && location.Type != stapi.LocationTypeWormhole {
				systemLocations[n] = location
				n++
			} else if isCurrent {
				currentLocation = location
			}
		}
		systemLocations = systemLocations[:n]
	}

	distancesTo := make(map[string]int)
	for _, location := range systemLocations {
		distancesTo[location.Symbol] = FindDistance(currentLocation, location)
	}

	marketplaces, err := GetAllMarketplaces()
	if err != nil {
		return nil, nil, err
	}

	var rankings []cargoDestination

	var currentLocationGoods []*stapi.MarketplaceGood
	for key, val := range marketplaces {
		if strings.EqualFold(currentLocationSymbol, key) {
				currentLocationGoods = val
		}
	}

	for _, currentLocationGood := range currentLocationGoods {
		// for every other market
		for marketLocation, market := range marketplaces {
			for _, marketGood := range market {
				if strings.EqualFold(currentLocationGood.Symbol, marketGood.Symbol) {
					// the destination also has this type of good

					if currentLocationGood.PurchasePricePerUnit > marketGood.SellPricePerUnit {
						// if we're going to lose money, we're not interested
						continue
					}

					rankings = append(rankings, cargoDestination{
						Cargo:       marketGood.Symbol,
						Destination: marketLocation,
						Value:       float64(currentLocationGood.PurchasePricePerUnit - marketGood.SellPricePerUnit) / float64(distancesTo[marketLocation]*marketGood.VolumePerUnit),
					})

				}
			}
		}

	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Value > rankings[j].Value
	})

	if len(rankings) < 1 {
		return nil, nil, errors.New("analysis: no suitable routes")
	}

	selected := rankings[0]

	var selectedLocation *stapi.Location
	for _, loc := range systemLocations {
		if strings.EqualFold(loc.Symbol, selected.Destination) {
			selectedLocation = loc
			break
		}
	}

	var selectedGood *stapi.MarketplaceGood
	for _, good := range currentLocationGoods {
		if strings.EqualFold(selected.Cargo, good.Symbol) {
			selectedGood = good
			break
		}
	}

	return selectedLocation, selectedGood, nil

}
