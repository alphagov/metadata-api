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
)

var _ = Describe("NewRequest", func() {
	It("sets the bearer token in the header when making requests", func() {
		bearerToken := "FOO"

		ts := testServer(bearerToken)
		defer ts.Close()

		response, err := NewRequest(ts.URL, bearerToken)
		defer response.Body.Close()
		Expect(err).To(BeNil())
		Expect(response.StatusCode).To(Equal(200))

		body, err := ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())
		Expect(strings.TrimSpace(string(body))).To(Equal(
			"You're authorised!"))
	})
})

func testServer(bearerToken string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+bearerToken {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Not authorised!")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "You're authorised!")
	}))
}
