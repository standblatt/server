# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.17.3] - 2026-10-17

### Added
- One print agent can serve several ranges: `STANDBLATT_AGENT_TOKEN` takes several tokens separated by commas (by @HB9HIL)
- Duplicate clubs can be merged on the Clubs page, including their registrations in closed events (by @HB9HIL)

### Changed
- Lane groups are gone: each distance is a range of its own with lanes from 1, e.g. *Pistol 25 m* and *Pistol 50 m*, so online booking, reservations and scorer links work per distance. Existing groups become ranges on upgrade, with registrations, reservations, scorer links, programmes and print settings moved along; add a print agent token for each new range, and check kiosk links limited to a range and scorer links whose lanes spanned both distances (by @HB9HIL)
- Club names are matched against the club list regardless of upper case and spaces; a new shooter's club, a CSV import and a confirmed lane booking suggest similar clubs and only create a club after an explicit choice (by @HB9HIL)

## [0.17.2] - 2026-10-17

### Added
- Club pages and helper links show a bar while there is no connection (by @HB9HIL)

### Changed
- Rankings offer PDF and CSV instead of Print; rankings are no longer sent to the range's printer (by @HB9HIL)
- The sidebar groups Shooters and Clubs under Master data (by @HB9HIL)
- On phones and tablets PDFs go to the share sheet instead of a new tab (by @HB9HIL)
- On phones tables show their main columns with the rest below the name, and the setup tabs become a list, so no page scrolls sideways (by @HB9HIL)

### Fixed
- PDF links on helper pages work (they were sent without the helper code) (by @HB9HIL)
- On phones the sidebar closes after choosing a page (by @HB9HIL)
- Installed on an iPhone, the app keeps its headers below the status bar and nothing scrolls behind it (by @HB9HIL)
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
