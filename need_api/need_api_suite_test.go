package need_api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNeedApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NeedApi Suite")
}
