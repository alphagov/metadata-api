package performanceclient

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPerformanceplatformClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PerformanceplatformClient.Go Suite")
}
