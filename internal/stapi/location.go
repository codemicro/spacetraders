package stapi

import (
	"encoding/json"
	"errors"
	"strings"
)

// {
//   "symbol": "OE-PM",
//   "type": "PLANET",
//   "name": "Prime",
//   "x": 20,
//   "y": -25,
//   "ships": []
// }

type Location struct {
	Symbol         string             `json:"symbol"`
	Type           LocationType       `json:"type"`
	Name           string             `json:"name"`
	XCoordinate    int                `json:"x"`
	YCoordinate    int                `json:"y"`
	Ships          []*Ship            `json:"ships"`
	AvailableGoods []*MarketplaceGood `json:"marketplace"`
}

type LocationType uint

func (l *LocationType) UnmarshalJSON(data []byte) error {

	var ds string
	if err := json.Unmarshal(data, &ds); err != nil {
		return err
	}

	ds = strings.ToLower(ds)

	switch ds {
	case "planet":
		*l = LocationTypePlanet
	case "moon":
		*l = LocationTypeMoon
	case "wormhole":
		*l = LocationTypeWormhole
	case "gas_giant":
		*l = LocationTypeGasGiant
	case "asteroid":
		*l = LocationTypeAsteroid
	}
	return nil
}

const (
	LocationTypeUnknown = iota
	LocationTypePlanet
	LocationTypeMoon
	LocationTypeWormhole
	LocationTypeAsteroid
	LocationTypeGasGiant
)

var (
	ErrorSystemNotFound   = errors.New("stapi: system not found")
	ErrorLocationNotFound = errors.New("stapi: location is not detectable or does not exist")
)

func GetSystemLocations(system string) ([]*Location, error) {
	url := URLSystemLocations(system)
	ts := struct {
		Locations []*Location `json:"locations"`
	}{}

	err := orchestrateRequest(
		request.Clone().Get(url),
		&ts,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorSystemNotFound},
	)

	return ts.Locations, err
}

func GetLocationInfo(location string) (*Location, error) {
	url := URLLocationInformation(location)
	ts := struct {
		Location *Location `json:"location"` // This is fumbling the dockedShips parameter of the response
	}{}

	err := orchestrateRequest(
		request.Clone().Get(url),
		&ts,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorLocationNotFound},
	)

	return ts.Location, err
}
