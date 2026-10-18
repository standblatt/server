# Drucken

Standblatt erzeugt alle Dokumente als PDF:

| Dokument | Format | Wann |
|---|---|---|
| Standblatt | A5 | bei der Anmeldung, mit QR-Code und Programm-Barcodes |
| Kontrollstreifen | Etikette 90 × 29 mm | bei der Anmeldung, wird auf die Scheibenkarte geklebt |
| Resultatblatt | A5 | wenn alle Programme Resultate haben |
| Kranzkarte | A6, auf vorgedruckte Karten gedruckt | nach Wahl einer Kranzkarte als Auszeichnung |
| Druckvorlage | jede Grösse, auf vorgedruckte Blätter gedruckt | bei Programmen mit einer [Druckvorlage](print-forms.md), sobald die Resultate da sind |
| Rangliste | A4 | von der Ranglistenseite aus, nur als PDF (geht nicht an den Drucker des Schiessstands) |

Jeder Druckknopf druckt auf dem Drucker des Schiessstands, wenn einer
eingerichtet ist und sein Druck-Agent online ist; sonst öffnet sich das PDF
im Browser und lässt sich von dort drucken. Auf Handy oder Tablet geht das
PDF stattdessen ins Teilen-Menü (drucken, in Dateien sichern, senden),
sodass die App offen bleibt; verlangt das Handy zuerst einen weiteren Tipp,
erscheint ein Hinweis mit dem Button **Teilen**.

## Drucker einrichten

**Schiessstände → Drucken** (Admins):

![Drucken](../images/screenshots/printing.png){ .screenshot }

1. **Druck-Agent hinzufügen**, gib ihm einen Namen und kopiere das
   Token. Installiere den Agenten auf einem Computer bei den Druckern (siehe
   [Druck-Agent](../admin-guide/print-agent.md)). Ein Computer kann mehrere
   Stände bedienen, z. B. *Pistole 25 m* und *Pistole 50 m*: Füge auf jedem
   Stand einen Agenten hinzu und gib dem Agenten beide Tokens, durch Komma
   getrennt.
2. Sobald der Agent online ist, meldet er seine Drucker. Wähle für jedes
   Dokument den Drucker (*Urkunden* sind die Urkunden aus den Ranglisten,
   *Druckvorlagen* alle anderen Vorlagen des Vereins) und seine Optionen
   (Fach, Papierformat, Medientyp, weitere `lp`-Optionen wie `fit-to-page`).
3. **Letzte Aufträge** zeigt, was gedruckt wurde und ob es geklappt hat.

Ein Auftrag, den der Agent nicht innert zwei Minuten bestätigt, wird
wiederholt; nach drei Versuchen schlägt er fehl. **widerrufen** entfernt
einen Agenten sofort.
