# Anlässe

**Anlässe** listet die Anlässe des Vereins in den Gruppen *Heute* (zwischen dem ersten und dem letzten Schiesstag), *Bevorstehend*, *Vergangen* und *Abgeschlossen*; Admins erstellen dort neue. Das Öffnen eines Anlasses zeigt in der Seitenleiste sein Dashboard, Reservationen, Ranglisten, Statistik und Einrichtung; *Anlass wechseln* oben in der Seitenleiste führt zurück zur Liste.

![Anlässe](../images/screenshots/events.png){ .screenshot }

## Einstellungen {#settings}

Die **Einrichtung** des Anlasses hat fünf Tabs: *Programme* (siehe
[Programme](programmes.md)), *Anlass*, *Kiosk-Anzeigen* (siehe [Ranglisten
und Kiosk](rankings-kiosk.md)), *Helfer* (siehe [Helfer](helpers.md)) und
*Kopieren & abschliessen*. Mitarbeiter sehen nur die Programme.

Der Tab *Anlass* enthält die Einstellungen (Admins):

| Einstellung | Bedeutung |
|---|---|
| Name, Datum, Schiessstände | das Datum ist der Haupttag; sein Jahr bestimmt die Alterskategorien (SSV: Wettkampfjahr minus Jahrgang, der Tag selbst spielt keine Rolle, was für Vorschiesstage im Dezember vor einem Anlass im Januar zählt); jeden Stand ankreuzen, auf dem der Anlass geschossen wird; oder keinen für einen Anlass ohne Scheiben (siehe unten) |
| Schiesstage | einzelne Tage mit optionalen Schiesszeiten, für Vorschiesstage und mehrtägige Anlässe (Anmeldungen zählen am Tag, an dem sie gemacht werden); der Haupttag ist immer einer |
| Zeitraum (Saison) | für Anlässe über Wochen oder Monate, z. B. das Feldschiessen von der Öffnung des Schützenhauses bis 31. August: erster und letzter Tag, dazu die Wochentage mit ihren Schiesszeiten; kein Wochentag angekreuzt heisst jeder Tag. Einzelne Tage, z. B. das offizielle Wochenende, ergänzen den Zeitraum oder ersetzen an diesem Datum seine Zeiten |
| Veranstalter | wird auf Dokumente gedruckt |
| Sprache der Dokumente und des Kiosks | Sprache von PDFs, CSV-Exporten und Kiosk-Bildschirmen, unabhängig von der Oberflächensprache der einzelnen Benutzer |
| Von der Vereinskonkurrenz ausgeschlossener Verein | meist der Gastgeberverein; aufgeführt als *ausser Konkurrenz* |
| Dokumente bei der Anmeldung | **Standblatt** und **Kontrollstreifen**: ausschalten, was der Anlass nicht braucht; die Anmeldung bietet es dann weder an noch druckt sie es |
| Scheibenstatus auf den Ranglisten zeigen | Kiosk-Bildschirme zeigen die Scheiben |
| Wert der Auszeichnung auf Kranzkarten drucken | Wert in Ziffern und Worten auf der Karte |
| Der Veranstalter gibt die Munition ab | zeigt den Knopf **Munition** im Dashboard und für die Munitionsausgabe (siehe [Helfer](helpers.md)); aus bei Anlässen, an denen die Schützen ihre eigene mitbringen, z. B. Luftpistole und Luftgewehr |

Die Schiesszeiten geben auch der Online-Reservation ihre Zeitblöcke, ausser
die Reservation legt für einen Tag andere Zeiten fest (siehe
[Online-Reservation](online-booking.md)).

## Anlässe ohne Scheiben {#no-lanes}

Ein Anlass braucht keinen Stand, z. B. ein Jahresprogramm oder ein
Feldschiessen, dessen Resultate auf Papier oder aus der Trefferanzeige
kommen. Ohne Stand gibt es kein Scheibenraster, keine Reservationen und
keine Online-Buchung; Schützen werden ohne Scheibe angemeldet, und ihre
Resultate werden wie gewohnt erfasst. Dokumente lassen sich weiterhin als
PDF öffnen, gehen aber nicht an einen Druck-Agenten, weil Druck-Agenten zu
einem Stand gehören.

## Für nächstes Jahr kopieren

Auf dem Tab *Kopieren & abschliessen* erstellt **Für nächstes Jahr
kopieren** einen neuen Anlass mit denselben Ständen, Organisator, Sprache,
Einstellungen und allen Programmen, inklusive Preise, Auszeichnungen,
Kombinationen und Spezialpreise; die Schiesstage und der Zeitraum wandern
mit dem Datum, wöchentliche Schiesszeiten bleiben auf ihren
Wochentagen. Anmeldungen, Resultate und Kiosk-Links werden nicht
kopiert. Der Dialog schlägt den Namen mit hochgezähltem Jahr (und einer
führenden Nummer wie *42.*) und das Datum ein Jahr später vor; passe beides
bei Bedarf an.

![Einen Anlass kopieren](../images/screenshots/copy-event.png){ .screenshot
}

## Abschliessen {#close}

Wenn der Anlass vorbei ist und die Ranglisten geprüft sind, **schliesst** du
ihn auf dem Tab *Kopieren & abschliessen* **ab** (**Anlass
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

**Löschen** (auf dem Tab *Kopieren & abschliessen*) entfernt einen Anlass, aber nur solange niemand dafür angemeldet ist.
