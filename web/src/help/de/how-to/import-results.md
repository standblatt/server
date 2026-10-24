# Resultate aus der Trefferanzeige importieren

Das brauchst du, wenn die Resultate vieler Schützen in einer Datei stehen,
z. B. aus der elektronischen Trefferanzeige exportiert oder in eine Tabelle
getippt. Es braucht das Recht *Anlässe konfigurieren* (siehe
[Rollen](../user-guide/users.md#roles)); der Anlass darf nicht abgeschlossen
sein.

## Die Datei vorbereiten {#file}

1. Exportiere oder speichere die Resultate als CSV mit einer Kopfzeile und
   einer Zeile pro Schütze und Programm. Die Spalten stehen unter
   [CSV-Import](../reference/csv-import.md#results): mindestens `license`
   (oder `firstname` und `lastname`), `programme` und die Passen `s1`, `s2`,
   ….
2. Verwende in der Spalte `programme` den Barcode des Programms oder seinen
   Namen, wo es Varianten gibt mit der Variante in Klammern,
   z. B. `Championship (Sport)`.
3. Hat die Datei eine Spalte `club`, prüfe, dass jeder Verein in der
   Vereinsliste unter **Stammdaten** → **Vereine** steht; füge fehlende dort
   hinzu.

## Die Datei importieren {#import}

1. Öffne den Anlass und klicke in der Seitenleiste auf **Resultate
   erfassen**.
2. Klicke auf **CSV importieren**, wähle die Datei und klicke auf
   **Importieren**.
3. *Importiert: … Resultate, … neue Schützen.* heisst, du bist fertig; die
   Resultate stehen unter **Ranglisten**.

Importiere eine Datei nur einmal. Ein zweiter Import wird bei Programmen,
die einmal geschossen werden, abgelehnt (die Schützen haben schon
Resultate), aber ein Programm mit mehreren
[Durchgängen](../user-guide/programmes.md#attempts) nimmt jede Zeile erneut
als weiteren Durchgang.

## Zeilenfehler beheben {#errors}

Bei einem Fehler wird nichts importiert, und der Dialog zeigt *Nichts
importiert, bitte diese Zeilen korrigieren:* mit den
Zeilennummern. Korrigiere jede aufgeführte Zeile in der Datei und importiere
die ganze Datei erneut.

| Meldung | Behebung |
|---|---|
| Programme "…" not found | schreibe das Programm wie auf seiner Programmseite, bei mehreren Varianten mit der Variante, oder verwende seinen Barcode |
| Club "…" is not in the club list | korrigiere die Schreibweise oder füge den Verein unter **Vereine** hinzu |
| No shooter with the licence …; add first and last name to create one | ergänze die Zeile mit `firstname` und `lastname`, damit der Schütze erstellt werden kann |
| Several shooters are called …; add the birthday or the licence | ergänze die Zeile mit `birthday` oder `license` |
| Series … is missing, … is not a number | fülle jede Passe des Programms aus; Dezimalstellen mit `.` oder `,` |
| Series …: points must be between 0 and … | korrigiere die Punkte; das Maximum bestimmen die Passen des Programms |
| Series …: enter all … shots | gib jeden Einzelschuss der Passe an, oder keinen und nur das Passentotal |
| The shooter already has results in … | entferne die Zeile; die Resultate sind schon in Standblatt |
| The same shooter and programme appear twice in the file | behalte eine der zwei Zeilen |
| The shooter has no attempts left in … | entferne die Zeile oder erhöhe die Durchgänge des Programms |

Eine Meldung zur Kopfzeile (*The header row needs the column programme*, *…
the column license or the columns firstname and lastname*) heisst, dass in
der ersten Zeile ein Spaltenname fehlt oder falsch geschrieben ist.

Schützen, die nicht in den Stammdaten sind, legt der Import an. Einzelne
Resultate lassen sich danach noch in den Resultaten der Anmeldung
korrigieren (siehe [Resultate](../user-guide/results.md)).
