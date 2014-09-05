package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/alphagov/metadata-api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info", func() {
	var (
		contentAPIResponse, needAPIResponse     string
		testServer, testContentApi, testNeedApi *httptest.Server

		contentAPIBearerToken = "some-secret-content-api-bearer-string"
		needAPIBearerToken    = "some-secret-need-api-bearer-string"
	)

	BeforeEach(func() {
		testContentApi = testHandlerServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+contentAPIBearerToken {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Not authorised!")
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, contentAPIResponse)
		})
		testNeedApi = testHandlerServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+needAPIBearerToken {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Not authorised!")
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, needAPIResponse)
		})

		testServer = testHandlerServer(InfoHandler(
			testContentApi.URL, testNeedApi.URL, contentAPIBearerToken, needAPIBearerToken))
	})

	AfterEach(func() {
		testServer.Close()
		testContentApi.Close()
		testNeedApi.Close()

		contentAPIResponse = `{"_response_info":{"status":"not found"}}`
		needAPIResponse = `{"_response_info":{"status":"not found"}}`
	})

	Describe("no slug provided", func() {
		It("returns a 404 without a trailing slash", func() {
			response, err := http.Get(testServer.URL + "/info")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())
			Expect(body).To(ContainSubstring(`"_response_info":{"status":"not found"}`))
		})

		It("returns a 404 with a trailing slash", func() {
			response, err := http.Get(testServer.URL + "/info/")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())
			Expect(body).To(ContainSubstring(`"_response_info":{"status":"not found"}`))
		})
	})

	Describe("fetching a valid slug", func() {
		BeforeEach(func() {
			contentAPIResponse = `{
  "id": "https://www.gov.uk/api/driving-licence-fees.json",
  "web_url": "https://www.gov.uk/driving-licence-fees",
  "title": "Driving licence fees",
  "format": "answer",
  "updated_at": "2014-06-27T14:21:48+01:00",
  "details": {
    "need_ids": ["100567"],
    "language": "en",
    "body": "foo"
  },
  "_response_info": {
    "status": "ok"
  }
}`
			needAPIResponse = `{
  "_response_info": {
    "status": "ok"
  },
  "id": 100019,
  "role": "Someone carrying out a clinical trial",
  "goal": "maintain my clinical trial authorisation",
  "benefit": "ensure that my clinical trial continues to meet MHRA requirements and the appropriate legal criteria",
  "organisation_ids": ["medicines-and-healthcare-products-regulatory-agency"],
  "organisations": [{
    "id": "medicines-and-healthcare-products-regulatory-agency",
    "name": "Medicines and Healthcare Products Regulatory Agency",
    "govuk_status": "joining",
    "abbreviation": "MHRA",
    "parent_ids": ["department-of-health"],
    "child_ids": []
  }],
  "applies_to_all_organisations": false,
  "justifications": ["The government is legally obliged to provide it", "It's something that people can do or it's something people need to know before they can do something that's regulated by/related to government"],
  "impact": null,
  "met_when": null,
  "yearly_user_contacts": null,
  "yearly_site_views": null,
  "yearly_need_views": null,
  "yearly_searches": null,
  "other_evidence": null,
  "legislation": null,
  "in_scope": null,
  "out_of_scope_reason": null,
  "duplicate_of": null
}`
		})

		It("returns a metadata response with the Artefact and Needs exposed", func() {
			response, err := getSlug(testServer.URL, "dummy-slug")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())

			metadata, err := ParseMetadataResponse([]byte(body))
			Expect(err).To(BeNil())

			Expect(metadata.ResponseInfo.Status).To(Equal("ok"))
			Expect(metadata.Artefact.Details.NeedIDs).To(Equal([]string{"100567"}))
			Expect(metadata.Needs).To(HaveLen(1))
			Expect(metadata.Needs[0].ID).To(Equal(100019))
		})
	})
})

func getSlug(serverURL, slug string) (*http.Response, error) {
	return http.Get(serverURL + "/info/" + slug)
}
