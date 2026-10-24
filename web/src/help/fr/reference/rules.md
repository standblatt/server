# Règles de concours

Cette page décrit les règles exactement comme Standblatt les implémente.

## Catégorie d’âge

La catégorie découle de *année de la manifestation − année de naissance* ;
voir le tableau dans
[Notions](../getting-started/concepts.md#age-categories). Sans date de
naissance, un tireur compte comme élite. La catégorie est fixée à
l’inscription.

## Inscription

Une inscription est refusée quand :

- un programme n’appartient pas à la manifestation, est choisi deux fois ou
  n’est pas actif, ou plus de tentatives sont choisies qu’il n’en permet ;
- un programme n’a pas de prix pour la catégorie d’âge du tireur (« pas pour
  cette catégorie d’âge ») ;
- deux variantes du même programme sont choisies ;
- un programme reprend des passes d’un autre programme qui n’est pas choisi,
  ou le nombre de passes choisies ne correspond pas ;
- la cible est occupée par une autre inscription active ou hors service ;
- le tireur est déjà inscrit à l’un des programmes de cette manifestation
  (un tireur tire chaque programme une fois par manifestation, sauf s’il
  permet plusieurs [tentatives](../user-guide/programmes.md#attempts) ; une
  autre variante dans une inscription séparée est permise).

Une cible réservée peut quand même être utilisée après confirmation.

## Résultats

Chaque passe doit être saisie. Les points sont entre 0 et le maximum de la
passe (par défaut coups × 10), les mouches entre 0 et le nombre de coups, et
seulement sur les passes qui comptent les mouches. L’enregistrement des
résultats attribue la distinction la plus précieuse atteinte ; un nouvel
enregistrement garde la distinction choisie tant que le nouveau total
l’atteint encore, sinon la plus précieuse atteinte. Dès que chaque programme
d’une inscription a des résultats, l’inscription est terminée et sa cible se
libère.

### Combinaisons

Un programme marqué *combinable avec* un autre programme peut en reprendre
des passes : à l’inscription, le tireur choisit lesquelles (par exemple les
passes 2 et 3 d’un programme à quatre passes pour un programme à deux
passes). Les passes choisies sont copiées dans l’ordre croissant et
renumérotées à partir de 1 ; le programme combiné ne peut pas être modifié
seul et suit chaque modification de la source.

## Distinctions

Un niveau de distinction n’est ouvert à une catégorie d’âge que s’il a une
limite pour elle. Un tireur peut choisir tout niveau dont son total atteint
la limite ; le niveau le plus précieux est en tête et attribué par défaut à
l’enregistrement des résultats. Le serveur vérifie le choix.

## Classement {#ranking}

Tous les résultats d’un programme sont ordonnés par :

1. total le plus élevé ;
2. les passes comparées de la dernière vers l’arrière (la meilleure dernière
   passe l’emporte, puis l’avant-dernière, …) ;
3. plus de mouches ;
4. âge : chez les juniors (U10–U21), le plus jeune d’abord, chez les autres
   le plus âgé d’abord ; les non-juniors passent avant les juniors.

Les rangs sont consécutifs (1, 2, 3, … même à résultat égal). Le même
classement alimente les pages du personnel, les écrans kiosque, le PDF et le
CSV du classement ainsi que les prix spéciaux.

## Concours de sections {#club-competition}

Seuls les programmes marqués *compte pour le concours de sections*
participent, et seulement les résultats terminés de tireurs ayant une
société.

- La valeur d’un résultat est la somme de ses deux premières passes.
- Une société avec moins de 5 participants est listée comme *non classée*.
- Sur *n* participants, les *max(5, ⌊n/2⌋)* meilleurs sont obligatoires, les
  autres facultatifs.
- Résultat = ⌊(Σ obligatoires + 0.02 × Σ facultatifs) / nombre
  d’obligatoires × 1000⌋ / 1000, c.-à-d. tronqué à trois décimales.
- Égalité : plus de participants d’abord, puis les valeurs triées comparées
  une à une.
- La société exclue dans les réglages du module *Concours de sections* (en
  général l’organisatrice) garde sa place dans la liste mais ne reçoit pas
  de rang (« hors concours »).

## Prix spéciaux {#special-prizes}

Un programme peut décerner des prix spéciaux aux meilleurs rangs de chaque
catégorie d’âge. Trois nombres les définissent :

- **Rangs par catégorie**, p. ex. 3 ;
- **participants minimum par catégorie** dans ce programme, p. ex. 10 ;
- **participants minimum de la société du tireur** dans le même programme,
  p. ex. 4.

Par catégorie d’âge comptant assez de participants, les premiers tireurs
dans l’ordre du classement dont la société a assez de participants gagnent ;
les tireurs dont la société en a trop peu sont sautés et le suivant
monte. Les tireurs sans société ne sont pas éligibles, sauf si le minimum de
société est 0. Le rapport liste les gagnants et les tireurs sautés avec la
raison.

!!! example "Du programme aux réglages"
    Un programme dit : *« Dans chaque catégorie, les rangs 1 à 3 reçoivent un prix spécial, pour autant que la catégorie compte au moins 10 participants et la société du tireur au moins 4. »* Cela donne rangs 3, minimum de catégorie 10, minimum de société 4.

## Listes de cibles

Les cibles s’écrivent comme une liste de nombres et de plages : `1, 5,
7-10`. Standblatt les affiche de façon compacte et transforme les suites de
trois ou plus en plages (`1-3, 5, 7, 8`).

## Montants et nombres

Les montants s’affichent `CHF 1'234.50`, les dates `dd.mm.yyyy`. Sur les
cartes-couronnes, la valeur peut être imprimée en toutes lettres dans la
langue de la manifestation, p. ex. *Fünfundzwanzig Franken*, *Septante et un
francs*, *Ventitré franchi*.
