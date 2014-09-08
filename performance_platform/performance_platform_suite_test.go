package performance_platform_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPerformancePlatform(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PerformancePlatform Suite")
}
