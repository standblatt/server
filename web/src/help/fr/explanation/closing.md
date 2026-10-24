# Clôturer une manifestation

Tant qu’une manifestation se déroule, tout peut encore y changer : un
résultat est corrigé, une inscription est déplacée, un programme reçoit un
nouveau prix. À un moment donné, l’organisateur doit dire *voici le
résultat*. C’est ce que fait **Clôturer la manifestation** : elle transforme
une manifestation en cours en une manifestation terminée.

## Ce que signifie clôturer {#meaning}

Clôturer fait deux choses à la fois :

- **Les classements deviennent définitifs et sont publiés.** Jusque-là, ils
  étaient un état de travail ; désormais, ils sont le résultat de la
  manifestation.
- **La manifestation passe en lecture seule.** Les inscriptions,
  participations, résultats, distinctions, programmes avec leurs passes,
  prix et paliers de distinctions, les classements généraux et les réglages
  de la manifestation ne peuvent plus être modifiés.

Tout reste visible : le tableau de bord, les résultats, les classements, les
statistiques, les rapports et les documents. Un bandeau sur chaque page de
la manifestation indique qu’elle est clôturée.

Les deux vont ensemble. Un classement ne peut être définitif que si plus
rien de ce qui se trouve en dessous ne peut changer, et Standblatt l’impose
dans la base de données elle-même, pas seulement en masquant des boutons :
chaque modification des données d’une manifestation clôturée est refusée,
qu’elle vienne de l’application, d’un lien d’aide ou d’un jeton API.

## Pourquoi les manifestations clôturées sont en lecture seule {#read-only}

Une fois une manifestation clôturée, d’autres choses s’appuient sur elle :

- Les tireurs copient leurs résultats et leurs rangs dans leur journal
  d’entraînement personnel.
- Le rapport financier, les distinctions et le concours de sections ont été
  remis et réglés.
- La manifestation de l’année prochaine a peut-être été copiée à partir
  d’elle.

Si les résultats pouvaient encore changer, le journal d’un tireur pourrait
afficher un rang que le classement n’a plus, et les rapports ne
correspondraient plus à ce qui a été payé et remis. C’est aussi pourquoi ce
qui se passe à l’accueil, comme marquer les distinctions comme remises, doit
être fait avant la clôture.

## Le lien avec Mon entraînement {#my-training}

Chaque utilisateur de Standblatt a un journal d’entraînement personnel, [Mon
entraînement](../my-training/index.md). Un tireur qui saisit un **numéro de
licence** et une **date de naissance** dans son profil est apparié aux
tireurs de toutes les sociétés : partout où les données de base d’une
société ont un tireur avec le même numéro de licence (les zéros initiaux ne
comptent pas) et la même date de naissance, les résultats de ce tireur
appartiennent à l’utilisateur.

Seuls les résultats des manifestations **clôturées** sont proposés, car seul
leur classement est définitif. Chacun est accompagné du rang dans son
programme et du nombre de participants, p. ex. *Rang 12 sur 87*. Rien n’est
ajouté tout seul : le tireur choisit une discipline et ajoute le résultat,
qui devient une séance de compétition du journal avec ses passes et son
rang.

L’appariement fonctionne entre sociétés sans qu’aucune société ne cède ses
données : les données de base des sociétés restent où elles sont,
l’appariement ne les lit que pour l’utilisateur à qui elles appartiennent,
et seuls les propres résultats du tireur sont affichés. Un tireur sans
numéro de licence ou sans date de naissance dans son profil, ou dont la
société a enregistré une autre date de naissance, ne reçoit aucun résultat
de compétition. Voir [Profil et résultats de
compétition](../my-training/account.md#competitions).

## Rouvrir {#reopen}

**Rouvrir** rend une manifestation clôturée de nouveau modifiable, p. ex. quand un résultat erroné apparaît après la clôture. C’est prévu pour des corrections, pas comme un second état de travail :

- Les résultats que des tireurs ont déjà ajoutés à leur journal y restent
  tels qu’ils ont été ajoutés ; une correction après coup ne les atteint
  pas.
- Les résultats pas encore ajoutés ne sont plus proposés dans *Mon
  entraînement* jusqu’à ce que la manifestation soit de nouveau clôturée.

Clôture de nouveau la manifestation une fois la correction faite.

Une manifestation clôturée ne peut pas être supprimée ; rouvre-la
d’abord. Elle peut être copiée pour l’année prochaine pendant qu’elle est
clôturée.
