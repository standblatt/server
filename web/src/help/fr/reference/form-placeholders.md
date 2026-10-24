# Champs des modèles

Un [modèle d’impression](../user-guide/print-forms.md) imprime à la place
d’un espace réservé comme `{{shooter.name}}` sa valeur. Un espace réservé
sans valeur n’imprime rien ; de même pour un espace réservé que la [source
de données](../user-guide/print-forms.md#source) du modèle ne connaît
pas. Dans un élément répété, `#` compte 1, 2, 3, … : `{{s1.shot#}}` imprime
`{{s1.shot1}}`, `{{s1.shot2}}`, ….

## Toutes les sources de données {#common}

| Espace réservé | Champ | Exemple |
|---|---|---|
| `{{shooter.name}}` | Nom (nom de famille d’abord) | Example Anna |
| `{{shooter.lastname}}` | Nom | Example |
| `{{shooter.firstname}}` | Prénom | Anna |
| `{{shooter.license}}` | Numéro de licence | 123456 |
| `{{shooter.birthday}}` | Date de naissance | 02.06.1985 |
| `{{shooter.birthyear}}` | Année de naissance | 1985 |
| `{{shooter.postalcode}}` | NPA | 9999 |
| `{{shooter.city}}` | Localité | Hilltown |
| `{{category}}` | Catégorie d’âge | Élite |
| `{{club}}` | Société | SG Hilltown |
| `{{club.ssv}}` | Numéro FST | 1.23.4.56 |
| `{{event.name}}` | Manifestation | Programme annuel 2027 |
| `{{event.date}}` | Date (jour principal) | 12.05.2027 |
| `{{event.organizer}}` | Organisateur | SG Hilltown |
| `{{today}}` | Aujourd’hui, dans le fuseau horaire de la manifestation | 12.05.2027 |

## Résultats d’un programme {#results}

Modèles de programme et cartes-couronnes.

| Espace réservé | Champ | Exemple |
|---|---|---|
| `{{programme}}` | Programme, avec sa variante | Obligatoire |
| `{{lane}}` | Cible | 7 |
| `{{registration}}` | Numéro d’inscription | 1042 |
| `{{total}}` | Total | 85 |
| `{{mouches}}` | Mouches | 3 |
| `{{s1}}`, `{{s2}}`, … | Passe 1, 2, … | 47 |
| `{{s1.mouches}}`, … | Mouches passe 1, … | 1 |
| `{{s1.shot1}}`, `{{s1.shot#}}`, … | coups individuels de la passe 1 (passes avec **Coups individuels**) | 5 |
| `{{award}}` | Distinction | Carte |
| `{{award.value}}` | Valeur de la distinction | CHF 12.00 |
| `{{award.value.words}}` | Valeur en lettres, dans la langue de la manifestation | douze francs |

`{{award.value}}` et `{{award.value.words}}` restent vides sauf si *Imprimer
la valeur de la distinction sur les cartes-couronnes* est activé et que la
distinction a une valeur.

## Une inscription {#registration}

Bandes de contrôle et documents d’inscription.

| Espace réservé | Champ | Exemple |
|---|---|---|
| `{{category.short}}` | Catégorie d’âge, abrégée | E |
| `{{lane}}` | Cible | 7 |
| `{{registration}}` | Numéro d’inscription | 1042 |
| `{{registration.code}}` | Code QR de l’inscription (le [code d’inscription](codes.md)) | REGCODEAAAA1042AAAAL123456 |
| `{{programmes}}` | Programmes, séparés par des virgules | Obligatoire, Tir en campagne |
| `{{passes}}` | Passes : `P1 P2` par programme, `P0` pour un programme qui reprend ses passes d’un autre, `*` sur les passes reprises | P1 P2 + P0 |

## Un rang {#rank}

Diplômes.

| Espace réservé | Champ | Exemple |
|---|---|---|
| `{{rank}}` | Rang | 1 |
| `{{total}}` | Total | 98 |
| `{{ranking}}` | Classement : le programme ou le classement général | Championnat de section |
| `{{class}}` | Classe du classement général ; vide pour un programme | Juniors |

Le diplôme d’un programme remplit aussi chaque espace réservé de [Résultats
d’un programme](#results) à partir de la tentative qui compte.

## Éléments {#elements}

| Élément | Imprime |
|---|---|
| *Texte libre* | son texte, avec d’éventuels espaces réservés |
| *Code QR* | sa valeur sous forme de code QR |
| *Code-barres* | sa valeur sous forme de code-barres Code 128 |
| *Ligne* | une ligne |
| image | une image de la [bibliothèque d’images](../user-guide/print-forms.md#images) de la société |
