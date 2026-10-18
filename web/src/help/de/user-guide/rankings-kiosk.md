# Ranglisten und Kiosk

## Ranglisten

**Ranglisten** zeigt eine Live-Rangliste pro Programm, geordnet nach den [Regeln des SSV](../reference/rules.md#ranking). Die ersten drei Ränge sind hervorgehoben.

![Ranglisten](../images/screenshots/rankings.png){ .screenshot }

- **PDF** öffnet die Rangliste (A4) zum Drucken oder Teilen, **CSV** lädt
  sie für Excel herunter. Ranglisten gehen nicht an den Drucker des
  Schiessstands.
- Programme, die für Kiosk-Anzeigen ausgeschaltet sind, werden hier trotzdem
  rangiert und als *nicht im Kiosk* markiert.
- [Wertungen](classifications.md) folgen auf die Programme, eine Tabelle pro
  Klasse; ihr **PDF** und **CSV** ist der Bericht *Wertungen*.

## Kiosk-Bildschirme

Ein Kiosk-Link zeigt die Ranglisten eines Anlasses auf einem Bildschirm in
der Halle, ohne Anmeldung. Erstelle einen im Tab *Kiosk-Anzeigen* der
**Einrichtung** des Anlasses:

1. Gib einen Namen für den Bildschirm ein (z. B. *Halle 1*), wähle
   **Dunkel** oder **Hell** und wahlweise eine Sprache.
2. **Kiosk-Link erstellen**; kopiere den Link. Er wird nur einmal angezeigt.
3. Öffne den Link in einem Browser auf dem Bildschirm, am besten im
   Vollbild.

![Kiosk-Bildschirm](../images/screenshots/kiosk.png){ .screenshot }

- Der Kiosk zeigt nur Programme mit eingeschaltetem *Auf Kiosk-Anzeigen
  zeigen*, und den Scheibenstatus, wenn die Einstellung *Scheibenstatus auf
  den Ranglisten zeigen* des Anlasses an ist.
- Die Sprache ist die Sprache des Anlasses, sofern der Link keine andere
  trägt (`?language=de`, `fr`, `it`, `en`); das Farbschema kommt von
  `?theme=dark` oder `?theme=light`. Bei einem Anlass auf mehreren Ständen
  beschränkt `?range=<id>` (beim Erstellen des Links gewählt) die Anzeige
  auf diesen Stand: seine Scheiben, seine Programme und die Programme ohne
  Stand.
- Der Bildschirm aktualisiert sich live und verbindet sich nach einem
  Netzunterbruch selbst wieder. Bei einer neuen Version von Standblatt lädt
  er sich selbst neu.
- **widerrufen** in der Liste stoppt einen Link sofort. *zuletzt gesehen*
  sagt, wann ein Bildschirm zuletzt Daten geholt hat.
