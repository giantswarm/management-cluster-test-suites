package standard

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/management-cluster-test-suites/internal/suite"
)

func TestCAPAStandard(t *testing.T) {
	suite.Setup()

	RegisterFailHandler(Fail)
	RunSpecs(t, "CAPA Standard Suite")
}
