# Impression

Standblatt génère tous les documents en PDF. Chaque document avec son
format, son contenu et son réglage d’impression :
[Documents](../reference/documents.md#documents).

Chaque bouton d’impression imprime sur l’imprimante du stand si elle est
configurée et que son agent d’impression est en ligne ; sinon, le PDF
s’ouvre (voir [Quand il n’y a pas
d’imprimante](../explanation/printing.md#fallback)).

Pourquoi un agent d’impression et comment il fonctionne : [Comment
fonctionne l’impression](../explanation/printing.md).

## Configurer les imprimantes

**Stands de tir → Configurer les imprimantes** (droit *Stands et imprimantes*) a trois parties :

![Impression](../images/screenshots/printing.png){ .screenshot }

- **Agents d’impression** : les ordinateurs qui impriment pour ce stand,
  chacun avec sa version et son système d’exploitation, *connecté* ou le
  moment où il a été vu pour la dernière fois, et son nombre
  d’imprimantes. **Ajouter un agent d’impression** ajoute un agent sans
  interface et affiche son jeton une fois, **Connecter par code** connecte
  l’[application](#app) par son [code](#connect), **révoquer** retire un
  agent immédiatement.
- **Documents** : par document, l’imprimante (*navigateur (PDF)* continue
  d’ouvrir le PDF) et les options que l’agent signale pour cette imprimante
  ; sous Linux, d’autres options `lp` comme `fit-to-page` vont dans le champ
  de texte. *Diplômes* sont les diplômes depuis les classements, *Modèles
  d’impression (préimprimés)* tous les autres modèles de la
  société. **Enregistrer** enregistre la ligne. **Page de test** imprime une
  petite page (90 × 29 mm) avec société, stand, imprimante, options et heure
  sur l’imprimante de cette ligne, même avant d’enregistrer ; un agent doit
  être en ligne.
- **Tâches récentes** montre ce qui a été imprimé et si cela a fonctionné.

Une tâche que l’agent ne confirme pas dans les deux minutes est relancée ;
après trois essais, elle échoue.

Pas à pas : [Connecter une imprimante au
stand](../how-to/connect-printer.md).

## L’application Standblatt Print Agent {#app}

L’application est la solution simple pour un portable ou un PC près des
imprimantes : pas de jeton ni de réglages à saisir.

La télécharger, l’installer et la connecter : [Connecter une imprimante au
stand](../how-to/connect-printer.md#app).

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

**Connecter un agent d’impression** (ouvert par l’application avec un code comme `K7QD-4MZP`) montre l’ordinateur et demande la société, le stand et le nom ; seuls les membres qui peuvent gérer les stands peuvent connecter un ordinateur, et un code est valable dix minutes. **Connecter par code** sous **Stands de tir → Configurer les imprimantes** le connecte à ce stand sous son propre nom. Étapes : [Connecter une imprimante au stand](../how-to/connect-printer.md#app).
