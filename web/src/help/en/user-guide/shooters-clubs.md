# Shooters and clubs

Each club keeps its own shooter master data; nothing is shared with other clubs. **Shooters** and **Clubs** are under **Master data** in the sidebar.

## Shooters

![Shooters](../images/screenshots/shooters.png){ .screenshot }

- **Search** works like on the dashboard: license number (leading zeros optional), part of a name, or first and last name in either order. The list shows 50 shooters per page, the club list 50 clubs; the arrows below turn the pages.
- Clicking a shooter opens the details: name, license number, birthday, address and the clubs the shooter belongs to.
- **Club** (new shooter, or adding a club in the details) suggests the clubs of the list. A name that matches a club regardless of upper case and spaces is that club. Any other name shows *is not in the club list yet* with similar clubs under *Did you mean*; a new club is only created after **Create new club**, so a typo does not become a second club.
- **Import** (admins) reads a CSV file, see [CSV import](../reference/csv-import.md).
- Shooters with registrations cannot be deleted.

## Clubs

![Clubs](../images/screenshots/clubs.png){ .screenshot }

**Clubs** are the clubs shooters belong to, with their SSV number. They appear in registrations, rankings and the club competition. A shooter can belong to several clubs and registers for one of them. Clubs used by registrations cannot be deleted, and two clubs cannot have the same name regardless of upper case and spaces.

**Merge into…** (admins) cleans up a duplicate, e.g. *S.G. Hiltown* next to *SG Hilltown*: choose the club that stays (clubs with a similar name come first). It keeps its name and gets the shooters, registrations and booking requests of the other one, including those of closed events, so their rankings show the club that stays; the SSV number is taken over when it has none. The other club is deleted.
