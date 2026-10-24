# Manifestazioni

**Manifestazioni** elenca le manifestazioni della società nei gruppi *Oggi* (tra il primo e l’ultimo giorno di tiro), *Prossime*, *Passate* e *Chiuse*; i ruoli con il diritto *Configurare le manifestazioni* ne creano di nuove. Aprendo una manifestazione la barra laterale mostra il suo cruscotto e le altre pagine dell’esercizio; *Cambia manifestazione* in alto nella barra laterale riporta all’elenco.

![Manifestazioni](../images/screenshots/events.png){ .screenshot }

## Configurazione e svolgimento {#configuration}

Una manifestazione ha due parti. Lo **svolgimento** è ciò con cui lavori il
giorno stesso: *Dashboard*, *Registra risultati*, *Classifiche*,
*Statistiche* e le pagine dei moduli attivi come
*Prenotazioni*. **Configurazione** in fondo alla barra laterale (solo admin)
apre l’impostazione della manifestazione; la barra laterale elenca allora le
sue pagine, e *Torna alla manifestazione* riporta allo svolgimento:

| Pagina | Contenuto |
|---|---|
| **Generale** | nome, giorni, stand, organizzatore e lingua (vedi sotto) |
| **Moduli** | gli elementi che la manifestazione usa (vedi sotto) |
| **Programmi** | cosa si può sparare (vedi [Programmi](programmes.md)) |
| Pagine dei moduli | p. es. **Prenotazioni** (con la prenotazione online), **Classifiche generali**, **Display chiosco**, **Link aiutante**, solo finché il modulo è attivo |
| **Copia e chiudi** | copiare per l’anno prossimo, chiudere, riaprire ed eliminare |

La configurazione richiede il permesso *Configurare le manifestazioni*, lo
svolgimento *Gestire le manifestazioni* (vedi
[Ruoli](users.md#roles)). Perché le due cose sono separate: [Configurazione
e svolgimento](../explanation/modules.md#configuration).

## Impostazioni {#settings}

**Generale** contiene ciò che ogni manifestazione ha:

| Impostazione | Significato |
|---|---|
| Nome, data, stand | la data è il giorno principale; il suo anno fissa le categorie d’età (FST: anno della gara meno anno di nascita, il giorno in sé non conta, il che vale per il tiro anticipato di dicembre prima di una manifestazione di gennaio); spuntare ogni stand su cui si tira la manifestazione; o nessuno per una manifestazione senza bersagli (vedi sotto) |
| Giorni di tiro | singoli giorni con orari di tiro facoltativi, per il tiro anticipato e le manifestazioni di più giorni (le iscrizioni contano il giorno in cui vengono fatte); il giorno principale è sempre uno di essi |
| Periodo (stagione) | per manifestazioni di settimane o mesi, p. es. il Tiro in campagna dall’apertura dello stand fino al 31 agosto: primo e ultimo giorno, più i giorni della settimana con i loro orari di tiro; nessun giorno spuntato significa ogni giorno. I singoli giorni, p. es. il fine settimana ufficiale, si aggiungono al periodo o ne sostituiscono gli orari in quella data |
| Organizzatore | stampato sui documenti |
| Lingua dei documenti e del chiosco | lingua dei PDF, delle esportazioni CSV e degli schermi chiosco, indipendente dalla lingua d’interfaccia di ogni utente |

Gli orari di tiro danno anche alla prenotazione online le sue fasce, a meno
che la prenotazione fissi altri orari per un giorno (vedi [prenotazione
online](online-booking.md)).

Passo per passo per una stagione: [Gestire un programma annuale durante la
stagione](../how-to/annual-programme.md), [Organizzare un Tiro in campagna
su un periodo](../how-to/feldschiessen.md).

## Moduli {#modules}

**Moduli** attiva e disattiva gli elementi costitutivi della manifestazione; perché le manifestazioni si distinguono solo per i loro moduli è spiegato in [Moduli e tipi di manifestazione](../explanation/modules.md).

![I moduli di una manifestazione](../images/screenshots/modules.png){
.screenshot }

Ogni modulo è una scheda con il suo interruttore; le impostazioni di un
modulo, p. es. *Stampa il valore della distinzione sulle carte corona* di
*Distinzioni*, si trovano sulla sua scheda. Cosa aggiunge ogni modulo, le
sue impostazioni e quando resta attivo:
[Moduli](../reference/modules.md#event-modules).

Quando crei una manifestazione, scegli di che tipo è; questo attiva solo dei
moduli (vedi [Tipi di manifestazione](../reference/modules.md#kinds)).

Un modulo che contiene dati resta attivo, e la scheda dice perché il suo
interruttore è bloccato (vedi [Perché un modulo che contiene dati resta
attivo](../explanation/modules.md#in-use)). Un modulo richiesto da un altro,
p. es. *Distinzioni* per *Premi speciali* o *Assegnazione dei bersagli* per
*Prenotazioni*, resta attivo finché lo è l’altro.

### Moduli della società {#club-modules}

Alcuni moduli appartengono alla società invece che a una
manifestazione. **Moduli della società** nella parte società della barra
laterale (permesso *Utenti e società*) li attiva e disattiva; per ora si
tratta di *Modelli di stampa* (vedi [Modelli di
stampa](print-forms.md)). Cosa aggiunge ognuno e quando resta attivo:
[Moduli della società](../reference/modules.md#club-modules); perché
appartengono alla società: [Moduli della società e moduli della
manifestazione](../explanation/modules.md#club-modules).

## Manifestazioni senza bersagli {#no-lanes}

Una manifestazione senza bersagli non ha stand, o disattiva *Assegnazione
dei bersagli*; cosa significa ogni scelta è spiegato in [Manifestazioni con
e senza bersagli](../explanation/modules.md#lanes).

## Copia per l’anno prossimo {#copy}

In **Copia e chiudi** nella configurazione, **copia** sotto *Copia per
l’anno prossimo* crea una nuova manifestazione con gli stessi stand,
organizzatore, lingua, impostazioni, moduli e tutti i programmi, compresi
prezzi, distinzioni, combinazioni e premi speciali, e le [classifiche
generali](classifications.md); i giorni di tiro e il periodo si spostano con
la data, gli orari di tiro settimanali restano sui loro giorni della
settimana. Iscrizioni, risultati, link chiosco, link aiutante e il link di
prenotazione non vengono copiati. Il dialogo *Copia manifestazione* propone
il nome e la data un anno dopo; **Crea copia** crea la manifestazione. Passo
per passo: [Riprendere la manifestazione dell’anno
scorso](../how-to/next-year.md#copy).

![Copiare una manifestazione](../images/screenshots/copy-event.png){
.screenshot }

## Modelli {#templates}

Una manifestazione ben configurata può diventare un **modello** della
società: **Salva come modello** in **Copia e chiudi** ne crea una copia con
i suoi moduli, impostazioni, stand, programmi e classifiche generali, con il
nome della manifestazione senza l’anno. I modelli sono elencati in *Modelli*
sotto le manifestazioni; un modello si apre nella sua configurazione e si
modifica come una manifestazione, ma non accetta iscrizioni. **Nuova
manifestazione** elenca i modelli della società sotto i tipi (vedi [Partire
da un modello della società](../how-to/next-year.md#template)). **Elimina**
in **Copia e chiudi** del modello lo rimuove.

![Nuova manifestazione con i tipi e un modello della
società](../images/screenshots/new-event.png){ .screenshot }

## Chiudere {#close}

Quando una manifestazione è finita e le classifiche sono controllate,
**Chiudi la manifestazione** in **Copia e chiudi** nella configurazione la
chiude; un banner su ogni pagina della manifestazione lo indica, e
**Riapri** la rende di nuovo modificabile. Cosa significa la chiusura per le
classifiche, Il mio allenamento e le modifiche: [Chiudere una
manifestazione](../explanation/closing.md).

## Eliminare

**Elimina** (in **Copia e chiudi**) rimuove una manifestazione, ma solo finché nessuno è iscritto. Una manifestazione chiusa non può essere eliminata; riaprila prima.
