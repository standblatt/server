# Drucken

Standblatt erzeugt alle Dokumente als PDF. Jedes Dokument mit Format, Inhalt
und Druckereinstellung: [Dokumente](../reference/documents.md#documents).

Jeder Druckknopf druckt auf dem Drucker des Schiessstands, wenn einer
eingerichtet und sein Druck-Agent online ist; sonst öffnet sich das PDF
(siehe [Wenn kein Drucker da ist](../explanation/printing.md#fallback)).

Warum ein Druck-Agent und wie er funktioniert: [Wie das Drucken
funktioniert](../explanation/printing.md).

## Drucker einrichten

**Schiessstände → Drucker einrichten** (Recht *Stände und Drucker*) hat drei Teile:

![Drucken](../images/screenshots/printing.png){ .screenshot }

- **Druck-Agents**: die Computer, die für diesen Schiessstand drucken, jeder
  mit Version und Betriebssystem, *verbunden* oder wann er zuletzt gesehen
  wurde, und seiner Anzahl Drucker. **Druck-Agent hinzufügen** fügt einen
  Agent ohne Oberfläche hinzu und zeigt sein Token einmal, **Mit Code
  verbinden** verbindet die [App](#app) über ihren [Code](#connect),
  **widerrufen** entfernt einen Agent sofort.
- **Dokumente**: pro Dokument der Drucker (*Browser (PDF)* öffnet weiterhin
  das PDF) und die Optionen, die der Agent für diesen Drucker meldet; unter
  Linux kommen weitere `lp`-Optionen wie `fit-to-page` ins
  Textfeld. *Urkunden* sind die Urkunden aus den Ranglisten, *Druckvorlagen
  (vorgedruckt)* jede andere Vorlage des Vereins. **Speichern** speichert
  die Zeile. **Testseite** druckt eine kleine Seite (90 × 29 mm) mit Verein,
  Schiessstand, Drucker, Optionen und Uhrzeit auf dem Drucker dieser Zeile,
  auch vor dem Speichern; dafür muss ein Agent online sein.
- **Letzte Aufträge** zeigt, was gedruckt wurde und ob es geklappt hat.

Ein Auftrag, den der Agent nicht innert zwei Minuten bestätigt, wird erneut
versucht; nach drei Versuchen schlägt er fehl.

Schritt für Schritt: [Einen Drucker am Schiessstand
verbinden](../how-to/connect-printer.md).

## Die App Standblatt Print Agent {#app}

Die App ist der einfache Weg für einen Laptop oder PC bei den Druckern: kein
Token und keine Einstellungen zum Eintippen.

Herunterladen, installieren und verbinden: [Einen Drucker am Schiessstand
verbinden](../how-to/connect-printer.md#app).

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

**Druck-Agent verbinden** (von der App mit einem Code wie `K7QD-4MZP` geöffnet) zeigt den Computer und fragt nach Verein, Schiessstand und Name; nur Mitglieder, die Schiessstände verwalten dürfen, können einen Computer verbinden, und ein Code ist zehn Minuten gültig. **Mit Code verbinden** unter **Schiessstände → Drucker einrichten** verbindet ihn unter seinem eigenen Namen mit diesem Schiessstand. Schritte: [Einen Drucker am Schiessstand verbinden](../how-to/connect-printer.md#app).
