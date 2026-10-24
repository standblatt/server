# Comment fonctionne l’impression

Une fête de tir imprime beaucoup : une feuille de stand et une bande de
contrôle pour chaque inscription, des feuilles de résultats, des
cartes-couronnes, des formulaires. Standblatt est une application web, et
cela détermine le fonctionnement de l’impression.

## Chaque document est un PDF {#pdf}

Le serveur génère chaque document en PDF, dans la langue de la manifestation
: feuilles de stand, bandes de contrôle, feuilles de résultats,
cartes-couronnes, modèles d’impression, diplômes et classements. Le même PDF
est imprimé au stand, ouvert dans le navigateur ou partagé depuis un
téléphone ; un document a donc le même aspect, où qu’il sorte. Les
classements sont seulement téléchargés, en PDF ou CSV, et jamais envoyés à
une imprimante du stand.

## Pourquoi un agent d’impression {#agent}

Un navigateur ne peut pas imprimer tout seul : chaque impression demande un
dialogue, une imprimante et un clic. À un accueil chargé, c’est trop lent,
et une feuille de stand doit sortir de l’imprimante à côté de l’accueil
pendant que la bande de contrôle sort de l’imprimante d’étiquettes.

L’**agent d’impression** résout cela. C’est un petit programme sur un
ordinateur à côté des imprimantes : l’application Standblatt Print Agent
sous Windows ou sur un Mac, ou un agent sans interface sur un ordinateur
Linux comme un Raspberry Pi. Quand quelqu’un appuie sur un bouton
d’impression, le serveur met le PDF dans une file d’attente ; l’agent le
récupère et l’imprime sur l’imprimante choisie sans rien demander. L’agent
se connecte à Standblatt, et non l’inverse ; rien au stand ne doit donc être
accessible depuis Internet.

Chaque tâche d’impression est confirmée par l’agent. Une tâche qui n’est pas
confirmée dans les deux minutes est redistribuée, jusqu’à trois essais ;
**Tâches récentes** montre ce qui s’est passé.

## Les agents et les imprimantes appartiennent à un stand {#range}

Les imprimantes se trouvent dans un stand, pas dans une manifestation ;
l’agent d’impression et les réglages d’impression appartiennent donc au
[stand](../getting-started/concepts.md#ranges-and-lanes). Chaque
manifestation sur ce stand les utilise sans configuration propre. Les
documents d’une inscription s’impriment sur le stand de sa cible ; une
manifestation sur plusieurs stands imprime donc au stand où le tireur
tire. Un ordinateur peut servir plusieurs stands, p. ex. *Pistolet 25 m* et
*Pistolet 50 m* dans la même halle.

Les réglages choisissent une imprimante **par document**, car les documents
demandent des imprimantes différentes : la bande de contrôle va à une
imprimante d’étiquettes, la feuille de stand à un bac A5, les
cartes-couronnes et les modèles d’impression à une alimentation manuelle
avec les feuilles préimprimées. Les options d’impression que l’agent signale
(bac, format du papier, recto verso, etc.) sont enregistrées avec le choix.

Un agent se connecte soit avec un jeton copié dans ses réglages, soit, avec
l’application, par un code affiché sur l’ordinateur et confirmé dans
Standblatt (voir [Impression](../user-guide/printing.md#connect)). Dans les
deux cas, le serveur ne garde qu’une empreinte du jeton, jamais le jeton
lui-même.

## Quand il n’y a pas d’imprimante {#fallback}

Un document ne va à l’imprimante que si son stand a une imprimante choisie
pour ce document et si un agent du stand était en ligne au cours de la
dernière minute. Sinon, le PDF s’ouvre à la place : dans un nouvel onglet
sur un ordinateur, dans la feuille de partage sur un téléphone ou une
tablette (imprimer, enregistrer dans les fichiers, envoyer), pour que
l’application reste ouverte. Rien n’est perdu quand l’agent est éteint ;
l’accueil imprime simplement à la main.

Il en va de même pour les manifestations sans stand : il n’y a pas d’agent à
qui envoyer, chaque document s’ouvre donc en PDF.

## Modèles d’impression sur feuilles préimprimées {#forms}

Certains documents ne sont pas conçus par Standblatt mais fournis sur papier
: les cartes-couronnes, les feuilles de stand que l’armée fournit pour le
programme obligatoire, les diplômes propres à une société. Un **modèle
d’impression** décrit une telle feuille : où va chaque valeur sur le papier,
en millimètres depuis le coin supérieur gauche. Imprimer un modèle ne pose
que les valeurs sur la feuille qui se trouve déjà dans l’imprimante.

Un modèle a une **source de données** : les résultats d’un programme, une
inscription ou un rang. La source détermine quelles valeurs le modèle peut
imprimer et où il est utilisé : les modèles de résultats s’impriment une
fois un programme terminé ou comme cartes-couronnes, ceux d’une inscription
comme bandes de contrôle ou cartes de participant, ceux d’un rang comme
diplômes depuis les classements. La carte-couronne et la bande de contrôle
sont elles-mêmes des modèles intégrés qu’une société peut copier et adapter.

Un scan de la feuille vide peut se trouver derrière le modèle dans
l’éditeur, pour placer les champs sur ses cases. Le scan ne sert qu’au
placement et n’est jamais imprimé.

## Le calibrage en principe {#calibration}

Les imprimantes impriment rarement exactement où elles le devraient :
l’alimentation du papier et les bacs décalent l’impression d’un millimètre
ou deux. Un modèle juste à l’écran peut donc atterrir à côté des cases sur
le papier. Le **décalage de l’imprimante** du modèle déplace chaque élément
de la même valeur pour compenser : imprime une feuille d’essai, mesure de
combien les valeurs sont décalées et saisis la différence. Le décalage fait
partie du modèle ; un modèle imprimé sur une autre imprimante qui décale
différemment a donc besoin d’une copie avec son propre décalage. L’aperçu
ignore le décalage et montre le modèle tel qu’il a été conçu. Les étapes se
trouvent dans [Imprimer les résultats sur les feuilles de stand de
l’armée](../how-to/army-sheet.md#calibrate).
