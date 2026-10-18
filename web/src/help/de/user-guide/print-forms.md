# Druckvorlagen

Eine Druckvorlage bringt Resultate auf ein vorgedrucktes Blatt, z. B. die
Standblätter, die die Armee für das Obligatorische Programm liefert: Die
Schüsse werden an der Scheibe oder nachträglich erfasst, und das Blatt wird
gedruckt statt von Hand ausgefüllt. Von Hand ausfüllen bleibt möglich;
Druckvorlagen sind freiwillig und lassen sich unter **Vereinsbausteine**
ausschalten (siehe [Vereinsbausteine](events.md#club-modules)).

## Eine Vorlage erstellen {#create}

**Druckvorlagen** (Admins, im Vereinsteil der Seitenleiste) listet die Vorlagen des Vereins. **Neue Druckvorlage** fragt nach einem Namen und einem Papierformat und öffnet dann den Editor. Er hat drei Spalten:

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
  **Schriftgrösse**, **Fett** und **Ausrichtung**. **Wiederholung** druckt
  es mehrmals, jedes Mal um **dx** und **dy** verschoben; `#` in einem
  Platzhalter zählt 1, 2, 3, …. Darunter legt die **Vorlage** den Namen, die
  Grösse des Blatts, den Druckversatz und den **Scan des leeren Formulars**
  fest: ein PNG oder JPEG des ganzen Blatts, das hinter dem Blatt erscheint,
  damit die Felder auf die Kästchen des Formulars gesetzt werden können; es
  wird nie gedruckt.

**Speichern** speichert die Vorlage; ein Punkt auf dem Knopf zeigt ungespeicherte Änderungen, und wer die Seite mit ihnen verlässt, wird zuerst gefragt. Mit **Automatisch speichern** (eingeschaltet, bis man es ausschaltet; der Browser merkt es sich) werden ungespeicherte Änderungen jede Minute gespeichert. **PDF-Vorschau** öffnet die Vorlage genau so, wie das PDF sie darstellt; das kann leicht vom Blatt im Browser abweichen.

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

| Taste | Wirkung |
|---|---|
| Pfeiltasten | verschieben die Auswahl um 0.5 mm, mit Shift um 5 mm |
| Ctrl+Z / Ctrl+Y | rückgängig / wiederholen |
| Ctrl+A | alles auswählen |
| Ctrl+D | duplizieren |
| Entf | Auswahl löschen |
| Esc | Auswahl aufheben |
| Ctrl+S | speichern |
| Ctrl+Mausrad | zoomen |

## Den Drucker kalibrieren {#calibrate}

Drucker drucken selten genau dort, wo sie sollen. Drucke ein Blatt mit
**Testdruck**, miss, wie weit die Werte daneben liegen, und trage es als
**Druckerversatz nach rechts** und **Druckerversatz nach unten** ein
(negativ für links oder oben). Der Versatz gilt für alle Elemente; die
Vorschau ignoriert ihn.

## Eine Vorlage verwenden {#use}

Öffne unter **Konfiguration** → **Programme** des Anlasses ein Programm und
wähle im gleichnamigen Abschnitt seine **Druckvorlage**; das geht auch, wenn
schon Schützen angemeldet sind. Jedes abgeschlossene Programm bietet dann in
den Resultaten **Auf Vorlage drucken** an: am Desk, auf der Seite
**Resultate erfassen** und über Warner-Links an der Scheibe (siehe
[Helfer](helpers.md)). Mit einem Druck-Agenten und einem Drucker für
*Druckvorlagen* unter **Schiessstände → Drucken** geht das Blatt an diesen
Drucker (z. B. an seinen manuellen Einzug), sonst öffnet sich das PDF (siehe
[Drucken](printing.md)).

**Exportieren** lädt eine Vorlage als Datei herunter und **Importieren** liest sie wieder ein, z. B. um eine Vorlage mit einem anderen Verein zu teilen.
