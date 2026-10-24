# Droits et rôles

Le rôle d’un membre dans une société est un ensemble de droits. Les rôles se
modifient sous **Utilisateurs** → **Rôles** (voir
[Rôles](../user-guide/users.md#roles)) ; les jetons API portent une
sélection des mêmes droits.

## Droits {#permissions}

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

- L’un des droits qu’une action nomme suffit.
- Sans *Gérer les manifestations*, une manifestation s’ouvre sur ses
  classements.
- L’application ne montre à chaque membre que les pages et les boutons que
  son rôle permet.

## Rôles par défaut {#default-roles}

| Rôle | Droits |
|---|---|
| **Admin** | tous ; la colonne ne peut pas être modifiée et le rôle ne peut pas être supprimé |
| **Comité** | tous sauf *Stands et imprimantes* et *Utilisateurs et société* |
| **Aide** | voir et gérer les manifestations, statistiques, voir et modifier les tireurs et les sociétés |
| **Membre** | voir les manifestations et les classements |

- **Nouveau rôle** commence avec seulement *Voir les manifestations et les
  classements*.
- Un rôle peut être supprimé dès que personne ne l’a ; le rôle Admin jamais.
- Une société garde toujours au moins un membre avec le rôle Admin.
- Les membres avec *Utilisateurs et société* mais sans le rôle Admin ne
  donnent, modifient et retirent que des rôles dont ils ont eux-mêmes les
  droits.

## Liens d’aide {#helper-links}

Les liens d’aide ne sont pas des rôles de la société ; chaque lien a l’un de
ces rôles (voir [Créer un lien](../user-guide/helpers.md#create)) :

| Rôle | Voit | Peut |
|---|---|---|
| **Cibarre** | la grille de ses propres cibles, les tireurs qui s’y trouvent | saisir résultats et distinctions, **Terminer** un tireur |
| **Distribution de munition** | le champ de recherche, la liste d’attente, les tireurs sur toutes les cibles | marquer la munition comme remise (**Munition**), attribuer une cible à un tireur en attente, imprimer la feuille de stand |
| **Saisie des résultats** | le champ de recherche, les tireurs sur toutes les cibles, la liste des terminés | saisir résultats et distinctions sur toutes les cibles |

Les cibarres et le bureau des résultats peuvent aussi imprimer les [modèles
d’impression](../user-guide/print-forms.md#use) de leurs cibles ; le bureau
des résultats peut marquer les distinctions comme remises.

## Portées des jetons API {#scopes}

Un jeton (**Compte** → **Jetons API**, voir [Jetons
API](../user-guide/login-language.md#api-tokens)) a une liste de sociétés et
de portées :

| Portée | Clé | Permet |
|---|---|---|
| *Lire mon journal d’entraînement* | `training:read` | lire son propre journal d’entraînement |
| *Saisir et modifier mes séances d’entraînement* | `training:write` | ajouter, modifier et supprimer ses propres séances, disciplines, armes et objectifs |
| un droit ci-dessus | `events.view`, `events.operate`, `events.configure`, `reports.view`, `shooters.view`, `shooters.edit`, `shooters.manage`, `ranges.manage`, `forms.manage`, `club.manage` | la même chose que dans un rôle, dans les sociétés du jeton |

- Un droit de société ne peut être choisi que si le rôle du propriétaire l’a
  dans l’une des sociétés du jeton.
- Un jeton n’a jamais plus que le rôle du propriétaire : un droit que le
  rôle perd, ou une société que le propriétaire quitte, disparaît aussitôt
  pour le jeton.
- Un jeton sans société n’atteint que le journal d’entraînement.

### Modèles prédéfinis {#presets}

| Modèle | Portées |
|---|---|
| **Classements pour un site web** | *Voir les manifestations et les classements* |
| **App d’entraînement** | *Lire mon journal d’entraînement*, *Saisir et modifier mes séances d’entraînement* |
| **Lecture seule** | *Voir les manifestations et les classements*, *Statistiques et rapports*, *Voir les tireurs et les sociétés*, *Lire mon journal d’entraînement* |

Un modèle ne remplit que les droits que les sociétés choisies permettent.

### Utiliser un jeton {#token-use}

| Quoi | Valeur |
|---|---|
| Format | `sb_` et 40 chiffres hexadécimaux, affiché une fois |
| En-tête | `Authorization: Bearer sb_…` ou `X-API-Key: sb_…` |
| Expiration | dans 30 ou 90 jours, dans un an, ou jamais |
| Points d’accès | `/api/tenants/…` (données de la société, selon les droits), `/api/my/…` (journal d’entraînement, selon la portée), `GET /api/token`, `/api/version`, `/api/auth/options` ; avec *Configurer les manifestations* aussi `/api/modules` et `/api/programme-templates` ; tout le reste répond 403 |
| `GET /api/token` | le nom, les portées et les sociétés du jeton |
| Limite de requêtes | par défaut 600 requêtes par jeton en cinq minutes, puis 429 ; l’exploitant la modifie avec `STANDBLATT_API_RATE_LIMIT` (0 = pas de limite) |
| Jetons inconnus | après 20 jetons inconnus ou révoqués depuis une adresse en 15 minutes, 429 |
| Description de l’API | `/api/docs` |
| Jamais accessible | mot de passe, e-mail, réglages du compte, autres jetons |
