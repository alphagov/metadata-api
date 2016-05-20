package performance_platform

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/alphagov/performanceplatform-client-go"
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
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
	Value     int       `json:"value"`
}

type PageViewsForDate struct {
	Timestamp time.Time `json:"_start_at"`
	PageViews float32   `json:"uniquePageviews:sum"`
}

type ProblemReportsForDate struct {
	Timestamp      time.Time `json:"_start_at"`
	ProblemReports float32   `json:"total:sum"`
}

type SearchUniquesForDate struct {
	Timestamp     time.Time `json:"_start_at"`
	SearchUniques float32   `json:"searchUniques:sum"`
}

func (terms SearchTerms) Len() int           { return len(terms) }
func (terms SearchTerms) Swap(i, j int)      { terms[i], terms[j] = terms[j], terms[i] }
func (terms SearchTerms) Less(i, j int) bool { return terms[i].TotalSearches > terms[j].TotalSearches }

func SlugStatistics(client performanceclient.DataClient, slug string, is_multipart bool) (*Statistics, error) {
	var pageViews, searches, problemReports []Statistic
	var searchTerms SearchTerms
	var waitGroup sync.WaitGroup

	errorChannel := make(chan error)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		query_params := performanceclient.QueryParams{
			Collect:  []string{"uniquePageviews:sum"},
			GroupBy:  []string{"pagePath"},
			Duration: 42,
			Period:   "day",
			EndAt:    now.BeginningOfDay().UTC(),
		}
		if !is_multipart {
			query_params.FilterBy = []string{"pagePath:" + slug}
		} else {
			query_params.FilterByPrefix = []string{"pagePath:" + slug}
		}

		if pageViewsResponse, err := client.Fetch("govuk-info", "page-statistics", query_params); err != nil {
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

		query_params := performanceclient.QueryParams{
			Collect:  []string{"searchUniques:sum"},
			GroupBy:  []string{"pagePath"},
			Duration: 42,
			Period:   "day",
			EndAt:    now.BeginningOfDay().UTC(),
		}
		if !is_multipart {
			query_params.FilterBy = []string{"pagePath:" + slug}
		} else {
			query_params.FilterByPrefix = []string{"pagePath:" + slug}
		}

		if searchesResponse, err := client.Fetch("govuk-info", "search-terms", query_params); err != nil {
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

		if searchTermsResponse, err := client.Fetch("govuk-info", "search-terms", performanceclient.QueryParams{
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

		query_params := performanceclient.QueryParams{
			Collect:  []string{"total:sum"},
			GroupBy:  []string{"pagePath"},
			Duration: 42,
			Period:   "day",
			EndAt:    now.BeginningOfDay().UTC(),
		}
		if !is_multipart {
			query_params.FilterBy = []string{"pagePath:" + slug}
		} else {
			query_params.FilterByPrefix = []string{"pagePath:" + slug}
		}

		if problemReportsResponse, err := client.Fetch("govuk-info", "page-contacts", query_params); err != nil {
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

func parsePageViews(response *performanceclient.BackdropResponse) ([]Statistic, error) {
	var datasetsPerPath []struct {
		Path   string             `json:"pagePath"`
		Values []PageViewsForDate `json:"values"`
	}

	err := json.Unmarshal(response.Data, &datasetsPerPath)

	if err != nil {
		return []Statistic{}, err
	}

	statistics := make([]Statistic, 0)
	for _, datasetPerPath := range datasetsPerPath {
		for _, datum := range datasetPerPath.Values {
			statistic := Statistic{
				Path:      datasetPerPath.Path,
				Timestamp: datum.Timestamp,
				Value:     int(datum.PageViews),
			}
			statistics = append(statistics, statistic)
		}
	}
	return statistics, nil
}

func parseSearches(response *performanceclient.BackdropResponse) ([]Statistic, error) {
	var datasetsPerPath []struct {
		Path   string                 `json:"pagePath"`
		Values []SearchUniquesForDate `json:"values"`
	}

	err := json.Unmarshal(response.Data, &datasetsPerPath)

	if err != nil {
		return []Statistic{}, err
	}

	statistics := make([]Statistic, 0)
	for _, datasetPerPath := range datasetsPerPath {
		for _, datum := range datasetPerPath.Values {
			statistic := Statistic{
				Path:      datasetPerPath.Path,
				Timestamp: datum.Timestamp,
				Value:     int(datum.SearchUniques),
			}
			statistics = append(statistics, statistic)
		}
	}
	return statistics, nil
}

func parseProblemReports(response *performanceclient.BackdropResponse) ([]Statistic, error) {
	var datasetsPerPath []struct {
		Path   string                  `json:"pagePath"`
		Values []ProblemReportsForDate `json:"values"`
	}

	err := json.Unmarshal(response.Data, &datasetsPerPath)

	if err != nil {
		return []Statistic{}, err
	}

	statistics := make([]Statistic, 0)
	for _, datasetPerPath := range datasetsPerPath {
		for _, datum := range datasetPerPath.Values {
			statistic := Statistic{
				Path:      datasetPerPath.Path,
				Timestamp: datum.Timestamp,
				Value:     int(datum.ProblemReports),
			}
			statistics = append(statistics, statistic)
		}
	}
	return statistics, nil
}

func parseSearchTerms(response *performanceclient.BackdropResponse) (SearchTerms, error) {
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
