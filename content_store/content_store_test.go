package content_store_test

import (
	"io/ioutil"
	"os"

	. "github.com/alphagov/metadata-api/content"
	"github.com/alphagov/metadata-api/content_store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type stubRequest struct {
	url         string
	bearerToken string
}

func (req stubRequest) GetJSON(url string, bearerToken string) (string, error) {
	base_url := "http://content-store.dev.gov.uk/content/"
	known_url := base_url + "known"
	unknown_url := base_url + "unknown"
	five_hundred_url := base_url + "five_hundred"
	invalid_response_url := base_url + "invalid_response"
	placeholder := base_url + "placeholder"

	validResponseBytes, _ := ioutil.ReadFile("../fixtures/content_store_response.json")
	validJSONResponse := string(validResponseBytes)

	invalidResponseBytes, _ := ioutil.ReadFile("../fixtures/content_store_response_invalid.json")
	invalidJSONResponse := string(invalidResponseBytes)

	placeholderResponseBytes, _ := ioutil.ReadFile("../fixtures/content_store_response_placeholder.json")
	placeholderJSONResponse := string(placeholderResponseBytes)

	if url == known_url {
		return validJSONResponse, nil
	} else if url == invalid_response_url {
		return invalidJSONResponse, StatusError{404}
	} else if url == unknown_url {
		return "", StatusError{404}
	} else if url == five_hundred_url {
		return "", StatusError{500}
	} else if url == placeholder {
		return placeholderJSONResponse, nil
	} else {
		return "", nil
	}
}

var _ = Describe("content_store", func() {
	Describe("GetArtefact", func() {
		stub := stubRequest{}

		Context("successful request", func() {
			It("requests and returns the the artefact", func() {
				os.Setenv("GOVUK_WEBSITE_ROOT", "http://dev.gov.uk")
				artefact, err := content_store.GetArtefact("known", stub)
				Expect(err).To(BeNil())
				Expect(artefact.ID).To(Equal("73940c62-2580-42b1-9c22-f8e85b71065d"))
				Expect(artefact.WebURL).To(Equal("http://dev.gov.uk/government/get-involved/take-part/volunteer"))
				Expect(artefact.Title).To(Equal("Volunteer"))
				Expect(artefact.Format).To(Equal("take_part"))
				Expect(artefact.Details.NeedIDs).To(Equal([]string{}))
				Expect(artefact.Details.BusinessProposition).To(Equal(false))
				Expect(artefact.Details.Description).To(Equal("Find out how to volunteer in your local community and give your time to help others."))
			})
		})

		Context("content not found", func() {
			It("returns a 404 if the content isn't found", func() {
				artefact, err := content_store.GetArtefact("unknown", stub)
				Expect(err).NotTo(BeNil())
				stErr, _ := err.(StatusError)
				Expect(stErr.StatusCode).To(Equal(404))
				Expect(artefact).To(BeNil())
			})
		})

		Context("request returns a 500", func() {
			It("returns a 500 if the request raises an error", func() {
				artefact, err := content_store.GetArtefact("five_hundred", stub)
				Expect(err).NotTo(BeNil())
				stErr, _ := err.(StatusError)
				Expect(stErr.StatusCode).To(Equal(500))
				Expect(artefact).To(BeNil())
			})
		})

		Context("an invalid content item", func() {
			It("returns an error", func() {
				artefact, err := content_store.GetArtefact("invalid_response", stub)
				Expect(err).NotTo(BeNil())
				Expect(artefact).To(BeNil())
			})
		})

		Context("a placeholder item is returned", func() {
			It("returns a 404 and a nil artefact", func() {
				artefact, err := content_store.GetArtefact("placeholder", stub)
				Expect(err).NotTo(BeNil())
				stErr, _ := err.(StatusError)
				Expect(stErr.StatusCode).To(Equal(404))
				Expect(artefact).To(BeNil())
			})
		})
	})
})
