# Anmeldung und Sprache

## Anmelden

Melde dich mit dem Benutzernamen (oder der E-Mail-Adresse) und dem Passwort
an, die du von deinem Vereins-Admin oder dem Plattformbetreiber erhalten
hast. Vereinskonten entstehen nicht durch Registrierung; nur [Mein
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
- **Konto löschen**: nur für Konten ohne Verein; Trainingslogbuch, Ziele und
  Einstellungen werden endgültig gelöscht. Vereinskonten entfernt der
  Vereins-Admin.

![Konto](../images/screenshots/account.png){ .screenshot }

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
