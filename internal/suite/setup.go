package suite

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2" // nolint
	. "github.com/onsi/gomega"    // nolint

	"github.com/giantswarm/clustertest/v3"
	"github.com/giantswarm/clustertest/v3/pkg/logger"

	"github.com/giantswarm/management-cluster-test-suites/internal/state"
)

// Setup handles the creation of the BeforeSuite and AfterSuite handlers. This covers the creations and cleanup of the test cluster.
// `clusterReadyFns` can be provided if the cluster requires custom checks for cluster-ready status. If not provided the cluster will
// be checked for at least a single control plane node being marked as ready.
func Setup() {
	BeforeSuite(func() {
		logger.LogWriter = GinkgoWriter

		state.SetContext(context.Background())

		framework, err := clustertest.New("")
		Expect(err).NotTo(HaveOccurred())
		state.SetFramework(framework)

		ctx := context.Background()
		_, cancelApplyCtx := context.WithTimeout(ctx, 20*time.Minute)
		defer cancelApplyCtx()

		// In certain cases, when connecting over the VPN, it is possible that the tunnel
		// isn't ready and can take a short while to become usable. This attempts to wait
		// for the connection to be usable before starting the tests.
		Eventually(func() error {
			logger.Log("Checking connection to MC is available.")
			logger.Log("MC API Endpoint: '%s'", state.GetFramework().MC().GetAPIServerEndpoint())
			logger.Log("MC name: '%s'", state.GetFramework().MC().GetClusterName())
			return state.GetFramework().MC().CheckConnection()
		}).
			WithTimeout(5 * time.Minute).
			WithPolling(5 * time.Second).
			Should(Succeed())
	})

	AfterSuite(func() {
		// Ensure we reset the context timeout to make sure we allow plenty of time to clean up
		ctx := state.GetContext()
		ctx, _ = context.WithTimeout(ctx, 1*time.Hour) // nolint
		state.SetContext(ctx)

		// TODO: Any cleanup ?
	})
}
