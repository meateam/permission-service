# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- hotfix: ([21](https://github.com/meateam/permission-service/pull/21)): changed updated_at, created_at to updatedAt, createdAt

### Added

- minor: ([21](https://github.com/meateam/permission-service/pull/21)): added updatedAt and createdAt fields, pagination sorting defaults to updated_at field

## [v3.0.0] - 2021-01-13

### Added

- major: ([10](https://github.com/meateam/permission-service/pull/10)): added pagination, and appID in create
- FEAT([17](https://github.com/meateam/permission-service/pull/17)): add ci and test file

### Refactor

- REFACTOR([95](https://github.com/meateam/drive-project/issues/96)): upgrade docker compose to v3 and use env_file

## [v2.0.0] - 2020-10-28

### Added

- FEAT([9](https://github.com/meateam/permission-service/pull/9)): RPC method GetPermissionByMongoID

[unreleased]: https://github.com/meateam/permission-service/compare/master...develop
[v3.0.0]: https://github.com/meateam/permission-service/compare/v2.0.0...v3.0.0
[v2.0.0]: https://github.com/meateam/permission-service/compare/v1.3...v2.0.0
