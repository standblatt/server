# Modules and kinds of events

Clubs run very different events: a shooting festival with a desk, lane assignment and wreath cards, a Feldschiessen on the official weekend, an annual programme shot over the season, a club championship. Standblatt has no separate code for each of them. Every event is the same kind of thing; events differ only in their programmes, their ranges and days, and the **modules** they switch on.

## The core and the modules {#core}

Every event has the same core: name, days, ranges and language, programmes with their series and prices, registrations, results, rankings, statistics and closing. Everything beyond that is a module: lane assignment, reservations with the online booking, the lane sheet and the control strip, the ammunition desk, attempts, awards, special prizes, the club competition, classifications, kiosk displays and helper links.

A module that is off is gone, not only hidden: it disappears from the sidebar, the dashboard, the programme editor, the rankings and the reports, and the server refuses its functions too, so neither a helper link nor an API token can use it. An annual programme without awards therefore has no award section in its programmes and no award reports; nobody has to skip past what the event does not use.

Some modules need others: *Special prizes* switch on *Awards*, because a special prize is an award, and *Reservations* need *Lane assignment*, because a reservation blocks lanes. A module stays on as long as a module that needs it is on.

## Kinds of events are only a starting point {#kinds}

**New event** asks what kind of event it is: shooting festival, annual programme or training, Feldschiessen, club championship or a custom selection. The kind only switches on a set of modules (the list is under [Event kinds](../reference/modules.md#kinds)); the event does not remember it. Afterwards, the modules can be changed like those of any other event, so a Feldschiessen that also wants classifications switches them on, and a shooting festival without a desk switches lane assignment off.

This is why a new kind of event needs no new feature: if a combination of modules covers it, it can be run today. A club that holds an event every year in the same form saves it as a [template](../user-guide/events.md#templates) instead; a template is an event with modules, programmes and classifications, not a kind.

## Why a module that holds data stays on {#in-use}

Switching off a module that holds data would hide data that still counts: award tiers that shooters have reached, kiosk links that still show the rankings in the hall, programmes that allow several attempts. Standblatt therefore keeps such a module on until its data is gone, and the module card says why its switch is locked:

| Module | Stays on while |
|---|---|
| Lane assignment | a registration has a lane |
| Reservations | the online booking link is active |
| Attempts | a programme allows more than one attempt |
| Awards | a programme has award tiers |
| Special prizes | a programme awards special prizes |
| Classifications | the event has a classification |
| Kiosk displays | the event has a kiosk link |
| Helper links | the event has a helper link |

The lane sheet, the control strip, the ammunition desk and the club competition hold no data of their own and can be switched at any time.

## Configuration and operation {#configuration}

An event is prepared once and run many times: the programmes, prices and modules are set up before the first shooter arrives, while the desk works on the day. Standblatt keeps the two apart. The **operation** (dashboard, results, rankings, statistics and the pages of modules such as reservations) needs the permission *Run events*; the **configuration** (general settings, modules, programmes, the set-up pages of modules, copying and closing) needs *Configure events*. A helper at the desk cannot change a price by accident, and the configuration is not in the way of the desk.

## Club modules and event modules {#club-modules}

Most modules belong to an event, because each event decides whether it hands out awards or uses kiosk screens. Some building blocks are used by every event of a club and therefore belong to the club: the **print forms** with their image library are drawn once and used by the programmes of every event. Such **club modules** are switched under [Club modules](../user-guide/events.md#club-modules). A new club has print forms switched on; like event modules, the module stays on while the club has forms or images.

## Events with and without lanes {#lanes}

Lanes matter only where a desk assigns them. Standblatt therefore separates the range from the lanes:

| The event has | What follows |
|---|---|
| a range and *Lane assignment* | the desk of a shooting festival: lane grid, reservations, online booking, scorer links and lane status on kiosk screens |
| a range, but no *Lane assignment* | shooters are registered without a lane; documents still go to the range's print agent and programmes can still be tied to the range |
| no range | no lanes and nothing that needs them; documents open as PDF, because print agents belong to a range |

An annual programme or a Feldschiessen whose results come on paper or from the scoring system needs no lanes, but the club may still want its printer at the range; that is the middle row. Modules that need lanes are switched off by themselves when the event has no range.

Results reach an event without lanes the same way as any other: through **Enter results**, a CSV file or the results of a registration (see [Enter results](../user-guide/results-entry.md)).

## Copying keeps the shape {#copy}

**Copy for next year** and templates take over the modules that are switched on and those that hold data, together with the programmes and classifications. An event copied from last year therefore starts with the same shape, without its registrations and results.
