# Einen Drucker am Schiessstand verbinden

Das brauchst du, wenn Standblätter, Kontrollstreifen und Kranzkarten aus dem
Drucker am Schiessstand kommen sollen, statt sich als PDF zu öffnen. Es
braucht das Recht *Stände und Drucker* (siehe
[Rollen](../user-guide/users.md#roles)).

Ein Computer neben den Druckern führt den Druck-Agent aus: ein Windows-PC
oder ein Mac mit der App Standblatt Print Agent oder ein Linux-Computer wie
ein Raspberry Pi mit dem Agent ohne Oberfläche. Warum das so funktioniert,
steht in [Wie das Drucken funktioniert](../explanation/printing.md).

## Mit der App unter Windows oder auf einem Mac {#app}

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
3. Klicke in der App auf **Schiessstand verbinden**. Sie zeigt einen Code
   wie `K7QD-4MZP` und öffnet **Druck-Agent verbinden** im Browser.
4. Melde dich an, falls verlangt, prüfe den auf der Seite genannten
   Computer, wähle den **Verein**, den **Schiessstand** und den **Name** des
   Agents und klicke auf **Verbinden**. Der Code ist zehn Minuten gültig.

![Ein Code in der App](../images/app/print-agent-code.png){ .screenshot }

Ist der Browser auf diesem Computer nicht zur Hand, öffne **Schiessstände →
Drucker einrichten** des Schiessstands auf einem beliebigen anderen
Computer, tippe den Code ins Feld neben **Mit Code verbinden** und klicke
darauf: Der Computer wird unter seinem eigenen Namen mit diesem Schiessstand
verbunden.

Für einen zweiten Schiessstand auf demselben Computer, z. B. *Pistole 25 m*
und *Pistole 50 m*, klicke in der App auf **Weiteren Schiessstand
verbinden** und wiederhole die Schritte 3 und 4.

## Mit dem Agent ohne Oberfläche unter Linux {#headless}

1. Öffne **Schiessstände** und klicke beim Schiessstand auf **Drucker
   einrichten**.
2. Gib unter **Druck-Agents** einen Namen ein (z. B. *Pi Halle 1*) und
   klicke auf **Druck-Agent hinzufügen**.
3. Klicke auf **Kopieren**, um die zwei Zeilen `STANDBLATT_URL` und
   `STANDBLATT_AGENT_TOKEN` zu kopieren. Das Token wird nur jetzt angezeigt.
4. Installiere und starte den Agent mit diesen Einstellungen, wie im [README
   des Druck-Agents](https://github.com/standblatt/print-agent#readme)
   beschrieben.

Für einen zweiten Schiessstand auf demselben Computer füge auch dort einen
Agent hinzu und gib dem Agent beide Tokens, durch ein Komma getrennt.

## Die Drucker wählen {#printers}

1. Öffne **Schiessstände → Drucker einrichten** des Schiessstands. Der Agent
   zeigt *verbunden* mit seiner Version und seinem Betriebssystem; seine
   Drucker erscheinen, sobald er sich verbunden hat.
2. Wähle für jedes Dokument den Drucker, dann seine Optionen (Fach, Papier,
   beidseitig …, wie der Agent sie für diesen Drucker meldet; unter Linux
   kommen weitere `lp`-Optionen wie `fit-to-page` ins Textfeld).
3. Klicke in der Zeile auf **Testseite**: Eine kleine Seite (90 × 29 mm) mit
   Verein, Schiessstand, Drucker, Optionen und Uhrzeit wird mit den Optionen
   dieser Zeile gedruckt. Ein Etikettendrucker druckt sie als Etikett,
   andere Drucker oben links auf dem Blatt.
4. Klicke in der Zeile auf **Speichern**. Ein Dokument, das auf *Browser
   (PDF)* bleibt, öffnet sich weiterhin als PDF.

![Drucken](../images/screenshots/printing.png){ .screenshot }

*Urkunden* sind die Urkunden aus den Ranglisten, *Druckvorlagen (vorgedruckt)* jede andere [Druckvorlage](../user-guide/print-forms.md) des Vereins. Ranglisten werden nie an einen Drucker geschickt.

## Prüfen, dass es druckt {#jobs}

Drucke etwas, z. B. ein Standblatt bei der nächsten Anmeldung, und sieh
unter **Letzte Aufträge** auf derselben Seite nach: *fertig* heisst, es
wurde gedruckt, *fehlgeschlagen* zeigt den Fehler daneben. Ein Auftrag, den
der Agent nicht innert zwei Minuten bestätigt, wird erneut versucht; nach
drei Versuchen schlägt er fehl.

Um einen Computer aufzugeben, klicke neben seinem Agent auf **widerrufen**;
mit der App entferne auch den Schiessstand unter **Erweitert** und
deinstalliere die App. Alle Ansichten der Druckseite beschreibt
[Drucken](../user-guide/printing.md).
