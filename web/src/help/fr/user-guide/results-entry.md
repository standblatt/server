# Saisir les résultats

Les résultats peuvent être saisis de trois façons, mélangeables au sein
d’une manifestation :

- **À la cible** sur une tablette ou un téléphone, par un lien de secrétaire
  (voir [Aides](helpers.md)) : le secrétaire saisit les passes, ou chaque
  coup, des tireurs de ses cibles.
- **Au desk** par l’inscription et ses résultats (voir
  [Résultats](results.md)), comme lors d’une fête de tir.
- **Après coup** sous **Saisir les résultats** dans la barre latérale de la
  manifestation, tireur par tireur, p. ex. depuis des feuilles papier ou
  l’affichage de l’installation de marquage, ou comme fichier exporté de
  l’installation. Cette page décrit les deux.

## Tireur par tireur {#quick}

1. Choisis le **Programme** ; il reste choisi pour les tireurs suivants.
2. Tape un numéro de licence ou un nom dans **Tireur** et appuie sur
   ++enter++ : le premier résultat est repris. Un lecteur de codes-barres
   fonctionne aussi. **Nouveau tireur** ajoute quelqu’un qui n’est pas
   encore dans les données de base.
3. Tape les passes, ou les coups individuels pour les passes qui les
   saisissent (voir [Programmes](programmes.md)), et appuie sur ++enter++.

Le résultat est enregistré aussitôt, et le tireur apparaît sous *Saisis à
l’instant* avec un lien vers les résultats. Si le tireur était inscrit au
guichet pour ce programme et n’a pas encore de résultats, les résultats vont
dans cette inscription ; sinon, une inscription sans cible est créée, déjà
terminée. Un tireur ne peut avoir des résultats qu’une fois par programme,
sauf si le programme permet plusieurs [tentatives](programmes.md#attempts) :
chaque résultat supplémentaire est alors une nouvelle tentative, jusqu’à ce
qu’il n’en reste plus.

Les distinctions se choisissent comme d’habitude : la distinction atteinte
la plus précieuse est fixée automatiquement et peut être changée dans les
résultats.

## Import depuis un fichier {#import}

**Importer un CSV** (admins) sous **Saisir les résultats** lit un fichier avec une ligne par tireur et programme, p. ex. exporté de l’installation de marquage ou saisi dans un tableur ; les colonnes sont décrites sous [Import CSV](../reference/csv-import.md#results). Chaque ligne est d’abord contrôlée : en cas d’erreur rien n’est importé et le dialogue liste les lignes à corriger. Les tireurs absents des données de base sont créés ; une société doit déjà figurer dans la liste.
