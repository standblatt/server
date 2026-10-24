# Afficher les classements sur un écran dans la halle

Tu en as besoin quand les tireurs et les spectateurs doivent suivre les
classements en direct sur un écran au stand. L’écran a besoin d’un
navigateur et d’Internet, sans connexion à un compte. La mise en place
demande le droit *Configurer les manifestations* (voir
[Rôles](../user-guide/users.md#roles)).

1. Sous **Configuration** → **Modules** de la manifestation, active
   *Affichages kiosque*.
2. Ouvre **Affichages kiosque** dans la configuration.
3. Saisis un nom pour l’écran (p. ex. *Halle 1*), choisis **Sombre** ou
   **Clair** et, si l’écran ne doit pas utiliser la langue de la
   manifestation, une langue. Pour une manifestation sur plusieurs stands,
   choisis le stand que montre l’écran, ou *Tous les stands*.
4. Clique sur **Créer un lien kiosque** et **copier** le lien. Il n’est
   affiché qu’une fois.
5. Ouvre le lien dans un navigateur sur l’écran et clique sur le bouton
   **Plein écran** en haut à droite.

![Écran kiosque](../images/screenshots/kiosk.png){ .screenshot }

L’écran se met à jour en direct et se reconnecte tout seul après une coupure
du réseau.

## Choisir ce que montre l’écran {#content}

- Un programme n’apparaît que tant que **Afficher sur les écrans kiosque**
  est activé sur sa page de programme (c’est le cas pour les nouveaux
  programmes). Les pages du personnel, le PDF et le CSV le classent quand
  même.
- Un classement général apparaît avec **Afficher sur les écrans kiosque**
  dans son éditeur (voir [Classements
  généraux](../user-guide/classifications.md)).
- Avec l’attribution des cibles, **Afficher l’état des cibles sur les
  classements** sous **Affichages kiosque** ajoute l’état des cibles
  au-dessus des classements.

## Arrêter un écran {#revoke}

Clique sur **révoquer** à côté de l’écran dans la liste sous **Affichages
kiosque** : le lien cesse aussitôt de fonctionner. *vu la dernière fois*
indique quand un écran a récupéré des données pour la dernière fois. Les
parties du lien sont listées dans [Liens sans
connexion](../reference/links.md).
