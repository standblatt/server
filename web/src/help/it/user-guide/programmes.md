# Programmi

I programmi appartengono a una manifestazione e si gestiscono nella scheda
*Programmi* della sua **Configurazione** (solo admin).

![Configurazione con l’elenco dei
programmi](../images/screenshots/setup.png){ .screenshot }

## Aggiungere programmi

- **Nuovo programma** apre l’editor.
- **Aggiungi da … → Biblioteca dei programmi**: un modello della biblioteca
  globale, curata dal gestore della piattaforma.
- **Aggiungi da … → Manifestazione precedente**: copia tutti i programmi di
  un’altra manifestazione della società, comprese combinazioni e premi
  speciali. I programmi il cui nome e variante esistono già vengono saltati.
- **Aggiungi da … → File di programma**: importa un [file di
  programma](../reference/programme-file.md) esportato da un’altra
  manifestazione o società.

## L’editor

![Editor del programma](../images/screenshots/programme-editor.png){
.screenshot }

- **Nome** e **Variante**: i programmi con lo stesso nome sono varianti
  l’uno dell’altro; un tiratore si iscrive a una di esse. Il campo del nome
  propone i nomi dei programmi della manifestazione, e appena il nome
  corrisponde a uno di essi l’editor dice che stai creando una variante e
  avvisa se il nome della variante è vuoto o già usato.
- **Codice a barre**: stampato sul foglio di stand per la segnalazione
  elettronica (vedi [Codici](../reference/codes.md)).
- **Mostra sugli schermi chiosco**.
- **Serie**: colpi, punti massimi, se si contano le mouche e se si
  registrano i **Colpi singoli**, per serie. Con i colpi singoli si può
  registrare ogni colpo e la serie ne è la somma, p. es. per stamparli su un
  modulo; il totale della serie da solo si può ancora registrare.
- **Tentativi** (vedi sotto): quante volte un tiratore può tirare il
  programma nella manifestazione.
- **Prezzi (CHF)**: solo le categorie d’età selezionate possono tirare il
  programma. Con più tentativi, un secondo campo per categoria è il prezzo
  di ogni tentativo supplementare (*rachat*); vuoto significa lo stesso
  prezzo, 0 significa che il programma è pagato una volta.
- **Distinzioni**: nome, tipo (carta corona, premio in natura o premio in
  denaro pagato in contanti), valore e totale minimo per categoria d’età;
  lascia vuota una categoria dove la distinzione non è disponibile.

Non appena dei tiratori sono iscritti a un programma, le **serie** sono
fisse perché i risultati mantengano il loro significato; l'editor lo segnala
e le disattiva. Tutto il resto resta modificabile: nome, variante,
descrizione, codice a barre, le opzioni, prezzi, nomi delle serie, nomi,
valori e limiti delle distinzioni. Categorie d'età e distinzioni possono
essere **aggiunte** ma non rimosse. Su un programma in uso è inoltre
disponibile:

- **Attivo**: i programmi inattivi non si possono scegliere per nuove
  iscrizioni; quelle esistenti li mantengono.
- **Combinabile con**: il programma di cui questo può prendere le serie
  all’iscrizione.
- **Stand di tiro** (manifestazioni su più stand): lo stand su cui si tira
  il programma, p. es. *Pistola 50 m*; l’iscrizione offre allora solo i suoi
  bersagli e gli schermi chiosco raggruppano le classifiche per
  stand. *qualsiasi* lascia la scelta aperta. Programmi di stand diversi
  richiedono iscrizioni separate.
- **Modello di stampa** (quando la società ha dei [modelli di
  stampa](print-forms.md)): il modello su cui vengono stampati i risultati
  del programma; *nessuno* non stampa alcun modello.
- **Premi speciali** (vedi sotto).

**Esporta** scarica il programma come file; **Elimina** è possibile solo finché nessuno è iscritto.

## Tentativi {#attempts}

Di norma un tiratore tira un programma una volta per
manifestazione. **Tentativi per tiratore** ne consente di più, p. es. 3 per
un bersaglio fortuna con rachat a una festa di tiro, o **Illimitato** per un
programma annuale tirato durante tutta la stagione.

- **Tentativo che conta** decide quale tentativo conta per la classifica, la
  gara delle società, i premi speciali e *Il mio allenamento*: **Il
  migliore** (secondo le regole della classifica), **Il primo** o
  **L’ultimo**. La classifica lo mostra come *2/3*: conta il secondo di tre
  tentativi.
- **Solo il tentativo che conta riceve una distinzione**: disattivato, ogni
  tentativo può ricevere la sua distinzione (una carta corona per rachat);
  attivato, solo il tentativo che conta ne riceve una, e la distinzione
  passa a un tentativo migliore appena viene tirato.
- All’iscrizione, un programma con più tentativi mostra un contatore: due
  tentativi insieme danno due schede di risultati, *Tentativo 1* e
  *Tentativo 2*. Altri tentativi si possono aggiungere anche più tardi in
  una nuova iscrizione, o registrando un altro risultato in **Registra
  risultati**, finché il tiratore non ha più tentativi.

Il numero di tentativi può essere modificato mentre ci sono tiratori
iscritti, ma non al di sotto dei tentativi che un tiratore ha già. I
programmi con più tentativi non possono essere combinati.

## Premi speciali {#special-prizes}

**Premi speciali** nel menu del programma premia i migliori ranghi di ogni categoria d’età, per esempio i ranghi 1–3 in ogni categoria con almeno 10 partecipanti, contando solo i tiratori la cui società ha almeno 4 partecipanti in questo programma. I tiratori saltati lasciano posto al successivo. La regola esatta è nelle [regole di gara](../reference/rules.md#special-prizes); i vincitori sono nel rapporto *Premi speciali* in [Statistiche](statistics-reports.md).

![Finestra dei premi speciali](../images/screenshots/special-prizes.png){
.screenshot }

Ranghi 0 disattiva i premi speciali. Le impostazioni si possono cambiare in
qualsiasi momento, anche dopo la manifestazione.
