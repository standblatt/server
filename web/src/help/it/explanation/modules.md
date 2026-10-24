# Moduli e tipi di manifestazione

Le società organizzano manifestazioni molto diverse: una festa di tiro con
accoglienza, assegnazione dei bersagli e carte corona, un Tiro in campagna
nel fine settimana ufficiale, un programma annuale tirato durante la
stagione, un campionato sociale. Standblatt non ha codice separato per
ognuna di esse. Ogni manifestazione è lo stesso tipo di cosa; le
manifestazioni si distinguono solo per i loro programmi, i loro stand e
giorni e i **moduli** che attivano.

## Il nucleo e i moduli {#core}

Ogni manifestazione ha lo stesso nucleo: nome, giorni, stand e lingua,
programmi con le loro serie e prezzi, iscrizioni, risultati, classifiche,
statistiche e chiusura. Tutto il resto è un modulo: assegnazione dei
bersagli, prenotazioni con la prenotazione online, il foglio di stand e la
striscia di controllo, la distribuzione munizioni, i tentativi, le
distinzioni, i premi speciali, il concorso di sezione, le classifiche
generali, i display chiosco e i link aiutante.

Un modulo disattivato sparisce, non viene solo nascosto: scompare dalla
barra laterale, dalla dashboard, dall’editor dei programmi, dalle
classifiche e dai rapporti, e anche il server rifiuta le sue funzioni, così
né un link aiutante né un token API possono usarlo. Un programma annuale
senza distinzioni non ha quindi una sezione distinzioni nei suoi programmi
né rapporti sulle distinzioni; nessuno deve saltare ciò che la
manifestazione non usa.

Alcuni moduli ne richiedono altri: *Premi speciali* attiva *Distinzioni*,
perché un premio speciale è una distinzione, e *Prenotazioni* richiede
*Assegnazione dei bersagli*, perché una prenotazione blocca dei bersagli. Un
modulo resta attivo finché è attivo un modulo che lo richiede.

## I tipi di manifestazione sono solo un punto di partenza {#kinds}

**Nuova manifestazione** chiede di che tipo di manifestazione si tratta: tiro di festa, programma annuale o allenamento, Tiro in campagna, campionato sociale o una selezione personalizzata. Il tipo attiva solo un insieme di moduli (l’elenco è in [Tipi di manifestazione](../reference/modules.md#kinds)); la manifestazione non lo ricorda. In seguito i moduli si possono cambiare come quelli di qualsiasi altra manifestazione, così un Tiro in campagna che vuole anche classifiche generali le attiva, e una festa di tiro senza accoglienza disattiva l’assegnazione dei bersagli.

Per questo un nuovo tipo di manifestazione non richiede una nuova funzione:
se una combinazione di moduli lo copre, si può organizzare già oggi. Una
società che tiene ogni anno una manifestazione nella stessa forma la salva
invece come [modello](../user-guide/events.md#templates); un modello è una
manifestazione con moduli, programmi e classifiche generali, non un tipo.

## Perché un modulo che contiene dati resta attivo {#in-use}

Disattivare un modulo che contiene dati nasconderebbe dati che contano
ancora: livelli di distinzione raggiunti dai tiratori, link chiosco che
mostrano ancora le classifiche nella sala, programmi che permettono più
tentativi. Standblatt tiene quindi attivo un tale modulo finché i suoi dati
non sono spariti, e la scheda del modulo dice perché il suo interruttore è
bloccato:

| Modulo | Resta attivo finché |
|---|---|
| Assegnazione dei bersagli | un’iscrizione ha un bersaglio |
| Prenotazioni | il link della prenotazione online è attivo |
| Tentativi | un programma permette più di un tentativo |
| Distinzioni | un programma ha livelli di distinzione |
| Premi speciali | un programma assegna premi speciali |
| Classifiche generali | la manifestazione ha una classifica generale |
| Display chiosco | la manifestazione ha un link chiosco |
| Link aiutante | la manifestazione ha un link aiutante |

Il foglio di stand, la striscia di controllo, la distribuzione munizioni e
il concorso di sezione non contengono dati propri e si possono attivare o
disattivare in qualsiasi momento.

## Configurazione e svolgimento {#configuration}

Una manifestazione si prepara una volta e si svolge molte volte: programmi,
prezzi e moduli vengono configurati prima che arrivi il primo tiratore,
mentre l’accoglienza lavora il giorno stesso. Standblatt tiene separate le
due cose. Lo **svolgimento** (dashboard, risultati, classifiche, statistiche
e le pagine di moduli come le prenotazioni) richiede il permesso *Gestire le
manifestazioni*; la **configurazione** (impostazioni generali, moduli,
programmi, le pagine di configurazione dei moduli, copia e chiusura)
richiede *Configurare le manifestazioni*. Un aiutante all’accoglienza non
può cambiare un prezzo per sbaglio, e la configurazione non intralcia
l’accoglienza.

## Moduli della società e moduli della manifestazione {#club-modules}

La maggior parte dei moduli appartiene a una manifestazione, perché ogni
manifestazione decide se consegna distinzioni o usa schermi chiosco. Alcuni
elementi costitutivi sono usati da ogni manifestazione di una società e
appartengono quindi alla società: i **modelli di stampa** con la loro
biblioteca di immagini vengono disegnati una volta e usati dai programmi di
ogni manifestazione. Tali **moduli della società** si attivano in [Moduli
della società](../user-guide/events.md#club-modules). Una nuova società ha i
modelli di stampa attivati; come i moduli della manifestazione, il modulo
resta attivo finché la società ha modelli o immagini.

## Manifestazioni con e senza bersagli {#lanes}

I bersagli contano solo dove un’accoglienza li assegna. Standblatt separa
quindi lo stand dai bersagli:

| La manifestazione ha | Cosa ne consegue |
|---|---|
| uno stand e *Assegnazione dei bersagli* | l’accoglienza di una festa di tiro: griglia dei bersagli, prenotazioni, prenotazione online, link per marcatori e stato dei bersagli sugli schermi chiosco |
| uno stand, ma nessuna *Assegnazione dei bersagli* | i tiratori vengono iscritti senza bersaglio; i documenti vanno comunque all’agente di stampa dello stand e i programmi possono comunque essere legati allo stand |
| nessuno stand | nessun bersaglio e nulla che ne abbia bisogno; i documenti si aprono come PDF, perché gli agenti di stampa appartengono a uno stand |

Un programma annuale o un Tiro in campagna i cui risultati arrivano su carta
o dal sistema di punteggio non ha bisogno di bersagli, ma la società può
comunque volere la sua stampante allo stand; è la riga di mezzo. I moduli
che richiedono bersagli si disattivano da soli quando la manifestazione non
ha uno stand.

I risultati arrivano a una manifestazione senza bersagli come a qualsiasi
altra: tramite **Registra risultati**, un file CSV o i risultati di
un’iscrizione (vedi [Registra risultati](../user-guide/results-entry.md)).

## La copia mantiene la forma {#copy}

**Copia per l’anno prossimo** e i modelli riprendono i moduli attivati e quelli che contengono dati, insieme ai programmi e alle classifiche generali. Una manifestazione copiata dall’anno scorso parte quindi con la stessa forma, senza le sue iscrizioni e i suoi risultati.
