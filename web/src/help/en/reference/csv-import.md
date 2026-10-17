# CSV import

Club admins import shooters on the **Shooters** page. The file needs a header row; the separator (`;` or `,`) is taken from it.

```text
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9050;Appenzell;PS Example
```

| Column | Required | Notes |
|---|---|---|
| `firstname` | yes | |
| `lastname` | yes | |
| `license` | no | also `licensenumber`; purely numeric numbers are padded to six digits |
| `birthday` | no | `YYYY-MM-DD` or `DD.MM.YYYY` |
| `postalcode`, `city` | no | |
| `club` | no | the shooter's club; missing clubs are created |

- One line per shooter and club; a shooter in two clubs appears twice.
- Shooters are matched by license number, or by name and birthday when there is none. Existing shooters are updated, nothing is deleted.
- A file with invalid lines is rejected as a whole, with the list of lines to fix.
