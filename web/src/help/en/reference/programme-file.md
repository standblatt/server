# Programme file

**Export** in the programme menu or on the programme page downloads a programme as `<name>.programme.json`; **Programme file** under *Add from …* imports it into another event or another club. The platform operator uploads the same files to the global programme library.

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

- `ageCategory` is one of `u10`, `u13`, `u15`, `u17`, `u19`, `u21`, `elite`, `senior`, `veteran`, `senior_veteran`.
- `kind` is `card` (wreath card, printed), `natural` (prize in kind) or `cash` (cash prize).
- `showOnRankings` is the *Show on kiosk displays* switch.
- Optional fields: `attempts` (0 or left out = once), `unlimitedAttempts`, `attemptRule` (`best`, `first` or `last`) and `awardCountingOnly` for [attempts](../user-guide/programmes.md#attempts); `repeatPrice` in a price (left out = the same price, 0 = paid once); `trackShots` in a series for single shots.
- Combinations and special prizes are settings of the programme within an event and are not part of the file; copying programmes from an earlier event keeps them.
