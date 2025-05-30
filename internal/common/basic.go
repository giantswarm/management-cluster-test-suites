package common

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2" // nolint
	. "github.com/onsi/gomega"    // nolint

	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/client"
	"github.com/giantswarm/clustertest/pkg/failurehandler"
	"github.com/giantswarm/clustertest/pkg/logger"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/wait"
	cr "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/management-cluster-test-suites/internal/state"
)

func RunBasic() {
	Context("basic", func() {
		var fakeWC *application.Cluster

		BeforeEach(func() {
			fakeWC = &application.Cluster{
				Name:         state.GetFramework().MC().GetClusterName(),
				Organization: organization.New("giantswarm"),
			}
		})

		It("should be able to connect to the management cluster", FlakeAttempts(3), func() {
			Expect(state.GetFramework().MC().CheckConnection()).To(Succeed())
		})

		It("has all the control-plane nodes running", func() {
			Eventually(
				wait.ConsistentWaitCondition(
					wait.AreNumNodesReady(state.GetContext(), state.GetFramework().MC(), int(3), &cr.MatchingLabels{"node-role.kubernetes.io/control-plane": ""}),
					5,
					5*time.Second,
				)).
				WithTimeout(5 * time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(Succeed())
		})

		It("has all the worker nodes running", func() {
			values := &application.ClusterValues{}
			err := state.GetFramework().MC().GetHelmValues(fakeWC.Name, fakeWC.GetNamespace(), values)
			Expect(err).NotTo(HaveOccurred())

			minNodes := 0
			maxNodes := 0
			for _, pool := range values.NodePools {
				if pool.Replicas > 0 {
					minNodes += pool.Replicas
					maxNodes += pool.Replicas
					continue
				}

				minNodes += pool.MinSize
				maxNodes += pool.MaxSize
			}
			expectedNodes := wait.Range{
				Min: minNodes,
				Max: maxNodes,
			}

			workersFunc := func() (bool, error) {
				ok, err := wait.AreNumNodesReadyWithinRange(state.GetContext(), state.GetFramework().MC(), expectedNodes, client.DoesNotHaveLabels{"node-role.kubernetes.io/control-plane"})()
				if err != nil {
					logger.Log("failed to get nodes: %s", err)
					return false, err
				}
				if !ok {
					return false, fmt.Errorf("unexpected number of nodes")
				}
				return true, nil
			}

			Eventually(wait.Consistent(
				func() error {
					ok, err := workersFunc()
					if err != nil {
						logger.Log("failed to get nodes: %s", err)
						return err
					}
					if !ok {
						return fmt.Errorf("unexpected number of nodes")
					}
					return nil
				},
				12, 5*time.Second)).
				WithTimeout(15 * time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(Succeed())
		})

		It("has all its Deployments Ready (means all replicas are running)", func() {
			Eventually(
				wait.ConsistentWaitCondition(
					wait.AreAllDeploymentsReady(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					Succeed(),
					failurehandler.DeploymentsNotReady(state.GetFramework(), fakeWC),
				)
		})

		It("has all its StatefulSets Ready (means all replicas are running)", func() {
			Eventually(
				wait.ConsistentWaitCondition(
					wait.AreAllStatefulSetsReady(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					Succeed(),
					failurehandler.StatefulSetsNotReady(state.GetFramework(), fakeWC),
				)
		})

		It("has all its DaemonSets Ready (means all daemon pods are running)", func() {
			Eventually(
				wait.ConsistentWaitCondition(
					wait.AreAllDaemonSetsReady(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					Succeed(),
					failurehandler.DaemonSetsNotReady(state.GetFramework(), fakeWC),
				)
		})

		It("has all its Jobs completed successfully", func() {
			Eventually(
				wait.ConsistentWaitCondition(
					wait.AreAllJobsSucceeded(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					Succeed(),
					failurehandler.JobsUnsuccessful(state.GetFramework(), fakeWC),
				)
		})

		It("has all of its Pods in the Running state", func() {
			Eventually(
				wait.ConsistentWaitCondition(
					wait.AreAllPodsInSuccessfulPhase(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					Succeed(),
					failurehandler.PodsNotReady(state.GetFramework(), fakeWC),
				)
		})

	})
}
