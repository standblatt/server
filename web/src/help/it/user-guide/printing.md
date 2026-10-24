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

1. Collega un computer vicino alle stampanti: un PC Windows o un Mac con
   l’[app Standblatt Print Agent](#app) tramite [codice](#connect), oppure
   un computer Linux come un Raspberry Pi con l’[agente senza
   interfaccia](https://github.com/standblatt/print-agent): **Aggiungi
   agente di stampa**, dagli un nome e copia il token. Un computer può
   servire più stand, p. es. *Pistola 25 m* e *Pistola 50 m*: collegalo a
   ogni stand (l’app con **Collega un altro stand di tiro**, l’agente senza
   interfaccia con entrambi i token, separati da una virgola).
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

## L’app Standblatt Print Agent {#app}

L’app è la via semplice per un portatile o un PC vicino alle stampanti:
nessun token e nessuna impostazione da digitare.

1. Scaricala dall’[ultima
   versione](https://github.com/standblatt/print-agent/releases/latest):
   `Standblatt-Print-Agent-Setup.exe` per Windows,
   `standblatt-print-agent-app-darwin-arm64.dmg` per un Mac con chip Apple
   (`-amd64.dmg` per un Mac Intel), il `.deb` o il `.rpm` per Linux.
2. Installala. L’app non è ancora firmata: Windows chiede una volta,
   **Ulteriori informazioni → Esegui comunque**; su un Mac trascinala in
   *Applicazioni* e aprila la prima volta con un clic destro → **Apri**. Su
   Windows si installa per l’utente attuale senza diritti di amministratore,
   con un collegamento nel menu Start e sul desktop, e si avvia con Windows.
3. Fai clic su **Collega uno stand di tiro** e conferma il codice in
   Standblatt (vedi [sotto](#connect)).

![L’app Standblatt Print Agent](../images/app/print-agent.png){ .screenshot
}

La finestra mostra per ogni stand se stampa, con le ultime stampe, e sotto
le stampanti del computer, ognuna con una **Pagina di prova** che stampa su
quel computer senza Standblatt. **Avanzate** contiene *Avvia con il
computer*, l’indirizzo di Standblatt (da cambiare solo per un server di
prova) e la rimozione di uno stand. Se chiudi la finestra, l’app continua
nell’area di notifica: l’icona è verde quando ogni stand stampa, altrimenti
rossa.

Su Windows e su un Mac l’app segnala quando c’è una nuova versione;
**Aggiorna ora** la installa e riavvia l’app. Per dismettere un computer,
rimuovi il suo stand in **Avanzate**, revoca il suo agente in **Stand di
tiro → Configura le stampanti** e disinstalla l’app (su Windows in
*Impostazioni → App*).

## Collegare un computer con un codice {#connect}

![Un codice nell’app](../images/app/print-agent-code.png){ .screenshot }

L’app Standblatt su un computer vicino alle stampanti non ha bisogno di un
token: mostra un codice come `K7QD-4MZP` e apre **Collegare un agente di
stampa** nel browser. Accedi se richiesto, controlla il computer indicato,
scegli società, stand di tiro e nome dell’agente e **Collega**. Solo i
membri che possono gestire gli stand di tiro possono collegare un computer;
un codice vale dieci minuti. L’app si collega poi da sola e l’agente compare
in **Stand di tiro → Configura le stampanti** con la sua versione e il
sistema operativo.

Se il browser di quel computer non è a portata di mano, inserisci invece il
codice in **Stand di tiro → Configura le stampanti** dello stand, in
**Collega con codice**: il computer viene collegato a quello stand con il
suo nome.
