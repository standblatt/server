# Stampare i risultati sui fogli di stand dell’esercito

Ti serve quando i risultati dei tiratori devono finire su fogli prestampati,
p. es. i fogli di stand che l’esercito fornisce per il programma
obbligatorio, e preferisci stamparli piuttosto che compilarli a
mano. Richiede i permessi *Modelli di stampa*, *Configurare le
manifestazioni* e *Stand e stampanti* (vedi
[Ruoli](../user-guide/users.md#roles)).

Prima di iniziare, scansiona un foglio vuoto come PNG o JPEG e tieni alcuni
fogli di riserva per le stampe di prova.

## 1. Preparare il programma {#programme}

1. Controlla che **Moduli della società** nella parte della società della
   barra laterale abbia *Modelli di stampa* attivato (è attivo a meno che
   qualcuno non l’abbia disattivato).
2. In **Configurazione** → **Programmi** della manifestazione, apri il
   programma e spunta **Colpi singoli** per ogni serie i cui colpi hanno
   caselle sul foglio. Fallo prima della prima iscrizione: appena ci sono
   tiratori iscritti, le serie sono fissate.
3. **Salva** il programma.

## 2. Creare il modello {#form}

1. Apri **Modelli di stampa** nella parte della società della barra laterale
   e clicca **Nuovo modello di stampa**.
2. Inserisci un nome (p. es. *Programma obbligatorio 300 m*), scegli
   **Stampa** *Risultati di un programma* e il formato carta più vicino al
   foglio, poi **Crea**. Si apre l’editor.
3. In **Modello** a destra, imposta **Larghezza (mm)** e **Altezza (mm)**
   alla misura del foglio, poi **Carica una scansione** in **Scansione del
   modulo vuoto**. La scansione appare dietro il foglio durante la modifica;
   non viene mai stampata.
4. Da **Campi** a sinistra, trascina i campi del tiratore
   (p. es. **Cognome**, **Nome**, **Anno di nascita**, **Numero di
   licenza**) sulle loro caselle nella scansione.
5. Dal gruppo **Colpi singoli**, trascina **Colpi serie 1** sulla prima
   casella dei colpi di quella serie. È una fila di caselle a 8 mm di
   distanza; trascina il punto rosa sulla sua ultima casella finché ogni
   valore sta nella sua casella, o imposta **dx (mm)** in
   **Proprietà**. Ripeti per le altre serie.
6. Da **Risultato**, aggiungi i totali delle serie (**Serie 1** …) e il
   **Totale** dove il foglio ha caselle per essi.
7. **Salva**.

![Modelli di stampa](../images/screenshots/print-forms.png){ .screenshot }

Come spostare, allineare e agganciare gli elementi è in [Modelli di
stampa](../user-guide/print-forms.md#arrange), ogni campo in [Segnaposto dei
modelli](../reference/form-placeholders.md).

## 3. Calibrare la stampante {#calibrate}

1. Metti un foglio vuoto nella stampante allo stand.
2. Clicca **Stampa di prova**: il modello salvato si apre come PDF con
   valori inventati e senza la scansione. Stampalo in dimensioni reali
   (senza adattare alla pagina).
3. Misura di quanto i valori sono spostati rispetto alle loro caselle e
   inseriscilo come **Spostamento della stampante a destra (mm)** e
   **Spostamento della stampante in basso (mm)** in **Modello** (negativo
   per spostare a sinistra o in alto).
4. **Salva** e ripeti dal passo 1 finché i valori stanno nelle caselle.

## 4. Assegnare il modello e la stampante {#assign}

1. In **Configurazione** → **Programmi** della manifestazione, apri il
   programma e scegli il modello nella sezione **Modello di stampa**. Viene
   salvato subito; funziona anche con tiratori già iscritti.
2. In **Stand di tiro → Configura le stampanti** dello stand, scegli la
   stampante per **Modelli di stampa (prestampati)**, p. es. con la sua
   alimentazione manuale come cassetto, e **Salva** (vedi [Collegare una
   stampante](connect-printer.md)).

## 5. Stampare i fogli {#print}

Appena tutte le serie del programma sono registrate, i suoi risultati
mostrano **Stampa sul modello**: nei risultati all’accoglienza e sui link
per marcatori al bersaglio (vedi [Aiutanti](../user-guide/helpers.md)). In
**Registra risultati**, ogni riga sotto *Appena registrati* ha il proprio
pulsante di stampa. Metti un foglio nella stampante e cliccalo. Senza un
agente di stampa online, il PDF si apre e viene stampato dal browser.
