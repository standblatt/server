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

**Schiessstände → Drucker einrichten** (Admins):

![Drucken](../images/screenshots/printing.png){ .screenshot }

1. **Druck-Agent hinzufügen**, gib ihm einen Namen und kopiere das
   Token. Installiere den Agenten auf einem Computer bei den Druckern (siehe
   [Druck-Agent](../admin-guide/print-agent.md)). Ein Computer kann mehrere
   Stände bedienen, z. B. *Pistole 25 m* und *Pistole 50 m*: Füge auf jedem
   Stand einen Agenten hinzu und gib dem Agenten beide Tokens, durch Komma
   getrennt.
2. Sobald der Agent online ist, meldet er seine Drucker, und die Liste zeigt
   seine Version und sein Betriebssystem. Wähle für jedes Dokument den
   Drucker (*Urkunden* sind die Urkunden aus den Ranglisten, *Druckvorlagen*
   alle anderen Vorlagen des Vereins) und seine Optionen: die Optionen, die
   der Agent für diesen Drucker meldet (unter Linux Fach, Papierformat und
   Medientyp, dazu weitere `lp`-Optionen wie `fit-to-page` im Textfeld;
   unter Windows Fach, Papier, beidseitig, Skalierung und Farbe).
3. **Testseite** druckt eine kleine Seite (90 × 29 mm, so gross wie der
   Kontrollstreifen) auf dem Drucker dieser Zeile mit den dort gewählten
   Optionen, auch vor dem Speichern: Verein, Schiessstand, Drucker, Optionen
   und Uhrzeit. Ein Etikettendrucker druckt sie als Etikett, andere Drucker
   oben links auf dem Blatt. Dafür muss ein Agent online sein.
4. **Letzte Aufträge** zeigt, was gedruckt wurde und ob es geklappt hat.

Ein Auftrag, den der Agent nicht innert zwei Minuten bestätigt, wird
wiederholt; nach drei Versuchen schlägt er fehl. **widerrufen** entfernt
einen Agenten sofort.

## Computer per Code verbinden {#connect}

Die Standblatt-App auf einem Computer bei den Druckern braucht kein Token:
Sie zeigt einen Code wie `K7QD-4MZP` und öffnet **Druck-Agent verbinden** im
Browser. Melde dich dort an, falls nötig, prüfe den angezeigten Computer,
wähle Verein, Schiessstand und den Namen des Agenten und klicke
**Verbinden**. Nur Mitglieder, die Schiessstände verwalten dürfen, können
einen Computer verbinden; ein Code gilt zehn Minuten. Danach verbindet sich
die App von selbst, und der Agent erscheint unter **Schiessstände →
Drucker einrichten** mit seiner Version und seinem Betriebssystem.
