# Wettkampfregeln

Diese Seite beschreibt die Regeln genau so, wie Standblatt sie umsetzt.

## Alterskategorie {#age-categories}

Die Kategorie ergibt sich aus *Jahr des Anlasses − Jahrgang*. Ohne
Geburtsdatum zählt ein Schütze als Elite. Die Kategorie wird bei der
Anmeldung festgelegt.

| Kategorie | Kurz | Alter im Jahr des Anlasses |
|---|---|---|
| U10 | U10 | bis 9 |
| U13 | U13 | 10–12 |
| U15 | U15 | 13–14 |
| U17 | U17 | 15–16 |
| U19 | U19 | 17–18 |
| U21 | U21 | 19–20 |
| Elite | E | 21–45 |
| Senioren | S | 46–59 |
| Veteranen | V | 60–69 |
| Seniorveteranen | SV | 70 und älter |

Das Jahr des Anlasses ist das Jahr seines Haupttags, auch für Schiesstage in
einem anderen Jahr.

## Anmeldung

Eine Anmeldung wird abgelehnt, wenn:

- ein Programm nicht zum Anlass gehört, doppelt gewählt oder nicht aktiv
  ist, oder mehr Durchgänge gewählt sind, als es erlaubt;
- ein Programm keinen Preis für die Alterskategorie des Schützen hat („nicht
  für diese Alterskategorie“);
- zwei Varianten desselben Programms gewählt sind;
- ein Programm Passen aus einem anderen Programm übernimmt, das nicht
  gewählt ist oder die Anzahl gewählter Passen nicht stimmt;
- die Scheibe von einer anderen aktiven Anmeldung belegt oder ausser Betrieb
  ist;
- der Schütze an diesem Anlass schon für eines der Programme angemeldet ist
  (ein Schütze schiesst jedes Programm einmal pro Anlass, ausser es erlaubt
  mehrere [Durchgänge](../user-guide/programmes.md#attempts); eine andere
  Variante in einer separaten Anmeldung ist erlaubt).

Eine reservierte Scheibe lässt sich nach einer Bestätigung trotzdem
verwenden.

## Resultate

Jede Passe muss erfasst sein. Punkte liegen zwischen 0 und dem Passenmaximum
(standardmässig Schüsse × 10), Mouchen zwischen 0 und der Anzahl Schüsse,
und nur bei Passen, die Mouchen erfassen. Beim Speichern der Resultate wird
die wertvollste erreichte Auszeichnung vergeben; erneutes Speichern behält
die gewählte Auszeichnung, solange das neue Total sie noch erreicht, sonst
die wertvollste erreichte. Sobald jedes Programm einer Anmeldung Resultate
hat, ist die Anmeldung abgeschlossen und ihre Scheibe wird frei.

### Kombinationen

Ein Programm, das als *kombinierbar mit* einem anderen Programm markiert
ist, kann Passen daraus übernehmen: Bei der Anmeldung wählt der Schütze,
welche Passen (zum Beispiel Passe 2 und 3 eines Programms mit vier Passen
für ein Programm mit zwei Passen). Die gewählten Passen werden in
aufsteigender Reihenfolge kopiert und ab 1 neu nummeriert; das kombinierte
Programm lässt sich nicht für sich bearbeiten und folgt jeder Änderung der
Quelle.

## Auszeichnungen

Eine Auszeichnungsstufe steht einer Alterskategorie nur offen, wenn sie
dafür eine Limite hat. Ein Schütze kann jede Stufe wählen, deren Limite sein
Total erreicht; die wertvollste Stufe steht zuoberst und wird beim Speichern
der Resultate standardmässig vergeben. Der Server prüft die Wahl.

## Rangierung {#ranking}

Alle Resultate eines Programms werden geordnet nach:

1. höherem Total;
2. den Passen, verglichen von der letzten rückwärts (die bessere letzte
   Passe gewinnt, dann die zweitletzte, …);
3. mehr Mouchen;
4. Alter: bei Junioren (U10–U21) zuerst der jüngere Schütze, bei den übrigen
   zuerst der ältere; Nicht-Junioren stehen vor Junioren.

Die Ränge sind fortlaufend (1, 2, 3, … auch bei gleichen
Resultaten). Dieselbe Rangierung speist die Seiten des Personals, die
Kiosk-Bildschirme, das Ranglisten-PDF und -CSV sowie die Sonderpreise.

## Vereinskonkurrenz {#club-competition}

Nur Programme mit *Zählt für die Vereinskonkurrenz* machen mit, und nur
abgeschlossene Resultate von Schützen mit Verein.

- Der Wert eines Resultats ist die Summe seiner ersten zwei Passen.
- Ein Verein mit weniger als 5 Teilnehmern wird als *nicht rangiert*
  aufgeführt.
- Von *n* Teilnehmern sind die besten *max(5, ⌊n/2⌋)* Pflicht, die übrigen
  freiwillig.
- Resultat = ⌊(Σ Pflicht + 0.02 × Σ freiwillig) / Anzahl Pflicht × 1000⌋ /
  1000, also auf drei Nachkommastellen abgeschnitten.
- Gleichstand: mehr Teilnehmer zuerst, dann die sortierten Werte einzeln
  verglichen.
- Der in den Einstellungen des Moduls *Vereinskonkurrenz* ausgeschlossene
  Verein (meist der Gastgeber) behält seinen Platz in der Liste, erhält aber
  keinen Rang („ausser Konkurrenz“).

## Sonderpreise {#special-prizes}

Ein Programm kann Sonderpreise an die besten Ränge jeder Alterskategorie
vergeben. Drei Zahlen legen sie fest:

- **Ränge pro Kategorie**, z. B. 3;
- **Mindestteilnehmer pro Kategorie** in diesem Programm, z. B. 10;
- **Mindestteilnehmer des Vereins des Schützen** im selben Programm,
  z. B. 4.

Pro Alterskategorie mit genügend Teilnehmern gewinnen die ersten Schützen in
Ranglistenreihenfolge, deren Verein genügend Teilnehmer hat; Schützen, deren
Verein zu wenige hat, werden übersprungen, und der nächste rückt
nach. Schützen ohne Verein sind nicht berechtigt, ausser das Vereinsminimum
ist 0. Der Bericht listet sowohl die Gewinner als auch die übersprungenen
Schützen mit dem Grund.

!!! example "Von der Ausschreibung zu den Einstellungen"
    In einer Ausschreibung steht: *„In jeder Kategorie erhalten die Ränge 1–3 einen Sonderpreis, sofern die Kategorie mindestens 10 Teilnehmer und der Verein des Schützen mindestens 4 hat.“* Das sind Ränge 3, Kategorie-Minimum 10, Vereins-Minimum 4.

## Scheibenlisten

Scheiben werden als Liste von Zahlen und Bereichen geschrieben: `1, 5,
7-10`. Standblatt gibt sie kompakt aus und macht aus Folgen von drei oder
mehr einen Bereich (`1-3, 5, 7, 8`).

## Geld und Zahlen

Beträge werden als `CHF 1'234.50` angezeigt, Daten als `dd.mm.yyyy`. Auf
Kranzkarten kann der Wert in Worten in der Sprache des Anlasses gedruckt
werden, z. B. *Fünfundzwanzig Franken*, *Septante et un francs*, *Ventitré
franchi*.
