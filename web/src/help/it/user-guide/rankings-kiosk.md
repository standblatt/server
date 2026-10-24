# Classifiche e chiosco

## Classifiche

**Classifiche** mostra una classifica in diretta per programma, ordinata secondo le [regole della FST](../reference/rules.md#ranking). I primi tre ranghi sono evidenziati.

Come si collegano classifiche dei programmi, classifiche generali, concorso
di sezione e premi speciali: [Come si combinano le
classifiche](../explanation/rankings.md).

![Classifiche](../images/screenshots/rankings.png){ .screenshot }

- **PDF** apre la classifica (A4) per stamparla o condividerla, **CSV** la
  scarica per Excel. Le classifiche non vanno alla stampante dello stand.
- I programmi disattivati per gli schermi chiosco sono comunque classificati
  qui e segnati *non sul chiosco*.
- Le [classifiche generali](classifications.md) seguono i programmi, una
  tabella per classe; il loro **PDF** e **CSV** sono il rapporto
  *Classifiche generali*.
- **Diplomi** stampa diplomi per i migliori ranghi su un [modello di
  stampa](print-forms.md#certificates) della società (vedi [Stampare i
  diplomi per i migliori ranghi](../how-to/certificates.md)); il pulsante
  compare non appena la società ha un modello che stampa un rango.

## Schermi chiosco

Un link chiosco mostra le classifiche di una manifestazione su uno schermo
della sala, senza accesso; configuralo con [Mostrare le classifiche su uno
schermo nella sala](../how-to/kiosk.md).

**Display chiosco** nella configurazione (modulo *Display chiosco*) elenca i link. Uno nuovo riceve un nome, **Scuro** o **Chiaro**, una lingua e, in una manifestazione con più stand, uno stand; **Crea link chiosco** mostra il link una sola volta.

![Schermo chiosco](../images/screenshots/kiosk.png){ .screenshot }

- Il chiosco mostra solo i programmi e le classifiche generali con *Mostra
  sugli schermi chiosco* attivo, e, con il modulo *Assegnazione dei
  bersagli*, lo stato dei bersagli se *Mostra lo stato dei bersagli sulle
  classifiche* è attivo in **Display chiosco**.
- Il link imposta tema, lingua e stand (vedi [Parametri
  dell’URL](../reference/links.md#parameters)).
- Il pulsante **Schermo intero** in alto a destra porta lo schermo a schermo
  intero.
- Lo schermo si aggiorna in diretta e si ricollega da solo dopo
  un’interruzione di rete. A una nuova versione di Standblatt si ricarica da
  solo.
- **revoca** nell’elenco ferma subito un link. *visto l’ultima volta* indica
  quando uno schermo ha preso dati l’ultima volta.
