# Druckvorlagen

Eine Druckvorlage bringt Resultate auf ein vorgedrucktes Blatt, z. B. die
Standblätter, die die Armee für das Obligatorische Programm liefert: Die
Schüsse werden an der Scheibe oder nachträglich erfasst, und das Blatt wird
gedruckt statt von Hand ausgefüllt. Von Hand ausfüllen bleibt möglich;
Druckvorlagen sind freiwillig und lassen sich unter **Vereinsmodule**
ausschalten (siehe [Vereinsmodule](events.md#club-modules)).

Warum und wie Druckvorlagen funktionieren: [Druckvorlagen auf vorgedruckten
Blättern](../explanation/printing.md#forms).

## Eine Vorlage erstellen {#create}

**Druckvorlagen** (Recht *Druckvorlagen*, im Vereinsteil der Seitenleiste) listet die Vorlagen des Vereins. **Neue Druckvorlage** fragt nach einem Namen, was die Vorlage **Druckt** (siehe [Datenquelle](#source)) und einem Papierformat und öffnet dann den Editor. Er hat drei Spalten:

- Links die **Felder**: alles, was eine Vorlage drucken kann, in Gruppen, je
  mit einem Beispiel, und eine Suche. Ein Klick fügt das Feld unter dem
  untersten Element ein; aufs Blatt gezogen landet es dort, wo man es
  loslässt.
  - *Eigener Text* druckt festen Text, *QR-Code* und *Barcode* drucken ihren
    Wert als Code, *Linie* zieht eine Linie.
  - *Schütze*, *Verein und Anlass* und *Resultat* drucken je einen Wert,
    z. B. den Namen, das Total oder die Punkte einer Passe.
  - *Einzelschüsse* fügt eine ganze Reihe Schusskästchen ein: *Schüsse Passe
    1* ist ein Element `{{s1.shot#}}`, einmal pro Schuss wiederholt, 8 mm
    auseinander.
- Das **Blatt** in der Mitte zeigt die Vorlage massstabgetreu über dem Scan,
  mit erfundenen Werten; **Feldnamen zeigen** zeigt stattdessen die
  Platzhalter. Oben und links laufen Lineale in Millimetern, der Zoom reicht
  von 50 bis 200 %.
- Rechts bearbeiten die **Eigenschaften** das gewählte Element genau: seinen
  **Wert** (Platzhalter wie `{{shooter.name}}` werden beim Drucken
  ausgefüllt), **x** und **y** in Millimetern ab der oberen linken Ecke,
  **Breite** und **Höhe** (bei Breite 0 läuft Text bis zum rechten Rand),
  **Schriftgrösse**, **Fett**, **Kursiv**, **Text umbrechen** (der Text
  bricht innerhalb der Breite in Zeilen um, jede so hoch wie die Höhe) und
  **Ausrichtung**. **Wiederholung** druckt es mehrmals, jedes Mal um **dx**
  und **dy** verschoben; `#` in einem Platzhalter zählt 1, 2, 3, …. Darunter
  legt die **Vorlage** den Namen, die Grösse des Blatts, den Druckversatz
  und den **Scan des leeren Formulars** fest: ein PNG oder JPEG des ganzen
  Blatts, das hinter dem Blatt erscheint, damit die Felder auf die Kästchen
  des Formulars gesetzt werden können; es wird nie gedruckt.

**Speichern** speichert die Vorlage; ein Punkt auf dem Knopf zeigt ungespeicherte Änderungen, und wer die Seite mit ihnen verlässt, wird zuerst gefragt. Mit **Automatisch speichern** (eingeschaltet, bis man es ausschaltet; der Browser merkt es sich) werden ungespeicherte Änderungen jede Minute gespeichert. **PDF-Vorschau** öffnet die gespeicherte Vorlage genau so, wie das PDF sie darstellt; das kann leicht vom Blatt im Browser abweichen.

## Datenquelle {#source}

Was eine Vorlage **Druckt**, bestimmt ihre Felder und wo sie gedruckt wird:

| Druckt | Felder | Gedruckt aus |
|---|---|---|
| Resultate eines Programms | Schütze, Verein, Anlass, Passen, Einzelschüsse, Total, Auszeichnung und ihr Wert | den Resultaten eines Programms (siehe [Vorlage verwenden](#use)) oder als Kranzkarte des Anlasses |
| Eine Anmeldung | Schütze, Verein, Anlass, Scheibe, Programme, Passen, QR-Code der Anmeldung | als Kontrollstreifen des Anlasses |
| Ein Rang | Schütze, Verein, Anlass, Rang, Total, Rangliste, Klasse | Urkunden aus den Ranglisten |

Die Quelle lässt sich im Editor unter **Vorlage** ändern; Felder, die die
neue Quelle nicht kennt, drucken nichts.

Jeder Platzhalter: [Platzhalter in
Vorlagen](../reference/form-placeholders.md).

## Eingebaute Vorlagen {#builtin}

Standblatt druckt die Kranzkarte und den Kontrollstreifen mit **eingebauten
Vorlagen**, aufgeführt unter den Vorlagen des Vereins. **Kopieren** macht
daraus eine Vorlage des Vereins, die sich wie jede andere ändern lässt,
z. B. um die Felder auf eine andere vorgedruckte Kranzkarte zu setzen oder
das Logo des Vereins hinzuzufügen. Ein Anlass wählt sie dann unter
**Konfiguration** → **Module**: *Vorlage der Kranzkarten* im Modul
*Auszeichnungen*, *Vorlage der Kontrollstreifen* im Modul
*Kontrollstreifen*. Ohne Wahl, oder wenn die gewählte Vorlage gelöscht ist,
druckt die eingebaute Vorlage.

![Druckvorlagen mit den eingebauten Vorlagen und der
Bildbibliothek](../images/screenshots/print-forms.png){ .screenshot }

## Anmeldedokumente und Urkunden {#certificates}

Eine Vorlage, die **Eine Anmeldung** druckt, erscheint neben dem
Kontrollstreifen des Anlasses in der Anmeldung als Knopf mit ihrem eigenen
Namen, z. B. eine Teilnehmerkarte oder eine Quittung; sie druckt auf dem
Drucker der *Druckvorlagen (vorgedruckt)*.

Eine Vorlage, die **Ein Rang** druckt, erzeugt Urkunden über **Urkunden**
unter **Ranglisten**, gedruckt auf dem Drucker der *Urkunden* des ersten
Schiessstands des Anlasses; siehe [Urkunden für die besten Ränge
drucken](../how-to/certificates.md).

![Urkunden aus einer Rangliste](../images/screenshots/certificates.png){
.screenshot }

## Elemente anordnen {#arrange}

- Ein Element wird durch Ziehen **verschoben** und an seiner Ecke **in der
  Grösse geändert**. Der pinke Punkt am letzten Kästchen eines wiederholten
  Elements zieht den Abstand auf.
- Beim Verschieben **rastet** das Element an den Kanten und Mitten der
  anderen Elemente und des Blatts ein, eine pinke Hilfslinie zeigt die
  Linie; sonst rastet es auf ganze Millimeter. Mit gedrückter Alt-Taste
  lässt es sich frei setzen.
- **Mehrere auswählen** mit Shift oder Ctrl und einem Klick, oder auf dem
  leeren Blatt einen Rahmen um sie ziehen. Ziehen verschiebt sie gemeinsam;
  Schriftgrösse, Fett und Ausrichtung gelten dann für alle.
- Ein **Rechtsklick** öffnet das Menü: *Duplizieren*, *In den Vordergrund*,
  *Ausrichtung* (Kanten, Mitten, gleichmässig verteilen), *Auf dem Blatt
  zentrieren* und *Löschen*.

Tasten: [Editor für Druckvorlagen](../reference/shortcuts.md#print-forms).

## Bilder {#images}

Logos, Briefköpfe und andere Bilder kommen aus der Bildbibliothek des
Vereins unter der Liste der **Druckvorlagen**: **Bild hochladen** fügt ein
PNG oder JPEG von höchstens 2 MB hinzu, benannt nach seiner Datei. Im Editor
listet die Gruppe **Bilder** der Felder sie auf; ein Klick oder Ziehen legt
eines aufs Blatt, 30 mm breit in den Proportionen des Bildes. Seine
Eigenschaften wählen ein anderes **Bild**, seine Position und Grösse, und
**Seitenverhältnis beibehalten** passt die Höhe wieder an die Breite an. Ein
Bild lässt sich wie jedes andere Element wiederholen.

Ein Bild lässt sich erst löschen, wenn keine Vorlage es mehr
zeigt. **Exportieren** nimmt die Bilder einer Vorlage in ihre Datei mit, und
**Importieren** fügt sie der Bibliothek hinzu; ein Bild gleichen Namens, das
der Verein schon hat, wird wiederverwendet.

## Den Drucker kalibrieren {#calibrate}

**Druckerversatz nach rechts (mm)** und **Druckerversatz nach unten (mm)** unter **Vorlage** verschieben jedes Element (negativ, um nach links oder oben zu verschieben). **Testdruck** öffnet die gespeicherte Vorlage mit dem Versatz und ohne den Scan; die Vorschau ignoriert den Versatz. Warum: [Kalibrieren im Grundsatz](../explanation/printing.md#calibration). Schritte: [Resultate auf die Standblätter der Armee drucken](../how-to/army-sheet.md#calibrate).

## Eine Vorlage verwenden {#use}

Ein Programm wählt seine Vorlage im Abschnitt **Druckvorlage** seiner
Programmseite, auch wenn schon Schützen angemeldet sind. Abgeschlossene
Resultate bieten dann **Auf Vorlage drucken** an: in den Resultaten am Desk,
unter *Soeben erfasst* auf **Resultate erfassen** und auf Helferlinks von
Warnern und der Resultaterfassung (siehe [Helfer](helpers.md)). Gedruckt
wird auf dem Drucker der *Druckvorlagen (vorgedruckt)* des Schiessstands der
Anmeldung, sonst öffnet sich das PDF (siehe
[Drucken](printing.md)). Schritte: [Resultate auf die Standblätter der Armee
drucken](../how-to/army-sheet.md).

**Exportieren** lädt eine Vorlage mit ihren Bildern als Datei herunter und **Importieren** liest sie wieder ein, z. B. um eine Vorlage mit einem anderen Verein zu teilen.
