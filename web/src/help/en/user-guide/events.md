# Events

**Events** lists the events of the club in the groups *Today* (between the first and the last shooting day), *Upcoming*, *Past* and *Closed*; admins create new ones there. Opening an event shows its dashboard and the other pages of the operation in the sidebar; *Change event* at the top of the sidebar returns to the list.

![Events](../images/screenshots/events.png){ .screenshot }

## Configuration and operation {#configuration}

An event has two areas. The **operation** is what you work with on the day: *Dashboard*, *Enter results*, *Rankings*, *Statistics* and the pages of switched-on modules such as *Reservations*. **Configuration** at the end of the sidebar (admins only) opens the set-up of the event; the sidebar then lists its pages, and *Back to the event* returns to the operation:

| Page | What it holds |
|---|---|
| **General** | name, days, ranges, organizer and language (see below) |
| **Modules** | the building blocks the event uses (see below) |
| **Programmes** | what can be shot (see [Programmes](programmes.md)) |
| Pages of modules | e.g. **Reservations** (with the online booking), **Classifications**, **Kiosk displays**, **Helper links**, only while the module is on |
| **Copy & close** | copy for next year, close, reopen and delete |

Staff work only in the operation; the configuration is for admins.

## Settings {#settings}

**General** holds what every event has:

| Setting | Meaning |
|---|---|
| Name, Date, Ranges | the date is the main day; its year sets the age categories (SSV: competition year minus year of birth, the day itself does not matter, which counts for pre-shooting days in December before a January event); tick every range the event is shot on, or none for an event without lanes (see below) |
| Shooting days | single days with optional shooting times, for pre-shooting days and multi-day events (registrations count on the day they are made); the main day is always one |
| Period (season) | for events over weeks or months, e.g. the Feldschiessen from the opening of the range until 31 August: first and last day, plus the weekdays with their shooting times; no weekday ticked means every day. Single days, e.g. the official weekend, add to the period or replace its times on that date |
| Organizer | printed on documents |
| Language of documents and kiosk | language of PDFs, CSV exports and kiosk screens, independent of each user's interface language |

## Modules {#modules}

Not every event needs everything: a shooting festival books lanes and hands out wreath cards, an annual programme only collects results. **Modules** switches the building blocks of the event on and off; what is off disappears from the sidebar, the dashboard, the programme editor, the rankings and the reports.

![The modules of an event](../images/screenshots/modules.png){ .screenshot }

| Module | What it adds |
|---|---|
| Lane assignment | shooters are registered on lanes: the lane grid on the dashboard, the lane status on kiosk displays and scorer links (needs a range; see below) |
| Reservations | lanes reserved for groups and clubs; its page in the configuration switches on the [online booking](online-booking.md), through which clubs ask for lanes with a public link (needs a range) |
| Lane sheet, Control strip | the documents printed at registration; the control strip can print with a [print form](print-forms.md#builtin) of the club |
| Ammunition desk | the **Ammo** button on the dashboard and for the ammunition desk (see [Helpers](helpers.md#ammunition)); leave it off where shooters bring their own, e.g. air pistol and air rifle |
| Attempts | programmes that can be shot several times (see [Programmes](programmes.md#attempts)) |
| Awards | wreath cards, prizes in kind and cash prizes; its setting *Print the award value on wreath cards* puts the value in figures and words on the card, *Print form of the wreath cards* chooses a [print form](print-forms.md#builtin) of the club |
| Special prizes | the best ranks of every age category (switches on Awards, see [Programmes](programmes.md#special-prizes)) |
| Club competition | clubs ranked by their shooters' results; its setting *Club excluded from the club competition* is usually the host club, listed as *out of competition* |
| Classifications | rankings over several programmes or attempts (see [Classifications](classifications.md)) |
| Kiosk displays | rankings on screens in the hall (see [Rankings and kiosk](rankings-kiosk.md)) |
| Helper links | scorers, ammunition desk and results clerk without an account (see [Helpers](helpers.md)) |

When you create an event, choose what kind it is; this only switches modules on:

| Kind | Modules |
|---|---|
| Shooting festival | Lane assignment, Reservations, Lane sheet, Control strip, Awards, Special prizes, Club competition, Kiosk displays, Helper links |
| Annual programme or training | Attempts, Classifications |
| Feldschiessen | Lane assignment, Awards, Club competition, Kiosk displays, Helper links |
| Club championship | Attempts, Classifications, Kiosk displays |
| Custom selection | none |

A module that holds data stays on: switch off the kiosk only after deleting its kiosk links, the awards only when no programme has award tiers, and so on; the card says why its switch is locked; Reservations stay on while the online booking link is active. A module another one needs (Awards for Special prizes) stays on as long as that one is.

Shooting times also give the online booking its slots, unless the booking sets other hours for a day (see [Online booking](online-booking.md)).

### Club modules {#club-modules}

Some modules belong to the club rather than to an event. **Club modules** in the club part of the sidebar (admins) switches them on and off:

| Module | What it adds |
|---|---|
| Print forms | forms for pre-printed sheets and the *Print form* section of every programme (see [Print forms](print-forms.md)) |

## Events without lanes {#no-lanes}

An event without lanes, e.g. an annual programme or a Feldschiessen whose results come on paper or from the scoring system, either has no range or switches the module *Lane assignment* off. Then there is no lane grid, no reservations and no scorer links; shooters are registered without a lane and their results are entered as usual. With a range but without lane assignment, documents still print on the range's print agent and programmes can still be tied to the range; without a range, documents open as PDF, because print agents belong to a range.

## Copy for next year {#copy}

On **Copy & close** in the configuration, **Copy for next year** creates a new event with the same ranges, organizer, language, settings, modules and all programmes, including prices, awards, combinations and special prizes, and the [classifications](classifications.md); the shooting days and the period move with the date, weekly shooting times stay on their weekdays. Registrations, results and kiosk links are not copied. The dialog suggests the name with the year (and a leading number like *42.*) counted up and the date one year later; adjust both as needed.

![Copy an event](../images/screenshots/copy-event.png){ .screenshot }

## Templates {#templates}

An event that is set up well can become a **template** of the club: **Save as template** on **Copy & close** stores a copy with its modules, settings, ranges, programmes and classifications, named after the event without its year. Templates are listed under *Templates* below the events; a template opens in its configuration and can be changed like an event, but it takes no registrations. **New event** offers the club's templates below the built-in kinds: choose one, then give the new event its name and date. **Delete** on the template's **Copy & close** removes it.

![New event with the kinds and a template of the club](../images/screenshots/new-event.png){ .screenshot }

## Close {#close}

When an event is over and the rankings are checked, **Close** it on **Copy & close** in the configuration. Closing releases the rankings: shooters with a [personal training log](../my-training/account.md#competitions) now find their results of this event, with their rank in the programme, and can add them to their log. Everything stays visible (dashboard, results, rankings, reports, documents), but registrations, results, awards, programmes and settings can no longer be changed. A banner on every page of the event says so. **Reopen** makes it editable again; results that shooters already added to their logs stay there. Copying a closed event for next year works.

## Delete

**Delete** (on **Copy & close**) removes an event, but only while nobody is registered for it.
