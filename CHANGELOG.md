# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Upcoming changes...

## [0.9.0] - 2025-03-31
### Added
- Added HTTP response modifier
- Added reflection option

## [0.8.0] - 2024-12-05
### Added
- Allow change CN on gRPC Gateway

## [0.7.0] - 2024-08-22
### Added
- Support for configuring custom CommonName (CN) values in TLS certificates, enabling more flexible certificate management and custom domain setups

## [0.6.0] - 2024-08-08
### Added
- Added 'sqlite' DB query string
- Switched to `modernc.org/sqlite` driver
- Upgraded to go 1.22

## [0.5.1] - 2024-06-24
### Added
- Added case-insensitive like operator lookup based on DB type

## [0.5.0] - 2024-06-24
### Added
- Added DB query tracing option

## [0.4.0] - 2024-06-20
### Added
- Upgraded to goland 1.20
- Added optional OTEL parameter initialisation

## [0.3.0] - 2023-11-27
### Added
- Added SQL query tracing support

## [0.2.0] - 2023-11-20
### Added
- Added Open Telemetry (OTEL) export support
- Upgraded to Go 1.20


## [0.1.0] - 2023-03-08
### Added
- Add database helpers
- Added file helpers
- Added gRPC helpers
- Added HTTP/REST gateway helper

[0.1.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.0.0...v0.1.0
[0.2.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.1.0...v0.2.0
[0.3.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.2.0...v0.3.0
[0.4.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.3.0...v0.4.0
[0.5.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.4.0...v0.5.0
[0.5.1]: https://github.com/scanoss/go-grpc-helper/compare/v0.5.0...v0.5.1
[0.6.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.5.1...v0.6.0
[0.7.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.6.0...v0.7.0
[0.8.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.7.0...v0.8.0
[0.9.0]: https://github.com/scanoss/go-grpc-helper/compare/v0.8.0...v0.9.0