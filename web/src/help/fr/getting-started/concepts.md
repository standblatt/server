# Notions

Cette page explique les mots qu’utilise Standblatt et comment ils
s’articulent. Les termes allemands, français et italiens se trouvent dans le
[Glossaire](../reference/glossary.md) ; les étapes de chaque tâche se
trouvent dans le guide d’utilisation.

En bref : une **société** a des **stands de tir**, des **tireurs** et des
**utilisateurs**. Elle organise des **manifestations** ; une manifestation
propose des **programmes** et active les **modules** dont elle a besoin. Une
**inscription** inscrit un tireur à la manifestation avec un ou plusieurs
programmes, chaque programme devient une **participation**, la participation
reçoit des **résultats**, et les résultats donnent des **distinctions** et
une place dans les **classements**. Quand la manifestation est **clôturée**,
les classements sont définitifs et parviennent aux journaux d’entraînement
personnels des tireurs.

## Sociétés, utilisateurs et rôles {#clubs-users-and-roles}

Une **société** est un mandant : tout dans Standblatt appartient à
exactement une société, et les sociétés ne voient jamais les données des
autres. Seul l’exploitant de la plateforme crée des sociétés, dans le
panneau d’administration. Les comptes viennent d’une invitation d’une
société, de l’inscription à [Mon entraînement](../my-training/account.md),
de **Ajouter un utilisateur** sur la page des utilisateurs d’une société ou
de l’exploitant.

Les **utilisateurs** sont des comptes personnels. Un compte peut appartenir à plusieurs sociétés, avec un rôle par société. Un rôle est un ensemble de droits comme *Gérer les manifestations* ou *Configurer les manifestations* ; chaque société commence avec les rôles *Admin*, *Comité*, *Aide* et *Membre* et peut les modifier ou ajouter les siens (voir [Rôles](../user-guide/users.md#roles)). Les sociétés invitent des personnes par e-mail ; la personne invitée rejoint la société en acceptant le lien.

Pourquoi les comptes sont globaux alors que toutes les données d’une société
restent dans la société est expliqué dans [Sociétés, comptes et
données](../explanation/data.md).

## Stands de tir et cibles {#ranges-and-lanes}

Un **stand de tir** (anglais *range*) est un stand de la société avec un
certain nombre de **cibles** (*lanes*). Les cibles peuvent être marquées
hors service. Les imprimantes et l’agent d’impression appartiennent à un
stand.

Chaque distance est un stand à part : une halle avec des cibles à 25 m et à
50 m représente deux stands, chacun avec des cibles à partir de 1. Une
manifestation peut se tirer sur plusieurs stands, sur un seul ou sur aucun
(voir [Manifestations avec et sans
cibles](../explanation/modules.md#lanes)).

## Manifestations et programmes {#events-and-programmes}

Une **manifestation** (anglais *event*) est un concours sur un ou plusieurs
stands de tir, sur un ou plusieurs jours (la date est le jour principal ;
les autres jours couvrent le tir anticipé et les manifestations sur
plusieurs jours). Elle a une langue, utilisée pour les documents imprimés et
les écrans kiosque, et des [modules](../user-guide/events.md#modules) qui
activent ce dont elle a besoin, comme l’attribution des cibles, les
distinctions ou le concours de sections.

Il n’y a pas de types de manifestations : une fête de tir, un tir en
campagne et un programme annuel sont le même genre de manifestation avec des
modules différents (voir [Modules et genres de
manifestations](../explanation/modules.md)).

Une manifestation propose des **programmes** : chaque programme a des
**passes** avec un nombre de coups et un maximum de points, des **prix** par
catégorie d’âge et des **paliers de distinctions** (cartes-couronnes, prix
en nature ou prix en espèces) avec un total minimal par catégorie d’âge. Les
programmes de même nom sont des **variantes** (par exemple *Sport* et
*Appuyé*) ; un tireur tire une seule variante. Un programme peut être
**combinable** avec un autre : il reprend alors certaines passes de cet
autre programme au lieu d’être tiré séparément.

Un programme peut permettre plusieurs **tentatives**, dont une compte. Un
**classement général** classe les tireurs sur plusieurs programmes ou
tentatives. Le lien entre classements de programme, tentatives, classements
généraux, concours de sections et prix spéciaux est expliqué dans [Comment
les classements s’articulent](../explanation/rankings.md).

## Tireurs, inscriptions et participations {#shooters-registrations-and-entries}

Les **tireurs** sont les personnes qui tirent. Chaque société tient ses propres données de base des tireurs, avec les **sociétés** auxquelles ils appartiennent (la société du tireur, à ne pas confondre avec la société qui utilise Standblatt).

Une **inscription** inscrit un tireur à la manifestation, si besoin sur une
cible, avec un ou plusieurs programmes. Chaque programme d’une inscription
est une **participation** (*entry*). Une inscription est *inscrite* (en
attente), *en tir* (sur une cible) ou *terminée* (tous les programmes ont
des résultats ; la cible est de nouveau libre).

## Résultats, distinctions et classements {#results-awards-and-rankings}

Les **résultats** sont les points par passe, avec ou sans mouches (dix intérieurs). Dès que toutes les passes sont saisies, le tireur peut recevoir l’une des **distinctions** que son total atteint dans sa catégorie d’âge. Les **classements** ordonnent tous les résultats d’un programme selon les [règles de la FST](../reference/rules.md#ranking).

Une distinction ne dépend que du propre total du tireur ; un rang dépend de
celui de tous les autres. Les classements sont en direct pendant la
manifestation et deviennent définitifs quand elle est
[clôturée](../explanation/closing.md).

## Catégories d’âge {#age-categories}

La catégorie d’âge suit la règle de la FST : l’année de la manifestation
moins l’année de naissance. Le jour de naissance ne compte pas, et l’année
est celle du jour principal de la manifestation ; un jour de tir anticipé en
décembre compte donc pour une manifestation de janvier. La catégorie est
fixée à l’inscription du tireur et détermine les prix, les distinctions et
les prix spéciaux qu’il peut obtenir. Un tireur sans date de naissance
compte comme Élite. Les catégories et leurs âges sont listés dans les
[Règles de concours](../reference/rules.md#age-categories).

## Documents et impression {#documents-and-printing}

Chaque document est un PDF généré par le serveur dans la langue de la
manifestation. Au stand, un petit agent d’impression l’envoie à la bonne
imprimante ; sans agent, le PDF s’ouvre dans le navigateur. Voir [Comment
fonctionne l’impression](../explanation/printing.md).
