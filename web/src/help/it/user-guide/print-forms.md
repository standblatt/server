# Modelli di stampa

Un modello di stampa porta i risultati su un foglio prestampato, p. es. i
fogli di stand che l’esercito fornisce per il programma obbligatorio: i
colpi si registrano al bersaglio o in seguito, e il foglio viene stampato
invece di essere compilato a mano. Compilare a mano resta possibile; i
modelli sono facoltativi e si possono disattivare in **Moduli della
società** (vedi [Moduli della società](events.md#club-modules)).

## Creare un modello {#create}

**Modelli di stampa** (admin, nella parte società della barra laterale) elenca i modelli della società. **Nuovo modello di stampa** chiede un nome, ciò che il modello **Stampa** (vedi [Fonte di dati](#source)) e un formato di carta, poi apre l’editor. Ha tre colonne:

- A sinistra i **Campi**: tutto ciò che un modello può stampare, in gruppi,
  ognuno con un esempio, e una ricerca. Un clic aggiunge il campo sotto
  l’elemento più in basso; trascinato sul foglio, finisce dove lo si
  rilascia.
  - *Testo libero* stampa un testo fisso, *Codice QR* e *Codice a barre*
    stampano il loro valore come codice, *Linea* traccia una linea.
  - *Tiratore*, *Società e manifestazione* e *Risultato* stampano un valore
    ciascuno, p. es. il nome, il totale o i punti di una serie.
  - *Colpi singoli* aggiunge un’intera fila di caselle per i colpi: *Colpi
    serie 1* è un elemento `{{s1.shot#}}` ripetuto una volta per colpo, a 8
    mm di distanza.
- Al centro il **Foglio** mostra il modello in scala sopra la scansione, con
  valori inventati; **Mostra i nomi dei campi** mostra invece i
  segnaposto. In alto e a sinistra corrono righelli in millimetri, lo zoom
  va dal 50 al 200 %.
- A destra le **Proprietà** modificano con precisione l’elemento
  selezionato: il suo **Valore** (i segnaposto come `{{shooter.name}}`
  vengono compilati in stampa), **x** e **y** in millimetri dall’angolo in
  alto a sinistra, **Larghezza** e **Altezza** (con larghezza 0 il testo
  arriva fino al bordo destro), **Dimensione carattere**, **Grassetto**,
  **Corsivo**, **A capo automatico** (il testo va a capo entro la larghezza,
  ogni riga alta quanto l’altezza) e **Allineamento**. **Ripetizione** lo
  stampa più volte, spostato ogni volta di **dx** e **dy**; `#` in un
  segnaposto conta 1, 2, 3, …. Sotto, il **Modello** imposta il nome, la
  dimensione del foglio, lo spostamento della stampante e la **Scansione del
  modulo vuoto**: un PNG o JPEG dell’intero foglio che appare dietro il
  foglio per posizionare i campi sulle caselle del modulo; non viene mai
  stampato.

**Salva** salva il modello; un punto sul pulsante indica modifiche non salvate, e lasciando la pagina con esse viene chiesta prima una conferma. Con **Salva automaticamente** (attivo finché non lo si disattiva; il browser lo ricorda) le modifiche non salvate vengono salvate ogni minuto. **Anteprima PDF** apre il modello esattamente come lo rende il PDF, che può differire leggermente dal foglio nel browser.

## Fonte di dati {#source}

Ciò che un modello **Stampa** ne determina i campi e dove viene stampato:

| Stampa | Campi | Stampato da |
|---|---|---|
| Risultati di un programma | tiratore, società, manifestazione, serie, colpi singoli, totale, distinzione e il suo valore | i risultati di un programma (vedi [Usare un modello](#use)) o come carta corona della manifestazione |
| Un'iscrizione | tiratore, società, manifestazione, bersaglio, programmi, serie, codice QR dell’iscrizione | come striscia di controllo della manifestazione |
| Un rango | tiratore, società, manifestazione, rango, totale, classifica, classe | diplomi dalle classifiche |

La fonte si può cambiare nell’editor sotto **Modello**; i campi che la nuova
fonte non conosce non stampano nulla.

## Modelli integrati {#builtin}

Standblatt stampa la carta corona e la striscia di controllo con **modelli
integrati**, elencati sotto i modelli della società. **Copia** ne fa un
modello della società che si modifica come ogni altro, p. es. per spostare i
campi su un’altra carta corona prestampata o per aggiungere il logo della
società. Una manifestazione lo sceglie poi in **Configurazione** →
**Moduli**: *Modello delle carte corona* nel modulo *Distinzioni*, *Modello
delle strisce di controllo* nel modulo *Striscia di controllo*. Senza
scelta, o se il modello scelto è eliminato, stampa il modello integrato.

## Disporre gli elementi {#arrange}

- Un elemento si **sposta** trascinandolo e si **ridimensiona** dal suo
  angolo. Il punto rosa sull’ultima casella di un elemento ripetuto regola
  la spaziatura.
- Durante lo spostamento l’elemento si **aggancia** ai bordi e ai centri
  degli altri elementi e del foglio, e una guida rosa mostra la linea;
  altrimenti si aggancia ai millimetri interi. Tenendo premuto Alt lo si
  posiziona liberamente.
- **Selezionarne diversi** con Maiusc o Ctrl e un clic, oppure tracciando un
  riquadro attorno a loro sul foglio vuoto. Trascinandoli si spostano
  insieme; dimensione carattere, grassetto e allineamento valgono allora per
  tutti.
- Un **clic destro** apre il menu: *Duplica*, *Porta in primo piano*,
  *Allineamento* (bordi, centri, distribuire uniformemente), *Centra sul
  foglio* ed *Elimina*.

| Tasto | Effetto |
|---|---|
| Frecce | spostano la selezione di 0.5 mm, con Maiusc di 5 mm |
| Ctrl+Z / Ctrl+Y | annulla / ripeti |
| Ctrl+A | seleziona tutto |
| Ctrl+D | duplica |
| Canc | elimina la selezione |
| Esc | annulla la selezione |
| Ctrl+S | salva |
| Ctrl+rotellina | zoom |

## Immagini {#images}

Loghi, intestazioni e altre immagini vengono dalla biblioteca di immagini
della società, sotto l'elenco dei **Modelli di stampa**: **Carica immagine**
aggiunge un PNG o JPEG di al massimo 2 MB, con il nome del suo
file. Nell'editor, il gruppo **Immagini** dei campi le elenca; un clic o un
trascinamento ne mette una sul foglio, larga 30 mm nelle proporzioni
dell'immagine. Le sue proprietà scelgono un'altra **Immagine**, la posizione
e la dimensione, e **Mantieni le proporzioni** riadatta l'altezza alla
larghezza. Un'immagine si può ripetere come ogni altro elemento.

Un'immagine si può eliminare solo finché nessun modello la
mostra. **Esporta** mette le immagini di un modello nel suo file, e
**Importa** le aggiunge alla biblioteca riusando un'immagine con lo stesso
nome che la società ha già.

## Calibrare la stampante {#calibrate}

Le stampanti raramente stampano esattamente dove dovrebbero. Stampa un
foglio con **Stampa di prova**, misura di quanto i valori sono spostati e
inseriscilo come **Spostamento della stampante a destra** e **Spostamento
della stampante in basso** (negativo per sinistra o alto). Lo spostamento
vale per tutti gli elementi; l’anteprima lo ignora.

## Usare un modello {#use}

In **Configurazione** → **Programmi** della manifestazione apri un programma
e scegli il suo **Modello di stampa** nella sezione omonima; funziona anche
con tiratori già iscritti. Ogni programma concluso offre allora **Stampa sul
modello** nei risultati: al desk, nella pagina **Registra risultati** e
tramite i link per segretari al bersaglio (vedi [Aiutanti](helpers.md)). Con
un agente di stampa e una stampante scelta per *Modelli di stampa* in
**Stand di tiro → Stampa**, il foglio va a quella stampante (p. es. al suo
alimentatore manuale), altrimenti si apre il PDF (vedi
[Stampa](printing.md)).

**Esporta** scarica un modello con le sue immagini come file e **Importa** lo rilegge, p. es. per condividere un modello con un’altra società.
