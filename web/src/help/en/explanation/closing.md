# Closing an event

While an event runs, everything in it can still change: a result is corrected, a registration moves, a programme gets a new price. At some point the organizer has to say *this is the result*. That is what **Close** does: it turns a running event into a finished one.

## What closing means {#meaning}

Closing does two things at once:

- **The rankings become final and are released.** Until now they were a state of work; from now on they are the result of the event.
- **The event becomes read-only.** Registrations, entries, results, awards, programmes with their series, prices and award tiers, classifications and the event's settings can no longer be changed.

Everything stays visible: the dashboard, the results, the rankings, the statistics, the reports and the documents. A banner on every page of the event says that it is closed.

The two belong together. A ranking can only be final if nothing beneath it can change any more, and Standblatt enforces this in the database itself, not only by hiding buttons: every change to the data of a closed event is refused, whether it comes from the app, a helper link or an API token.

## Why closed events are read-only {#read-only}

Once an event is closed, other things rely on it:

- Shooters copy their results and ranks into their personal training logs.
- The financial report, the awards and the club competition have been handed out and settled.
- Next year's event may have been copied from it.

If results could still change, a shooter's log could show a rank the ranking no longer has, and the reports would no longer match what was paid and handed out. That is also why things that happen at the desk, such as marking awards as handed out, have to be done before closing.

## The link to My training {#my-training}

Every Standblatt user has a personal training log, [My training](../my-training/index.md). A shooter who enters a **licence number** and a **birthday** in the profile is matched with the shooters of every club: wherever a club's master data has a shooter with the same licence number (leading zeros do not matter) and the same birthday, that shooter's results belong to the user.

Only results of **closed** events are offered, because only their ranking is final. Each one comes with the rank in its programme and the number of participants, e.g. *Rank 12 of 87*. Nothing is added by itself: the shooter chooses a discipline and adds the result, and it becomes a competition session of the log with its series and its rank.

The match works across clubs without any club giving up its data: the clubs' master data stays where it is, the match only reads it for the user it belongs to, and only the shooter's own results are shown. A shooter without a licence number or a birthday in the profile, or whose club recorded a different birthday, gets no competition results. See [Profile and competition results](../my-training/account.md#competitions).

## Reopening {#reopen}

**Reopen** makes a closed event editable again, e.g. when a wrong result shows up after closing. It is meant for corrections, not as a second state of work:

- Results that shooters already added to their logs stay there as they were added; a correction afterwards does not reach them.
- Results not added yet are no longer offered in *My training* until the event is closed again.

Close the event again once the correction is done.

A closed event cannot be deleted; reopen it first. It can be copied for next year while it is closed.
