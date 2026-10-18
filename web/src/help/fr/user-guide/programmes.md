# Programmes

Les programmes appartiennent à une manifestation et se gèrent sous
**Configuration** → **Programmes** (admins seulement).

![La liste des programmes](../images/screenshots/setup.png){ .screenshot }

## Ajouter des programmes

- **Nouveau programme** ouvre une page de programme vide.
- **Ajouter depuis … → Bibliothèque de programmes** : un modèle de la
  bibliothèque globale, tenue par l’exploitant de la plateforme.
- **Ajouter depuis … → Manifestation précédente** : copie tous les
  programmes d’une autre manifestation de la société, y compris combinaisons
  et prix spéciaux. Les programmes dont le nom et la variante existent déjà
  sont sautés.
- **Ajouter depuis … → Fichier de programme** : importe un [fichier de
  programme](../reference/programme-file.md) exporté d’une autre
  manifestation ou société.

## La page du programme {#page}

Un clic sur un programme dans la liste ouvre sa page. Elle contient les
sections du programme, plus une section pour chaque
[module](events.md#modules) activé qui complète les programmes ; **Créer**
ou **Enregistrer** en bas les enregistre ensemble.

![Page du programme](../images/screenshots/programme-editor.png){
.screenshot }

- **Nom** et **Variante** : les programmes de même nom sont des variantes
  les uns des autres ; un tireur s’inscrit à l’une d’elles. Le champ du nom
  propose les noms des programmes de la manifestation, et dès que le nom
  correspond à l’un d’eux, la page indique que tu crées une variante et
  avertit si le nom de variante est vide ou déjà pris.
- **Code-barres** : imprimé sur la feuille de stand pour le marquage
  électronique (voir [Codes](../reference/codes.md)).
- **Compte pour le concours de sections** et **Afficher sur les écrans
  kiosque** : seulement tant que les modules *Concours de sections* et
  *Affichages kiosque* sont activés.
- **Passes** : coups, points maximum, si les mouches sont comptées et si les
  **Coups individuels** sont saisis, par passe. Avec les coups individuels,
  chaque coup peut être saisi et la passe en est la somme, p. ex. pour les
  imprimer sur un formulaire ; le total de la passe seul peut toujours être
  saisi.
- **Tentatives** (module *Tentatives*, voir plus bas) : combien de fois un
  tireur peut tirer le programme lors de la manifestation.
- **Prix (CHF)** : seules les catégories d’âge cochées peuvent tirer le
  programme. Avec plusieurs tentatives, un second champ par catégorie est le
  prix de chaque tentative supplémentaire (*rachat*) ; vide signifie le même
  prix, 0 signifie que le programme est payé une fois.
- **Distinctions** (module *Distinctions*) : nom, type (carte-couronne, prix
  en nature ou prix en espèces versé en argent), valeur et total minimal par
  catégorie d’âge ; laisse une catégorie vide là où la distinction n’est pas
  disponible.

Dès que des tireurs sont inscrits à un programme, les **passes** sont figées
pour que les résultats gardent leur sens ; la page l'indique et les
grise. Tout le reste reste modifiable : nom, variante, description,
code-barres, les options, prix, noms des passes, noms, valeurs et limites
des distinctions. Les catégories d'âge et les distinctions peuvent être
**ajoutées** mais pas retirées.

La section **Utilisation** d’un programme enregistré est enregistrée
immédiatement à chaque modification, même quand des tireurs sont inscrits :

- **Actif** : les programmes inactifs ne peuvent pas être choisis pour de
  nouvelles inscriptions ; les inscriptions existantes les gardent.
- **Combinable avec** : le programme dont celui-ci peut reprendre les passes
  à l’inscription (pas pour les programmes à plusieurs tentatives).
- **Stand de tir** (manifestations sur plusieurs stands) : le stand sur
  lequel le programme se tire, p. ex. *Pistolet 50 m* ; l’inscription ne
  propose alors que ses cibles et les écrans kiosque groupent les
  classements par stand. *indifférent* laisse le choix ouvert. Des
  programmes de stands différents demandent des inscriptions séparées.
- **Modèle d’impression** (quand la société a des [modèles
  d’impression](print-forms.md)) : le modèle sur lequel les résultats du
  programme sont imprimés ; *aucun* n’imprime pas de modèle.

Les **Prix spéciaux** (module *Prix spéciaux*, voir plus bas) suivent dans une section à part. **Exporter** en haut de la page télécharge le programme sous forme de fichier ; **Supprimer** dans le menu du programme dans la liste n’est possible que tant que personne n’est inscrit.

## Tentatives {#attempts}

Par défaut, un tireur tire un programme une fois par
manifestation. **Tentatives par tireur** en permet davantage, p. ex. 3 pour
une cible chance avec rachats à une fête de tir, ou **Illimité** pour un
programme annuel tiré sur toute la saison.

- **Tentative qui compte** détermine quelle tentative compte pour le
  classement, le concours de sections, les prix spéciaux et *Mon
  entraînement* : **La meilleure** (selon les règles du classement), **La
  première** ou **La dernière**. Le classement l’affiche comme *2/3* : la
  deuxième de trois tentatives compte.
- **Seule la tentative qui compte reçoit une distinction** : désactivé,
  chaque tentative peut recevoir sa distinction (une carte-couronne par
  rachat) ; activé, seule la tentative qui compte en reçoit une, et la
  distinction passe à une meilleure tentative dès qu’elle est tirée.
- À l’inscription, un programme à plusieurs tentatives affiche un compteur :
  deux tentatives à la fois donnent deux cartes de résultats, *Tentative 1*
  et *Tentative 2*. D’autres tentatives peuvent aussi être ajoutées plus
  tard dans une nouvelle inscription, ou en saisissant un autre résultat
  sous **Saisir les résultats**, jusqu’à ce que le tireur n’ait plus de
  tentatives.

Le nombre de tentatives peut être modifié pendant que des tireurs sont
inscrits, mais pas en dessous des tentatives qu’un tireur a déjà. Les
programmes à plusieurs tentatives ne peuvent pas être combinés.

## Prix spéciaux {#special-prizes}

La section **Prix spéciaux** de la page du programme récompense les
meilleurs rangs de chaque catégorie d’âge, par exemple les rangs 1 à 3 dans
chaque catégorie d’au moins 10 participants, en ne comptant que les tireurs
dont la société a au moins 4 participants dans ce programme. Les tireurs
sautés laissent la place au suivant. La règle exacte figure dans les [règles
de concours](../reference/rules.md#special-prizes) ; les gagnants sont dans
le rapport *Prix spéciaux* sous [Statistiques](statistics-reports.md).

![Prix spéciaux d’un programme](../images/screenshots/special-prizes.png){
.screenshot }

**Enregistrer les prix spéciaux** les enregistre ; rangs 0 désactive les prix spéciaux. Les réglages peuvent être modifiés à tout moment, même après la manifestation.
