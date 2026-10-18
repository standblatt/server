# Benutzer

**Benutzer** listet die Konten mit Zugang zum Verein, die offenen Einladungen und die Rollen des Vereins. Die Seite steht Rollen mit dem Recht *Benutzer und Verein* offen, also mindestens den Admins.

![Benutzer](../images/screenshots/users.png){ .screenshot }

## Mitglieder und Einladungen {#invitations}

- **Per E-Mail einladen**: eine oder viele Adressen einfügen (durch Kommas
  getrennt oder eine pro Zeile, bis 300) und die Rolle wählen. Jede Adresse
  erhält eine E-Mail mit einem Link, der sieben Tage gültig ist; die Liste
  unter dem Formular zeigt, was mit jeder Adresse geschehen ist. Die Mails
  gehen mit etwa vier pro Sekunde raus, eine grosse Liste dauert also eine
  Minute.
- **Annehmen**: Der Link öffnet eine Seite, die sagt, welcher Verein mit
  welcher Rolle einlädt. Ohne Konto wählt die Person ein Passwort und ist
  sofort angemeldet; die E-Mail-Adresse ist das Login. Wer schon ein Konto
  mit der Adresse hat, meldet sich zuerst an und nimmt dann an. Erst danach
  erscheint die Person bei den Mitgliedern.
- **Einladungen** unter den Mitgliedern listet die offenen und die
  abgelaufenen: **Erneut senden** schickt einen neuen Link (der alte
  funktioniert nicht mehr), **Zurückziehen** storniert die Einladung.
- **Benutzer hinzufügen**: für ein Konto ohne E-Mail, etwa ein gemeinsames
  Konto am Stand: ein bestehendes Konto wird über seinen Benutzernamen
  hinzugefügt; für ein neues Konto zusätzlich ein Passwort eingeben
  (mindestens 10 Zeichen).
- **Rolle**: die Rolle eines Mitglieds in der Liste wählen. Ein Verein
  behält immer mindestens einen Admin.
- **Entfernen** nimmt den Zugang zu diesem Verein weg; das Konto selbst
  bleibt und behält seine anderen Vereine.

## Rollen {#roles}

Was ein Mitglied darf, hängt von seiner Rolle ab. Die Matrix unter
**Rollen** hat eine Spalte pro Rolle und eine Zeile pro Recht; ein Häkchen
gibt das Recht sofort.

![Rollen](../images/screenshots/roles.png){ .screenshot }

| Recht | Erlaubt |
|---|---|
| **Anlässe und Ranglisten sehen** | Anlässe, Programme, Ranglisten und Wertungen |
| **Anlässe durchführen** | Anmeldungen, Resultate, Auszeichnungen, Reservationen, Dokumente und Drucken |
| **Anlässe konfigurieren** | Anlässe anlegen, konfigurieren, kopieren, abschliessen und löschen, mit Programmen, Wertungen, Kiosk- und Helferlinks und Online-Buchung |
| **Statistik und Berichte** | Statistik, Berichte mit Beträgen und Ranglisten als CSV |
| **Schützen und Vereine sehen** | die Personendaten der Stammdaten |
| **Schützen und Vereine bearbeiten** | Schützen und Vereine erfassen und ändern |
| **Stammdaten verwalten** | Schützen importieren, Vereine umbenennen, zusammenführen und löschen, Schützen löschen |
| **Stände und Drucker** | Schiessstände, Druck-Agents und Druckereinstellungen |
| **Druckvorlagen** | Druckvorlagen und Bilder (nur solange das Vereinsmodul eingeschaltet ist) |
| **Benutzer und Verein** | Benutzer, Rollen, Einladungen und Vereinsmodule |

Jeder Verein startet mit vier Rollen:

| Rolle | Rechte |
|---|---|
| **Admin** | alle; die Spalte lässt sich nicht ändern und die Rolle nicht löschen |
| **Vorstand** | alle ausser *Stände und Drucker* und *Benutzer und Verein* |
| **Helfer** | Anlässe sehen und durchführen, Statistik, Schützen und Vereine sehen und bearbeiten |
| **Mitglied** | Anlässe und Ranglisten sehen |

**Neue Rolle** fügt eine eigene Rolle hinzu, die nur die Anlässe sieht, bis du mehr anhakst. Der Stift benennt eine Rolle um (ein leerer Name behält den Vorgabenamen), der Papierkorb löscht sie, sobald niemand sie mehr hat. Die App zeigt jedem Mitglied nur die Seiten und Knöpfe, die seine Rolle erlaubt: Ein Mitglied ohne *Anlässe durchführen* öffnet einen Anlass bei seinen Ranglisten.

Konten sind global: Ein Konto löscht der Plattformbetreiber im
[Admin-Panel](../admin-guide/admin-panel.md). Ein vergessenes Passwort setzt
man mit **Passwort vergessen?** auf der Login-Seite zurück, wenn das Konto
eine bestätigte E-Mail-Adresse hat, sonst der Betreiber.
