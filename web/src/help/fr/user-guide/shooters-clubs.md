# Tireurs et sociétés

Chaque société tient ses propres données de base des tireurs ; rien n’est
partagé avec d’autres sociétés (voir [Données de base des tireurs par
société](../explanation/data.md#shooters)). **Tireurs** et **Sociétés** se
trouvent sous **Données de base** dans la barre latérale.

## Tireurs

![Tireurs](../images/screenshots/shooters.png){ .screenshot }

- La **recherche** fonctionne comme dans le tableau de bord : numéro de
  licence (zéros initiaux facultatifs), partie d’un nom, ou prénom et nom
  dans n’importe quel ordre. La liste affiche 50 tireurs par page, la liste
  des sociétés 50 sociétés ; les flèches en dessous tournent les pages.
- Un clic sur un tireur ouvre les détails : nom, numéro de licence, date de
  naissance, adresse et les sociétés auxquelles il appartient.
- **Société** (nouveau tireur ou ajout d’une société dans les détails)
  propose les sociétés de la liste. Un nom qui correspond à une société sans
  tenir compte des majuscules ni des espaces est cette société. Tout autre
  nom affiche *n’est pas encore dans la liste des sociétés* avec les
  sociétés semblables sous *Voulais-tu dire* ; une nouvelle société n’est
  créée qu’avec **Créer la nouvelle société**, ainsi une faute de frappe ne
  devient pas une deuxième société.
- **Importer un CSV** (*Gérer les données de base*) lit un fichier CSV, voir
  [Import CSV](../reference/csv-import.md).
- Les tireurs avec des inscriptions ne peuvent pas être supprimés.

## Sociétés

![Sociétés](../images/screenshots/clubs.png){ .screenshot }

Les **sociétés** sont les sociétés auxquelles appartiennent les tireurs, avec leur numéro FST. Elles apparaissent dans les inscriptions, les classements et le concours de sections. Un tireur peut appartenir à plusieurs sociétés et s’inscrit pour l’une d’elles. Les sociétés utilisées par des inscriptions ne peuvent pas être supprimées, et deux sociétés ne peuvent pas porter le même nom, même avec d’autres majuscules ou espaces. Pourquoi ce ne sont pas la société qui utilise Standblatt : [La société du tireur et la société qui utilise Standblatt](../explanation/data.md#shooter-club).

**Fusionner avec…** (*Gérer les données de base*) supprime un doublon, p. ex. *S.G. Hiltown* à côté de *SG Hilltown* : choisis la société qui reste (les sociétés au nom semblable viennent en premier). Elle garde son nom et reçoit les tireurs, inscriptions et demandes de réservation de l’autre, y compris ceux des manifestations clôturées, si bien que leurs classements affichent la société qui reste ; le numéro FST est repris si elle n’en a pas. L’autre société est supprimée.
