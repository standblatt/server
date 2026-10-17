# Registra risultati

I risultati si possono registrare in tre modi, combinabili all’interno di
una manifestazione:

- **Al bersaglio** su tablet o telefono, tramite un link per segretari (vedi
  [Aiutanti](helpers.md)): il segretario registra le serie, o ogni colpo,
  dei tiratori dei suoi bersagli.
- **Al desk** tramite l’iscrizione e i suoi risultati (vedi
  [Risultati](results.md)), come in una festa di tiro.
- **In seguito** in **Registra risultati** nella barra laterale della
  manifestazione, tiratore per tiratore, p. es. da fogli di carta o dal
  display dell’impianto di segnalazione, oppure come file esportato
  dall’impianto. Questa pagina descrive entrambi.

## Tiratore per tiratore {#quick}

1. Scegli il **Programma**; resta scelto per i tiratori successivi.
2. Digita un numero di licenza o un nome in **Tiratore** e premi ++enter++:
   viene preso il primo risultato. Funziona anche un lettore di codici a
   barre. **Nuovo tiratore** aggiunge qualcuno che non è ancora nei dati di
   base.
3. Digita le serie, o i colpi singoli per le serie che li registrano (vedi
   [Programmi](programmes.md)), e premi ++enter++.

Un tiratore può avere risultati in un programma una sola volta; una riga per
risultati già presenti è un errore.

Le distinzioni si scelgono come di consueto: la distinzione raggiunta più
preziosa viene impostata automaticamente e si può cambiare nei risultati.

## Importazione da un file {#import}

**Importa CSV** (admin) in **Registra risultati** legge un file con una riga per tiratore e programma, p. es. esportato dall’impianto di segnalazione o digitato in un foglio di calcolo; le colonne sono descritte in [Importazione CSV](../reference/csv-import.md#results). Ogni riga viene controllata prima: con un errore non viene importato nulla e il dialogo elenca le righe da correggere. I tiratori che non sono nei dati di base vengono creati; una società deve già essere nell’elenco.
