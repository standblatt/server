# Manifestazioni

**Manifestazioni** elenca le manifestazioni della società nei gruppi *Oggi* (tra il primo e l’ultimo giorno di tiro), *Prossime*, *Passate* e *Chiuse*; gli admin ne creano di nuove qui. Aprendo una manifestazione la barra laterale mostra la dashboard e le altre pagine dello svolgimento; *Cambia manifestazione* in alto nella barra laterale riporta all’elenco.

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
| Pagine dei moduli | p. es. **Classifiche generali**, **Display chiosco**, **Link aiutante**, **Prenotazione online**, solo finché il modulo è attivo |
| **Copia e chiudi** | copiare per l’anno prossimo, chiudere, riaprire ed eliminare |

Il personale lavora solo nello svolgimento; la configurazione è riservata
agli admin.

## Impostazioni {#settings}

**Generale** contiene ciò che ogni manifestazione ha:

| Impostazione | Significato |
|---|---|
| Nome, data, stand | la data è il giorno principale; il suo anno fissa le categorie d’età (FST: anno della gara meno anno di nascita, il giorno in sé non conta, il che vale per il tiro anticipato di dicembre prima di una manifestazione di gennaio); spuntare ogni stand su cui si tira la manifestazione; o nessuno per una manifestazione senza bersagli (vedi sotto) |
| Giorni di tiro | singoli giorni con orari di tiro facoltativi, per il tiro anticipato e le manifestazioni di più giorni (le iscrizioni contano il giorno in cui vengono fatte); il giorno principale è sempre uno di essi |
| Periodo (stagione) | per manifestazioni di settimane o mesi, p. es. il Tiro in campagna dall’apertura dello stand fino al 31 agosto: primo e ultimo giorno, più i giorni della settimana con i loro orari di tiro; nessun giorno spuntato significa ogni giorno. I singoli giorni, p. es. il fine settimana ufficiale, si aggiungono al periodo o ne sostituiscono gli orari in quella data |
| Organizzatore | stampato sui documenti |
| Lingua dei documenti e del chiosco | lingua dei PDF, delle esportazioni CSV e degli schermi chiosco, indipendente dalla lingua d’interfaccia di ogni utente |

## Moduli {#modules}

Non ogni manifestazione ha bisogno di tutto: una festa di tiro prenota
bersagli e consegna carte corona, un programma annuale raccoglie solo
risultati. **Moduli** attiva e disattiva gli elementi della manifestazione;
ciò che è disattivato scompare dalla barra laterale, dalla dashboard,
dall’editor dei programmi, dalle classifiche e dai rapporti.

| Modulo | Cosa aggiunge |
|---|---|
| Prenotazioni | bersagli riservati a gruppi e società (richiede uno stand) |
| Prenotazione online | le società chiedono bersagli tramite un link pubblico (richiede uno stand e attiva Prenotazioni, vedi [Prenotazione online](online-booking.md)) |
| Foglio di stand, Striscia di controllo | i documenti stampati all’iscrizione |
| Distribuzione munizioni | il pulsante **Munizione** sulla dashboard e per la distribuzione munizioni (vedi [Aiutanti](helpers.md#ammunition)); lascialo disattivato dove i tiratori portano le proprie, p. es. pistola e carabina ad aria compressa |
| Tentativi | programmi che si possono sparare più volte (vedi [Programmi](programmes.md#attempts)) |
| Distinzioni | carte corona, premi in natura e premi in denaro; la sua impostazione *Stampa il valore della distinzione sulle carte corona* mette il valore in cifre e lettere sulla carta |
| Premi speciali | i migliori ranghi di ogni categoria d’età (attiva Distinzioni, vedi [Programmi](programmes.md#special-prizes)) |
| Concorso di sezione | le società classificate secondo i risultati dei loro tiratori; la sua impostazione *Società esclusa dal concorso di sezione* è di solito la società organizzatrice, elencata *fuori concorso* |
| Classifiche generali | classifiche su più programmi o tentativi (vedi [Classifiche generali](classifications.md)) |
| Display chiosco | classifiche su schermi nella sala (vedi [Classifiche e chiosco](rankings-kiosk.md)) |
| Link aiutante | marcatori, distribuzione munizioni e registrazione risultati senza account (vedi [Aiutanti](helpers.md)) |

Quando crei una manifestazione, scegli di che tipo è; questo attiva solo dei
moduli:

| Tipo | Moduli |
|---|---|
| Tiro di festa | Prenotazioni, Prenotazione online, Foglio di stand, Striscia di controllo, Distinzioni, Premi speciali, Concorso di sezione, Display chiosco, Link aiutante |
| Programma annuale o allenamento | Tentativi, Classifiche generali |
| Tiro in campagna | Distinzioni, Concorso di sezione, Display chiosco, Link aiutante |
| Campionato sociale | Tentativi, Classifiche generali, Display chiosco |
| Selezione personalizzata | nessuno |

Un modulo che contiene dati resta attivo: disattiva il chiosco solo dopo
averne eliminato i link chiosco, le distinzioni solo quando nessun programma
ha livelli di distinzione, e così via; la scheda indica perché il suo
interruttore è bloccato. Un modulo di cui un altro ha bisogno (Prenotazioni
per Prenotazione online) resta attivo finché lo è quest’ultimo.

Gli orari di tiro danno anche alla prenotazione online le sue fasce, a meno
che la prenotazione fissi altri orari per un giorno (vedi [prenotazione
online](online-booking.md)).

### Moduli della società {#club-modules}

Alcuni moduli appartengono alla società invece che a una
manifestazione. **Moduli della società** nella parte società della barra
laterale (admin) li attiva e disattiva:

| Modulo | Cosa aggiunge |
|---|---|
| Modelli di stampa | modelli per fogli prestampati e la sezione *Modello di stampa* di ogni programma (vedi [Modelli di stampa](print-forms.md)) |

## Manifestazioni senza bersagli {#no-lanes}

Una manifestazione non ha bisogno di uno stand, p. es. un programma annuale
o un Tiro in campagna i cui risultati arrivano su carta o dall’impianto di
segnalazione. Senza stand non c’è la griglia dei bersagli, e i moduli
Prenotazioni e Prenotazione online non sono disponibili; i tiratori vengono
iscritti senza bersaglio e i loro risultati si registrano come di
consueto. I documenti si aprono ancora in PDF, ma non vanno a un agente di
stampa, perché gli agenti di stampa appartengono a uno stand.

## Copia per l’anno prossimo {#copy}

In **Copia e chiudi** nella configurazione, **Copia per l’anno prossimo**
crea una nuova manifestazione con gli stessi stand, organizzatore, lingua,
impostazioni, moduli e tutti i programmi, compresi prezzi, distinzioni,
combinazioni e premi speciali, e le [classifiche
generali](classifications.md); i giorni di tiro e il periodo si spostano con
la data, gli orari di tiro settimanali restano sui loro giorni della
settimana. Iscrizioni, risultati e link chiosco non vengono copiati. La
finestra propone il nome con l’anno (e un numero iniziale come *42.*)
incrementato e la data un anno dopo; adatta entrambi se necessario.

![Copiare una manifestazione](../images/screenshots/copy-event.png){
.screenshot }

## Chiudere {#close}

Quando la manifestazione è finita e le classifiche sono controllate,
**chiudila** in **Copia e chiudi** nella configurazione (**Chiudi la
manifestazione**). La chiusura pubblica le classifiche: i tiratori con un
[diario di allenamento personale](../my-training/account.md#competitions)
trovano ora i loro risultati di questa manifestazione, con il rango nel
programma, e possono aggiungerli al diario. Tutto resta visibile (dashboard,
risultati, classifiche, rapporti, documenti), ma iscrizioni, risultati,
distinzioni, programmi e impostazioni non possono più essere modificati. Un
banner su ogni pagina della manifestazione lo segnala. **Riapri** la rende
di nuovo modificabile; i risultati che i tiratori hanno già aggiunto al
diario restano. Copiare una manifestazione chiusa per l’anno prossimo
funziona.

## Eliminare

**Elimina** (in **Copia e chiudi**) rimuove una manifestazione, ma solo finché nessuno è iscritto.
