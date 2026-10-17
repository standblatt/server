# Reservationen

**Reservationen** sperrt Scheiben des Schiessstands des Anlasses für eine Zeit, zum Beispiel für ein Training oder eine Gruppe.

![Reservationen](../images/screenshots/reservations.png){ .screenshot }

- **Scheiben**: eine Liste wie `5` oder `1, 5, 7-10`.
- **Von** und **Bis**: der Zeitraum. Reservationen auf derselben Scheibe
  dürfen sich nicht überschneiden; direkt hintereinander geht.
- **Als reserviert zeigen (Minuten vorher)**: Vorlaufzeit, ab der die
  Scheiben im Dashboard und im Kiosk als reserviert erscheinen.

Reservationen gehören zum Schiessstand und erscheinen darum in jedem Anlass
auf diesem Stand. Reservierte Scheiben lassen sich nach einer Bestätigung
trotzdem einem Schützen zuweisen. Der Server kündigt Beginn und Ende von
Reservationen selbst an, so wechseln die Bildschirme ohne Neuladen.

Anfragen, die externe Vereine über die
[Online-Reservation](online-booking.md) schicken, stehen über den
Reservationen; nach der Bestätigung werden sie zu gewöhnlichen Reservationen
mit dem Namen des Vereins.
