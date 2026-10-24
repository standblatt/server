# Sociétés, comptes et données

Standblatt sert de nombreuses sociétés sur un seul serveur. Cette page
explique qui possède quelles données, pourquoi chacun a un seul compte pour
tout, et pourquoi les tireurs d’une société ne sont pas ceux d’une autre.

## Les sociétés sont des mandants {#tenants}

Dans Standblatt, une **société** est un mandant : ses stands, tireurs, liste
de sociétés, manifestations, programmes, modèles d’impression et
utilisateurs lui appartiennent, et à personne d’autre. Les sociétés ne
voient jamais les données des autres, pas même qu’un tireur ou une
manifestation existe ailleurs. La séparation fait partie de chaque requête
du serveur, ce n’est pas un filtre de l’interface.

## Un compte, un rôle par société {#accounts}

Les personnes ne sont pas des données : la même personne peut être
secrétaire d’une société, aider à l’accueil de la fête de tir d’une autre
société et tenir un journal d’entraînement personnel. Standblatt a donc des
**comptes globaux**. Un compte se connecte une fois et atteint chaque
société dont il fait partie, avec un rôle par société qui décide de ce qu’il
peut y faire (voir [Rôles](../user-guide/users.md#roles)).

Une société ne possède pas le compte, seulement l’affiliation. Retirer un
membre lui enlève l’accès à cette société ; le compte, son journal
d’entraînement et ses autres sociétés restent. Les sociétés invitent des
personnes par e-mail, et la personne ne devient membre qu’en acceptant le
lien ; ainsi personne n’est ajouté à une société à son insu.

Les aides d’une manifestation n’ont besoin d’aucun compte : un [lien
d’aide](../user-guide/helpers.md) donne accès à une manifestation et à une
tâche. Un [jeton API](../user-guide/login-language.md#api-tokens) agit pour
son propriétaire et jamais avec plus que le rôle du propriétaire dans une
société.

## Données de base des tireurs par société {#shooters}

Un **tireur** n’est pas un compte. C’est un enregistrement dans les données
de base d’une société : nom, numéro de licence, date de naissance, adresse
et les sociétés auxquelles le tireur appartient. Chaque société tient ses
propres tireurs, importés d’un fichier CSV, saisis à l’accueil ou repris
d’une réservation en ligne.

Il n’existe pas de registre commun de tous les tireurs licenciés : les
données des membres de la fédération ne font pas partie de Standblatt, et
comme les sociétés ne voient jamais les données des autres, les tireurs
saisis par une société restent invisibles pour les autres. Le prix à payer :
un tireur qui tire aux manifestations de trois sociétés est enregistré trois
fois, une fois dans chaque société. Là où les enregistrements doivent se
rejoindre, pour le journal d’entraînement d’un tireur, le numéro de licence
et la date de naissance s’en chargent (voir [Clôturer une
manifestation](closing.md#my-training)).

## La société du tireur et la société qui utilise Standblatt {#shooter-club}

Deux choses différentes s’appellent société :

- La **société qui utilise Standblatt** (*tenant club*) est celle qui
  organise la manifestation.
- La **société du tireur** est une entrée de la liste **Sociétés** de cette
  société : la société à laquelle un tireur appartient, p. ex. une société
  invitée à une fête de tir. Ces sociétés n’ont pas accès à Standblatt ; ce
  sont des noms avec un numéro FST.

Un tireur peut appartenir à plusieurs sociétés et s’inscrit pour l’une
d’elles. Cette société est imprimée sur les documents et compte dans les
classements et le concours de sections. Comme les sociétés sont reconnues
par leur nom sans tenir compte des majuscules et des espaces, et que de
nouvelles ne sont créées que volontairement, une faute de frappe ne fait pas
d’une société deux ; si c’est arrivé quand même, une fusion corrige cela,
même dans les manifestations clôturées.

## Mises à jour en direct {#live}

Plusieurs personnes travaillent en même temps sur la même manifestation :
l’accueil inscrit, les cibarres saisissent les résultats aux cibles, les
écrans kiosque affichent les classements dans la halle. Dès que quelque
chose est enregistré, le serveur prévient chaque page ouverte de cette
manifestation, et chaque page récupère ce qui a changé. Un résultat saisi à
une cible apparaît sur le tableau de bord et les écrans kiosque sans
recharger. Après une coupure du réseau, une page se remet à jour toute
seule, rien n’est donc manqué.

## Hors ligne {#offline}

Le travail de la société a besoin d’une connexion. Une cible ne peut être
attribuée qu’une fois, un tireur ne peut tirer un programme qu’une fois, et
un résultat enregistré sur un appareil doit se trouver sur tous les autres ;
seul le serveur peut en décider pour tout le monde en même temps. Sans
connexion, les pages de la société et les liens d’aide affichent une barre
rouge et n’enregistrent rien, mais les dialogues ouverts gardent ce qui a
été tapé, pour qu’on puisse l’enregistrer dès que la connexion revient.

*Mon entraînement* est différent, car personne d’autre ne travaille sur un journal personnel. Une séance enregistrée sans connexion reste sur l’appareil et est envoyée plus tard, sans jamais être créée deux fois, et l’application installée s’ouvre sans connexion avec les dernières données. Voir [Séances inachevées et hors ligne](../my-training/logging.md#unfinished-sessions-and-offline).
