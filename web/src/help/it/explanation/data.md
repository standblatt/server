# Società, account e dati

Standblatt serve molte società su un unico server. Questa pagina spiega a
chi appartengono quali dati, perché le persone hanno un solo account per
tutto e perché i tiratori di una società non sono quelli di un’altra.

## Le società sono mandanti {#tenants}

Una **società** in Standblatt è un mandante: i suoi stand di tiro, tiratori,
elenco delle società, manifestazioni, programmi, modelli di stampa e utenti
appartengono a lei e a nessun altro. Le società non vedono mai i dati delle
altre, nemmeno che un tiratore o una manifestazione esiste altrove. La
separazione fa parte di ogni interrogazione del server, non è un filtro
dell’interfaccia.

## Un account, un ruolo per società {#accounts}

Le persone sono diverse dai dati: la stessa persona può essere segretaria di
una società, aiutare all’accoglienza della festa di tiro di un’altra società
e tenere un diario di allenamento personale. Standblatt ha quindi **account
globali**. Un account accede una volta sola e raggiunge ogni società di cui
fa parte, con un ruolo per società che decide cosa può fare lì (vedi
[Ruoli](../user-guide/users.md#roles)).

Una società non possiede l’account, solo l’affiliazione. Rimuovere un membro
toglie l’accesso a quella società; l’account, il suo diario di allenamento e
le sue altre società restano. Le società invitano le persone via e-mail, e
la persona diventa membro solo accettando il link, così nessuno viene
aggiunto a una società a sua insaputa.

Gli aiutanti di una manifestazione non hanno bisogno di alcun account: un
[link aiutante](../user-guide/helpers.md) dà accesso a una manifestazione e
a un compito. Un [token API](../user-guide/login-language.md#api-tokens)
agisce per il suo proprietario e mai con più del ruolo del proprietario in
una società.

## Dati di base dei tiratori per società {#shooters}

Un **tiratore** non è un account. È una scheda nei dati di base di una
società: nome, numero di licenza, data di nascita, indirizzo e le società a
cui il tiratore appartiene. Ogni società tiene i propri tiratori, importati
da un file CSV, digitati all’accoglienza o ripresi da una prenotazione
online.

Non esiste un registro comune di tutti i tiratori con licenza: i dati dei
membri della federazione non fanno parte di Standblatt, e poiché le società
non vedono mai i dati delle altre, i tiratori inseriti da una società
restano invisibili alle altre. Il prezzo è che un tiratore che tira alle
manifestazioni di tre società viene registrato tre volte, una per
società. Dove le schede devono riunirsi, per il diario di allenamento del
tiratore, lo fanno il numero di licenza e la data di nascita (vedi [Chiudere
una manifestazione](closing.md#my-training)).

## La società del tiratore e la società titolare {#shooter-club}

Due cose diverse si chiamano società:

- La **società titolare** è la società che usa Standblatt e organizza la
  manifestazione.
- La **società del tiratore** è una voce nell’elenco **Società** della
  società titolare: la società a cui appartiene un tiratore, p. es. una
  società ospite a una festa di tiro. Queste società non hanno accesso a
  Standblatt; sono nomi con un numero FST.

Un tiratore può appartenere a più società e si iscrive per una di
esse. Quella società viene stampata sui documenti e conta nelle classifiche
e nel concorso di sezione. Poiché le società vengono riconosciute per nome
indipendentemente da maiuscole e spazi, e quelle nuove vengono create solo
di proposito, un errore di battitura non trasforma una società in due; dove
è successo comunque, l’unione rimette ordine, anche nelle manifestazioni
chiuse.

## Aggiornamenti in diretta {#live}

Più persone lavorano contemporaneamente sulla stessa manifestazione:
l’accoglienza iscrive, i marcatori registrano i risultati ai bersagli, gli
schermi chiosco mostrano le classifiche nella sala. Ogni volta che qualcosa
viene salvato, il server lo comunica a ogni pagina aperta di quella
manifestazione, e ogni pagina carica ciò che è cambiato. Un risultato
registrato a un bersaglio compare sulla dashboard e sugli schermi chiosco
senza ricaricare. Dopo un’interruzione della rete una pagina si rimette in
pari da sola, così nulla va perso.

## Offline {#offline}

Il lavoro della società richiede una connessione. Un bersaglio può essere
assegnato una sola volta, un tiratore può tirare un programma una sola
volta, e un risultato salvato su un dispositivo deve essere su tutti gli
altri; solo il server può deciderlo per tutti nello stesso momento. Senza
connessione, le pagine della società e i link aiutante mostrano una barra
rossa e non salvano nulla, ma i dialoghi aperti conservano quanto digitato,
così può essere salvato appena la connessione torna.

*Il mio allenamento* è diverso, perché nessun altro lavora su un diario personale. Una sessione salvata senza connessione viene conservata sul dispositivo e caricata più tardi, senza mai crearla due volte, e l’app installata si apre senza connessione con gli ultimi dati. Vedi [Sessioni incompiute e offline](../my-training/logging.md#unfinished-sessions-and-offline).
