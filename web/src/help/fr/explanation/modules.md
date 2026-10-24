# Modules et genres de manifestations

Les sociétés organisent des manifestations très différentes : une fête de
tir avec un accueil, l’attribution des cibles et des cartes-couronnes, un
tir en campagne le week-end officiel, un programme annuel tiré pendant la
saison, un championnat de société. Standblatt n’a pas de code séparé pour
chacune d’elles. Chaque manifestation est le même genre de chose ; les
manifestations ne diffèrent que par leurs programmes, leurs stands et leurs
jours, et les **modules** qu’elles activent.

## Le noyau et les modules {#core}

Chaque manifestation a le même noyau : nom, jours, stands et langue,
programmes avec leurs passes et leurs prix, inscriptions, résultats,
classements, statistiques et clôture. Tout le reste est un module :
l’attribution des cibles, les réservations avec la réservation en ligne, la
feuille de stand et la bande de contrôle, la distribution de munition, les
tentatives, les distinctions, les prix spéciaux, le concours de sections,
les classements généraux, les affichages kiosque et les liens d’aide.

Un module désactivé a disparu, il n’est pas seulement masqué : il disparaît
de la barre latérale, du tableau de bord, de l’éditeur de programmes, des
classements et des rapports, et le serveur refuse aussi ses fonctions, si
bien que ni un lien d’aide ni un jeton API ne peuvent l’utiliser. Un
programme annuel sans distinctions n’a donc pas de section de distinctions
dans ses programmes ni de rapports de distinctions ; personne n’a à passer
par-dessus ce que la manifestation n’utilise pas.

Certains modules en demandent d’autres : *Prix spéciaux* active
*Distinctions*, car un prix spécial est une distinction, et *Réservations*
demande *Attribution des cibles*, car une réservation bloque des cibles. Un
module reste activé tant qu’un module qui en a besoin est activé.

## Les genres de manifestations ne sont qu’un point de départ {#kinds}

**Nouvelle manifestation** demande de quel genre de manifestation il s’agit : tir de fête, programme annuel ou entraînement, tir en campagne, championnat de société ou sélection personnalisée. Le genre ne fait qu’activer un ensemble de modules (la liste se trouve sous [Genres de manifestations](../reference/modules.md#kinds)) ; la manifestation ne s’en souvient pas. Ensuite, les modules se modifient comme ceux de n’importe quelle autre manifestation : un tir en campagne qui veut aussi des classements généraux les active, et une fête de tir sans accueil désactive l’attribution des cibles.

C’est pourquoi un nouveau genre de manifestation ne demande aucune nouvelle
fonction : si une combinaison de modules le couvre, il peut être organisé
dès aujourd’hui. Une société qui organise chaque année une manifestation
sous la même forme l’enregistre plutôt comme
[modèle](../user-guide/events.md#templates) ; un modèle est une
manifestation avec des modules, des programmes et des classements généraux,
pas un genre.

## Pourquoi un module qui contient des données reste activé {#in-use}

Désactiver un module qui contient des données masquerait des données qui
comptent encore : des paliers de distinctions que des tireurs ont atteints,
des liens kiosque qui affichent encore les classements dans la halle, des
programmes qui permettent plusieurs tentatives. Standblatt garde donc un tel
module activé jusqu’à ce que ses données aient disparu, et la carte du
module indique pourquoi son interrupteur est verrouillé :

| Module | Reste activé tant que |
|---|---|
| Attribution des cibles | une inscription a une cible |
| Réservations | le lien de réservation en ligne est actif |
| Tentatives | un programme permet plus d’une tentative |
| Distinctions | un programme a des paliers de distinctions |
| Prix spéciaux | un programme décerne des prix spéciaux |
| Classements généraux | la manifestation a un classement général |
| Affichages kiosque | la manifestation a un lien kiosque |
| Liens d'aide | la manifestation a un lien d’aide |

La feuille de stand, la bande de contrôle, la distribution de munition et le
concours de sections ne contiennent pas de données propres et peuvent être
activés ou désactivés à tout moment.

## Configuration et déroulement {#configuration}

Une manifestation se prépare une fois et se déroule de nombreuses fois : les
programmes, les prix et les modules sont mis en place avant l’arrivée du
premier tireur, tandis que l’accueil travaille le jour même. Standblatt
sépare les deux. Le **déroulement** (tableau de bord, résultats,
classements, statistiques et pages des modules comme les réservations)
demande le droit *Gérer les manifestations* ; la **configuration** (réglages
généraux, modules, programmes, pages de configuration des modules, copie et
clôture) demande *Configurer les manifestations*. Une aide à l’accueil ne
peut pas modifier un prix par erreur, et la configuration ne gêne pas
l’accueil.

## Modules de la société et modules de la manifestation {#club-modules}

La plupart des modules appartiennent à une manifestation, car chaque
manifestation décide si elle remet des distinctions ou utilise des écrans
kiosque. Certains éléments sont utilisés par toutes les manifestations d’une
société et appartiennent donc à la société : les **modèles d’impression**
avec leur bibliothèque d’images sont dessinés une fois et utilisés par les
programmes de chaque manifestation. Ces **modules de la société** s’activent
sous [Modules de la société](../user-guide/events.md#club-modules). Une
nouvelle société a les modèles d’impression activés ; comme les modules de
manifestation, le module reste activé tant que la société a des modèles ou
des images.

## Manifestations avec et sans cibles {#lanes}

Les cibles ne comptent que là où un accueil les attribue. Standblatt sépare
donc le stand des cibles :

| La manifestation a | Ce qui en découle |
|---|---|
| un stand et *Attribution des cibles* | l’accueil d’une fête de tir : grille des cibles, réservations, réservation en ligne, liens de cibarre et état des cibles sur les écrans kiosque |
| un stand, mais pas *Attribution des cibles* | les tireurs sont inscrits sans cible ; les documents vont toujours à l’agent d’impression du stand et les programmes peuvent toujours être liés au stand |
| pas de stand | pas de cibles ni rien qui en a besoin ; les documents s’ouvrent en PDF, car les agents d’impression appartiennent à un stand |

Un programme annuel ou un tir en campagne dont les résultats arrivent sur
papier ou depuis le système de marquage n’a pas besoin de cibles, mais la
société peut quand même vouloir son imprimante au stand ; c’est la ligne du
milieu. Les modules qui ont besoin de cibles se désactivent d’eux-mêmes
quand la manifestation n’a pas de stand.

Les résultats arrivent dans une manifestation sans cibles de la même façon
que dans toute autre : par **Saisir les résultats**, un fichier CSV ou les
résultats d’une inscription (voir [Saisir les
résultats](../user-guide/results-entry.md)).

## La copie garde la forme {#copy}

**Copier pour l’année prochaine** et les modèles reprennent les modules activés et ceux qui contiennent des données, avec les programmes et les classements généraux. Une manifestation copiée de l’année passée commence donc avec la même forme, sans ses inscriptions ni ses résultats.
