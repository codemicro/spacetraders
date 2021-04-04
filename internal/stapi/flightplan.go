package stapi

import "time"

// {
//   "id": "ckn10o6p79173851bs61c341aog",
//   "shipId": "ckn10c5mj8363091bs69fldut5a",
//   "createdAt": "2021-04-03T00:47:50.251Z",
//   "arrivesAt": "2021-04-03T00:49:30.236Z",
//   "destination": "OE-PM",
//   "departure": "OE-PM-TR",
//   "distance": 4,
//   "fuelConsumed": 2,
//   "fuelRemaining": 18,
//   "terminatedAt": null,
//   "timeRemainingInSeconds": 99
// }

type Flightplan struct {
	ID                  string     `json:"id"`
	ShipID              string     `json:"shipId"`
	CreatedAt           *time.Time `json:"createdAt"`
	ArrivesAt           *time.Time `json:"arrivesAt"`
	Destination         string     `json:"destination"`
	Departure           string     `json:"departure"`
	Distance            int        `json:"distance"`
	FuelConsumed        int        `json:"fuelConsumed"`
	FuelRemaining       int        `json:"fuelRemaining"`
	TerminatedAt        *time.Time `json:"terminatedAt"`
	FlightTimeRemaining int        `json:"timeRemainingInSeconds"`
}
