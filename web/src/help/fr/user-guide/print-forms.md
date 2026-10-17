# Modèles d’impression

Un modèle d’impression place les résultats sur une feuille préimprimée,
p. ex. les feuilles de stand que l’armée fournit pour le programme
obligatoire : les coups sont saisis à la cible ou après coup, et la feuille
est imprimée au lieu d’être remplie à la main. Remplir à la main reste
possible ; les modèles sont facultatifs.

## Créer un modèle {#create}

**Modèles d’impression** (admins, dans la partie société de la barre latérale) liste les modèles de la société. **Nouveau modèle d’impression** demande un nom et un format de papier, puis ouvre l’éditeur :

- **Largeur** et **Hauteur** de la feuille en millimètres.
- **Scan du formulaire vide** : un PNG ou JPEG de la feuille entière. Il
  apparaît derrière l’**Aperçu** pour placer les champs sur les cases du
  formulaire ; il n’est jamais imprimé.
- **Éléments** : ce qui est imprimé et où. **x** et **y** sont des
  millimètres depuis le coin supérieur gauche, **Largeur** et **Hauteur** la
  taille (une largeur de 0 laisse le texte aller jusqu’au bord droit),
  **Taille** la taille de police en points.
  - *Texte* imprime sa **Valeur**, dans laquelle les espaces réservés sont
    remplacés, p. ex. `{{shooter.name}}` ou `{{total}}`.
  - *QR* et *Code-barres* impriment leur valeur sous forme de code, *Ligne*
    trace une ligne depuis x/y sur la largeur/hauteur.
  - **Répétition** imprime l’élément plusieurs fois, décalé chaque fois de
    **dx** et **dy** ; `#` dans un espace réservé compte 1, 2, 3, …. Un
    élément `{{s1.shot#}}` avec une répétition de 10 et dx 8 remplit les dix
    cases de coups de la passe 1.
- **Espaces réservés** liste chaque valeur avec un exemple ; un clic
  l’ajoute à l’élément modifié en dernier.

**Enregistrer** actualise l’aperçu, qui montre des valeurs inventées.

## Calibrer l’imprimante {#calibrate}

Les imprimantes impriment rarement exactement là où elles devraient. Imprime
une feuille avec **Impression d’essai**, mesure de combien les valeurs sont
décalées et saisis-le comme **Décalage de l’imprimante vers la droite** et
**Décalage de l’imprimante vers le bas** (négatif pour la gauche ou le
haut). Le décalage s’applique à tous les éléments ; l’aperçu l’ignore.

## Utiliser un modèle {#use}

Dans l’onglet *Programmes* de la **Configuration** de la manifestation,
choisis le **Modèle d’impression** d’un programme ; cela fonctionne aussi
quand des tireurs sont inscrits. Chaque programme terminé propose alors
**Imprimer sur le modèle** dans les résultats : au desk, sur la page
**Saisir les résultats** et par les liens de secrétaire à la cible (voir
[Aides](helpers.md)). Avec un agent d’impression et une imprimante choisie
pour *Modèles d’impression* sous **Stands de tir → Impression**, la feuille
part vers cette imprimante (p. ex. son alimentation manuelle), sinon le PDF
s’ouvre (voir [Impression](printing.md)).

**Exporter** télécharge un modèle sous forme de fichier et **Importer** le relit, p. ex. pour partager un modèle avec une autre société.
