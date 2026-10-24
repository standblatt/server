# CSV-Import

Mitglieder mit dem Recht *Stammdaten verwalten* importieren Schützen auf der
Seite **Schützen**. Die Datei braucht eine Kopfzeile; das Trennzeichen (`;`
oder `,`) wird daraus übernommen.

```text
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9999;Hilltown;PS Example
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

## Resultate {#results}

Mitglieder mit dem Recht *Anlässe konfigurieren* importieren Resultate mit
**CSV importieren** auf der Seite **Resultate erfassen** eines Anlasses
(siehe [Resultate erfassen](../user-guide/results-entry.md#import)):

```text
license;firstname;lastname;birthday;club;programme;s1;s2;m2;s1_1;s1_2;s1_3
013756;;;;;Championship;96;93;2;;;
;Max;Muster;1990-02-02;SG Hilltown;Obligatory;;45;;9;10;8
```

| Spalte | Pflicht | Hinweise |
|---|---|---|
| `license` oder `firstname` + `lastname` | ja | der Schütze: über die Lizenz gefunden, sonst über den Namen (und `birthday`, wenn mehrere so heissen); angelegt, wenn nicht gefunden, wofür es `firstname` und `lastname` braucht: eine Lizenz allein, die nicht gefunden wird, ist ein Zeilenfehler |
| `birthday` | nein | `YYYY-MM-DD` oder `DD.MM.YYYY` |
| `club` | nein | muss in der Vereinsliste stehen (ohne Rücksicht auf Gross-/Kleinschreibung und Leerzeichen); der Schütze kommt dazu |
| `programme` | ja | Barcode, Name mit Variante wie `Championship (Sport)` oder der Name allein |
| `s1`, `s2`, … | ja | Punkte pro Passe, Dezimalstellen mit `.` oder `,` |
| `m1`, `m2`, … | nein | Mouchen pro Passe |
| `s1_1`, `s1_2`, … | nein | Einzelschüsse der Passe 1, für Passen, die sie erfassen; die Passe ist deren Summe |

- Jede Zeile wird zuerst geprüft; bei einem Zeilenfehler wird nichts
  importiert.
- Ein Schütze kann in einem Programm nur einmal Resultate haben; eine Zeile
  für schon vorhandene Resultate ist ein Fehler. Programme mit mehreren
  [Durchgängen](../user-guide/programmes.md#attempts) nehmen eine Zeile pro
  Durchgang, bis der Schütze keinen mehr übrig hat.
- Bei Programmen mit mehreren Durchgängen fügt ein erneuter Import derselben
  Datei ihre Zeilen als weitere Durchgänge hinzu.
- Ist der Schütze ohne Resultate für das Programm angemeldet, gehen die
  Resultate in diese Anmeldung.
