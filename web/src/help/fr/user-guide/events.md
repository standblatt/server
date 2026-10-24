# Manifestations

**Manifestations** liste les manifestations de la société dans les groupes *Aujourd’hui* (entre le premier et le dernier jour de tir), *À venir*, *Passées* et *Clôturées* ; les rôles avec le droit *Configurer les manifestations* y en créent de nouvelles. Ouvrir une manifestation affiche dans la barre latérale son tableau de bord et les autres pages de l’exploitation ; *Changer de manifestation* en haut de la barre latérale ramène à la liste.

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

La configuration demande le droit *Configurer les manifestations*, le
déroulement *Gérer les manifestations* (voir
[Rôles](users.md#roles)). Pourquoi les deux sont séparés : [Configuration et
déroulement](../explanation/modules.md#configuration).

## Réglages {#settings}

**Général** contient ce que toute manifestation possède :

| Réglage | Signification |
|---|---|
| Nom, date, stands | la date est le jour principal ; son année fixe les catégories d’âge (FST : année du concours moins année de naissance, le jour lui-même ne compte pas, ce qui importe pour le tir anticipé de décembre avant une manifestation de janvier) ; cocher chaque stand sur lequel la manifestation se tire ; ou aucun pour une manifestation sans cibles (voir plus bas) |
| Jours de tir | jours isolés avec heures de tir facultatives, pour le tir anticipé et les manifestations sur plusieurs jours (les inscriptions comptent le jour où elles sont faites) ; le jour principal en fait toujours partie |
| Période (saison) | pour les manifestations sur des semaines ou des mois, p. ex. le Tir en campagne de l’ouverture du stand jusqu’au 31 août : premier et dernier jour, plus les jours de la semaine avec leurs heures de tir ; aucun jour coché signifie chaque jour. Les jours isolés, p. ex. le week-end officiel, s’ajoutent à la période ou remplacent ses heures à cette date |
| Organisateur | imprimé sur les documents |
| Langue des documents et du kiosque | langue des PDF, des exports CSV et des écrans kiosque, indépendamment de la langue d’interface de chaque utilisateur |

Les heures de tir donnent aussi ses créneaux à la réservation en ligne, sauf
si la réservation fixe d’autres horaires pour un jour (voir [réservation en
ligne](online-booking.md)).

Pas à pas pour une saison : [Organiser un programme annuel sur la
saison](../how-to/annual-programme.md), [Organiser un tir en campagne sur
une période](../how-to/feldschiessen.md).

## Modules {#modules}

**Modules** active et désactive les éléments de la manifestation ; pourquoi les manifestations ne diffèrent que par leurs modules est expliqué dans [Modules et genres de manifestations](../explanation/modules.md).

![Les modules d'une manifestation](../images/screenshots/modules.png){
.screenshot }

Chaque module est une carte avec son interrupteur ; les réglages d’un
module, p. ex. *Imprimer la valeur de la distinction sur les
cartes-couronnes* de *Distinctions*, se trouvent sur sa carte. Ce que chaque
module ajoute, ses réglages et quand il reste activé :
[Modules](../reference/modules.md#event-modules).

En créant une manifestation, choisis son genre ; cela n’active que des
modules (voir [Genres de manifestations](../reference/modules.md#kinds)).

Un module qui contient des données reste activé, et la carte indique
pourquoi son interrupteur est verrouillé (voir [Pourquoi un module qui
contient des données reste activé](../explanation/modules.md#in-use)). Un
module dont un autre a besoin, p. ex. *Distinctions* pour *Prix spéciaux* ou
*Attribution des cibles* pour *Réservations*, reste activé tant que celui-ci
l’est.

### Modules de la société {#club-modules}

Certains modules appartiennent à la société plutôt qu’à une
manifestation. **Modules de la société** dans la partie société de la barre
latérale (droit *Utilisateurs et société*) les active et les désactive ;
pour l’instant, il s’agit de *Modèles d’impression* (voir [Modèles
d’impression](print-forms.md)). Ce que chacun ajoute et quand il reste
activé : [Modules de la société](../reference/modules.md#club-modules) ;
pourquoi ils appartiennent à la société : [Modules de la société et modules
de la manifestation](../explanation/modules.md#club-modules).

## Manifestations sans cibles {#no-lanes}

Une manifestation sans cibles n’a pas de stand, ou désactive *Attribution
des cibles* ; ce que signifie chaque choix est expliqué dans [Manifestations
avec et sans cibles](../explanation/modules.md#lanes).

## Copier pour l’année prochaine {#copy}

Sous **Copier & clôturer** dans la configuration, **copier** sous *Copier
pour l’année prochaine* crée une nouvelle manifestation avec les mêmes
stands, organisateur, langue, réglages, modules et tous les programmes, y
compris prix, distinctions, combinaisons et prix spéciaux, ainsi que les
[classements généraux](classifications.md) ; les jours de tir et la période
suivent la date, les heures de tir hebdomadaires restent sur leurs jours de
la semaine. Les inscriptions, les résultats, les liens kiosque, les liens
d’aide et le lien de réservation ne sont pas copiés. Le dialogue *Copier la
manifestation* propose le nom et la date un an plus tard ; **Créer la
copie** crée la manifestation. Pas à pas : [Reprendre la manifestation de
l’année passée](../how-to/next-year.md#copy).

![Copier une manifestation](../images/screenshots/copy-event.png){
.screenshot }

## Modèles {#templates}

Une manifestation bien configurée peut devenir un **modèle** de la société :
**Enregistrer comme modèle** sur **Copier & clôturer** en crée une copie
avec ses modules, réglages, stands, programmes et classements généraux,
nommée d’après la manifestation sans son année. Les modèles sont listés sous
*Modèles* sous les manifestations ; un modèle s’ouvre dans sa configuration
et se modifie comme une manifestation, mais il n’accepte pas
d’inscriptions. **Nouvelle manifestation** liste les modèles de la société
sous les genres (voir [Partir d’un modèle de la
société](../how-to/next-year.md#template)). **Supprimer** sur **Copier &
clôturer** du modèle le supprime.

![Nouvelle manifestation avec les genres et un modèle de la
société](../images/screenshots/new-event.png){ .screenshot }

## Clôturer {#close}

Quand une manifestation est terminée et les classements vérifiés, **Clôturer
la manifestation** sur **Copier & clôturer** dans la configuration la
clôture ; un bandeau sur chaque page de la manifestation l’indique alors, et
**Rouvrir** la rend de nouveau modifiable. Ce que la clôture signifie pour
les classements, Mon entraînement et les modifications : [Clôturer une
manifestation](../explanation/closing.md).

## Supprimer

**Supprimer** (sous **Copier & clôturer**) retire une manifestation, mais seulement tant que personne n’y est inscrit. Une manifestation clôturée ne peut pas être supprimée ; rouvre-la d’abord.
