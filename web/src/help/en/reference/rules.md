# Competition rules

This page describes the rules exactly as Standblatt implements them. They live in `internal/domain` and are covered by unit tests.

## Age category

The category follows from *event year − birth year*; see the table in [Concepts](../getting-started/concepts.md#age-categories). Without a birthday a shooter counts as Elite. The category is fixed when the shooter registers.

## Registration

A registration is refused when:

- a programme does not belong to the event, is chosen twice or is not active, or more attempts are chosen than it allows;
- a programme has no price for the shooter's age category ("not open to the age category");
- two variants of the same programme are chosen;
- a programme takes series from another programme that is not chosen, or the number of picked series does not match;
- the lane is taken by another active registration or is out of service;
- the shooter is already registered for one of the programmes at this event (a shooter shoots each programme once per event unless it allows several [attempts](../user-guide/programmes.md#attempts); another variant in a separate registration is fine).

A reserved lane can still be used after a confirmation.

## Results

Every series must be entered. Points lie between 0 and the series maximum (by default shots × 10), mouches between 0 and the number of shots, and only on series that track mouches. Saving results gives the most valuable award the total reaches; saving them again keeps the chosen award as long as the new total still reaches it, otherwise the most valuable one reached. Once every programme of a registration has results, the registration is completed and its lane becomes free.

### Combinations

A programme marked *combinable with* another programme can take series from it: at registration the shooter picks which series (for example series 2 and 3 of a four-series programme for a two-series programme). The picked series are copied in ascending order and renumbered from 1; the combined programme cannot be edited on its own and follows every change of the source.

## Awards

An award tier is available to an age category only if it has a limit for it. A shooter can choose any tier whose limit the total reaches; the most valuable tier is listed first and given by default when the results are saved. The server checks the choice.

## Ranking {#ranking}

All results of a programme are ordered by:

1. higher total;
2. the series compared from the last one backwards (the better last series wins, then the second-to-last, …);
3. more mouches;
4. age: among juniors (U10–U21) the younger shooter first, among the others the older shooter first; non-juniors rank before juniors.

Ranks are consecutive (1, 2, 3, … even on equal results). The same ranking feeds the staff pages, kiosk screens, the ranking PDF and CSV, and special prizes. It was checked against the published rankings of the 2026 Appenzell air pistol championship: identical on all 260 lines.

## Club competition {#club-competition}

Only programmes marked *counts for the club competition* take part, and only completed results of shooters with a club.

- The score of a result is the sum of its first two series.
- A club with fewer than 5 participants is listed as *not ranked*.
- Of *n* participants the best *max(5, ⌊n/2⌋)* are mandatory, the rest are optional.
- Result = ⌊(Σ mandatory + 0.02 × Σ optional) / number of mandatory × 1000⌋ / 1000, i.e. truncated to three decimals.
- Ties: more participants first, then the sorted scores compared one by one.
- The club excluded in the settings of the module *Club competition* (usually the host) keeps its place in the list but gets no rank ("out of competition").

## Special prizes {#special-prizes}

A programme can award special prizes to the best ranks of every age category. Three numbers define them:

- **Ranks per category**, e.g. 3;
- **minimum participants per category** in this programme, e.g. 10;
- **minimum participants of the shooter's club** in the same programme, e.g. 4.

Per age category with enough participants, the first shooters in ranking order whose club has enough participants win; shooters whose club has too few are skipped and the next one moves up. Shooters without a club are not eligible unless the club minimum is 0. The report lists both the winners and the skipped shooters with the reason.

!!! example "Appenzell air pistol championship"
    The brochure says: *"In every category, ranks 1–3 receive a special prize, provided the category has at least 10 participants and the shooter's club at least 4."* That is ranks 3, category minimum 10, club minimum 4. Applied to the 2026 results, the rule reproduces the hand-made winner list.

## Lane lists

Lanes are written as a list of numbers and ranges: `1, 5, 7-10`. Standblatt prints them compactly, turning runs of three or more into ranges (`1-3, 5, 7, 8`).

## Money and numbers

Amounts are shown as `CHF 1'234.50`, dates as `dd.mm.yyyy`. On wreath cards the value can be printed in words in the event's language, e.g. *Fünfundzwanzig Franken*, *Septante et un francs*, *Ventitré franchi*.
