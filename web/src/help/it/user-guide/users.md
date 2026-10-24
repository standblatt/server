# Utenti

**Utenti** elenca gli account con accesso alla società, gli inviti aperti e i ruoli della società. La pagina è aperta ai ruoli con il diritto *Utenti e società*, quindi almeno agli admin.

![Utenti](../images/screenshots/users.png){ .screenshot }

## Membri e inviti {#invitations}

- **Invita via e-mail**: incolla uno o più indirizzi (separati da virgole o
  uno per riga, fino a 300) e scegli il ruolo. Ogni indirizzo riceve
  un’e-mail con un link valido sette giorni; l’elenco sotto il modulo mostra
  cosa è successo a ogni indirizzo. Le e-mail partono a circa quattro al
  secondo, quindi un elenco grande richiede un minuto.
- **Accettare**: il link apre una pagina che dice quale società invita e con
  quale ruolo. Senza account la persona sceglie una password ed è subito
  connessa; l’indirizzo e-mail è il login. Chi ha già un account con
  l’indirizzo accede prima e poi accetta. Solo allora la persona compare tra
  i membri.
- **Inviti** sotto i membri elenca quelli aperti e scaduti: **Invia di
  nuovo** manda un nuovo link (il vecchio non funziona più), **Ritira**
  annulla l’invito.
- **Aggiungi utente**: per un account senza e-mail, per esempio un account
  comune allo stand: un account esistente viene aggiunto tramite il nome
  utente; per un nuovo account inserisci anche una password (almeno 10
  caratteri).
- **Ruolo**: scegli il ruolo di un membro nell’elenco. Una società mantiene
  sempre almeno un admin.
- **Rimuovi** toglie l’accesso a questa società; l’account stesso resta e
  mantiene le sue altre società.

## Ruoli {#roles}

Cosa può fare un membro dipende dal suo ruolo. La matrice sotto **Ruoli** ha
una colonna per ruolo e una riga per diritto; una spunta dà subito il
diritto.

![Ruoli](../images/screenshots/roles.png){ .screenshot }

Ogni autorizzazione e cosa permette:
[Autorizzazioni](../reference/permissions.md#permissions). Ogni società
parte con i ruoli **Admin**, **Comitato**, **Aiutante** e **Membro** (vedi
[Ruoli predefiniti](../reference/permissions.md#default-roles)).

**Nuovo ruolo** aggiunge un ruolo proprio, che vede solo le manifestazioni finché non ne spunti di più. La matita rinomina un ruolo (un nome vuoto mantiene il nome predefinito), il cestino lo elimina quando nessuno lo ha più. L’app mostra a ogni membro solo le pagine e i pulsanti che il suo ruolo consente: un membro senza *Gestire le manifestazioni* apre una manifestazione sulle sue classifiche.

I membri con *Utenti e società* ma senza il ruolo Admin restano entro i
propri diritti: assegnano, modificano, rimuovono e invitano solo con ruoli
di cui hanno essi stessi i diritti, e creano o modificano ruoli solo entro
questi limiti. Solo gli admin nominano, modificano o rimuovono quindi degli
admin.

Gli account sono globali (vedi [Un account, un ruolo per
società](../explanation/data.md#accounts)): un account viene eliminato dal
gestore della piattaforma nel pannello di amministrazione. Per una password
dimenticata vedi [Password dimenticata](login-language.md#password).
