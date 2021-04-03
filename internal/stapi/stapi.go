package stapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codemicro/spacetraders/internal/config"
	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	requestDelay = time.Second
	numWorkers = 2
)

var request = gorequest.New().Timeout(10*time.Second).AppendHeader("Authorization", "Bearer "+config.C.Token)

func init() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numWorkers; i += 1 {
		go requestWorker()
	}

}

// trackedRequest represents a pending request
type trackedRequest struct {
	responseNotifier chan *completedRequest
	request          *gorequest.SuperAgent
}

// completedRequest represents a response and corresponding error as a result of a HTTP request
type completedRequest struct {
	response *http.Response
	body     []byte
	err      error
}

type ApiError struct {
	StatusCode   int
	ResponseBody []byte
}

func (err *ApiError) Error() string {
	return fmt.Sprintf("stapi: the API returned a non-okay status code, %d", err.StatusCode)
}

func newApiError(statusCode int, responseBody []byte) *ApiError {
	return &ApiError{
		StatusCode:   statusCode,
		ResponseBody: responseBody,
	}
}

var requestQueue = make(chan trackedRequest, 1024)

func orchestrateRequest(req *gorequest.SuperAgent, output interface{}, isStatusCodeOk func(int) bool, errorsByStatusCode map[int]error) error {

	responseNotifier := make(chan *completedRequest)

	requestQueue <- trackedRequest{
		responseNotifier: responseNotifier,
		request:          req,
	}

	completed := <-responseNotifier

	if completed.err != nil {
		return completed.err
	}

	// check status code map
	for code, err := range errorsByStatusCode {
		if completed.response.StatusCode == code {
			return err
		}
	}

	// check status function
	if !isStatusCodeOk(completed.response.StatusCode) {
		return newApiError(completed.response.StatusCode, completed.body)
	}

	// parse response and return error or nil
	return json.Unmarshal(completed.body, &output)
}

var ErrorFailedRatelimit = errors.New("stapi: unable to make request (too many responses with 429)")

func requestWorker() {
	var retries int
	for {

		rq := <-requestQueue

		for {
			resp, body, errs := rq.request.Clone().EndBytes()

			if resp.StatusCode == 429 {
				retries += 1

				// wait for the retry after and a random duration between 0 and 5 seconds extra
				retryAfter, _ := strconv.Atoi(resp.Header.Get("retry-after"))
				n := rand.Intn(5) + retryAfter

				time.Sleep(time.Duration(n) * time.Second)

				continue
			}

			if retries == 3 {
				rq.responseNotifier <- &completedRequest{
					response: resp,
					body:     body,
					err:      ErrorFailedRatelimit,
				}
				continue
			}

			var err error
			if errs != nil {
				err = multierror.Append(err, errs...)
			}

			rq.responseNotifier <- &completedRequest{
				body:     body,
				response: resp,
				err:      err,
			}
			close(rq.responseNotifier)
			break
		}

		time.Sleep(requestDelay)
	}
}
