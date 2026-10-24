# Concepts

This page explains the words Standblatt uses and how they fit together. The German, French and Italian terms are in the [Glossary](../reference/glossary.md); the steps for each task are in the user guide.

In short: a **club** has **ranges**, **shooters** and **users**. It runs **events**; an event offers **programmes** and switches on the **modules** it needs. A **registration** puts a shooter on the event with one or more programmes, each programme becomes an **entry**, the entry gets **results**, and the results give **awards** and a place in the **rankings**. When the event is **closed**, the rankings are final and reach the shooters' personal training logs.

## Clubs, users and roles {#clubs-users-and-roles}

A **club** is a tenant: everything in Standblatt belongs to exactly one club, and clubs never see each other's data. The platform operator creates clubs and user accounts in the admin panel.

**Users** are personal accounts. One account can belong to several clubs, with one role per club. A role is a set of permissions such as *Run events* or *Configure events*; every club starts with the roles *Admin*, *Board*, *Helper* and *Member* and can change them or add its own (see [Roles](../user-guide/users.md#roles)). Clubs invite people by email; the invited person joins by accepting the link.

Why accounts are global while all club data stays in the club is explained in [Clubs, accounts and data](../explanation/data.md).

## Ranges and lanes {#ranges-and-lanes}

A **range** (German *Stand*) is a shooting range of the club with a number of **lanes** (*Scheiben*). Lanes can be marked out of service. Printers and the print agent belong to a range.

Each distance is a range of its own: a hall with 25 m and 50 m lanes is two ranges, each with lanes from 1. An event can be shot on several ranges, on one, or on none at all (see [Events with and without lanes](../explanation/modules.md#lanes)).

## Events and programmes {#events-and-programmes}

An **event** (*Anlass*) is a competition on one or more ranges on one or more days (the date is the main day; further days cover pre-shooting and multi-day events). It has a language, which is used for printed documents and kiosk screens, and [modules](../user-guide/events.md#modules) that switch on what it needs, such as lane assignment, awards or the club competition.

There are no event types: a shooting festival, a Feldschiessen and an annual programme are the same kind of event with different modules (see [Modules and kinds of events](../explanation/modules.md)).

An event offers **programmes** (*Stiche*): each programme has **series** (*Passen*) with a number of shots and maximum points, **prices** per age category and **award tiers** (wreath cards, prizes in kind or cash prizes) with a minimum total per age category. Programmes with the same name are **variants** (for example *Sport* and *Supported*); a shooter shoots one variant. A programme can be **combinable** with another programme: it then takes some of its series from that programme instead of being shot separately.

A programme can allow several **attempts**, of which one counts. A **classification** ranks shooters over several programmes or attempts. How programme rankings, attempts, classifications, the club competition and special prizes relate is explained in [How rankings fit together](../explanation/rankings.md).

## Shooters, registrations and entries {#shooters-registrations-and-entries}

**Shooters** are the people who shoot. Each club keeps its own shooter master data, together with the **clubs** shooters belong to (the shooter's own club, not to be confused with the tenant club).

A **registration** puts a shooter on the event, optionally on a lane, with one or more programmes. Each programme of a registration is an **entry**. A registration is *registered* (waiting), *shooting* (on a lane) or *completed* (all programmes have results; the lane is free again).

## Results, awards and rankings {#results-awards-and-rankings}

**Results** are the points per series, optionally with mouches (inner tens). Once all series are in, the shooter may receive one of the **awards** the total reaches in the shooter's age category. **Rankings** order all results of a programme by the [SSV rules](../reference/rules.md#ranking).

An award depends only on the shooter's own total; a rank depends on everybody else's. Rankings are live while the event runs and become final when it is [closed](../explanation/closing.md).

## Age categories {#age-categories}

The age category follows the SSV rule: the year of the event minus the year of birth. The day of birth does not matter, and the year is that of the event's main day, so a pre-shooting day in December counts for a January event. The category is fixed when the shooter registers and decides the prices, awards and special prizes the shooter can get. A shooter without a birthday counts as Elite. The categories and their ages are listed in [Competition rules](../reference/rules.md#age-categories).

## Documents and printing {#documents-and-printing}

Every document is a PDF rendered by the server in the event's language. At the range, a small print agent sends it to the right printer; without one, the PDF opens in the browser. See [How printing works](../explanation/printing.md).
