# File di programma

**Esporta** nel menu del programma scarica un programma come `<name>.programme.json`; **File di programma** sotto *Aggiungi da …* lo importa in un’altra manifestazione o in un’altra società. Il gestore della piattaforma carica gli stessi file nella biblioteca globale dei programmi.

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

- `ageCategory` è uno tra `u10`, `u13`, `u15`, `u17`, `u19`, `u21`, `elite`,
  `senior`, `veteran`, `senior_veteran`.
- `kind` è `card` (carta corona, stampata) o `natural` (premio in natura).
- `showOnRankings` è l’interruttore *Mostra sugli schermi chiosco*.
- Combinazioni e premi speciali sono impostazioni del programma all’interno
  di una manifestazione e non fanno parte del file; copiare i programmi da
  una manifestazione precedente li mantiene.
