# Aides

Les aides de la manifestation (cibarres, distribution de munition, bureau
des résultats) travaillent avec un **lien d’aide** au lieu d’un login. Un
lien d’aide ouvre sur n’importe quel téléphone ou ordinateur portable une
vue réduite d’une manifestation ; les admins activent le module *Liens
d'aide* sous **Configuration** → **Modules**, créent les liens sur **Liens
d'aide** dans la configuration et les révoquent quand elle est terminée.

![Page d’aide du cibarre sur un
téléphone](../images/screenshots/helper-scorer.png){ .screenshot .phone }

## Rôles {#roles}

| Rôle | Voit | Peut |
|---|---|---|
| **Cibarre** | la grille de ses propres cibles, les tireurs qui s’y trouvent | saisir résultats et distinctions, **Terminer** un tireur |
| **Distribution de munition** | le champ de recherche, la liste d’attente, les tireurs sur toutes les cibles | marquer la munition comme remise (**Munition**), attribuer une cible à un tireur en attente, imprimer la feuille de stand |
| **Saisie des résultats** | le champ de recherche, les tireurs sur toutes les cibles, la liste des terminés | saisir résultats et distinctions sur toutes les cibles |

Tout le reste reste aux utilisateurs de la société : inscrire des tireurs,
modifier des programmes, classements, statistiques, supprimer.

## Créer un lien {#create}

1. Saisis un nom (qui reçoit le lien, p. ex. *Cibarre stand A*), choisis le
   rôle et, pour un cibarre, le stand et les cibles (`1-5, 8`). Un cibarre
   s’occupe généralement de plusieurs cibles. Les cibarres n’existent que
   tant que la manifestation attribue des cibles (module *Attribution des
   cibles*).
2. **Créer un lien d’aide** ; copie le lien et envoie-le à l’aide
   (p. ex. par messagerie). Il n’est affiché qu’une fois.
3. L’aide ouvre le lien ; il reste valable jusqu’à sa révocation ou la
   suppression de la manifestation. `?language=de` (`fr`, `it`, `en`)
   remplace la langue de la manifestation.

**révoquer** dans la liste arrête un lien immédiatement ; *vu la dernière fois* indique quand il a été utilisé pour la dernière fois. Un cibarre qui ouvre un tireur sur une autre cible obtient *Cette cible n’est pas la tienne*.

## Munition {#ammunition}

**Munition** sur une inscription en attente ou en tir note l’heure à laquelle la munition a été remise ; le bouton n’existe que si le module *Distribution de munition* est activé (les tireurs au pistolet et à la carabine à air apportent généralement la leur). Le tableau de bord montre un badge vert *munition* sur ces inscriptions, ainsi la distribution et le rangeur voient d’un coup d’œil qui a été servi.
