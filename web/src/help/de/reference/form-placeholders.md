# Platzhalter in Vorlagen

Eine [Druckvorlage](../user-guide/print-forms.md) druckt anstelle eines
Platzhalters wie `{{shooter.name}}` dessen Wert. Ein Platzhalter ohne Wert
druckt nichts; ebenso einer, den die
[Datenquelle](../user-guide/print-forms.md#source) der Vorlage nicht
kennt. In einem wiederholten Element zählt `#` 1, 2, 3, …: `{{s1.shot#}}`
druckt `{{s1.shot1}}`, `{{s1.shot2}}`, ….

## Jede Datenquelle {#common}

| Platzhalter | Feld | Beispiel |
|---|---|---|
| `{{shooter.name}}` | Name (Nachname zuerst) | Example Anna |
| `{{shooter.lastname}}` | Nachname | Example |
| `{{shooter.firstname}}` | Vorname | Anna |
| `{{shooter.license}}` | Lizenznummer | 123456 |
| `{{shooter.birthday}}` | Geburtsdatum | 02.06.1985 |
| `{{shooter.birthyear}}` | Jahrgang | 1985 |
| `{{shooter.postalcode}}` | PLZ | 9999 |
| `{{shooter.city}}` | Ort | Hilltown |
| `{{category}}` | Alterskategorie | Elite |
| `{{club}}` | Verein | SG Hilltown |
| `{{club.ssv}}` | SSV-Nummer | 1.23.4.56 |
| `{{event.name}}` | Anlass | Jahresprogramm 2027 |
| `{{event.date}}` | Datum (Haupttag) | 12.05.2027 |
| `{{event.organizer}}` | Organisator | SG Hilltown |
| `{{today}}` | Heute, in der Zeitzone des Anlasses | 12.05.2027 |

## Resultate eines Programms {#results}

Vorlagen von Programmen und Kranzkarten.

| Platzhalter | Feld | Beispiel |
|---|---|---|
| `{{programme}}` | Programm, mit seiner Variante | Obligatorisch |
| `{{lane}}` | Scheibe | 7 |
| `{{registration}}` | Anmeldenummer | 1042 |
| `{{total}}` | Total | 85 |
| `{{mouches}}` | Mouches | 3 |
| `{{s1}}`, `{{s2}}`, … | Passe 1, 2, … | 47 |
| `{{s1.mouches}}`, … | Mouches Passe 1, … | 1 |
| `{{s1.shot1}}`, `{{s1.shot#}}`, … | Einzelschüsse der Passe 1 (Passen mit **Einzelschüsse**) | 5 |
| `{{award}}` | Auszeichnung | Karte |
| `{{award.value}}` | Wert der Auszeichnung | CHF 12.00 |
| `{{award.value.words}}` | Wert in Worten, in der Sprache des Anlasses | Zwölf Franken |

`{{award.value}}` und `{{award.value.words}}` bleiben leer, ausser *Wert der
Auszeichnung auf Kranzkarten drucken* ist an und die Auszeichnung hat einen
Wert.

## Eine Anmeldung {#registration}

Kontrollstreifen und Anmeldedokumente.

| Platzhalter | Feld | Beispiel |
|---|---|---|
| `{{category.short}}` | Alterskategorie, kurz | E |
| `{{lane}}` | Scheibe | 7 |
| `{{registration}}` | Anmeldenummer | 1042 |
| `{{registration.code}}` | QR-Code der Anmeldung (der [Anmeldecode](codes.md)) | REGCODEAAAA1042AAAAL123456 |
| `{{programmes}}` | Programme, durch Komma getrennt | Obligatorisch, Feldschiessen |
| `{{passes}}` | Passen: `P1 P2` pro Programm, `P0` für ein Programm, das seine Passen aus einem anderen übernimmt, `*` bei den übernommenen Passen | P1 P2 + P0 |

## Ein Rang {#rank}

Urkunden.

| Platzhalter | Feld | Beispiel |
|---|---|---|
| `{{rank}}` | Rang | 1 |
| `{{total}}` | Total | 98 |
| `{{ranking}}` | Rangliste: das Programm oder die Wertung | Sektionsmeisterschaft |
| `{{class}}` | Klasse der Wertung; leer bei einem Programm | Junioren |

Eine Urkunde eines Programms füllt auch jeden Platzhalter von [Resultate
eines Programms](#results) aus dem zählenden Durchgang.

## Elemente {#elements}

| Element | Druckt |
|---|---|
| *Eigener Text* | seinen Text, mit allfälligen Platzhaltern |
| *QR-Code* | seinen Wert als QR-Code |
| *Barcode* | seinen Wert als Barcode Code 128 |
| *Linie* | eine Linie |
| Bild | ein Bild aus der [Bildbibliothek](../user-guide/print-forms.md#images) des Vereins |
