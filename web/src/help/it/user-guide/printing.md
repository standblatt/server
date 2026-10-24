# Stampa

Standblatt genera tutti i documenti come PDF. Ogni documento con il suo
formato, contenuto e impostazione della stampante:
[Documenti](../reference/documents.md#documents).

Ogni pulsante di stampa stampa sulla stampante dello stand se ne è
configurata una e il suo agente di stampa è online; altrimenti si apre il
PDF (vedi [Quando non c’è una
stampante](../explanation/printing.md#fallback)).

Perché un agente di stampa e come funziona: [Come funziona la
stampa](../explanation/printing.md).

## Configurare le stampanti

**Stand di tiro → Configura le stampanti** (permesso *Stand e stampanti*) ha tre parti:

![Stampa](../images/screenshots/printing.png){ .screenshot }

- **Agenti di stampa**: i computer che stampano per questo stand, ognuno con
  la sua versione e il sistema operativo, *collegato* o quando è stato visto
  l’ultima volta, e il suo numero di stampanti. **Aggiungi agente di
  stampa** aggiunge un agente senza interfaccia e mostra una volta il suo
  token, **Collega con codice** collega l’[app](#app) tramite il suo
  [codice](#connect), **revoca** rimuove subito un agente.
- **Documenti**: per ogni documento la stampante (*browser (PDF)* continua
  ad aprire il PDF) e le opzioni che l’agente riporta per quella stampante;
  su Linux altre opzioni `lp` come `fit-to-page` vanno nel campo di
  testo. *Diplomi* sono i diplomi dalle classifiche, *Modelli di stampa
  (prestampati)* ogni altro modello della società. **Salva** salva la
  riga. **Pagina di prova** stampa una piccola pagina (90 × 29 mm) con
  società, stand, stampante, opzioni e ora sulla stampante di quella riga,
  anche prima di salvare; serve un agente online.
- **Lavori recenti** mostra cosa è stato stampato e se ha funzionato.

Un lavoro che l’agente non conferma entro due minuti viene ripetuto; dopo
tre tentativi fallisce.

Passo per passo: [Collegare una stampante allo
stand](../how-to/connect-printer.md).

## L’app Standblatt Print Agent {#app}

L’app è la via semplice per un portatile o un PC vicino alle stampanti:
nessun token e nessuna impostazione da digitare.

Scaricarla, installarla e collegarla: [Collegare una stampante allo
stand](../how-to/connect-printer.md#app).

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

**Collegare un agente di stampa** (aperto dall’app con un codice come `K7QD-4MZP`) mostra il computer e chiede società, stand e nome; solo i membri che possono gestire gli stand possono collegare un computer, e un codice è valido dieci minuti. **Collega con codice** in **Stand di tiro → Configura le stampanti** lo collega a quello stand con il suo nome. Passi: [Collegare una stampante allo stand](../how-to/connect-printer.md#app).
