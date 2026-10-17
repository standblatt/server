# Codes und Barcodes

## Anmeldecode (QR)

Standblätter und Kontrollstreifen tragen einen QR-Code mit

```text
REGCODEAAAA<Anmeldung>AAAAL<Lizenz oder 0>
```

zum Beispiel `REGCODEAAAA42AAAAL100109`. Ein Handscanner tippt ihn ins
Suchfeld des Dashboards und drückt Enter; Standblatt öffnet die Resultate
dieser Anmeldung. Der Lizenzteil schützt vor dem Code eines anderen
Schützen; ein Code eines anderen Anlasses wird abgelehnt.

## Programm-Barcode (Interleaved 2 of 5)

Ein Programm kann einen Barcode haben, der als Interleaved 2 of 5 neben dem
Programm auf dem Standblatt gedruckt wird (nur die Ziffern zählen, auf eine
gerade Anzahl Ziffern aufgefüllt). Elektronische Trefferanzeigen wählen
damit das Programm.
