# Online-Reservation

Mit der **Online-Reservation** fragen externe Vereine selbst Scheiben
an. Der Organisator setzt einen Link in die Einladung; ein Verein wählt die
gewünschten Zeitblöcke, gibt an, wie viele Scheiben er jeweils braucht, und
meldet seine Schützen. Ein Admin bestätigt oder lehnt jede Anfrage ab.

## Einschalten {#setup}

Schalte das Modul *Online-Reservation* unter **Konfiguration** → **Module**
ein (es braucht einen Stand); dann genügt auf **Online-Reservation** in der
Konfiguration **Einschalten**. Die Seite zeigt dann den Link mit einem
QR-Code für die Einladung und die angebotenen Zeitblöcke.

![Online-Reservation in der
Konfiguration](../images/screenshots/booking-setup.png){ .screenshot }

Ohne Änderungen ist jeder Tag des Anlasses während seiner Schiesszeiten
offen (siehe [Anlässe](events.md#settings)), ohne Schiesszeiten von 08:00
bis 17:00, in Zeitblöcken von 30 Minuten, und Anfragen sind bis 14 Tage vor
dem ersten Tag möglich. **Anpassen** ändert das:

- **Schiessen von** und **Schiessen bis**: die Zeiten jedes Tages ohne
  Schiesszeiten.
- **Pause von** und **Pause bis**: eine optionale Pause, z. B. über Mittag.
- **Länge eines Zeitblocks (Minuten)**: wie lange ein Zeitblock dauert.
- **Scheiben pro Zeitblock und Verein**: eine Grenze pro Anfrage; leer
  heisst keine Grenze.
- **Anmeldeschluss**: der letzte Tag, an dem Vereine anfragen oder ihre
  Anfrage ändern können.
- **Tage mit anderen Zeiten**: eine Zeile pro Tag des Anlasses,
  z. B. Freitag 17:00–20:00. Die Tage kommen aus dem Anlass und folgen ihm,
  wenn du einen hinzufügst oder entfernst.

**Neuer Link** ersetzt den Link; der alte funktioniert nicht mehr. **Ausschalten** schliesst die Online-Reservation; Anfragen und bestätigte Reservationen bleiben. *Für nächstes Jahr kopieren* nimmt die Zeiten mit, aber weder den Link noch den Anmeldeschluss.

## Was ein Verein sieht {#club}

Der Link öffnet eine Seite ohne Login. Sie zeigt den Anlass, den
Anmeldeschluss und jeden Zeitblock mit den noch freien Scheiben pro
Stand. Nur bestätigte Reservationen zählen; Namen anderer Vereine werden nie
angezeigt.

![Die Reservationsseite eines
Vereins](../images/screenshots/booking-public.png){ .screenshot }

Der Verein gibt die Anzahl Scheiben pro Zeitblock ein, seinen Namen, eine
Kontaktperson, E-Mail und Telefon sowie seine **Schützen**: Name, Vorname,
Geburtsdatum und optional die Lizenz. **Liste einfügen** übernimmt Zeilen
aus Excel (Name, Vorname, Geburtsdatum, Lizenz). Die Schützenliste kann
später folgen; die Programme werden wie gewohnt am Stand gewählt.

Nach **Anfrage senden** erhält der Verein eine Mail mit einem persönlichen
Link. Darüber kann er bis zum Anmeldeschluss den Status sehen, die Anfrage
ändern oder stornieren. Ändert er die Zeitblöcke einer bestätigten Anfrage,
werden ihre Scheiben frei und die Anfrage wartet auf eine neue Bestätigung;
neue Schützen lassen sie bestätigt.

## Anfragen bestätigen {#confirm}

Jeder Admin mit E-Mail erhält zu jeder neuen, geänderten oder stornierten
Anfrage eine Mail. Die Anfragen stehen zuoberst unter **Reservationen**; die
Seitenleiste zeigt, wie viele warten.

![Reservationsanfragen auf der Seite
Reservationen](../images/screenshots/booking-requests.png){ .screenshot }

- **Bestätigen** gibt jedem Zeitblock freie Scheiben seines Standes, wenn
  möglich nebeneinander, als gewöhnliche Reservation mit dem Namen des
  Vereins, und schickt dem Verein die Scheibennummern per Mail. Hat ein
  Zeitblock zu wenige freie Scheiben, wird nichts reserviert und die Meldung
  nennt den Zeitblock; verschiebe eine Reservation oder lehne ab.
- **Bestätigen** fragt zuerst, zu welchem **Verein** der Stammdaten die
  Anfrage gehört, weil der Verein seinen Namen selbst eingetippt hat: Der
  Verein mit gleichem Namen ist schon ausgewählt, Vereine mit ähnlichem
  Namen werden vorgeschlagen, und *Neuer Verein* legt ihn an. Die Anfrage
  zeigt danach den zugeordneten Verein.
- **Bestätigen** übernimmt auch die Schützen in die Stammdaten, in diesen
  Verein: ein Schütze wird über die Lizenz gefunden, sonst über Name und
  Geburtsdatum, und andernfalls angelegt. Bestehende Schützen werden nie
  überschrieben, weil die Liste aus einem öffentlichen Formular stammt; nur
  ein fehlendes Geburtsdatum oder eine fehlende Lizenz wird
  ergänzt. Schützen, die nach der Bestätigung dazukommen, übernimmt
  **Schützen übernehmen**.
- **Ablehnen** fragt nach einer optionalen Begründung, die der Verein per
  Mail erhält. Scheiben, die die Anfrage hatte, sind wieder frei.

Die reservierten Scheiben sind gewöhnliche [Reservationen](reservations.md):
ändere dort bei Bedarf die Scheibennummern. Am Anlass findet die Anmeldung
die Schützen mit der normalen Suche. Der Bericht *Scheibenreservationen*
unter [Statistik](statistics-reports.md) listet jeden Verein mit seinen
Zeitblöcken, Scheiben und Schützen.
