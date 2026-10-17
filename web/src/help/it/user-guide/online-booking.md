# Prenotazione online

Con la **prenotazione online** le società esterne richiedono da sé i
bersagli. L’organizzatore inserisce un link nell’invito; una società sceglie
le fasce orarie desiderate, indica quanti bersagli le servono in ciascuna e
annuncia i suoi tiratori. Un admin conferma o rifiuta ogni richiesta.

## Attivare {#setup}

Nella scheda *Prenotazione online* della **Configurazione** della
manifestazione basta **Attiva**. La scheda mostra poi il link con un codice
QR per l’invito e le fasce orarie offerte.

![Prenotazione online nella
configurazione](../images/screenshots/booking-setup.png){ .screenshot }

Senza modifiche, ogni giorno della manifestazione è aperto durante i suoi
orari di tiro (vedi [Manifestazioni](events.md#settings)), altrimenti dalle
08:00 alle 17:00, in fasce di 30 minuti, e le richieste sono accettate fino
a 14 giorni prima del primo giorno. **Adatta** lo cambia:

- **Tiro dalle** e **Tiro fino alle**: gli orari di ogni giorno senza orari
  di tiro.
- **Pausa dalle** e **Pausa fino alle**: una pausa facoltativa, p. es. a
  mezzogiorno.
- **Durata di una fascia oraria (minuti)**: quanto dura una fascia oraria.
- **Bersagli per fascia oraria e società**: un limite per richiesta; vuoto
  significa nessun limite.
- **Termine d’iscrizione**: l’ultimo giorno in cui le società possono
  richiedere o modificare la richiesta.
- **Giorni con altri orari**: una riga per giorno della manifestazione,
  p. es. venerdì 17:00–20:00. I giorni vengono dalla manifestazione e la
  seguono quando ne aggiungi o togli uno.

**Nuovo link** sostituisce il link; quello vecchio smette di funzionare. **Disattiva** chiude la prenotazione; le richieste e le prenotazioni confermate restano. *Copia per l’anno prossimo* riprende gli orari, ma né il link né il termine.

## Cosa vede una società {#club}

Il link apre una pagina senza login. Mostra la manifestazione, il termine e
ogni fascia oraria con i bersagli ancora liberi per poligono. Contano solo
le prenotazioni confermate; i nomi delle altre società non vengono mai
mostrati.

![La pagina di prenotazione di una
società](../images/screenshots/booking-public.png){ .screenshot }

La società inserisce il numero di bersagli per fascia oraria, il suo nome,
una persona di contatto, e-mail e telefono, e i suoi **Tiratori**: cognome,
nome, data di nascita e facoltativamente la licenza. **Incolla lista**
riprende righe copiate da Excel (cognome, nome, data di nascita,
licenza). La lista dei tiratori può seguire più tardi; i programmi si
scelgono al poligono come sempre.

Dopo **Invia la richiesta** la società riceve un’e-mail con un link
personale. Tramite questo può vedere lo stato, modificare o annullare la
richiesta fino al termine. Cambiare le fasce orarie di una richiesta
confermata libera i suoi bersagli e la richiesta attende una nuova conferma;
aggiungere tiratori la lascia confermata.

## Confermare le richieste {#confirm}

Ogni admin con un’e-mail riceve un messaggio per ogni richiesta nuova,
modificata o annullata. Le richieste compaiono in cima a **Prenotazioni**;
la barra laterale indica quante sono in attesa.

![Richieste di prenotazione sulla pagina delle
prenotazioni](../images/screenshots/booking-requests.png){ .screenshot }

- **Conferma** assegna a ogni fascia oraria bersagli liberi del suo
  poligono, possibilmente vicini, come prenotazione ordinaria a nome della
  società, e invia i numeri dei bersagli alla società. Se una fascia ha
  troppo pochi bersagli liberi non viene riservato nulla e il messaggio
  indica la fascia; sposta una prenotazione o rifiuta.
- **Conferma** riprende anche i tiratori nei dati di base: un tiratore viene
  trovato tramite la licenza, altrimenti tramite nome e data di nascita, e
  creato in caso contrario; la società viene creata se serve. I tiratori
  esistenti non vengono mai sovrascritti, perché la lista viene da un modulo
  pubblico; si completa solo una data di nascita o una licenza mancante. I
  tiratori aggiunti dopo la conferma si riprendono con **Riprendi
  tiratori**.
- **Rifiuta** chiede un motivo facoltativo, che la società riceve via
  e-mail. I bersagli occupati dalla richiesta tornano liberi.

I bersagli riservati sono normali [prenotazioni](reservations.md): modifica
lì i numeri dei bersagli se serve. Durante la manifestazione l’iscrizione
trova i tiratori con la ricerca normale. Il rapporto *Prenotazioni di
bersagli* in [Statistiche e rapporti](statistics-reports.md) elenca ogni
società con le sue fasce orarie, i bersagli e i tiratori.
