# Events

**Events** lists the events of the club in the groups *Today* (between the first and the last shooting day), *Upcoming*, *Past* and *Closed*; admins create new ones there. Opening an event shows its dashboard, reservations, rankings, statistics and setup in the sidebar; *Change event* at the top of the sidebar returns to the list.

![Events](../images/screenshots/events.png){ .screenshot }

## Settings {#settings}

The event's **Setup** has five tabs: *Programmes* (see [Programmes](programmes.md)), *Event*, *Kiosk displays* (see [Rankings and kiosk](rankings-kiosk.md)), *Helpers* (see [Helpers](helpers.md)) and *Copy & close*. Staff only see the programmes.

The *Event* tab holds the settings (admins):

| Setting | Meaning |
|---|---|
| Name, Date, Ranges | the date is the main day; its year sets the age categories (SSV: competition year minus year of birth, the day itself does not matter, which counts for pre-shooting days in December before a January event); tick every range the event is shot on, or none for an event without lanes (see below) |
| Shooting days | single days with optional shooting times, for pre-shooting days and multi-day events (registrations count on the day they are made); the main day is always one |
| Period (season) | for events over weeks or months, e.g. the Feldschiessen from the opening of the range until 31 August: first and last day, plus the weekdays with their shooting times; no weekday ticked means every day. Single days, e.g. the official weekend, add to the period or replace its times on that date |
| Organizer | printed on documents |
| Language of documents and kiosk | language of PDFs, CSV exports and kiosk screens, independent of each user's interface language |
| Club excluded from the club competition | usually the host club; listed as *out of competition* |
| Documents at registration | **Lane sheet** and **Control strip**: switch off what the event does not use; the registration then neither offers nor prints it |
| Show lane status on the rankings | kiosk screens show the lanes |
| Print the award value on wreath cards | value in figures and words on the card |
| The organizer hands out the ammunition | shows the **Ammo** button on the dashboard and for the ammunition desk (see [Helpers](helpers.md)); off for events where shooters bring their own, e.g. air pistol and air rifle |

Shooting times also give the online booking its slots, unless the booking sets other hours for a day (see [Online booking](online-booking.md)).

## Events without lanes {#no-lanes}

An event needs no range, e.g. an annual programme or a Feldschiessen whose results come on paper or from the scoring system. Without a range there is no lane grid, no reservations and no online booking; shooters are registered without a lane and their results are entered as usual. Documents can still be opened as PDF, but they do not go to a print agent, because print agents belong to a range.

## Copy for next year

On the *Copy & close* tab, **Copy for next year** creates a new event with the same ranges, organizer, language, settings and all programmes, including prices, awards, combinations and special prizes, and the [classifications](classifications.md); the shooting days and the period move with the date, weekly shooting times stay on their weekdays. Registrations, results and kiosk links are not copied. The dialog suggests the name with the year (and a leading number like *42.*) counted up and the date one year later; adjust both as needed.

![Copy an event](../images/screenshots/copy-event.png){ .screenshot }

## Close {#close}

When an event is over and the rankings are checked, **Close** it on the *Copy & close* tab. Closing releases the rankings: shooters with a [personal training log](../my-training/account.md#competitions) now find their results of this event, with their rank in the programme, and can add them to their log. Everything stays visible (dashboard, results, rankings, reports, documents), but registrations, results, awards, programmes and settings can no longer be changed. A banner on every page of the event says so. **Reopen** makes it editable again; results that shooters already added to their logs stay there. Copying a closed event for next year works.

## Delete

**Delete** (on the *Copy & close* tab) removes an event, but only while nobody is registered for it.
