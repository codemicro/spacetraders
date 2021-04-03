package stapi

import (
	"errors"
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
	Username string   `json:"username"`
	Credits  int      `json:"credits"`
	Ships    []*Ship  `json:"ships"`
	Loans    []*Loan `json:"loans"`
}

var (
	ErrorUserNotFound = errors.New("stapi: user not found")
)

func GetUserInfo(username string) (*User, error) {
	url := URLUserInfo(username)
	ts := struct{User *User `json:"user"`}{}

	err := orchestrateRequest(
		request.Clone().Get(url),
		&ts,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorUserNotFound},
	)

	return ts.User, err
}