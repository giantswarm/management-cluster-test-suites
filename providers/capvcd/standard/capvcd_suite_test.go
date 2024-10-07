package standard

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/management-cluster-test-suites/internal/suite"
)

func TestCAPVCDStandard(t *testing.T) {
	suite.Setup()

	RegisterFailHandler(Fail)
	RunSpecs(t, "CAPVCD Standard Suite")
}
