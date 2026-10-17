# Compte

Chaque utilisateur de Standblatt a Mon entraînement : les utilisateurs d’une
société l’ouvrent depuis la page d’accueil ou le menu utilisateur, les
tireurs sans société s’inscrivent eux-mêmes. Adresse e-mail, mot de passe,
langue, deuxième facteur et suppression du compte se trouvent sur la page
[Compte](../user-guide/login-language.md#account), qu’ouvrent **Compte**
dans la barre latérale ou **Réglages → Compte**.

## S’inscrire

Ouvre **Créer un compte** sur la page de connexion. Il te faut ton nom, ton
adresse e-mail et un mot de passe d’au moins 10 caractères. Tu es connecté
tout de suite.

![Inscription](../images/screenshots/my-signup.png){ .screenshot .phone }

Standblatt t’envoie un lien pour confirmer ton adresse e-mail. Tant que tu
ne l’as pas confirmée, une note sur la page Aujourd’hui te le rappelle ;
**Renvoyer** envoie un nouveau lien. Le lien est valable sept jours.

Mot de passe oublié ? Utilise **Mot de passe oublié ?** sur la page de
connexion. Tu reçois un lien valable une heure. Définir un nouveau mot de
passe te déconnecte partout. Connecté, tu changes le mot de passe sous
[Compte](../user-guide/login-language.md#account).

Les comptes créés par un admin de société se connectent avec leur nom
d’utilisateur au lieu d’une adresse e-mail ; ils peuvent ajouter une adresse
e-mail sous **Réglages → E-mail**.

## Profil et résultats de compétition {#competitions}

Sous **Réglages → Profil**, tu peux saisir ta **date de naissance** et ton
**numéro de licence**. Avec les deux, Standblatt trouve tes résultats des
manifestations organisées avec Standblatt par n’importe quelle société, par
exemple le championnat de société. Ils apparaissent dans le **Journal** sous
*Tes résultats de compétition* dès que la société a clôturé la
manifestation, c’est-à-dire dès que le classement est définitif. Chaque
résultat porte ton **rang** dans le programme (par exemple *Rang 12 sur
87*). Choisis la discipline et touche **Ajouter** ou **Tout ajouter**. Les
résultats ajoutés deviennent des séances de compétition avec leurs séries et
gardent le rang.

## Disciplines personnalisées

Six disciplines sont intégrées : pistolet et carabine à air comprimé 10 m,
pistolet 25 m et 50 m, fusil 50 m et 300 m. Pour tout le reste, p. ex. un
tir en campagne, ajoute une **Discipline personnalisée** dans les réglages
avec :

- le nombre habituel de passes et de coups par passe,
- un comptage entier ou au dixième (jusqu’à 10.9) et la valeur maximale d’un
  coup,
- en option une cible pour toucher les impacts.

## Armes {#weapons}

Tu tires avec plus d’une arme, p. ex. deux pistolets à air ? Ajoute-les sous
**Réglages → Armes** avec un nom de ton choix ; le calibre ne joue aucun
rôle. Choisis l’arme en enregistrant une séance et la
[Progression](progress.md#weapons-compared) les compare. Renommer une arme
garde ses séances. La supprimer aussi, elles n’ont simplement plus d’arme.

## Jetons API

Tu veux importer des séances d’une autre application ou analyser ton journal
avec ton propre script ? Crée un jeton sous **Réglages → Jetons API** :

- un **nom** qui te dira plus tard qui l’utilise,
- les **autorisations** : *Lire les séances* et/ou *Enregistrer et modifier
  des séances*,
- quand il **expire** : dans 30 ou 90 jours, dans un an ou jamais.

Le jeton (`sb_…`) n’est affiché qu’une seule fois ; copie-le tout de
suite. Les scripts l’envoient comme `Authorization: Bearer sb_…` (ou
`X-API-Key: tl_…`). Il vaut pour tes séances, disciplines, armes, objectifs,
statistiques et l’export, mais jamais pour ton mot de passe, ton e-mail, ton
compte ou d’autres jetons, ni pour les données de société. La liste montre
quand chaque jeton a été utilisé pour la dernière fois ; le révoquer
l’arrête immédiatement. Les points d’accès figurent dans la [référence de
l’API](../developer/api.md).

## Installer l’application

Standblatt est une application web qui peut être installée :

- **Android / Chrome** : Réglages → **Installer**, ou *Ajouter à l’écran
  d’accueil* dans le menu du navigateur.
- **iPhone** : Partager → *Sur l’écran d’accueil*.

Installée, elle démarre en plein écran. Un appui long sur l’icône propose
*Démarrer l’entraînement*.

## Tes données

- **Exporter toutes les séances (CSV)** télécharge chaque passe de chaque
  séance.
- **Supprimer le compte** supprime le compte avec toutes les séances,
  objectifs et réglages après la saisie de ton mot de passe. Les comptes qui
  appartiennent à une société doivent d’abord être retirés de la société par
  son admin.
