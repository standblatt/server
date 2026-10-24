# Link senza accesso

Questi link aprono una pagina senza accedere. Chi ha il link può usarlo,
quindi invialo solo alle persone a cui è destinato.

| Link | Creato con | Mostra | Valido |
|---|---|---|---|
| Link chiosco `/kiosk/<secret>` | **Crea link chiosco** in **Display chiosco** nella configurazione | le classifiche in diretta della manifestazione: programmi e classifiche generali con *Mostra sugli schermi chiosco*, e lo stato dei bersagli con *Mostra lo stato dei bersagli sulle classifiche* | fino a **revoca** o all’eliminazione della manifestazione |
| Link aiutante `/h/<secret>` | **Crea link aiutante** in **Link aiutante** nella configurazione | un’accoglienza ridotta della manifestazione per un [ruolo aiutante](permissions.md#helper-links) | fino a **revoca** o all’eliminazione della manifestazione |
| Link di prenotazione `/b/<code>` | **Attiva** in **Prenotazioni** nella configurazione | le fasce orarie della manifestazione con i bersagli liberi per stand, e il modulo di richiesta | fino a **Nuovo link** o **Disattiva**; richieste solo fino al termine |
| Link di gestione della prenotazione `/b/manage/<secret>` | inviato via e-mail alla società dopo **Invia la richiesta** | lo stato della richiesta della società; modificarla o annullarla | modifiche e annullamento fino al termine |
| Link d’invito `/invite/<secret>` | **Invita via e-mail** in **Utenti** | la società e il ruolo; accettare con un account o una nuova password | sette giorni; termina quando accettato, ritirato o inviato di nuovo |
| Pagina di collegamento `/connect?code=<code>` | aperta dall’app Standblatt Print Agent | il computer del codice; società, stand e nome dell’agente | il codice: dieci minuti; richiede un accesso con *Stand e stampanti* |
| Conferma dell’e-mail `/verify?token=…` | aggiunta o modifica dell’indirizzo e-mail, registrazione | conferma l’indirizzo | sette giorni |
| Link per la password `/reset-password?token=…` | **Password dimenticata?**; il gestore nel pannello di amministrazione | imposta una nuova password | un’ora; dal gestore sette giorni |

La parte segreta dei link chiosco, aiutante, di gestione e d’invito viene
mostrata o inviata una volta sola e salvata solo come hash: un link perso
non può essere mostrato di nuovo, solo sostituito.

## Parametri dell’URL {#parameters}

| Parametro | Link | Effetto |
|---|---|---|
| `?theme=dark`, `?theme=light` | chiosco | visualizzazione scura o chiara |
| `?language=de` (`fr`, `it`, `en`) | chiosco, aiutante | lingua della pagina invece della lingua della manifestazione |
| `?range=<id>` | chiosco | in una manifestazione con più stand: solo i bersagli di quello stand, i suoi programmi e i programmi senza stand |
| `?code=` | collegamento | il codice da confermare, compilato dall’app |

**Crea link chiosco** scrive nel link il tema, la lingua e lo stand scelti lì.
