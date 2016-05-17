package performanceclient

import (
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"io/ioutil"
	"net/http"
	"time"
)

// Option is a self-referential function used to configure a RequestOptions struct.
// See http://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html
type Option func(*http.Request, *RequestOptions) Option

// RequestOptions is the container for tweaking how NewRequest functions.
type RequestOptions struct {
	// MaxElapsedTime is the duration allowed to try to get a response from the origin server.
	MaxElapsedTime time.Duration
}

func (ro *RequestOptions) option(req *http.Request, opts ...Option) (previous Option) {
	for _, opt := range opts {
		previous = opt(req, ro)
	}
	return previous
}

var (
	// ErrNotFound is an error indicating that the server returned a 404.
	ErrNotFound = errors.New("not found")
	// ErrBadRequest is an error indicating that the client request had a problem.
	ErrBadRequest = errors.New("bad request")
)

// NewRequest tries to make a request to the URL, returning the http.Response if it was successful, or an error if there was a problem.
// Optional Option arguments can be passed to specify contextual behaviour for this request. See MaxElapsedTime.
func NewRequest(url string, options ...Option) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Performance-Platform-Client/1.0")

	requestOptions := RequestOptions{MaxElapsedTime: 5 * time.Second}
	requestOptions.option(request, options...)

	response, err := tryGet(request, requestOptions)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if response.StatusCode == http.StatusBadRequest {
		return response, ErrBadRequest
	}

	return response, err
}

// MaxElapsedTime specifies the maximum duration that we should use to retry requests to the origin server. The default value is 5 seconds.
func MaxElapsedTime(duration time.Duration) Option {
	return func(req *http.Request, ro *RequestOptions) Option {
		previous := ro.MaxElapsedTime
		ro.MaxElapsedTime = duration
		return MaxElapsedTime(previous)
	}
}

// BearerToken sets the Authorization header for a request
func BearerToken(bearerToken string) Option {
	return func(req *http.Request, ro *RequestOptions) Option {
		previous := req.Header.Get("Authorization")
		if bearerToken != "" {
			req.Header.Add("Authorization", "Bearer "+bearerToken)
		} else {
			req.Header.Del("Authorization")
		}
		return BearerToken(previous)
	}
}

// ReadResponseBody reads the response body stream and returns a byte array, or an error if there was a problem.
func ReadResponseBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func tryGet(req *http.Request, options RequestOptions) (res *http.Response, err error) {
	// Use a channel to communicate between the goroutines. We use a channel rather
	// than simple variable closure since that's how Go works :)
	c := make(chan *http.Response, 1)

	operation := func() error {
		response, httpError := http.DefaultClient.Do(req)
		if httpError != nil {
			return httpError
		}
		switch response.StatusCode {
		case 502, 503:
			// Oh dear, we'll retry that one
			return fmt.Errorf("Server unavailable")
		}

		// We're good, keep the returned response
		c <- response
		return nil
	}

	expo := backoff.NewExponentialBackOff()
	expo.MaxElapsedTime = options.MaxElapsedTime

	err = backoff.Retry(operation, expo)

	if err != nil {
		// Operation has failed, repeatedly got a problem or server unavailable
		return nil, err
	}

	// Got a good response, take it out of the channel
	res = <-c

	return res, err
}
