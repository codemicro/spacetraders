package analysis

import (
	"encoding/json"
	"github.com/codemicro/spacetraders/internal/stapi"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// This file is for things related to tracking the state of markets in systems

type Markets map[string]*MarketInfo

type MarketInfo struct {
	Symbol string
	Time   time.Time
	Goods  []stapi.MarketplaceGood
}

const marketTrackerFile = "markets.json"

var (
	currentMarketState = make(Markets)
	marketTrackerLock  = new(sync.RWMutex)
)

func init() {
	// load markets file if exists
	var fileExists bool
	{
		_, err := os.Stat(marketTrackerFile)
		fileExists = err == nil
	}

	if fileExists {
		rawData, err := ioutil.ReadFile(marketTrackerFile)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(rawData, &currentMarketState)
		if err != nil {
			panic(err)
		}
	}

}

func RecordMarketplaceAtLocation(location string, marketplace []*stapi.MarketplaceGood) error {

	newMarketplace := make([]stapi.MarketplaceGood, len(marketplace))
	for _, x := range marketplace {
		newMarketplace = append(newMarketplace, *x)
	}

	marketTrackerLock.Lock()
	defer marketTrackerLock.Unlock()

	currentMarketState[location] = &MarketInfo{
		Symbol: location,
		Time:   time.Now(),
		Goods:  newMarketplace,
	}

	jsonData, err := json.MarshalIndent(currentMarketState, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(marketTrackerFile, jsonData, 0644)
}

func GetMarketplaceAtLocation(location string) ([]*stapi.MarketplaceGood, bool) {
	marketTrackerLock.RLock()

	curr, ok := currentMarketState[location]
	if !ok {
		return nil, false
	}

	marketTrackerLock.RUnlock()

	goods := make([]*stapi.MarketplaceGood, len(curr.Goods))
	for _, x := range curr.Goods{
		goods = append(goods, &x)
	}

	return goods, true
}
