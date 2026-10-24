# Imprimer des diplômes pour les meilleurs rangs

Tu en as besoin quand les gagnants d’un programme ou d’un classement général
reçoivent un diplôme avec leur nom et leur rang, p. ex. les rangs 1 à 3 à la
distribution des prix. Créer le modèle demande le droit *Modèles
d’impression*, imprimer *Gérer les manifestations* (voir
[Rôles](../user-guide/users.md#roles)).

## 1. Créer le modèle de diplôme {#form}

1. Ouvre **Modèles d’impression** dans la partie société de la barre
   latérale et clique sur **Nouveau modèle d’impression**.
2. Saisis un nom (p. ex. *Diplôme*), choisis **Imprime** *Un rang* et un
   format de papier comme *A4 portrait (210 × 297 mm)*, puis **Créer**.
3. Depuis les **Champs**, ajoute ce que le diplôme affiche : **Nom** depuis
   *Tireur*, **Rang**, **Classement** et **Classe** depuis *Rang*, **Total**
   depuis *Résultat*, et **Texte libre** pour des mots fixes comme
   *Diplôme*.
4. Pour un logo, utilise d’abord **Téléverser une image** dans la
   bibliothèque d’images sous la liste des **Modèles d’impression** ; elle
   apparaît ensuite dans le groupe **Images** des champs (voir
   [Images](../user-guide/print-forms.md#images)).
5. **Enregistrer**, puis vérifie le résultat avec **Aperçu PDF**.

Sur du papier avec un motif préimprimé, ajoute un scan de celui-ci sous
**Scan du formulaire vide** et calibre l’imprimante comme pour les [feuilles
de stand de l’armée](army-sheet.md#calibrate).

## 2. Choisir l’imprimante {#printer}

Sous **Stands de tir → Configurer les imprimantes** du premier stand de la
manifestation, choisis l’imprimante pour **Diplômes** et **Enregistrer**
(voir [Connecter une imprimante](connect-printer.md)). Sans cela, les
diplômes s’ouvrent en un seul PDF.

## 3. Imprimer les diplômes {#print}

1. Ouvre les **Classements** de la manifestation.
2. Clique sur **Diplômes** dans l’en-tête du programme ou du classement
   général. Le bouton apparaît dès que la société a un modèle qui imprime un
   rang.
3. Choisis le **Modèle**, pour un classement général à plusieurs classes la
   **Classe**, et les rangs avec **Du rang** et **Au rang** (1 à 3 sauf
   modification).
4. Clique sur **Imprimer les diplômes**, ou sur **PDF** pour les ouvrir.

![Diplômes depuis un classement](../images/screenshots/certificates.png){
.screenshot }

Une page s’imprime par tireur classé. Le diplôme d’un programme connaît
aussi les résultats du tireur, p. ex. les passes.
