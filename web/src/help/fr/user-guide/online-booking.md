# Réservation en ligne

Avec la **réservation en ligne**, les sociétés externes demandent
elles-mêmes des cibles. L’organisateur met un lien dans l’invitation ; une
société choisit les cibles souhaitées dans chaque créneau et annonce ses
tireurs. Un admin confirme ou refuse chaque demande.

## Activer {#setup}

La réservation en ligne fait partie du module *Réservations* (il nécessite
un stand et *Attribution des cibles*) : sur **Réservations** dans la
configuration, **Activer** suffit. La page montre ensuite le lien avec un
code QR pour l’invitation, si les sociétés peuvent faire une demande en ce
moment (*Ouvre le*, *Ouvert jusqu’au* ou *Délai passé*) et les créneaux
proposés, chaque jour avec le début de chaque créneau et ses pauses.

![Réservation en ligne dans la
configuration](../images/screenshots/booking-setup.png){ .screenshot }

Sans modification, chaque jour de la manifestation est ouvert pendant ses
heures de tir (voir [Manifestations](events.md#settings)), sinon de 08:00 à
17:00, en créneaux de 30 minutes, et les demandes sont acceptées tout de
suite jusqu’à 14 jours avant le premier jour. Un créneau peut être demandé
jusqu’à deux heures avant son début. **Adapter** change cela :

- **Tir de** et **Tir jusqu’à** : les horaires de chaque jour sans heures de
  tir.
- **Pause de** et **Pause jusqu’à** : une pause facultative, p. ex. à midi.
- **Durée d’un créneau (minutes)** : la durée d’un créneau.
- **Cibles par créneau et société** : une limite par demande ; vide signifie
  sans limite.
- **Ouverture des inscriptions** : le premier jour où les sociétés peuvent
  demander ; jusque-là, le lien indique seulement à partir de quand.
- **Délai d’inscription** : le dernier jour où les sociétés peuvent demander
  ou modifier leur demande, au plus tard le dernier jour de tir.
- **Jours avec d’autres horaires** : une ligne par jour de la manifestation,
  p. ex. vendredi 17:00–20:00. Les jours viennent de la manifestation et la
  suivent quand tu en ajoutes ou en retires un.

**Nouveau lien** remplace le lien ; l’ancien ne fonctionne plus. **Désactiver** ferme la réservation ; les demandes et les réservations confirmées restent. *Copier pour l’année prochaine* reprend les horaires, mais ni le lien, ni l’ouverture des inscriptions, ni le délai.

## Ce que voit une société {#club}

Le lien ouvre une page sans connexion. Elle explique le déroulement en trois
étapes et montre chaque créneau qui commence dans plus de deux heures avec
les cibles de chaque stand sous forme de boutons numérotés : vert est libre,
jaune est libre mais demandé par une autre société dont la demande n’est pas
encore confirmée, rouge pâle est occupé, et le choix de la société est foncé
avec un cadre vert. Avant le jour d’ouverture, elle indique seulement à
partir de quand la réservation est ouverte. Seules les réservations
confirmées occupent une cible ; les noms des autres sociétés ne sont jamais
affichés. Avec plus de deux jours, chaque jour est replié et indique combien
de créneaux la société y a choisis.

![La page de réservation d’une
société](../images/screenshots/booking-public.png){ .screenshot }

La société touche les cibles souhaitées dans chaque créneau (une deuxième
touche en retire une ; avec une limite par créneau, les autres boutons se
bloquent dès qu’elle est atteinte), puis saisit son nom, une personne de
contact, l’e-mail et le téléphone, ainsi que ses **Tireurs** : nom, prénom,
date de naissance et éventuellement la licence. **Coller une liste** reprend
des lignes copiées depuis Excel (nom, prénom, date de naissance,
licence). La liste des tireurs peut suivre plus tard ; les programmes sont
choisis au stand comme d’habitude.

Une barre en bas de l’écran compte les créneaux et les cibles choisis et
contient **Envoyer la demande**. Ensuite, la société reçoit un e-mail avec
un lien personnel. Elle peut y voir le statut, modifier ou annuler la
demande jusqu’au délai. Modifier les créneaux d’une demande confirmée libère
ses cibles et la demande attend une nouvelle confirmation ; ajouter des
tireurs la laisse confirmée.

## Confirmer les demandes {#confirm}

Chaque admin avec un e-mail reçoit un mail pour chaque demande nouvelle,
modifiée ou annulée. Les demandes figurent en haut de **Réservations**, les
ouvertes d’abord et, dans chaque statut, dans l’ordre d’arrivée ; la barre
latérale indique par un compteur rouge combien attendent. Une demande
ouverte indique les cibles souhaitées par la société. Chaque demande indique
à la seconde près quand elle a été reçue, une demande ouverte aussi quand la
société l’a modifiée pour la dernière fois ; ainsi, quand deux sociétés
veulent les mêmes cibles, la première peut les obtenir.

![Demandes de réservation sur la page des
réservations](../images/screenshots/booking-requests.png){ .screenshot }

- **Confirmer** attribue à chaque créneau les cibles choisies par la
  société, comme réservation ordinaire au nom de la société qui apparaît
  comme réservée un créneau avant son début, et envoie les numéros de cibles
  à la société par mail. Une cible choisie occupée entre-temps est remplacée
  par une autre cible libre du stand, voisine si possible. Si un créneau a
  trop peu de cibles libres, rien n’est réservé et le message nomme le
  créneau ; déplacez une réservation ou refusez.
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
