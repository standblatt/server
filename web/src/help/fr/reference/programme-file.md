# Fichier de programme

**Exporter** dans le menu du programme ou sur la page du programme télécharge un programme en `<name>.programme.json` ; **Fichier de programme** sous *Ajouter depuis …* l’importe dans une autre manifestation ou une autre société. L’exploitant de la plateforme téléverse les mêmes fichiers dans la bibliothèque globale de programmes.

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

- `ageCategory` est l’un de `u10`, `u13`, `u15`, `u17`, `u19`, `u21`,
  `elite`, `senior`, `veteran`, `senior_veteran`.
- `kind` vaut `card` (carte-couronne, imprimée) ou `natural` (prix en
  nature).
- `showOnRankings` est l’interrupteur *Afficher sur les écrans kiosque*.
- Les combinaisons et les prix spéciaux sont des réglages du programme au
  sein d’une manifestation et ne font pas partie du fichier ; copier des
  programmes d’une manifestation précédente les conserve.
