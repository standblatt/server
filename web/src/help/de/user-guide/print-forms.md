# Druckvorlagen

Eine Druckvorlage bringt Resultate auf ein vorgedrucktes Blatt, z. B. die
Standblätter, die die Armee für das Obligatorische Programm liefert: Die
Schüsse werden an der Scheibe oder nachträglich erfasst, und das Blatt wird
gedruckt statt von Hand ausgefüllt. Von Hand ausfüllen bleibt möglich;
Druckvorlagen sind freiwillig.

## Eine Vorlage erstellen {#create}

**Druckvorlagen** (Admins, im Vereinsteil der Seitenleiste) listet die Vorlagen des Vereins. **Neue Druckvorlage** fragt nach Name und Papiergrösse und öffnet dann den Editor:

- **Breite** und **Höhe** des Blatts in Millimetern.
- **Scan des leeren Formulars**: ein PNG oder JPEG des ganzen Blatts. Er
  erscheint hinter der **Vorschau**, damit sich die Felder auf die Kästchen
  des Formulars setzen lassen; gedruckt wird er nie.
- **Elemente**: was wo gedruckt wird. **x** und **y** sind Millimeter ab der
  oberen linken Ecke, **Breite** und **Höhe** die Grösse (Breite 0 lässt
  Text bis zum rechten Rand laufen), **Grösse** die Schriftgrösse in Punkt.
  - *Text* druckt seinen **Wert**, in dem Platzhalter ersetzt werden,
    z. B. `{{shooter.name}}` oder `{{total}}`.
  - *QR* und *Barcode* drucken ihren Wert als Code, *Linie* zieht eine Linie
    von x/y um Breite/Höhe.
  - **Wiederholung** druckt das Element mehrmals, jedes Mal um **dx** und
    **dy** verschoben; `#` in einem Platzhalter zählt 1, 2, 3, …. Ein
    Element `{{s1.shot#}}` mit Wiederholung 10 und dx 8 füllt die zehn
    Schusskästchen der Passe 1.
- **Platzhalter** listet jeden Wert mit einem Beispiel; ein Klick fügt ihn
  dem zuletzt bearbeiteten Element hinzu.

Das **Blatt** zeigt die Vorlage massstäblich über dem Scan, mit erfundenen
Werten. Ziehe ein Element, um es zu verschieben, und seine Ecke, um die
Grösse zu ändern; die Pfeiltasten verschieben das gewählte Element um einen
halben Millimeter, mit Shift um fünf, und ein Doppelklick aufs leere Blatt
fügt dort einen Text ein. Positionen rasten auf halbe Millimeter
ein. Tabelle und Blatt zeigen dieselben Elemente: Wähle eines auf dem Blatt
und passe seine genauen Werte in der Tabelle an, oder umgekehrt.

**Speichern** speichert die Vorlage. **PDF-Vorschau** zeigt die Vorlage genau so, wie das PDF sie darstellt; das kann leicht vom Blatt im Browser abweichen.

## Den Drucker kalibrieren {#calibrate}

Drucker drucken selten genau dort, wo sie sollen. Drucke ein Blatt mit
**Testdruck**, miss, wie weit die Werte daneben liegen, und trage es als
**Druckerversatz nach rechts** und **Druckerversatz nach unten** ein
(negativ für links oder oben). Der Versatz gilt für alle Elemente; die
Vorschau ignoriert ihn.

## Eine Vorlage verwenden {#use}

Wähle im Tab *Programme* der **Einrichtung** des Anlasses die
**Druckvorlage** eines Programms; das geht auch, wenn schon Schützen
angemeldet sind. Jedes abgeschlossene Programm bietet dann in den Resultaten
**Auf Vorlage drucken** an: am Desk, auf der Seite **Resultate erfassen**
und über Schreiber-Links an der Scheibe (siehe [Helfer](helpers.md)). Mit
einem Druck-Agenten und einem Drucker für *Druckvorlagen* unter
**Schiessstände → Drucken** geht das Blatt an diesen Drucker (z. B. an
seinen manuellen Einzug), sonst öffnet sich das PDF (siehe
[Drucken](printing.md)).

**Exportieren** lädt eine Vorlage als Datei herunter und **Importieren** liest sie wieder ein, z. B. um eine Vorlage mit einem anderen Verein zu teilen.
