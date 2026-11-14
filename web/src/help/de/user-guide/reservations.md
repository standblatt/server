# Reservationen

**Reservationen** sperrt Scheiben des Schiessstands des Anlasses für eine Zeit, zum Beispiel für ein Training oder eine Gruppe. Es ist die Seite des Moduls *Reservationen*, das einen Schiessstand und die *Scheibenzuteilung* braucht.

![Reservationen](../images/screenshots/reservations.png){ .screenshot }

**Neue Reservation** fragt nach:

- **Beschreibung** und, bei mehreren Schiessständen, dem **Schiessstand**.
- **Scheiben**: eine Liste wie `5` oder `1, 5, 7-10`.
- **Von** und **Bis**: der Zeitraum. Reservationen auf derselben Scheibe
  dürfen sich nicht überschneiden; direkt hintereinander geht.
- **Als reserviert zeigen (Minuten vorher)**: Vorlaufzeit, ab der die
  Scheiben im Dashboard und im Kiosk als reserviert erscheinen.

Die Liste zeigt pro Tag einen Block in zeitlicher Reihenfolge. Sie markiert
jede Reservation als *geplant*, *aktiv* oder *vorbei*, und ein Streifen
kleiner Kästchen zeigt, welche Scheiben des Stands sie belegt (ab
Tablet-Breite); **vergangene zeigen** fügt die vergangenen hinzu. Ihr Menü
bietet **Bearbeiten**, **Kopieren** und **Löschen**.

Reservationen gehören zum Schiessstand und erscheinen darum in jedem Anlass
auf diesem Stand. Reservierte Scheiben lassen sich nach einer Bestätigung
trotzdem einem Schützen zuweisen. Der Server kündigt Beginn und Ende von
Reservationen selbst an, so wechseln die Bildschirme ohne Neuladen.

Offene Anfragen, die externe Vereine über die
[Online-Reservation](online-booking.md) schicken, stehen in derselben Liste,
eine Zeile pro Zeitblock mit der Markierung *angefragt* und den gewünschten
Scheiben in Amber. Nach der Bestätigung werden sie zu gewöhnlichen
Reservationen mit dem Namen des Vereins und der Markierung
*Online-Reservation*; abgelehnte und stornierte Anfragen sind unter der
Liste eingeklappt.
