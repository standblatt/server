# Collegare una stampante allo stand

Ti serve quando fogli di stand, strisce di controllo e carte corona devono
uscire dalla stampante allo stand invece di aprirsi come PDF. Richiede il
permesso *Stand e stampanti* (vedi [Ruoli](../user-guide/users.md#roles)).

Un computer accanto alle stampanti esegue l’agente di stampa: un PC Windows
o un Mac con l’app Standblatt Print Agent, o un computer Linux come un
Raspberry Pi con l’agente senza interfaccia. Perché funziona così è spiegato
in [Come funziona la stampa](../explanation/printing.md).

## Con l’app su Windows o su un Mac {#app}

1. Scaricala dall’[ultima
   versione](https://github.com/standblatt/print-agent/releases/latest):
   `Standblatt-Print-Agent-Setup.exe` per Windows,
   `standblatt-print-agent-app-darwin-arm64.dmg` per un Mac con chip Apple
   (`-amd64.dmg` per un Mac Intel), il `.deb` o il `.rpm` per Linux.
2. Installala. L’app non è ancora firmata: Windows chiede una volta,
   **Ulteriori informazioni → Esegui comunque**; su un Mac trascinala in
   *Applicazioni* e aprila la prima volta con un clic destro → **Apri**. Su
   Windows si installa per l’utente attuale senza diritti di amministratore,
   con un collegamento nel menu Start e sul desktop, e si avvia con Windows.
3. Nell’app clicca **Collega uno stand di tiro**. Mostra un codice come
   `K7QD-4MZP` e apre **Collegare un agente di stampa** nel browser.
4. Accedi se richiesto, controlla il computer indicato sulla pagina, scegli
   la **Società**, lo **Stand di tiro** e il **Nome** dell’agente, e
   **Collega**. Il codice è valido dieci minuti.

![Un codice nell’app](../images/app/print-agent-code.png){ .screenshot }

Se il browser di quel computer non è a portata di mano, apri **Stand di tiro
→ Configura le stampanti** dello stand su un qualsiasi altro computer,
digita il codice nel campo accanto a **Collega con codice** e cliccalo: il
computer viene collegato a quello stand con il suo nome.

Per un secondo stand sullo stesso computer, p. es. *Pistola 25 m* e *Pistola
50 m*, clicca **Collega un altro stand di tiro** nell’app e ripeti i passi 3
e 4.

## Con l’agente senza interfaccia su Linux {#headless}

1. Apri **Stand di tiro** e clicca **Configura le stampanti** dello stand.
2. In **Agenti di stampa** inserisci un nome (p. es. *Pi hall 1*) e clicca
   **Aggiungi agente di stampa**.
3. Clicca **copia** per copiare le due righe `STANDBLATT_URL` e
   `STANDBLATT_AGENT_TOKEN`. Il token viene mostrato solo ora.
4. Installa e avvia l’agente con queste impostazioni come descritto nel
   [README dell’agente di
   stampa](https://github.com/standblatt/print-agent#readme).

Per un secondo stand sullo stesso computer, aggiungi anche lì un agente e
dai all’agente entrambi i token, separati da una virgola.

## Scegliere le stampanti {#printers}

1. Apri **Stand di tiro → Configura le stampanti** dello stand. L’agente
   mostra *collegato* con la sua versione e il sistema operativo; le sue
   stampanti compaiono appena si è collegato.
2. Per ogni documento scegli la stampante, poi le sue opzioni (cassetto,
   carta, fronte-retro …, come l’agente le riporta per quella stampante; su
   Linux altre opzioni `lp` come `fit-to-page` vanno nel campo di testo).
3. Clicca **Pagina di prova** nella riga: una piccola pagina (90 × 29 mm)
   con società, stand, stampante, opzioni e ora viene stampata con le
   opzioni di quella riga. Una stampante di etichette la stampa come
   etichetta, le altre stampanti nell’angolo in alto a sinistra del foglio.
4. Clicca **Salva** nella riga. Un documento lasciato su *browser (PDF)*
   continua ad aprirsi come PDF.

![Stampa](../images/screenshots/printing.png){ .screenshot }

*Diplomi* sono i diplomi dalle classifiche, *Modelli di stampa (prestampati)* ogni altro [modello di stampa](../user-guide/print-forms.md) della società. Le classifiche non vengono mai inviate a una stampante.

## Controllare che stampi {#jobs}

Stampa qualcosa, p. es. un foglio di stand alla prossima iscrizione, e
guarda **Lavori recenti** sulla stessa pagina: *fatto* significa che è stato
stampato, *fallito* mostra l’errore accanto. Un lavoro che l’agente non
conferma entro due minuti viene ritentato; dopo tre tentativi fallisce.

Per dismettere un computer, clicca **revoca** accanto al suo agente; con
l’app, rimuovi anche lo stand in **Avanzate** e disinstalla l’app. Tutte le
schermate della pagina di stampa sono descritte in
[Stampa](../user-guide/printing.md).
