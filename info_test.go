package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/alphagov/metadata-api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kr/pretty"
)

type stubbedJSONRequest struct {
	Response string
}

func (apiRequest stubbedJSONRequest) GetJSON(url string, bearerToken string) (string, error) {
	return apiRequest.Response, nil
}

var _ = Describe("Info", func() {
	var (
		contentAPIResponse, contentStoreResponse, needAPIResponse, pageviewsResponse,
		searchesResponse, problemReportsResponse, termsResponse string

		testServer, testContentAPI, testNeedAPI, testPerformanceAPI *httptest.Server

		testApiRequest stubbedJSONRequest

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
			if strings.Contains(r.URL.Path, "page-statistics") {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, pageviewsResponse)
			} else if strings.Contains(r.URL.Path, "search-terms") &&
				r.URL.Query().Get("group_by") == "searchKeyword" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, termsResponse)
			} else if strings.Contains(r.URL.Path, "search-terms") {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, searchesResponse)
			} else if strings.Contains(r.URL.Path, "page-contacts") {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, problemReportsResponse)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		})

		testApiRequest = stubbedJSONRequest{
			contentStoreResponse,
		}

		testServer = testHandlerServer(InfoHandler(
			testContentAPI.URL, testNeedAPI.URL, testPerformanceAPI.URL, testApiRequest, config))
	})

	AfterEach(func() {
		testServer.Close()
		testContentAPI.Close()
		testNeedAPI.Close()

		contentAPIResponse = `{"_response_info":{"status":"not found"}}`
		contentStoreResponse = ``
		needAPIResponse = `{"_response_info":{"status":"not found"}}`
		searchesResponse = `{"data":[]}`
		pageviewsResponse = `{"data":[]}`
		termsResponse = `{"data":[]}`
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
			pageviewsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_pageviews_response.json")
			searchesResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_searches_response.json")
			problemReportsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_problem_reports_response.json")
			termsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_terms_response.json")

			contentAPIResponse = string(contentAPIResponseBytes)
			needAPIResponse = string(needAPIResponseBytes)
			pageviewsResponse = string(pageviewsResponseBytes)
			searchesResponse = string(searchesResponseBytes)
			problemReportsResponse = string(problemReportsResponseBytes)
			termsResponse = string(termsResponseBytes)
		})

		It("returns a metadata response with the Artefact, Needs and Performance data exposed", func() {
			response, err := getSlug(testServer.URL, "dummy-slug")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())

			expectedResultBytes, _ := ioutil.ReadFile("fixtures/info_response_single_page.json")
			diff := Diff(string(expectedResultBytes), body)
			Expect(diff).To(BeNil())
		})
	})

	Describe("fetching a valid slug from the content store", func() {
		BeforeEach(func() {
			contentStoreResponseBytes, _ := ioutil.ReadFile("fixtures/content_store_response.json")
			needAPIResponseBytes, _ := ioutil.ReadFile("fixtures/need_api_response.json")
			pageviewsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_pageviews_response.json")
			searchesResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_searches_response.json")
			problemReportsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_problem_reports_response.json")
			termsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_terms_response.json")

			contentStoreResponse = string(contentStoreResponseBytes)
			needAPIResponse = string(needAPIResponseBytes)
			pageviewsResponse = string(pageviewsResponseBytes)
			searchesResponse = string(searchesResponseBytes)
			problemReportsResponse = string(problemReportsResponseBytes)
			termsResponse = string(termsResponseBytes)
		})

		It("returns a metadata response with the Artefact, Needs and Performance data exposed", func() {
			response, err := getSlug(testServer.URL, "dummy-slug")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())

			expectedResultBytes, _ := ioutil.ReadFile("fixtures/info_response_content_store.json")
			trimmedString := strings.TrimSpace(string(expectedResultBytes))
			diff := Diff(trimmedString, body)
			Expect(diff).To(BeNil())
		})
	})

	Describe("fetching a slug without need_ids", func() {
		BeforeEach(func() {
			contentAPIResponseBytes, _ := ioutil.ReadFile("fixtures/content_api_response_without_need_ids.json")
			pageviewsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_pageviews_response.json")
			searchesResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_searches_response.json")
			problemReportsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_problem_reports_response.json")
			termsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_terms_response.json")

			contentAPIResponse = string(contentAPIResponseBytes)
			pageviewsResponse = string(pageviewsResponseBytes)
			searchesResponse = string(searchesResponseBytes)
			problemReportsResponse = string(problemReportsResponseBytes)
			termsResponse = string(termsResponseBytes)
		})

		It("returns a metadata response with the an empty Needs array", func() {
			response, err := getSlug(testServer.URL, "dummy-slug")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())

			expectedResponseBytes, _ := ioutil.ReadFile("fixtures/info_response_empty_needs.json")
			diff := Diff(string(expectedResponseBytes), body)
			Expect(diff).To(BeNil())
		})

	})

	Describe("fetching a valid slug with a multipart format", func() {
		BeforeEach(func() {
			contentAPIResponseBytes, _ := ioutil.ReadFile("fixtures/content_api_response_with_parts.json")
			pageviewsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_pageviews_multipart_response.json")
			searchesResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_searches_multipart_response.json")
			problemReportsResponseBytes, _ := ioutil.ReadFile("fixtures/performance_platform_problem_reports_multipart_response.json")

			contentAPIResponse = string(contentAPIResponseBytes)
			pageviewsResponse = string(pageviewsResponseBytes)
			searchesResponse = string(searchesResponseBytes)
			problemReportsResponse = string(problemReportsResponseBytes)
		})

		It("returns a metadata response with a parts array, and handles multipart data correctly", func() {
			response, err := getSlug(testServer.URL, "dummy-slug")
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			body, err := readResponseBody(response)
			Expect(err).To(BeNil())

			expectedResultBytes, _ := ioutil.ReadFile("fixtures/info_response_multipart.json")
			diff := Diff(string(expectedResultBytes), body)
			Expect(diff).To(BeNil())
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
				testContentAPI.URL, testNeedAPI.URL, testPerformanceAPI.URL, testApiRequest, config))
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
			Expect(body).To(ContainSubstring("\"status\":\"not found\""))
		})
	})
})

func getSlug(serverURL, slug string) (*http.Response, error) {
	return http.Get(serverURL + "/info/" + slug)
}
