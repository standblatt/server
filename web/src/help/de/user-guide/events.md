# Anlässe

**Anlässe** listet die Anlässe des Vereins in den Gruppen *Heute* (zwischen dem ersten und dem letzten Schiesstag), *Bevorstehend*, *Vergangen* und *Abgeschlossen*; Admins erstellen dort neue. Das Öffnen eines Anlasses zeigt in der Seitenleiste sein Dashboard und die übrigen Seiten des Betriebs; *Anlass wechseln* oben in der Seitenleiste führt zurück zur Liste.

![Anlässe](../images/screenshots/events.png){ .screenshot }

## Konfiguration und Betrieb {#configuration}

Ein Anlass hat zwei Bereiche. Der **Betrieb** ist das, womit du am Anlass
arbeitest: *Dashboard*, *Resultate erfassen*, *Ranglisten*, *Statistik* und
die Seiten eingeschalteter Bausteine wie *Reservationen*. **Konfiguration**
am Ende der Seitenleiste (nur Admins) öffnet die Einrichtung des Anlasses;
die Seitenleiste zeigt dann deren Seiten, und *Zurück zum Anlass* führt
zurück in den Betrieb:

| Seite | Was sie enthält |
|---|---|
| **Allgemein** | Name, Tage, Stände, Veranstalter und Sprache (siehe unten) |
| **Bausteine** | die Bausteine, die der Anlass verwendet (siehe unten) |
| **Programme** | was geschossen werden kann (siehe [Programme](programmes.md)) |
| Seiten der Bausteine | z. B. **Wertungen**, **Kiosk-Anzeigen**, **Helferlinks**, **Online-Reservation**, nur solange der Baustein an ist |
| **Kopieren & abschliessen** | für nächstes Jahr kopieren, abschliessen, wiedereröffnen und löschen |

Personal arbeitet nur im Betrieb; die Konfiguration ist für Admins.

## Einstellungen {#settings}

**Allgemein** enthält, was jeder Anlass hat:

| Einstellung | Bedeutung |
|---|---|
| Name, Datum, Schiessstände | das Datum ist der Haupttag; sein Jahr bestimmt die Alterskategorien (SSV: Wettkampfjahr minus Jahrgang, der Tag selbst spielt keine Rolle, was für Vorschiesstage im Dezember vor einem Anlass im Januar zählt); jeden Stand ankreuzen, auf dem der Anlass geschossen wird; oder keinen für einen Anlass ohne Scheiben (siehe unten) |
| Schiesstage | einzelne Tage mit optionalen Schiesszeiten, für Vorschiesstage und mehrtägige Anlässe (Anmeldungen zählen am Tag, an dem sie gemacht werden); der Haupttag ist immer einer |
| Zeitraum (Saison) | für Anlässe über Wochen oder Monate, z. B. das Feldschiessen von der Öffnung des Schützenhauses bis 31. August: erster und letzter Tag, dazu die Wochentage mit ihren Schiesszeiten; kein Wochentag angekreuzt heisst jeder Tag. Einzelne Tage, z. B. das offizielle Wochenende, ergänzen den Zeitraum oder ersetzen an diesem Datum seine Zeiten |
| Veranstalter | wird auf Dokumente gedruckt |
| Sprache der Dokumente und des Kiosks | Sprache von PDFs, CSV-Exporten und Kiosk-Bildschirmen, unabhängig von der Oberflächensprache der einzelnen Benutzer |

## Bausteine {#modules}

Nicht jeder Anlass braucht alles: Ein Schützenfest reserviert Scheiben und
gibt Kranzkarten ab, ein Jahresprogramm sammelt nur Resultate. **Bausteine**
schaltet die Bausteine des Anlasses ein und aus; was aus ist, verschwindet
aus Seitenleiste, Dashboard, Programm-Editor, Ranglisten und Berichten.

| Baustein | Was er bringt |
|---|---|
| Reservationen | für Gruppen und Vereine reservierte Scheiben (braucht einen Stand) |
| Online-Reservation | Vereine fragen über einen öffentlichen Link Scheiben an (braucht einen Stand und schaltet Reservationen ein, siehe [Online-Reservation](online-booking.md)) |
| Standblatt, Kontrollstreifen | die Dokumente, die bei der Anmeldung gedruckt werden |
| Munitionsausgabe | der Knopf **Munition** im Dashboard und für die Munitionsausgabe (siehe [Helfer](helpers.md#ammunition)); aus lassen, wo die Schützen ihre eigene mitbringen, z. B. Luftpistole und Luftgewehr |
| Durchgänge | Programme, die mehrmals geschossen werden können (siehe [Programme](programmes.md#attempts)) |
| Auszeichnungen | Kranzkarten, Naturalgaben und Geldpreise; seine Einstellung *Wert der Auszeichnung auf Kranzkarten drucken* setzt den Wert in Ziffern und Worten auf die Karte |
| Sonderpreise | die besten Ränge jeder Alterskategorie (schaltet Auszeichnungen ein, siehe [Programme](programmes.md#special-prizes)) |
| Vereinskonkurrenz | Vereine nach den Resultaten ihrer Schützen rangiert; seine Einstellung *Von der Vereinskonkurrenz ausgeschlossener Verein* ist meist der Gastgeberverein, aufgeführt als *ausser Konkurrenz* |
| Wertungen | Ranglisten über mehrere Programme oder Durchgänge (siehe [Wertungen](classifications.md)) |
| Kiosk-Anzeigen | Ranglisten auf Bildschirmen in der Halle (siehe [Ranglisten und Kiosk](rankings-kiosk.md)) |
| Helferlinks | Warner, Munitionsausgabe und Resultaterfassung ohne Konto (siehe [Helfer](helpers.md)) |

Beim Erstellen eines Anlasses wählst du, was für ein Anlass es ist; das
schaltet nur Bausteine ein:

| Art | Bausteine |
|---|---|
| Schützenfest | Reservationen, Online-Reservation, Standblatt, Kontrollstreifen, Auszeichnungen, Sonderpreise, Vereinskonkurrenz, Kiosk-Anzeigen, Helferlinks |
| Jahresprogramm oder Training | Durchgänge, Wertungen |
| Feldschiessen | Auszeichnungen, Vereinskonkurrenz, Kiosk-Anzeigen, Helferlinks |
| Vereinsmeisterschaft | Durchgänge, Wertungen, Kiosk-Anzeigen |
| Eigene Auswahl | keine |

Ein Baustein mit Daten bleibt an: Den Kiosk schaltest du erst aus, wenn
seine Kiosk-Links gelöscht sind, die Auszeichnungen erst, wenn kein Programm
mehr Auszeichnungsstufen hat, und so weiter; die Karte sagt, warum ihr
Schalter gesperrt ist. Ein Baustein, den ein anderer braucht (Reservationen
für die Online-Reservation), bleibt an, solange jener an ist.

Die Schiesszeiten geben auch der Online-Reservation ihre Zeitblöcke, ausser
die Reservation legt für einen Tag andere Zeiten fest (siehe
[Online-Reservation](online-booking.md)).

### Vereinsbausteine {#club-modules}

Manche Bausteine gehören zum Verein statt zu einem
Anlass. **Vereinsbausteine** im Vereinsteil der Seitenleiste (Admins)
schaltet sie ein und aus:

| Baustein | Was er bringt |
|---|---|
| Druckvorlagen | Vorlagen für vorgedruckte Blätter und der Abschnitt *Druckvorlage* jedes Programms (siehe [Druckvorlagen](print-forms.md)) |

## Anlässe ohne Scheiben {#no-lanes}

Ein Anlass braucht keinen Stand, z. B. ein Jahresprogramm oder ein
Feldschiessen, dessen Resultate auf Papier oder aus der Trefferanzeige
kommen. Ohne Stand gibt es kein Scheibenraster, und die Bausteine
Reservationen und Online-Reservation stehen nicht zur Verfügung; Schützen
werden ohne Scheibe angemeldet, und ihre Resultate werden wie gewohnt
erfasst. Dokumente lassen sich weiterhin als PDF öffnen, gehen aber nicht an
einen Druck-Agenten, weil Druck-Agenten zu einem Stand gehören.

## Für nächstes Jahr kopieren {#copy}

Auf **Kopieren & abschliessen** in der Konfiguration erstellt **Für nächstes
Jahr kopieren** einen neuen Anlass mit denselben Ständen, Veranstalter,
Sprache, Einstellungen, Bausteinen und allen Programmen, inklusive Preise,
Auszeichnungen, Kombinationen und Sonderpreise, sowie den
[Wertungen](classifications.md); die Schiesstage und der Zeitraum wandern
mit dem Datum, wöchentliche Schiesszeiten bleiben auf ihren
Wochentagen. Anmeldungen, Resultate und Kiosk-Links werden nicht
kopiert. Der Dialog schlägt den Namen mit hochgezähltem Jahr (und einer
führenden Nummer wie *42.*) und das Datum ein Jahr später vor; passe beides
bei Bedarf an.

![Einen Anlass kopieren](../images/screenshots/copy-event.png){ .screenshot
}

## Abschliessen {#close}

Wenn der Anlass vorbei ist und die Ranglisten geprüft sind, **schliesst** du
ihn auf **Kopieren & abschliessen** in der Konfiguration **ab** (**Anlass
abschliessen**). Das Abschliessen gibt die Ranglisten frei: Schützen mit
einem [persönlichen Trainingslog](../my-training/account.md#competitions)
finden jetzt ihre Resultate dieses Anlasses mit ihrem Rang im Programm und
können sie in ihr Log übernehmen. Alles bleibt sichtbar (Dashboard,
Resultate, Ranglisten, Berichte, Dokumente), aber Anmeldungen, Resultate,
Auszeichnungen, Programme und Einstellungen können nicht mehr geändert
werden. Ein Banner auf jeder Seite des Anlasses weist darauf
hin. **Wiedereröffnen** macht ihn wieder bearbeitbar; Resultate, die
Schützen bereits in ihr Log übernommen haben, bleiben dort. Das Kopieren
eines abgeschlossenen Anlasses für nächstes Jahr funktioniert.

## Löschen

**Löschen** (auf **Kopieren & abschliessen**) entfernt einen Anlass, aber nur solange niemand dafür angemeldet ist.
