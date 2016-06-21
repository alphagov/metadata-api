package content_test

import (
	. "github.com/alphagov/metadata-api/content"
	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ApiRequest", func() {
	Describe("Get", func() {
		var apiRequest = ApiRequest{}

		Context("url contains content", func() {
			It("returns a string of the body content at the url", func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				var url = "https://test.com"
				var response = `{"id": 1}`
				httpmock.RegisterResponder("GET", url,
					httpmock.NewStringResponder(200, response))

				responseString, err := apiRequest.GetJSON(url, "")
				Expect(err).To(BeNil())
				Expect(responseString).To(Equal(response))
			})

			It("returns a 404 for not found", func() {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				var url = "https://test.com"
				httpmock.RegisterResponder("GET", url,
					httpmock.NewStringResponder(404, ""))

				responseString, err := apiRequest.GetJSON(url, "")
				stErr, _ := err.(StatusError)
				Expect(responseString).To(Equal(""))
				Expect(stErr.StatusCode).To(Equal(404))
			})
		})
	})
})
