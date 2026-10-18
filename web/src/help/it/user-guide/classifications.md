# Classifiche generali

Una classifica di programma classifica un programma. Una **classifica
generale** classifica i tiratori su più programmi o su più
[tentativi](programmes.md#attempts) di un programma, p. es. un campionato di
sezione come somma di tre programmi, o un campionato di società sui tre
migliori tentativi di un programma annuale. Le classifiche generali appaiono
in **Classifiche**, sugli schermi chiosco e come rapporto; non modificano
mai risultati o distinzioni.

## Creare una classifica generale {#create}

Nella scheda *Classifiche generali* della **Configurazione** della
manifestazione (admin), **Nuova classifica generale** apre l’editor:

- **Programmi**: ogni programma che conta, con i suoi **Tentativi che
  contano**. Con 3 contano i tre tentativi migliori (o i primi o gli ultimi
  tre, secondo il *Tentativo che conta* del programma). Il totale di un
  tiratore è la somma di tutti i tentativi contati.
- **Contare solo i programmi migliori**: con più programmi contano solo
  altrettanti programmi con più punti, p. es. i due migliori su tre; 0 li
  conta tutti.
- **Classificare solo i tiratori con tutti i risultati che contano**: i
  tiratori a cui manca ancora un programma o un tentativo sono elencati
  sotto la classifica, senza rango. Senza questa opzione ognuno è
  classificato secondo ciò che ha.
- **Spareggio**: a parità di totale decidono i criteri selezionati nel loro
  ordine; le frecce lo cambiano. *Serie dall’ultima all’indietro* confronta
  tutte le serie contate, *Più colpi singoli più alti* richiede i [colpi
  singoli](programmes.md) (più dieci, poi più nove …), *Età* segue la
  [regola FST](../reference/rules.md#ranking). Di default l’ordine FST.
- **Classi**: classifiche separate, ognuna per alcune categorie d’età,
  p. es. *Juniori* con U17–U21 e *Élite* con le altre. Una categoria d’età
  appartiene al massimo a una classe; i tiratori la cui categoria non è in
  nessuna classe non sono elencati. Senza classi tutti sono in un’unica
  classifica.
- **Mostra sugli schermi chiosco**.

## Esempi {#examples}

| Manifestazione | Programmi | Impostazioni |
|---|---|---|
| Campionato di sezione | Programma obbligatorio, Tiro in campagna, Bersaglio sezione | tentativi che contano 1 ciascuno, *Classificare solo i tiratori con tutti i risultati che contano* |
| Campionato di società sulla stagione | programma annuale con tentativi illimitati | tentativi che contano 3 |
| Festa di tiro, due migliori bersagli | tre bersagli | *Contare solo i programmi migliori* 2 |

**Copia per l’anno prossimo** copia le classifiche generali con i programmi. Eliminando un programma, questo scompare dalle classifiche generali.
