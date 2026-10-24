# Anlässe

**Anlässe** listet die Anlässe des Vereins in den Gruppen *Heute* (zwischen dem ersten und dem letzten Schiesstag), *Bevorstehend*, *Vergangen* und *Abgeschlossen*; Rollen mit dem Recht *Anlässe konfigurieren* erstellen dort neue. Das Öffnen eines Anlasses zeigt in der Seitenleiste sein Dashboard und die übrigen Seiten des Betriebs; *Anlass wechseln* oben in der Seitenleiste führt zurück zur Liste.

![Anlässe](../images/screenshots/events.png){ .screenshot }

## Konfiguration und Betrieb {#configuration}

Ein Anlass hat zwei Bereiche. Der **Betrieb** ist das, womit du am Anlass
arbeitest: *Dashboard*, *Resultate erfassen*, *Ranglisten*, *Statistik* und
die Seiten eingeschalteter Module wie *Reservationen*. **Konfiguration** am
Ende der Seitenleiste (nur Admins) öffnet die Einrichtung des Anlasses; die
Seitenleiste zeigt dann deren Seiten, und *Zurück zum Anlass* führt zurück
in den Betrieb:

| Seite | Was sie enthält |
|---|---|
| **Allgemein** | Name, Tage, Stände, Veranstalter und Sprache (siehe unten) |
| **Module** | die Module, die der Anlass verwendet (siehe unten) |
| **Programme** | was geschossen werden kann (siehe [Programme](programmes.md)) |
| Seiten der Module | z. B. **Reservationen** (mit der Online-Reservation), **Wertungen**, **Kiosk-Anzeigen**, **Helferlinks**, nur solange das Modul an ist |
| **Kopieren & abschliessen** | für nächstes Jahr kopieren, abschliessen, wiedereröffnen und löschen |

Die Konfiguration braucht das Recht *Anlässe konfigurieren*, der Betrieb
*Anlässe durchführen* (siehe [Rollen](users.md#roles)). Warum beides
getrennt ist: [Konfiguration und
Betrieb](../explanation/modules.md#configuration).

## Einstellungen {#settings}

**Allgemein** enthält, was jeder Anlass hat:

| Einstellung | Bedeutung |
|---|---|
| Name, Datum, Schiessstände | das Datum ist der Haupttag; sein Jahr bestimmt die Alterskategorien (SSV: Wettkampfjahr minus Jahrgang, der Tag selbst spielt keine Rolle, was für Vorschiesstage im Dezember vor einem Anlass im Januar zählt); jeden Stand ankreuzen, auf dem der Anlass geschossen wird; oder keinen für einen Anlass ohne Scheiben (siehe unten) |
| Schiesstage | einzelne Tage mit optionalen Schiesszeiten, für Vorschiesstage und mehrtägige Anlässe (Anmeldungen zählen am Tag, an dem sie gemacht werden); der Haupttag ist immer einer |
| Zeitraum (Saison) | für Anlässe über Wochen oder Monate, z. B. das Feldschiessen von der Öffnung des Schützenhauses bis 31. August: erster und letzter Tag, dazu die Wochentage mit ihren Schiesszeiten; kein Wochentag angekreuzt heisst jeder Tag. Einzelne Tage, z. B. das offizielle Wochenende, ergänzen den Zeitraum oder ersetzen an diesem Datum seine Zeiten |
| Veranstalter | wird auf Dokumente gedruckt |
| Sprache der Dokumente und des Kiosks | Sprache von PDFs, CSV-Exporten und Kiosk-Bildschirmen, unabhängig von der Oberflächensprache der einzelnen Benutzer |

Die Schiesszeiten geben auch der Online-Reservation ihre Zeitblöcke, ausser
die Reservation legt für einen Tag andere Zeiten fest (siehe
[Online-Reservation](online-booking.md)).

Schritt für Schritt für eine Saison: [Ein Jahresprogramm über die Saison
durchführen](../how-to/annual-programme.md), [Ein Feldschiessen über einen
Zeitraum einrichten](../how-to/feldschiessen.md).

## Module {#modules}

**Module** schaltet die Bausteine des Anlasses ein und aus; warum sich Anlässe nur durch ihre Module unterscheiden, erklärt [Module und Arten von Anlässen](../explanation/modules.md).

![Die Module eines Anlasses](../images/screenshots/modules.png){ .screenshot
}

Jedes Modul ist eine Karte mit seinem Schalter; die Einstellungen eines
Moduls, z. B. *Wert der Auszeichnung auf Kranzkarten drucken* bei
*Auszeichnungen*, stehen auf seiner Karte. Was jedes Modul bringt, seine
Einstellungen und wann es an bleibt:
[Module](../reference/modules.md#event-modules).

Beim Erstellen eines Anlasses wählst du, was für ein Anlass es ist; das
schaltet nur Module ein (siehe [Arten von
Anlässen](../reference/modules.md#kinds)).

Ein Modul mit Daten bleibt an, und die Karte sagt, warum ihr Schalter
gesperrt ist (siehe [Warum ein Modul mit Daten an
bleibt](../explanation/modules.md#in-use)). Ein Modul, das ein anderes
braucht, z. B. *Auszeichnungen* für *Sonderpreise* oder *Scheibenzuteilung*
für *Reservationen*, bleibt an, solange das andere an ist.

### Vereinsmodule {#club-modules}

Manche Module gehören zum Verein statt zu einem Anlass. **Vereinsmodule** im
Vereinsteil der Seitenleiste (Recht *Benutzer und Verein*) schaltet sie ein
und aus; vorerst sind das die *Druckvorlagen* (siehe
[Druckvorlagen](print-forms.md)). Was jedes bringt und wann es an bleibt:
[Vereinsmodule](../reference/modules.md#club-modules); warum sie zum Verein
gehören: [Vereinsmodule und
Anlassmodule](../explanation/modules.md#club-modules).

## Anlässe ohne Scheiben {#no-lanes}

Ein Anlass ohne Scheiben hat keinen Schiessstand oder schaltet die
*Scheibenzuteilung* aus; was jede Wahl bedeutet, erklärt [Anlässe mit und
ohne Scheiben](../explanation/modules.md#lanes).

## Für nächstes Jahr kopieren {#copy}

Auf **Kopieren & abschliessen** in der Konfiguration erstellt **Kopieren**
unter *Für nächstes Jahr kopieren* einen neuen Anlass mit denselben
Schiessständen, Organisator, Sprache, Einstellungen, Modulen und allen
Programmen, inklusive Preise, Auszeichnungen, Kombinationen und
Sonderpreise, sowie den [Wertungen](classifications.md); die Schiesstage und
der Zeitraum wandern mit dem Datum, wöchentliche Schiesszeiten bleiben auf
ihren Wochentagen. Anmeldungen, Resultate, Kiosk-Links, Helferlinks und der
Reservationslink werden nicht kopiert. Der Dialog *Anlass kopieren* schlägt
den Namen und das Datum ein Jahr später vor; **Kopie erstellen** erstellt
den Anlass. Schritt für Schritt: [Den Anlass des letzten Jahres
übernehmen](../how-to/next-year.md#copy).

![Einen Anlass kopieren](../images/screenshots/copy-event.png){ .screenshot
}

## Vorlagen {#templates}

Ein gut eingerichteter Anlass kann zu einer **Vorlage** des Vereins werden:
**Als Vorlage speichern** auf **Kopieren & abschliessen** legt eine Kopie
mit seinen Modulen, Einstellungen, Schiessständen, Programmen und Wertungen
an, benannt nach dem Anlass ohne Jahreszahl. Vorlagen stehen unter
*Vorlagen* unterhalb der Anlässe; eine Vorlage öffnet sich in ihrer
Konfiguration und lässt sich wie ein Anlass ändern, nimmt aber keine
Anmeldungen an. **Neuer Anlass** listet die Vorlagen des Vereins unter den
Arten auf (siehe [Mit einer Vorlage des Vereins
starten](../how-to/next-year.md#template)). **Löschen** auf **Kopieren &
abschliessen** der Vorlage entfernt sie.

![Neuer Anlass mit den Arten und einer Vorlage des
Vereins](../images/screenshots/new-event.png){ .screenshot }

## Abschliessen {#close}

Wenn ein Anlass vorbei ist und die Ranglisten geprüft sind, schliesst ihn
**Anlass abschliessen** auf **Kopieren & abschliessen** in der Konfiguration
ab; ein Banner auf jeder Seite des Anlasses sagt das dann, und
**Wiedereröffnen** macht ihn wieder bearbeitbar. Was das Abschliessen für
die Ranglisten, Mein Training und Änderungen bedeutet: [Einen Anlass
abschliessen](../explanation/closing.md).

## Löschen

**Löschen** (auf **Kopieren & abschliessen**) entfernt einen Anlass, aber nur solange niemand dafür angemeldet ist. Ein abgeschlossener Anlass lässt sich nicht löschen; eröffne ihn zuerst wieder.
