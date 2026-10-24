# Classer un championnat de section sur plusieurs programmes

Tu en as besoin quand le classement additionne plusieurs programmes d’une
manifestation, p. ex. un championnat de section comme somme de
*Obligatoire*, *Tir en campagne* et *Cible section*. Tous les programmes
doivent appartenir à la même manifestation. Il faut le droit *Configurer les
manifestations* (voir [Rôles](../user-guide/users.md#roles)).

1. Ouvre la manifestation et va dans **Configuration** → **Modules**. Active
   **Classements généraux**, sauf si le genre de la manifestation l’a déjà
   fait.
2. Sous **Programmes**, assure-toi que chaque programme du championnat
   existe (voir [Programmes](../user-guide/programmes.md)).
3. Ouvre **Classements généraux** dans la configuration et clique sur
   **Nouveau classement général**.
4. Saisis le **Nom**, p. ex. *Championnat de section*.
5. Choisis le premier programme comme **Programme 1**, puis **Ajouter un
   programme** pour chacun des suivants. Laisse **Tentatives qui comptent**
   à 1 pour chaque programme.
6. Active **Classer seulement les tireurs avec tous les résultats qui
   comptent**. Les tireurs à qui il manque encore un programme sont listés
   sous le classement, sans rang.
7. Laisse **Ne compter que les meilleurs programmes** à 0, pour que chaque
   programme compte. Pour *les deux meilleurs sur trois*, saisis 2.
8. Vérifie le **Départage** : par défaut, c’est l’ordre FST ; les flèches
   modifient l’ordre des critères.
9. Ajoute des **Classes** si le championnat classe des groupes d’âge
   séparément, p. ex. *Juniors* avec U17 à U21 et *Élite* avec les
   autres. Sans classes, tout le monde est dans un seul classement.
10. Active **Afficher sur les écrans kiosque** si le classement doit passer
    dans la halle (seulement avec le module *Affichages kiosque*).
11. Clique sur **Créer**.

Les résultats sont saisis comme d’habitude pour chaque
programme. **Classements** affiche le championnat après les classements de
programme, un tableau par classe ; son **PDF** et son **CSV** sont le
rapport *Classements généraux*, et **Diplômes** imprime des diplômes pour
les meilleurs rangs sur un [modèle
d’impression](../user-guide/print-forms.md#certificates).

Comment les totaux et les égalités sont décidés est expliqué dans [Comment
les classements s’articulent](../explanation/rankings.md).
