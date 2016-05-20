package performanceclient

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-querystring/query"
)

// DataClient is a client for the Backdrop Read API
type DataClient interface {
	// BuildURL creates the URL to make a request to the read API
	BuildURL(dataGroup, dataType string, dataQuery QueryParams) string
	// Fetch makes a request to the read API and returns a non-nil BackdropResponse or an error
	Fetch(dataGroup, dataType string, dataQuery QueryParams) (*BackdropResponse, error)
}

type defaultDataClient struct {
	URL     string
	log     *logrus.Logger
	options []Option
}

// NewDataClient creates a new DataClient for talking to the Read API
func NewDataClient(url string, logger *logrus.Logger, options ...Option) DataClient {
	return &defaultDataClient{
		URL:     url,
		log:     logger,
		options: options,
	}
}

// BackdropResponse is a response from the Backdrop Read API
type BackdropResponse struct {
	Data    json.RawMessage `json:"data"`
	Warning string          `json:"warning,omitempty"`
	Status  string          `json:"status,omitempty"`
	Message string          `json:"message,omitempty"`
}

func (client *defaultDataClient) BuildURL(dataGroup, dataType string, dataQuery QueryParams) string {
	url := fmt.Sprintf("%s/data/%s/%s", client.URL, dataGroup, dataType)

	values, _ := query.Values(dataQuery)
	queryParameters := values.Encode()

	if len(queryParameters) > 1 {
		url += "?" + queryParameters
	}

	return url
}

func (client *defaultDataClient) Fetch(dataGroup, dataType string, dataQuery QueryParams) (*BackdropResponse, error) {
	url := client.BuildURL(dataGroup, dataType, dataQuery)

	client.log.WithFields(logrus.Fields{
		"url": url,
	}).Debug("Requesting performance data for slug")

	backdropResponse, err := NewRequest(url, client.options...)

	if err != nil {
		switch err {
		case ErrBadRequest:
			if body, readErr := ReadResponseBody(backdropResponse); readErr == nil {
				_, backdropErr := parseBackdropResponse(body)
				client.log.Errorf("Bad request to URL %q %v %q", url, err, backdropErr)
			} else {
				client.log.Errorf("Bad request to URL %q", url)
			}
		case ErrNotFound:
			client.log.Errorf("Not found: %q", url)
		}
		return nil, err
	}

	backdropBody, err := ReadResponseBody(backdropResponse)
	if err != nil {
		return nil, err
	}

	backdrop, err := parseBackdropResponse([]byte(backdropBody))
	if err != nil {
		return nil, err
	}

	return backdrop, nil
}

func parseBackdropResponse(response []byte) (*BackdropResponse, error) {
	backdropResponse := &BackdropResponse{}
	if err := json.Unmarshal(response, &backdropResponse); err != nil {
		return nil, err
	}

	if backdropResponse.Status == "error" {
		return nil, errors.New(backdropResponse.Message)
	}

	return backdropResponse, nil
}
