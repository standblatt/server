# Dokumente und Berichte

Jedes Dokument ist ein PDF in der Sprache des Anlasses. Ein Druck geht an
den Drucker, der unter **Schiessstände → Drucker einrichten** des
Schiessstands der Anmeldung für das Dokument eingestellt ist, wenn ein
Druck-Agent des Schiessstands online ist; sonst öffnet sich das PDF (siehe
[Drucken](../user-guide/printing.md)).

## Dokumente {#documents}

| Dokument | Format | Inhalt | Gedruckt aus | Druckereinstellung | Modul |
|---|---|---|---|---|---|
| Standblatt | A5 | Anlass, Schütze, QR-[Anmeldecode](codes.md), Schiessstand (bei Anlässen auf mehreren Schiessständen), Scheibe, Programme mit Schüssen und [Barcodes](codes.md) | *Standblatt drucken* bei der Anmeldung; die Anmeldung | *Standblatt (A5)* | Standblatt |
| Kontrollstreifen | Etikette 90 × 29 mm oder eine [Druckvorlage](../user-guide/print-forms.md#builtin) des Vereins | Anlass, Schütze, Lizenz, Verein, Alterskategorie, Passen, QR-Anmeldecode | *Kontrollstreifen drucken* bei der Anmeldung; die Anmeldung | *Kontrollstreifen (Etikette 90 × 29 mm)* | Kontrollstreifen |
| Resultatblatt | A5 | Anlass, Schütze, jedes Programm mit seinen Passen, Total und Auszeichnung | **Resultatblatt drucken** in den Resultaten, sobald jedes Programm Resultate hat | *Resultatblatt (A5)* | – |
| Kranzkarte | A6, auf vorgedruckte Karten gedruckt, oder eine Druckvorlage des Vereins | Schütze, Jahrgang, Verein, Datum, Organisator, Wert in Ziffern und Worten (Einstellung *Wert der Auszeichnung auf Kranzkarten drucken*) | **Kranzkarte drucken** in den Resultaten, bei einer Kranzkarte als Auszeichnung | *Kranzkarte (A6, vorgedruckt)* | Auszeichnungen |
| Druckvorlage eines Programms | Grösse der Vorlage | siehe [Platzhalter in Vorlagen](form-placeholders.md#results) | **Auf Vorlage drucken** in den Resultaten eines abgeschlossenen Programms | *Druckvorlagen (vorgedruckt)* | Vereinsmodul Druckvorlagen |
| Anmeldedokument | Grösse der Vorlage | siehe [Platzhalter in Vorlagen](form-placeholders.md#registration) | ein nach der Vorlage benannter Knopf in der Anmeldung | *Druckvorlagen (vorgedruckt)* | Vereinsmodul Druckvorlagen |
| Urkunden | Grösse der Vorlage, eine Seite pro rangiertem Schützen | siehe [Platzhalter in Vorlagen](form-placeholders.md#rank) | **Urkunden** unter **Ranglisten** | *Urkunden* | Vereinsmodul Druckvorlagen |
| Rangliste | A4 | Rang, Name, Jahrgang, Kategorie, Verein, Passen, Total; Ränge 1–3 hervorgehoben; *(2/3)* für den zählenden Durchgang | **PDF** unter **Ranglisten** und **Statistik** | nie gedruckt, nur PDF | – |
| Testseite | Etikette 90 × 29 mm | Verein, Schiessstand, Drucker, Optionen, Uhrzeit | **Testseite** unter **Schiessstände → Drucker einrichten** | der Drucker dieser Zeile | – |

- Ohne Schiessstand öffnen sich Dokumente nur als PDF.
- Helferlinks drucken nur das Standblatt (Munitionsausgabe) oder
  Druckvorlagen (Warner, Resultaterfassung), und nur für ihre Scheiben.
- Ein Druckauftrag, den der Agent nicht innert zwei Minuten bestätigt, wird
  erneut versucht, höchstens dreimal.

## Berichte {#reports}

Unter **Statistik** (Recht *Statistik und Berichte*). Jeder Bericht steht
als **PDF** und **CSV** zur Verfügung (UTF-8 mit BOM, durch `;` getrennt,
öffnet sich direkt in Excel). Die Texte folgen der Sprache des Anlasses.

| Bericht | Inhalt | Modul |
|---|---|---|
| Finanzbericht | Einnahmen pro Programm, Auszeichnungen mit ihrem Wert | – |
| Einnahmen nach Programm und Alterskategorie | Teilnehmer, Preis und Zwischensumme pro Programm und Kategorie | – |
| Vereinsteilnehmer | Teilnehmer pro Verein, Programm und Alterskategorie | – |
| Auszeichnungen | jede vergebene Auszeichnung mit Schütze, Verein, Programm, Wert und Resultat | Auszeichnungen |
| Abzugebende Auszeichnungen | noch nicht abgeholte Auszeichnungen mit Schütze, Verein, Programm, Art und Wert, dazu das noch auszuzahlende Bargeld (siehe [Resultate](../user-guide/results.md#handout)) | Auszeichnungen |
| Vereinskonkurrenz | Rangliste, Vereine ausser Konkurrenz und nicht rangierte Vereine | Vereinskonkurrenz |
| Sonderpreise | Gewinner pro Programm und Alterskategorie, übersprungene Schützen mit Grund (siehe [Programme](../user-guide/programmes.md#special-prizes)) | Sonderpreise |
| Wertungen | jede [Wertung](../user-guide/classifications.md) pro Klasse, mit den Punkten pro Programm; nicht zählende Teile stehen in Klammern | Wertungen |
| Scheibenreservationen | die Vereine, die über die [Online-Reservation](../user-guide/online-booking.md) Scheiben angefragt haben, bestätigt oder noch offen, mit Kontakt, Zeitblöcken, Scheibennummern und angemeldeten Schützen | Reservationen |

Das PDF eines Berichts ist A4, *Auszeichnungen* im Querformat. Die Rangliste
jedes Programms gibt es auch als CSV-Datei: Rang, Nachname, Vorname,
Jahrgang, Kategorie, Verein, Passen, Total und Mouchen, bei Programmen mit
mehreren Durchgängen auch der zählende Durchgang und die Anzahl Durchgänge.
