# Ranglisten und Kiosk

## Ranglisten

**Ranglisten** zeigt eine Live-Rangliste pro Programm, geordnet nach den [Regeln des SSV](../reference/rules.md#ranking). Die ersten drei Ränge sind hervorgehoben.

Wie Programm-Ranglisten, Wertungen, Vereinskonkurrenz und Sonderpreise
zusammenhängen: [Wie die Ranglisten
zusammenhängen](../explanation/rankings.md).

![Ranglisten](../images/screenshots/rankings.png){ .screenshot }

- **PDF** öffnet die Rangliste (A4) zum Drucken oder Teilen, **CSV** lädt
  sie für Excel herunter. Ranglisten gehen nicht an den Drucker des
  Schiessstands.
- Programme, die für Kiosk-Anzeigen ausgeschaltet sind, werden hier trotzdem
  rangiert und als *nicht im Kiosk* markiert.
- [Wertungen](classifications.md) folgen auf die Programme, eine Tabelle pro
  Klasse; ihr **PDF** und **CSV** ist der Bericht *Wertungen*.
- **Urkunden** druckt Urkunden für die besten Ränge auf eine
  [Druckvorlage](print-forms.md#certificates) des Vereins (siehe [Urkunden
  für die besten Ränge drucken](../how-to/certificates.md)); der Knopf
  erscheint, sobald der Verein eine Vorlage hat, die einen Rang druckt.

## Kiosk-Bildschirme

Ein Kiosk-Link zeigt die Ranglisten eines Anlasses auf einem Bildschirm in
der Halle, ohne Anmeldung; einrichten mit [Die Ranglisten auf einem
Bildschirm in der Halle zeigen](../how-to/kiosk.md).

**Kiosk-Anzeigen** in der Konfiguration (Modul *Kiosk-Anzeigen*) listet die Links. Ein neuer braucht einen Namen, **Dunkel** oder **Hell**, eine Sprache und, bei einem Anlass mit mehreren Schiessständen, einen Schiessstand; **Kiosk-Link erstellen** zeigt den Link nur einmal.

![Kiosk-Bildschirm](../images/screenshots/kiosk.png){ .screenshot }

- Der Kiosk zeigt nur Programme und Wertungen mit eingeschaltetem *Auf
  Kiosk-Anzeigen zeigen* und, mit dem Modul *Scheibenzuteilung*, den
  Scheibenstatus, wenn *Scheibenstatus auf den Ranglisten zeigen* unter
  **Kiosk-Anzeigen** an ist.
- Der Link legt Design, Sprache und Schiessstand fest (siehe
  [URL-Parameter](../reference/links.md#parameters)).
- Der Knopf **Vollbild** oben rechts schaltet den Bildschirm ins Vollbild.
- Der Bildschirm aktualisiert sich live und verbindet sich nach einem
  Netzunterbruch selbst wieder. Bei einer neuen Version von Standblatt lädt
  er sich selbst neu.
- **widerrufen** in der Liste stoppt einen Link sofort. *zuletzt gesehen*
  sagt, wann ein Bildschirm zuletzt Daten geholt hat.
