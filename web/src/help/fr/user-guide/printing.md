# Impression

Standblatt produit tous les documents en PDF :

| Document | Format | Quand |
|---|---|---|
| Feuille de stand | A5 | à l’inscription, avec code QR et codes-barres des programmes |
| Bande de contrôle | étiquette 90 × 29 mm | à l’inscription, collée sur la carte-cible |
| Feuille de résultats | A5 | quand tous les programmes ont des résultats |
| Carte-couronne | A6, imprimée sur des cartes préimprimées | après le choix d’une carte-couronne comme distinction |
| Classement | A4 | depuis la page des classements, en PDF seulement (pas envoyé à l’imprimante du stand) |

Chaque bouton d’impression imprime sur l’imprimante du stand si elle est
configurée et que son agent d’impression est en ligne ; sinon le PDF s’ouvre
dans le navigateur et peut être imprimé de là. Sur un téléphone ou une
tablette, le PDF passe plutôt par le menu de partage (imprimer, enregistrer
dans Fichiers, envoyer), l’app reste donc ouverte ; si le téléphone demande
d’abord un autre toucher, une note avec le bouton **Partager** apparaît.

## Configurer les imprimantes

**Stands de tir → Impression** (admins) :

![Impression](../images/screenshots/printing.png){ .screenshot }

1. **Ajouter un agent d’impression**, donne-lui un nom et copie le
   jeton. Installe l’agent sur un ordinateur près des imprimantes (voir
   [Agent d’impression](../admin-guide/print-agent.md)).
2. Dès que l’agent est en ligne, il annonce ses imprimantes. Choisis pour
   chaque document l’imprimante et ses options (bac, format de papier, type
   de support, autres options `lp` comme `fit-to-page`).
3. **Tâches récentes** montre ce qui a été imprimé et si cela a fonctionné.

Une tâche que l’agent ne confirme pas dans les deux minutes est relancée ;
après trois tentatives, elle échoue. **révoquer** retire un agent
immédiatement.
