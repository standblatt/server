# Enter results

Results can be entered in three ways; they can be mixed within one event:

- **At the lane** on a tablet or phone, through a scorer link (see [Helpers](helpers.md)): the scorer enters the series, or every shot, of the shooters on their lanes.
- **At the desk** through the registration and its results (see [Results](results.md)), the way of a shooting festival.
- **Afterwards** on **Enter results** in the event's sidebar, shooter by shooter, e.g. from paper sheets or the display of the scoring system, or as a file exported from the scoring system. This page describes both.

## Shooter by shooter {#quick}

1. Choose the **Programme**; it stays chosen for the next shooters.
2. Type a licence number or a name into **Shooter** and press ++enter++: the first hit is taken. A barcode scanner works too. **New shooter** adds someone who is not in the master data yet.
3. Type the series, or the single shots for series that record them (see [Programmes](programmes.md)), and press ++enter++.

The result is saved at once and the shooter is listed under *Entered now* with a link to the results. If the shooter was registered at the desk for this programme and has no results yet, the results go into that registration; otherwise a registration without a lane is created, already completed. A shooter can have results in a programme only once, unless the programme allows several [attempts](programmes.md#attempts): then every further result is a new attempt until none is left.

Awards are chosen as usual: the most valuable award reached is set automatically and can be changed in the results.

## Import from a file {#import}

**Import CSV** (admins) on **Enter results** reads a file with one line per shooter and programme, e.g. exported from the scoring system or typed in a spreadsheet; the columns are described under [CSV import](../reference/csv-import.md#results). Every line is checked first: with any error nothing is imported, and the dialog lists the lines to fix. Shooters who are not in the master data are created; a club must already be in the club list.
