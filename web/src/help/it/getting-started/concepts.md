# Concetti

Questa pagina spiega le parole che usa Standblatt e come si combinano. I
termini tedeschi, francesi e italiani sono nel
[Glossario](../reference/glossary.md); i passi per ogni compito sono nella
guida utente.

In breve: una **società** ha **stand di tiro**, **tiratori** e
**utenti**. Organizza **manifestazioni**; una manifestazione offre
**programmi** e attiva i **moduli** di cui ha bisogno. Un’**iscrizione**
porta un tiratore alla manifestazione con uno o più programmi, ogni
programma diventa una **partecipazione**, la partecipazione riceve
**risultati**, e i risultati danno **distinzioni** e un posto nelle
**classifiche**. Quando la manifestazione è **chiusa**, le classifiche sono
definitive e raggiungono i diari di allenamento personali dei tiratori.

## Società, utenti e ruoli {#clubs-users-and-roles}

Una **società** è un mandante: tutto in Standblatt appartiene esattamente a
una società, e le società non vedono mai i dati delle altre. Solo il gestore
della piattaforma crea società, nel pannello di amministrazione. Gli account
nascono da un invito di una società, dalla registrazione a [Il mio
allenamento](../my-training/account.md), da **Aggiungi utente** nella pagina
utenti di una società o dal gestore.

Gli **utenti** sono account personali. Un account può appartenere a più società, con un ruolo per società. Un ruolo è un insieme di diritti come *Gestire le manifestazioni* o *Configurare le manifestazioni*; ogni società inizia con i ruoli *Admin*, *Comitato*, *Aiutante* e *Membro* e può modificarli o aggiungerne di propri (vedi [Ruoli](../user-guide/users.md#roles)). Le società invitano le persone via e-mail; la persona invitata entra accettando il link.

Perché gli account sono globali mentre tutti i dati della società restano
nella società è spiegato in [Società, account e
dati](../explanation/data.md).

## Stand di tiro e bersagli {#ranges-and-lanes}

Uno **stand di tiro** (inglese *range*) è uno stand della società con un
certo numero di **bersagli** (*lanes*). I bersagli possono essere segnati
fuori servizio. Le stampanti e l’agente di stampa appartengono a uno stand.

Ogni distanza è uno stand a sé: una sala con bersagli a 25 m e a 50 m sono
due stand, ognuno con bersagli da 1. Una manifestazione può essere tirata su
più stand, su uno o su nessuno (vedi [Manifestazioni con e senza
bersagli](../explanation/modules.md#lanes)).

## Manifestazioni e programmi {#events-and-programmes}

Una **manifestazione** (inglese *event*) è una gara su uno o più stand di
tiro in uno o più giorni (la data è il giorno principale; gli altri giorni
coprono il tiro anticipato e le manifestazioni di più giorni). Ha una
lingua, usata per i documenti stampati e gli schermi chiosco, e dei
[moduli](../user-guide/events.md#modules) che attivano ciò di cui ha
bisogno, come l’assegnazione dei bersagli, le distinzioni o il concorso di
sezione.

Non esistono tipi di manifestazione: una festa di tiro, un Tiro in campagna
e un programma annuale sono lo stesso tipo di manifestazione con moduli
diversi (vedi [Moduli e tipi di manifestazione](../explanation/modules.md)).

Una manifestazione offre **programmi**: ogni programma ha **serie** con un
numero di colpi e un punteggio massimo, **prezzi** per categoria d’età e
**livelli di distinzione** (carte corona, premi in natura o premi in denaro)
con un totale minimo per categoria d’età. I programmi con lo stesso nome
sono **varianti** (per esempio *Sport* e *Appoggiato*); un tiratore tira una
sola variante. Un programma può essere **combinabile** con un altro: prende
allora alcune serie da quel programma invece di essere tirato separatamente.

Un programma può permettere più **tentativi**, di cui uno conta. Una
**classifica generale** classifica i tiratori su più programmi o
tentativi. Come si collegano classifiche dei programmi, tentativi,
classifiche generali, concorso di sezione e premi speciali è spiegato in
[Come si combinano le classifiche](../explanation/rankings.md).

## Tiratori, iscrizioni e partecipazioni {#shooters-registrations-and-entries}

I **tiratori** sono le persone che tirano. Ogni società gestisce i propri dati di base dei tiratori, insieme alle **società** a cui appartengono (la società del tiratore, da non confondere con la società che usa Standblatt).

Un’**iscrizione** porta un tiratore alla manifestazione, eventualmente su un
bersaglio, con uno o più programmi. Ogni programma di un’iscrizione è una
**partecipazione** (*entry*). Un’iscrizione è *iscritta* (in attesa), *al
tiro* (su un bersaglio) o *conclusa* (tutti i programmi hanno risultati; il
bersaglio è di nuovo libero).

## Risultati, distinzioni e classifiche {#results-awards-and-rankings}

I **risultati** sono i punti per serie, facoltativamente con mouche (dieci interni). Quando tutte le serie sono registrate, il tiratore può ricevere una delle **distinzioni** che il suo totale raggiunge nella sua categoria d’età. Le **classifiche** ordinano tutti i risultati di un programma secondo le [regole della FST](../reference/rules.md#ranking).

Una distinzione dipende solo dal totale del tiratore stesso; un rango
dipende da quello di tutti gli altri. Le classifiche sono in diretta durante
la manifestazione e diventano definitive quando viene
[chiusa](../explanation/closing.md).

## Categorie d’età {#age-categories}

La categoria d’età segue la regola FST: l’anno della manifestazione meno
l’anno di nascita. Il giorno di nascita non conta, e l’anno è quello del
giorno principale della manifestazione, così un giorno di tiro anticipato in
dicembre conta per una manifestazione di gennaio. La categoria viene fissata
quando il tiratore si iscrive e decide i prezzi, le distinzioni e i premi
speciali che il tiratore può ottenere. Un tiratore senza data di nascita
conta come Élite. Le categorie e le loro età sono elencate in [Regole di
gara](../reference/rules.md#age-categories).

## Documenti e stampa {#documents-and-printing}

Ogni documento è un PDF generato dal server nella lingua della
manifestazione. Allo stand, un piccolo agente di stampa lo invia alla
stampante giusta; senza agente, il PDF si apre nel browser. Vedi [Come
funziona la stampa](../explanation/printing.md).
