package performance_platform_test

import (
	"net/http"
	"time"

	. "github.com/alphagov/metadata-api/performance_platform"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Sirupsen/logrus"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Statistics", func() {

	var server *ghttp.Server
	var client *Client

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = NewClient(server.URL(), logrus.New())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("SlugStatistics", func() {
		It("Should return formatted data", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/data/govuk-info/page-statistics"),
					ghttp.RespondWith(http.StatusOK, `
{
"data": [
  {
    "_count": 1,
    "_end_at": "2014-07-04T00:00:00+00:00",
    "_start_at": "2014-07-03T00:00:00+00:00",
    "uniquePageviews:sum": 25931
  }
]
}`),
				),
			)

			statistics, err := client.SlugStatistics("/foo")
			Expect(err).To(BeNil())
			Expect(statistics).ToNot(BeNil())
			Expect(len(statistics.PageViews)).To(Equal(1))
			Expect(statistics.PageViews[0].Value).To(Equal(25931))

			timestamp, err := time.Parse(time.RFC3339, "2014-07-03T00:00:00+00:00")
			Expect(err).To(BeNil())
			Expect(statistics.PageViews[0].Timestamp).
				To(Equal(timestamp))
		})
	})

})
