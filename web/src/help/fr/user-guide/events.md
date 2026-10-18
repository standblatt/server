# Manifestations

**Manifestations** liste les manifestations de la société dans les groupes *Aujourd’hui* (entre le premier et le dernier jour de tir), *À venir*, *Passées* et *Clôturées* ; les admins en créent de nouvelles ici. Ouvrir une manifestation affiche dans la barre latérale son tableau de bord et les autres pages du déroulement ; *Changer de manifestation* en haut de la barre latérale ramène à la liste.

![Manifestations](../images/screenshots/events.png){ .screenshot }

## Configuration et déroulement {#configuration}

Une manifestation a deux parties. Le **déroulement** est ce avec quoi vous
travaillez le jour même : *Tableau de bord*, *Saisir les résultats*,
*Classements*, *Statistiques* et les pages des modules activés comme
*Réservations*. **Configuration** au bas de la barre latérale (admins
seulement) ouvre la mise en place de la manifestation ; la barre latérale
liste alors ses pages, et *Retour à la manifestation* ramène au déroulement
:

| Page | Contenu |
|---|---|
| **Général** | nom, jours, stands, organisateur et langue (voir plus bas) |
| **Modules** | les éléments dont la manifestation se sert (voir plus bas) |
| **Programmes** | ce qui peut être tiré (voir [Programmes](programmes.md)) |
| Pages des modules | p. ex. **Réservations** (avec la réservation en ligne), **Classements généraux**, **Affichages kiosque**, **Liens d'aide**, seulement tant que le module est activé |
| **Copier & clôturer** | copier pour l’année prochaine, clôturer, rouvrir et supprimer |

Le personnel travaille seulement dans le déroulement ; la configuration est
réservée aux admins.

## Réglages {#settings}

**Général** contient ce que toute manifestation possède :

| Réglage | Signification |
|---|---|
| Nom, date, stands | la date est le jour principal ; son année fixe les catégories d’âge (FST : année du concours moins année de naissance, le jour lui-même ne compte pas, ce qui importe pour le tir anticipé de décembre avant une manifestation de janvier) ; cocher chaque stand sur lequel la manifestation se tire ; ou aucun pour une manifestation sans cibles (voir plus bas) |
| Jours de tir | jours isolés avec heures de tir facultatives, pour le tir anticipé et les manifestations sur plusieurs jours (les inscriptions comptent le jour où elles sont faites) ; le jour principal en fait toujours partie |
| Période (saison) | pour les manifestations sur des semaines ou des mois, p. ex. le Tir en campagne de l’ouverture du stand jusqu’au 31 août : premier et dernier jour, plus les jours de la semaine avec leurs heures de tir ; aucun jour coché signifie chaque jour. Les jours isolés, p. ex. le week-end officiel, s’ajoutent à la période ou remplacent ses heures à cette date |
| Organisateur | imprimé sur les documents |
| Langue des documents et du kiosque | langue des PDF, des exports CSV et des écrans kiosque, indépendamment de la langue d’interface de chaque utilisateur |

## Modules {#modules}

Toutes les manifestations n’ont pas besoin de tout : une fête de tir réserve
des cibles et remet des cartes-couronnes, un programme annuel ne fait que
réunir des résultats. **Modules** active et désactive les éléments de la
manifestation ; ce qui est désactivé disparaît de la barre latérale, du
tableau de bord, de l’éditeur de programmes, des classements et des
rapports.

| Module | Ce qu’il ajoute |
|---|---|
| Réservations | cibles réservées pour des groupes et des sociétés ; sa page dans la configuration active la [réservation en ligne](online-booking.md), par laquelle les sociétés demandent des cibles avec un lien public (nécessite un stand) |
| Feuille de stand, Bande de contrôle | les documents imprimés à l’inscription ; la bande de contrôle peut s’imprimer avec un [modèle d’impression](print-forms.md#builtin) de la société |
| Distribution de munition | le bouton **Munition** sur le tableau de bord et pour la distribution de munition (voir [Aides](helpers.md#ammunition)) ; laissez-le désactivé là où les tireurs apportent la leur, p. ex. pistolet et carabine à air |
| Tentatives | des programmes qui peuvent être tirés plusieurs fois (voir [Programmes](programmes.md#attempts)) |
| Distinctions | cartes-couronnes, prix en nature et prix en espèces ; son réglage *Imprimer la valeur de la distinction sur les cartes-couronnes* met la valeur en chiffres et en lettres sur la carte, *Modèle des cartes-couronnes* choisit un [modèle d’impression](print-forms.md#builtin) de la société |
| Prix spéciaux | les meilleurs rangs de chaque catégorie d’âge (active Distinctions, voir [Programmes](programmes.md#special-prizes)) |
| Concours de sections | les sociétés classées selon les résultats de leurs tireurs ; son réglage *Société exclue du concours de sections* est en général la société organisatrice, listée *hors concours* |
| Classements généraux | classements sur plusieurs programmes ou tentatives (voir [Classements généraux](classifications.md)) |
| Affichages kiosque | classements sur des écrans dans la halle (voir [Classements et kiosque](rankings-kiosk.md)) |
| Liens d'aide | cibarres, distribution de munition et saisie des résultats sans compte (voir [Aides](helpers.md)) |

En créant une manifestation, choisissez son genre ; cela n'active que des
modules :

| Genre | Modules |
|---|---|
| Tir de fête | Réservations, Feuille de stand, Bande de contrôle, Distinctions, Prix spéciaux, Concours de sections, Affichages kiosque, Liens d'aide |
| Programme annuel ou entraînement | Tentatives, Classements généraux |
| Tir en campagne | Distinctions, Concours de sections, Affichages kiosque, Liens d'aide |
| Championnat de société | Tentatives, Classements généraux, Affichages kiosque |
| Sélection personnalisée | aucun |

Un module qui contient des données reste activé : désactivez le kiosque
seulement après avoir supprimé ses liens kiosque, les distinctions seulement
quand aucun programme n’a de niveaux de distinction, et ainsi de suite ; la
carte indique pourquoi son interrupteur est bloqué ; Réservations reste
activé tant que le lien de réservation en ligne est actif. Un module dont un
autre a besoin (Distinctions pour Prix spéciaux) reste activé tant que
celui-ci l’est.

Les heures de tir donnent aussi ses créneaux à la réservation en ligne, sauf
si la réservation fixe d’autres horaires pour un jour (voir [réservation en
ligne](online-booking.md)).

### Modules de la société {#club-modules}

Certains modules appartiennent à la société plutôt qu'à une
manifestation. **Modules de la société** dans la partie société de la barre
latérale (admins) les active et les désactive :

| Module | Ce qu'il ajoute |
|---|---|
| Modèles d’impression | des modèles pour feuilles préimprimées et la section *Modèle d’impression* de chaque programme (voir [Modèles d’impression](print-forms.md)) |

## Manifestations sans cibles {#no-lanes}

Une manifestation n’a pas besoin de stand, p. ex. un programme annuel ou un
Tir en campagne dont les résultats arrivent sur papier ou depuis
l’installation de marquage. Sans stand, il n’y a pas de grille des cibles,
et le module Réservations n’est pas disponible ; les tireurs sont inscrits
sans cible et leurs résultats sont saisis comme d’habitude. Les documents
s’ouvrent toujours en PDF, mais ne partent pas vers un agent d’impression,
car les agents d’impression appartiennent à un stand.

## Copier pour l’année prochaine {#copy}

Sous **Copier & clôturer** dans la configuration, **Copier pour l’année
prochaine** crée une nouvelle manifestation avec les mêmes stands,
organisateur, langue, réglages, modules et tous les programmes, y compris
prix, distinctions, combinaisons et prix spéciaux, ainsi que les
[classements généraux](classifications.md) ; les jours de tir et la période
suivent la date, les heures de tir hebdomadaires restent sur leurs jours de
la semaine. Les inscriptions, les résultats et les liens kiosque ne sont pas
copiés. Le dialogue propose le nom avec l’année (et un numéro en tête comme
*42.*) incrémentée et la date un an plus tard ; adaptez les deux au besoin.

![Copier une manifestation](../images/screenshots/copy-event.png){
.screenshot }

## Clôturer {#close}

Quand la manifestation est terminée et les classements vérifiés,
**clôturez**-la sous **Copier & clôturer** dans la configuration (**Clôturer
la manifestation**). La clôture publie les classements : les tireurs
disposant d’un [journal d’entraînement
personnel](../my-training/account.md#competitions) trouvent alors leurs
résultats de cette manifestation, avec leur rang dans le programme, et
peuvent les ajouter à leur journal. Tout reste visible (tableau de bord,
résultats, classements, rapports, documents), mais les inscriptions, les
résultats, les distinctions, les programmes et les réglages ne peuvent plus
être modifiés. Une bannière sur chaque page de la manifestation
l’indique. **Rouvrir** la rend à nouveau modifiable ; les résultats déjà
ajoutés par les tireurs à leur journal y restent. Copier une manifestation
clôturée pour l’année prochaine fonctionne.

## Supprimer

**Supprimer** (sous **Copier & clôturer**) retire une manifestation, mais seulement tant que personne n’y est inscrit.
