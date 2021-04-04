package stapi

import "errors"

// {
//   "symbol": "SHIP_PARTS",
//   "volumePerUnit": 4,
//   "pricePerUnit": 1006,
//   "spread": 3,
//   "purchasePricePerUnit": 1009,
//   "sellPricePerUnit": 1003,
//   "quantityAvailable": 43577
// }

// found at game/locations/OE-PM/marketplace

type MarketplaceGood struct {
	Symbol               string `json:"symbol"`
	VolumePerUnit        int    `json:"volumePerUnit"`
	PricePerUnit         int    `json:"pricePerUnit"`
	Spread               int    `json:"spread"`
	PurchasePricePerUnit int    `json:"purchasePricePerUnit"`
	SellPricePerUnit     int    `json:"sellPricePerUnit"`
	QuantityAvailable    int    `json:"quantityAvailable"`
}

var (
	ErrorCannotViewMarketplace = errors.New("stapi: marketplace listings are only visible to docked ships at this location")
)

func GetMarketplaceAtLocation(location string) ([]*MarketplaceGood, error) {
	url := URLMarketplaceAtLocation(location)
	ts := struct {
		Location *Location `json:"location"`
	}{}

	err := orchestrateRequest(
		request.Clone().Get(url),
		&ts,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorLocationNotFound, 400: ErrorCannotViewMarketplace},
	)

	if err != nil {
		return nil, err
	}
	return ts.Location.AvailableGoods, nil
}