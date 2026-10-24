# Come si combinano le classifiche

Una manifestazione può classificare i suoi tiratori in più modi
contemporaneamente: per programma, su più programmi o tentativi, per società
e per categoria d’età per i premi speciali. Non sono gare separate con dati
separati. Partono tutte dagli stessi risultati, e la maggior parte parte
dalla stessa classifica del programma.

## La classifica del programma {#programme}

La classifica del programma è la base. Ogni programma ha una classifica di
tutti quelli che lo hanno tirato, ordinata secondo le [regole
FST](../reference/rules.md#ranking): totale, poi le serie dall’ultima
all’indietro, poi le mouche, poi l’età. Esiste una sola implementazione di
questo ordine, e alimenta la pagina delle classifiche, gli schermi chiosco,
il PDF e il CSV della classifica, i premi speciali e il rango che un
tiratore vede in *Il mio allenamento*. Il posto di un tiratore è quindi lo
stesso ovunque compaia.

I ranghi sono consecutivi (1, 2, 3, … anche a parità di risultato): lo
spareggio decide quasi ogni parità, così due tiratori non condividono mai un
rango.

Le **varianti** sono programmi separati con lo stesso nome, per esempio *Sport* e *Appoggiato*, e ognuna ha la propria classifica: il tiro appoggiato e il tiro a braccio libero non vengono confrontati.

Un programma **combinato** prende serie da un altro programma invece di
essere tirato separatamente. I suoi risultati sono copie delle serie scelte
della sua fonte, quindi viene classificato come qualsiasi altro programma;
quando la fonte cambia, la copia segue, e con essa la classifica.

## I tentativi e il tentativo che conta {#attempts}

Un programma che permette più tentativi conserva ogni tentativo come
partecipazione con i propri risultati. Una classifica, però, elenca ogni
tiratore una sola volta. Il **tentativo che conta** del programma (il
migliore, il primo o l’ultimo) decide quale tentativo rappresenta il
tiratore, e la classifica mostra quale sia, p. es. *2/3*.

Lo stesso tentativo che conta viene usato ovunque serva un solo risultato
per tiratore: la classifica, il concorso di sezione, i premi speciali e il
rango in *Il mio allenamento*. Le distinzioni sono diverse, perché una
distinzione appartiene a un risultato e non a una classifica: ogni tentativo
può ottenere la propria distinzione, a meno che *Solo il tentativo che conta
riceve una distinzione* sia attivo; in quel caso la distinzione passa al
tentativo che conta.

Le combinazioni prendono le serie di un risultato; con più tentativi non
sarebbe chiaro di quale, quindi i programmi con più tentativi non possono
essere combinati.

## Classifiche generali {#classifications}

Una classifica generale classifica i tiratori su più programmi o su più
tentativi di un programma: la somma di tre programmi per un campionato di
sezione, i tre migliori tentativi di un programma annuale per un campionato
sociale. Prende i tentativi che contano di ogni programma, somma ciò che
conta, può contare solo i programmi migliori e classifica con il proprio
spareggio e le proprie classi di categorie d’età.

Una classifica generale legge soltanto i risultati. Non modifica mai un
risultato, una distinzione o una classifica di programma, e una
manifestazione può avere più classifiche generali accanto alle sue
classifiche dei programmi. La classifica del programma stessa si può vedere
come la classifica generale più semplice: un programma, un tentativo, lo
spareggio FST e nessuna classe. Vedi [Classifiche
generali](../user-guide/classifications.md) per le impostazioni.

## Concorso di sezione {#club-competition}

Il concorso di sezione classifica le società, non i tiratori. Usa solo i
programmi segnati con *Conta per il concorso di sezione* e solo i risultati
completi dei tiratori che si sono iscritti per una società; la società è
quella scelta all’iscrizione, non la società titolare. Ogni risultato conta
con le sue prime due serie, e il risultato di una società è la media dei
suoi migliori risultati, con un piccolo bonus per ogni ulteriore
partecipante (vedi [Concorso di
sezione](../reference/rules.md#club-competition)). Poiché ha un proprio
punteggio, il concorso di sezione non segue la classifica del programma.

## Premi speciali {#special-prizes}

I premi speciali vanno ai migliori ranghi di ogni categoria d’età. Seguono
invece la classifica del programma: la percorrono per categoria d’età e
saltano i tiratori la cui società o categoria ha troppo pochi partecipanti
(vedi [Premi speciali](../reference/rules.md#special-prizes)). Le loro
impostazioni possono essere modificate durante la manifestazione, perché
leggono soltanto la classifica.

## Quale classifica usa cosa {#overview}

| Classifica | Classifica | Da | Tentativo che conta | Ordine |
|---|---|---|---|---|
| Classifica del programma | tiratori | un programma | sì | regole FST |
| Classifica generale | tiratori, per classe | più programmi o tentativi | un proprio numero per programma | proprio spareggio, FST come predefinito |
| Concorso di sezione | società | programmi che contano per il concorso di sezione | sì | risultato della società |
| Premi speciali | tiratori, per categoria d’età | la classifica del programma | sì | la classifica del programma |

Le distinzioni non sono una classifica: un tiratore raggiunge una
distinzione con il proprio totale nella sua categoria d’età, comunque tirino
gli altri.

## Perché le serie sono bloccate {#locked}

I risultati vengono salvati per serie. Appena ci sono tiratori iscritti a un
programma, una modifica al numero di serie, ai loro colpi, al loro massimo o
al fatto che registrino mouche o colpi singoli cambierebbe il significato
dei risultati salvati: un totale di 96 su 100 non è un totale su 200, e ogni
classifica che li confronta sarebbe sbagliata. Standblatt fissa quindi le
serie appena un tiratore è iscritto.

Tutto ciò che non cambia il significato di un risultato resta aperto: nomi,
prezzi, valori e limiti delle distinzioni, le opzioni. Categorie d’età e
livelli di distinzione si possono aggiungere ma non togliere, perché le
iscrizioni sono state prezzate e le distinzioni assegnate in base a essi.

## In diretta e definitive {#live}

Tutte le classifiche vengono calcolate dai risultati ogni volta che vengono
mostrate, così un risultato salvato compare subito sulla pagina delle
classifiche e sugli schermi chiosco. Durante la manifestazione possono
ancora cambiare. [Chiudere](closing.md) la manifestazione le rende
definitive.
