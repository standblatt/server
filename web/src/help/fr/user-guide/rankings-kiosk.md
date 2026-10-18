# Classements et kiosque

## Classements

**Classements** affiche un classement en direct par programme, ordonné selon les [règles de la FST](../reference/rules.md#ranking). Les trois premiers rangs sont mis en évidence.

![Classements](../images/screenshots/rankings.png){ .screenshot }

- **PDF** ouvre le classement (A4) pour l’imprimer ou le partager, **CSV**
  le télécharge pour Excel. Les classements ne sont pas envoyés à
  l’imprimante du stand.
- Les programmes désactivés pour les écrans kiosque sont quand même classés
  ici et marqués *pas sur le kiosque*.
- Les [classements généraux](classifications.md) suivent les programmes, un
  tableau par classe ; leur **PDF** et **CSV** sont le rapport *Classements
  généraux*.

## Écrans kiosque

Un lien kiosque affiche les classements d’une manifestation sur un écran de
la halle, sans connexion. Active le module *Affichages kiosque* sous
**Configuration** → **Modules**, puis crée-en un sur **Affichages kiosque**
dans la configuration :

1. Saisis un nom pour l’écran (p. ex. *Halle 1*), choisis **Sombre** ou
   **Clair** et, si besoin, une langue.
2. **Créer un lien kiosque** ; copie le lien. Il n’est affiché qu’une fois.
3. Ouvre le lien dans un navigateur sur l’écran, idéalement en plein écran.

![Écran kiosque](../images/screenshots/kiosk.png){ .screenshot }

- Le kiosque n’affiche que les programmes avec *Afficher sur les écrans
  kiosque* activé, et l’état des cibles si *Afficher l’état des cibles sur
  les classements* est activé (sur la même page).
- La langue est celle de la manifestation, sauf si le lien en indique une
  autre (`?language=de`, `fr`, `it`, `en`) ; le thème vient de `?theme=dark`
  ou `?theme=light`. Pour une manifestation sur plusieurs stands,
  `?range=<id>` (choisi à la création du lien) limite l’affichage à ce stand
  : ses cibles, ses programmes et les programmes sans stand.
- L’écran se met à jour en direct et se reconnecte tout seul après une
  coupure du réseau. À une nouvelle version de Standblatt, il se recharge
  tout seul.
- **révoquer** dans la liste arrête un lien immédiatement. *vu la dernière
  fois* indique quand un écran a récupéré des données pour la dernière fois.
