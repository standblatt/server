# How rankings fit together

An event can rank its shooters in several ways at once: per programme, over several programmes or attempts, by club, and by age category for special prizes. They are not separate competitions with separate data. All of them start from the same results, and most of them start from the same programme ranking.

## The programme ranking {#programme}

The programme ranking is the base. Every programme has one ranking of everybody who shot it, ordered by the [SSV rules](../reference/rules.md#ranking): total, then the series from the last one backwards, then mouches, then age. There is only one implementation of this order, and it feeds the rankings page, the kiosk screens, the ranking PDF and CSV, the special prizes and the rank a shooter sees in *My training*. A shooter's place is therefore the same wherever it appears.

Ranks are consecutive (1, 2, 3, … even on equal results): the tie-break decides almost every tie, so two shooters never share a rank.

**Variants** are separate programmes with the same name, for example *Sport* and *Supported*, and each one has its own ranking: shooting supported and shooting freehand are not compared.

A **combined** programme takes series from another programme instead of being shot separately. Its results are copies of the picked series of its source, so it is ranked like any other programme; when the source changes, the copy follows, and with it the ranking.

## Attempts and the counting attempt {#attempts}

A programme that allows several attempts keeps every attempt as an entry with its own results. A ranking, however, lists each shooter once. The programme's **counting attempt** (the best, the first or the last) decides which attempt stands for the shooter, and the ranking shows which one it is, e.g. *2/3*.

The same counting attempt is used wherever one result per shooter is needed: the ranking, the club competition, the special prizes and the rank in *My training*. Awards are different, because an award belongs to a result rather than to a ranking: every attempt can earn its own award, unless *Only the counting attempt receives an award* is on, in which case the award moves to whichever attempt counts.

Combinations take the series of one result; with several attempts it would be unclear which one, so programmes with several attempts cannot be combined.

## Classifications {#classifications}

A classification ranks shooters over several programmes or over several attempts of one programme: the sum of three programmes for a section championship, the best three attempts of an annual programme for a club championship. It takes the counting attempts of each programme, adds up what counts, may count only the best programmes, and ranks with its own tie-break and its own classes of age categories.

A classification only reads results. It never changes a result, an award or a programme ranking, and an event can have several classifications next to its programme rankings. The programme ranking itself can be seen as the simplest classification: one programme, one attempt, the SSV tie-break and no classes. See [Classifications](../user-guide/classifications.md) for the settings.

## Club competition {#club-competition}

The club competition ranks clubs, not shooters. It uses only programmes marked *Counts for the club competition* and only completed results of shooters who registered for a club; the club is the one chosen at registration, not the tenant club. Each result scores with its first two series, and a club's result is the average of its best results, with a small bonus for every further participant (see [Club competition](../reference/rules.md#club-competition)). Because it has its own score, the club competition does not follow the programme ranking.

## Special prizes {#special-prizes}

Special prizes go to the best ranks of every age category. They do follow the programme ranking: they walk down it per age category and skip shooters whose club or category has too few participants (see [Special prizes](../reference/rules.md#special-prizes)). Their settings can be changed while the event runs, because they only read the ranking.

## Which ranking uses what {#overview}

| Ranking | Ranks | From | Counting attempt | Order |
|---|---|---|---|---|
| Programme ranking | shooters | one programme | yes | SSV rules |
| Classification | shooters, per class | several programmes or attempts | its own number per programme | own tie-break, SSV by default |
| Club competition | clubs | programmes counting for the club competition | yes | club result |
| Special prizes | shooters, per age category | the programme ranking | yes | the programme ranking |

Awards are not a ranking: a shooter reaches an award with their own total in their age category, however the others shoot.

## Why the series are locked {#locked}

Results are stored per series. Once shooters are registered for a programme, a change to the number of series, their shots, their maximum or whether they track mouches or single shots would change what the stored results mean: a total of 96 out of 100 is not a total out of 200, and every ranking that compares them would be wrong. Standblatt therefore fixes the series as soon as a shooter is registered.

Everything that does not change the meaning of a result stays open: names, prices, award values and limits, the flags. Age categories and award tiers can be added but not removed, because registrations were priced and awards were given by them.

## Live and final {#live}

All rankings are computed from the results whenever they are shown, so a saved result appears at once on the rankings page and the kiosk screens. While the event runs they can still change. [Closing](closing.md) the event makes them final.
