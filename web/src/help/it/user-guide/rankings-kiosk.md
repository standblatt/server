# Classifiche e chiosco

## Classifiche

**Classifiche** mostra una classifica in diretta per programma, ordinata secondo le [regole della FST](../reference/rules.md#ranking). I primi tre ranghi sono evidenziati.

![Classifiche](../images/screenshots/rankings.png){ .screenshot }

- **PDF** apre la classifica (A4) per stamparla o condividerla, **CSV** la
  scarica per Excel. Le classifiche non vanno alla stampante dello stand.
- I programmi disattivati per gli schermi chiosco sono comunque classificati
  qui e segnati *non sul chiosco*.
- Le [classifiche generali](classifications.md) seguono i programmi, una
  tabella per classe; il loro **PDF** e **CSV** sono il rapporto
  *Classifiche generali*.

## Schermi chiosco

Un link chiosco mostra le classifiche di una manifestazione su uno schermo
della sala, senza accesso. Creane uno nella scheda *Display chiosco* della
**Configurazione** della manifestazione:

1. Inserisci un nome per lo schermo (p. es. *Sala 1*), scegli **Scuro** o
   **Chiaro** e, se vuoi, una lingua.
2. **Crea link chiosco**; copia il link. Viene mostrato una sola volta.
3. Apri il link in un browser sullo schermo, idealmente a schermo intero.

![Schermo chiosco](../images/screenshots/kiosk.png){ .screenshot }

- Il chiosco mostra solo i programmi con *Mostra sugli schermi chiosco*
  attivo, e lo stato dei bersagli se l’impostazione *Mostra lo stato dei
  bersagli sulle classifiche* della manifestazione è attiva.
- La lingua è quella della manifestazione, salvo che il link ne indichi
  un’altra (`?language=de`, `fr`, `it`, `en`); il tema viene da
  `?theme=dark` o `?theme=light`. Per una manifestazione su più stand,
  `?range=<id>` (scelto alla creazione del link) limita la visualizzazione a
  quello stand: i suoi bersagli, i suoi programmi e i programmi senza stand.
- Lo schermo si aggiorna in diretta e si ricollega da solo dopo
  un’interruzione di rete. A una nuova versione di Standblatt si ricarica da
  solo.
- **revoca** nell’elenco ferma subito un link. *visto l’ultima volta* indica
  quando uno schermo ha preso dati l’ultima volta.
