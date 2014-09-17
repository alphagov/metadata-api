package performance_platform

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/alphagov/metadata-api/request"
	"github.com/google/go-querystring/query"
)

var (
	pageStatisticsURL = "/data/govuk-info/page-statistics"
)

type Data struct {
	HumanID          string  `json:"humanId"`
	PagePath         string  `json:"pagePath"`
	SearchKeyword    string  `json:"searchKeyword"`
	SearchUniques    float32 `json:"searchUniques"`
	SearchUniquesSum float32 `json:"searchUniques:sum"`
	TimeSpan         string  `json:"timeSpan"`
	Type             string  `json:"dataType"`
	UniquePageViews  float32 `json:"uniquePageViews"`

	// Underscore fields mean something in backdrop?
	ID        string    `json:"_id"`
	Count     float32   `json:"_count"`
	Timestamp time.Time `json:"_timestamp"`
}

type Backdrop struct {
	Data    []Data `json:"data"`
	Warning string `json:"warning"`
}

type Query struct {
	FilterBy []string `url:"filter_by,omitempty"`
	Collect  []string `url:"collect,omitempty"`
	SortBy   []string `url:"sort_by,omitempty"`

	Duration int    `url:"duration,omitempty"`
	Period   string `url:"period,omitempty"`

	StartAt time.Time `url:"start_at,omitempty"`
	EndAt   time.Time `url:"start_at,omitempty"`
}

func ParseBackdropResponse(response []byte) (*Backdrop, error) {
	backdropResponse := &Backdrop{}
	if err := json.Unmarshal(response, &backdropResponse); err != nil {
		return nil, err
	}

	return backdropResponse, nil
}

func BuildURL(base, dataGroup, dataType string, backdropQuery Query) string {
	path := fmt.Sprintf("/data/%s/%s", dataGroup, dataType)
	values, _ := query.Values(backdropQuery)
	queryParameters := values.Encode()

	url := base + path
	if len(queryParameters) > 1 {
		url += "?" + queryParameters
	}

	return url
}

func FetchSlugStatistics(performanceAPI, slug string, log *logrus.Logger) (*Backdrop, error) {
	query := Query{
		FilterBy: []string{"pagePath:" + slug},
	}
	statisticsURL := BuildURL(performanceAPI, "govuk-info", "page-statistics", query)

	log.WithFields(logrus.Fields{
		"escapedSlug":   slug,
		"statisticsURL": statisticsURL,
	}).Debug("Requesting performance data for slug")

	backdropResponse, err := request.NewRequest(statisticsURL, "EMPTY")
	if err != nil {
		return nil, err
	}

	backdropBody, err := request.ReadResponseBody(backdropResponse)
	if err != nil {
		return nil, err
	}

	backdrop, err := ParseBackdropResponse([]byte(backdropBody))
	if err != nil {
		return nil, err
	}

	return backdrop, nil
}
