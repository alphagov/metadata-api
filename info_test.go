package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/alphagov/metadata-api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info", func() {
	var (
		contentAPIResponse, needAPIResponse, performanceAPIResponse string
		testServer, testContentAPI, testNeedAPI, testPerformanceAPI *httptest.Server

		config = &Config{
			BearerTokenContentAPI: "some-secret-content-api-bearer-string",
			BearerTokenNeedAPI:    "some-secret-need-api-bearer-string",
		}
	)

	BeforeEach(func() {
		testContentAPI = testHandlerServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+config.BearerTokenContentAPI {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Not authorised!")
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, contentAPIResponse)
		})
		testNeedAPI = testHandlerServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+config.BearerTokenNeedAPI {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Not authorised!")
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, needAPIResponse)
		})
		testPerformanceAPI = testHandlerServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, performanceAPIResponse)
		})

		testServer = testHandlerServer(InfoHandler(
			testContentAPI.URL, testNeedAPI.URL, testPerformanceAPI.URL, config))
	})

	AfterEach(func() {
		testServer.Close()
		testContentAPI.Close()
		testNeedAPI.Close()

		contentAPIResponse = `{"_response_info":{"status":"not found"}}`
		needAPIResponse = `{"_response_info":{"status":"not found"}}`
		performanceAPIResponse = `{
          "data": [],
          "warning": "Warning: This data-set is unpublished. Data may be subject to change or be inaccurate."
        }`
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
			contentAPIResponseBytes, _ := ioutil.ReadFile("fixtures/content_api_response.json")
			needAPIResponseBytes, _ := ioutil.ReadFile("fixtures/need_api_response.json")
			performanceAPIResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_response.json")

			contentAPIResponse = string(contentAPIResponseBytes)
			needAPIResponse = string(needAPIResponseBytes)
			performanceAPIResponse = string(performanceAPIResponseBytes)
		})

		It("returns a metadata response with the Artefact, Needs, and Performance Data exposed", func() {
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

			Expect(metadata.Performance.Data[0].PagePath).To(Equal("/intellectual-property-an-overview"))
			Expect(metadata.Performance.Data[0].UniquePageViews).To(Equal(float32(102)))
		})
	})

	Describe("querying for a slug that doesn't exist", func() {
		BeforeEach(func() {
			testContentAPI = testHandlerServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer "+config.BearerTokenContentAPI {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintln(w, "Not authorised!")
					return
				}

				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, contentAPIResponse)
			})

			testServer = testHandlerServer(InfoHandler(
				testContentAPI.URL, testNeedAPI.URL, testPerformanceAPI.URL, config))
		})

		AfterEach(func() {
			testContentAPI.Close()
		})

		It("returns with a status of not found if there's no slug in the Content API", func() {
			response, err := getSlug(testServer.URL, "not-found-slug")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())

			metadata, err := ParseMetadataResponse([]byte(body))
			Expect(err).To(BeNil())

			Expect(metadata.ResponseInfo.Status).To(Equal("not found"))
			Expect(metadata.Artefact).To(BeNil())
			Expect(metadata.Needs).To(BeNil())
		})
	})
})

func getSlug(serverURL, slug string) (*http.Response, error) {
	return http.Get(serverURL + "/info/" + slug)
}
