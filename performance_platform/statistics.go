package performance_platform

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/jinzhu/now"
)

type Statistics struct {
	PageViews      []Statistic `json:"page_views"`
	Searches       []Statistic `json:"searches"`
	ProblemReports []Statistic `json:"problem_reports"`
	SearchTerms    SearchTerms `json:"search_terms"`
}

type SearchTerms []SearchTerm

type SearchTerm struct {
	Keyword       string
	TotalSearches int
	Searches      []Statistic
}

type Statistic struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int       `json:"value"`
}

func (terms SearchTerms) Len() int           { return len(terms) }
func (terms SearchTerms) Swap(i, j int)      { terms[i], terms[j] = terms[j], terms[i] }
func (terms SearchTerms) Less(i, j int) bool { return terms[i].TotalSearches > terms[j].TotalSearches }

func (client *Client) SlugStatistics(slug string) (*Statistics, error) {
	var pageViews, searches, problemReports []Statistic
	var searchTerms SearchTerms
	var waitGroup sync.WaitGroup

	errorChannel := make(chan error)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if pageViewsResponse, err := client.Fetch("govuk-info", "page-statistics", Query{
			FilterBy: []string{"pagePath:" + slug},
			Collect:  []string{"uniquePageviews:sum"},
			Duration: 42,
			Period:   "day",
			EndAt:    now.BeginningOfDay().UTC(),
		}); err != nil {
			errorChannel <- err
		} else {
			if pageViews, err = parsePageViews(pageViewsResponse); err != nil {
				errorChannel <- err
			}
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if searchesResponse, err := client.Fetch("govuk-info", "search-terms", Query{
			FilterBy: []string{"pagePath:" + slug},
			Collect:  []string{"searchUniques:sum"},
			Duration: 42,
			Period:   "day",
			EndAt:    now.BeginningOfDay().UTC(),
		}); err != nil {
			errorChannel <- err
		} else {
			if searches, err = parseSearches(searchesResponse); err != nil {
				errorChannel <- err
			}
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if searchTermsResponse, err := client.Fetch("govuk-info", "search-terms", Query{
			FilterBy: []string{"pagePath:" + slug},
			GroupBy:  []string{"searchKeyword"},
			Collect:  []string{"searchUniques:sum"},
			Duration: 42,
			Period:   "day",
			EndAt:    now.BeginningOfDay().UTC(),
		}); err != nil {
			errorChannel <- err
		} else {
			if searchTerms, err = parseSearchTerms(searchTermsResponse); err != nil {
				errorChannel <- err
			} else {
				sort.Sort(searchTerms)
				if len(searchTerms) > 10 {
					searchTerms = searchTerms[0:10]
				}
			}
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		if problemReportsResponse, err := client.Fetch("govuk-info", "page-contacts", Query{
			FilterBy: []string{"pagePath:" + slug},
			Collect:  []string{"total:sum"},
			Duration: 42,
			Period:   "day",
			EndAt:    now.BeginningOfDay().UTC(),
		}); err != nil {
			errorChannel <- err
		} else {
			if problemReports, err = parseProblemReports(problemReportsResponse); err != nil {
				errorChannel <- err
			}
		}
	}()

	waitGroup.Wait()

	if len(errorChannel) > 0 {
		return nil, <-errorChannel
	}

	return &Statistics{
		PageViews:      pageViews,
		Searches:       searches,
		ProblemReports: problemReports,
		SearchTerms:    searchTerms,
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

func parseProblemReports(response *BackdropResponse) ([]Statistic, error) {
	var data []struct {
		Timestamp      time.Time `json:"_start_at"`
		ProblemReports float32   `json:"total:sum"`
	}

	err := json.Unmarshal(response.Data, &data)

	if err != nil {
		return []Statistic{}, err
	}

	statistics := make([]Statistic, len(data))
	for i, datum := range data {
		statistics[i] = Statistic{
			Timestamp: datum.Timestamp,
			Value:     int(datum.ProblemReports),
		}
	}
	return statistics, nil
}

func parseSearchTerms(response *BackdropResponse) (SearchTerms, error) {
	var data []struct {
		Keyword       string  `json:"searchKeyword"`
		TotalSearches float32 `json:"searchUniques:sum"`

		Values []struct {
			Timestamp     time.Time `json:"_start_at"`
			SearchUniques float32   `json:"searchUniques:sum"`
		} `json:"values"`
	}

	err := json.Unmarshal(response.Data, &data)

	if err != nil {
		return SearchTerms{}, err
	}

	terms := make(SearchTerms, len(data))
	for i, datum := range data {
		statistics := make([]Statistic, len(datum.Values))
		for j, value := range datum.Values {
			statistics[j] = Statistic{
				Timestamp: value.Timestamp,
				Value:     int(value.SearchUniques),
			}
		}
		terms[i] = SearchTerm{
			Keyword:       datum.Keyword,
			TotalSearches: int(datum.TotalSearches),
			Searches:      statistics,
		}
	}
	return terms, nil
}
