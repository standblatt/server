# Events

**Events** lists the events of the club in the groups *Today* (between the first and the last shooting day), *Upcoming*, *Past* and *Closed*; roles with the permission *Configure events* create new ones there. Opening an event shows its dashboard and the other pages of the operation in the sidebar; *Change event* at the top of the sidebar returns to the list.

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

The configuration needs the permission *Configure events*, the operation *Run events* (see [Roles](users.md#roles)). Why the two are separate: [Configuration and operation](../explanation/modules.md#configuration).

## Settings {#settings}

**General** holds what every event has:

| Setting | Meaning |
|---|---|
| Name, Date, Ranges | the date is the main day; its year sets the age categories (SSV: competition year minus year of birth, the day itself does not matter, which counts for pre-shooting days in December before a January event); tick every range the event is shot on, or none for an event without lanes (see below) |
| Shooting days | single days with optional shooting times, for pre-shooting days and multi-day events (registrations count on the day they are made); the main day is always one |
| Period (season) | for events over weeks or months, e.g. the Feldschiessen from the opening of the range until 31 August: first and last day, plus the weekdays with their shooting times; no weekday ticked means every day. Single days, e.g. the official weekend, add to the period or replace its times on that date |
| Organizer | printed on documents |
| Language of documents and kiosk | language of PDFs, CSV exports and kiosk screens, independent of each user's interface language |

Shooting times also give the online booking its slots, unless the booking sets other hours for a day (see [Online booking](online-booking.md)).

Step by step for a season: [Run an annual programme over the season](../how-to/annual-programme.md), [Set up a Feldschiessen over a period](../how-to/feldschiessen.md).

## Modules {#modules}

**Modules** switches the building blocks of the event on and off; why events differ only by their modules is explained in [Modules and kinds of events](../explanation/modules.md).

![The modules of an event](../images/screenshots/modules.png){ .screenshot }

Every module is a card with its switch; the settings of a module, e.g. *Print the award value on wreath cards* of *Awards*, sit on its card. What every module adds, its settings and when it stays on: [Modules](../reference/modules.md#event-modules).

When you create an event, choose what kind it is; this only switches modules on (see [Event kinds](../reference/modules.md#kinds)).

A module that holds data stays on, and the card says why its switch is locked (see [Why a module that holds data stays on](../explanation/modules.md#in-use)). A module another one needs, e.g. *Awards* for *Special prizes* or *Lane assignment* for *Reservations*, stays on as long as that one is.

### Club modules {#club-modules}

Some modules belong to the club rather than to an event. **Club modules** in the club part of the sidebar (permission *Users and club*) switches them on and off; for now that is *Print forms* (see [Print forms](print-forms.md)). What each adds and when it stays on: [Club modules](../reference/modules.md#club-modules); why they belong to the club: [Club modules and event modules](../explanation/modules.md#club-modules).

## Events without lanes {#no-lanes}

An event without lanes has no range, or switches *Lane assignment* off; what each choice means is explained in [Events with and without lanes](../explanation/modules.md#lanes).

## Copy for next year {#copy}

On **Copy & close** in the configuration, **Copy** under *Copy for next year* creates a new event with the same ranges, organizer, language, settings, modules and all programmes, including prices, awards, combinations and special prizes, and the [classifications](classifications.md); the shooting days and the period move with the date, weekly shooting times stay on their weekdays. Registrations, results, kiosk links, helper links and the booking link are not copied. The dialog *Copy event* suggests the name and the date one year later; **Create copy** creates the event. Step by step: [Take over last year's event](../how-to/next-year.md#copy).

![Copy an event](../images/screenshots/copy-event.png){ .screenshot }

## Templates {#templates}

An event that is set up well can become a **template** of the club: **Save as template** on **Copy & close** stores a copy with its modules, settings, ranges, programmes and classifications, named after the event without its year. Templates are listed under *Templates* below the events; a template opens in its configuration and can be changed like an event, but it takes no registrations. **New event** lists the club's templates below the kinds (see [Start from a template of the club](../how-to/next-year.md#template)). **Delete** on the template's **Copy & close** removes it.

![New event with the kinds and a template of the club](../images/screenshots/new-event.png){ .screenshot }

## Close {#close}

When an event is over and the rankings are checked, **Close event** on **Copy & close** in the configuration closes it; a banner on every page of the event then says so, and **Reopen** makes it editable again. What closing means for the rankings, My training and changes: [Closing an event](../explanation/closing.md).

## Delete

**Delete** (on **Copy & close**) removes an event, but only while nobody is registered for it. A closed event cannot be deleted; reopen it first.
