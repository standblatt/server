# Import CSV

Les membres avec le droit *Gérer les données de base* importent les tireurs
sur la page **Tireurs**. Le fichier a besoin d’une ligne d’en-tête ; le
séparateur (`;` ou `,`) en est déduit.

```text
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9999;Hilltown;PS Example
```

| Colonne | Obligatoire | Remarques |
|---|---|---|
| `firstname` | oui | |
| `lastname` | oui | |
| `license` | non | aussi `licensenumber` ; les numéros purement numériques sont complétés à six chiffres |
| `birthday` | non | `YYYY-MM-DD` ou `DD.MM.YYYY` |
| `postalcode`, `city` | non | |
| `club` | non | la société du tireur, comparée sans tenir compte des majuscules ni des espaces |

- Une ligne par tireur et société ; un tireur dans deux sociétés apparaît
  deux fois.
- Les tireurs sont reconnus par leur numéro de licence, ou par nom et date
  de naissance s’il n’y en a pas. Les tireurs existants sont mis à jour,
  rien n’est supprimé.
- Un fichier avec des lignes invalides est refusé en entier, avec la liste
  des lignes à corriger.
- Les sociétés absentes de la liste arrêtent l’import avant que quoi que ce
  soit ne soit enregistré : choisis pour chacune une société existante (les
  noms semblables sont proposés) ou une nouvelle, puis clique à nouveau sur
  **Importer**.

## Résultats {#results}

Les membres ayant le droit *Configurer les manifestations* importent des
résultats avec **Importer un CSV** sur la page **Saisir les résultats**
d’une manifestation (voir [Saisir les
résultats](../user-guide/results-entry.md#import)) :

```text
license;firstname;lastname;birthday;club;programme;s1;s2;m2;s1_1;s1_2;s1_3
013756;;;;;Championship;96;93;2;;;
;Max;Muster;1990-02-02;SG Hilltown;Obligatory;;45;;9;10;8
```

| Colonne | Obligatoire | Remarques |
|---|---|---|
| `license` ou `firstname` + `lastname` | oui | le tireur : trouvé par sa licence, sinon par son nom (et `birthday` si plusieurs portent ce nom) ; créé s’il n’est pas trouvé, ce qui demande `firstname` et `lastname` : une licence seule qui n’est pas trouvée est une erreur de ligne |
| `birthday` | non | `YYYY-MM-DD` ou `DD.MM.YYYY` |
| `club` | non | doit figurer dans la liste des sociétés (sans tenir compte des majuscules ni des espaces) ; le tireur la rejoint |
| `programme` | oui | code-barres, nom avec variante comme `Championship (Sport)` ou le nom seul |
| `s1`, `s2`, … | oui | points par passe, décimales avec `.` ou `,` |
| `m1`, `m2`, … | non | mouches par passe |
| `s1_1`, `s1_2`, … | non | coups individuels de la passe 1, pour les passes qui les saisissent ; la passe en est la somme |

- Chaque ligne est d’abord vérifiée ; à la moindre erreur de ligne, rien
  n’est importé.
- Un tireur ne peut avoir des résultats qu’une fois par programme ; une
  ligne pour des résultats déjà présents est une erreur. Les programmes à
  plusieurs [tentatives](../user-guide/programmes.md#attempts) prennent une
  ligne par tentative, jusqu’à ce que le tireur n’en ait plus.
- Pour les programmes à plusieurs tentatives, importer à nouveau le même
  fichier ajoute ses lignes comme tentatives supplémentaires.
- Si le tireur est inscrit au programme sans résultats, les résultats vont
  dans cette inscription.
