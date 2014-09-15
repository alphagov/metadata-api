package performance_platform_test

import (
	"strings"

	. "github.com/alphagov/metadata-api/performance_platform"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data", func() {
	Describe("ParseBackdropResponse", func() {
		It("returns an error when it can't parse the string", func() {
			backdrop, err := ParseBackdropResponse([]byte(""))

			Expect(err).ToNot(BeNil())
			Expect(backdrop).To(BeNil())
		})

		It("parses a Backdrop API response into a Backdrop struct", func() {
			exampleResponse := strings.TrimSpace(`
{
  "data": [
    {
      "_count": 1.0,
      "searchKeyword": "2014 level thresholds",
      "searchUniques": 4,
      "searchUniques:sum": 12.0
    }
  ],
  "warning": "Warning: This data-set is unpublished. Data may be subject to change or be inaccurate."
}`)

			backdrop, err := ParseBackdropResponse([]byte(exampleResponse))
			Expect(err).To(BeNil())
			Expect(backdrop).To(Equal(&Backdrop{
				Data: []Data{Data{
					Count:            1.0,
					SearchKeyword:    "2014 level thresholds",
					SearchUniques:    4,
					SearchUniquesSum: 12,
				}},
				Warning: "Warning: This data-set is unpublished. Data may be subject to change or be inaccurate.",
			}))
		})
	})

	Describe("FetchSlugStatistics", func() {

	})
})
