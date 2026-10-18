# Manifestazioni

**Manifestazioni** elenca le manifestazioni della società nei gruppi *Oggi* (tra il primo e l’ultimo giorno di tiro), *Prossime*, *Passate* e *Chiuse*; gli admin ne creano di nuove qui. Aprendo una manifestazione la barra laterale mostra dashboard, prenotazioni, classifiche, statistiche e configurazione; *Cambia manifestazione* in alto nella barra laterale riporta all’elenco.

![Manifestazioni](../images/screenshots/events.png){ .screenshot }

## Impostazioni {#settings}

La **Configurazione** della manifestazione ha cinque schede: *Programmi*
(vedi [Programmi](programmes.md)), *Manifestazione*, *Display chiosco* (vedi
[Classifiche e chiosco](rankings-kiosk.md)), *Aiutanti* (vedi
[Aiutanti](helpers.md)) e *Copia e chiudi*. I collaboratori vedono solo i
programmi.

La scheda *Manifestazione* contiene le impostazioni (admin):

| Impostazione | Significato |
|---|---|
| Nome, data, stand | la data è il giorno principale; il suo anno fissa le categorie d’età (FST: anno della gara meno anno di nascita, il giorno in sé non conta, il che vale per il tiro anticipato di dicembre prima di una manifestazione di gennaio); spuntare ogni stand su cui si tira la manifestazione; o nessuno per una manifestazione senza bersagli (vedi sotto) |
| Giorni di tiro | singoli giorni con orari di tiro facoltativi, per il tiro anticipato e le manifestazioni di più giorni (le iscrizioni contano il giorno in cui vengono fatte); il giorno principale è sempre uno di essi |
| Periodo (stagione) | per manifestazioni di settimane o mesi, p. es. il Tiro in campagna dall’apertura dello stand fino al 31 agosto: primo e ultimo giorno, più i giorni della settimana con i loro orari di tiro; nessun giorno spuntato significa ogni giorno. I singoli giorni, p. es. il fine settimana ufficiale, si aggiungono al periodo o ne sostituiscono gli orari in quella data |
| Organizzatore | stampato sui documenti |
| Lingua dei documenti e del chiosco | lingua dei PDF, delle esportazioni CSV e degli schermi chiosco, indipendente dalla lingua d’interfaccia di ogni utente |
| Società esclusa dal concorso di sezione | di solito la società organizzatrice; elencata come *fuori concorso* |
| Documenti all’iscrizione | **Foglio di stand** e **Striscia di controllo**: disattivare ciò che la manifestazione non usa; l’iscrizione allora non lo propone né lo stampa |
| Mostra lo stato dei bersagli sulle classifiche | gli schermi chiosco mostrano i bersagli |
| Stampa il valore della distinzione sulle carte corona | valore in cifre e in lettere sulla carta |
| L’organizzatore consegna le munizioni | mostra il pulsante **Munizione** nella dashboard e per la distribuzione munizioni (vedi [Aiutanti](helpers.md)); disattivato per le manifestazioni in cui i tiratori portano le proprie, ad es. pistola e carabina ad aria compressa |

Gli orari di tiro danno anche alla prenotazione online le sue fasce, a meno
che la prenotazione fissi altri orari per un giorno (vedi [prenotazione
online](online-booking.md)).

## Manifestazioni senza bersagli {#no-lanes}

Una manifestazione non ha bisogno di uno stand, p. es. un programma annuale
o un Tiro in campagna i cui risultati arrivano su carta o dall’impianto di
segnalazione. Senza stand non ci sono griglia dei bersagli, prenotazioni né
prenotazione online; i tiratori vengono iscritti senza bersaglio e i loro
risultati si registrano come di consueto. I documenti si aprono ancora in
PDF, ma non vanno a un agente di stampa, perché gli agenti di stampa
appartengono a uno stand.

## Copia per l’anno prossimo

Nella scheda *Copia e chiudi*, **Copia per l’anno prossimo** crea una nuova
manifestazione con gli stessi stand, organizzatore, lingua, impostazioni e
tutti i programmi, compresi prezzi, distinzioni, combinazioni e premi
speciali, e le [classifiche generali](classifications.md); i giorni di tiro
e il periodo si spostano con la data, gli orari di tiro settimanali restano
sui loro giorni della settimana. Iscrizioni, risultati e link chiosco non
vengono copiati. La finestra propone il nome con l’anno (e un numero
iniziale come *42.*) incrementato e la data un anno dopo; adatta entrambi se
necessario.

![Copiare una manifestazione](../images/screenshots/copy-event.png){
.screenshot }

## Chiudere {#close}

Quando la manifestazione è finita e le classifiche sono controllate,
**chiudila** nella scheda *Copia e chiudi* (**Chiudi la
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

**Elimina** (nella scheda *Copia e chiudi*) rimuove una manifestazione, ma solo finché nessuno è iscritto.
