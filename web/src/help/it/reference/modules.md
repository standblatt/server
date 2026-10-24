# Moduli

I moduli sono gli elementi costitutivi che una manifestazione o una società
attiva. Ciò che è disattivato scompare dalla barra laterale, dalla
dashboard, dall’editor dei programmi, dalle classifiche e dai rapporti. Come
attivarli: [Moduli](../user-guide/events.md#modules) e [Moduli della
società](../user-guide/events.md#club-modules).

## Moduli della manifestazione {#event-modules}

In **Configurazione** → **Moduli** di una manifestazione.

| Modulo | Cosa aggiunge |
|---|---|
| Assegnazione dei bersagli | i tiratori vengono iscritti su bersagli: la griglia dei bersagli sulla dashboard, lo stato dei bersagli sui display chiosco e i link per marcatori (richiede uno stand; vedi sotto) |
| Prenotazioni | bersagli riservati a gruppi e società; la sua pagina nella configurazione attiva la [prenotazione online](../user-guide/online-booking.md), con cui le società chiedono bersagli tramite un link pubblico (richiede uno stand) |
| Foglio di stand | il foglio di stand stampato all’iscrizione: tiratore, bersaglio, programmi e codici a barre |
| Striscia di controllo | la striscia di controllo stampata all’iscrizione: le serie che il tiratore può tirare |
| Distribuzione munizioni | il pulsante **Munizione** sulla dashboard e per la distribuzione munizioni (vedi [Aiutanti](../user-guide/helpers.md#ammunition)); lascialo disattivato dove i tiratori portano le proprie, p. es. pistola e carabina ad aria compressa |
| Tentativi | programmi che si possono tirare più volte (vedi [Programmi](../user-guide/programmes.md#attempts)) |
| Distinzioni | carte corona, premi in natura e premi in denaro secondo i punti per categoria d’età |
| Premi speciali | i migliori ranghi di ogni categoria d’età (attiva Distinzioni, vedi [Programmi](../user-guide/programmes.md#special-prizes)) |
| Concorso di sezione | le società classificate secondo i risultati dei loro tiratori |
| Classifiche generali | classifiche su più programmi o tentativi (vedi [Classifiche generali](../user-guide/classifications.md)) |
| Display chiosco | classifiche su schermi nella sala (vedi [Classifiche e chiosco](../user-guide/rankings-kiosk.md)) |
| Link aiutante | marcatori, distribuzione munizioni e registrazione risultati senza account (vedi [Aiutanti](../user-guide/helpers.md)) |

### Dove compare un modulo {#places}

| Modulo | Pagine e sezioni |
|---|---|
| Assegnazione dei bersagli | griglia dei bersagli sulla dashboard, bersagli nell’iscrizione, stato dei bersagli sui display chiosco, ruolo *Marcatore* dei link aiutante |
| Prenotazioni | pagina **Prenotazioni** della manifestazione, con le richieste di prenotazione e il loro numero nella barra laterale; pagina **Prenotazioni** nella configurazione (prenotazione online); rapporto *Prenotazioni di bersagli* |
| Foglio di stand | *Stampa foglio di stand* nell’iscrizione |
| Striscia di controllo | *Stampa striscia di controllo* nell’iscrizione |
| Distribuzione munizioni | pulsante **Munizione** e badge *munizione* sulla dashboard e sui link aiutante |
| Tentativi | sezione **Tentativi** della pagina del programma; contatore per programma nell’iscrizione |
| Distinzioni | sezione **Distinzioni** della pagina del programma; distinzioni nei risultati; cifre in **Statistiche**; rapporti *Distinzioni* e *Distinzioni da consegnare* |
| Premi speciali | sezione **Premi speciali** della pagina del programma; rapporto *Premi speciali* |
| Concorso di sezione | *Conta per il concorso di sezione* sulla pagina del programma; concorso di sezione in **Statistiche**; rapporto *Concorso di sezione* |
| Classifiche generali | pagina **Classifiche generali** nella configurazione; tabelle in **Classifiche** e sui display chiosco; rapporto *Classifiche generali* |
| Display chiosco | pagina **Display chiosco** nella configurazione; *Mostra sugli schermi chiosco* sulla pagina del programma |
| Link aiutante | pagina **Link aiutante** nella configurazione |

### Regole {#rules}

| Modulo | Richiede | Richiede uno stand | Impostazioni | Non si può disattivare finché |
|---|---|---|---|---|
| Assegnazione dei bersagli | – | sì | – | un’iscrizione ha un bersaglio |
| Prenotazioni | Assegnazione dei bersagli | sì | la prenotazione online nella sua pagina di configurazione | il link della prenotazione online è attivo |
| Foglio di stand | – | – | – | – |
| Striscia di controllo | – | – | *Modello delle strisce di controllo* | – |
| Distribuzione munizioni | – | – | – | – |
| Tentativi | – | – | – | un programma permette più di un tentativo |
| Distinzioni | – | – | *Stampa il valore della distinzione sulle carte corona* (attivo per una nuova manifestazione), *Modello delle carte corona* | un programma ha livelli di distinzione |
| Premi speciali | Distinzioni | – | – | un programma assegna premi speciali |
| Concorso di sezione | – | – | *Società esclusa dal concorso di sezione* | – |
| Classifiche generali | – | – | – | la manifestazione ha una classifica generale |
| Display chiosco | – | – | *Mostra lo stato dei bersagli sulle classifiche* (con Assegnazione dei bersagli) | la manifestazione ha un link chiosco |
| Link aiutante | – | – | – | la manifestazione ha un link aiutante |

- Un modulo richiesto da un altro resta attivo finché quell’altro è attivo.
- I moduli che richiedono uno stand sono disattivati in una manifestazione
  senza stand, qualunque cosa sia scelta (vedi [Manifestazioni senza
  bersagli](../user-guide/events.md#no-lanes)).
- Le impostazioni si trovano sulla scheda del modulo in **Moduli**; *Mostra
  lo stato dei bersagli sulle classifiche* è sulla pagina **Display
  chiosco**.
- **Copia per l’anno prossimo** e i modelli riprendono i moduli scelti e i
  moduli in uso.
- Una manifestazione creata tramite l’API senza un elenco di moduli riceve
  ogni modulo della manifestazione tranne *Distribuzione munizioni*.

## Moduli della società {#club-modules}

In **Moduli della società** nella parte società della barra laterale
(permesso *Utenti e società*).

| Modulo | Cosa aggiunge | Non si può disattivare finché |
|---|---|---|
| Modelli di stampa | modelli per fogli prestampati e la sezione *Modello di stampa* di ogni programma (vedi [Modelli di stampa](../user-guide/print-forms.md)); pagina **Modelli di stampa** e il permesso *Modelli di stampa* | la società ha un modello di stampa o un’immagine |

Una nuova società parte con *Modelli di stampa* attivato. Se è disattivato,
le manifestazioni stampano la carta corona e la striscia di controllo
integrate.

## Tipi di manifestazione {#kinds}

**Nuova manifestazione** chiede il tipo; attiva solo dei moduli, tutto si può cambiare in seguito.

| Tipo | Moduli |
|---|---|
| Tiro di festa | Assegnazione dei bersagli, Prenotazioni, Foglio di stand, Striscia di controllo, Distinzioni, Premi speciali, Concorso di sezione, Display chiosco, Link aiutante |
| Programma annuale o allenamento | Tentativi, Classifiche generali |
| Tiro in campagna | Assegnazione dei bersagli, Distinzioni, Concorso di sezione, Display chiosco, Link aiutante |
| Campionato sociale | Tentativi, Classifiche generali, Display chiosco |
| Selezione personalizzata | nessuno |

Un [modello](../user-guide/events.md#templates) della società porta invece i
propri moduli.
