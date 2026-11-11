# Prenotazione online

Con la **prenotazione online**, le società esterne richiedono da sole i
bersagli. L’organizzatore inserisce un link nell’invito; una società sceglie
i bersagli desiderati in ogni fascia oraria e annuncia i suoi tiratori. Un
admin conferma o rifiuta ogni richiesta.

## Attivare {#setup}

La prenotazione online fa parte del modulo *Prenotazioni* (richiede uno
stand e *Assegnazione dei bersagli*): in **Prenotazioni** nella
configurazione basta **Attiva**. La pagina mostra poi il link con un codice
QR per l’invito, se le società possono fare richiesta in questo momento
(*Apre il*, *Aperta fino al* o *Termine scaduto*) e le fasce orarie offerte,
per ogni giorno con l’inizio di ogni fascia oraria e le pause.

![Prenotazione online nella
configurazione](../images/screenshots/booking-setup.png){ .screenshot }

Senza modifiche, ogni giorno della manifestazione è aperto durante i suoi
orari di tiro (vedi [Manifestazioni](events.md#settings)), altrimenti dalle
08:00 alle 17:00, in fasce di 30 minuti, e le richieste sono accettate da
subito fino a 14 giorni prima del primo giorno. Una fascia oraria può essere
richiesta fino a due ore prima dell’inizio. **Adatta** lo cambia:

- **Tiro dalle** e **Tiro fino alle**: gli orari di ogni giorno senza orari
  di tiro.
- **Pausa dalle** e **Pausa fino alle**: una pausa facoltativa, p. es. a
  mezzogiorno.
- **Durata di una fascia oraria (minuti)**: quanto dura una fascia oraria.
- **Bersagli per fascia oraria e società**: un limite per richiesta; vuoto
  significa nessun limite.
- **Iscrizioni dal**: il primo giorno in cui le società possono richiedere;
  fino ad allora il link indica solo da quando.
- **Termine d’iscrizione**: l’ultimo giorno in cui le società possono
  richiedere o modificare la richiesta, al più tardi l’ultimo giorno di
  tiro.
- **Giorni con altri orari**: una riga per giorno della manifestazione,
  p. es. venerdì 17:00–20:00. I giorni vengono dalla manifestazione e la
  seguono quando ne aggiungi o togli uno.

**Nuovo link** sostituisce il link; quello vecchio smette di funzionare. **Disattiva** chiude la prenotazione; le richieste e le prenotazioni confermate restano. *Copia per l’anno prossimo* riprende gli orari, ma né il link, né l’apertura delle iscrizioni, né il termine.

## Cosa vede una società {#club}

Il link apre una pagina senza login. Spiega la procedura in tre passi e
mostra ogni fascia oraria che inizia tra più di due ore con i bersagli di
ogni stand come pulsanti numerati: verde è libero, giallo è libero ma
richiesto da un’altra società la cui richiesta non è ancora confermata,
rosso sbiadito è occupato, e la scelta della società è scura con una cornice
verde. Prima del giorno di apertura indica solo da quando la prenotazione è
aperta. Solo le prenotazioni confermate occupano un bersaglio; i nomi delle
altre società non vengono mai mostrati. Con più di due giorni, ogni giorno è
chiuso e mostra quante fasce orarie la società vi ha scelto.

![La pagina di prenotazione di una
società](../images/screenshots/booking-public.png){ .screenshot }

La società tocca i bersagli desiderati in ogni fascia oraria (un secondo
tocco ne toglie uno; con un limite per fascia oraria gli altri pulsanti si
bloccano appena è raggiunto) e inserisce il suo nome, una persona di
contatto, e-mail e telefono, e i suoi **Tiratori**: cognome, nome, data di
nascita e facoltativamente la licenza. **Incolla elenco** riprende righe
copiate da Excel (cognome, nome, data di nascita, licenza). L’elenco dei
tiratori può seguire più tardi; i programmi si scelgono allo stand come al
solito.

Una barra in fondo allo schermo conta le fasce orarie e i bersagli scelti e
contiene **Invia la richiesta**. Poi la società riceve un’e-mail con un link
personale. Tramite questo può vedere lo stato, modificare o annullare la
richiesta fino al termine. Cambiare le fasce orarie di una richiesta
confermata libera i suoi bersagli e la richiesta attende una nuova conferma;
aggiungere tiratori la lascia confermata.

## Confermare le richieste {#confirm}

Ogni admin con un’e-mail riceve un messaggio per ogni richiesta nuova,
modificata o annullata. Le richieste sono elencate in cima a
**Prenotazioni**, prima quelle aperte e, per ogni stato, nell’ordine di
arrivo; la barra laterale mostra con un contatore rosso quante sono in
attesa. Una richiesta aperta indica i bersagli desiderati dalla
società. Ogni richiesta mostra al secondo quando è stata ricevuta, una
aperta anche quando la società l’ha modificata l’ultima volta; così, quando
due società vogliono gli stessi bersagli, la prima può ottenerli.

![Richieste di prenotazione sulla pagina delle
prenotazioni](../images/screenshots/booking-requests.png){ .screenshot }

- **Conferma** assegna a ogni fascia oraria i bersagli scelti dalla società,
  come prenotazione ordinaria con il nome della società che appare come
  riservata una fascia oraria prima dell’inizio, e invia i numeri dei
  bersagli alla società per e-mail. Un bersaglio scelto che nel frattempo è
  occupato viene sostituito da un altro bersaglio libero dello stand, se
  possibile vicino. Se una fascia oraria ha troppo pochi bersagli liberi,
  non viene prenotato nulla e il messaggio indica la fascia; spostate una
  prenotazione o rifiutate.
- **Conferma** chiede prima a quale **Società** dei dati di base appartiene
  la richiesta, perché la società ha digitato il nome da sé: la società con
  lo stesso nome è già scelta, le società con un nome simile vengono
  proposte e *Nuova società* la crea. La richiesta mostra poi la società
  attribuita.
- **Conferma** riprende anche i tiratori nei dati di base, in quella
  società: un tiratore viene trovato tramite la licenza, altrimenti tramite
  nome e data di nascita, e creato in caso contrario. I tiratori esistenti
  non vengono mai sovrascritti, perché la lista viene da un modulo pubblico;
  si completa solo una data di nascita o una licenza mancante. I tiratori
  aggiunti dopo la conferma si riprendono con **Riprendi tiratori**.
- **Rifiuta** chiede un motivo facoltativo, che la società riceve via
  e-mail. I bersagli occupati dalla richiesta tornano liberi.

I bersagli riservati sono normali [prenotazioni](reservations.md): modifica
lì i numeri dei bersagli se serve. Durante la manifestazione l’iscrizione
trova i tiratori con la ricerca normale. Il rapporto *Prenotazioni di
bersagli* in [Statistiche e rapporti](statistics-reports.md) elenca ogni
società con le sue fasce orarie, i bersagli e i tiratori.
