# Utilisateurs

**Utilisateurs** liste les comptes ayant accès à la société, les invitations ouvertes et les rôles de la société. La page est ouverte aux rôles avec le droit *Utilisateurs et société*, donc au moins aux admins.

![Utilisateurs](../images/screenshots/users.png){ .screenshot }

## Membres et invitations {#invitations}

- **Inviter par e-mail** : collez une ou plusieurs adresses (séparées par
  des virgules ou une par ligne, jusqu’à 300) et choisissez le rôle. Chaque
  adresse reçoit un e-mail avec un lien valable sept jours ; la liste sous
  le formulaire indique ce qui s’est passé pour chaque adresse. Les e-mails
  partent à environ quatre par seconde, une grande liste prend donc une
  minute.
- **Accepter** : le lien ouvre une page qui indique quelle société invite
  avec quel rôle. Sans compte, la personne choisit un mot de passe et est
  connectée aussitôt ; l’adresse e-mail sert d’identifiant. Qui a déjà un
  compte avec cette adresse se connecte d’abord, puis accepte. Ce n’est
  qu’ensuite que la personne apparaît parmi les membres.
- **Invitations** sous les membres liste celles qui sont ouvertes et
  expirées : **Renvoyer** envoie un nouveau lien (l’ancien ne fonctionne
  plus), **Retirer** annule l’invitation.
- **Ajouter un utilisateur** : pour un compte sans e-mail, par exemple un
  compte commun au stand : un compte existant est ajouté par son nom
  d’utilisateur ; pour un nouveau compte, saisissez aussi un mot de passe
  (au moins 10 caractères).
- **Rôle** : choisissez le rôle d’un membre dans la liste. Une société garde
  toujours au moins un admin.
- **Retirer** supprime l’accès à cette société ; le compte lui-même reste et
  garde ses autres sociétés.

## Rôles {#roles}

Ce qu’un membre peut faire dépend de son rôle. La matrice sous **Rôles** a
une colonne par rôle et une ligne par droit ; une coche donne le droit
aussitôt.

![Rôles](../images/screenshots/roles.png){ .screenshot }

| Droit | Permet |
|---|---|
| **Voir les manifestations et les classements** | manifestations, programmes, classements et classements généraux |
| **Gérer les manifestations** | inscriptions, résultats, distinctions, réservations, documents et impression |
| **Configurer les manifestations** | créer, configurer, copier, clôturer et supprimer des manifestations, avec programmes, classements généraux, liens kiosque et d’aide et réservation en ligne |
| **Statistiques et rapports** | statistiques, rapports avec montants et classements en CSV |
| **Voir les tireurs et les sociétés** | les données personnelles des données de base |
| **Modifier les tireurs et les sociétés** | saisir et modifier des tireurs et des sociétés |
| **Gérer les données de base** | importer des tireurs, renommer, fusionner et supprimer des sociétés, supprimer des tireurs |
| **Stands et imprimantes** | stands de tir, agents d’impression et réglages des imprimantes |
| **Modèles d’impression** | modèles d’impression et images (seulement tant que le module de la société est activé) |
| **Utilisateurs et société** | utilisateurs, rôles, invitations et modules de la société |

Chaque société commence avec quatre rôles :

| Rôle | Droits |
|---|---|
| **Admin** | tous ; la colonne ne peut pas être modifiée et le rôle ne peut pas être supprimé |
| **Comité** | tous sauf *Stands et imprimantes* et *Utilisateurs et société* |
| **Aide** | voir et gérer les manifestations, statistiques, voir et modifier les tireurs et les sociétés |
| **Membre** | voir les manifestations et les classements |

**Nouveau rôle** ajoute un rôle propre, qui ne voit que les manifestations tant que vous ne cochez pas davantage. Le crayon renomme un rôle (un nom vide garde le nom par défaut), la corbeille le supprime dès que plus personne ne l’a. L’application ne montre à chaque membre que les pages et boutons que son rôle permet : un membre sans *Gérer les manifestations* ouvre une manifestation sur ses classements.

Les membres avec *Utilisateurs et société* mais sans le rôle Admin restent
dans leurs propres droits : ils donnent, modifient, retirent et invitent
seulement avec des rôles dont ils ont eux-mêmes les droits, et ne créent ou
modifient des rôles que dans ces limites. Seuls les admins nomment,
modifient ou retirent donc des admins.

Les comptes sont globaux : un compte est supprimé par l’exploitant de la
plateforme dans le [panneau
d’administration](../admin-guide/admin-panel.md). Un mot de passe oublié se
réinitialise avec **Mot de passe oublié ?** sur la page de connexion si le
compte a une adresse e-mail confirmée, sinon par l’exploitant.
