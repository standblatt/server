# Clubs, accounts and data

Standblatt serves many clubs on one server. This page explains who owns which data, why people have one account for everything, and why the shooters of one club are not those of another.

## Clubs are tenants {#tenants}

A **club** in Standblatt is a tenant: its ranges, shooters, clubs list, events, programmes, print forms and users belong to it and to nobody else. Clubs never see each other's data, not even that a shooter or an event exists elsewhere. The separation is part of every query the server makes, not a filter of the interface.

## One account, a role per club {#accounts}

People are different from data: the same person can be the secretary of one club, help at the desk of another club's festival and keep a personal training log. Standblatt therefore has **global accounts**. One account logs in once and reaches every club it belongs to, with a role per club that decides what it may do there (see [Roles](../user-guide/users.md#roles)).

A club does not own the account, only the membership. Removing a member takes away the access to that club; the account, its training log and its other clubs stay. Clubs invite people by email, and the person becomes a member only by accepting the link, so nobody is added to a club without knowing.

Helpers at an event do not need an account at all: a [helper link](../user-guide/helpers.md) gives access to one event and one task. An [API token](../user-guide/login-language.md#api-tokens) acts for its owner and never with more than the owner's role in a club.

## Shooter master data per club {#shooters}

A **shooter** is not an account. It is a record in a club's master data: name, licence number, birthday, address and the clubs the shooter belongs to. Each club keeps its own shooters, imported from a CSV file, typed at the desk, or taken over from an online booking.

There is no shared register of all licensed shooters: the federation's member data is not part of Standblatt, and since clubs never see each other's data, the shooters one club entered stay invisible to the others. The price is that a shooter who shoots at three clubs' events is recorded three times, once in each club. Where the records have to come together, for a shooter's own training log, the licence number and the birthday do it (see [Closing an event](closing.md#my-training)).

## The shooter's club and the tenant club {#shooter-club}

Two different things are called a club:

- The **tenant club** is the club that uses Standblatt and runs the event.
- The **shooter's club** is an entry in the tenant's **Clubs** list: the club a shooter belongs to, e.g. a guest club at a shooting festival. These clubs have no access to Standblatt; they are names with an SSV number.

A shooter can belong to several clubs and registers for one of them. That club is printed on the documents and counts in the rankings and the club competition. Because clubs are matched by name regardless of upper case and spaces, and new ones are created only on purpose, a typo does not turn one club into two; where it happened anyway, merging cleans it up, even in closed events.

## Live updates {#live}

Several people work on the same event at once: the desk registers, scorers enter results at the lanes, kiosk screens show the rankings in the hall. Whenever something is saved, the server tells every open page of that event, and each page fetches what changed. A result entered at a lane appears on the dashboard and the kiosk screens without reloading. After a network interruption a page catches up by itself, so nothing is missed.

## Offline {#offline}

Club work needs a connection. A lane can only be given once, a shooter can shoot a programme only once, and a result saved on one device must be on all the others; only the server can decide that for everybody at the same time. Without a connection, club pages and helper links show a red bar and save nothing, but open dialogs keep what was typed, so it can be saved once the connection is back.

*My training* is different, because nobody else works on a personal log. A session saved without a connection is kept on the device and uploaded later, without ever creating it twice, and the installed app opens without a connection with the last data. See [Unfinished sessions and offline](../my-training/logging.md#unfinished-sessions-and-offline).
