# Chiudere una manifestazione

Mentre una manifestazione è in corso, tutto può ancora cambiare: un
risultato viene corretto, un’iscrizione si sposta, un programma riceve un
nuovo prezzo. A un certo punto l’organizzatore deve dire *questo è il
risultato*. È ciò che fa **Chiudi la manifestazione**: trasforma una
manifestazione in corso in una manifestazione conclusa.

## Cosa significa chiudere {#meaning}

La chiusura fa due cose contemporaneamente:

- **Le classifiche diventano definitive e vengono pubblicate.** Finora erano
  uno stato di lavoro; d’ora in poi sono il risultato della manifestazione.
- **La manifestazione diventa di sola lettura.** Iscrizioni, partecipazioni,
  risultati, distinzioni, programmi con le loro serie, prezzi e livelli di
  distinzione, classifiche generali e impostazioni della manifestazione non
  possono più essere modificati.

Tutto resta visibile: la dashboard, i risultati, le classifiche, le
statistiche, i rapporti e i documenti. Un banner su ogni pagina della
manifestazione indica che è chiusa.

Le due cose vanno insieme. Una classifica può essere definitiva solo se
nulla di ciò che sta sotto può più cambiare, e Standblatt lo impone nella
banca dati stessa, non solo nascondendo i pulsanti: ogni modifica ai dati di
una manifestazione chiusa viene rifiutata, che venga dall’app, da un link
aiutante o da un token API.

## Perché le manifestazioni chiuse sono di sola lettura {#read-only}

Una volta chiusa una manifestazione, altre cose si basano su di essa:

- I tiratori copiano i loro risultati e ranghi nel loro diario di
  allenamento personale.
- Il rapporto finanziario, le distinzioni e il concorso di sezione sono
  stati consegnati e regolati.
- La manifestazione dell’anno prossimo può essere stata copiata da essa.

Se i risultati potessero ancora cambiare, il diario di un tiratore potrebbe
mostrare un rango che la classifica non ha più, e i rapporti non
corrisponderebbero più a quanto è stato pagato e consegnato. Per questo
anche ciò che avviene all’accoglienza, come segnare le distinzioni come
consegnate, va fatto prima della chiusura.

## Il legame con Il mio allenamento {#my-training}

Ogni utente di Standblatt ha un diario di allenamento personale, [Il mio
allenamento](../my-training/index.md). Un tiratore che inserisce un **numero
di licenza** e una **data di nascita** nel profilo viene confrontato con i
tiratori di tutte le società: ovunque i dati di base di una società
contengano un tiratore con lo stesso numero di licenza (gli zeri iniziali
non contano) e la stessa data di nascita, i risultati di quel tiratore
appartengono all’utente.

Vengono offerti solo i risultati di manifestazioni **chiuse**, perché solo
la loro classifica è definitiva. Ognuno arriva con il rango nel suo
programma e il numero di partecipanti, p. es. *Rango 12 di 87*. Nulla viene
aggiunto da sé: il tiratore sceglie una disciplina e aggiunge il risultato,
che diventa una sessione di gara del diario con le sue serie e il suo rango.

Il confronto funziona tra società senza che nessuna società ceda i suoi
dati: i dati di base delle società restano dove sono, il confronto li legge
solo per l’utente a cui appartengono e mostra solo i risultati del tiratore
stesso. Un tiratore senza numero di licenza o data di nascita nel profilo, o
la cui società ha registrato un’altra data di nascita, non riceve risultati
di gara. Vedi [Profilo e risultati di
gara](../my-training/account.md#competitions).

## Riaprire {#reopen}

**Riapri** rende di nuovo modificabile una manifestazione chiusa, p. es. quando un risultato errato emerge dopo la chiusura. Serve per le correzioni, non come secondo stato di lavoro:

- I risultati che i tiratori hanno già aggiunto al loro diario restano lì
  come sono stati aggiunti; una correzione successiva non li raggiunge.
- I risultati non ancora aggiunti non vengono più offerti in *Il mio
  allenamento* finché la manifestazione non viene chiusa di nuovo.

Chiudi di nuovo la manifestazione appena la correzione è fatta.

Una manifestazione chiusa non può essere eliminata; riaprila prima. Può
essere copiata per l’anno prossimo mentre è chiusa.
