package stapi

// {
//   "id": "ckn10c5mj8363091bs69fldut5a",
//   "location": "OE-PM",
//   "x": 20,
//   "y": -25,
//   "cargo": [],
//   "spaceAvailable": 82,
//   "type": "JW-MK-I",
//   "class": "MK-I",
//   "maxCargo": 100,
//   "speed": 1,
//   "manufacturer": "Jackshaw",
//   "plating": 5,
//   "weapons": 5
// }

type Ship struct {
	ID             string   `json:"id"`
	Location       string   `json:"location"`
	XCoordinate    int      `json:"x"`
	YCoordinate    int      `json:"y"`
	Cargo          []*Cargo `json:"cargo"`
	SpaceAvailable int      `json:"spaceAvailable"`
	Type           string   `json:"type"`
	Class          string   `json:"class"`
	MaxCargo       int      `json:"maxCargo"`
	Speed          int      `json:"speed"`
	Manufacturer   string   `json:"manufacturer"`
	Plating        int      `json:"plating"`
	Weapons        int      `json:"weapons"`
}
