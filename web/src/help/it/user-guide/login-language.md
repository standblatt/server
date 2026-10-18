# Accesso e lingua

## Accedere

Accedi con il tuo indirizzo e-mail (o nome utente) e la tua password: la
password l’hai scelta accettando l’invito della tua società, oppure l’hai
ricevuta dall’admin della tua società o dal gestore della piattaforma. Gli
account della società non nascono da una registrazione; solo [Il mio
allenamento](../my-training/account.md) ha una sua registrazione. Password
dimenticata? Se il tuo account ha un indirizzo e-mail confermato, usa
**Password dimenticata?** nella pagina di accesso; altrimenti chiedi al
gestore di impostare una nuova password nel [pannello di
amministrazione](../admin-guide/admin-panel.md). Da connesso la cambi da
solo sotto [Account](#account).

Dopo l’accesso arrivi alla pagina iniziale con una scheda per ogni società a
cui appartieni e una per [Il mio allenamento](../my-training/index.md), che
ogni utente ha; **Account** ed **Esci** si trovano sotto le schede. Senza
società la pagina iniziale porta direttamente a Il mio allenamento.

![Login](../images/screenshots/login.png){ .screenshot }

![Pagina iniziale](../images/screenshots/home.png){ .screenshot }

Un accesso resta valido sette giorni senza attività e al massimo 30 giorni;
con **Resta connesso per 90 giorni** spuntato, dura 90 giorni su quel
dispositivo. Dopo troppi tentativi falliti, gli accessi dallo stesso
indirizzo vengono bloccati per un po’.

Se il tuo account appartiene a più società, scegli la società dopo
l’accesso; **Cambia società** nel menu utente riporta a questa scelta.

## Account {#account}

**Account** nel menu utente (in fondo alla barra laterale), nella pagina iniziale e nella barra laterale di Il mio allenamento apre le impostazioni del tuo account:

- **Accesso**: il tuo nome utente (non si può cambiare) e il tuo indirizzo
  e-mail. **Aggiungi** o **Cambia** chiede la tua password e invia un link
  di conferma; finché non confermi, l’indirizzo conta come *non confermato*
  e **Invia di nuovo** manda un nuovo link. Un indirizzo confermato ti
  permette di reimpostare da solo una password dimenticata.
- **Password**: inserisci la password attuale e quella nuova (almeno 10
  caratteri). Gli altri dispositivi vengono disconnessi, quello che stai
  usando resta connesso.
- **Lingua e aspetto**: la lingua dell’interfaccia del tuo account e la
  modalità chiara o scura.
- **Secondo fattore**: vedi sotto.
- **Elimina account**: solo per account senza società; diario di
  allenamento, obiettivi e impostazioni vengono eliminati
  definitivamente. Gli account di società li rimuove l’admin della società.

![Account](../images/screenshots/account.png){ .screenshot }

## Secondo fattore {#totp}

Un’app di autenticazione (Google Authenticator, Aegis, 1Password e simili)
può proteggere il tuo accesso: dopo la password, Standblatt chiede il codice
a sei cifre mostrato dall’app. Configurala sotto [Account](#account):
**Configura**, scansiona il codice QR, inserisci il codice attuale,
**Attiva**. **Disattiva** chiede la tua password. Telefono perso? Il gestore
può azzerare il secondo fattore nel [pannello di
amministrazione](../admin-guide/admin-panel.md). Il gestore può anche
disattivare la funzione per tutto il server.

## Lingua

L’interfaccia è disponibile in inglese, tedesco, francese e italiano. Senza
una scelta, Standblatt segue la lingua del browser. **Lingua** nel menu
utente (in fondo alla barra laterale) la imposta per il tuo account su ogni
dispositivo.

I documenti stampati, le esportazioni CSV e gli schermi chiosco non seguono
la tua scelta ma la lingua della manifestazione, così ogni stampa di una
manifestazione ha lo stesso aspetto (vedi
[Manifestazioni](events.md#settings)).

## Versione

In fondo alla barra laterale c’è la versione del server, utile per segnalare
un problema; in Il mio allenamento sul telefono si trova in
**Impostazioni**. Quando esce una nuova versione, le pagine aperte offrono
**Ricarica** entro un minuto e la caricano al più tardi al prossimo cambio
di pagina.
