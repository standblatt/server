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

Chaque droit et ce qu’il permet :
[Droits](../reference/permissions.md#permissions). Chaque société commence
avec les rôles **Admin**, **Comité**, **Aide** et **Membre** (voir [Rôles
par défaut](../reference/permissions.md#default-roles)).

**Nouveau rôle** ajoute un rôle propre, qui ne voit que les manifestations tant que vous ne cochez pas davantage. Le crayon renomme un rôle (un nom vide garde le nom par défaut), la corbeille le supprime dès que plus personne ne l’a. L’application ne montre à chaque membre que les pages et boutons que son rôle permet : un membre sans *Gérer les manifestations* ouvre une manifestation sur ses classements.

Les membres avec *Utilisateurs et société* mais sans le rôle Admin restent
dans leurs propres droits : ils donnent, modifient, retirent et invitent
seulement avec des rôles dont ils ont eux-mêmes les droits, et ne créent ou
modifient des rôles que dans ces limites. Seuls les admins nomment,
modifient ou retirent donc des admins.

Les comptes sont globaux (voir [Un compte, un rôle par
société](../explanation/data.md#accounts)) : un compte est supprimé par
l’exploitant de la plateforme dans le panneau d’administration. Pour un mot
de passe oublié, voir [Mot de passe oublié](login-language.md#password).
