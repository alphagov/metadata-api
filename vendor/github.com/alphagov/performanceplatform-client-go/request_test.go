package performanceclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewRequest", func() {
	It("sets the bearer token in the header when making requests", func() {
		bearerToken := "FOO"

		ts := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+bearerToken {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Not authorized!")
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "You're authorized!")
		})

		defer ts.Close()

		response, err := NewRequest(ts.URL, BearerToken(bearerToken))
		Expect(err).To(BeNil())
		Expect(response.StatusCode).To(Equal(200))

		body, err := ReadResponseBody(response)
		Expect(err).To(BeNil())
		Expect(strings.TrimSpace(string(body))).To(Equal(
			"You're authorized!"))
	})

	It("handles empty bearer tokens", func() {
		bearerToken := ""

		ts := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+bearerToken {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Not authorized!")
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "You're authorized!")
		})

		defer ts.Close()

		response, err := NewRequest(ts.URL, BearerToken(bearerToken))
		Expect(err).To(BeNil())
		Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))

		body, err := ReadResponseBody(response)
		Expect(err).To(BeNil())
		Expect(strings.TrimSpace(string(body))).To(Equal(
			"Not authorized!"))
	})

	It("handles bad networking from the origin server", func() {
		ts := testServer(func(w http.ResponseWriter, r *http.Request) {
			hj, ok := w.(http.Hijacker)
			if !ok {
				panic("webserver doesn't support hijacking – failing the messy way")
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				panic("webserver doesn't support hijacking – failing the messy way")
			}
			// Fail in a clean way so that we don't clutter the output
			conn.Close()
		})
		defer ts.Close()
		// Ensure this isn't a slow test by restricting how many retries happen
		response, err := NewRequest(ts.URL, BearerToken("FOO"), MaxElapsedTime(5*time.Millisecond))
		Expect(response).To(BeNil())
		Expect(err).ShouldNot(BeNil())
	})

	It("retries server unavailable in a forgiving manner", func() {
		semaphore := make(chan struct{})

		ts := testServer(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-semaphore:
				// Second time through, the channel is closed, so we succeed
				w.WriteHeader(http.StatusOK)
			default:
				// First time through, channel gives nothing so we error
				w.WriteHeader(http.StatusServiceUnavailable)
				close(semaphore)
			}
		})
		defer ts.Close()
		response, err := NewRequest(ts.URL, BearerToken("FOO"))
		Expect(response).ShouldNot(BeNil())
		Expect(err).Should(BeNil())
	})

	It("propagates 404s", func() {
		ts := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer ts.Close()
		response, err := NewRequest(ts.URL, BearerToken("FOO"))
		Expect(response).Should(BeNil())
		Expect(err).ShouldNot(BeNil())
		Expect(err).Should(Equal(ErrNotFound))
	})

	It("propagates 400s", func() {
		ts := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{
message: "Either 'duration' or both 'start_at' and 'end_at' are required for a period query",
status: "error"
}`)
		})
		defer ts.Close()
		response, err := NewRequest(ts.URL, BearerToken("FOO"))
		Expect(response).ShouldNot(BeNil())
		Expect(err).ShouldNot(BeNil())
		Expect(err).Should(Equal(ErrBadRequest))
	})
})

func testServer(handler interface{}) *httptest.Server {
	var h http.Handler
	switch handler := handler.(type) {
	case http.Handler:
		h = handler
	case func(http.ResponseWriter, *http.Request):
		h = http.HandlerFunc(handler)
	default:
		// error
		panic("handler cannot be used in an HTTP Server")
	}
	return httptest.NewServer(h)
}
