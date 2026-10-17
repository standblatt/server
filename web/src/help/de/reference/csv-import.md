# CSV-Import

Vereins-Admins importieren Schützen auf der Seite **Schützen**. Die Datei
braucht eine Kopfzeile; das Trennzeichen (`;` oder `,`) wird daraus
übernommen.

```text
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9050;Appenzell;PS Example
```

| Spalte | Pflicht | Hinweise |
|---|---|---|
| `firstname` | ja | |
| `lastname` | ja | |
| `license` | nein | auch `licensenumber`; rein numerische Nummern werden auf sechs Stellen aufgefüllt |
| `birthday` | nein | `YYYY-MM-DD` oder `DD.MM.YYYY` |
| `postalcode`, `city` | nein | |
| `club` | nein | der Verein des Schützen, abgeglichen ohne Rücksicht auf Gross-/Kleinschreibung und Leerzeichen |

- Eine Zeile pro Schütze und Verein; ein Schütze in zwei Vereinen erscheint
  zweimal.
- Schützen werden über die Lizenznummer zugeordnet oder über Name und
  Geburtsdatum, wenn es keine gibt. Bestehende Schützen werden aktualisiert,
  gelöscht wird nichts.
- Eine Datei mit ungültigen Zeilen wird als Ganzes abgelehnt, mit der Liste
  der Zeilen, die zu korrigieren sind.
- Vereine, die nicht in der Vereinsliste stehen, halten den Import an, bevor
  etwas gespeichert wird: Wähle für jeden einen bestehenden Verein (ähnliche
  Namen werden vorgeschlagen) oder einen neuen und klicke dann nochmals auf
  **Importieren**.
