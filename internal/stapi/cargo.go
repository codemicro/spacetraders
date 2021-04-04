package stapi

// {
//   "good": "FUEL",
//   "quantity": 18,
//   "totalVolume": 18
// }

type Cargo struct {
	Good        string `json:"good"`
	Quantity    int    `json:"quantity"`
	TotalVolume int    `json:"totalVolume"`
}
