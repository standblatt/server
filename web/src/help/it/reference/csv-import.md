# Importazione CSV

Gli admin della società importano i tiratori nella pagina **Tiratori**. Il
file ha bisogno di una riga d’intestazione; il separatore (`;` o `,`) viene
preso da essa.

```text
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9050;Appenzell;PS Example
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
