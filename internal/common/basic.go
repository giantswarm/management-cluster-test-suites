package common

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2" // nolint
	. "github.com/onsi/gomega"    // nolint
	"k8s.io/apimachinery/pkg/types"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/clustertest/v3/pkg/application"
	"github.com/giantswarm/clustertest/v3/pkg/client"
	"github.com/giantswarm/clustertest/v3/pkg/failurehandler"
	"github.com/giantswarm/clustertest/v3/pkg/logger"
	"github.com/giantswarm/clustertest/v3/pkg/organization"
	"github.com/giantswarm/clustertest/v3/pkg/wait"
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
			By(fmt.Sprintf("fetching values for Helm release %s/%s", fakeWC.GetNamespace(), fakeWC.Name))
			Eventually(func() error {
				return state.GetFramework().MC().GetHelmValues(fakeWC.Name, fakeWC.GetNamespace(), values)
			}).
				WithTimeout(10 * time.Minute).
				Should(Succeed())

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
				3, 5*time.Second)).
				WithTimeout(15 * time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(Succeed())
		})

		It("has all default Apps installed successfully", func() {
			defaultAppsSelectorLabels := cr.MatchingLabels{
				"giantswarm.io/cluster":        fakeWC.Name,
				"app.kubernetes.io/managed-by": "Helm",
			}

			appList := &v1alpha1.AppList{}
			err := state.GetFramework().MC().List(state.GetContext(), appList, cr.InNamespace(fakeWC.Organization.GetNamespace()), defaultAppsSelectorLabels)
			Expect(err).NotTo(HaveOccurred())

			appNamespacedNames := []types.NamespacedName{}
			for _, app := range appList.Items {
				appNamespacedNames = append(appNamespacedNames, types.NamespacedName{Name: app.Name, Namespace: app.Namespace})
			}

			Eventually(wait.AreAllAppDeployedSlice(state.GetContext(), state.GetFramework().MC(), appNamespacedNames)).
				WithTimeout(20*time.Minute).
				WithPolling(10*time.Second).
				Should(
					BeEmpty(),
					failurehandler.Bundle(
						failurehandler.AppIssues(state.GetFramework(), fakeWC),
						// TODO: enable once we have a way to report owning teams without coping code across from cluster-test-suites
						// reportOwningTeams(),
					),
				)
		})

		It("all observability-bundle apps are deployed without issues", func() {
			// We need to wait for the observability-bundle app to be deployed before we can check the apps it deploys.
			observabilityAppsAppName := fmt.Sprintf("%s-%s", fakeWC.Name, "observability-bundle")

			Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), observabilityAppsAppName, fakeWC.GetNamespace())).
				WithTimeout(30 * time.Second).
				WithPolling(50 * time.Millisecond).
				Should(BeTrue())

			// Wait for all observability-bundle apps to be deployed
			appList := &v1alpha1.AppList{}
			err := state.GetFramework().MC().List(state.GetContext(), appList, cr.InNamespace(fakeWC.Organization.GetNamespace()), cr.MatchingLabels{"giantswarm.io/managed-by": observabilityAppsAppName})
			Expect(err).NotTo(HaveOccurred())

			appNamespacedNames := []types.NamespacedName{}
			for _, app := range appList.Items {
				appNamespacedNames = append(appNamespacedNames, types.NamespacedName{Name: app.Name, Namespace: app.Namespace})
			}

			Eventually(wait.AreAllAppDeployedSlice(state.GetContext(), state.GetFramework().MC(), appNamespacedNames)).
				WithTimeout(8*time.Minute).
				WithPolling(10*time.Second).
				Should(
					BeEmpty(),
					failurehandler.AppIssues(state.GetFramework(), fakeWC),
				)
		})

		It("all security-bundle apps are deployed without issues", func() {
			// We need to wait for the security-bundle app to be deployed before we can check the apps it deploys.
			securityAppsAppName := fmt.Sprintf("%s-%s", fakeWC.Name, "security-bundle")

			Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), securityAppsAppName, fakeWC.GetNamespace())).
				WithTimeout(30 * time.Second).
				WithPolling(50 * time.Millisecond).
				Should(BeTrue())

			// Wait for all security-bundle apps to be deployed
			appList := &v1alpha1.AppList{}
			err := state.GetFramework().MC().List(state.GetContext(), appList, cr.InNamespace(fakeWC.Organization.GetNamespace()), cr.MatchingLabels{"giantswarm.io/managed-by": securityAppsAppName})
			Expect(err).NotTo(HaveOccurred())

			appNamespacedNames := []types.NamespacedName{}
			for _, app := range appList.Items {
				appNamespacedNames = append(appNamespacedNames, types.NamespacedName{Name: app.Name, Namespace: app.Namespace})
			}

			Eventually(wait.AreAllAppDeployedSlice(state.GetContext(), state.GetFramework().MC(), appNamespacedNames)).
				WithTimeout(10*time.Minute).
				WithPolling(10*time.Second).
				Should(
					BeEmpty(),
					failurehandler.AppIssues(state.GetFramework(), fakeWC),
				)
		})

		It("has all its Deployments Ready (means all replicas are running)", func() {
			Eventually(
				wait.ConsistentWaitConditionSlice(
					wait.AreAllDeploymentsReadySlice(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				),
			).
				WithTimeout(1*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					BeEmpty(),
					failurehandler.DeploymentsNotReady(state.GetFramework(), fakeWC),
				)
		})

		It("has all its StatefulSets Ready (means all replicas are running)", func() {
			Eventually(
				wait.ConsistentWaitConditionSlice(
					wait.AreAllStatefulSetsReadySlice(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					BeEmpty(),
					failurehandler.StatefulSetsNotReady(state.GetFramework(), fakeWC),
				)
		})

		It("has all its DaemonSets Ready (means all daemon pods are running)", func() {
			Eventually(
				wait.ConsistentWaitConditionSlice(
					wait.AreAllDaemonSetsReadySlice(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					BeEmpty(),
					failurehandler.DaemonSetsNotReady(state.GetFramework(), fakeWC),
				)
		})

		It("has all its Jobs completed successfully", func() {
			Eventually(
				wait.ConsistentWaitConditionSlice(
					wait.AreAllJobsSucceededSlice(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					BeEmpty(),
					failurehandler.JobsUnsuccessful(state.GetFramework(), fakeWC),
				)
		})

		It("has all of its Pods in the Running state", func() {
			Eventually(
				wait.ConsistentWaitConditionSlice(
					wait.AreAllPodsInSuccessfulPhaseSlice(state.GetContext(), state.GetFramework().MC()),
					5,
					time.Second,
				)).
				WithTimeout(5*time.Minute).
				WithPolling(wait.DefaultInterval).
				Should(
					BeEmpty(),
					failurehandler.PodsNotReady(state.GetFramework(), fakeWC),
				)
		})

	})
}
