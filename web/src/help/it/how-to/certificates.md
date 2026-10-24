# Stampare i diplomi per i migliori ranghi

Ti serve quando i vincitori di un programma o di una classifica generale
ricevono un diploma con il loro nome e rango, p. es. i ranghi da 1 a 3 alla
premiazione. Creare il modello richiede il permesso *Modelli di stampa*,
stampare *Gestire le manifestazioni* (vedi
[Ruoli](../user-guide/users.md#roles)).

## 1. Creare il modello del diploma {#form}

1. Apri **Modelli di stampa** nella parte della società della barra laterale
   e clicca **Nuovo modello di stampa**.
2. Inserisci un nome (p. es. *Diploma*), scegli **Stampa** *Un rango* e un
   formato carta come *A4 verticale (210 × 297 mm)*, poi **Crea**.
3. Da **Campi**, aggiungi ciò che mostra il diploma: **Nome** da *Tiratore*,
   **Rango**, **Classifica** e **Classe** da *Rango*, **Totale** da
   *Risultato*, e **Testo libero** per parole fisse come *Diploma*.
4. Per un logo, prima **Carica immagine** nella biblioteca di immagini sotto
   l’elenco dei **Modelli di stampa**; compare poi nel gruppo **Immagini**
   dei campi (vedi [Immagini](../user-guide/print-forms.md#images)).
5. **Salva**, poi controlla il risultato con **Anteprima PDF**.

Su carta con un disegno prestampato, aggiungine una scansione in **Scansione
del modulo vuoto** e calibra la stampante come per i [fogli di stand
dell’esercito](army-sheet.md#calibrate).

## 2. Scegliere la stampante {#printer}

In **Stand di tiro → Configura le stampanti** del primo stand della
manifestazione, scegli la stampante per **Diplomi** e **Salva** (vedi
[Collegare una stampante](connect-printer.md)). Senza stampante, i diplomi
si aprono come un unico PDF.

## 3. Stampare i diplomi {#print}

1. Apri **Classifiche** della manifestazione.
2. Clicca **Diplomi** nell’intestazione del programma o della classifica
   generale. Il pulsante compare non appena la società ha un modello che
   stampa un rango.
3. Scegli il **Modello di stampa**, per una classifica generale con più
   classi la **Classe**, e i ranghi con **Dal rango** e **Al rango** (da 1 a
   3 se non modificati).
4. Clicca **Stampa i diplomi**, o **PDF** per aprirli.

![Diplomi da una classifica](../images/screenshots/certificates.png){
.screenshot }

Si stampa una pagina per ogni tiratore classificato. Il diploma di un
programma conosce anche i risultati del tiratore, p. es. le serie.
