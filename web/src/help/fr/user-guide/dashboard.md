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

Tout se met à jour en direct ; le point à côté de **en direct** dans
l’en-tête indique la connexion (voir [Mises à jour en
direct](../explanation/data.md#live)). Sans connexion, une barre rouge le
signale et rien ne peut être enregistré (voir [Hors
ligne](../explanation/data.md#offline)).

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

Toutes les touches : [Clavier et
scanner](../reference/shortcuts.md#dashboard).

## Terminer le tir

**Terminer** clôt une inscription avant que tous les résultats soient là, par exemple quand un tireur abandonne ou rend ses cibles plus tard. L’inscription passe à *Terminé* et sa cible se libère ; les résultats manquants peuvent toujours être saisis.
