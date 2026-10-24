# Imprimer les résultats sur les feuilles de stand de l’armée

Tu en as besoin quand les résultats des tireurs doivent figurer sur des
feuilles préimprimées, p. ex. les feuilles de stand que l’armée fournit pour
le programme obligatoire, et que tu préfères les imprimer plutôt que de les
remplir à la main. Il faut les droits *Modèles d’impression*, *Configurer
les manifestations* et *Stands et imprimantes* (voir
[Rôles](../user-guide/users.md#roles)).

Avant de commencer, scanne une feuille vide en PNG ou JPEG et garde quelques
feuilles de réserve pour les impressions d’essai.

## 1. Préparer le programme {#programme}

1. Vérifie que **Modules de la société**, dans la partie société de la barre
   latérale, a *Modèles d’impression* activé (c’est le cas, sauf si
   quelqu’un l’a désactivé).
2. Sous **Configuration** → **Programmes** de la manifestation, ouvre le
   programme et coche **Coups individuels** pour chaque passe dont la
   feuille a des cases pour les coups. Fais-le avant la première inscription
   : dès que des tireurs sont inscrits, les passes sont fixées.
3. **Enregistrer** le programme.

## 2. Créer le modèle {#form}

1. Ouvre **Modèles d’impression** dans la partie société de la barre
   latérale et clique sur **Nouveau modèle d’impression**.
2. Saisis un nom (p. ex. *Programme obligatoire 300 m*), choisis **Imprime**
   *Résultats d'un programme* et le format de papier le plus proche de la
   feuille, puis **Créer**. L’éditeur s’ouvre.
3. Dans **Modèle** à droite, règle **Largeur (mm)** et **Hauteur (mm)** sur
   la taille de la feuille, puis **Téléverser un scan** sous **Scan du
   formulaire vide**. Le scan apparaît derrière la feuille pendant l’édition
   ; il n’est jamais imprimé.
4. Depuis les **Champs** à gauche, glisse les champs du tireur
   (p. ex. **Nom**, **Prénom**, **Année de naissance**, **Numéro de
   licence**) sur leurs cases du scan.
5. Depuis le groupe **Coups individuels**, glisse **Coups passe 1** sur la
   première case de coup de cette passe. C’est une rangée de cases à 8 mm
   d’intervalle ; tire le point rose sur sa dernière case jusqu’à ce que
   chaque valeur soit dans sa case, ou règle **dx (mm)** sous
   **Propriétés**. Répète pour les autres passes.
6. Depuis **Résultat**, ajoute les totaux des passes (**Passe 1** …) et le
   **Total** là où la feuille a des cases pour eux.
7. **Enregistrer**.

![Modèles d’impression](../images/screenshots/print-forms.png){ .screenshot
}

Comment déplacer, aligner et aimanter les éléments est expliqué dans
[Modèles d’impression](../user-guide/print-forms.md#arrange), chaque champ
dans [Champs des modèles](../reference/form-placeholders.md).

## 3. Calibrer l’imprimante {#calibrate}

1. Mets une feuille vide dans l’imprimante au stand.
2. Clique sur **Impression d’essai** : le modèle enregistré s’ouvre en PDF
   avec des valeurs inventées et sans le scan. Imprime-le en taille réelle
   (sans ajustement à la page).
3. Mesure de combien les valeurs sont décalées par rapport à leurs cases et
   saisis-le comme **Décalage de l’imprimante vers la droite (mm)** et
   **Décalage de l’imprimante vers le bas (mm)** sous **Modèle** (en négatif
   pour déplacer vers la gauche ou vers le haut).
4. **Enregistrer** et recommence à l’étape 1 jusqu’à ce que les valeurs
   soient dans les cases.

## 4. Attribuer le modèle et l’imprimante {#assign}

1. Sous **Configuration** → **Programmes** de la manifestation, ouvre le
   programme et choisis le modèle dans la section **Modèle
   d’impression**. Il est enregistré aussitôt ; cela fonctionne même quand
   des tireurs sont inscrits.
2. Sous **Stands de tir → Configurer les imprimantes** du stand, choisis
   l’imprimante pour **Modèles d’impression (préimprimés)**, p. ex. avec son
   alimentation manuelle comme bac, et **Enregistrer** (voir [Connecter une
   imprimante](connect-printer.md)).

## 5. Imprimer les feuilles {#print}

Dès que toutes les passes du programme sont saisies, ses résultats affichent
**Imprimer sur le modèle** : dans les résultats à l’accueil et sur les liens
de cibarre à la cible (voir [Aides](../user-guide/helpers.md)). Sur **Saisir
les résultats**, chaque ligne sous *Saisis à l’instant* a son propre bouton
d’impression. Mets une feuille dans l’imprimante et clique dessus. Sans
agent d’impression en ligne, le PDF s’ouvre et s’imprime depuis le
navigateur.
