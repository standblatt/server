# Links ohne Anmeldung

Diese Links öffnen eine Seite ohne Anmeldung. Wer den Link hat, kann ihn
nutzen, schicke ihn deshalb nur den Personen, für die er gedacht ist.

| Link | Erstellt mit | Zeigt | Gültig |
|---|---|---|---|
| Kiosk-Link `/kiosk/<secret>` | **Kiosk-Link erstellen** unter **Kiosk-Anzeigen** in der Konfiguration | die Live-Ranglisten des Anlasses: Programme und Wertungen mit *Auf Kiosk-Anzeigen zeigen* und den Scheibenstatus mit *Scheibenstatus auf den Ranglisten zeigen* | bis **widerrufen** oder bis der Anlass gelöscht wird |
| Helferlink `/h/<secret>` | **Helferlink erstellen** unter **Helferlinks** in der Konfiguration | einen reduzierten Desk des Anlasses für eine [Helferrolle](permissions.md#helper-links) | bis **widerrufen** oder bis der Anlass gelöscht wird |
| Reservationslink `/b/<code>` | **Einschalten** unter **Reservationen** in der Konfiguration | die Zeitblöcke des Anlasses mit den freien Scheiben pro Schiessstand und das Anfrageformular | bis **Neuer Link** oder **Ausschalten**; Anfragen nur bis zur Frist |
| Verwaltungslink der Reservation `/b/manage/<secret>` | dem Verein nach **Anfrage senden** gemailt | den Status der Anfrage des Vereins; ändern oder stornieren | Ändern und Stornieren bis zur Frist |
| Einladungslink `/invite/<secret>` | **Per E-Mail einladen** unter **Benutzer** | den Verein und die Rolle; annehmen mit einem Konto oder einem neuen Passwort | sieben Tage; endet beim Annehmen, Zurückziehen oder erneuten Senden |
| Verbindungsseite `/connect?code=<code>` | von der App Standblatt Print Agent geöffnet | den Computer des Codes; Verein, Schiessstand und Name des Agents | der Code: zehn Minuten; braucht eine Anmeldung mit *Stände und Drucker* |
| E-Mail-Bestätigung `/verify?token=…` | Hinzufügen oder Ändern der E-Mail-Adresse, Registrierung | bestätigt die Adresse | sieben Tage |
| Passwort-Link `/reset-password?token=…` | **Passwort vergessen?**; der Betreiber im Admin-Panel | setzt ein neues Passwort | eine Stunde; vom Betreiber sieben Tage |

Der geheime Teil von Kiosk-, Helfer-, Verwaltungs- und Einladungslinks wird
einmal angezeigt oder gemailt und nur als Hash gespeichert: Ein verlorener
Link lässt sich nicht mehr anzeigen, nur ersetzen.

## URL-Parameter {#parameters}

| Parameter | Link | Wirkung |
|---|---|---|
| `?theme=dark`, `?theme=light` | Kiosk | dunkle oder helle Anzeige |
| `?language=de` (`fr`, `it`, `en`) | Kiosk, Helfer | Sprache der Seite statt der Sprache des Anlasses |
| `?range=<id>` | Kiosk | bei einem Anlass mit mehreren Schiessständen: nur die Scheiben dieses Schiessstands, seine Programme und die Programme ohne Schiessstand |
| `?code=` | Verbinden | der zu bestätigende Code, von der App ausgefüllt |

**Kiosk-Link erstellen** schreibt das dort gewählte Design, die Sprache und den Schiessstand in den Link.
