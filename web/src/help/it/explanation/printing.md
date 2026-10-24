# Come funziona la stampa

Una festa di tiro stampa molto: un foglio di stand e una striscia di
controllo per ogni iscrizione, fogli dei risultati, carte corona,
modelli. Standblatt è un’app web, e questo determina come funziona la
stampa.

## Ogni documento è un PDF {#pdf}

Il server genera ogni documento come PDF, nella lingua della manifestazione:
fogli di stand, strisce di controllo, fogli dei risultati, carte corona,
modelli di stampa, diplomi e classifiche. Lo stesso PDF viene stampato allo
stand, aperto nel browser o condiviso da un telefono, così un documento ha
lo stesso aspetto ovunque esca. Le classifiche vengono solo scaricate, come
PDF o CSV, e mai inviate a una stampante dello stand.

## Perché un agente di stampa {#agent}

Un browser non può stampare da solo: ogni stampa chiede un dialogo, una
stampante e un clic. A un’accoglienza affollata è troppo lento, e un foglio
di stand dovrebbe uscire dalla stampante accanto all’accoglienza mentre la
striscia di controllo esce dalla stampante di etichette.

L’**agente di stampa** risolve il problema. È un piccolo programma su un
computer accanto alle stampanti: l’app Standblatt Print Agent su Windows o
su un Mac, o un agente senza interfaccia su un computer Linux come un
Raspberry Pi. Quando qualcuno preme un pulsante di stampa, il server mette
il PDF in una coda; l’agente lo preleva e lo stampa sulla stampante scelta
senza chiedere. L’agente si collega a Standblatt, non il contrario, così
nulla allo stand deve essere raggiungibile da internet.

Ogni lavoro di stampa viene confermato dall’agente. Un lavoro non confermato
entro due minuti viene riassegnato, fino a tre tentativi; **Lavori recenti**
mostra cosa è successo.

## Agenti e stampanti appartengono a uno stand {#range}

Le stampanti stanno in uno stand, non in una manifestazione, quindi l’agente
di stampa e le impostazioni delle stampanti appartengono allo [stand di
tiro](../getting-started/concepts.md#ranges-and-lanes). Ogni manifestazione
su quello stand li usa senza una configurazione propria. I documenti di
un’iscrizione si stampano sullo stand del suo bersaglio; una manifestazione
su più stand stampa quindi nello stand dove tira il tiratore. Un computer
può servire più stand, p. es. *Pistola 25 m* e *Pistola 50 m* nella stessa
sala.

Le impostazioni scelgono una stampante **per documento**, perché i documenti
richiedono stampanti diverse: la striscia di controllo va a una stampante di
etichette, il foglio di stand a un cassetto A5, le carte corona e i modelli
di stampa a un’alimentazione manuale con i fogli prestampati. Le opzioni di
stampa riportate dall’agente (cassetto, formato carta, fronte-retro e così
via) vengono salvate con la scelta.

Un agente viene collegato o con un token copiato nelle sue impostazioni o,
con l’app, tramite un codice mostrato sul computer e confermato in
Standblatt (vedi [Stampa](../user-guide/printing.md#connect)). In entrambi i
casi il server conserva solo un’impronta del token, mai il token stesso.

## Quando non c’è una stampante {#fallback}

Un documento va alla stampante solo se il suo stand ha una stampante scelta
per quel documento e un agente dello stand era online nell’ultimo
minuto. Altrimenti si apre il PDF: in una nuova scheda su un computer, nel
menu di condivisione su un telefono o tablet (stampare, salvare nei file,
inviare), così l’app resta aperta. Nulla va perso quando l’agente è spento;
l’accoglienza stampa semplicemente a mano.

Lo stesso vale per le manifestazioni senza stand: non c’è nessun agente a
cui inviare, quindi ogni documento si apre come PDF.

## Modelli di stampa su fogli prestampati {#forms}

Alcuni documenti non sono disegnati da Standblatt ma forniti su carta: carte
corona, i fogli di stand che l’esercito fornisce per il programma
obbligatorio, i diplomi propri di una società. Un **modello di stampa**
descrive un tale foglio: dove va ogni valore sulla carta, in millimetri
dall’angolo in alto a sinistra. Stampare un modello mette solo i valori sul
foglio che è già nella stampante.

Un modello ha una **fonte di dati**: i risultati di un programma,
un’iscrizione o un rango. La fonte decide quali valori il modello può
stampare e dove viene usato: i modelli dei risultati si stampano dopo che un
programma è concluso o come carte corona, i modelli di un’iscrizione come
strisce di controllo o tessere di partecipazione, i modelli di un rango come
diplomi dalle classifiche. La carta corona e la striscia di controllo sono a
loro volta modelli integrati che una società può copiare e adattare.

Una scansione del foglio vuoto può stare dietro il modello nell’editor, così
i campi si possono posizionare sulle sue caselle. La scansione serve solo
per posizionare e non viene mai stampata.

## La calibrazione in linea di principio {#calibration}

Le stampanti raramente stampano esattamente dove dovrebbero: alimentatori e
cassetti spostano la stampa di un millimetro o due. Un modello giusto sullo
schermo può quindi finire accanto alle caselle sulla carta. Lo **spostamento
della stampante** del modello sposta ogni elemento della stessa misura per
compensare: stampa un foglio di prova, misura di quanto sono spostati i
valori e inserisci la differenza. Lo spostamento fa parte del modello,
quindi un modello stampato su un’altra stampante che sposta diversamente
richiede una copia con il proprio spostamento. L’anteprima ignora lo
spostamento e mostra il modello come disegnato. I passi sono in [Stampare i
risultati sui fogli di stand
dell’esercito](../how-to/army-sheet.md#calibrate).
