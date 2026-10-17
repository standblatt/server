# Modelli di stampa

Un modello di stampa porta i risultati su un foglio prestampato, p. es. i
fogli di stand che l’esercito fornisce per il programma obbligatorio: i
colpi si registrano al bersaglio o in seguito, e il foglio viene stampato
invece di essere compilato a mano. Compilare a mano resta possibile; i
modelli sono facoltativi.

## Creare un modello {#create}

**Modelli di stampa** (admin, nella parte società della barra laterale) elenca i modelli della società. **Nuovo modello di stampa** chiede nome e formato della carta, poi apre l’editor:

- **Larghezza** e **Altezza** del foglio in millimetri.
- **Scansione del modulo vuoto**: un PNG o JPEG dell’intero foglio. Appare
  dietro l’**Anteprima** per posizionare i campi sulle caselle del modulo;
  non viene mai stampata.
- **Elementi**: cosa viene stampato e dove. **x** e **y** sono millimetri
  dall’angolo superiore sinistro, **Larghezza** e **Altezza** la dimensione
  (una larghezza di 0 lascia il testo fino al bordo destro), **Dimensione**
  la dimensione del carattere in punti.
  - *Testo* stampa il suo **Valore**, in cui i segnaposto vengono
    sostituiti, p. es. `{{shooter.name}}` o `{{total}}`.
  - *QR* e *Codice a barre* stampano il loro valore come codice, *Linea*
    traccia una linea da x/y per larghezza/altezza.
  - **Ripetizione** stampa l’elemento più volte, spostato ogni volta di
    **dx** e **dy**; `#` in un segnaposto conta 1, 2, 3, …. Un elemento
    `{{s1.shot#}}` con ripetizione 10 e dx 8 riempie le dieci caselle dei
    colpi della serie 1.
- **Segnaposto** elenca ogni valore con un esempio; un clic lo aggiunge
  all’ultimo elemento modificato.

Il **Foglio** mostra il modello in scala sopra la scansione, con valori
inventati. Trascina un elemento per spostarlo e il suo angolo per
ridimensionarlo; le frecce spostano l’elemento scelto di mezzo millimetro,
con Maiusc di cinque, e un doppio clic sul foglio vuoto vi aggiunge un
testo. Le posizioni si agganciano al mezzo millimetro. Tabella e foglio
mostrano gli stessi elementi: scegline uno sul foglio e modifica i suoi
valori esatti nella tabella, o viceversa.

**Salva** salva il modello. **Anteprima PDF** mostra il modello esattamente come lo rende il PDF, che può differire leggermente dal foglio nel browser.

## Calibrare la stampante {#calibrate}

Le stampanti raramente stampano esattamente dove dovrebbero. Stampa un
foglio con **Stampa di prova**, misura di quanto i valori sono spostati e
inseriscilo come **Spostamento della stampante a destra** e **Spostamento
della stampante in basso** (negativo per sinistra o alto). Lo spostamento
vale per tutti gli elementi; l’anteprima lo ignora.

## Usare un modello {#use}

Nella scheda *Programmi* della **Configurazione** della manifestazione
scegli il **Modello di stampa** di un programma; funziona anche con tiratori
già iscritti. Ogni programma concluso offre allora **Stampa sul modello**
nei risultati: al desk, nella pagina **Registra risultati** e tramite i link
per segretari al bersaglio (vedi [Aiutanti](helpers.md)). Con un agente di
stampa e una stampante scelta per *Modelli di stampa* in **Stand di tiro →
Stampa**, il foglio va a quella stampante (p. es. al suo alimentatore
manuale), altrimenti si apre il PDF (vedi [Stampa](printing.md)).

**Esporta** scarica un modello come file e **Importa** lo rilegge, p. es. per condividere un modello con un’altra società.
