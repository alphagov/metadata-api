package content_api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestContentApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContentApi Suite")
}
