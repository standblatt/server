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

1. Verbinde einen Computer bei den Druckern: einen Windows-PC oder einen Mac
   mit der [App Standblatt Print Agent](#app) per [Code](#connect), oder
   einen Linux-Computer wie einen Raspberry Pi mit dem [Agenten ohne
   Oberfläche](../admin-guide/print-agent.md): **Druck-Agent hinzufügen**,
   gib ihm einen Namen und kopiere das Token. Ein Computer kann mehrere
   Stände bedienen, z. B. *Pistole 25 m* und *Pistole 50 m*: Verbinde ihn
   mit jedem Stand (die App mit **Weiteren Schiessstand verbinden**, den
   Agenten ohne Oberfläche mit beiden Tokens, durch Komma getrennt).
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

## Die App Standblatt Print Agent {#app}

Die App ist der einfache Weg für einen Laptop oder PC bei den Druckern: kein
Token und keine Einstellungen zum Eintippen.

1. Lade sie vom [neusten
   Release](https://github.com/standblatt/print-agent/releases/latest)
   herunter: `Standblatt-Print-Agent-Setup.exe` für Windows,
   `standblatt-print-agent-app-darwin-arm64.dmg` für einen Mac mit
   Apple-Chip (`-amd64.dmg` für einen Intel-Mac), das `.deb` oder `.rpm` für
   Linux.
2. Installiere sie. Die App ist noch nicht signiert: Windows fragt einmal
   nach, **Weitere Informationen → Trotzdem ausführen**; auf dem Mac ziehst
   du sie in *Programme* und öffnest sie beim ersten Mal mit Rechtsklick →
   **Öffnen**. Unter Windows installiert sie sich für den angemeldeten
   Benutzer ohne Admin-Rechte, mit einer Verknüpfung im Startmenü und auf
   dem Desktop, und startet mit Windows.
3. Klicke **Schiessstand verbinden** und bestätige den Code in Standblatt
   (siehe [unten](#connect)).

![Die App Standblatt Print Agent](../images/app/print-agent.png){
.screenshot }

Das Fenster zeigt für jeden Schiessstand, ob er druckt, mit seinen letzten
Aufträgen, und darunter die Drucker des Computers, je mit einer
**Testseite**, die ohne Standblatt auf diesem Computer druckt. **Erweitert**
enthält *Mit dem Computer starten*, die Adresse von Standblatt (nur für
einen Testserver ändern) und das Entfernen eines Schiessstands. Schliesst du
das Fenster, läuft die App im Infobereich der Taskleiste weiter: Das Symbol
ist grün, wenn jeder Schiessstand druckt, sonst rot.

Unter Windows und auf dem Mac meldet die App, wenn eine neue Version da ist;
**Jetzt aktualisieren** installiert sie und startet die App neu. Brauchst du
einen Computer nicht mehr, entferne seinen Schiessstand unter **Erweitert**,
widerrufe seinen Agenten unter **Schiessstände → Drucker einrichten** und
deinstalliere die App (unter Windows bei *Einstellungen → Apps*).

## Computer per Code verbinden {#connect}

![Ein Code in der App](../images/app/print-agent-code.png){ .screenshot }

Die Standblatt-App auf einem Computer bei den Druckern braucht kein Token:
Sie zeigt einen Code wie `K7QD-4MZP` und öffnet **Druck-Agent verbinden** im
Browser. Melde dich dort an, falls nötig, prüfe den angezeigten Computer,
wähle Verein, Schiessstand und den Namen des Agenten und klicke
**Verbinden**. Nur Mitglieder, die Schiessstände verwalten dürfen, können
einen Computer verbinden; ein Code gilt zehn Minuten. Danach verbindet sich
die App von selbst, und der Agent erscheint unter **Schiessstände → Drucker
einrichten** mit seiner Version und seinem Betriebssystem.

Ist der Browser auf diesem Computer nicht zur Hand, gib den Code stattdessen
unter **Schiessstände → Drucker einrichten** des Schiessstands bei **Mit
Code verbinden** ein: Der Computer wird unter seinem eigenen Namen mit
diesem Schiessstand verbunden.
