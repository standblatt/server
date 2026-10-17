# Modèles d’impression

Un modèle d’impression place les résultats sur une feuille préimprimée,
p. ex. les feuilles de stand que l’armée fournit pour le programme
obligatoire : les coups sont saisis à la cible ou après coup, et la feuille
est imprimée au lieu d’être remplie à la main. Remplir à la main reste
possible ; les modèles sont facultatifs.

## Créer un modèle {#create}

**Modèles d’impression** (admins, dans la partie société de la barre latérale) liste les modèles de la société. **Nouveau modèle d’impression** demande un nom et un format de papier, puis ouvre l’éditeur. Il a trois colonnes :

- À gauche, les **Champs** : tout ce qu’un modèle peut imprimer, par
  groupes, chacun avec un exemple, et une recherche. Un clic ajoute le champ
  sous l’élément le plus bas ; glissé sur la feuille, il se place là où on
  le lâche.
  - *Texte libre* imprime un texte fixe, *Code QR* et *Code-barres*
    impriment leur valeur sous forme de code, *Ligne* trace une ligne.
  - *Tireur*, *Société et manifestation* et *Résultat* impriment chacun une
    valeur, p. ex. le nom, le total ou les points d’une passe.
  - *Coups individuels* ajoute toute une rangée de cases de coups : *Coups
    passe 1* est un élément `{{s1.shot#}}` répété une fois par coup, à 8 mm
    d’intervalle.
- Au milieu, la **Feuille** montre le modèle à l’échelle sur le scan, avec
  des valeurs inventées ; **Afficher les noms des champs** montre les
  espaces réservés à la place. Des règles en millimètres longent le haut et
  la gauche, le zoom va de 50 à 200 %.
- À droite, les **Propriétés** modifient précisément l’élément sélectionné :
  sa **Valeur** (les espaces réservés comme `{{shooter.name}}` sont remplis
  à l’impression), **x** et **y** en millimètres depuis le coin supérieur
  gauche, **Largeur** et **Hauteur** (avec une largeur de 0, le texte va
  jusqu’au bord droit), **Taille de police**, **Gras** et
  **Alignement**. **Répétition** l’imprime plusieurs fois, décalé chaque
  fois de **dx** et **dy** ; `#` dans un espace réservé compte 1, 2, 3,
  …. En dessous, le **Modèle** fixe le nom, la taille de la feuille, le
  décalage de l’imprimante et le **Scan du formulaire vide** : un PNG ou
  JPEG de toute la feuille qui apparaît derrière la feuille pour placer les
  champs sur les cases du formulaire ; il n’est jamais imprimé.

**Enregistrer** enregistre le modèle ; un point sur le bouton signale des modifications non enregistrées, et quitter la page avec elles demande d’abord une confirmation. **Aperçu PDF** ouvre le modèle exactement tel que le PDF le rend, ce qui peut différer légèrement de la feuille dans le navigateur.

## Disposer les éléments {#arrange}

- On **déplace** un élément en le glissant et on le **redimensionne** par
  son coin. Le point rose sur la dernière case d’un élément répété règle
  l’espacement.
- Pendant le déplacement, l’élément **s’aimante** aux bords et aux centres
  des autres éléments et de la feuille, et un repère rose montre la ligne ;
  ailleurs, il s’aligne sur des millimètres entiers. Avec la touche Alt
  enfoncée, on le place librement.
- **Sélectionner plusieurs** avec Maj ou Ctrl et un clic, ou en traçant un
  cadre autour d’eux sur la feuille vide. Les glisser les déplace ensemble ;
  taille de police, gras et alignement s’appliquent alors à tous.
- Un **clic droit** ouvre le menu : *Dupliquer*, *Mettre au premier plan*,
  *Alignement* (bords, centres, répartir uniformément), *Centrer sur la
  feuille* et *Supprimer*.

| Touche | Effet |
|---|---|
| Flèches | déplacent la sélection de 0.5 mm, avec Maj de 5 mm |
| Ctrl+Z / Ctrl+Y | annuler / rétablir |
| Ctrl+A | tout sélectionner |
| Ctrl+D | dupliquer |
| Suppr | supprimer la sélection |
| Échap | annuler la sélection |
| Ctrl+S | enregistrer |
| Ctrl+molette | zoomer |

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
