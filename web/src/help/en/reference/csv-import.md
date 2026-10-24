# CSV import

Members with the permission *Manage master data* import shooters on the **Shooters** page. The file needs a header row; the separator (`;` or `,`) is taken from it.

```text
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9999;Hilltown;PS Example
```

| Column | Required | Notes |
|---|---|---|
| `firstname` | yes | |
| `lastname` | yes | |
| `license` | no | also `licensenumber`; purely numeric numbers are padded to six digits |
| `birthday` | no | `YYYY-MM-DD` or `DD.MM.YYYY` |
| `postalcode`, `city` | no | |
| `club` | no | the shooter's club, matched regardless of upper case and spaces |

- One line per shooter and club; a shooter in two clubs appears twice.
- Shooters are matched by license number, or by name and birthday when there is none. Existing shooters are updated, nothing is deleted.
- A file with invalid lines is rejected as a whole, with the list of lines to fix.
- Clubs that are not in the club list stop the import before anything is stored: for each one, choose an existing club (similar names are suggested) or a new one, then **Import** again.

## Results {#results}

Members with the permission *Configure events* import results with **Import CSV** on an event's **Enter results** page (see [Enter results](../user-guide/results-entry.md#import)):

```text
license;firstname;lastname;birthday;club;programme;s1;s2;m2;s1_1;s1_2;s1_3
013756;;;;;Championship;96;93;2;;;
;Max;Muster;1990-02-02;SG Hilltown;Obligatory;;45;;9;10;8
```

| Column | Required | Notes |
|---|---|---|
| `license` or `firstname` + `lastname` | yes | the shooter: found by licence, else by name (and `birthday`, if several have the name); created when not found, which needs `firstname` and `lastname`: a licence alone that is not found is a line error |
| `birthday` | no | `YYYY-MM-DD` or `DD.MM.YYYY` |
| `club` | no | must be in the club list (regardless of upper case and spaces); the shooter joins it |
| `programme` | yes | barcode, name with variant like `Championship (Sport)`, or the name alone |
| `s1`, `s2`, … | yes | points per series, decimals with `.` or `,` |
| `m1`, `m2`, … | no | mouches per series |
| `s1_1`, `s1_2`, … | no | single shots of series 1, for series that record them; the series is their sum |

- Every line is checked first; with any line error nothing is imported.
- A shooter can have results in a programme only once; a line for results that exist already is an error. Programmes with several [attempts](../user-guide/programmes.md#attempts) take one line per attempt until the shooter has none left.
- For programmes with several attempts, importing the same file again adds its lines as further attempts.
- If the shooter is registered for the programme without results, the results go into that registration.
