# Classements et kiosque

## Classements

**Classements** affiche un classement en direct par programme, ordonné selon les [règles de la FST](../reference/rules.md#ranking). Les trois premiers rangs sont mis en évidence.

Le lien entre classements de programme, classements généraux, concours de
sections et prix spéciaux : [Comment les classements
s’articulent](../explanation/rankings.md).

![Classements](../images/screenshots/rankings.png){ .screenshot }

- **PDF** ouvre le classement (A4) pour l’imprimer ou le partager, **CSV**
  le télécharge pour Excel. Les classements ne sont pas envoyés à
  l’imprimante du stand.
- Les programmes désactivés pour les écrans kiosque sont quand même classés
  ici et marqués *pas sur le kiosque*.
- Les [classements généraux](classifications.md) suivent les programmes, un
  tableau par classe ; leur **PDF** et **CSV** sont le rapport *Classements
  généraux*.
- **Diplômes** imprime des diplômes pour les meilleurs rangs sur un [modèle
  d’impression](print-forms.md#certificates) de la société (voir [Imprimer
  des diplômes pour les meilleurs rangs](../how-to/certificates.md)) ; le
  bouton apparaît dès que la société a un modèle qui imprime un rang.

## Écrans kiosque

Un lien kiosque affiche les classements d’une manifestation sur un écran de
la halle, sans connexion ; mets-en un en place avec [Afficher les
classements sur un écran dans la halle](../how-to/kiosk.md).

**Affichages kiosque** dans la configuration (module *Affichages kiosque*) liste les liens. Un nouveau lien demande un nom, **Sombre** ou **Clair**, une langue et, pour une manifestation sur plusieurs stands, un stand ; **Créer un lien kiosque** n’affiche le lien qu’une fois.

![Écran kiosque](../images/screenshots/kiosk.png){ .screenshot }

- Le kiosque n’affiche que les programmes et les classements généraux avec
  *Afficher sur les écrans kiosque* activé et, avec le module *Attribution
  des cibles*, l’état des cibles si *Afficher l’état des cibles sur les
  classements* est activé sous **Affichages kiosque**.
- Le lien fixe le thème, la langue et le stand (voir [Paramètres
  d’URL](../reference/links.md#parameters)).
- Le bouton **Plein écran** en haut à droite passe l’écran en plein écran.
- L’écran se met à jour en direct et se reconnecte tout seul après une
  coupure du réseau. À une nouvelle version de Standblatt, il se recharge
  tout seul.
- **révoquer** dans la liste arrête un lien immédiatement. *vu la dernière
  fois* indique quand un écran a récupéré des données pour la dernière fois.
