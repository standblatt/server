# Import results from the scoring system

You need this when the results of many shooters are in a file, e.g. exported from the electronic scoring system or typed into a spreadsheet. It needs the permission *Configure events* (see [Roles](../user-guide/users.md#roles)); the event must not be closed.

## Prepare the file {#file}

1. Export or save the results as CSV with a header row and one line per shooter and programme. The columns are listed under [CSV import](../reference/csv-import.md#results): at least `license` (or `firstname` and `lastname`), `programme` and the series `s1`, `s2`, ….
2. In the `programme` column, use the programme's barcode or its name, with the variant in parentheses where there are variants, e.g. `Championship (Sport)`.
3. If the file has a `club` column, check that every club is in the club list under **Master data** → **Clubs**; add missing ones there.

## Import it {#import}

1. Open the event and click **Enter results** in the sidebar.
2. Click **Import CSV**, choose the file and click **Import**.
3. *Imported: … results, … new shooters.* means you are done; the results are on **Rankings**.

Import a file only once. A second import is refused for programmes that are shot once (the shooters already have results), but a programme with several [attempts](../user-guide/programmes.md#attempts) takes every line again as a further attempt.

## Fix line errors {#errors}

With any error, nothing is imported and the dialog shows *Nothing imported, please fix these lines:* with the line numbers. Fix every listed line in the file and import the whole file again.

| Message | Fix |
|---|---|
| Programme "…" not found | write the programme as on its programme page, with the variant if there are several, or use its barcode |
| Club "…" is not in the club list | correct the spelling or add the club under **Clubs** |
| No shooter with the licence …; add first and last name to create one | add `firstname` and `lastname` to the line, so the shooter can be created |
| Several shooters are called …; add the birthday or the licence | add `birthday` or `license` to the line |
| Series … is missing, … is not a number | fill in every series of the programme; decimals take `.` or `,` |
| Series …: points must be between 0 and … | correct the points; the maximum is set by the programme's series |
| Series …: enter all … shots | give every single shot of the series, or none and only the series total |
| The shooter already has results in … | remove the line; the results are in Standblatt already |
| The same shooter and programme appear twice in the file | keep one of the two lines |
| The shooter has no attempts left in … | remove the line, or raise the attempts of the programme |

A message about the header row (*The header row needs the column programme*, *… the column license or the columns firstname and lastname*) means a column name is missing or misspelt in the first line.

Shooters who are not in the master data are created by the import. Single results can still be corrected afterwards in the results of the registration (see [Results](../user-guide/results.md)).
