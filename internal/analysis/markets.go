package analysis

import (
	"encoding/json"
	"github.com/codemicro/spacetraders/internal/db"
	"github.com/codemicro/spacetraders/internal/stapi"
)

// This file is for things related to tracking the state of markets in systems

type Markets map[string][]*stapi.MarketplaceGood

func RecordMarketplaceAtLocation(location string, marketplace []*stapi.MarketplaceGood) error {

	marketplaceData, err := json.Marshal(marketplace)
	if err != nil {
		return err
	}

	return db.RecordMarketData(location, string(marketplaceData))
}

func GetMarketplaceAtLocation(location string) ([]*stapi.MarketplaceGood, bool, error) {

	dat, found, err := db.GetLatestDataForLocation(location)
	if !found || err != nil {
		return nil, found, err
	}

	var goods []*stapi.MarketplaceGood
	err = json.Unmarshal([]byte(dat.Data), &goods)
	if err != nil {
		return nil, true, err
	}

	return goods, true, nil
}

func GetAllMarketplaces() (Markets, error) {

	marketLocations, err := db.GetMarketLocations()
	if err != nil {
		return nil, err
	}

	markets := make(Markets)
	for _, marketLocation := range marketLocations {
		data, found, err := db.GetLatestDataForLocation(marketLocation)
		if !found || err != nil {
			return nil, err
		}

		var x []*stapi.MarketplaceGood
		err = json.Unmarshal([]byte(data.Data), &x)
		if err != nil {
			return nil, err
		}
		markets[marketLocation] = x
	}

	return markets, nil
}
