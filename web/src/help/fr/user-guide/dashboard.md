# Tableau de bord

Le tableau de bord est la vue de l’accueil d’une manifestation : qui tire
sur quelle cible, qui attend, qui a terminé.

![Dashboard](../images/screenshots/dashboard.png){ .screenshot }

- **En tir** : inscriptions sur une cible, avec le temps écoulé depuis le
  début, les boutons **Résultats**, **Modifier** et **Terminer**. Si le
  module *Distribution de munition* est activé, **Munition** note qu’elle a
  été remise et un badge vert *munition* le montre ensuite (voir
  [Aides](helpers.md)).
- **En attente** : inscriptions sans cible (par exemple quand toutes les
  cibles sont occupées). Sans [attribution des cibles](events.md#modules)
  (ou sans stand), il n’y a ni cibles ni *En tir* : la liste s’appelle alors
  **Inscrits** et propose **Résultats**, **Modifier** et **Terminer** ; une
  inscription la quitte dès que ses résultats sont complets.
- **Terminés** : inscriptions dont tous les programmes ont des résultats,
  avec le badge *distinctions à remettre*, qui ouvre le dialogue des
  résultats pour la remise, ou ensuite *distinctions remises* (voir
  [Résultats](results.md#handout)).
- **Cibles** : chaque cible est *libre*, *occupée* (avec le tireur),
  *réservée* (avec la réservation et sa fin) ou *hors service*. Une
  réservation apparaît dès son délai d’avance. Une manifestation sur
  plusieurs stands affiche une grille par stand, et la fenêtre d’inscription
  groupe les cibles par stand ; un tireur qui passe sur un autre stand est
  terminé puis inscrit à nouveau là-bas.
- Les chiffres en haut comptent les inscriptions en attente, en tir et
  terminées.
- *En attente* et *Terminé* affichent 50 inscriptions par page ; les flèches
  sous la liste tournent les pages.

Tout se met à jour en direct : un résultat saisi sur un autre ordinateur ou
une réservation qui commence apparaît sans recharger. Le point à côté de
**live** dans l’en-tête montre la connexion ; après une interruption du
réseau, la page se met à jour d’elle-même. Sans connexion, une barre rouge
le signale ; rien ne peut alors être enregistré, mais les dialogues ouverts
gardent ce que tu as saisi, pour que tu puisses enregistrer à nouveau dès
que tu es de retour en ligne.

## Rechercher et scanner

Appuie sur ++f4++ n’importe où dans le tableau de bord pour aller au champ
de recherche.

![Recherche](../images/screenshots/search.png){ .screenshot }

- Tape un **numéro de licence** (chiffres seulement, zéros initiaux
  facultatifs), un **nom** (partie du prénom ou du nom) ou **prénom et nom**
  dans n’importe quel ordre.
- Chaque résultat propose **Inscrire &lt;société&gt;** pour chaque société
  du tireur, ou **Ouvrir l’inscription** si le tireur est actuellement
  inscrit.
- **Nouveau tireur** crée un tireur et ouvre directement l’inscription.
- Un scanner manuel qui lit le [code d’inscription](../reference/codes.md)
  sur une feuille de stand ou une bande de contrôle tape le code et Entrée :
  les résultats de cette inscription s’ouvrent directement.

## Terminer le tir

**Terminer** clôt une inscription avant que tous les résultats soient là, par exemple quand un tireur abandonne ou rend ses cibles plus tard. L’inscription passe à *Terminé* et sa cible se libère ; les résultats manquants peuvent toujours être saisis.

## Mode sombre

Le bouton soleil/lune dans l’en-tête bascule entre mode clair et sombre ; le
choix est conservé dans le navigateur.

![Tableau de bord en mode sombre](../images/screenshots/dashboard-dark.png){
.screenshot }
