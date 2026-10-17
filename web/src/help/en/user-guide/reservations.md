# Reservations

**Reservations** blocks lanes of the event's range for a time, for example for a training or a group.

![Reservations](../images/screenshots/reservations.png){ .screenshot }

- **Lanes**: a list like `5` or `1, 5, 7-10`.
- **From** and **To**: the time span. Reservations on the same lane may not overlap; back-to-back is fine.
- **Show as reserved (minutes before)**: lead time from which the lanes show as reserved on the dashboard and kiosk.

Reservations belong to the range, so they appear in every event on that range. Reserved lanes can still be assigned to a shooter after a confirmation. The server announces the start and end of reservations by itself, so screens change without reloading.

Requests that external clubs send through [online booking](online-booking.md) are listed above the reservations; once confirmed, they become ordinary reservations named after the club.
