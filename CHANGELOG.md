# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- On iPhones, date and time fields keep the width and height of the fields next to them instead of reaching into them, e.g. *When* and *Where* when saving a training session (by @HB9HIL)
- The server refuses changing requests a browser sends from another origin, so a page on another subdomain can no longer submit a file upload (shooter or result import, images) with the session of a logged-in user (by @HB9HIL)
- The server no longer starts with `STANDBLATT_SMTP_HOST` but without `STANDBLATT_PUBLIC_URL`, so links in mails (password reset, invitations) never take the host from a forged request (by @HB9HIL)
- Built with Go 1.27.3, which distrusts compromised intermediate certificates in `crypto/x509` (GO-2026-6809) when sending mails over TLS (by @HB9HIL)

## [0.22.2] - 2026-10-19

### Added
- Help: *Printing* describes the Standblatt Print Agent app for Windows, macOS and Linux (download, install, connect, the window, updates, removing it) with screenshots of the app (by @HB9HIL)

### Fixed
- The German, French and Italian help keeps the new button name *Set up the printers* and the code field when the help translations are regenerated (by @HB9HIL)

## [0.22.1] - 2026-10-19

### Added
- *Ranges → Set up the printers* has *Connect by code*: type the code the Standblatt app shows to connect that computer to the range, without the link from the app (by @HB9HIL)

## [0.22.0] - 2026-10-19

### Removed
- The print agent (`cmd/print-agent`) and its binaries in the releases: it lives in [standblatt/print-agent](https://github.com/standblatt/print-agent) with a desktop app for Windows, macOS and Linux; the documentation points there (by @HB9HIL)

## [0.21.1] - 2026-10-19

### Changed
- The button to a range's printers on *Ranges* is called *Set up the printers* instead of *Printing* (by @HB9HIL)

### Fixed
- *Ranges → Set up the printers* shows a print agent that stops calling in as *not connected* in red, with when it was last seen, after a minute and without reloading the page (by @HB9HIL)

## [0.21.0] - 2026-10-18

### Added
- **Connecting a print agent by code**: the agent asks for a code (`POST /api/agent/pair`), a member who may manage ranges confirms it for a range on the page `/connect`, and the agent fetches its token once (`POST /api/agent/pair/status`); codes are valid for ten minutes and limited per IP (by @HB9HIL)
- Print agents report their version and operating system; *Ranges → Printing* shows them (by @HB9HIL)
- **Test page** in *Ranges → Printing*: each document row sends a small page (90 × 29 mm) with club, range, printer, options and time to its printer with the options of that row (by @HB9HIL)

### Changed
- *Ranges → Printing* offers every option a print agent reports for a printer, with plain names for the options of the Windows agent; the free-text `lp` options remain for printers of other agents, and a saved setting keeps only the options of its printer (by @HB9HIL)
- The print agent's hello requires `version` and `os` and answers with the club, the range and the routed printers the agent does not report (by @HB9HIL)

## [0.20.1] - 2026-10-18

### Changed
- An API token that lacks a scope its owner's role has is refused with "Not allowed for this API token" instead of naming the role (by @HB9HIL)
- API tokens may call `/api/version` and `/api/auth/options`, and with the permission *Configure events* the module list and the programme library (by @HB9HIL)
- Every API operation is tested with a session and with an API token; the tests fail when an operation has no test that calls it (by @HB9HIL)

## [0.20.0] - 2026-10-18

### Added
- **Roles per club**: the *Users* page has a permission matrix where a club makes its own roles and decides what each may do (see events, run events, configure events, statistics, master data, ranges, print forms, users); admins may always do everything (by @HB9HIL)

### Changed
- Club roles are made of permissions: admin and staff became the roles *Admin* and *Helper* (unchanged rights), new default roles are *Board* and *Member* (sees only the events and rankings); the app shows pages and buttons by permission (by @HB9HIL)
- **Invitations** wait for an answer: the invited person accepts through the link in the mail, with a new account and password or, after logging in, with the account they have; the *Users* page lists open and expired invitations to send again or withdraw. Inviting no longer creates accounts (by @HB9HIL)
- Every API operation of a club declares its permissions, checked before the request is read; tests call every operation with and without them and keep the table of who may call what under review (by @HB9HIL)
- Members who manage the club without the admin role stay within their own permissions: only admins make, change or remove admins, and nobody gives a role more than their own (by @HB9HIL)
- The *Users* page shows members without the admin role only the roles, members, invitations and permissions they may give or change (by @HB9HIL)
- **API tokens** moved from *My training* to the account and reach the whole API: a token gets permissions of the owner's role (and the training scopes) and the clubs it may use, never more than the role allows; presets for a website, a training app and read-only use; `GET /api/token` tells a script what its token may do; at most 600 requests per 5 minutes per token (`STANDBLATT_API_RATE_LIMIT`). Existing tokens stop working and need to be created again (by @HB9HIL)
- The version in the sidebar links to the GitHub repository (by @HB9HIL)

### Removed
- The redirect of the old setup page links (`setup?tab=`) to the configuration pages (by @HB9HIL)

### Fixed
- Uploads larger than 8 KB (shooter and result imports, images, form scans) failed on pods with a read-only file system; they are now kept in memory (by @HB9HIL)
- `make dev` seeded no test event since events got several ranges (by @HB9HIL)

## [0.19.1] - 2026-10-18

### Added
- `make clean` stops the containers and removes installed packages and build output (by @HB9HIL)
- staticcheck runs in `make prepare` and CI (pinned as a Go tool in `go.mod`) (by @HB9HIL)

### Changed
- Go tests copy a migrated template database and hash passwords cheaply, so they run in about 15 instead of 90 seconds (by @HB9HIL)
- Print forms load the scan of the empty form only for the preview, not for every printout (by @HB9HIL)
- The getting-started checklist no longer asks for a range once the club has an event, since events without lanes need none; the help group *Setup* is called *Configuration* like the app (by @HB9HIL)

### Fixed
- A control strip or result sheet that failed to render was answered as an empty PDF instead of an error (by @HB9HIL)
- Print forms of an entry on an event with several ranges printed on the first range instead of the shooter's (by @HB9HIL)
- Certificates of a classification left the club's SSV number empty (by @HB9HIL)
- The getting-started checklist and the shooter searches of the dashboard and the result entry shared cache entries with the shooters page; a programme page now follows live changes (by @HB9HIL)

## [0.19.0] - 2026-10-18

### Added
- **Attempts**: a programme can be shot several times per event, a fixed number or unlimited (further attempts at a shooting festival, an annual programme); the best, first or last attempt counts in rankings, club competition, special prizes and My training, further attempts can have their own price, awards go to every attempt or only the counting one, and the registration takes several attempts at once (by @HB9HIL)
- **Classifications**: rankings over several programmes or attempts of an event, e.g. a section championship as a sum of programmes or a club championship over the best three attempts; only the best programmes can count, incomplete shooters can be left unranked, the tie-break order is configurable (series, single shots, mouches, age) and classes rank age categories separately; shown on Rankings, kiosk displays and as a PDF/CSV report, copied for next year (by @HB9HIL)
- **Event templates**: a club saves an event as a template (modules, settings, ranges, programmes, classifications) and starts new events from it; templates take no registrations and are edited like events (by @HB9HIL)
- **Registration documents and certificates**: print forms that print a registration (e.g. a participant card) get a button in the registration; forms that print a rank make certificates from every ranking and classification, for a range of ranks and one class, as one PDF or on their own printer (by @HB9HIL)
- **Wreath card and control strip as print forms**: both print with built-in print forms that a club can copy and change (fields, logo); an event chooses its forms in the modules *Awards* and *Control strip*. Print forms have a data source (results of a programme, a registration, a rank) that decides their fields, and text can be italic or wrap (by @HB9HIL)
- **Images in print forms**: each club has an image library (PNG or JPEG) for logos and letterheads; the form editor places images like any other element, and exported forms carry their images (by @HB9HIL)

### Changed
- **Configuration**: an event's set-up has its own area, separate from the pages used on the day: *General*, *Modules*, *Programmes*, the pages of switched-on modules (reservations with the online booking, classifications, kiosk displays, helper links) and *Copy & close*, with their own sidebar; the operation sidebar shows only what the event uses. The **Modules** page switches building blocks on and off, says why one is locked and holds their small settings (award value on wreath cards, club excluded from the club competition). Staff no longer see the configuration; old setup links lead to the new pages (by @HB9HIL)
- **Club modules**: clubs switch their own building blocks on and off under *Club modules*, starting with print forms; the choice of a programme's print form is the module's section on the programme page (by @HB9HIL)
- **Programme page**: a programme opens as a page instead of a side sheet, with sections for the programme, series, prices and the switched-on modules (attempts, awards, club competition, kiosk); the section *Use* (active, range, combination, print form) and the special prizes save at once, also while shooters are registered. The programme list only shows the programmes with their status (by @HB9HIL)
- A new event starts from a kind (shooting festival, annual programme or training, Feldschiessen, club championship or a custom selection) that switches on the matching modules (by @HB9HIL)
- **Modules**: every event switches its building blocks on and off (lane assignment, so an event can use a range and its print agent without lanes, reservations with the online booking, lane sheet, control strip, ammunition desk, attempts, awards, special prizes, club competition, classifications, kiosk displays, helper links); the API refuses the operations, documents and reports of a module that is off, and a module that holds data (e.g. kiosk links) cannot be switched off. Existing events keep everything they had (by @HB9HIL)
- API: `EventSettings` lost `noLaneSheet`, `noControlStrip` and `ammoIssue`; events have `modules` and `modulesInUse` instead, `create-event` without `modules` switches on all but the ammunition desk and `GET /api/modules` lists them (by @HB9HIL)

### Fixed
- The README is short again and points to the help and the docs; the unused setting `STANDBLATT_PRICE_LABEL` is gone from the docs and `.env.example` (by @HB9HIL)
- The dashboard of an event without lane assignment or range showed registrations as waiting with *Assign lane*; they are listed as *Registered* with the time they were made and *Results*, *Edit* and *Finish* (by @HB9HIL)
- Copying an event that was itself a copy failed when the first copy had no booking days (by @HB9HIL)
- The German, French and Italian help of *Enter results* showed a paragraph of the CSV import instead of how results are saved (by @HB9HIL)

## [0.18.0] - 2026-10-18

### Added
- Events without a range, e.g. an annual programme whose results are entered afterwards: no lanes, reservations or online booking (by @HB9HIL)
- Lane sheet and control strip can be switched off per event (by @HB9HIL)
- Series can record single shots: the series is their sum, and the total alone can still be entered (by @HB9HIL)
- **Enter results**: results from paper or the scoring system shooter by shooter, without the desk; registrations are created as needed (by @HB9HIL)
- Results can be imported from a CSV file, e.g. exported from the scoring system; unknown shooters are created (by @HB9HIL)
- **Print forms**: results, down to single shots, printed onto pre-printed sheets such as the army's lane sheets; an editor modelled on Wavelog's QSL card designer (fields palette, rows of shot boxes, snapping with guides, multiple selection, align and distribute, undo, rulers and zoom, autosave every minute that can be switched off) over a scan of the empty form, with repeats, a printer offset and an exact PDF preview, export and import; a programme chooses its form and scorer links can print it at the lane (by @HB9HIL)

### Fixed
- The event line in the sidebar no longer ends in a dangling separator (by @HB9HIL)

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
