# Stampa

Standblatt produce tutti i documenti in PDF:

| Documento | Formato | Quando |
|---|---|---|
| Foglio di stand | A5 | all’iscrizione, con codice QR e codici a barre dei programmi |
| Striscia di controllo | etichetta 90 × 29 mm | all’iscrizione, incollata sulla carta bersaglio |
| Foglio dei risultati | A5 | quando tutti i programmi hanno risultati |
| Carta corona | A6, stampata su carte prestampate | dopo la scelta di una carta corona come distinzione |
| Modello di stampa | qualsiasi formato, stampato su fogli prestampati | per i programmi con un [modello di stampa](print-forms.md), appena ci sono i risultati |
| Classifica | A4 | dalla pagina delle classifiche, solo in PDF (non inviata alla stampante dello stand) |

Ogni pulsante di stampa stampa sulla stampante dello stand se è configurata
e il suo agente di stampa è online; altrimenti il PDF si apre nel browser e
si può stampare da lì. Su telefono o tablet il PDF passa invece al menu di
condivisione (stampare, salvare in File, inviare), così l’app resta aperta;
se il telefono chiede prima un altro tocco, appare un avviso con il pulsante
**Condividi**.

## Configurare le stampanti

**Stand di tiro → Configura le stampanti** (admin):

![Stampa](../images/screenshots/printing.png){ .screenshot }

1. **Aggiungi agente di stampa**, dagli un nome e copia il token. Installa
   l’agente su un computer vicino alle stampanti (vedi [Agente di
   stampa](../admin-guide/print-agent.md)). Un computer può servire più
   stand, p. es. *Pistola 25 m* e *Pistola 50 m*: aggiungi un agente su ogni
   stand e dai all’agente entrambi i token, separati da una virgola.
2. Appena l’agente è online, comunica le sue stampanti e l’elenco mostra la
   sua versione e il sistema operativo. Scegli per ogni documento la
   stampante (*Diplomi* sono i diplomi delle classifiche, *Modelli di
   stampa* tutti gli altri modelli della società) e le sue opzioni: quelle
   che l’agente comunica per quella stampante (su Linux cassetto, formato
   carta e tipo di supporto, più altre opzioni `lp` come `fit-to-page` nel
   campo di testo; su Windows cassetto, carta, fronte retro,
   ridimensionamento e colore).
3. **Pagina di prova** stampa una piccola pagina (90 × 29 mm, le dimensioni
   della striscia di controllo) sulla stampante di quella riga con le
   opzioni scelte, anche prima di salvare: società, stand di tiro,
   stampante, opzioni e ora. Una stampante di etichette la stampa come
   etichetta, le altre stampanti in alto a sinistra sul foglio. Serve un
   agente online.
4. **Lavori recenti** mostra cosa è stato stampato e se ha funzionato.

Un lavoro che l’agente non conferma entro due minuti viene ripetuto; dopo
tre tentativi fallisce. **revoca** rimuove subito un agente.

## Collegare un computer con un codice {#connect}

L’app Standblatt su un computer vicino alle stampanti non ha bisogno di un
token: mostra un codice come `K7QD-4MZP` e apre **Collegare un agente di
stampa** nel browser. Accedi se richiesto, controlla il computer indicato,
scegli società, stand di tiro e nome dell’agente e **Collega**. Solo i
membri che possono gestire gli stand di tiro possono collegare un computer;
un codice vale dieci minuti. L’app si collega poi da sola e l’agente compare
in **Stand di tiro → Configura le stampanti** con la sua versione e il sistema operativo.
