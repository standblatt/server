# Classements généraux

Un classement de programme classe un programme. Un **classement général**
classe les tireurs sur plusieurs programmes ou sur plusieurs
[tentatives](programmes.md#attempts) d’un programme, p. ex. un championnat
de section comme somme de trois programmes, ou un championnat de société sur
les trois meilleures tentatives d’un programme annuel. Les classements
généraux apparaissent sous **Classements**, sur les écrans kiosque et comme
rapport ; ils ne modifient jamais les résultats ni les distinctions.

## Créer un classement général {#create}

Dans l’onglet *Classements généraux* de la **Configuration** de la
manifestation (admins), **Nouveau classement général** ouvre l’éditeur :

- **Programmes** : chaque programme qui compte, avec ses **Tentatives qui
  comptent**. Avec 3, les trois meilleures tentatives comptent (ou les trois
  premières ou dernières, selon la *Tentative qui compte* du programme). Le
  total d’un tireur est la somme de toutes les tentatives comptées.
- **Ne compter que les meilleurs programmes** : avec plusieurs programmes,
  seuls autant de programmes avec le plus de points comptent, p. ex. les
  deux meilleurs sur trois ; 0 les compte tous.
- **Classer seulement les tireurs avec tous les résultats qui comptent** :
  les tireurs à qui il manque encore un programme ou une tentative sont
  listés sous le classement, sans rang. Sans cette option, chacun est classé
  selon ce qu’il a.
- **Départage** : en cas d’égalité, les critères cochés décident dans leur
  ordre ; les flèches le modifient. *Passes de la dernière en arrière*
  compare toutes les passes comptées, *Plus de coups les plus hauts*
  nécessite des [coups individuels](programmes.md) (plus de dix, puis plus
  de neuf …), *Âge* suit la [règle FST](../reference/rules.md#ranking). Par
  défaut, l’ordre FST.
- **Classes** : classements séparés, chacun pour certaines catégories d’âge,
  p. ex. *Juniors* avec U17 à U21 et *Élite* avec les autres. Une catégorie
  d’âge appartient au plus à une classe ; les tireurs dont la catégorie
  n’est dans aucune classe ne sont pas listés. Sans classes, tout le monde
  est dans un seul classement.
- **Afficher sur les écrans kiosque**.

## Exemples {#examples}

| Manifestation | Programmes | Réglages |
|---|---|---|
| Championnat de section | Programme obligatoire, Tir en campagne, Cible section | tentatives qui comptent 1 chacun, *Classer seulement les tireurs avec tous les résultats qui comptent* |
| Championnat de société sur la saison | programme annuel à tentatives illimitées | tentatives qui comptent 3 |
| Fête de tir, deux meilleures cibles | trois cibles | *Ne compter que les meilleurs programmes* 2 |

**Copier pour l’année prochaine** copie les classements généraux avec les programmes. Supprimer un programme le retire des classements généraux.
