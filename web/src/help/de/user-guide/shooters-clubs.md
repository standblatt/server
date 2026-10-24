# Schützen und Vereine

Jeder Verein führt seine eigenen Schützen-Stammdaten; nichts wird mit
anderen Vereinen geteilt (siehe [Schützen-Stammdaten pro
Verein](../explanation/data.md#shooters)). **Schützen** und **Vereine**
stehen in der Seitenleiste unter **Stammdaten**.

## Schützen

![Schützen](../images/screenshots/shooters.png){ .screenshot }

- **Suchen** funktioniert wie im Dashboard: Lizenznummer (führende Nullen
  optional), Teil eines Namens oder Vor- und Nachname in beliebiger
  Reihenfolge. Die Liste zeigt 50 Schützen pro Seite, die Vereinsliste 50
  Vereine; die Pfeile darunter blättern weiter.
- Ein Klick auf einen Schützen öffnet die Details: Name, Lizenznummer,
  Geburtsdatum, Adresse und die Vereine, zu denen der Schütze gehört.
- **Verein** (neuer Schütze oder Verein in den Details hinzufügen) schlägt
  die Vereine der Liste vor. Ein Name, der ohne Rücksicht auf
  Gross-/Kleinschreibung und Leerzeichen zu einem Verein passt, ist dieser
  Verein. Jeder andere Name zeigt *ist noch nicht in der Vereinsliste* mit
  ähnlichen Vereinen unter *Meintest du*; ein neuer Verein entsteht erst mit
  **Neuen Verein erfassen**, so wird ein Tippfehler kein zweiter Verein.
- **CSV importieren** (*Stammdaten verwalten*) liest eine CSV-Datei, siehe
  [CSV-Import](../reference/csv-import.md).
- Schützen mit Anmeldungen lassen sich nicht löschen.

## Vereine

![Vereine](../images/screenshots/clubs.png){ .screenshot }

**Vereine** sind die Vereine, zu denen Schützen gehören, mit ihrer SSV-Nummer. Sie erscheinen in Anmeldungen, Ranglisten und der Vereinskonkurrenz. Ein Schütze kann zu mehreren Vereinen gehören und meldet sich für einen davon an. Vereine, die von Anmeldungen verwendet werden, lassen sich nicht löschen, und zwei Vereine können nicht denselben Namen tragen, auch nicht mit anderer Gross-/Kleinschreibung oder anderen Leerzeichen. Warum das nicht der Verein ist, der Standblatt nutzt: [Der Verein des Schützen und der Mandantenverein](../explanation/data.md#shooter-club).

**Zusammenführen mit…** (*Stammdaten verwalten*) räumt ein Duplikat auf, z. B. *S.G. Hiltown* neben *SG Hilltown*: Wähle den Verein, der bleibt (Vereine mit ähnlichem Namen stehen zuoberst). Er behält seinen Namen und erhält die Schützen, Anmeldungen und Reservationsanfragen des anderen, auch die aus abgeschlossenen Anlässen, sodass deren Ranglisten den bleibenden Verein zeigen; die SSV-Nummer wird übernommen, wenn er keine hat. Der andere Verein wird gelöscht.
