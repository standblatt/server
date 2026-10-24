# Form placeholders

A [print form](../user-guide/print-forms.md) prints the value of a placeholder such as `{{shooter.name}}` in place of it. A placeholder without a value prints nothing; so does one the form's [data source](../user-guide/print-forms.md#source) does not know. Inside a repeated element, `#` counts 1, 2, 3, …: `{{s1.shot#}}` prints `{{s1.shot1}}`, `{{s1.shot2}}`, ….

## Every data source {#common}

| Placeholder | Field | Example |
|---|---|---|
| `{{shooter.name}}` | Name (last name first) | Example Anna |
| `{{shooter.lastname}}` | Last name | Example |
| `{{shooter.firstname}}` | First name | Anna |
| `{{shooter.license}}` | License number | 123456 |
| `{{shooter.birthday}}` | Birthday | 02.06.1985 |
| `{{shooter.birthyear}}` | Year of birth | 1985 |
| `{{shooter.postalcode}}` | Postal code | 9999 |
| `{{shooter.city}}` | City | Hilltown |
| `{{category}}` | Age category | Elite |
| `{{club}}` | Club | SG Hilltown |
| `{{club.ssv}}` | SSV number | 1.23.4.56 |
| `{{event.name}}` | Event | Annual programme 2027 |
| `{{event.date}}` | Date (main day) | 12.05.2027 |
| `{{event.organizer}}` | Organizer | SG Hilltown |
| `{{today}}` | Today, in the event's time zone | 12.05.2027 |

## Results of a programme {#results}

Programme forms and wreath cards.

| Placeholder | Field | Example |
|---|---|---|
| `{{programme}}` | Programme, with its variant | Obligatory |
| `{{lane}}` | Lane | 7 |
| `{{registration}}` | Registration number | 1042 |
| `{{total}}` | Total | 85 |
| `{{mouches}}` | Mouches | 3 |
| `{{s1}}`, `{{s2}}`, … | Series 1, 2, … | 47 |
| `{{s1.mouches}}`, … | Mouches series 1, … | 1 |
| `{{s1.shot1}}`, `{{s1.shot#}}`, … | single shots of series 1 (series with **Single shots**) | 5 |
| `{{award}}` | Award | Card |
| `{{award.value}}` | Value of the award | CHF 12.00 |
| `{{award.value.words}}` | Value in words, in the event's language | twelve francs |

`{{award.value}}` and `{{award.value.words}}` stay empty unless *Print the award value on wreath cards* is on and the award has a value.

## A registration {#registration}

Control strips and registration documents.

| Placeholder | Field | Example |
|---|---|---|
| `{{category.short}}` | Age category, short | E |
| `{{lane}}` | Lane | 7 |
| `{{registration}}` | Registration number | 1042 |
| `{{registration.code}}` | QR code of the registration (the [registration code](codes.md)) | REGCODEAAAA1042AAAAL123456 |
| `{{programmes}}` | Programmes, comma separated | Obligatory, Field match |
| `{{passes}}` | Passes: `P1 P2` per programme, `P0` for a programme that takes its series from another one, `*` on the series that are taken | P1 P2 + P0 |

## A rank {#rank}

Certificates.

| Placeholder | Field | Example |
|---|---|---|
| `{{rank}}` | Rank | 1 |
| `{{total}}` | Total | 98 |
| `{{ranking}}` | Ranking: the programme or the classification | Section championship |
| `{{class}}` | Class of the classification; empty for a programme | Juniors |

A certificate of a programme also fills every placeholder of [Results of a programme](#results) from the counting attempt.

## Elements {#elements}

| Element | Prints |
|---|---|
| *Own text* | its text, with any placeholders |
| *QR code* | its value as a QR code |
| *Barcode* | its value as a Code 128 barcode |
| *Line* | a line |
| image | an image of the club's [image library](../user-guide/print-forms.md#images) |
