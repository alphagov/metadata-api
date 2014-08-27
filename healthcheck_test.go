package main_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/alphagov/metadata-api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Healthcheck", func() {
	It("responds simply with 'OK'", func() {
		testServer := testHandlerServer(HealthCheckHandler)
		defer testServer.Close()

		response, err := http.Get(testServer.URL)
		Expect(err).To(BeNil())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		body, err := readResponseBody(response)
		Expect(err).To(BeNil())
		Expect(body).To(Equal("OK"))
	})
})

func testHandlerServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func readResponseBody(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	return strings.TrimSpace(string(body)), err
}
