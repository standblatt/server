# Documents et rapports

Chaque document est un PDF dans la langue de la manifestation. Une
impression va à l’imprimante définie pour le document sous **Stands de tir →
Configurer les imprimantes** du stand de l’inscription quand un agent
d’impression du stand est en ligne ; sinon, le PDF s’ouvre (voir
[Impression](../user-guide/printing.md)).

## Documents {#documents}

| Document | Format | Contenu | Imprimé depuis | Réglage d’impression | Module |
|---|---|---|---|---|---|
| Feuille de stand | A5 | manifestation, tireur, QR du [code d’inscription](codes.md), stand (manifestations sur plusieurs stands), cible, programmes avec coups et [codes-barres](codes.md) | *Imprimer la feuille de stand* à l’inscription ; l’inscription | *Feuille de stand (A5)* | Feuille de stand |
| Bande de contrôle | étiquette 90 × 29 mm, ou un [modèle d’impression](../user-guide/print-forms.md#builtin) de la société | manifestation, tireur, licence, société, catégorie d’âge, passes, QR du code d’inscription | *Imprimer la bande de contrôle* à l’inscription ; l’inscription | *Bande de contrôle (étiquette 90 × 29 mm)* | Bande de contrôle |
| Feuille de résultats | A5 | manifestation, tireur, chaque programme avec ses passes, son total et sa distinction | **Imprimer la feuille de résultats** dans les résultats, dès que chaque programme a des résultats | *Feuille de résultats (A5)* | – |
| Carte-couronne | A6, imprimée sur des cartes préimprimées, ou un modèle d’impression de la société | tireur, année de naissance, société, date, organisateur, valeur en chiffres et en lettres (réglage *Imprimer la valeur de la distinction sur les cartes-couronnes*) | **Imprimer la carte-couronne** dans les résultats, pour une distinction carte-couronne | *Carte-couronne (A6, préimprimée)* | Distinctions |
| Modèle d’impression d’un programme | taille du modèle | voir [Champs des modèles](form-placeholders.md#results) | **Imprimer sur le modèle** dans les résultats d’un programme terminé | *Modèles d’impression (préimprimés)* | module de la société Modèles d’impression |
| Document d’inscription | taille du modèle | voir [Champs des modèles](form-placeholders.md#registration) | un bouton portant le nom du modèle dans l’inscription | *Modèles d’impression (préimprimés)* | module de la société Modèles d’impression |
| Diplômes | taille du modèle, une page par tireur classé | voir [Champs des modèles](form-placeholders.md#rank) | **Diplômes** sous **Classements** | *Diplômes* | module de la société Modèles d’impression |
| Classement | A4 | rang, nom, année, catégorie, société, passes, total ; rangs 1–3 mis en évidence ; *(2/3)* pour la tentative qui compte | **PDF** sous **Classements** et **Statistiques** | jamais imprimé, PDF seulement | – |
| Page de test | étiquette 90 × 29 mm | société, stand, imprimante, options, heure | **Page de test** sous **Stands de tir → Configurer les imprimantes** | l’imprimante de cette ligne | – |

- Sans stand, les documents s’ouvrent seulement en PDF.
- Les liens d’aide n’impriment que la feuille de stand (distribution de
  munition) ou les modèles d’impression (cibarre, bureau des résultats), et
  seulement pour leurs cibles.
- Une tâche d’impression que l’agent ne confirme pas dans les deux minutes
  est relancée, au plus trois fois.

## Rapports {#reports}

Sous **Statistiques** (droit *Statistiques et rapports*). Chaque rapport est
disponible en **PDF** et en **CSV** (UTF-8 avec BOM, séparé par `;`, s’ouvre
directement dans Excel). Les textes suivent la langue de la manifestation.

| Rapport | Contenu | Module |
|---|---|---|
| Rapport financier | recettes par programme, distinctions avec leur valeur | – |
| Recettes par programme et catégorie d’âge | participants, prix et sous-total par programme et catégorie | – |
| Participants par société | participants par société, programme et catégorie d’âge | – |
| Distinctions | chaque distinction remise, avec tireur, société, programme, valeur et résultat | Distinctions |
| Distinctions à remettre | distinctions pas encore retirées, avec tireur, société, programme, type et valeur, plus les espèces encore à verser (voir [Résultats](../user-guide/results.md#handout)) | Distinctions |
| Concours de sections | classement, sociétés hors concours et non classées | Concours de sections |
| Prix spéciaux | gagnants par programme et catégorie d’âge, tireurs sautés avec la raison (voir [Programmes](../user-guide/programmes.md#special-prizes)) | Prix spéciaux |
| Classements généraux | chaque [classement général](../user-guide/classifications.md) par classe, avec les points par programme ; les parties qui ne comptent pas sont entre parenthèses | Classements généraux |
| Réservations de cibles | les sociétés qui ont demandé des cibles par la [réservation en ligne](../user-guide/online-booking.md), confirmées ou en attente, avec contact, créneaux, numéros de cible et tireurs annoncés | Réservations |

Le PDF d’un rapport est en A4, *Distinctions* en paysage. Le classement de
chaque programme est aussi un fichier CSV : rang, nom, prénom, année,
catégorie, société, passes, total et mouches, pour les programmes à
plusieurs tentatives aussi la tentative qui compte et le nombre de
tentatives.
