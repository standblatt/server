# Gestire un programma annuale durante la stagione

Ti serve quando i membri tirano un programma quante volte vogliono durante
la stagione e un campionato sociale classifica i tre migliori risultati di
ogni tiratore. Richiede il permesso *Configurare le manifestazioni* (vedi
[Ruoli](../user-guide/users.md#roles)); per registrare i risultati serve
*Gestire le manifestazioni*.

## Creare la manifestazione {#event}

1. In **Manifestazioni** clicca **Nuova manifestazione** e scegli
   **Programma annuale o allenamento**. Attiva i moduli *Tentativi* e
   *Classifiche generali*; scegli invece **Campionato sociale** se la
   classifica deve apparire anche su un display chiosco.
2. Inserisci il **Nome**, p. es. *Programma annuale 2027*, e come **Data**
   il primo giorno della stagione. Il suo anno determina le categorie d’età.
3. Spunta lo stand in **Stand di tiro**, o nessuno se i risultati arrivano
   solo su carta (vedi [Manifestazioni senza
   bersagli](../user-guide/events.md#no-lanes)).
4. Attiva **Periodo (stagione)**, inserisci il **Primo giorno** e l’**Ultimo
   giorno**, e in **Orari di tiro** spunta i giorni della settimana
   dell’allenamento della società con i loro orari. Nessun giorno spuntato
   significa ogni giorno del periodo.
5. Clicca **Crea**.

## Aggiungere il programma {#programme}

1. Apri **Configurazione** → **Programmi** e clicca **Nuovo programma** (o
   **Aggiungi da …** per prenderlo dalla biblioteca, da un file o da una
   manifestazione precedente).
2. Inserisci il nome e le serie.
3. In **Tentativi** spunta **Illimitato** e imposta **Tentativo che conta**
   su **Il migliore**.
4. In **Prezzi (CHF)** spunta ogni categoria d’età che può tirare e
   inserisci il prezzo. Il secondo campo è il prezzo di ogni tentativo
   supplementare: vuoto significa lo stesso prezzo, 0 significa che il
   programma è pagato una volta.
5. Clicca **Crea**.

## Classificare i tre migliori tentativi {#classification}

1. Apri **Configurazione** → **Classifiche generali** e clicca **Nuova
   classifica generale**.
2. Inserisci un **Nome**, p. es. *Campionato sociale*, scegli il programma
   annuale come **Programma 1** e imposta **Tentativi che contano** su 3.
3. Attiva **Classificare solo i tiratori con tutti i risultati che contano**
   se solo i tiratori con tre risultati devono ricevere un rango; gli altri
   vengono allora elencati sotto la classifica.
4. Aggiungi **Classi** se juniori e adulti vengono classificati
   separatamente, p. es. *Juniori* con U17–U21 ed *Élite* con gli altri.
5. Clicca **Crea**.

Tutte le impostazioni dell’editor sono descritte in [Classifiche
generali](../user-guide/classifications.md).

## Registrare i risultati durante la stagione {#results}

1. Apri la manifestazione e clicca **Registra risultati** nella barra
   laterale.
2. Scegli il **Programma**, digita un numero di licenza o un nome in
   **Tiratore** e premi ++enter++.
3. Digita le serie e premi ++enter++. Ogni risultato di un tiratore è un
   nuovo tentativo.

I risultati del sistema di punteggio si possono anche importare come file
(vedi [Importare i risultati dal sistema di punteggio](import-results.md)).

**Classifiche** mostra il programma con il miglior tentativo di ogni tiratore (*2/5*: conta il secondo di cinque tentativi) e il campionato sociale con la somma dei tre migliori.

## Concludere la stagione {#close}

Quando gli ultimi risultati sono arrivati e controllati, apri
**Configurazione** → **Copia e chiudi** e clicca **Chiudi la
manifestazione**. Cosa cambia la chiusura è spiegato in [Chiudere una
manifestazione](../explanation/closing.md); la stagione successiva inizia
con [Riprendere la manifestazione dell’anno scorso](next-year.md).
