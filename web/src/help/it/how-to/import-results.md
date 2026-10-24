# Importare i risultati dal sistema di punteggio

Ti serve quando i risultati di molti tiratori sono in un file,
p. es. esportato dal sistema elettronico di punteggio o digitato in un
foglio di calcolo. Richiede il permesso *Configurare le manifestazioni*
(vedi [Ruoli](../user-guide/users.md#roles)); la manifestazione non deve
essere chiusa.

## Preparare il file {#file}

1. Esporta o salva i risultati come CSV con una riga d’intestazione e una
   riga per tiratore e programma. Le colonne sono elencate in [Importazione
   CSV](../reference/csv-import.md#results): almeno `license` (o `firstname`
   e `lastname`), `programme` e le serie `s1`, `s2`, ….
2. Nella colonna `programme` usa il codice a barre del programma o il suo
   nome, con la variante tra parentesi dove ci sono varianti,
   p. es. `Championship (Sport)`.
3. Se il file ha una colonna `club`, controlla che ogni società sia
   nell’elenco delle società in **Dati di base** → **Società**; aggiungi lì
   quelle mancanti.

## Importarlo {#import}

1. Apri la manifestazione e clicca **Registra risultati** nella barra
   laterale.
2. Clicca **Importa CSV**, scegli il file e clicca **Importa**.
3. *Importati: … risultati, … nuovi tiratori.* significa che hai finito; i
   risultati sono in **Classifiche**.

Importa un file una sola volta. Un secondo import viene rifiutato per i
programmi che si tirano una volta (i tiratori hanno già risultati), ma un
programma con più [tentativi](../user-guide/programmes.md#attempts) riprende
ogni riga come un ulteriore tentativo.

## Correggere gli errori di riga {#errors}

Con un qualsiasi errore non viene importato nulla e il dialogo mostra *Nulla
importato, correggere queste righe:* con i numeri di riga. Correggi ogni
riga elencata nel file e importa di nuovo l’intero file.

| Messaggio | Correzione |
|---|---|
| Programme "…" not found | scrivi il programma come sulla sua pagina del programma, con la variante se ce ne sono più, o usa il suo codice a barre |
| Club "…" is not in the club list | correggi l’ortografia o aggiungi la società in **Società** |
| No shooter with the licence …; add first and last name to create one | aggiungi `firstname` e `lastname` alla riga, così il tiratore può essere creato |
| Several shooters are called …; add the birthday or the licence | aggiungi `birthday` o `license` alla riga |
| Series … is missing, … is not a number | compila ogni serie del programma; i decimali accettano `.` o `,` |
| Series …: points must be between 0 and … | correggi i punti; il massimo è fissato dalle serie del programma |
| Series …: enter all … shots | indica ogni singolo colpo della serie, o nessuno e solo il totale della serie |
| The shooter already has results in … | rimuovi la riga; i risultati sono già in Standblatt |
| The same shooter and programme appear twice in the file | tieni una delle due righe |
| The shooter has no attempts left in … | rimuovi la riga, o aumenta i tentativi del programma |

Un messaggio sulla riga d’intestazione (*The header row needs the column
programme*, *… the column license or the columns firstname and lastname*)
significa che nella prima riga manca un nome di colonna o è scritto male.

I tiratori che non sono nei dati di base vengono creati dall’import. I
singoli risultati si possono ancora correggere in seguito nei risultati
dell’iscrizione (vedi [Risultati](../user-guide/results.md)).
