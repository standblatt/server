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

| Diritto | Permette |
|---|---|
| **Vedere manifestazioni e classifiche** | manifestazioni, programmi, classifiche e classifiche generali |
| **Gestire le manifestazioni** | iscrizioni, risultati, distinzioni, prenotazioni, documenti e stampa |
| **Configurare le manifestazioni** | creare, configurare, copiare, chiudere ed eliminare manifestazioni, con programmi, classifiche generali, link chiosco e aiutante e prenotazione online |
| **Statistiche e rapporti** | statistiche, rapporti con importi e classifiche in CSV |
| **Vedere tiratori e società** | i dati personali dei dati di base |
| **Modificare tiratori e società** | registrare e modificare tiratori e società |
| **Gestire i dati di base** | importare tiratori, rinominare, unire ed eliminare società, eliminare tiratori |
| **Stand e stampanti** | stand di tiro, agenti di stampa e impostazioni delle stampanti |
| **Modelli di stampa** | modelli di stampa e immagini (solo finché il modulo della società è attivo) |
| **Utenti e società** | utenti, ruoli, inviti e moduli della società |

Ogni società inizia con quattro ruoli:

| Ruolo | Diritti |
|---|---|
| **Admin** | tutti; la colonna non si può modificare e il ruolo non si può eliminare |
| **Comitato** | tutti tranne *Stand e stampanti* e *Utenti e società* |
| **Aiutante** | vedere e gestire le manifestazioni, statistiche, vedere e modificare tiratori e società |
| **Membro** | vedere manifestazioni e classifiche |

**Nuovo ruolo** aggiunge un ruolo proprio, che vede solo le manifestazioni finché non ne spunti di più. La matita rinomina un ruolo (un nome vuoto mantiene il nome predefinito), il cestino lo elimina quando nessuno lo ha più. L’app mostra a ogni membro solo le pagine e i pulsanti che il suo ruolo consente: un membro senza *Gestire le manifestazioni* apre una manifestazione sulle sue classifiche.

Gli account sono globali: un account viene eliminato dal gestore della
piattaforma nel [pannello di
amministrazione](../admin-guide/admin-panel.md). Una password dimenticata si
reimposta con **Password dimenticata?** nella pagina di accesso se l’account
ha un indirizzo e-mail confermato, altrimenti lo fa il gestore.
