# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Go: Update dependencies.

## [0.7.6] - 2025-08-28

### Changed

- Go: Update dependencies.

## [0.7.5] - 2025-08-23

### Changed

- Go: Upgrade `clustertest` to v1.39.2, downgrade Cluster API to v1.10.5.

## [0.7.4] - 2025-08-22

### Changed

- Go: Update dependencies.

## [0.7.3] - 2025-08-07

### Changed

- Go: Update dependencies.

## [0.7.2] - 2025-07-30

### Changed

- Go: Update dependencies.

## [0.7.1] - 2025-06-13

### Fixed

- Fix Slice-based tests incorrectly still using `Succeed()` check

## [0.7.0] - 2025-06-12

### Changed

- Switched to using WaitConditionSlice functions to provide more information in the test failure messages

## [0.6.0] - 2025-06-03

### Added

- Added tests to check for Apps being installed successfully including default apps, observability bundle and security bundle

### Changed

- Reduced consistent check for worker nodes from 12 down to 3

## [0.5.0] - 2025-05-30

### Added

- Added standard test for all worker nodes being ready

### Fixed

- Fix linting issues.

## [0.4.0] - 2024-10-07

### Added

- Added test suites for all supported providers

## [0.3.0] - 2024-10-03

### Added

- Added failure handler for not-ready Pods test

## [0.2.1] - 2024-10-01

### Fixed

- Bump `clustertest` with fixes for `GetLogs` and excluding running Job when checkking for unsuccessful

## [0.2.0] - 2024-10-01

### Added

- Added failure handlers to test cases to get more debug info on failure

## [0.1.1] - 2024-09-24

### Fixed

- Added missing Dockerfile

## [0.1.0] - 2024-09-23

## Added

- Initial basic cluster-wide tests
- CAPA standard test suite

[Unreleased]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.7.6...HEAD
[0.7.6]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.7.5...v0.7.6
[0.7.5]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.7.4...v0.7.5
[0.7.4]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.7.3...v0.7.4
[0.7.3]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.7.2...v0.7.3
[0.7.2]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.7.1...v0.7.2
[0.7.1]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/giantswarm/management-cluster-test-suites/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/management-cluster-test-suites/releases/tag/v0.1.0
