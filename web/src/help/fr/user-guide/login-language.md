# Connexion et langue

## Se connecter

Connecte-toi avec le nom d’utilisateur (ou l’adresse e-mail) et le mot de
passe reçus de l’admin de ta société ou de l’exploitant de la
plateforme. Les comptes de société ne se créent pas par inscription ; seul
[Mon entraînement](../my-training/account.md) a sa propre inscription. Mot
de passe oublié ? Si ton compte a une adresse e-mail confirmée, utilise
**Mot de passe oublié ?** sur la page de connexion ; sinon demande à
l’exploitant de définir un nouveau mot de passe dans le [panneau
d’administration](../admin-guide/admin-panel.md). Connecté, tu le changes
toi-même sous [Compte](#account).

Après la connexion, tu arrives sur la page d’accueil avec une carte pour
chaque société dont tu fais partie et une pour [Mon
entraînement](../my-training/index.md), que chaque utilisateur possède ;
**Compte** et **Se déconnecter** se trouvent sous les cartes. Sans société,
la page d’accueil mène directement à Mon entraînement.

![Login](../images/screenshots/login.png){ .screenshot }

![Page d’accueil](../images/screenshots/home.png){ .screenshot }

Une connexion reste valable sept jours sans activité et 30 jours au plus ;
avec **Rester connecté 90 jours** coché, elle dure 90 jours sur cet
appareil. Après trop d’échecs, les connexions depuis la même adresse sont
bloquées un moment.

Si ton compte appartient à plusieurs sociétés, tu choisis la société après
la connexion ; **Changer de société** dans le menu utilisateur ramène à ce
choix.

## Compte {#account}

**Compte** dans le menu utilisateur (bas de la barre latérale), sur la page d’accueil et dans la barre latérale de Mon entraînement ouvre les réglages de ton compte :

- **Connexion** : ton nom d’utilisateur (il ne peut pas être changé) et ton
  adresse e-mail. **Ajouter** ou **Changer** demande ton mot de passe et
  envoie un lien de confirmation ; tant que tu ne confirmes pas, l’adresse
  compte comme *non confirmée* et **Renvoyer** envoie un nouveau lien. Une
  adresse confirmée te permet de réinitialiser toi-même un mot de passe
  oublié.
- **Mot de passe** : saisis le mot de passe actuel et le nouveau (au moins
  10 caractères). Les autres appareils sont déconnectés, celui que tu
  utilises reste connecté.
- **Langue et apparence** : la langue de l’interface de ton compte et le
  mode clair ou sombre.
- **Deuxième facteur** : voir ci-dessous.
- **Supprimer le compte** : seulement pour les comptes sans société ; le
  journal d’entraînement, les objectifs et les réglages sont supprimés
  définitivement. Les comptes de société sont supprimés par l’admin de la
  société.

![Compte](../images/screenshots/account.png){ .screenshot }

## Deuxième facteur {#totp}

Une app d’authentification (Google Authenticator, Aegis, 1Password et
autres) peut protéger ta connexion : après le mot de passe, Standblatt
demande le code à six chiffres affiché par l’app. Configure-la sous
[Compte](#account) : **Configurer**, scanner le code QR, saisir le code
actuel, **Activer**. **Désactiver** demande ton mot de passe. Téléphone
perdu ? L’exploitant peut réinitialiser le deuxième facteur dans le [panneau
d’administration](../admin-guide/admin-panel.md). Il peut aussi désactiver
la fonction pour tout le serveur.

## Langue

L’interface est disponible en anglais, allemand, français et italien. Sans
choix, Standblatt suit la langue du navigateur. **Langue** dans le menu
utilisateur (en bas de la barre latérale) la définit pour ton compte sur
tous les appareils.

Les documents imprimés, les exports CSV et les écrans kiosque ne suivent pas
ton choix mais la langue de la manifestation, pour que chaque impression
d’une manifestation soit identique (voir
[Manifestations](events.md#settings)).

## Version

Le bas de la barre latérale indique la version du serveur, ce qui aide à
signaler un problème ; dans Mon entraînement sur un téléphone, elle se
trouve sous **Réglages**. Quand une nouvelle version sort, les pages
ouvertes proposent **Recharger** dans la minute et la chargent au plus tard
au prochain changement de page.
