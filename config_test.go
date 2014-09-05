package main_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/alphagov/metadata-api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("ReadConfig", func() {
		It("can read and parse a config file from the path", func() {
			configFile := "test_config.json"
			configData := `{"bearer_token_content_api": "foo", "bearer_token_need_api": "bar"}`

			workingDirectory, _ := filepath.Abs(filepath.Dir(os.Args[0]))

			err := ioutil.WriteFile(workingDirectory+"/"+configFile, []byte(configData), 0644)
			Expect(err).To(BeNil())

			config, err := ReadConfig(configFile)
			Expect(err).To(BeNil())
			Expect(config).To(Equal(&Config{
				BearerTokenContentAPI: "foo",
				BearerTokenNeedAPI:    "bar",
			}))
		})
	})
})
