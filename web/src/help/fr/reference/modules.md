# Modules

Les modules sont les éléments qu’une manifestation ou une société active. Ce
qui est désactivé disparaît de la barre latérale, du tableau de bord, de
l’éditeur de programmes, des classements et des rapports. Comment les
activer : [Modules](../user-guide/events.md#modules) et [Modules de la
société](../user-guide/events.md#club-modules).

## Modules de la manifestation {#event-modules}

Sous **Configuration** → **Modules** d’une manifestation.

| Module | Ce qu’il ajoute |
|---|---|
| Attribution des cibles | les tireurs sont inscrits sur des cibles : la grille des cibles sur le tableau de bord, l’état des cibles sur les affichages kiosque et les liens de cibarre (nécessite un stand ; voir plus bas) |
| Réservations | cibles réservées pour des groupes et des sociétés ; sa page dans la configuration active la [réservation en ligne](../user-guide/online-booking.md), par laquelle les sociétés demandent des cibles avec un lien public (nécessite un stand) |
| Feuille de stand | la feuille de stand imprimée à l’inscription : tireur, cible, programmes et codes-barres |
| Bande de contrôle | la bande de contrôle imprimée à l’inscription : les passes que le tireur peut tirer |
| Distribution de munition | le bouton **Munition** sur le tableau de bord et pour la distribution de munition (voir [Aides](../user-guide/helpers.md#ammunition)) ; laisse-le désactivé là où les tireurs apportent la leur, p. ex. pistolet et carabine à air |
| Tentatives | des programmes qui peuvent être tirés plusieurs fois (voir [Programmes](../user-guide/programmes.md#attempts)) |
| Distinctions | cartes-couronnes, prix en nature et prix en espèces selon les points par catégorie d’âge |
| Prix spéciaux | les meilleurs rangs de chaque catégorie d’âge (active Distinctions, voir [Programmes](../user-guide/programmes.md#special-prizes)) |
| Concours de sections | les sociétés classées selon les résultats de leurs tireurs |
| Classements généraux | classements sur plusieurs programmes ou tentatives (voir [Classements généraux](../user-guide/classifications.md)) |
| Affichages kiosque | classements sur des écrans dans la halle (voir [Classements et kiosque](../user-guide/rankings-kiosk.md)) |
| Liens d'aide | cibarres, distribution de munition et saisie des résultats sans compte (voir [Aides](../user-guide/helpers.md)) |

### Où apparaît un module {#places}

| Module | Pages et sections |
|---|---|
| Attribution des cibles | grille des cibles sur le tableau de bord, cibles dans l’inscription, état des cibles sur les affichages kiosque, rôle *Cibarre* des liens d’aide |
| Réservations | page **Réservations** de la manifestation, avec les demandes de réservation et leur nombre dans la barre latérale ; page **Réservations** dans la configuration (réservation en ligne) ; rapport *Réservations de cibles* |
| Feuille de stand | *Imprimer la feuille de stand* dans l’inscription |
| Bande de contrôle | *Imprimer la bande de contrôle* dans l’inscription |
| Distribution de munition | bouton **Munition** et badge *munition* sur le tableau de bord et sur les liens d’aide |
| Tentatives | section **Tentatives** de la page du programme ; compteur par programme dans l’inscription |
| Distinctions | section **Distinctions** de la page du programme ; distinctions dans les résultats ; chiffres sous **Statistiques** ; rapports *Distinctions* et *Distinctions à remettre* |
| Prix spéciaux | section **Prix spéciaux** de la page du programme ; rapport *Prix spéciaux* |
| Concours de sections | *Compte pour le concours de sections* sur la page du programme ; concours de sections sous **Statistiques** ; rapport *Concours de sections* |
| Classements généraux | page **Classements généraux** dans la configuration ; tableaux sous **Classements** et sur les affichages kiosque ; rapport *Classements généraux* |
| Affichages kiosque | page **Affichages kiosque** dans la configuration ; *Afficher sur les écrans kiosque* sur la page du programme |
| Liens d'aide | page **Liens d'aide** dans la configuration |

### Règles {#rules}

| Module | Demande | Nécessite un stand | Réglages | Ne peut pas être désactivé tant que |
|---|---|---|---|---|
| Attribution des cibles | – | oui | – | une inscription a une cible |
| Réservations | Attribution des cibles | oui | la réservation en ligne sur sa page de configuration | le lien de réservation en ligne est actif |
| Feuille de stand | – | – | – | – |
| Bande de contrôle | – | – | *Modèle des bandes de contrôle* | – |
| Distribution de munition | – | – | – | – |
| Tentatives | – | – | – | un programme permet plus d’une tentative |
| Distinctions | – | – | *Imprimer la valeur de la distinction sur les cartes-couronnes* (activé pour une nouvelle manifestation), *Modèle des cartes-couronnes* | un programme a des paliers de distinctions |
| Prix spéciaux | Distinctions | – | – | un programme décerne des prix spéciaux |
| Concours de sections | – | – | *Société exclue du concours de sections* | – |
| Classements généraux | – | – | – | la manifestation a un classement général |
| Affichages kiosque | – | – | *Afficher l’état des cibles sur les classements* (avec Attribution des cibles) | la manifestation a un lien kiosque |
| Liens d'aide | – | – | – | la manifestation a un lien d’aide |

- Un module demandé par un autre reste activé tant que celui-ci est activé.
- Les modules qui nécessitent un stand sont désactivés dans une
  manifestation sans stand, quel que soit le choix (voir [Manifestations
  sans cibles](../user-guide/events.md#no-lanes)).
- Les réglages se trouvent sur la carte du module sous **Modules** ;
  *Afficher l’état des cibles sur les classements* se trouve sur la page
  **Affichages kiosque**.
- **Copier pour l’année prochaine** et les modèles reprennent les modules
  choisis et les modules utilisés.
- Une manifestation créée par l’API sans liste de modules reçoit tous les
  modules de manifestation sauf *Distribution de munition*.

## Modules de la société {#club-modules}

Sous **Modules de la société** dans la partie société de la barre latérale
(droit *Utilisateurs et société*).

| Module | Ce qu’il ajoute | Ne peut pas être désactivé tant que |
|---|---|---|
| Modèles d’impression | des modèles pour feuilles préimprimées et la section *Modèle d’impression* de chaque programme (voir [Modèles d’impression](../user-guide/print-forms.md)) ; page **Modèles d’impression** et droit *Modèles d’impression* | la société a un modèle d’impression ou une image |

Une nouvelle société commence avec *Modèles d’impression* activé. Désactivé,
les manifestations impriment la carte-couronne et la bande de contrôle
intégrées.

## Genres de manifestations {#kinds}

**Nouvelle manifestation** demande le genre ; il ne fait qu’activer des modules, tout peut être modifié plus tard.

| Genre | Modules |
|---|---|
| Tir de fête | Attribution des cibles, Réservations, Feuille de stand, Bande de contrôle, Distinctions, Prix spéciaux, Concours de sections, Affichages kiosque, Liens d'aide |
| Programme annuel ou entraînement | Tentatives, Classements généraux |
| Tir en campagne | Attribution des cibles, Distinctions, Concours de sections, Affichages kiosque, Liens d'aide |
| Championnat de société | Tentatives, Classements généraux, Affichages kiosque |
| Sélection personnalisée | aucun |

Un [modèle](../user-guide/events.md#templates) de la société apporte plutôt
ses propres modules.
