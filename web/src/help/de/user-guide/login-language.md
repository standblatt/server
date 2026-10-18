# Anmeldung und Sprache

## Anmelden

Melde dich mit deiner E-Mail-Adresse (oder deinem Benutzernamen) und deinem
Passwort an: Das Passwort hast du beim Annehmen der Einladung deines Vereins
gewählt oder von deinem Vereins-Admin oder dem Plattformbetreiber
erhalten. Vereinskonten entstehen nicht durch Registrierung; nur [Mein
Training](../my-training/account.md) hat eine eigene Registrierung. Passwort
vergessen? Hat dein Konto eine bestätigte E-Mail-Adresse, nutze **Passwort
vergessen?** auf der Login-Seite; sonst bitte den Betreiber, im
[Admin-Panel](../admin-guide/admin-panel.md) ein neues Passwort zu
setzen. Angemeldet änderst du es selbst unter [Konto](#account).

Nach dem Login landest du auf der Startseite mit einer Karte für jeden
Verein, dem du angehörst, und einer für [Mein
Training](../my-training/index.md), das jeder Benutzer hat; **Konto** und
**Abmelden** stehen unter den Karten. Ohne Vereinsmitgliedschaft führt die
Startseite direkt zu Mein Training.

![Login](../images/screenshots/login.png){ .screenshot }

![Startseite](../images/screenshots/home.png){ .screenshot }

Eine Anmeldung bleibt sieben Tage ohne Aktivität und höchstens 30 Tage
gültig; mit angekreuztem **90 Tage angemeldet bleiben** gilt sie auf diesem
Gerät 90 Tage. Nach zu vielen Fehlversuchen werden Anmeldungen von derselben
Adresse eine Weile gesperrt.

Gehört dein Konto zu mehreren Vereinen, wählst du den Verein nach dem
Anmelden; **Verein wechseln** im Benutzermenü führt zu dieser Auswahl
zurück.

## Konto {#account}

**Konto** im Benutzermenü (unten in der Seitenleiste), auf der Startseite und in der Seitenleiste von Mein Training öffnet deine Kontoeinstellungen:

- **Anmeldung**: dein Benutzername (er lässt sich nicht ändern) und deine
  E-Mail-Adresse. **Hinzufügen** oder **Ändern** verlangt dein Passwort und
  schickt einen Bestätigungslink; bis du bestätigst, gilt die Adresse als
  *nicht bestätigt*, und **Erneut senden** schickt einen neuen Link. Mit
  einer bestätigten Adresse kannst du ein vergessenes Passwort selbst
  zurücksetzen.
- **Passwort**: aktuelles und neues Passwort eingeben (mindestens 10
  Zeichen). Andere Geräte werden abgemeldet, das aktuelle bleibt angemeldet.
- **Sprache und Darstellung**: die Oberflächensprache deines Kontos und
  heller oder dunkler Modus.
- **Zweiter Faktor**: siehe unten.
- **API-Tokens**: siehe unten.
- **Konto löschen**: nur für Konten ohne Verein; Trainingslogbuch, Ziele und
  Einstellungen werden endgültig gelöscht. Vereinskonten entfernt der
  Vereins-Admin.

![Konto](../images/screenshots/account.png){ .screenshot }

## API-Tokens {#api-tokens}

Willst du die Ranglisten auf der Vereinswebsite zeigen, Einheiten aus einer
anderen App importieren oder Daten mit einem eigenen Skript auswerten?
Erstelle ein Token unter **Konto → API-Tokens** mit **Neues Token**:

- einen **Namen**, an dem du später erkennst, was es verwendet,
- eine **Vorlage** setzt die Rechte für einen üblichen Zweck: **Ranglisten
  für eine Website**, **Training-App** oder **Nur lesen**,
- die **Vereine**, die das Token nutzen darf (keiner: nur dein
  Trainingslog),
- die **Rechte**: *Mein Trainingslog lesen*, *Meine Trainings erfassen und
  ändern* und für die gewählten Vereine die [Rechte der
  Rollen](users.md#roles), die deine eigene Rolle dort hat,
- wann es **abläuft**: in 30 oder 90 Tagen, in einem Jahr oder nie.

Das Token (`sb_…`) wird nur einmal angezeigt; kopiere es sofort. Skripte
senden es als `Authorization: Bearer sb_…` (oder `X-API-Key: sb_…`). Es
darf, was du gewählt hast, und nie mehr als deine Rolle im Verein: Verlierst
du ein Recht oder verlässt du den Verein, verliert es das Token sofort. Es
erreicht nie dein Passwort, deine E-Mail, dein Konto oder andere
Tokens. `GET /api/token` sagt einem Skript Name, Rechte und Vereine seines
Tokens. Ein Token macht höchstens 600 Anfragen in fünf Minuten. Die Liste
zeigt, wann jedes Token zuletzt benutzt wurde; Widerrufen stoppt es
sofort. Entwickler finden die Endpunkte unter `/api/docs`.

## Zweiter Faktor {#totp}

Eine Authenticator-App (Google Authenticator, Aegis, 1Password und ähnliche)
kann deine Anmeldung schützen: nach dem Passwort verlangt Standblatt den
sechsstelligen Code, den die App anzeigt. Einrichten unter
[Konto](#account): **Einrichten**, QR-Code scannen, aktuellen Code eingeben,
**Einschalten**. **Ausschalten** verlangt dein Passwort. Handy verloren? Der
Betreiber kann den zweiten Faktor im
[Admin-Panel](../admin-guide/admin-panel.md) zurücksetzen. Der Betreiber
kann die Funktion auch für den ganzen Server ausschalten.

## Sprache

Die Oberfläche gibt es auf Englisch, Deutsch, Französisch und
Italienisch. Ohne Wahl folgt Standblatt der Sprache des
Browsers. **Sprache** im Benutzermenü (unten in der Seitenleiste) setzt sie
für dein Konto auf jedem Gerät.

Gedruckte Dokumente, CSV-Exporte und Kiosk-Bildschirme folgen nicht deiner
Wahl, sondern der Sprache des Anlasses, damit jeder Ausdruck eines Anlasses
gleich aussieht (siehe [Anlässe](events.md#settings)).

## Version

Unten in der Seitenleiste steht die Version des Servers, was beim Melden
eines Problems hilft; in Mein Training auf dem Handy steht sie unter
**Einstellungen**. Wird eine neue Version veröffentlicht, bieten offene
Seiten innerhalb einer Minute **Neu laden** an und laden sie spätestens beim
nächsten Seitenwechsel.
