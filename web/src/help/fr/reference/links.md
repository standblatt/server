# Liens sans connexion

Ces liens ouvrent une page sans connexion. Quiconque a le lien peut
l’utiliser ; ne l’envoie donc qu’aux personnes à qui il est destiné.

| Lien | Créé par | Affiche | Valable |
|---|---|---|---|
| Lien kiosque `/kiosk/<secret>` | **Créer un lien kiosque** sous **Affichages kiosque** dans la configuration | les classements en direct de la manifestation : programmes et classements généraux avec *Afficher sur les écrans kiosque*, et l’état des cibles avec *Afficher l’état des cibles sur les classements* | jusqu’à **révoquer** ou la suppression de la manifestation |
| Lien d’aide `/h/<secret>` | **Créer un lien d'aide** sous **Liens d'aide** dans la configuration | un accueil réduit de la manifestation pour un [rôle d’aide](permissions.md#helper-links) | jusqu’à **révoquer** ou la suppression de la manifestation |
| Lien de réservation `/b/<code>` | **Activer** sous **Réservations** dans la configuration | les créneaux de la manifestation avec les cibles libres par stand, et le formulaire de demande | jusqu’à **Nouveau lien** ou **Désactiver** ; demandes seulement jusqu’au délai |
| Lien de gestion de la réservation `/b/manage/<secret>` | envoyé par e-mail à la société après **Envoyer la demande** | le statut de la demande de la société ; la modifier ou l’annuler | modifications et annulation jusqu’au délai |
| Lien d’invitation `/invite/<secret>` | **Inviter par e-mail** sous **Utilisateurs** | la société et le rôle ; accepter avec un compte ou un nouveau mot de passe | sept jours ; prend fin à l’acceptation, au retrait ou à un nouvel envoi |
| Page de connexion de l’agent `/connect?code=<code>` | ouverte par l’application Standblatt Print Agent | l’ordinateur du code ; société, stand et nom de l’agent | le code : dix minutes ; demande une connexion avec *Stands et imprimantes* |
| Confirmation de l’e-mail `/verify?token=…` | ajout ou modification de l’adresse e-mail, inscription | confirme l’adresse | sept jours |
| Lien de mot de passe `/reset-password?token=…` | **Mot de passe oublié ?** ; l’exploitant dans le panneau d’administration | définit un nouveau mot de passe | une heure ; de l’exploitant sept jours |

La partie secrète des liens kiosque, d’aide, de gestion et d’invitation est
affichée ou envoyée une seule fois et enregistrée uniquement sous forme de
hachage : un lien perdu ne peut pas être réaffiché, seulement remplacé.

## Paramètres d’URL {#parameters}

| Paramètre | Lien | Effet |
|---|---|---|
| `?theme=dark`, `?theme=light` | kiosque | affichage sombre ou clair |
| `?language=de` (`fr`, `it`, `en`) | kiosque, aide | langue de la page au lieu de la langue de la manifestation |
| `?range=<id>` | kiosque | pour une manifestation sur plusieurs stands : seulement les cibles de ce stand, ses programmes et les programmes sans stand |
| `?code=` | connect | le code à confirmer, rempli par l’application |

**Créer un lien kiosque** inscrit dans le lien le thème, la langue et le stand qui y ont été choisis.
