# Dashboard

Das Dashboard ist die Ansicht am Empfang eines Anlasses: wer auf welcher
Scheibe schiesst, wer wartet, wer fertig ist.

![Dashboard](../images/screenshots/dashboard.png){ .screenshot }

- **Am Schiessen**: Anmeldungen auf einer Scheibe, mit der Zeit seit ihrem
  Start, den Knöpfen **Resultate**, **Bearbeiten** und **Beenden**. Ist das
  Modul *Munitionsausgabe* an, hält **Munition** fest, dass sie ausgegeben
  wurde, und danach zeigt es ein grünes *Munition*-Abzeichen (siehe
  [Helfer](helpers.md)).
- **Wartend**: Anmeldungen ohne Scheibe (zum Beispiel, wenn alle Scheiben
  belegt sind). Ohne [Scheibenzuteilung](events.md#modules) (oder ohne
  Stand) gibt es keine Scheiben und kein *Am Schiessen*: Die Liste heisst
  dann **Angemeldet** und bietet **Resultate**, **Bearbeiten** und
  **Beenden**; eine Anmeldung verlässt sie, sobald ihre Resultate
  vollständig sind.
- **Beendet**: Anmeldungen, deren Programme alle Resultate haben, mit dem
  Abzeichen *Auszeichnungen abzugeben*, das den Resultatdialog für die
  Abgabe öffnet oder danach *Auszeichnungen abgegeben* (siehe
  [Resultate](results.md#handout)).
- **Scheiben**: Jede Scheibe ist *frei*, *belegt* (mit dem Schützen),
  *reserviert* (mit der Reservation und ihrem Ende) oder *ausser
  Betrieb*. Eine Reservation erscheint ab ihrer Vorlaufzeit. Ein Anlass auf
  mehreren Ständen zeigt ein Raster pro Stand, und der Anmeldedialog
  gruppiert die Scheiben nach Stand; ein Schütze, der auf einen anderen
  Stand wechselt, wird abgeschlossen und dort neu angemeldet.
- Die Zahlen oben zählen wartende, schiessende und abgeschlossene
  Anmeldungen.
- *Wartend* und *Abgeschlossen* zeigen 50 Anmeldungen pro Seite; die Pfeile
  unter der Liste blättern weiter.

Alles aktualisiert sich live: Ein an einem anderen Computer erfasstes
Resultat oder eine beginnende Reservation erscheint ohne Neuladen. Der Punkt
neben **live** in der Kopfzeile zeigt die Verbindung; nach einem
Netzunterbruch holt die Seite selbst auf. Ohne Verbindung weist ein roter
Balken darauf hin; dann lässt sich nichts speichern, aber offene Dialoge
behalten deine Eingaben, sodass du wieder speichern kannst, sobald du online
bist.

## Suchen und scannen

Drücke irgendwo im Dashboard ++f4++, um ins Suchfeld zu springen.

![Suche](../images/screenshots/search.png){ .screenshot }

- Tippe eine **Lizenznummer** (nur Ziffern, führende Nullen optional), einen
  **Namen** (Teil von Vor- oder Nachname) oder **Vor- und Nachname** in
  beliebiger Reihenfolge.
- Jeder Treffer bietet **&lt;Verein&gt; anmelden** für jeden Verein des
  Schützen an oder **Anmeldung öffnen**, wenn der Schütze gerade angemeldet
  ist.
- **Neuer Schütze** legt einen Schützen an und öffnet gleich die Anmeldung.
- Ein Handscanner, der den [Anmeldecode](../reference/codes.md) auf einem
  Standblatt oder Kontrollstreifen liest, tippt den Code und Enter: Die
  Resultate dieser Anmeldung öffnen sich direkt.

## Schiessen beenden

**Beenden** schliesst eine Anmeldung ab, bevor alle Resultate da sind, zum Beispiel wenn ein Schütze aufgibt oder die Scheiben später abgibt. Die Anmeldung wechselt zu *Abgeschlossen* und ihre Scheibe wird frei; fehlende Resultate lassen sich weiterhin erfassen.

## Dunkles Design

Der Sonne/Mond-Knopf in der Kopfzeile wechselt zwischen hellem und dunklem
Design; die Wahl bleibt im Browser gespeichert.

![Dashboard im dunklen Design](../images/screenshots/dashboard-dark.png){
.screenshot }
