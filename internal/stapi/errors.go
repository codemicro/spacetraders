package stapi

import (
	"encoding/json"
	"fmt"
)

// {"error":{"message":"User has insufficient credits for transaction.","code":2004}}

type ErrorResponse struct {
	Message string `json:"message"`
	Code int `json:"code"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("error response from API: code %d, %s", e.Code, e.Message)
}

// Known error codes
const (
	ErrorCodeInsufficientFunds = 2004
	ErrorCodeNotFound = 404
)

func newAPIError(httpStatusCode int, responseBody []byte) error {
	ts := struct {ef *ErrorResponse `json:"error"`}{}
	if err := json.Unmarshal(responseBody, &ts); err != nil {
		// TODO: should this return the JSON error or this other error?
		return err
		//return &ErrorResponse{
		//	Message: string(responseBody),
		//	Code:    httpStatusCode,
		//}
	}
	return ts.ef
}