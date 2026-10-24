# Organiser un programme annuel sur la saison

Tu en as besoin quand les membres tirent un programme aussi souvent qu’ils
le veulent pendant la saison et qu’un championnat de société classe les
trois meilleurs résultats de chaque tireur. Il faut le droit *Configurer les
manifestations* (voir [Rôles](../user-guide/users.md#roles)) ; saisir les
résultats demande *Gérer les manifestations*.

## Créer la manifestation {#event}

1. Sous **Manifestations**, clique sur **Nouvelle manifestation** et choisis
   **Programme annuel ou entraînement**. Cela active les modules
   *Tentatives* et *Classements généraux* ; choisis plutôt **Championnat de
   société** si le classement doit aussi passer sur un affichage kiosque.
2. Saisis le **Nom**, p. ex. *Programme annuel 2027*, et comme **Date** le
   premier jour de la saison. Son année fixe les catégories d’âge.
3. Coche le stand sous **Stands de tir**, ou aucun si les résultats arrivent
   seulement sur papier (voir [Manifestations sans
   cibles](../user-guide/events.md#no-lanes)).
4. Active **Période (saison)**, saisis le **Premier jour** et le **Dernier
   jour**, et sous **Heures de tir** coche les jours de semaine de
   l’entraînement de la société avec leurs horaires. Aucun jour coché
   signifie chaque jour de la période.
5. Clique sur **Créer**.

## Ajouter le programme {#programme}

1. Ouvre **Configuration** → **Programmes** et clique sur **Nouveau
   programme** (ou sur **Ajouter depuis …** pour le reprendre de la
   bibliothèque, d’un fichier ou d’une manifestation précédente).
2. Saisis le nom et les passes.
3. Sous **Tentatives**, coche **Illimité** et règle **Tentative qui compte**
   sur **La meilleure**.
4. Sous **Prix (CHF)**, coche chaque catégorie d’âge qui peut tirer et
   saisis le prix. Le second champ est le prix de chaque tentative
   supplémentaire : vide signifie le même prix, 0 signifie que le programme
   est payé une fois.
5. Clique sur **Créer**.

## Classer les trois meilleures tentatives {#classification}

1. Ouvre **Configuration** → **Classements généraux** et clique sur
   **Nouveau classement général**.
2. Saisis un **Nom**, p. ex. *Championnat de société*, choisis le programme
   annuel comme **Programme 1** et règle **Tentatives qui comptent** sur 3.
3. Active **Classer seulement les tireurs avec tous les résultats qui
   comptent** si seuls les tireurs avec trois résultats doivent recevoir un
   rang ; les autres sont alors listés sous le classement.
4. Ajoute des **Classes** si les juniors et les adultes sont classés
   séparément, p. ex. *Juniors* avec U17 à U21 et *Élite* avec les autres.
5. Clique sur **Créer**.

Tous les réglages de l’éditeur sont décrits dans [Classements
généraux](../user-guide/classifications.md).

## Saisir les résultats pendant la saison {#results}

1. Ouvre la manifestation et clique sur **Saisir les résultats** dans la
   barre latérale.
2. Choisis le **Programme**, tape un numéro de licence ou un nom dans
   **Tireur** et appuie sur ++enter++.
3. Tape les passes et appuie sur ++enter++. Chaque résultat d’un tireur est
   une nouvelle tentative.

Les résultats du système de marquage peuvent aussi être importés sous forme
de fichier (voir [Importer les résultats du système de
marquage](import-results.md)).

**Classements** montre le programme avec la meilleure tentative de chaque tireur (*2/5* : la deuxième de cinq tentatives compte) et le championnat de société avec la somme des trois meilleures.

## Terminer la saison {#close}

Quand les derniers résultats sont saisis et vérifiés, ouvre
**Configuration** → **Copier & clôturer** et clique sur **Clôturer la
manifestation**. Ce que la clôture change est expliqué dans [Clôturer une
manifestation](../explanation/closing.md) ; la saison suivante commence avec
[Reprendre la manifestation de l’année passée](next-year.md).
