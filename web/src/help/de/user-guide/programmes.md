# Programme

Programme gehören zu einem Anlass und werden unter **Konfiguration** →
**Programme** verwaltet (nur Admins).

![Die Programmliste](../images/screenshots/setup.png){ .screenshot }

## Programme hinzufügen

- **Neues Programm** öffnet den Editor.
- **Übernehmen aus … → Programm-Bibliothek**: eine Vorlage der globalen
  Bibliothek, die der Plattformbetreiber pflegt.
- **Übernehmen aus … → Früherer Anlass**: kopiert alle Programme eines
  anderen Anlasses des Vereins, inklusive Kombinationen und
  Sonderpreisen. Programme, deren Name und Variante schon existieren, werden
  übersprungen.
- **Übernehmen aus … → Programmdatei**: importiert eine
  [Programmdatei](../reference/programme-file.md), die aus einem anderen
  Anlass oder Verein exportiert wurde.

## Der Editor

![Programm-Editor](../images/screenshots/programme-editor.png){ .screenshot
}

- **Name** und **Variante**: Programme mit gleichem Namen sind Varianten
  voneinander; ein Schütze meldet sich für eine davon an. Das Namensfeld
  schlägt die Namen der Programme des Anlasses vor, und sobald der Name zu
  einem passt, sagt der Editor, dass du eine Variante anlegst, und warnt,
  wenn der Variantenname leer oder schon vergeben ist.
- **Barcode**: wird für die elektronische Trefferanzeige auf das Standblatt
  gedruckt (siehe [Codes](../reference/codes.md)).
- **Auf Kiosk-Anzeigen zeigen**.
- **Passen**: Schüsse, Maximalpunkte, ob Mouchen gezählt und ob
  **Einzelschüsse** erfasst werden, pro Passe. Mit Einzelschüssen lässt sich
  jeder Schuss erfassen, und die Passe ist deren Summe, z. B. um sie auf ein
  Formular zu drucken; das Passentotal allein lässt sich weiterhin erfassen.
- **Durchgänge** (siehe unten): wie oft ein Schütze das Programm am Anlass
  schiessen darf.
- **Preise (CHF)**: Nur die angekreuzten Alterskategorien dürfen das
  Programm schiessen. Bei mehreren Durchgängen ist ein zweites Feld pro
  Kategorie der Preis jedes weiteren Durchgangs (*Nachdoppel*); leer heisst
  gleicher Preis, 0 heisst, das Programm wird einmal bezahlt.
- **Auszeichnungen**: Name, Art (Kranzkarte, Naturalgabe oder Geldpreis, der
  in bar ausbezahlt wird), Wert und das Mindesttotal pro Alterskategorie;
  lass eine Kategorie leer, wo die Auszeichnung nicht verfügbar ist.

Sobald Schützen für ein Programm angemeldet sind, sind die **Passen** fest,
damit die Resultate ihre Bedeutung behalten; der Editor sagt das und graut
sie aus. Alles andere lässt sich weiterhin ändern: Name, Variante,
Beschreibung, Barcode, die Schalter, Preise, Passennamen, Namen, Werte und
Limiten der Auszeichnungen. Alterskategorien und Auszeichnungen können
**hinzugefügt**, aber nicht entfernt werden. Bei einem Programm in Gebrauch
geht ausserdem:

- **Aktiv**: Inaktive Programme lassen sich für neue Anmeldungen nicht
  wählen; bestehende behalten sie.
- **Kombinierbar mit**: das Programm, dessen Passen dieses bei der Anmeldung
  übernehmen darf.
- **Schiessstand** (Anlässe auf mehreren Ständen): der Stand, auf dem das
  Programm geschossen wird, z. B. *Pistole 50 m*; die Anmeldung bietet dann
  nur seine Scheiben an, und Kiosk-Bildschirme gruppieren die Ranglisten
  nach Stand. *beliebig* lässt es offen. Programme verschiedener Stände
  brauchen getrennte Anmeldungen.
- **Druckvorlage** (wenn der Verein [Druckvorlagen](print-forms.md) hat):
  die Vorlage, auf die die Resultate des Programms gedruckt werden; *keine*
  druckt keine Vorlage.
- **Sonderpreise** (siehe unten).

**Exportieren** lädt das Programm als Datei herunter; **Löschen** geht nur, solange niemand angemeldet ist.

## Durchgänge {#attempts}

Normalerweise schiesst ein Schütze ein Programm einmal pro
Anlass. **Durchgänge pro Schütze** erlaubt mehr, z. B. 3 für einen
Glücksstich mit Nachdoppel am Schützenfest, oder **Unbegrenzt** für ein
Jahresprogramm über die ganze Saison.

- **Zählender Durchgang** bestimmt, welcher Durchgang in Rangliste,
  Vereinswettkampf, Sonderpreisen und *Mein Training* zählt: **Der beste**
  (nach den Ranglistenregeln), **Der erste** oder **Der letzte**. Die
  Rangliste zeigt ihn als *2/3*: Der zweite von drei Durchgängen zählt.
- **Nur der zählende Durchgang erhält eine Auszeichnung**: ausgeschaltet
  kann jeder Durchgang seine Auszeichnung erhalten (eine Kranzkarte pro
  Nachdoppel); eingeschaltet erhält nur der zählende Durchgang eine, und die
  Auszeichnung wandert, wenn ein besserer Durchgang geschossen wird.
- Bei der Anmeldung zeigt ein Programm mit mehreren Durchgängen einen
  Zähler: Zwei Durchgänge auf einmal sind zwei Resultatkarten, *Durchgang 1*
  und *Durchgang 2*. Weitere Durchgänge lassen sich auch später in einer
  neuen Anmeldung hinzufügen oder indem man unter **Resultate erfassen** ein
  weiteres Resultat erfasst, bis der Schütze keine Durchgänge mehr hat.

Die Anzahl Durchgänge lässt sich ändern, solange Schützen angemeldet sind,
aber nicht unter die Durchgänge, die ein Schütze schon hat. Programme mit
mehreren Durchgängen können nicht kombiniert werden.

## Sonderpreise {#special-prizes}

**Sonderpreise** im Programmmenü zeichnet die besten Ränge jeder Alterskategorie aus, zum Beispiel die Ränge 1–3 in jeder Kategorie mit mindestens 10 Teilnehmern, wobei nur Schützen zählen, deren Verein mindestens 4 Teilnehmer in diesem Programm hat. Übersprungene Schützen machen dem nächsten Platz. Die genaue Regel steht in den [Wettkampfregeln](../reference/rules.md#special-prizes); die Gewinner stehen im Bericht *Sonderpreise* unter [Statistik](statistics-reports.md).

![Dialog Sonderpreise](../images/screenshots/special-prizes.png){
.screenshot }

Ränge 0 schaltet die Sonderpreise aus. Die Einstellungen lassen sich
jederzeit ändern, auch nach dem Anlass.
