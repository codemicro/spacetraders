package stapi

import (
	"encoding/json"
	"errors"
	"github.com/codemicro/spacetraders/internal/config"
	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	coreLog "log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	requestDelay = time.Second
	numWorkers   = 2
)

var request = gorequest.New().
	Timeout(10*time.Second).
	Set("Authorization", "Bearer "+config.C.Token).
	SetDebug(config.C.DebugMode)

func init() {

	if config.C.DebugMode {
		f, err := os.OpenFile("request.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		request.SetLogger(coreLog.New(f, "", 0))
	}

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

type cachePolicy struct {
	Allow bool
	CacheDuration time.Duration
}

var (
	requestQueue = make(chan trackedRequest, 1024)
	responseCache = cache.New(time.Minute * 5, time.Minute * 5)

	requestsInProgressLock sync.RWMutex
	requestsInProgress = make(map[string]*sync.WaitGroup)
)

func orchestrateRequest(req *gorequest.SuperAgent, output interface{}, isStatusCodeOk func(int) bool, errorsByStatusCode map[int]error, cPolicy cachePolicy) error {

	allowCache := cPolicy.Allow && req.Method == "GET"

	var wg *sync.WaitGroup

	if allowCache {

		// Because of the nature of the queue system used in this program, multiple requests for the same cache-able resource can be in the queue at any one time.
		// To remedy this, a map of currently queued requests is held with a wait group.
		// If a request for an in-progress resource comes in, execution will be blocked until the request has been completed and the response cached.

		requestsInProgressLock.RLock()
		wg = requestsInProgress[req.Url]
		requestsInProgressLock.RUnlock()

		if wg != nil {
			wg.Wait()
		}

		if dat, found := responseCache.Get(req.Url); found {
			return json.Unmarshal(*dat.(*[]byte), &output)
		}

		requestsInProgressLock.Lock()
		wg = new(sync.WaitGroup)
		wg.Add(1)
		requestsInProgress[req.Url] = wg
		requestsInProgressLock.Unlock()

		defer func() {
			wg.Done()

			requestsInProgressLock.Lock()
			delete(requestsInProgress, req.Url)
			requestsInProgressLock.Unlock()
		}()

	}

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
		return newAPIError(completed.response.StatusCode, completed.body)
	}

	// at this point we can cache the response, since it's all ok
	if allowCache {
		responseCache.Set(req.Url, &completed.body, cPolicy.CacheDuration)
	}

	// parse response and return error or nil
	return json.Unmarshal(completed.body, &output)
}

var ErrorFailedRatelimit = errors.New("stapi: unable to make request (too many responses with 429)")

func requestWorker() {
	for rq := range requestQueue {

		var ratelimitRetries int
		var conflictRetries int
		var internalServerErrorRetries int

		for {
			resp, body, errs := rq.request.Clone().EndBytes()

			if resp.StatusCode == 429 {
				ratelimitRetries += 1

				// wait for the retry after and a random duration between 0 and 5 seconds extra
				retryAfter, _ := strconv.Atoi(resp.Header.Get("retry-after"))
				n := rand.Intn(5) + retryAfter

				time.Sleep(time.Duration(n) * time.Second)

				continue
			} else if resp.StatusCode == 409 {
				conflictRetries += 1

				if conflictRetries != 3 {
					log.Info().Msg("got 409 conflict on " + rq.request.Url)
					time.Sleep(time.Second * 2)
					continue
				} // otherwise return error as normal

			} else if resp.StatusCode == 500 {
				internalServerErrorRetries += 1

				if internalServerErrorRetries != 5 {
					log.Info().Msg("got 500 internal server error on " + rq.request.Url + " - will retry in 10 seconds")
					time.Sleep(time.Second * 10)
					continue
				} // otherwise return error as normal
			}

			if ratelimitRetries == 3 {
				rq.responseNotifier <- &completedRequest{
					response: resp,
					body:     body,
					err:      ErrorFailedRatelimit,
				}
				close(rq.responseNotifier)
				break
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
