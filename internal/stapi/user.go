package stapi

import (
	"errors"
	"sync"
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
	)

	return ts.User, err
}

func (u *User) SubmitPurchaseOrder(shipID, good string, quantity int) (*Ship, error) {
	url := URLSubmitPurchaseOrder(u.Username)
	ts := struct {
		// This is fumbling the "order" parameter of the response
		Ship *Ship `json:"ship"`
		Credits int `json:"credits"`
	}{}

	requestMap := map[string]interface{}{
		"shipId": shipID,
		"good": good,
		"quantity": quantity,
	}

	err := orchestrateRequest(
		request.Clone().Post(url).Type("form").SendMap(requestMap),
		&ts,
		func(i int) bool { return i == 201 },
		map[int]error{404: ErrorUserNotFound},
	)

	if err != nil {
		return nil, err
	}

	u.updateLock.Lock()
	u.Credits = ts.Credits
	u.updateLock.Unlock()

	return ts.Ship, nil
}