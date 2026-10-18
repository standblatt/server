# Modelli di stampa

Un modello di stampa porta i risultati su un foglio prestampato, p. es. i
fogli di stand che l’esercito fornisce per il programma obbligatorio: i
colpi si registrano al bersaglio o in seguito, e il foglio viene stampato
invece di essere compilato a mano. Compilare a mano resta possibile; i
modelli sono facoltativi.

## Creare un modello {#create}

**Modelli di stampa** (admin, nella parte società della barra laterale) elenca i modelli della società. **Nuovo modello di stampa** chiede un nome e un formato di carta, poi apre l’editor. Ha tre colonne:

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
  arriva fino al bordo destro), **Dimensione carattere**, **Grassetto** e
  **Allineamento**. **Ripetizione** lo stampa più volte, spostato ogni volta
  di **dx** e **dy**; `#` in un segnaposto conta 1, 2, 3, …. Sotto, il
  **Modello** imposta il nome, la dimensione del foglio, lo spostamento
  della stampante e la **Scansione del modulo vuoto**: un PNG o JPEG
  dell’intero foglio che appare dietro il foglio per posizionare i campi
  sulle caselle del modulo; non viene mai stampato.

**Salva** salva il modello; un punto sul pulsante indica modifiche non salvate, e lasciando la pagina con esse viene chiesta prima una conferma. Con **Salva automaticamente** (attivo finché non lo si disattiva; il browser lo ricorda) le modifiche non salvate vengono salvate ogni minuto. **Anteprima PDF** apre il modello esattamente come lo rende il PDF, che può differire leggermente dal foglio nel browser.

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

## Calibrare la stampante {#calibrate}

Le stampanti raramente stampano esattamente dove dovrebbero. Stampa un
foglio con **Stampa di prova**, misura di quanto i valori sono spostati e
inseriscilo come **Spostamento della stampante a destra** e **Spostamento
della stampante in basso** (negativo per sinistra o alto). Lo spostamento
vale per tutti gli elementi; l’anteprima lo ignora.

## Usare un modello {#use}

In **Configurazione** → **Programmi** della manifestazione scegli il
**Modello di stampa** di un programma; funziona anche con tiratori già
iscritti. Ogni programma concluso offre allora **Stampa sul modello** nei
risultati: al desk, nella pagina **Registra risultati** e tramite i link per
segretari al bersaglio (vedi [Aiutanti](helpers.md)). Con un agente di
stampa e una stampante scelta per *Modelli di stampa* in **Stand di tiro →
Stampa**, il foglio va a quella stampante (p. es. al suo alimentatore
manuale), altrimenti si apre il PDF (vedi [Stampa](printing.md)).

**Esporta** scarica un modello come file e **Importa** lo rilegge, p. es. per condividere un modello con un’altra società.
