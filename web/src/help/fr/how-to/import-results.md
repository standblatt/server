# Importer les résultats du système de marquage

Tu en as besoin quand les résultats de nombreux tireurs se trouvent dans un
fichier, p. ex. exporté du système de marquage électronique ou saisi dans un
tableur. Il faut le droit *Configurer les manifestations* (voir
[Rôles](../user-guide/users.md#roles)) ; la manifestation ne doit pas être
clôturée.

## Préparer le fichier {#file}

1. Exporte ou enregistre les résultats en CSV avec une ligne d’en-tête et
   une ligne par tireur et programme. Les colonnes sont listées sous [Import
   CSV](../reference/csv-import.md#results) : au moins `license` (ou
   `firstname` et `lastname`), `programme` et les passes `s1`, `s2`, ….
2. Dans la colonne `programme`, utilise le code-barres du programme ou son
   nom, avec la variante entre parenthèses quand il y a des variantes,
   p. ex. `Championship (Sport)`.
3. Si le fichier a une colonne `club`, vérifie que chaque société figure
   dans la liste des sociétés sous **Données de base** → **Sociétés** ;
   ajoute-y celles qui manquent.

## L’importer {#import}

1. Ouvre la manifestation et clique sur **Saisir les résultats** dans la
   barre latérale.
2. Clique sur **Importer un CSV**, choisis le fichier et clique sur
   **Importer**.
3. *Importé : … résultats, … nouveaux tireurs.* signifie que tu as terminé ;
   les résultats sont sous **Classements**.

N’importe un fichier qu’une fois. Un deuxième import est refusé pour les
programmes tirés une seule fois (les tireurs ont déjà des résultats), mais
un programme à plusieurs [tentatives](../user-guide/programmes.md#attempts)
reprend chaque ligne comme tentative supplémentaire.

## Corriger les erreurs de ligne {#errors}

À la moindre erreur, rien n’est importé et le dialogue affiche *Rien n’a été
importé, veuillez corriger ces lignes :* avec les numéros de ligne. Corrige
chaque ligne listée dans le fichier et importe de nouveau tout le fichier.

| Message | Correction |
|---|---|
| Programme "…" not found | écris le programme comme sur sa page de programme, avec la variante s’il y en a plusieurs, ou utilise son code-barres |
| Club "…" is not in the club list | corrige l’orthographe ou ajoute la société sous **Sociétés** |
| No shooter with the licence …; add first and last name to create one | ajoute `firstname` et `lastname` à la ligne, pour que le tireur puisse être créé |
| Several shooters are called …; add the birthday or the licence | ajoute `birthday` ou `license` à la ligne |
| Series … is missing, … is not a number | remplis chaque passe du programme ; les décimales prennent `.` ou `,` |
| Passe … : les points doivent être entre 0 et … | corrige les points ; le maximum est fixé par les passes du programme |
| Passe … : saisir les … coups | indique chaque coup de la passe, ou aucun et seulement le total de la passe |
| Le tireur a déjà des résultats dans … | supprime la ligne ; les résultats sont déjà dans Standblatt |
| The same shooter and programme appear twice in the file | garde l’une des deux lignes |
| Le tireur n’a plus de tentatives dans … | supprime la ligne, ou augmente les tentatives du programme |

Un message sur la ligne d’en-tête (*The header row needs the column
programme*, *… the column license or the columns firstname and lastname*)
signifie qu’un nom de colonne manque ou est mal orthographié dans la
première ligne.

Les tireurs qui ne sont pas dans les données de base sont créés par
l’import. Des résultats isolés peuvent encore être corrigés après coup dans
les résultats de l’inscription (voir [Résultats](../user-guide/results.md)).
