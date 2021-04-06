package stapi

import (
	"errors"
	"sync"
	"time"
)

// {
//   "user": {
//     "username": "akpa",
//     "credits": 106180,
//     "ships": [],
//     "loans": []
//   }
// }

type User struct {
	Username string  `json:"username"`
	Credits  int     `json:"credits"`
	Ships    []*Ship `json:"ships"`
	Loans    []*Loan `json:"loans"`

	updateLock sync.RWMutex
}

var (
	ErrorUserNotFound = errors.New("stapi: user not found")
)

// TODO: the errors returned by these functions are not adequate

func GetUserInfo(username string) (*User, error) {
	url := URLUserInfo(username)
	ts := struct {
		User *User `json:"user"`
	}{}

	err := orchestrateRequest(
		request.Clone().Get(url),
		&ts,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorUserNotFound},
		cachePolicy{ true, time.Minute * 10 },
	)

	return ts.User, err
}

func (u *User) coreTradeOrder(url string, output interface{}, shipID, good string, quantity int) error {
	requestMap := map[string]interface{}{
		"shipId":   shipID,
		"good":     good,
		"quantity": quantity,
	}

	return orchestrateRequest(
		request.Clone().Post(url).Type("form").SendMap(requestMap),
		&output,
		func(i int) bool { return i == 201 },
		map[int]error{404: ErrorUserNotFound},
		cachePolicy{ false, 0 },
	)
}

func (u *User) SubmitPurchaseOrder(shipID, good string, quantity int) (*Ship, error) {
	url := URLSubmitPurchaseOrder(u.Username)
	ts := struct {
		// This is fumbling the "order" parameter of the response
		Ship    *Ship `json:"ship"`
		Credits int   `json:"credits"`
	}{}

	err := u.coreTradeOrder(url, &ts, shipID, good, quantity)
	if err != nil {
		return nil, err
	}

	u.updateLock.Lock()
	u.Credits = ts.Credits
	u.updateLock.Unlock()

	return ts.Ship, nil
}

func (u *User) SubmitSellOrder(shipID, good string, quantity int) (*Ship, *Order, error) {
	url := URLSubmitSellOrder(u.Username)
	ts := struct {
		Ship    *Ship  `json:"ship"`
		Credits int    `json:"credits"`
		Order   *Order `json:"order"`
	}{}

	err := u.coreTradeOrder(url, &ts, shipID, good, quantity)
	if err != nil {
		return nil, nil, err
	}

	u.updateLock.Lock()
	u.Credits = ts.Credits
	u.updateLock.Unlock()

	return ts.Ship, ts.Order, nil
}

func (u *User) SubmitFlightplan(shipID, destination string) (*Flightplan, error) {
	url := URLSubmitFlightplan(u.Username)
	ts := struct {
		Flightplan *Flightplan `json:"flightPlan"`
	}{}

	requestMap := map[string]interface{}{
		"shipId":      shipID,
		"destination": destination,
	}

	err := orchestrateRequest(
		request.Clone().Post(url).Type("form").SendMap(requestMap),
		&ts,
		func(i int) bool { return i == 201 },
		map[int]error{404: ErrorUserNotFound},
		cachePolicy{ false, 0 },
	)

	if err != nil {
		return nil, err
	}

	return ts.Flightplan, nil
}

func (u *User) GetFlightplan(flightplanID string) (*Flightplan, error) {
	url := URLGetFlightplanInformation(u.Username, flightplanID)
	ts := struct {
		Flightplan *Flightplan `json:"flightPlan"`
	}{}

	err := orchestrateRequest(
		request.Clone().Get(url),
		&ts,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorUserNotFound},
		cachePolicy{ false, 0 },
	)

	if err != nil {
		return nil, err
	}

	return ts.Flightplan, nil
}

func (u *User) GetShipInfo(shipID string) (*Ship, error) {
	url := URLGetShipInfo(u.Username, shipID)
	ts := struct {
		Ship *Ship `json:"ship"`
	}{}

	err := orchestrateRequest(
		request.Clone().Get(url),
		&ts,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorUserNotFound},
		cachePolicy{ false, 0 },
	)

	if err != nil {
		return nil, err
	}

	return ts.Ship, nil
}
