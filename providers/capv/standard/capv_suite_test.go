package standard

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/management-cluster-test-suites/v2/internal/suite"
)

func TestCAPVStandard(t *testing.T) {
	suite.Setup()

	RegisterFailHandler(Fail)
	RunSpecs(t, "CAPV Standard Suite")
}
