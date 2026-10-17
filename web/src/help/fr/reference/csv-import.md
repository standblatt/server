# Import CSV

Les admins de société importent les tireurs sur la page **Tireurs**. Le
fichier a besoin d’une ligne d’en-tête ; le séparateur (`;` ou `,`) en est
déduit.

```text
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9050;Appenzell;PS Example
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
