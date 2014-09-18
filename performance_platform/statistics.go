package performance_platform

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/now"
)

type Statistics struct {
	PageViews   []Statistic `json:"page_views"`
	Searches    []Statistic `json:"searches"`
	SearchTerms []string    `json:"search_terms"`
}

type Statistic struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int       `json:"value"`
}

func (client *Client) SlugStatistics(slug string) (*Statistics, error) {
	var pageViews, searches []Statistic

	if pageViewsResponse, err := client.Fetch("govuk-info", "page-statistics", Query{
		FilterBy: []string{"pagePath:" + slug},
		Collect:  []string{"uniquePageviews:sum"},
		Duration: 42,
		Period:   "day",
		EndAt:    now.BeginningOfDay().UTC(),
	}); err != nil {
		return nil, err
	} else {
		if pageViews, err = parsePageViews(pageViewsResponse); err != nil {
			return nil, err
		}
	}

	if searchesResponse, err := client.Fetch("govuk-info", "search-terms", Query{
		FilterBy: []string{"pagePath:" + slug},
		Collect:  []string{"searchUniques:sum"},
		Duration: 42,
		Period:   "day",
		EndAt:    now.BeginningOfDay().UTC(),
	}); err != nil {
		return nil, err
	} else {
		if searches, err = parseSearches(searchesResponse); err != nil {
			return nil, err
		}
	}

	return &Statistics{
		PageViews: pageViews,
		Searches:  searches,
	}, nil
}

func parsePageViews(response *BackdropResponse) ([]Statistic, error) {
	var data []struct {
		Timestamp time.Time `json:"_start_at"`
		PageViews float32   `json:"uniquePageviews:sum"`
	}

	err := json.Unmarshal(response.Data, &data)

	if err != nil {
		return []Statistic{}, err
	}

	statistics := make([]Statistic, len(data))
	for i, datum := range data {
		statistics[i] = Statistic{
			Timestamp: datum.Timestamp,
			Value:     int(datum.PageViews),
		}
	}
	return statistics, nil
}

func parseSearches(response *BackdropResponse) ([]Statistic, error) {
	var data []struct {
		Timestamp     time.Time `json:"_start_at"`
		SearchUniques float32   `json:"searchUniques:sum"`
	}

	err := json.Unmarshal(response.Data, &data)

	if err != nil {
		return []Statistic{}, err
	}

	statistics := make([]Statistic, len(data))
	for i, datum := range data {
		statistics[i] = Statistic{
			Timestamp: datum.Timestamp,
			Value:     int(datum.SearchUniques),
		}
	}
	return statistics, nil
}
