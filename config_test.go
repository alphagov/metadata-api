package main_test

import (
	"os"

	. "github.com/alphagov/metadata-api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("InitConfig", func() {
		It("can read and parse a config file from the path", func() {
			os.Setenv("CONTENT_API_BEARER_TOKEN", "foo")
			os.Setenv("NEED_API_BEARER_TOKEN", "bar")

			config := InitConfig()
			Expect(config).To(Equal(&Config{
				BearerTokenContentAPI: "foo",
				BearerTokenNeedAPI:    "bar",
			}))

			os.Unsetenv("CONTENT_API_BEARER_TOKEN")
			os.Unsetenv("NEED_API_BEARER_TOKEN")
		})
	})
})
