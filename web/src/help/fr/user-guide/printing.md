# Impression

Standblatt produit tous les documents en PDF :

| Document | Format | Quand |
|---|---|---|
| Feuille de stand | A5 | à l’inscription, avec code QR et codes-barres des programmes |
| Bande de contrôle | étiquette 90 × 29 mm | à l’inscription, collée sur la carte-cible |
| Feuille de résultats | A5 | quand tous les programmes ont des résultats |
| Carte-couronne | A6, imprimée sur des cartes préimprimées | après le choix d’une carte-couronne comme distinction |
| Modèle d’impression | toute taille, imprimé sur des feuilles préimprimées | pour les programmes avec un [modèle d’impression](print-forms.md), dès que les résultats sont là |
| Classement | A4 | depuis la page des classements, en PDF seulement (pas envoyé à l’imprimante du stand) |

Chaque bouton d’impression imprime sur l’imprimante du stand si elle est
configurée et que son agent d’impression est en ligne ; sinon le PDF s’ouvre
dans le navigateur et peut être imprimé de là. Sur un téléphone ou une
tablette, le PDF passe plutôt par le menu de partage (imprimer, enregistrer
dans Fichiers, envoyer), l’app reste donc ouverte ; si le téléphone demande
d’abord un autre toucher, une note avec le bouton **Partager** apparaît.

## Configurer les imprimantes

**Stands de tir → Configurer les imprimantes** (admins) :

![Impression](../images/screenshots/printing.png){ .screenshot }

1. **Ajouter un agent d’impression**, donne-lui un nom et copie le
   jeton. Installe l’agent sur un ordinateur près des imprimantes (voir
   [Agent d’impression](../admin-guide/print-agent.md)). Un ordinateur peut
   servir plusieurs stands, p. ex. *Pistolet 25 m* et *Pistolet 50 m* :
   ajoute un agent sur chaque stand et donne à l’agent les deux jetons,
   séparés par une virgule.
2. Dès que l’agent est en ligne, il annonce ses imprimantes et la liste
   indique sa version et son système d’exploitation. Choisis pour chaque
   document l’imprimante (*Diplômes* sont les diplômes des classements,
   *Modèles d’impression* tous les autres modèles de la société) et ses
   options : celles que l’agent annonce pour cette imprimante (sous Linux
   bac, format de papier et type de support, plus d’autres options `lp`
   comme `fit-to-page` dans le champ de texte ; sous Windows bac, papier,
   recto verso, mise à l’échelle et couleur).
3. **Page de test** imprime une petite page (90 × 29 mm, la taille de la
   bande de contrôle) sur l’imprimante de cette ligne avec les options
   choisies, même avant d’enregistrer : société, stand de tir, imprimante,
   options et heure. Une imprimante d’étiquettes l’imprime comme étiquette,
   les autres imprimantes en haut à gauche de la feuille. Un agent doit être
   en ligne.
4. **Tâches récentes** montre ce qui a été imprimé et si cela a fonctionné.

Une tâche que l’agent ne confirme pas dans les deux minutes est relancée ;
après trois tentatives, elle échoue. **révoquer** retire un agent
immédiatement.

## Connecter un ordinateur par code {#connect}

L’application Standblatt sur un ordinateur près des imprimantes n’a pas
besoin de jeton : elle affiche un code comme `K7QD-4MZP` et ouvre
**Connecter un agent d’impression** dans le navigateur. Connecte-toi si
nécessaire, vérifie l’ordinateur indiqué, choisis la société, le stand de
tir et le nom de l’agent, puis **Connecter**. Seuls les membres qui peuvent
gérer les stands de tir peuvent connecter un ordinateur ; un code est
valable dix minutes. L’application se connecte ensuite d’elle-même, et
l’agent apparaît sous **Stands de tir → Configurer les imprimantes** avec sa version et son
système d’exploitation.
