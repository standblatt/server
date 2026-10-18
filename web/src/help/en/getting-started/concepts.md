# Concepts

## Clubs, users and roles {#clubs-users-and-roles}

A **club** is a tenant: everything in Standblatt belongs to exactly one club, and clubs never see each other's data. The platform operator creates clubs and user accounts in the [admin panel](../admin-guide/admin-panel.md).

**Users** are personal accounts. One account can belong to several clubs, with one role per club. A role is a set of permissions such as *Run events* or *Configure events*; every club starts with the roles *Admin*, *Board*, *Helper* and *Member* and can change them or add its own (see [Roles](../user-guide/users.md#roles)). Clubs invite people by email; the invited person joins by accepting the link.

## Ranges and lanes

A **range** (German *Stand*) is a shooting range of the club with a number of **lanes** (*Scheiben*). Lanes can be marked out of service. Printers and the print agent belong to a range.

## Events and programmes

An **event** (*Anlass*) is a competition on one or more ranges on one or more days (the date is the main day; further days cover pre-shooting and multi-day events). It has a language, which is used for printed documents and kiosk screens, and [modules](../user-guide/events.md#modules) that switch on what it needs, such as lane assignment, awards or the club competition.

An event offers **programmes** (*Stiche*): each programme has **series** (*Passen*) with a number of shots and maximum points, **prices** per age category and **award tiers** (wreath cards or prizes in kind) with a minimum total per age category. Programmes with the same name are **variants** (for example *Sport* and *Supported*); a shooter shoots one variant. A programme can be **combinable** with another programme: it then takes some of its series from that programme instead of being shot separately.

## Shooters, registrations and entries

**Shooters** are the people who shoot. Each club keeps its own shooter master data, together with the **clubs** shooters belong to (the shooter's own club, not to be confused with the tenant club).

A **registration** puts a shooter on the event, optionally on a lane, with one or more programmes. Each programme of a registration is an **entry**. A registration is *registered* (waiting), *shooting* (on a lane) or *completed* (all programmes have results; the lane is free again).

## Results, awards and rankings

**Results** are the points per series, optionally with mouches (inner tens). Once all series are in, the shooter may receive one of the **awards** the total reaches in the shooter's age category. **Rankings** order all results of a programme by the [SSV rules](../reference/rules.md#ranking).

## Age categories {#age-categories}

The age category follows from the birth year and the year of the event:

| Category | Short | Age in the event year |
|---|---|---|
| U10 | U10 | up to 9 |
| U13 | U13 | 10–12 |
| U15 | U15 | 13–14 |
| U17 | U17 | 15–16 |
| U19 | U19 | 17–18 |
| U21 | U21 | 19–20 |
| Elite | E | 21–45 |
| Senior | S | 46–59 |
| Veteran | V | 60–69 |
| Senior veteran | SV | 70 and older |

A shooter without a birthday counts as Elite.

## Glossary

| English (Standblatt) | German | French | Italian |
|---|---|---|---|
| Range | Stand | Stand | Stand |
| Lane | Scheibe | Cible | Bersaglio |
| Event | Anlass | Manifestation | Manifestazione |
| Programme | Stich / Programm | Programme | Programma |
| Series | Passe | Passe | Serie |
| Award | Auszeichnung | Distinction | Distinzione |
| Wreath card | Kranzkarte | Carte-couronne | Carta corona |
| Prize in kind | Naturalgabe | Prix en nature | Premio in natura |
| Club competition | Vereinskonkurrenz | Concours de sections | Concorso di sezione |
| Lane sheet | Standblatt | Feuille de stand | Foglio di stand |
| Control strip | Kontrollstreifen | Bande de contrôle | Striscia di controllo |
