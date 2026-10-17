# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Club pages and helper links show a bar while there is no connection (by @HB9HIL)

### Fixed
- Training sessions saved offline are no longer lost when the login expired or the upload is refused, stay with the account that typed them, and the app opens offline (by @HB9HIL)
- A failed save without a connection says so instead of "Failed to fetch" (by @HB9HIL)

## [0.17.1] - 2026-10-17

### Added
- CodeQL code scanning and a `govulncheck` job in CI (by @HB9HIL)

### Fixed
- Built with Go 1.27.2, which fixes an HTTP/2 crash in `net/http` (GO-2026-6617) (by @HB9HIL)
- The admin logout cookie is `Secure` and `SameSite=Strict` like the login cookie (by @HB9HIL)

## [0.17.0] - 2026-10-17

First public release.

### Added
- Events with programmes, age categories (SSV rules), online booking and lane reservations (by @HB9HIL)
- Registration desk with lanes, result entry, awards and the award handout (by @HB9HIL)
- Rankings, club competition, special prizes, kiosk displays and live updates (by @HB9HIL)
- Documents (lane sheet, control strip, result sheet, wreath card, rankings) and printing through a print agent at the range (by @HB9HIL)
- Helper links for scorers, ammunition and results desks; personal training log with API tokens (by @HB9HIL)
- User interface, documents and help in English, German, French and Italian (by @HB9HIL)
