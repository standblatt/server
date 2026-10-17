# Tiratori e società

Ogni società gestisce i propri dati di base dei tiratori; niente è condiviso
con altre società. **Tiratori** e **Società** si trovano sotto **Dati di
base** nella barra laterale.

## Tiratori

![Tiratori](../images/screenshots/shooters.png){ .screenshot }

- La **ricerca** funziona come nella dashboard: numero di licenza (zeri
  iniziali facoltativi), parte di un nome, o nome e cognome in qualsiasi
  ordine. L’elenco mostra 50 tiratori per pagina, l’elenco delle società 50
  società; le frecce sotto sfogliano le pagine.
- Un clic su un tiratore apre i dettagli: nome, numero di licenza, data di
  nascita, indirizzo e le società a cui appartiene.
- **Società** (nuovo tiratore o aggiunta di una società nei dettagli)
  propone le società dell’elenco. Un nome che corrisponde a una società
  senza tenere conto di maiuscole e spazi è quella società. Ogni altro nome
  mostra *non è ancora nell’elenco delle società* con le società simili
  sotto *Intendevi*; una nuova società nasce solo con **Crea la nuova
  società**, così un errore di battitura non diventa una seconda società.
- **Importa** (admin) legge un file CSV, vedi [Importazione
  CSV](../reference/csv-import.md).
- I tiratori con iscrizioni non si possono eliminare.

## Società

![Società](../images/screenshots/clubs.png){ .screenshot }

Le **società** sono le società a cui appartengono i tiratori, con il loro numero FST. Compaiono nelle iscrizioni, nelle classifiche e nel concorso di sezione. Un tiratore può appartenere a più società e si iscrive per una di esse. Le società usate da iscrizioni non si possono eliminare, e due società non possono avere lo stesso nome, nemmeno con altre maiuscole o spazi.

**Unisci a…** (admin) elimina un doppione, p. es. *S.G. Hiltown* accanto a *SG Hilltown*: scegli la società che resta (le società con un nome simile vengono prima). Mantiene il suo nome e riceve i tiratori, le iscrizioni e le richieste di prenotazione dell’altra, anche quelli delle manifestazioni chiuse, così le loro classifiche mostrano la società che resta; il numero FST viene ripreso se non ne ha uno. L’altra società viene eliminata.
