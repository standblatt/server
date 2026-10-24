# Die Ranglisten auf einem Bildschirm in der Halle zeigen

Das brauchst du, wenn Schützen und Zuschauer die Live-Ranglisten auf einem
Bildschirm am Schiessstand verfolgen sollen. Der Bildschirm braucht einen
Browser und Internet, keine Anmeldung. Für das Einrichten braucht es das
Recht *Anlässe konfigurieren* (siehe
[Rollen](../user-guide/users.md#roles)).

1. Schalte unter **Konfiguration** → **Module** des Anlasses
   *Kiosk-Anzeigen* ein.
2. Öffne **Kiosk-Anzeigen** in der Konfiguration.
3. Gib einen Namen für den Bildschirm ein (z. B. *Halle 1*), wähle
   **Dunkel** oder **Hell** und, wenn der Bildschirm nicht die Sprache des
   Anlasses verwenden soll, eine Sprache. Bei einem Anlass mit mehreren
   Schiessständen wähle den Schiessstand, den der Bildschirm zeigt, oder
   *Alle Stände*.
4. Klicke auf **Kiosk-Link erstellen** und kopiere den Link mit
   **Kopieren**. Er wird nur einmal angezeigt.
5. Öffne den Link in einem Browser auf dem Bildschirm und klicke oben rechts
   auf den Knopf **Vollbild**.

![Kiosk-Bildschirm](../images/screenshots/kiosk.png){ .screenshot }

Der Bildschirm aktualisiert sich live und verbindet sich nach einem
Netzunterbruch selbst wieder.

## Wählen, was der Bildschirm zeigt {#content}

- Ein Programm erscheint nur, solange **Auf Kiosk-Anzeigen zeigen** auf
  seiner Programmseite an ist (bei neuen Programmen ist es an). Seiten des
  Personals, PDF und CSV rangieren es trotzdem.
- Eine Wertung erscheint mit **Auf Kiosk-Anzeigen zeigen** in ihrem Editor
  (siehe [Wertungen](../user-guide/classifications.md)).
- Mit Scheibenzuteilung fügt **Scheibenstatus auf den Ranglisten zeigen**
  unter **Kiosk-Anzeigen** den Scheibenstatus über den Ranglisten hinzu.

## Einen Bildschirm stoppen {#revoke}

Klicke in der Liste unter **Kiosk-Anzeigen** neben dem Bildschirm auf
**widerrufen**: Der Link funktioniert sofort nicht mehr. *zuletzt gesehen*
sagt, wann ein Bildschirm zuletzt Daten geholt hat. Die Teile des Links
stehen unter [Links ohne Anmeldung](../reference/links.md).
