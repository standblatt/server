# Modèles d’impression

Un modèle d’impression place les résultats sur une feuille préimprimée,
p. ex. les feuilles de stand que l’armée fournit pour le programme
obligatoire : les coups sont saisis à la cible ou après coup, et la feuille
est imprimée au lieu d’être remplie à la main. Remplir à la main reste
possible ; les modèles sont facultatifs et peuvent être désactivés sous
**Modules de la société** (voir [Modules de la
société](events.md#club-modules)).

## Créer un modèle {#create}

**Modèles d’impression** (admins, dans la partie société de la barre latérale) liste les modèles de la société. **Nouveau modèle d’impression** demande un nom, ce que le modèle **Imprime** (voir [Source de données](#source)) et un format de papier, puis ouvre l’éditeur. Il a trois colonnes :

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
  jusqu’au bord droit), **Taille de police**, **Gras**, **Italique**,
  **Renvoyer à la ligne** (le texte se coupe en lignes dans la largeur,
  chacune aussi haute que la hauteur) et **Alignement**. **Répétition**
  l’imprime plusieurs fois, décalé chaque fois de **dx** et **dy** ; `#`
  dans un espace réservé compte 1, 2, 3, …. En dessous, le **Modèle** fixe
  le nom, la taille de la feuille, le décalage de l’imprimante et le **Scan
  du formulaire vide** : un PNG ou JPEG de toute la feuille qui apparaît
  derrière la feuille pour placer les champs sur les cases du formulaire ;
  il n’est jamais imprimé.

**Enregistrer** enregistre le modèle ; un point sur le bouton signale des modifications non enregistrées, et quitter la page avec elles demande d’abord une confirmation. Avec **Enregistrer automatiquement** (activé tant qu’on ne le désactive pas ; le navigateur s’en souvient), les modifications non enregistrées sont enregistrées chaque minute. **Aperçu PDF** ouvre le modèle exactement tel que le PDF le rend, ce qui peut différer légèrement de la feuille dans le navigateur.

## Source de données {#source}

Ce qu’un modèle **Imprime** détermine ses champs et l’endroit où il est
imprimé :

| Imprime | Champs | Imprimé depuis |
|---|---|---|
| Résultats d'un programme | tireur, société, manifestation, passes, coups isolés, total, distinction et sa valeur | les résultats d’un programme (voir [Utiliser un modèle](#use)) ou comme carte-couronne de la manifestation |
| Une inscription | tireur, société, manifestation, cible, programmes, passes, code QR de l’inscription | comme bande de contrôle de la manifestation |
| Un rang | tireur, société, manifestation, rang, total, classement, classe | diplômes depuis les classements |

La source peut être changée dans l’éditeur sous **Modèle** ; les champs que
la nouvelle source ne connaît pas n’impriment rien.

## Modèles intégrés {#builtin}

Standblatt imprime la carte-couronne et la bande de contrôle avec des
**modèles intégrés**, listés sous les modèles de la société. **Copier** en
fait un modèle de la société qui se modifie comme tout autre, p. ex. pour
déplacer les champs sur une autre carte-couronne préimprimée ou pour ajouter
le logo de la société. Une manifestation le choisit ensuite sous
**Configuration** → **Modules** : *Modèle des cartes-couronnes* dans le
module *Distinctions*, *Modèle des bandes de contrôle* dans le module *Bande
de contrôle*. Sans choix, ou si le modèle choisi est supprimé, le modèle
intégré imprime.

![Modèles d’impression avec les modèles intégrés et la bibliothèque
d'images](../images/screenshots/print-forms.png){ .screenshot }

## Documents d’inscription et diplômes {#certificates}

Un modèle qui imprime **Une inscription**, en plus de la bande de contrôle
de la manifestation, p. ex. une carte de participant ou une quittance,
apparaît dans l’inscription avec un bouton portant son propre nom ; il
s’imprime sur l’imprimante des *Modèles d’impression*.

Un modèle qui imprime **Un rang** produit des diplômes : **Diplômes** dans
l’en-tête de chaque classement et de chaque classement général sous
**Classements** demande le modèle, la classe d’un classement général et les
rangs (1 à 3 sauf modification), puis imprime une page par tireur classé, ou
les ouvre en un seul PDF. Les diplômes partent vers l’imprimante des
*Diplômes* sous **Stands de tir → Configurer les imprimantes**. Les champs
d’un rang sont le rang, le total, le classement (le programme ou le
classement général) et la classe ; le diplôme d’un programme connaît aussi
ses résultats.

![Diplômes depuis un classement](../images/screenshots/certificates.png){
.screenshot }

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

## Images {#images}

Les logos, en-têtes et autres images viennent de la bibliothèque d'images de
la société, sous la liste des **Modèles d’impression** : **Téléverser une
image** ajoute un PNG ou un JPEG de 2 Mo au plus, nommé d'après son
fichier. Dans l'éditeur, le groupe **Images** des champs les liste ; un clic
ou un glisser en pose une sur la feuille, large de 30 mm dans les
proportions de l'image. Ses propriétés choisissent une autre **Image**, sa
position et sa taille, et **Garder les proportions** ajuste à nouveau la
hauteur à la largeur. Une image peut être répétée comme tout autre élément.

Une image ne peut être supprimée que tant qu'aucun modèle ne
l'affiche. **Exporter** met les images d'un modèle dans son fichier, et
**Importer** les ajoute à la bibliothèque en réutilisant une image du même
nom que la société a déjà.

## Calibrer l’imprimante {#calibrate}

Les imprimantes impriment rarement exactement là où elles devraient. Imprime
une feuille avec **Impression d’essai**, mesure de combien les valeurs sont
décalées et saisis-le comme **Décalage de l’imprimante vers la droite** et
**Décalage de l’imprimante vers le bas** (négatif pour la gauche ou le
haut). Le décalage s’applique à tous les éléments ; l’aperçu l’ignore.

## Utiliser un modèle {#use}

Sous **Configuration** → **Programmes** de la manifestation, ouvre un
programme et choisis son **Modèle d’impression** dans la section du même nom
; cela fonctionne aussi quand des tireurs sont inscrits. Chaque programme
terminé propose alors **Imprimer sur le modèle** dans les résultats : au
desk, sur la page **Saisir les résultats** et par les liens de secrétaire à
la cible (voir [Aides](helpers.md)). Avec un agent d’impression et une
imprimante choisie pour *Modèles d’impression* sous **Stands de tir →
Configurer les imprimantes**, la feuille part vers cette imprimante
(p. ex. son alimentation manuelle), sinon le PDF s’ouvre (voir
[Impression](printing.md)).

**Exporter** télécharge un modèle avec ses images sous forme de fichier et **Importer** le relit, p. ex. pour partager un modèle avec une autre société.
