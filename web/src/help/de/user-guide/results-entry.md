# Resultate erfassen

Resultate lassen sich auf drei Arten erfassen; innerhalb eines Anlasses
lassen sie sich mischen:

- **An der Scheibe** auf einem Tablet oder Handy, über einen Schreiber-Link
  (siehe [Helfer](helpers.md)): Der Schreiber erfasst die Passen oder jeden
  Schuss der Schützen auf seinen Scheiben.
- **Am Desk** über die Anmeldung und ihre Resultate (siehe
  [Resultate](results.md)), wie an einem Schützenfest.
- **Nachträglich** unter **Resultate erfassen** in der Seitenleiste des
  Anlasses, Schütze für Schütze, z. B. von Papierblättern oder der Anzeige
  der Trefferanzeige, oder als Datei aus der Trefferanzeige. Diese Seite
  beschreibt beides.

## Schütze für Schütze {#quick}

1. Wähle das **Programm**; es bleibt für die nächsten Schützen gewählt.
2. Tippe eine Lizenznummer oder einen Namen bei **Schütze** ein und drücke
   ++enter++: Der erste Treffer wird übernommen. Ein Barcode-Scanner
   funktioniert auch. **Neuer Schütze** erfasst jemanden, der noch nicht in
   den Stammdaten ist.
3. Tippe die Passen ein, oder die Einzelschüsse bei Passen, die sie erfassen
   (siehe [Programme](programmes.md)), und drücke ++enter++.

Das Resultat wird sofort gespeichert, und der Schütze erscheint unter
*Soeben erfasst* mit einem Link zu den Resultaten; bei einem Programm mit
[Druckvorlage](print-forms.md) bietet die Zeile **Auf Vorlage drucken**
an. War der Schütze am Desk für dieses Programm angemeldet und hat noch
keine Resultate, kommen die Resultate in diese Anmeldung; sonst wird eine
Anmeldung ohne Scheibe erstellt, bereits abgeschlossen. Ein Schütze kann in
einem Programm nur einmal Resultate haben, ausser das Programm erlaubt
mehrere [Durchgänge](programmes.md#attempts): Dann ist jedes weitere
Resultat ein neuer Durchgang, bis keiner mehr übrig ist.

Auszeichnungen werden wie gewohnt gewählt: Die wertvollste erreichte
Auszeichnung wird automatisch gesetzt und lässt sich in den Resultaten
ändern.

Tasten und Scanner: [Tastatur und
Scanner](../reference/shortcuts.md#quick-entry).

## Import aus einer Datei {#import}

**CSV importieren** (Recht *Anlässe konfigurieren*) unter **Resultate erfassen** liest eine Datei mit einer Zeile pro Schütze und Programm, z. B. aus der Trefferanzeige exportiert oder in einer Tabellenkalkulation getippt; die Spalten stehen unter [CSV-Import](../reference/csv-import.md#results). Jede Zeile wird zuerst geprüft; bei einem Fehler wird nichts importiert. Schritt für Schritt und wie man Zeilenfehler behebt: [Resultate aus der Trefferanzeige importieren](../how-to/import-results.md).
