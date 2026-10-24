# Documenti e rapporti

Ogni documento è un PDF nella lingua della manifestazione. Una stampa va
alla stampante impostata per il documento in **Stand di tiro → Configura le
stampanti** dello stand dell’iscrizione quando un agente di stampa dello
stand è online; altrimenti si apre il PDF (vedi
[Stampa](../user-guide/printing.md)).

## Documenti {#documents}

| Documento | Formato | Contenuto | Stampato da | Impostazione della stampante | Modulo |
|---|---|---|---|---|---|
| Foglio di stand | A5 | manifestazione, tiratore, [codice d’iscrizione](codes.md) QR, stand (manifestazioni su più stand), bersaglio, programmi con colpi e [codici a barre](codes.md) | *Stampa foglio di stand* all’iscrizione; l’iscrizione | *Foglio di stand (A5)* | Foglio di stand |
| Striscia di controllo | etichetta 90 × 29 mm, o un [modello di stampa](../user-guide/print-forms.md#builtin) della società | manifestazione, tiratore, licenza, società, categoria d’età, serie, codice d’iscrizione QR | *Stampa striscia di controllo* all’iscrizione; l’iscrizione | *Striscia di controllo (etichetta 90 × 29 mm)* | Striscia di controllo |
| Foglio dei risultati | A5 | manifestazione, tiratore, ogni programma con le sue serie, totale e distinzione | **Stampa foglio dei risultati** nei risultati, appena ogni programma ha risultati | *Foglio dei risultati (A5)* | – |
| Carta corona | A6, stampata su carte prestampate, o un modello di stampa della società | tiratore, anno di nascita, società, data, organizzatore, valore in cifre e in lettere (impostazione *Stampa il valore della distinzione sulle carte corona*) | **Stampa carta corona** nei risultati, per una distinzione carta corona | *Carta corona (A6, prestampata)* | Distinzioni |
| Modello di stampa di un programma | misura del modello | vedi [Segnaposto dei modelli](form-placeholders.md#results) | **Stampa sul modello** nei risultati di un programma concluso | *Modelli di stampa (prestampati)* | modulo della società Modelli di stampa |
| Documento d’iscrizione | misura del modello | vedi [Segnaposto dei modelli](form-placeholders.md#registration) | un pulsante con il nome del modello nell’iscrizione | *Modelli di stampa (prestampati)* | modulo della società Modelli di stampa |
| Diplomi | misura del modello, una pagina per tiratore classificato | vedi [Segnaposto dei modelli](form-placeholders.md#rank) | **Diplomi** in **Classifiche** | *Diplomi* | modulo della società Modelli di stampa |
| Classifica | A4 | rango, nome, anno, categoria, società, serie, totale; ranghi 1–3 evidenziati; *(2/3)* per il tentativo che conta | **PDF** in **Classifiche** e **Statistiche** | mai stampata, solo PDF | – |
| Pagina di prova | etichetta 90 × 29 mm | società, stand, stampante, opzioni, ora | **Pagina di prova** in **Stand di tiro → Configura le stampanti** | la stampante di quella riga | – |

- Senza stand, i documenti si aprono solo come PDF.
- I link aiutante stampano solo il foglio di stand (distribuzione munizioni)
  o i modelli di stampa (marcatore, registrazione risultati), e solo dei
  loro bersagli.
- Un lavoro di stampa che l’agente non conferma entro due minuti viene
  ripetuto, al massimo tre volte.

## Rapporti {#reports}

In **Statistiche** (permesso *Statistiche e rapporti*). Ogni rapporto è
disponibile in **PDF** e **CSV** (UTF-8 con BOM, separato da `;`, si apre
direttamente in Excel). I testi seguono la lingua della manifestazione.

| Rapporto | Contenuto | Modulo |
|---|---|---|
| Rapporto finanziario | incassi per programma, distinzioni con il loro valore | – |
| Incassi per programma e categoria d’età | partecipanti, prezzo e subtotale per programma e categoria | – |
| Partecipanti per società | partecipanti per società, programma e categoria d’età | – |
| Distinzioni | ogni distinzione assegnata, con tiratore, società, programma, valore e risultato | Distinzioni |
| Distinzioni da consegnare | distinzioni non ancora ritirate, con tiratore, società, programma, tipo e valore, più il denaro ancora da versare (vedi [Risultati](../user-guide/results.md#handout)) | Distinzioni |
| Concorso di sezione | classifica, società fuori concorso e non classificate | Concorso di sezione |
| Premi speciali | vincitori per programma e categoria d’età, tiratori saltati con il motivo (vedi [Programmi](../user-guide/programmes.md#special-prizes)) | Premi speciali |
| Classifiche generali | ogni [classifica generale](../user-guide/classifications.md) per classe, con i punti per programma; le parti che non contano sono tra parentesi | Classifiche generali |
| Prenotazioni di bersagli | le società che hanno richiesto bersagli tramite la [prenotazione online](../user-guide/online-booking.md), confermate o in attesa, con contatto, fasce orarie, numeri dei bersagli e tiratori annunciati | Prenotazioni |

Il PDF di un rapporto è A4, *Distinzioni* in orizzontale. La classifica di
ogni programma è anche un file CSV: rango, cognome, nome, anno, categoria,
società, serie, totale e mouche, per i programmi con più tentativi anche il
tentativo che conta e il numero di tentativi.
