package performanceclient

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("MetaClient", func() {
	var server *ghttp.Server
	var client MetaClient

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = NewMetaClient(server.URL(), logrus.New())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Fetch", func() {
		It("Should do a thing", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/"),
					ghttp.RespondWith(http.StatusOK, `{
  "modules": [
    {
      "info":[
      "Data source: Department for Work and Pensions"
      ],
      "value-attribute":"number_of_transactions",
      "description":"",
      "module-type":"kpi",
      "title":"Transactions per year",
      "format":{
        "sigfigs":3,
        "magnitude":true,
        "type":"number"
      },
      "classes":"cols3",
      "slug":"transactions-per-year",
      "data-source":{
        "data-group":"transactional-services",
        "data-type":"summaries",
        "query-params":{
          "sort_by":"_timestamp:descending",
          "filter_by":[
            "service_id:dwp-carers-allowance-new-claims",
            "type:seasonally-adjusted"
          ],
          "group_by": ["a-thing"]
        }
      }
    }
  ],
  "department": {"abbr":"DWP","title":"Department for Work and Pensions"}
}`),
				),
			)
			dashboard, err := client.Fetch("carers-allowance")
			Expect(err).To(BeNil())
			Expect(dashboard).ToNot(BeNil())
			Expect(dashboard.Modules).To(HaveLen(1))
			dataSource := dashboard.Modules[0].DataSource
			Expect(dataSource.DataGroup).To(Equal("transactional-services"))
			Expect(dataSource.DataType).To(Equal("summaries"))
			Expect(dataSource.QueryParams.SortBy).To(Equal("_timestamp:descending"))
		})
	})

	Describe("FetchDashboards", func() {
		It("Should do a thing", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/"),
					ghttp.RespondWith(http.StatusOK, `{
 "items": [
  {
   "department": {
    "abbr": "Home Office", 
    "title": "Home Office"
   }, 
   "agency": {
    "abbr": "DBS", 
    "title": "Disclosure and Barring Service"
   }, 
   "dashboard-type": "high-volume-transaction", 
   "slug": "home-office-enhanced-criminal-records-checks", 
   "title": "Enhanced criminal records checks"
  }, 
  {
   "department": {
    "abbr": "DFT", 
    "title": "Department for Transport"
   }, 
   "agency": {
    "abbr": "DVLA", 
    "title": "Driver and Vehicle Licensing Agency"
   }, 
   "dashboard-type": "high-volume-transaction", 
   "slug": "dft-request-vehicle-tax-refund", 
   "title": "Vehicle tax: refunds"
  }]
}`),
				),
			)
			dashboards, err := client.FetchDashboards()
			Expect(err).To(BeNil())
			Expect(dashboards).ToNot(BeNil())
		})
	})

})
