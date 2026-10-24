# Importazione CSV

I membri con il diritto *Gestire i dati di base* importano i tiratori nella
pagina **Tiratori**. Il file ha bisogno di una riga d’intestazione; il
separatore (`;` o `,`) viene preso da essa.

```text
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9999;Hilltown;PS Example
```

| Colonna | Obbligatoria | Note |
|---|---|---|
| `firstname` | sì | |
| `lastname` | sì | |
| `license` | no | anche `licensenumber`; i numeri puramente numerici vengono completati a sei cifre |
| `birthday` | no | `YYYY-MM-DD` o `DD.MM.YYYY` |
| `postalcode`, `city` | no | |
| `club` | no | la società del tiratore, confrontata senza tenere conto di maiuscole e spazi |

- Una riga per tiratore e società; un tiratore in due società compare due
  volte.
- I tiratori vengono riconosciuti dal numero di licenza, o da nome e data di
  nascita se manca. I tiratori esistenti vengono aggiornati, niente viene
  eliminato.
- Un file con righe non valide viene rifiutato per intero, con l’elenco
  delle righe da correggere.
- Le società che non sono nell’elenco fermano l’importazione prima che venga
  salvato qualcosa: scegli per ognuna una società esistente (i nomi simili
  vengono proposti) o una nuova, poi clicca di nuovo su **Importa**.

## Risultati {#results}

I membri con il permesso *Configurare le manifestazioni* importano i
risultati con **Importa CSV** nella pagina **Registra risultati** di una
manifestazione (vedi [Registra
risultati](../user-guide/results-entry.md#import)):

```text
license;firstname;lastname;birthday;club;programme;s1;s2;m2;s1_1;s1_2;s1_3
013756;;;;;Championship;96;93;2;;;
;Max;Muster;1990-02-02;SG Hilltown;Obligatory;;45;;9;10;8
```

| Colonna | Obbligatoria | Note |
|---|---|---|
| `license` oppure `firstname` + `lastname` | sì | il tiratore: trovato tramite la licenza, altrimenti tramite il nome (e `birthday`, se più persone hanno quel nome); creato se non trovato, il che richiede `firstname` e `lastname`: una licenza da sola che non viene trovata è un errore di riga |
| `birthday` | no | `YYYY-MM-DD` o `DD.MM.YYYY` |
| `club` | no | deve essere nell’elenco delle società (senza tenere conto di maiuscole e spazi); il tiratore vi entra |
| `programme` | sì | codice a barre, nome con variante come `Championship (Sport)` o il solo nome |
| `s1`, `s2`, … | sì | punti per serie, decimali con `.` o `,` |
| `m1`, `m2`, … | no | mouche per serie |
| `s1_1`, `s1_2`, … | no | colpi singoli della serie 1, per le serie che li registrano; la serie ne è la somma |

- Ogni riga viene prima controllata; con un qualsiasi errore di riga non
  viene importato nulla.
- Un tiratore può avere risultati in un programma una sola volta; una riga
  per risultati già presenti è un errore. I programmi con più
  [tentativi](../user-guide/programmes.md#attempts) prendono una riga per
  tentativo, finché il tiratore non ne ha più.
- Per i programmi con più tentativi, importare di nuovo lo stesso file
  aggiunge le sue righe come ulteriori tentativi.
- Se il tiratore è iscritto al programma senza risultati, i risultati vanno
  in quell’iscrizione.
