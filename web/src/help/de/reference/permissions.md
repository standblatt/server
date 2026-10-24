# Rechte und Rollen

Die Rolle eines Mitglieds in einem Verein ist eine Sammlung von
Rechten. Rollen werden unter **Benutzer** → **Rollen** bearbeitet (siehe
[Rollen](../user-guide/users.md#roles)); API-Tokens tragen eine Auswahl
derselben Rechte.

## Rechte {#permissions}

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

- Eines der Rechte, die eine Aktion nennt, genügt.
- Ohne *Anlässe durchführen* öffnet sich ein Anlass auf seinen Ranglisten.
- Die App zeigt jedem Mitglied nur die Seiten und Knöpfe, die seine Rolle
  erlaubt.

## Standardrollen {#default-roles}

| Rolle | Rechte |
|---|---|
| **Admin** | alle; die Spalte lässt sich nicht ändern und die Rolle nicht löschen |
| **Vorstand** | alle ausser *Stände und Drucker* und *Benutzer und Verein* |
| **Helfer** | Anlässe sehen und durchführen, Statistik, Schützen und Vereine sehen und bearbeiten |
| **Mitglied** | Anlässe und Ranglisten sehen |

- **Neue Rolle** startet nur mit *Anlässe und Ranglisten sehen*.
- Eine Rolle lässt sich löschen, sobald niemand sie hat; die Rolle Admin
  nie.
- Ein Verein behält immer mindestens ein Mitglied mit der Rolle Admin.
- Mitglieder mit *Benutzer und Verein*, aber ohne die Rolle Admin vergeben,
  ändern und entfernen nur Rollen, deren Rechte sie selbst haben.

## Helferlinks {#helper-links}

Helferlinks sind keine Rollen des Vereins; jeder Link hat eine dieser Rollen
(siehe [Einen Link erstellen](../user-guide/helpers.md#create)):

| Rolle | Sieht | Darf |
|---|---|---|
| **Warner** | das Scheibenraster der eigenen Scheiben, die Schützen darauf | Resultate und Auszeichnungen erfassen, einen Schützen **Beenden** |
| **Munitionsausgabe** | das Suchfeld, die Warteliste, die Schützen auf allen Scheiben | die Munition als ausgegeben markieren (**Munition**), einem wartenden Schützen eine Scheibe zuteilen, das Standblatt drucken |
| **Resultaterfassung** | das Suchfeld, die Schützen auf allen Scheiben, die Liste der Beendeten | Resultate und Auszeichnungen auf jeder Scheibe erfassen |

Warner und Resultaterfassung dürfen auch
[Druckvorlagen](../user-guide/print-forms.md#use) ihrer Scheiben drucken;
die Resultaterfassung darf Auszeichnungen als abgegeben markieren.

## Berechtigungen von API-Tokens {#scopes}

Ein Token (**Konto** → **API-Tokens**, siehe
[API-Tokens](../user-guide/login-language.md#api-tokens)) hat eine Liste von
Vereinen und von Berechtigungen:

| Berechtigung | Schlüssel | Erlaubt |
|---|---|---|
| *Mein Trainingslog lesen* | `training:read` | das eigene Trainingslog lesen |
| *Meine Trainings erfassen und ändern* | `training:write` | die eigenen Einheiten, Disziplinen, Waffen und Ziele hinzufügen, ändern und löschen |
| ein Recht von oben | `events.view`, `events.operate`, `events.configure`, `reports.view`, `shooters.view`, `shooters.edit`, `shooters.manage`, `ranges.manage`, `forms.manage`, `club.manage` | dasselbe wie in einer Rolle, in den Vereinen des Tokens |

- Ein Vereinsrecht lässt sich nur wählen, wenn die Rolle des Besitzers es in
  einem der Vereine des Tokens hat.
- Ein Token hat nie mehr als die Rolle des Besitzers: Ein Recht, das die
  Rolle verliert, oder ein Verein, den der Besitzer verlässt, ist für das
  Token sofort weg.
- Ein Token ohne Vereine erreicht nur das Trainingslog.

### Vorlagen {#presets}

| Vorlage | Berechtigungen |
|---|---|
| **Ranglisten für eine Website** | *Anlässe und Ranglisten sehen* |
| **Training-App** | *Mein Trainingslog lesen*, *Meine Trainings erfassen und ändern* |
| **Nur lesen** | *Anlässe und Ranglisten sehen*, *Statistik und Berichte*, *Schützen und Vereine sehen*, *Mein Trainingslog lesen* |

Eine Vorlage füllt nur die Rechte aus, die die gewählten Vereine erlauben.

### Ein Token verwenden {#token-use}

| Was | Wert |
|---|---|
| Format | `sb_` und 40 Hexadezimalziffern, einmal angezeigt |
| Header | `Authorization: Bearer sb_…` oder `X-API-Key: sb_…` |
| Ablauf | in 30 oder 90 Tagen, in einem Jahr oder nie |
| Endpunkte | `/api/tenants/…` (Vereinsdaten, nach Recht), `/api/my/…` (Trainingslog, nach Berechtigung), `GET /api/token`, `/api/version`, `/api/auth/options`; mit *Anlässe konfigurieren* auch `/api/modules` und `/api/programme-templates`; alles andere antwortet 403 |
| `GET /api/token` | Name, Berechtigungen und Vereine des Tokens |
| Ratenlimit | standardmässig 600 Anfragen pro Token in fünf Minuten, dann 429; der Betreiber ändert es mit `STANDBLATT_API_RATE_LIMIT` (0 = kein Limit) |
| Unbekannte Tokens | nach 20 unbekannten oder widerrufenen Tokens von einer Adresse in 15 Minuten 429 |
| API-Beschreibung | `/api/docs` |
| Nie erreichbar | Passwort, E-Mail, Kontoeinstellungen, andere Tokens |
