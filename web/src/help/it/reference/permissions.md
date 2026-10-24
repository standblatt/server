# Autorizzazioni e ruoli

Il ruolo di un membro in una società è un insieme di autorizzazioni. I ruoli
si modificano in **Utenti** → **Ruoli** (vedi
[Ruoli](../user-guide/users.md#roles)); i token API portano una selezione
delle stesse autorizzazioni.

## Autorizzazioni {#permissions}

| Diritto | Permette |
|---|---|
| **Vedere manifestazioni e classifiche** | manifestazioni, programmi, classifiche e classifiche generali |
| **Gestire le manifestazioni** | iscrizioni, risultati, distinzioni, prenotazioni, documenti e stampa |
| **Configurare le manifestazioni** | creare, configurare, copiare, chiudere ed eliminare manifestazioni, con programmi, classifiche generali, link chiosco e aiutante e prenotazione online |
| **Statistiche e rapporti** | statistiche, rapporti con importi e classifiche in CSV |
| **Vedere tiratori e società** | i dati personali dei dati di base |
| **Modificare tiratori e società** | registrare e modificare tiratori e società |
| **Gestire i dati di base** | importare tiratori, rinominare, unire ed eliminare società, eliminare tiratori |
| **Stand e stampanti** | stand di tiro, agenti di stampa e impostazioni delle stampanti |
| **Modelli di stampa** | modelli di stampa e immagini (solo finché il modulo della società è attivo) |
| **Utenti e società** | utenti, ruoli, inviti e moduli della società |

- Basta una delle autorizzazioni indicate per un’azione.
- Senza *Gestire le manifestazioni* una manifestazione si apre sulle sue
  classifiche.
- L’app mostra a ogni membro solo le pagine e i pulsanti che il suo ruolo
  permette.

## Ruoli predefiniti {#default-roles}

| Ruolo | Diritti |
|---|---|
| **Admin** | tutti; la colonna non si può modificare e il ruolo non si può eliminare |
| **Comitato** | tutti tranne *Stand e stampanti* e *Utenti e società* |
| **Aiutante** | vedere e gestire le manifestazioni, statistiche, vedere e modificare tiratori e società |
| **Membro** | vedere manifestazioni e classifiche |

- **Nuovo ruolo** parte solo con *Vedere manifestazioni e classifiche*.
- Un ruolo si può eliminare quando nessuno lo ha più; il ruolo admin mai.
- Una società mantiene sempre almeno un membro con il ruolo admin.
- I membri con *Utenti e società* ma senza il ruolo admin assegnano,
  modificano e rimuovono solo ruoli di cui hanno essi stessi le
  autorizzazioni.

## Link aiutante {#helper-links}

I link aiutante non sono ruoli della società; ogni link ha uno di questi
ruoli (vedi [Creare un link](../user-guide/helpers.md#create)):

| Ruolo | Vede | Può |
|---|---|---|
| **Marcatore** | la griglia dei propri bersagli, i tiratori su di essi | registrare risultati e distinzioni, **Terminare** un tiratore |
| **Distribuzione munizioni** | il campo di ricerca, la lista d’attesa, i tiratori su tutti i bersagli | segnare le munizioni come consegnate (**Munizione**), assegnare un bersaglio a un tiratore in attesa, stampare il foglio di stand |
| **Registrazione risultati** | il campo di ricerca, i tiratori su tutti i bersagli, l’elenco dei terminati | registrare risultati e distinzioni su ogni bersaglio |

Marcatori e registrazione risultati possono anche stampare i [modelli di
stampa](../user-guide/print-forms.md#use) dei loro bersagli; la
registrazione risultati può segnare le distinzioni come consegnate.

## Ambiti dei token API {#scopes}

Un token (**Account** → **Token API**, vedi [Token
API](../user-guide/login-language.md#api-tokens)) ha un elenco di società e
di ambiti:

| Ambito | Chiave | Permette |
|---|---|---|
| *Leggere il mio diario di allenamento* | `training:read` | leggere il proprio diario di allenamento |
| *Registrare e modificare le mie sessioni di allenamento* | `training:write` | aggiungere, modificare ed eliminare le proprie sessioni, discipline, armi e obiettivi |
| un’autorizzazione qui sopra | `events.view`, `events.operate`, `events.configure`, `reports.view`, `shooters.view`, `shooters.edit`, `shooters.manage`, `ranges.manage`, `forms.manage`, `club.manage` | lo stesso che in un ruolo, nelle società del token |

- Un’autorizzazione della società si può scegliere solo se il ruolo del
  proprietario la ha in una delle società del token.
- Un token non ha mai più del ruolo del proprietario: un’autorizzazione che
  il ruolo perde, o una società che il proprietario lascia, sparisce subito
  per il token.
- Un token senza società raggiunge solo il diario di allenamento.

### Preimpostazioni {#presets}

| Preimpostazione | Ambiti |
|---|---|
| **Classifiche per un sito web** | *Vedere manifestazioni e classifiche* |
| **App di allenamento** | *Leggere il mio diario di allenamento*, *Registrare e modificare le mie sessioni di allenamento* |
| **Solo lettura** | *Vedere manifestazioni e classifiche*, *Statistiche e rapporti*, *Vedere tiratori e società*, *Leggere il mio diario di allenamento* |

Una preimpostazione compila solo le autorizzazioni che le società scelte
permettono.

### Usare un token {#token-use}

| Cosa | Valore |
|---|---|
| Formato | `sb_` e 40 cifre esadecimali, mostrato una volta |
| Intestazione | `Authorization: Bearer sb_…` o `X-API-Key: sb_…` |
| Scadenza | tra 30 o 90 giorni, tra un anno, o mai |
| Endpoint | `/api/tenants/…` (dati della società, secondo l’autorizzazione), `/api/my/…` (diario di allenamento, secondo l’ambito), `GET /api/token`, `/api/version`, `/api/auth/options`; con *Configurare le manifestazioni* anche `/api/modules` e `/api/programme-templates`; tutto il resto risponde 403 |
| `GET /api/token` | nome, ambiti e società del token |
| Limite di richieste | per default 600 richieste per token in cinque minuti, poi 429; il gestore lo cambia con `STANDBLATT_API_RATE_LIMIT` (0 = nessun limite) |
| Token sconosciuti | dopo 20 token sconosciuti o revocati da un indirizzo in 15 minuti, 429 |
| Descrizione dell’API | `/api/docs` |
| Mai raggiungibili | password, e-mail, impostazioni dell’account, altri token |
