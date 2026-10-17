# Dashboard

La dashboard è la vista dell’accoglienza di una manifestazione: chi tira su
quale bersaglio, chi aspetta, chi ha finito.

![Dashboard](../images/screenshots/dashboard.png){ .screenshot }

- **Al tiro**: iscrizioni su un bersaglio, con il tempo dall’inizio, i
  pulsanti **Risultati**, **Modifica** e **Termina**. Se l’impostazione
  della manifestazione *L’organizzatore consegna le munizioni* è attiva,
  **Munizione** registra che sono state consegnate e poi lo mostra un badge
  verde *munizione* (vedi [Aiutanti](helpers.md)).
- **In attesa**: iscrizioni senza bersaglio (per esempio quando tutti i
  bersagli sono occupati).
- **Terminati**: iscrizioni i cui programmi hanno tutti i risultati, con il
  badge *distinzioni da consegnare*, che apre il dialogo dei risultati per
  la consegna, o poi *distinzioni consegnate* (vedi
  [Risultati](results.md#handout)).
- **Bersagli**: ogni bersaglio è *libero*, *occupato* (con il tiratore),
  *prenotato* (con la prenotazione e la sua fine) o *fuori servizio*. Una
  prenotazione compare a partire dal suo anticipo. Una manifestazione su più
  stand mostra una griglia per stand, e la finestra d’iscrizione raggruppa i
  bersagli per stand; un tiratore che passa a un altro stand viene concluso
  e iscritto di nuovo lì.
- Le cifre in alto contano le iscrizioni in attesa, al tiro e concluse.
- *In attesa* e *Concluso* mostrano 50 iscrizioni per pagina; le frecce
  sotto l’elenco sfogliano le pagine.

Tutto si aggiorna in diretta: un risultato registrato su un altro computer o
una prenotazione che inizia compare senza ricaricare. Il punto accanto a
**in diretta** nell’intestazione mostra la connessione; dopo un’interruzione
di rete la pagina si aggiorna da sola.

## Cercare e scansionare

Premi ++f4++ in un punto qualsiasi della dashboard per saltare al campo di
ricerca.

![Ricerca](../images/screenshots/search.png){ .screenshot }

- Digita un **numero di licenza** (solo cifre, zeri iniziali facoltativi),
  un **nome** (parte del nome o del cognome) o **nome e cognome** in
  qualsiasi ordine.
- Ogni risultato offre **Iscrivi &lt;società&gt;** per ogni società del
  tiratore, o **Apri iscrizione** se il tiratore è attualmente iscritto.
- **Nuovo tiratore** crea un tiratore e apre subito l’iscrizione.
- Uno scanner manuale che legge il [codice
  d’iscrizione](../reference/codes.md) su un foglio di stand o una striscia
  di controllo digita il codice e Invio: i risultati di quell’iscrizione si
  aprono direttamente.

## Terminare il tiro

**Termina** chiude un’iscrizione prima che tutti i risultati siano presenti, per esempio quando un tiratore rinuncia o consegna i bersagli più tardi. L’iscrizione passa a *Concluso* e il suo bersaglio si libera; i risultati mancanti si possono ancora registrare.

## Tema scuro

Il pulsante sole/luna nell’intestazione passa dal tema chiaro a quello
scuro; la scelta resta salvata nel browser.

![Dashboard con tema scuro](../images/screenshots/dashboard-dark.png){
.screenshot }
