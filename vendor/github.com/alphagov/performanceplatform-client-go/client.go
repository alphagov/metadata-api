package performanceclient

import (
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"
)

// Dashboards is a collection of Dashboard instances
type Dashboards struct {
	Items []Dashboard `json:"Items"`
}

// Dashboard represents a collection of modules for a given organisation
type Dashboard struct {
	Department    Organisation `json:"department"`
	Agency        Organisation `json:"agency,omitempty"`
	DashboardType string       `json:"dashboard-type"`
	Slug          string       `json:"slug"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Modules       []Module     `json:"modules"`
	Published     bool         `json:"published"`
	PageType      string       `json:"page-type"`
	Costs         string       `json:"costs"`
}

// Organisation represents a government organisational unit, such as department or agency
type Organisation struct {
	Abbreviation string `json:"abbr"`
	Title        string `json:"title"`
}

// Module represents a visualisation with a Dashboard
type Module struct {
	Info       []string   `json:"info"`
	DataSource DataSource `json:"data-source"`
	Tabs       []Tab      `json:"tabs"`
	Title      string     `json:"title"`
}

// Tab is a UI construct exposed within the meta data API
type Tab struct {
	Description string     `json:"description"`
	DataSource  DataSource `json:"data-source"`
}

// DataSource is the data source for a module
type DataSource struct {
	DataGroup   string      `json:"data-group"`
	DataType    string      `json:"data-type"`
	QueryParams QueryParams `json:"query-params"`
}

// QueryParams represents the possible parameters that the Backdrop RPC read API supports
type QueryParams struct {
	FilterBy       []string  `json:"filter_by,omitempty" url:"filter_by,omitempty"`
	FilterByPrefix []string  `json:"filter_by_prefix,omitempty" url:"filter_by_prefix,omitempty"`
	GroupBy        []string  `json:"group_by,omitempty" url:"group_by,omitempty"`
	Collect        []string  `json:"collect,omitempty" url:"collect,omitempty"`
	SortBy         string    `json:"sort_by,omitempty" url:"sort_by,omitempty"`
	Duration       int       `json:"duration,omitempty" url:"duration,omitempty"`
	Period         string    `json:"period,omitempty" url:"period,omitempty"`
	Limit          int       `json:"limit,omitempty" url:"limit,omitempty"`
	StartAt        time.Time `json:"start_at,omitempty" url:"start_at,omitempty"`
	EndAt          time.Time `json:"end_at,omitempty" url:"end_at,omitempty"`
}

// MetaClient defines the interface that we need to talk to the meta data API
type MetaClient interface {
	Fetch(slug string) (Dashboard, error)
	FetchDashboards() (Dashboards, error)
}

type defaultMetaClient struct {
	baseURL string
	log     *logrus.Logger
}

// NewMetaClient returns a new MetaClient implementation with sensible defaults.
func NewMetaClient(baseURL string, log *logrus.Logger) MetaClient {
	return &defaultMetaClient{baseURL, log}
}

func (c *defaultMetaClient) Fetch(slug string) (dashboard Dashboard, err error) {
	url := c.baseURL + "?slug=" + slug

	c.log.WithFields(logrus.Fields{
		"url": url,
	}).Debug("Requesting meta data for slug")

	resp, err := NewRequest(url)

	if err != nil {
		switch err {
		case ErrBadRequest:
			if body, readErr := ReadResponseBody(resp); readErr == nil {
				c.log.Errorf("Bad request to URL %q with result %q", url, body)
			} else {
				c.log.Errorf("Bad request to URL %q", url)
			}
		case ErrNotFound:
			c.log.Errorf("Not found: %q", url)
		}
		return
	}

	body, err := ReadResponseBody(resp)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &dashboard)
	return
}

func (c *defaultMetaClient) FetchDashboards() (results Dashboards, err error) {
	url := c.baseURL
	resp, err := NewRequest(url)

	if err != nil {
		return
	}

	body, err := ReadResponseBody(resp)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &results)
	return
}
