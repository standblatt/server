# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
