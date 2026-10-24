# Comment les classements s’articulent

Une manifestation peut classer ses tireurs de plusieurs façons à la fois :
par programme, sur plusieurs programmes ou tentatives, par société, et par
catégorie d’âge pour les prix spéciaux. Ce ne sont pas des concours séparés
avec des données séparées. Tous partent des mêmes résultats, et la plupart
partent du même classement de programme.

## Le classement de programme {#programme}

Le classement de programme est la base. Chaque programme a un classement de
tous ceux qui l’ont tiré, ordonné selon les [règles de la
FST](../reference/rules.md#ranking) : total, puis les passes de la dernière
en arrière, puis les mouches, puis l’âge. Il n’existe qu’une seule
implémentation de cet ordre, et elle alimente la page des classements, les
écrans kiosque, le PDF et le CSV du classement, les prix spéciaux et le rang
qu’un tireur voit dans *Mon entraînement*. La place d’un tireur est donc la
même partout où elle apparaît.

Les rangs sont consécutifs (1, 2, 3, … même à résultat égal) : le départage
tranche presque toutes les égalités, deux tireurs ne partagent donc jamais
un rang.

Les **variantes** sont des programmes séparés portant le même nom, par exemple *Sport* et *Appuyé*, et chacune a son propre classement : le tir appuyé et le tir à bras franc ne sont pas comparés.

Un programme **combiné** reprend des passes d’un autre programme au lieu
d’être tiré séparément. Ses résultats sont des copies des passes choisies de
sa source ; il est donc classé comme tout autre programme ; quand la source
change, la copie suit, et avec elle le classement.

## Les tentatives et la tentative qui compte {#attempts}

Un programme qui permet plusieurs tentatives garde chaque tentative comme
participation avec ses propres résultats. Un classement, en revanche, liste
chaque tireur une seule fois. La **tentative qui compte** du programme (la
meilleure, la première ou la dernière) décide quelle tentative représente le
tireur, et le classement montre laquelle c’est, p. ex. *2/3*.

La même tentative qui compte est utilisée partout où un seul résultat par
tireur est nécessaire : le classement, le concours de sections, les prix
spéciaux et le rang dans *Mon entraînement*. Les distinctions sont
différentes, car une distinction appartient à un résultat plutôt qu’à un
classement : chaque tentative peut obtenir sa propre distinction, sauf si
*Seule la tentative qui compte reçoit une distinction* est activé ; dans ce
cas, la distinction passe à la tentative qui compte.

Les combinaisons reprennent les passes d’un seul résultat ; avec plusieurs
tentatives, on ne saurait pas lequel, c’est pourquoi les programmes à
plusieurs tentatives ne peuvent pas être combinés.

## Classements généraux {#classifications}

Un classement général classe les tireurs sur plusieurs programmes ou sur
plusieurs tentatives d’un programme : la somme de trois programmes pour un
championnat de section, les trois meilleures tentatives d’un programme
annuel pour un championnat de société. Il prend les tentatives qui comptent
de chaque programme, additionne ce qui compte, peut ne compter que les
meilleurs programmes, et classe avec son propre départage et ses propres
classes de catégories d’âge.

Un classement général ne fait que lire des résultats. Il ne modifie jamais
un résultat, une distinction ou un classement de programme, et une
manifestation peut avoir plusieurs classements généraux à côté de ses
classements de programme. Le classement de programme lui-même peut être vu
comme le classement général le plus simple : un programme, une tentative, le
départage FST et pas de classes. Voir [Classements
généraux](../user-guide/classifications.md) pour les réglages.

## Concours de sections {#club-competition}

Le concours de sections classe des sociétés, pas des tireurs. Il n’utilise
que les programmes marqués *Compte pour le concours de sections* et
seulement les résultats complets des tireurs inscrits pour une société ; la
société est celle choisie à l’inscription, pas la société qui utilise
Standblatt. Chaque résultat compte avec ses deux premières passes, et le
résultat d’une société est la moyenne de ses meilleurs résultats, avec un
petit bonus pour chaque participant supplémentaire (voir [Concours de
sections](../reference/rules.md#club-competition)). Comme il a son propre
score, le concours de sections ne suit pas le classement de programme.

## Prix spéciaux {#special-prizes}

Les prix spéciaux vont aux meilleurs rangs de chaque catégorie d’âge. Eux
suivent bien le classement de programme : ils le parcourent par catégorie
d’âge et sautent les tireurs dont la société ou la catégorie a trop peu de
participants (voir [Prix
spéciaux](../reference/rules.md#special-prizes)). Leurs réglages peuvent
être modifiés pendant la manifestation, car ils ne font que lire le
classement.

## Quel classement utilise quoi {#overview}

| Classement | Classe | À partir de | Tentative qui compte | Ordre |
|---|---|---|---|---|
| Classement de programme | tireurs | un programme | oui | règles de la FST |
| Classement général | tireurs, par classe | plusieurs programmes ou tentatives | son propre nombre par programme | propre départage, FST par défaut |
| Concours de sections | sociétés | programmes qui comptent pour le concours de sections | oui | résultat de la société |
| Prix spéciaux | tireurs, par catégorie d’âge | le classement de programme | oui | le classement de programme |

Les distinctions ne sont pas un classement : un tireur obtient une
distinction avec son propre total dans sa catégorie d’âge, quels que soient
les résultats des autres.

## Pourquoi les passes sont verrouillées {#locked}

Les résultats sont enregistrés par passe. Dès que des tireurs sont inscrits
à un programme, modifier le nombre de passes, leurs coups, leur maximum ou
le fait qu’elles comptent les mouches ou les coups individuels changerait le
sens des résultats enregistrés : un total de 96 sur 100 n’est pas un total
sur 200, et chaque classement qui les compare serait faux. Standblatt fixe
donc les passes dès qu’un tireur est inscrit.

Tout ce qui ne change pas le sens d’un résultat reste ouvert : noms, prix,
valeurs et limites des distinctions, les options. Des catégories d’âge et
des paliers de distinctions peuvent être ajoutés mais pas retirés, car les
inscriptions ont été facturées et les distinctions décernées selon eux.

## En direct et définitif {#live}

Tous les classements sont calculés à partir des résultats chaque fois qu’ils
sont affichés ; un résultat enregistré apparaît donc aussitôt sur la page
des classements et sur les écrans kiosque. Pendant la manifestation, ils
peuvent encore changer. [Clôturer](closing.md) la manifestation les rend
définitifs.
