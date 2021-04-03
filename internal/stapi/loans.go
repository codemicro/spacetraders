package stapi

import (
	"time"
)

// {
//   "id": "ckn0zuhvz7120191bs6o6h4w57i",
//   "due": "2021-04-05T00:24:45.070Z",
//   "repaymentAmount": 280000,
//   "status": "CURRENT",
//   "type": "STARTUP"
// }

type Loan struct {
	// TODO: can status and type be represented by iotas?

	ID              string     `json:"id"`
	Due             *time.Time `json:"due"`
	RepaymentAmount int        `json:"repaymentAmount"`
	Status          string     `json:"status"`
	Type            string     `json:"type"`
}
