# Programmdatei

**Exportieren** im Programmmenü oder auf der Programmseite lädt ein Programm als `<name>.programme.json` herunter; **Programmdatei** unter *Übernehmen aus …* importiert es in einen anderen Anlass oder einen anderen Verein. Der Plattformbetreiber lädt dieselben Dateien in die globale Programmbibliothek hoch.

```json
{
  "format": "standblatt.programme",
  "version": 1,
  "programme": {
    "name": "Championship",
    "variant": "Sport",
    "description": "",
    "barcode": "700005002334864165956154",
    "countsForClubCompetition": true,
    "showOnRankings": true,
    "series": [
      { "shots": 10, "maxPoints": 100, "trackMouches": false, "name": "" }
    ],
    "prices": [
      { "ageCategory": "elite", "price": 20 }
    ],
    "awards": [
      {
        "name": "Wreath card CHF 10",
        "kind": "card",
        "value": 10,
        "limits": [{ "ageCategory": "elite", "points": 370 }]
      }
    ]
  }
}
```

- `ageCategory` ist eines von `u10`, `u13`, `u15`, `u17`, `u19`, `u21`,
  `elite`, `senior`, `veteran`, `senior_veteran`.
- `kind` ist `card` (Kranzkarte, gedruckt) oder `natural` (Naturalgabe).
- `showOnRankings` ist der Schalter *Auf Kiosk-Anzeigen zeigen*.
- Kombinationen und Sonderpreise sind Einstellungen des Programms innerhalb
  eines Anlasses und gehören nicht zur Datei; das Kopieren von Programmen
  aus einem früheren Anlass behält sie.
