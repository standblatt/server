# Codes et codes-barres

## Code d’inscription (QR)

Les feuilles de stand et les bandes de contrôle portent un code QR avec

```text
REGCODEAAAA<inscription>AAAAL<licence ou 0>
```

par exemple `REGCODEAAAA42AAAAL100109`. Un scanner manuel le tape dans le
champ de recherche du tableau de bord et appuie sur Entrée ; Standblatt
ouvre les résultats de cette inscription. La partie licence protège contre
le code d’un autre tireur ; un code d’une autre manifestation est refusé.

## Code-barres du programme (Interleaved 2 of 5)

Un programme peut avoir un code-barres, imprimé en Interleaved 2 of 5 à côté
du programme sur la feuille de stand (seuls les chiffres comptent, complétés
à un nombre pair). Les systèmes de marquage électronique s’en servent pour
choisir le programme.
