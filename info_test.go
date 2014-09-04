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
		bearerToken string = "some-secret-bearer-string"

		contentAPIResponse string
		testServer         *httptest.Server
		testContentApi     *httptest.Server
	)

	BeforeEach(func() {
		testContentApi = testHandlerServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+bearerToken {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Not authorised!")
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, contentAPIResponse)
		})
		testServer = testHandlerServer(InfoHandler(testContentApi.URL, bearerToken))
	})

	AfterEach(func() {
		testServer.Close()
		testContentApi.Close()

		contentAPIResponse = `{"_response_info":{"status":"not found"}}`
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
		})

		It("returns the NeedID", func() {
			response, err := getSlug(testServer.URL, "dummy-slug")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())
			Expect(body).To(ContainSubstring(`"_response_info":{"status":"ok"}`))
		})
	})
})

func getSlug(serverURL, slug string) (*http.Response, error) {
	return http.Get(serverURL + "/info/" + slug)
}
