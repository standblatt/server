# Stampa

Standblatt produce tutti i documenti in PDF:

| Documento | Formato | Quando |
|---|---|---|
| Foglio di stand | A5 | all’iscrizione, con codice QR e codici a barre dei programmi |
| Striscia di controllo | etichetta 90 × 29 mm | all’iscrizione, incollata sulla carta bersaglio |
| Foglio dei risultati | A5 | quando tutti i programmi hanno risultati |
| Carta corona | A6, stampata su carte prestampate | dopo la scelta di una carta corona come distinzione |
| Classifica | A4 | dalla pagina delle classifiche, solo in PDF (non inviata alla stampante dello stand) |

Ogni pulsante di stampa stampa sulla stampante dello stand se è configurata
e il suo agente di stampa è online; altrimenti il PDF si apre nel browser e
si può stampare da lì. Su telefono o tablet il PDF passa invece al menu di
condivisione (stampare, salvare in File, inviare), così l’app resta aperta;
se il telefono chiede prima un altro tocco, appare un avviso con il pulsante
**Condividi**.

## Configurare le stampanti

**Stand di tiro → Stampa** (admin):

![Stampa](../images/screenshots/printing.png){ .screenshot }

1. **Aggiungi agente di stampa**, dagli un nome e copia il token. Installa
   l’agente su un computer vicino alle stampanti (vedi [Agente di
   stampa](../admin-guide/print-agent.md)).
2. Appena l’agente è online, comunica le sue stampanti. Scegli per ogni
   documento la stampante e le sue opzioni (cassetto, formato carta, tipo di
   supporto, altre opzioni `lp` come `fit-to-page`).
3. **Lavori recenti** mostra cosa è stato stampato e se ha funzionato.

Un lavoro che l’agente non conferma entro due minuti viene ripetuto; dopo
tre tentativi fallisce. **revoca** rimuove subito un agente.
