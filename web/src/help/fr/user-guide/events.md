# Manifestations

**Manifestations** liste les manifestations de la société dans les groupes *Aujourd’hui* (entre le premier et le dernier jour de tir), *À venir*, *Passées* et *Clôturées* ; les admins en créent de nouvelles ici. Ouvrir une manifestation affiche dans la barre latérale son tableau de bord, les réservations, les classements, les statistiques et la configuration ; *Changer de manifestation* en haut de la barre latérale ramène à la liste.

![Manifestations](../images/screenshots/events.png){ .screenshot }

## Réglages {#settings}

La **Configuration** de la manifestation a cinq onglets : *Programmes* (voir
[Programmes](programmes.md)), *Manifestation*, *Écrans kiosque* (voir
[Classements et kiosque](rankings-kiosk.md)), *Auxiliaires* (voir
[Auxiliaires](helpers.md)) et *Copier & clôturer*. Les collaborateurs ne
voient que les programmes.

L’onglet *Manifestation* contient les réglages (admins) :

| Réglage | Signification |
|---|---|
| Nom, date, stands | la date est le jour principal ; son année fixe les catégories d’âge (FST : année du concours moins année de naissance, le jour lui-même ne compte pas, ce qui importe pour le tir anticipé de décembre avant une manifestation de janvier) ; cocher chaque stand sur lequel la manifestation se tire |
| Jours de tir | jours isolés avec heures de tir facultatives, pour le tir anticipé et les manifestations sur plusieurs jours (les inscriptions comptent le jour où elles sont faites) ; le jour principal en fait toujours partie |
| Période (saison) | pour les manifestations sur des semaines ou des mois, p. ex. le Tir en campagne de l’ouverture du stand jusqu’au 31 août : premier et dernier jour, plus les jours de la semaine avec leurs heures de tir ; aucun jour coché signifie chaque jour. Les jours isolés, p. ex. le week-end officiel, s’ajoutent à la période ou remplacent ses heures à cette date |
| Organisateur | imprimé sur les documents |
| Langue des documents et du kiosque | langue des PDF, des exports CSV et des écrans kiosque, indépendamment de la langue d’interface de chaque utilisateur |
| Société exclue du concours de sections | en général la société organisatrice ; listée *hors concours* |
| Afficher l’état des cibles sur les classements | les écrans kiosque montrent les cibles |
| Imprimer la valeur de la distinction sur les cartes-couronnes | valeur en chiffres et en lettres sur la carte |
| L’organisateur remet la munition | affiche le bouton **Munition** sur le tableau de bord et pour la distribution de munition (voir [Aides](helpers.md)) ; désactivé pour les manifestations où les tireurs apportent la leur, p. ex. pistolet et carabine à air |

Les heures de tir donnent aussi ses créneaux à la réservation en ligne, sauf
si la réservation fixe d’autres horaires pour un jour (voir [réservation en
ligne](online-booking.md)).

## Copier pour l’année prochaine

Dans l’onglet *Copier & clôturer*, **Copier pour l’année prochaine** crée
une nouvelle manifestation avec les mêmes stands, organisateur, langue,
réglages et tous les programmes, y compris prix, distinctions, combinaisons
et prix spéciaux ; les jours de tir et la période suivent la date, les
heures de tir hebdomadaires restent sur leurs jours de la semaine. Les
inscriptions, les résultats et les liens kiosque ne sont pas copiés. Le
dialogue propose le nom avec l’année (et un numéro en tête comme *42.*)
incrémentée et la date un an plus tard ; adaptez les deux au besoin.

![Copier une manifestation](../images/screenshots/copy-event.png){
.screenshot }

## Clôturer {#close}

Quand la manifestation est terminée et les classements vérifiés,
**clôturez**-la dans l’onglet *Copier & clôturer* (**Clôturer la
manifestation**). La clôture publie les classements : les tireurs disposant
d’un [journal d’entraînement
personnel](../my-training/account.md#competitions) trouvent alors leurs
résultats de cette manifestation, avec leur rang dans le programme, et
peuvent les ajouter à leur journal. Tout reste visible (tableau de bord,
résultats, classements, rapports, documents), mais les inscriptions, les
résultats, les distinctions, les programmes et les réglages ne peuvent plus
être modifiés. Une bannière sur chaque page de la manifestation
l’indique. **Rouvrir** la rend à nouveau modifiable ; les résultats déjà
ajoutés par les tireurs à leur journal y restent. Copier une manifestation
clôturée pour l’année prochaine fonctionne.

## Supprimer

**Supprimer** (dans l’onglet *Copier & clôturer*) retire une manifestation, mais seulement tant que personne n’y est inscrit.
