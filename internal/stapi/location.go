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
	Symbol      string       `json:"symbol"`
	Type        LocationType `json:"type"`
	Name        string       `json:"name"`
	XCoordinate int          `json:"x"`
	YCoordinate int          `json:"y"`
	Ships       []*Ship      `json:"ships"`
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
	}
	return nil
}

const (
	LocationTypeUnknown = iota
	LocationTypePlanet
	LocationTypeMoon
	LocationTypeWormhole
)

var (
	ErrorSystemNotFound = errors.New("stapi: system not found")
)

func GetSystemLocations(system string) ([]*Location, error) {
	url := URLSystemLocations(system)
	ts := struct{Locations []*Location `json:"locations"`}{}

	err := orchestrateRequest(
		request.Clone().Get(url),
		&ts,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorSystemNotFound},
	)

	return ts.Locations, err
}
