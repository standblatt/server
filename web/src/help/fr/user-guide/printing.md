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

1. Connecte un ordinateur près des imprimantes : un PC Windows ou un Mac
   avec l’[application Standblatt Print Agent](#app) par [code](#connect),
   ou un ordinateur Linux comme un Raspberry Pi avec l’[agent sans
   interface](../admin-guide/print-agent.md) : **Ajouter un agent
   d’impression**, donne-lui un nom et copie le jeton. Un ordinateur peut
   servir plusieurs stands, p. ex. *Pistolet 25 m* et *Pistolet 50 m* :
   connecte-le à chaque stand (l’application avec **Connecter un autre stand
   de tir**, l’agent sans interface avec les deux jetons, séparés par une
   virgule).
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

## L’application Standblatt Print Agent {#app}

L’application est la solution simple pour un portable ou un PC près des
imprimantes : pas de jeton ni de réglages à saisir.

1. Télécharge-la depuis la [dernière
   version](https://github.com/standblatt/print-agent/releases/latest) :
   `Standblatt-Print-Agent-Setup.exe` pour Windows,
   `standblatt-print-agent-app-darwin-arm64.dmg` pour un Mac avec puce Apple
   (`-amd64.dmg` pour un Mac Intel), le `.deb` ou le `.rpm` pour Linux.
2. Installe-la. L’application n’est pas encore signée : Windows demande une
   fois, **Informations complémentaires → Exécuter quand même** ; sur un
   Mac, glisse-la dans *Applications* et ouvre-la la première fois avec un
   clic droit → **Ouvrir**. Sous Windows, elle s’installe pour l’utilisateur
   actuel sans droits d’administrateur, avec un raccourci dans le menu
   Démarrer et sur le bureau, et démarre avec Windows.
3. Clique sur **Connecter un stand de tir** et confirme le code dans
   Standblatt (voir [plus bas](#connect)).

![L’application Standblatt Print Agent](../images/app/print-agent.png){
.screenshot }

La fenêtre montre pour chaque stand s’il imprime, avec ses dernières
impressions, et en dessous les imprimantes de l’ordinateur, chacune avec une
**Page de test** qui imprime sur cet ordinateur sans Standblatt. **Avancé**
contient *Démarrer avec l’ordinateur*, l’adresse de Standblatt (à modifier
seulement pour un serveur de test) et le retrait d’un stand. Si tu fermes la
fenêtre, l’application continue dans la zone de notification : l’icône est
verte quand chaque stand imprime, sinon rouge.

Sous Windows et sur un Mac, l’application signale quand une nouvelle version
est disponible ; **Mettre à jour maintenant** l’installe et redémarre
l’application. Pour ne plus utiliser un ordinateur, retire son stand sous
**Avancé**, révoque son agent sous **Stands de tir → Configurer les
imprimantes** et désinstalle l’application (sous Windows dans *Paramètres →
Applications*).

## Connecter un ordinateur par code {#connect}

![Un code dans l’application](../images/app/print-agent-code.png){
.screenshot }

L’application Standblatt sur un ordinateur près des imprimantes n’a pas
besoin de jeton : elle affiche un code comme `K7QD-4MZP` et ouvre
**Connecter un agent d’impression** dans le navigateur. Connecte-toi si
nécessaire, vérifie l’ordinateur indiqué, choisis la société, le stand de
tir et le nom de l’agent, puis **Connecter**. Seuls les membres qui peuvent
gérer les stands de tir peuvent connecter un ordinateur ; un code est
valable dix minutes. L’application se connecte ensuite d’elle-même, et
l’agent apparaît sous **Stands de tir → Configurer les imprimantes** avec sa
version et son système d’exploitation.

Si le navigateur de cet ordinateur n’est pas à portée de main, saisis plutôt
le code sous **Stands de tir → Configurer les imprimantes** du stand, dans
**Connecter par code** : l’ordinateur est connecté à ce stand sous son
propre nom.
