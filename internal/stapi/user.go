package stapi

import (
	"github.com/hashicorp/go-multierror"
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

func GetUserInfo(username string) (*User, error) {
	url := URLUserInfo(username)
	ts := struct{User User `json:"user"`}{}

	r := request.Clone().Get(url)

	if _, _, errs := r.EndStruct(&ts); errs != nil {
		return nil, multierror.Append(nil, errs...)
	}

	return &ts.User, nil
}