package performance_platform_test

import (
	"net/http"
	"time"

	. "github.com/alphagov/metadata-api/performance_platform"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Sirupsen/logrus"
	"github.com/alphagov/performanceplatform-client-go"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Statistics", func() {

	var server *ghttp.Server
	var client performanceclient.DataClient

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = performanceclient.NewDataClient(server.URL(), logrus.New())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("SlugStatistics", func() {
		It("Should return formatted data", func() {
			server.RouteToHandler("GET", "/data/govuk-info/page-statistics",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/data/govuk-info/page-statistics"),
					ghttp.RespondWith(http.StatusOK, `
{
"data": [
  {
    "pagePath": "/tax-disc",
    "values": [
      {
        "_count": 1,
        "_end_at": "2014-07-04T00:00:00+00:00",
        "_start_at": "2014-07-03T00:00:00+00:00",
        "uniquePageviews:sum": 25931
      }
    ]
  }
]
}`)))

			server.RouteToHandler("GET", "/data/govuk-info/search-terms",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/data/govuk-info/search-terms"),
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Query().Get("group_by") != "searchKeyword" {
							ghttp.RespondWith(http.StatusOK, `
{
"data": [
	{
		"pagePath": "/tax-disc",
		"values": [
			{
				"_count": 4,
				"_end_at": "2014-09-03T00:00:00+00:00",
				"_start_at": "2014-09-02T00:00:00+00:00",
				"searchUniques:sum": 71
			}
		]
	}
]
}`)(w, r)
						} else {
							ghttp.RespondWith(http.StatusOK, `
{
"data": [
  {
    "_count": 8,
    "_group_count": 8,
    "searchKeyword": "employer access",
    "searchUniques:sum": 126,
    "values": [{
      "_count": 1,
      "searchUniques:sum": 126,
      "_end_at": "2014-09-03T00:00:00+00:00",
      "_start_at": "2014-09-02T00:00:00+00:00"
    }]
  },
  {
    "_count": 3,
    "_group_count": 3,
    "searchKeyword": "pupil premium",
    "searchUniques:sum": 45,
    "values": [{
      "_count": 1,
      "searchUniques:sum": 45,
      "_end_at": "2014-09-03T00:00:00+00:00",
      "_start_at": "2014-09-02T00:00:00+00:00"
    }]
  },
  {
    "_count": 4,
    "_group_count": 4,
    "searchKeyword": "s2s",
    "searchUniques:sum": 104,
    "values": [{
      "_count": 1,
      "searchUniques:sum": 104,
      "_end_at": "2014-09-03T00:00:00+00:00",
      "_start_at": "2014-09-02T00:00:00+00:00"
    }]
  }
]
}`)(w, r)
						}
					})))

			server.RouteToHandler("GET", "/data/govuk-info/page-contacts",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/data/govuk-info/page-contacts"),
					ghttp.RespondWith(http.StatusOK, `
{
"data": [
	{
		"pagePath": "/tax-disc",
		"values": [
			{
				"_count": 4,
				"_end_at": "2014-09-03T00:00:00+00:00",
				"_start_at": "2014-09-02T00:00:00+00:00",
				"total:sum": 71
			}
		]
	}
]
}`)))

			statistics, err := SlugStatistics(client, "/foo", false)
			Expect(err).To(BeNil())
			Expect(statistics).ToNot(BeNil())
			Expect(len(statistics.PageViews)).To(Equal(1))
			Expect(len(statistics.Searches)).To(Equal(1))
			Expect(len(statistics.SearchTerms)).To(Equal(3))
			Expect(len(statistics.ProblemReports)).To(Equal(1))
			Expect(statistics.PageViews[0].Value).To(Equal(25931))
			Expect(statistics.PageViews[0].Path).To(Equal("/tax-disc"))
			Expect(statistics.Searches[0].Value).To(Equal(71))
			Expect(statistics.Searches[0].Path).To(Equal("/tax-disc"))
			Expect(statistics.ProblemReports[0].Value).To(Equal(71))
			Expect(statistics.ProblemReports[0].Path).To(Equal("/tax-disc"))

			pageViewTimestamp, err := time.Parse(time.RFC3339, "2014-07-03T00:00:00+00:00")
			Expect(err).To(BeNil())
			Expect(statistics.PageViews[0].Timestamp).
				To(Equal(pageViewTimestamp))

			searchesTimestamp, err := time.Parse(time.RFC3339, "2014-09-02T00:00:00+00:00")
			Expect(err).To(BeNil())
			Expect(statistics.Searches[0].Timestamp).
				To(Equal(searchesTimestamp))

			Expect(statistics.SearchTerms[0].Keyword).To(Equal("employer access"))
			Expect(statistics.SearchTerms[0].TotalSearches).To(Equal(126))
			Expect(statistics.SearchTerms[1].Keyword).To(Equal("s2s"))
			Expect(statistics.SearchTerms[1].TotalSearches).To(Equal(104))
			Expect(statistics.SearchTerms[2].Keyword).To(Equal("pupil premium"))
			Expect(statistics.SearchTerms[2].TotalSearches).To(Equal(45))

			Expect(len(statistics.SearchTerms[0].Searches)).To(Equal(1))
			Expect(statistics.SearchTerms[0].Searches[0].Value).To(Equal(126))
			Expect(statistics.SearchTerms[0].Searches[0].Timestamp).To(Equal(searchesTimestamp))
		})

	})

	Describe("SlugStatisticsMultiPartFormat", func() {
		It("Should return formatted data for a multi-part format", func() {
			server.RouteToHandler("GET", "/data/govuk-info/page-statistics",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/data/govuk-info/page-statistics"),
					ghttp.RespondWith(http.StatusOK, `
{
"data": [
  {
    "pagePath": "/tax-disc",
    "values": [
      {
        "_count": 1,
        "_end_at": "2014-07-03T00:00:00+00:00",
        "_start_at": "2014-07-02T00:00:00+00:00",
        "uniquePageviews:sum": 25931
      }
    ]
  },
  {
    "pagePath": "/tax-disc/page2",
    "values": [
      {
        "_count": 1,
        "_end_at": "2014-07-03T00:00:00+00:00",
        "_start_at": "2014-07-02T00:00:00+00:00",
        "uniquePageviews:sum": 25735
      }
    ]
  }
]
}`)))

			server.RouteToHandler("GET", "/data/govuk-info/search-terms",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/data/govuk-info/search-terms"),
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Query().Get("group_by") != "searchKeyword" {
							ghttp.RespondWith(http.StatusOK, `
{
"data": [
	{
		"pagePath": "/tax-disc",
		"values": [
			{
				"_count": 4,
				"_end_at": "2014-09-01T00:00:00+00:00",
				"_start_at": "2014-09-01T00:00:00+00:00",
				"searchUniques:sum": 71
			}
		]
	},
	{
		"pagePath": "/tax-disc/page2",
		"values": [
			{
				"_count": 4,
				"_end_at": "2014-09-02T00:00:00+00:00",
				"_start_at": "2014-09-01T00:00:00+00:00",
				"searchUniques:sum": 75
			}
		]
	}
]
}`)(w, r)
						} else {
							ghttp.RespondWith(http.StatusOK, `
{
"data": [
  {
    "_count": 8,
    "_group_count": 8,
    "searchKeyword": "employer access",
    "searchUniques:sum": 126,
    "values": [{
      "_count": 1,
      "searchUniques:sum": 126,
      "_end_at": "2014-09-02T00:00:00+00:00",
      "_start_at": "2014-09-01T00:00:00+00:00"
    }]
  },
  {
    "_count": 3,
    "_group_count": 3,
    "searchKeyword": "pupil premium",
    "searchUniques:sum": 45,
    "values": [{
      "_count": 1,
      "searchUniques:sum": 45,
      "_end_at": "2014-09-02T00:00:00+00:00",
      "_start_at": "2014-09-01T00:00:00+00:00"
    }]
  },
  {
    "_count": 4,
    "_group_count": 4,
    "searchKeyword": "s2s",
    "searchUniques:sum": 104,
    "values": [{
      "_count": 1,
      "searchUniques:sum": 104,
      "_end_at": "2014-09-02T00:00:00+00:00",
      "_start_at": "2014-09-03T00:00:00+00:00"
    }]
  }
]
}`)(w, r)
						}
					})))

			server.RouteToHandler("GET", "/data/govuk-info/page-contacts",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/data/govuk-info/page-contacts"),
					ghttp.RespondWith(http.StatusOK, `
{
"data": [
	{
		"pagePath": "/tax-disc",
		"values": [
			{
				"_count": 4,
				"_end_at": "2014-09-03T00:00:00+00:00",
				"_start_at": "2014-09-02T00:00:00+00:00",
				"total:sum": 71
			}
		]
	},
	{
		"pagePath": "/tax-disc/page2",
		"values": [
			{
				"_count": 4,
				"_end_at": "2014-09-03T00:00:00+00:00",
				"_start_at": "2014-09-02T00:00:00+00:00",
				"total:sum": 73
			}
		]
	}
]
}`)))

			statistics, err := SlugStatistics(client, "/foo", true)
			Expect(err).To(BeNil())
			Expect(statistics).ToNot(BeNil())
			Expect(len(statistics.PageViews)).To(Equal(2))
			Expect(len(statistics.Searches)).To(Equal(2))
			Expect(len(statistics.SearchTerms)).To(Equal(3))
			Expect(len(statistics.ProblemReports)).To(Equal(2))
			Expect(statistics.PageViews[1].Value).To(Equal(25735))
			Expect(statistics.PageViews[1].Path).To(Equal("/tax-disc/page2"))
			Expect(statistics.Searches[1].Value).To(Equal(75))
			Expect(statistics.Searches[1].Path).To(Equal("/tax-disc/page2"))
			Expect(statistics.ProblemReports[1].Value).To(Equal(73))
			Expect(statistics.ProblemReports[1].Path).To(Equal("/tax-disc/page2"))

			pageViewTimestamp, err := time.Parse(time.RFC3339, "2014-07-02T00:00:00+00:00")
			Expect(err).To(BeNil())
			Expect(statistics.PageViews[1].Timestamp).
				To(Equal(pageViewTimestamp))

			searchesTimestamp, err := time.Parse(time.RFC3339, "2014-09-01T00:00:00+00:00")
			Expect(err).To(BeNil())
			Expect(statistics.Searches[0].Timestamp).
				To(Equal(searchesTimestamp))
		})

	})

})
