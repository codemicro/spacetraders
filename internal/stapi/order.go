package stapi

// {
//   "good": "METALS",
//   "quantity": 80,
//   "pricePerUnit": 8,
//   "total": 640
// }

type Order struct {
	Good         string `json:"good"`
	Quantity     int    `json:"quantity"`
	PricePerUnit int    `json:"pricePerUnit"`
	Total        int    `json:"total"`
}
