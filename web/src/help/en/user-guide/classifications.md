# Classifications

A programme ranking ranks one programme. A **classification** ranks shooters over several programmes or over several [attempts](programmes.md#attempts) of one programme, e.g. a section championship as the sum of three programmes, or a club championship over the best three attempts of an annual programme. Classifications appear on **Rankings**, on kiosk displays and as a report; how they relate to the programme rankings is explained in [How rankings fit together](../explanation/rankings.md#classifications).

## Create a classification {#create}

With the module *Classifications* switched on (**Configuration** → **Modules**), **New classification** on **Classifications** in the configuration (admins) opens the editor:

- **Programmes**: every programme that counts, with its **Counting attempts**. With 3, the best three attempts count (or the first or last three, as set by the programme's *Counting attempt*). A shooter's total is the sum of all counted attempts.
- **Count only the best programmes**: with several programmes, only that many with the highest points count, e.g. the best two of three; 0 counts all.
- **Rank only shooters with all counting results**: shooters who are still missing a programme or an attempt are listed below the ranking without a rank. Without it everybody is ranked by what they have.
- **Tie-break**: equal totals are decided by the ticked criteria in their order; the arrows change it. *Series from the last one backwards* compares all counted series, *More of the highest single shots* needs [single shots](programmes.md) (more tens, then more nines …), *Age* follows the [SSV rule](../reference/rules.md#ranking). The default is the SSV order.
- **Classes**: separate rankings, each for some age categories, e.g. *Juniors* with U17 to U21 and *Elite* with the others. An age category belongs to one class at most; shooters whose category is in no class are not listed. Without classes everybody is in one ranking.
- **Show on kiosk displays**.

## Examples {#examples}

| Event | Programmes | Settings |
|---|---|---|
| Shooting festival, best two stiche | three stiche | *Count only the best programmes* 2 |

Step by step: [Rank a section championship over several programmes](../how-to/section-championship.md) and [Run an annual programme over the season](../how-to/annual-programme.md#classification).

**Copy for next year** copies the classifications with the programmes. Deleting a programme removes it from the classifications.
