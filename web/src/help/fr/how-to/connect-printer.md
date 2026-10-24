# Connecter une imprimante au stand

Tu en as besoin quand les feuilles de stand, les bandes de contrôle et les
cartes-couronnes doivent sortir de l’imprimante du stand au lieu de s’ouvrir
en PDF. Il faut le droit *Stands et imprimantes* (voir
[Rôles](../user-guide/users.md#roles)).

Un ordinateur à côté des imprimantes fait tourner l’agent d’impression : un
PC Windows ou un Mac avec l’application Standblatt Print Agent, ou un
ordinateur Linux comme un Raspberry Pi avec l’agent sans interface. Pourquoi
cela fonctionne ainsi est expliqué dans
[Impression](../explanation/printing.md).

## Avec l’application sous Windows ou sur un Mac {#app}

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
3. Dans l’application, clique sur **Connecter un stand de tir**. Elle
   affiche un code comme `K7QD-4MZP` et ouvre **Connecter un agent
   d’impression** dans le navigateur.
4. Connecte-toi si demandé, vérifie l’ordinateur nommé sur la page, choisis
   la **Société**, le **Stand de tir** et le **Nom** de l’agent, puis
   **Connecter**. Le code est valable dix minutes.

![Un code dans l’application](../images/app/print-agent-code.png){
.screenshot }

Si le navigateur de cet ordinateur n’est pas à portée de main, ouvre
**Stands de tir → Configurer les imprimantes** du stand sur n’importe quel
autre ordinateur, tape le code dans le champ à côté de **Connecter par
code** et clique dessus : l’ordinateur est connecté à ce stand sous son
propre nom.

Pour un deuxième stand sur le même ordinateur, p. ex. *Pistolet 25 m* et
*Pistolet 50 m*, clique sur **Connecter un autre stand de tir** dans
l’application et répète les étapes 3 et 4.

## Avec l’agent sans interface sous Linux {#headless}

1. Ouvre **Stands de tir** et clique sur **Configurer les imprimantes** du
   stand.
2. Sous **Agents d’impression**, saisis un nom (p. ex. *Pi halle 1*) et
   clique sur **Ajouter un agent d’impression**.
3. Clique sur **copier** pour copier les deux lignes `STANDBLATT_URL` et
   `STANDBLATT_AGENT_TOKEN`. Le jeton n’est affiché que maintenant.
4. Installe et démarre l’agent avec ces réglages comme décrit dans le
   [README de l’agent
   d’impression](https://github.com/standblatt/print-agent#readme).

Pour un deuxième stand sur le même ordinateur, ajoute aussi un agent là-bas
et donne les deux jetons à l’agent, séparés par une virgule.

## Choisir les imprimantes {#printers}

1. Ouvre **Stands de tir → Configurer les imprimantes** du stand. L’agent
   affiche *connecté* avec sa version et son système d’exploitation ; ses
   imprimantes apparaissent dès qu’il s’est connecté.
2. Pour chaque document, choisis l’imprimante, puis ses options (bac,
   papier, recto verso …, telles que l’agent les signale pour cette
   imprimante ; sous Linux, d’autres options `lp` comme `fit-to-page` vont
   dans le champ de texte).
3. Clique sur **Page de test** dans la ligne : une petite page (90 × 29 mm)
   avec société, stand, imprimante, options et heure s’imprime avec les
   options de cette ligne. Une imprimante d’étiquettes l’imprime comme
   étiquette, les autres imprimantes en haut à gauche de la feuille.
4. Clique sur **Enregistrer** dans la ligne. Un document laissé sur
   *navigateur (PDF)* continue de s’ouvrir en PDF.

![Impression](../images/screenshots/printing.png){ .screenshot }

*Diplômes* sont les diplômes depuis les classements, *Modèles d’impression (préimprimés)* tous les autres [modèles d’impression](../user-guide/print-forms.md) de la société. Les classements ne sont jamais envoyés à une imprimante.

## Vérifier que ça imprime {#jobs}

Imprime quelque chose, p. ex. une feuille de stand à la prochaine
inscription, et regarde **Tâches récentes** sur la même page : *terminé*
signifie que c’est imprimé, *échoué* affiche l’erreur à côté. Une tâche que
l’agent ne confirme pas dans les deux minutes est relancée ; après trois
essais, elle échoue.

Pour ne plus utiliser un ordinateur, clique sur **révoquer** à côté de son
agent ; avec l’application, retire aussi le stand sous **Avancé** et
désinstalle l’application. Tous les écrans de la page d’impression sont
décrits dans [Impression](../user-guide/printing.md).
