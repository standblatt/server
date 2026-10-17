# Regole di gara

Questa pagina descrive le regole esattamente come Standblatt le
implementa. Si trovano in `internal/domain` e sono coperte da test unitari.

## Categoria d’età

La categoria risulta da *anno della manifestazione − anno di nascita*; vedi
la tabella in
[Concetti](../getting-started/concepts.md#age-categories). Senza data di
nascita un tiratore conta come élite. La categoria viene fissata
all’iscrizione.

## Iscrizione

Un’iscrizione viene rifiutata quando:

- un programma non appartiene alla manifestazione, è scelto due volte o non
  è attivo;
- un programma non ha un prezzo per la categoria d’età del tiratore («non
  per questa categoria d’età»);
- sono scelte due varianti dello stesso programma;
- un programma prende serie da un altro programma che non è scelto, o il
  numero di serie scelte non corrisponde;
- il bersaglio è occupato da un’altra iscrizione attiva o fuori servizio;
- il tiratore è già iscritto a uno dei programmi di questa manifestazione
  (un tiratore tira ogni programma una volta per manifestazione; un’altra
  variante in un’iscrizione separata è ammessa).

Un bersaglio prenotato si può comunque usare dopo una conferma.

## Risultati

Ogni serie deve essere registrata. I punti sono tra 0 e il massimo della
serie (per default colpi × 10), le mouche tra 0 e il numero di colpi, e solo
sulle serie che contano le mouche. Il salvataggio dei risultati assegna la
distinzione più preziosa raggiunta; un nuovo salvataggio mantiene la
distinzione scelta finché il nuovo totale la raggiunge ancora, altrimenti la
più preziosa raggiunta. Appena ogni programma di un’iscrizione ha risultati,
l’iscrizione è conclusa e il suo bersaglio si libera.

### Combinazioni

Un programma segnato come *combinabile con* un altro programma può prenderne
delle serie: all’iscrizione il tiratore sceglie quali (per esempio le serie
2 e 3 di un programma a quattro serie per un programma a due serie). Le
serie scelte vengono copiate in ordine crescente e rinumerate da 1; il
programma combinato non si può modificare da solo e segue ogni modifica
della fonte.

## Distinzioni

Un livello di distinzione è disponibile per una categoria d’età solo se ha
un limite per essa. Un tiratore può scegliere qualsiasi livello di cui il
suo totale raggiunge il limite; il livello più prezioso è in cima e viene
assegnato per default al salvataggio dei risultati. Il server verifica la
scelta.

## Classifica {#ranking}

Tutti i risultati di un programma sono ordinati per:

1. totale più alto;
2. le serie confrontate dall’ultima all’indietro (vince l’ultima serie
   migliore, poi la penultima, …);
3. più mouche;
4. età: tra i juniori (U10–U21) prima il più giovane, tra gli altri prima il
   più anziano; i non juniori vengono prima dei juniori.

I ranghi sono consecutivi (1, 2, 3, … anche a parità di risultato). La
stessa classifica alimenta le pagine del personale, gli schermi chiosco, il
PDF e il CSV della classifica e i premi speciali. È stata verificata con le
classifiche pubblicate del campionato appenzellese di pistola ad aria 2026:
identica su tutte le 260 righe.

## Concorso di sezione {#club-competition}

Partecipano solo i programmi segnati *conta per il concorso di sezione*, e
solo i risultati conclusi di tiratori con una società.

- Il valore di un risultato è la somma delle sue prime due serie.
- Una società con meno di 5 partecipanti è elencata come *non classificata*.
- Su *n* partecipanti i migliori *max(5, ⌊n/2⌋)* sono obbligatori, gli altri
  facoltativi.
- Risultato = ⌊(Σ obbligatori + 0.02 × Σ facoltativi) / numero di
  obbligatori × 1000⌋ / 1000, cioè troncato a tre decimali.
- Parità: prima più partecipanti, poi i valori ordinati confrontati uno a
  uno.
- La società esclusa nelle impostazioni della manifestazione (di solito
  l’organizzatrice) mantiene il suo posto nell’elenco ma non riceve un rango
  («fuori concorso»).

## Premi speciali {#special-prizes}

Un programma può assegnare premi speciali ai migliori ranghi di ogni
categoria d’età. Tre numeri li definiscono:

- **Ranghi per categoria**, p. es. 3;
- **partecipanti minimi per categoria** in questo programma, p. es. 10;
- **partecipanti minimi della società del tiratore** nello stesso programma,
  p. es. 4.

Per ogni categoria d’età con abbastanza partecipanti vincono i primi
tiratori in ordine di classifica la cui società ha abbastanza partecipanti;
i tiratori la cui società ne ha troppo pochi vengono saltati e sale il
successivo. I tiratori senza società non sono ammessi, salvo che il minimo
per società sia 0. Il rapporto elenca sia i vincitori sia i tiratori saltati
con il motivo.

!!! example "Campionato appenzellese di pistola ad aria"
    Il bando dice: *«In ogni categoria i ranghi 1–3 ricevono un premio speciale, purché la categoria abbia almeno 10 partecipanti e la società del tiratore almeno 4.»* Cioè ranghi 3, minimo di categoria 10, minimo di società 4. Applicata ai risultati 2026, la regola riproduce la lista dei vincitori fatta a mano.

## Elenchi di bersagli

I bersagli si scrivono come elenco di numeri e intervalli: `1, 5,
7-10`. Standblatt li mostra in forma compatta e trasforma le sequenze di tre
o più in intervalli (`1-3, 5, 7, 8`).

## Importi e numeri

Gli importi sono mostrati come `CHF 1'234.50`, le date come
`dd.mm.yyyy`. Sulle carte corona il valore può essere stampato in lettere
nella lingua della manifestazione, p. es. *Fünfundzwanzig Franken*,
*Septante et un francs*, *Ventitré franchi*.
