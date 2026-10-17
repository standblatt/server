# Réservation en ligne

Avec la **réservation en ligne**, les sociétés externes demandent
elles-mêmes des cibles. L’organisateur met un lien dans l’invitation ; une
société choisit les créneaux voulus, indique combien de cibles il lui faut
dans chacun et annonce ses tireurs. Un admin confirme ou refuse chaque
demande.

## Activer {#setup}

Dans l’onglet *Réservation en ligne* de la **Configuration** de la
manifestation, **Activer** suffit. L’onglet montre ensuite le lien avec un
code QR pour l’invitation et les créneaux proposés.

![Réservation en ligne dans la
configuration](../images/screenshots/booking-setup.png){ .screenshot }

Sans modification, chaque jour de la manifestation est ouvert pendant ses
heures de tir (voir [Manifestations](events.md#settings)), sinon de 08:00 à
17:00, en créneaux de 30 minutes, et les demandes sont acceptées jusqu’à 14
jours avant le premier jour. **Adapter** change cela :

- **Tir de** et **Tir jusqu’à** : les horaires de chaque jour sans heures de
  tir.
- **Pause de** et **Pause jusqu’à** : une pause facultative, p. ex. à midi.
- **Durée d’un créneau (minutes)** : la durée d’un créneau.
- **Cibles par créneau et société** : une limite par demande ; vide signifie
  sans limite.
- **Délai d’inscription** : le dernier jour où les sociétés peuvent demander
  ou modifier leur demande.
- **Jours avec d’autres horaires** : une ligne par jour de la manifestation,
  p. ex. vendredi 17:00–20:00. Les jours viennent de la manifestation et la
  suivent quand tu en ajoutes ou en retires un.

**Nouveau lien** remplace le lien ; l’ancien ne fonctionne plus. **Désactiver** ferme la réservation ; les demandes et les réservations confirmées restent. *Copier pour l’année prochaine* reprend les horaires, mais ni le lien ni le délai.

## Ce que voit une société {#club}

Le lien ouvre une page sans connexion. Elle montre la manifestation, le
délai et chaque créneau avec les cibles encore libres par stand. Seules les
réservations confirmées comptent ; les noms des autres sociétés ne sont
jamais affichés.

![La page de réservation d’une
société](../images/screenshots/booking-public.png){ .screenshot }

La société saisit le nombre de cibles par créneau, son nom, une personne de
contact, l’e-mail et le téléphone, ainsi que ses **Tireurs** : nom, prénom,
date de naissance et facultativement la licence. **Coller une liste**
reprend des lignes copiées depuis Excel (nom, prénom, date de naissance,
licence). La liste des tireurs peut suivre plus tard ; les programmes se
choisissent au stand comme d’habitude.

Après **Envoyer la demande**, la société reçoit un e-mail avec un lien
personnel. Elle peut y voir le statut, modifier ou annuler la demande
jusqu’au délai. Modifier les créneaux d’une demande confirmée libère ses
cibles et la demande attend une nouvelle confirmation ; ajouter des tireurs
la laisse confirmée.

## Confirmer les demandes {#confirm}

Chaque admin avec un e-mail reçoit un message pour chaque demande nouvelle,
modifiée ou annulée. Les demandes figurent en haut de **Réservations** ; la
barre latérale indique combien attendent.

![Demandes de réservation sur la page des
réservations](../images/screenshots/booking-requests.png){ .screenshot }

- **Confirmer** attribue à chaque créneau des cibles libres de son stand,
  voisines si possible, comme réservation ordinaire au nom de la société, et
  envoie les numéros de cible à la société. Si un créneau a trop peu de
  cibles libres, rien n’est réservé et le message nomme le créneau ; déplace
  une réservation ou refuse.
- **Confirmer** demande d’abord à quelle **Société** des données de base la
  demande appartient, car la société a saisi son nom elle-même : la société
  du même nom est déjà choisie, les sociétés au nom semblable sont proposées
  et *Nouvelle société* la crée. La demande affiche ensuite la société
  attribuée.
- **Confirmer** reprend aussi les tireurs dans les données de base, dans
  cette société : un tireur est trouvé par sa licence, sinon par nom et date
  de naissance, et créé autrement. Les tireurs existants ne sont jamais
  écrasés, car la liste vient d’un formulaire public ; seule une date de
  naissance ou une licence manquante est complétée. Les tireurs ajoutés
  après la confirmation sont repris avec **Reprendre des tireurs**.
- **Refuser** demande un motif facultatif, que la société reçoit par
  e-mail. Les cibles que la demande occupait sont de nouveau libres.

Les cibles réservées sont des [réservations](reservations.md) ordinaires :
modifie-y les numéros de cible au besoin. Pendant la manifestation,
l’inscription trouve les tireurs avec la recherche habituelle. Le rapport
*Réservations de cibles* sous [Statistiques](statistics-reports.md) liste
chaque société avec ses créneaux, cibles et tireurs.
