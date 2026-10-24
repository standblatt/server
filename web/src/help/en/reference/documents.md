# Documents and reports

Every document is a PDF in the language of the event. A print goes to the printer set for the document on **Ranges → Set up the printers** of the registration's range when a print agent of the range is online; otherwise the PDF opens (see [Printing](../user-guide/printing.md)).

## Documents {#documents}

| Document | Format | Content | Printed from | Printer setting | Module |
|---|---|---|---|---|---|
| Lane sheet | A5 | event, shooter, QR [registration code](codes.md), range (events on several ranges), lane, programmes with shots and [barcodes](codes.md) | *Print lane sheet* at registration; the registration | *Lane sheet (A5)* | Lane sheet |
| Control strip | label 90 × 29 mm, or a [print form](../user-guide/print-forms.md#builtin) of the club | event, shooter, licence, club, age category, passes, QR registration code | *Print control strip* at registration; the registration | *Control strip (label 90 × 29 mm)* | Control strip |
| Result sheet | A5 | event, shooter, every programme with its series, total and award | **Print result sheet** in the results, once every programme has results | *Result sheet (A5)* | – |
| Wreath card | A6, printed onto pre-printed cards, or a print form of the club | shooter, year of birth, club, date, organizer, value in figures and words (setting *Print the award value on wreath cards*) | **Print wreath card** in the results, for a wreath card award | *Wreath card (A6, pre-printed)* | Awards |
| Print form of a programme | size of the form | see [Form placeholders](form-placeholders.md#results) | **Print onto form** in the results of a completed programme | *Print forms (pre-printed)* | club module Print forms |
| Registration document | size of the form | see [Form placeholders](form-placeholders.md#registration) | a button named after the form in the registration | *Print forms (pre-printed)* | club module Print forms |
| Certificates | size of the form, one page per ranked shooter | see [Form placeholders](form-placeholders.md#rank) | **Certificates** on **Rankings** | *Certificates* | club module Print forms |
| Ranking | A4 | rank, name, year, category, club, series, total; ranks 1–3 highlighted; *(2/3)* for the counting attempt | **PDF** on **Rankings** and **Statistics** | never printed, PDF only | – |
| Test page | label 90 × 29 mm | club, range, printer, options, time | **Test page** on **Ranges → Set up the printers** | the printer of that row | – |

- Without a range, documents open as PDF only.
- Helper links print only the lane sheet (ammunition desk) or print forms (scorer, results clerk), and only of their lanes.
- A print job that the agent does not confirm within two minutes is retried, at most three times.

## Reports {#reports}

Under **Statistics** (permission *Statistics and reports*). Every report is available as **PDF** and **CSV** (UTF-8 with BOM, `;` separated, opens directly in Excel). Texts follow the language of the event.

| Report | Content | Module |
|---|---|---|
| Financial report | revenue per programme, awards with their value | – |
| Revenue by programme and age category | participants, price and subtotal per programme and category | – |
| Club participants | participants per club, programme and age category | – |
| Awards | every award given, with shooter, club, programme, value and result | Awards |
| Awards to hand out | awards not yet collected, with shooter, club, programme, kind and value, plus the cash still to pay out (see [Results](../user-guide/results.md#handout)) | Awards |
| Club competition | ranking, out of competition and not ranked clubs | Club competition |
| Special prizes | winners per programme and age category, skipped shooters with the reason (see [Programmes](../user-guide/programmes.md#special-prizes)) | Special prizes |
| Classifications | every [classification](../user-guide/classifications.md) per class, with the points per programme; parts that do not count are in parentheses | Classifications |
| Lane bookings | the clubs that asked for lanes through [online booking](../user-guide/online-booking.md), confirmed or still waiting, with contact, slots, lane numbers and announced shooters | Reservations |

The PDF of a report is A4, *Awards* in landscape. The ranking of every programme is also a CSV file: rank, last name, first name, year, category, club, series, total and mouches, for programmes with several attempts also the counting attempt and the number of attempts.
