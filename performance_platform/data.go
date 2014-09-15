package performance_platform

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/alphagov/metadata-api/request"
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

func ParseBackdropResponse(response []byte) (*Backdrop, error) {
	backdropResponse := &Backdrop{}
	if err := json.Unmarshal(response, &backdropResponse); err != nil {
		return nil, err
	}

	return backdropResponse, nil
}

func FetchSlugStatistics(performanceAPI, slug string, log *logrus.Logger) (*Backdrop, error) {
	escapedSlug := url.QueryEscape(slug)
	statisticsURL := performanceAPI + pageStatisticsURL + "?filter_by=pagePath:" + escapedSlug

	log.WithFields(logrus.Fields{
		"escapedSlug":   escapedSlug,
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
