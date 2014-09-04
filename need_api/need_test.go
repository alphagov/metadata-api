package need_api_test

import (
	. "github.com/alphagov/metadata-api/need_api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Need", func() {
	Describe("ParseNeed", func() {
		It("returns an error when it can't parse the string", func() {
			need, err := ParseNeedResponse([]byte(""))

			Expect(err).ToNot(BeNil())
			Expect(need).To(BeNil())
		})

		It("can parse a Need API response and marshall it into a Need struct", func() {
			exampleNeedApiResponse := `{
    "_response_info": {
        "status": "ok"
    },
    "id": 100019,
    "role": "To test need code",
    "goal": "provide a test need",
    "benefit": "test",
    "organisation_ids": ["foo-id"],
    "organisations": [{
        "id": "foo-id",
        "name": "Foo Name",
        "govuk_status": "joining",
        "abbreviation": "MHRA",
        "parent_ids": ["department-of-health"]
    }],
    "applies_to_all_organisations": false,
    "justifications": ["This is a test need"],
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

			need, err := ParseNeedResponse([]byte(exampleNeedApiResponse))

			Expect(err).To(BeNil())
			Expect(need).To(Equal(&Need{
				ID:              100019,
				Role:            "To test need code",
				Goal:            "provide a test need",
				Benefit:         "test",
				OrganisationIDs: []string{"foo-id"},
				Organisations: []Organisation{
					Organisation{
						ID:           "foo-id",
						Name:         "Foo Name",
						Status:       "joining",
						Abbreviation: "MHRA",
						ParentIDs:    []string{"department-of-health"},
					},
				},
				Justifications: []string{"This is a test need"},
			}))
		})
	})
})
