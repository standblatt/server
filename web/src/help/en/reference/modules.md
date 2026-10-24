# Modules

Modules are the building blocks an event or a club switches on. What is off disappears from the sidebar, the dashboard, the programme editor, the rankings and the reports. How to switch them: [Modules](../user-guide/events.md#modules) and [Club modules](../user-guide/events.md#club-modules).

## Event modules {#event-modules}

Under **Configuration** → **Modules** of an event.

| Module | What it adds |
|---|---|
| Lane assignment | shooters are registered on lanes: the lane grid on the dashboard, the lane status on kiosk displays and scorer links (needs a range; see below) |
| Reservations | lanes reserved for groups and clubs; its page in the configuration switches on the [online booking](../user-guide/online-booking.md), through which clubs ask for lanes with a public link (needs a range) |
| Lane sheet | the lane sheet printed at registration: shooter, lane, programmes and barcodes |
| Control strip | the control strip printed at registration: the passes the shooter may shoot |
| Ammunition desk | the **Ammo** button on the dashboard and for the ammunition desk (see [Helpers](../user-guide/helpers.md#ammunition)); leave it off where shooters bring their own, e.g. air pistol and air rifle |
| Attempts | programmes that can be shot several times (see [Programmes](../user-guide/programmes.md#attempts)) |
| Awards | wreath cards, prizes in kind and cash prizes by points per age category |
| Special prizes | the best ranks of every age category (switches on Awards, see [Programmes](../user-guide/programmes.md#special-prizes)) |
| Club competition | clubs ranked by their shooters' results |
| Classifications | rankings over several programmes or attempts (see [Classifications](../user-guide/classifications.md)) |
| Kiosk displays | rankings on screens in the hall (see [Rankings and kiosk](../user-guide/rankings-kiosk.md)) |
| Helper links | scorers, ammunition desk and results clerk without an account (see [Helpers](../user-guide/helpers.md)) |

### Where a module shows {#places}

| Module | Pages and sections |
|---|---|
| Lane assignment | lane grid on the dashboard, lanes in the registration, lane status on kiosk displays, role *Scorer* of helper links |
| Reservations | page **Reservations** of the event, with the booking requests and their count in the sidebar; page **Reservations** in the configuration (online booking); report *Lane bookings* |
| Lane sheet | *Print lane sheet* in the registration |
| Control strip | *Print control strip* in the registration |
| Ammunition desk | **Ammo** button and *ammo* badge on the dashboard and on helper links |
| Attempts | section **Attempts** of the programme page; counter per programme in the registration |
| Awards | section **Awards** of the programme page; awards in the results; figures on **Statistics**; reports *Awards* and *Awards to hand out* |
| Special prizes | section **Special prizes** of the programme page; report *Special prizes* |
| Club competition | *Counts for the club competition* on the programme page; club competition on **Statistics**; report *Club competition* |
| Classifications | page **Classifications** in the configuration; tables on **Rankings** and kiosk displays; report *Classifications* |
| Kiosk displays | page **Kiosk displays** in the configuration; *Show on kiosk displays* on the programme page |
| Helper links | page **Helper links** in the configuration |

### Rules {#rules}

| Module | Requires | Needs a range | Settings | Cannot be switched off while |
|---|---|---|---|---|
| Lane assignment | – | yes | – | a registration has a lane |
| Reservations | Lane assignment | yes | the online booking on its configuration page | the online booking link is active |
| Lane sheet | – | – | – | – |
| Control strip | – | – | *Print form of the control strips* | – |
| Ammunition desk | – | – | – | – |
| Attempts | – | – | – | a programme allows more than one attempt |
| Awards | – | – | *Print the award value on wreath cards* (on for a new event), *Print form of the wreath cards* | a programme has award tiers |
| Special prizes | Awards | – | – | a programme awards special prizes |
| Club competition | – | – | *Club excluded from the club competition* | – |
| Classifications | – | – | – | the event has a classification |
| Kiosk displays | – | – | *Show lane status on the rankings* (with Lane assignment) | the event has a kiosk link |
| Helper links | – | – | – | the event has a helper link |

- A module that another one requires stays on as long as that one is on.
- Modules that need a range are off in an event without a range, whatever is chosen (see [Events without lanes](../user-guide/events.md#no-lanes)).
- The settings sit on the module's card under **Modules**; *Show lane status on the rankings* is on the page **Kiosk displays**.
- **Copy for next year** and templates take over the chosen modules and the modules in use.
- An event created through the API without a list of modules gets every event module except *Ammunition desk*.

## Club modules {#club-modules}

Under **Club modules** in the club part of the sidebar (permission *Users and club*).

| Module | What it adds | Cannot be switched off while |
|---|---|---|
| Print forms | forms for pre-printed sheets and the *Print form* section of every programme (see [Print forms](../user-guide/print-forms.md)); page **Print forms** and the permission *Print forms* | the club has a print form or an image |

A new club starts with *Print forms* on. With it off, events print the built-in wreath card and control strip.

## Event kinds {#kinds}

**New event** asks for the kind; it only switches modules on, everything can be changed later.

| Kind | Modules |
|---|---|
| Shooting festival | Lane assignment, Reservations, Lane sheet, Control strip, Awards, Special prizes, Club competition, Kiosk displays, Helper links |
| Annual programme or training | Attempts, Classifications |
| Feldschiessen | Lane assignment, Awards, Club competition, Kiosk displays, Helper links |
| Club championship | Attempts, Classifications, Kiosk displays |
| Custom selection | none |

A [template](../user-guide/events.md#templates) of the club brings its own modules instead.
