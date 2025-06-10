[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/cluster-test-suites/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/cluster-test-suites/tree/main)

# management-cluster-test-suites

## ☑️ Requirements

* Install [ginkgo](https://onsi.github.io/ginkgo/) on your machine: `go install github.com/onsi/ginkgo/v2/ginkgo`.
* A valid Kubeconfig, pointing at a running `ephemeral` MC. (Note: While this can run against any MC there is no guarantee that there wont be any leftovers or changes made to the MC)
* The `E2E_KUBECONFIG` environment variable set to point to the path of the above kubeconfig.

## 🏃 Running Tests

> [!IMPORTANT]
> The test suites are designed to be run against `ephemeral` MCs and possibly require some config or resources that already exists on those MCs.
>
> If you require running the tests against a different MC please reach out to [#Team-Tenet](https://gigantic.slack.com/archives/C07KSM2E51A) to discuss any pre-requisites that might be needed.

Assuming the above requirements are fulfilled:

* Running all the test suites:

  ```sh
  E2E_KUBECONFIG=/path/to/kubeconfig.yaml ginkgo --timeout 4h -v -r .
  ```

* Running a single provider (e.g. `capa`):

  ```sh
  E2E_KUBECONFIG=/path/to/kubeconfig.yaml ginkgo --timeout 4h -v -r ./providers/capa
  ```

* Running a single test suite (e.g. the `capa` `standard` test suite)

  ```sh
  E2E_KUBECONFIG=/path/to/kubeconfig.yaml ginkgo --timeout 4h -v -r ./providers/capa/standard
  ```

* Running with Docker:

  ```sh
  docker run --rm -it -v /path/to/kubeconfig.yaml:/kubeconfig.yaml -e E2E_KUBECONFIG=/kubeconfig.yaml gsoci.azurecr.io/giantswarm/management-cluster-test-suites ./
  ```

### Testing changes to `clustertest`

To test out changes to [clustertest](https://github.com/giantswarm/clustertest) without needing to create a new release you can add a `replace` directive to your `go.mod` to point to your local copy of `clustertest`. E.g.:

```
module github.com/giantswarm/cluster-test-suites

go 1.24.3

replace github.com/giantswarm/clustertest => /path/to/clustertest
```

### ⚙️ Running Tests in CI

> [!IMPORTANT]
> Currently it's only supported running these tests in CI as part of the [mc-bootstrap](https://github.com/giantswarm/mc-bootstrap) CI PipelineRuns.
>
> There is no support _yet_ for running changes in a PR against this repo.

## ➕ Adding Tests

> See the Ginkgo docs for specifics on how to write tests: https://onsi.github.io/ginkgo/#writing-specs

Where possible, new tests cases should be added that are compatible with all providers so that all benefit. This is obviously not always possible and some provider-specific tests may be required.

All tests make use of our [clustertest](https://github.com/giantswarm/clustertest) test framework. Please refer to the [documentation](https://pkg.go.dev/github.com/giantswarm/clustertest) for more information.

### Adding cross-provider tests

New cross-provider tests should be added to the [`./internal/common/basic.go`](./internal/common/basic.go) package as part of the `RunBasic` function.

A new test case can be included by adding a new `It` block within the `RunBasic` functions. E.g.

```go
It("should test something", func() {
  // Do things...

  Expect(something).To(BeTrue())

  // Cleanup if needed...
})
```

To add a new grouping of common tests you can create a new file with a function similar to `runMyNewGrouping()` and then add a call to this from the [`./internal/common/basic.go`](./internal/common/basic.go) `RunBasic()` function.

### Adding provider-specific tests

Each CAPI provider has its own subdirectory under [`./providers/`](./providers/) that specific tests can be added to.

Each directory under the provider directory is a test suite and each consists to a specific workload cluster configuration variation. All providers should at least contain a `standard` test suite that runs tests for a "default" cluster configuration.

New tests can be added to these provider-specific suites using any Ginkgo context nodes that make sense. Please refer to the [Ginkgo docs](https://onsi.github.io/ginkgo/) for more details.

### Adding Test Suites (new cluster variations in an existing provider)

As mentioned above, test suites are scoped to a single workload cluster configuration variant. To test the different possible configuration options of clusters, for example private networking, we need to create multiple test suites.

A new test suite is added by creating a new module under the provider directory containing the following:

* a `suite_test.go` file
* at least one `*_test.go` file

The `suite_test.go` should mostly be the same across test suites so it will likely be enough to copy the function over from the `standard` test suite and update the names used to represent the test suite being created. This file mostly relies on the [`suite`](./internal/suite/) module to handle the test suite setup and clean up logic.

## Resources

* [`clustertest` documentation](https://pkg.go.dev/github.com/giantswarm/clustertest)
* [CI Tekton Pipeline](https://github.com/giantswarm/mc-bootstrap/blob/main/tekton/pipelines/build-image-generate-mc.yaml)
* [Ginkgo docs](https://onsi.github.io/ginkgo/)
