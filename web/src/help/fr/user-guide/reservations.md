# Réservations

**Réservations** bloque des cibles du stand de la manifestation pour une durée, par exemple pour un entraînement ou un groupe. C’est la page du module *Réservations*, qui nécessite un stand et *Attribution des cibles*.

![Réservations](../images/screenshots/reservations.png){ .screenshot }

**Nouvelle réservation** demande :

- **Description** et, sur plusieurs stands, le **Stand de tir**.
- **Cibles** : une liste comme `5` ou `1, 5, 7-10`.
- **De** et **À** : la période. Les réservations sur la même cible ne
  peuvent pas se chevaucher ; à la suite, c’est possible.
- **Afficher comme réservée (minutes avant)** : délai à partir duquel les
  cibles apparaissent réservées dans le tableau de bord et le kiosque.

La liste affiche un bloc par jour dans l’ordre chronologique. Elle marque
chaque réservation comme *prévue*, *actif* ou *terminée*, et une bande de
petites cases montre quelles cibles du stand elle occupe (à partir de la
largeur d’une tablette) ; **afficher les passées** ajoute celles qui sont
terminées. Son menu propose **Modifier**, **copier** et **Supprimer**.

Les réservations appartiennent au stand et apparaissent donc dans chaque
manifestation sur ce stand. Les cibles réservées peuvent quand même être
attribuées à un tireur après confirmation. Le serveur annonce lui-même le
début et la fin des réservations, les écrans changent donc sans recharger.

Les demandes ouvertes que les sociétés externes envoient par la [réservation
en ligne](online-booking.md) figurent dans la même liste, une ligne par
créneau marquée *demandée*, avec les cibles souhaitées en ambre. Une fois
confirmées, elles deviennent des réservations ordinaires au nom de la
société, marquées *réservation en ligne* ; les demandes refusées et annulées
sont repliées sous la liste.
