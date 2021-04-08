package analysis

import (
	"errors"
	"fmt"
	"github.com/codemicro/spacetraders/internal/stapi"
	"github.com/codemicro/spacetraders/internal/tool"
	"sort"
	"strings"
)

type cargoDestination struct {
	Cargo         string
	Destination   string
	Value         int
	NumberOfUnits int
}

var ErrorNoSuitableRoutes = errors.New("analysis: no suitable routes")

func FindCombinedRouteAndCargo(currentLocationSymbol string, cargoCapacity, spendLimit int) (*stapi.Location, *stapi.MarketplaceGood, int, int, error) {

	systemLocations, err := stapi.GetSystemLocations(tool.SystemFromSymbol(currentLocationSymbol))
	if err != nil {
		return nil, nil, 0, 0, err
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

	if currentLocation == nil {
		return nil, nil, 0, 0, errors.New("could not locate " + currentLocationSymbol)
	}

	distancesTo := make(map[string]int)
	for _, location := range systemLocations {
		distancesTo[location.Symbol] = FindDistance(currentLocation, location)
	}

	marketplaces, err := GetAllMarketplaces()
	if err != nil {
		return nil, nil, 0, 0, err
	}

	var rankings []cargoDestination

	var currentLocationGoods []*stapi.MarketplaceGood
	{
		curr, err := stapi.GetMarketplaceAtLocation(currentLocationSymbol)
		if err != nil {
			return nil, nil, 0, 0, err
		}
		currentLocationGoods = curr
	}

	var currentFuelCost int
	// if there is no fuel at the current place, it will ignore the minimum profit required
	for _, x := range currentLocationGoods {
		if strings.EqualFold(x.Symbol, "FUEL") {
			currentFuelCost = x.PurchasePricePerUnit
			break
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

					requiredFuel := CalculateFuelForDistance(distancesTo[marketLocation], currentLocation.Type)

					profitPerUnit := marketGood.SellPricePerUnit - currentLocationGood.PurchasePricePerUnit

					unitsToBuy := (cargoCapacity - requiredFuel) / marketGood.VolumePerUnit
					for {
						cost := currentLocationGood.PurchasePricePerUnit * unitsToBuy
						if cost > spendLimit {
							unitsToBuy -= 1
						} else {
							break
						}
					}

					totalProfit := profitPerUnit * unitsToBuy

					if totalProfit < requiredFuel*currentFuelCost {
						continue
					}

					rankings = append(rankings, cargoDestination{
						Cargo:         marketGood.Symbol,
						Destination:   marketLocation,
						Value:         totalProfit,
						NumberOfUnits: unitsToBuy,
					})

				}
			}
		}

	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Value > rankings[j].Value
	})

	if len(rankings) < 1 {
		return nil, nil, 0, 0, ErrorNoSuitableRoutes
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

	if selectedLocation == nil || selectedGood == nil {
		return nil, nil, 0, 0, ErrorNoSuitableRoutes
	}

	return selectedLocation, selectedGood, selected.NumberOfUnits, selected.Value, nil

}
