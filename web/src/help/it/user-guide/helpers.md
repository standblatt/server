# Aiutanti

Gli aiutanti alla manifestazione (marcatori, distribuzione munizioni,
ufficio risultati) lavorano con un **link aiutante** invece di un login. Un
link aiutante apre su qualsiasi telefono o portatile una vista ridotta di
una manifestazione; gli admin attivano il modulo *Link aiutante* in
**Configurazione** → **Moduli**, creano i link in **Link aiutante** nella
configurazione e li revocano quando la manifestazione è finita.

![Pagina aiutante del marcatore su un
telefono](../images/screenshots/helper-scorer.png){ .screenshot .phone }

## Ruoli {#roles}

| Ruolo | Vede | Può |
|---|---|---|
| **Marcatore** | la griglia dei propri bersagli, i tiratori su di essi | registrare risultati e distinzioni, **Terminare** un tiratore |
| **Distribuzione munizioni** | il campo di ricerca, la lista d’attesa, i tiratori su tutti i bersagli | segnare le munizioni come consegnate (**Munizione**), assegnare un bersaglio a un tiratore in attesa, stampare il foglio di stand |
| **Registrazione risultati** | il campo di ricerca, i tiratori su tutti i bersagli, l’elenco dei terminati | registrare risultati e distinzioni su ogni bersaglio |

Tutto il resto resta agli utenti della società: iscrivere tiratori,
modificare programmi, classifiche, statistiche, eliminare.

## Creare un link {#create}

1. Inserisci un nome (chi riceve il link, ad es. *Marcatore poligono A*),
   scegli il ruolo e, per un marcatore, il poligono e i bersagli (`1-5,
   8`). Un marcatore di solito segue più bersagli.
2. **Crea link aiutante**; copia il link e mandalo all’aiutante (ad es. via
   messaggistica). Viene mostrato una sola volta.
3. L’aiutante apre il link; resta valido finché non viene revocato o la
   manifestazione eliminata. `?language=de` (`fr`, `it`, `en`) sovrascrive
   la lingua della manifestazione.

**revoca** nell’elenco ferma subito un link; *visto l’ultima volta* indica quando è stato usato l’ultima volta. Un marcatore che apre un tiratore su un altro bersaglio riceve *Questo bersaglio non è il tuo*.

## Munizioni {#ammunition}

**Munizione** su un’iscrizione in attesa o al tiro registra l’ora in cui le munizioni sono state consegnate; il pulsante esiste solo se il modulo *Distribuzione munizioni* è attivo (i tiratori di pistola e carabina ad aria compressa di solito portano le proprie). La dashboard mostra un badge verde *munizione* su queste iscrizioni, così distribuzione e rangeur vedono a colpo d’occhio chi è stato servito.
