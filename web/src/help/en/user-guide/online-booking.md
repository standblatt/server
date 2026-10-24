# Online booking

With **online booking**, external clubs ask for lanes themselves. The organizer puts a link into the invitation; a club picks the time slots it wants, says how many lanes it needs in each and announces its shooters. An admin confirms or declines every request.

## Switching it on {#setup}

The online booking is part of the module *Reservations* (it needs a range and *Lane assignment*): on **Reservations** in the configuration, **Switch on** is all it takes. The page then shows the link with a QR code for the invitation and the slots it offers.

![Online booking in the configuration](../images/screenshots/booking-setup.png){ .screenshot }

Without changes, every day of the event is open during its shooting times (see [Events](events.md#settings)), or from 08:00 to 17:00 when it has none, in slots of 30 minutes, and requests are accepted until 14 days before the first day. **Adjust** changes this:

- **Shooting from** and **Shooting until**: the hours of every day without shooting times.
- **Break from** and **Break until**: an optional break, e.g. for lunch.
- **Slot length (minutes)**: how long one slot is.
- **Lanes per slot and club**: a limit per request; empty means no limit.
- **Deadline**: the last day clubs can ask or change their request.
- **Days with other hours**: one line per day of the event, e.g. Friday 17:00–20:00. The days come from the event and follow it when you add or remove one.

**New link** replaces the link; the old one stops working. **Switch off** closes the booking; requests and confirmed reservations stay. *Copy for next year* takes the hours along, but not the link and not the deadline.

## What a club sees {#club}

The link opens a page without login. It shows the event, the deadline and every slot with the lanes still free on each range. Only confirmed reservations count; names of other clubs are never shown.

![The booking page of a club](../images/screenshots/booking-public.png){ .screenshot }

The club enters the number of lanes per slot, its name, a contact person, email and phone, and its **Shooters**: last name, first name, birthday and optionally the licence. **Paste list** takes rows copied from Excel (last name, first name, birthday, licence). The shooter list can follow later; programmes are chosen at the range as usual.

After **Send request**, the club gets a mail with a personal link. Through it, it can check the status, change the request or cancel it until the deadline. Changing the slots of a confirmed request releases its lanes and the request waits for a new confirmation; adding shooters keeps it confirmed.

## Confirming requests {#confirm}

Every admin with an email gets a mail for each new, changed or cancelled request. The requests are listed at the top of **Reservations**; the sidebar shows how many are waiting.

![Booking requests on the reservations page](../images/screenshots/booking-requests.png){ .screenshot }

- **Confirm** gives each slot free lanes of its range, neighbours where possible, as an ordinary reservation named after the club, and mails the lane numbers to the club. When a slot has too few free lanes, nothing is reserved and the message names the slot; move a reservation or decline.
- **Confirm** first asks which **Club** of the master data the request belongs to, because the club typed its name itself: the club with the same name is chosen already, clubs with a similar name are suggested, and *New club* creates it. The request then shows the club it was matched to.
- **Confirm** also takes the shooters into the master data, in that club: a shooter is found by licence, else by name and birthday, and created otherwise. Existing shooters are never overwritten, because the list comes from a public form; only a missing birthday or licence is filled in. Shooters added after the confirmation are taken over with **Take over shooters**.
- **Decline** asks for an optional reason, which goes to the club by mail. Lanes the request held are free again.

The reserved lanes are ordinary [reservations](reservations.md): change their lane numbers there if needed. At the event, the shooters are found with the normal search at the registration desk. The report *Lane bookings* under [Statistics and reports](statistics-reports.md) lists every club with its slots, lanes and shooters.
