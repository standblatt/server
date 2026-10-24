# Vereine, Konten und Daten

Standblatt bedient viele Vereine auf einem Server. Diese Seite erklärt, wem
welche Daten gehören, warum Personen ein Konto für alles haben und warum die
Schützen eines Vereins nicht die eines anderen sind.

## Vereine sind Mandanten {#tenants}

Ein **Verein** ist in Standblatt ein Mandant: Seine Schiessstände, Schützen,
Vereinsliste, Anlässe, Programme, Druckvorlagen und Benutzer gehören ihm und
niemandem sonst. Vereine sehen die Daten der anderen nie, nicht einmal, dass
es anderswo einen Schützen oder einen Anlass gibt. Die Trennung ist Teil
jeder Abfrage des Servers, kein Filter der Oberfläche.

## Ein Konto, eine Rolle pro Verein {#accounts}

Personen sind etwas anderes als Daten: Dieselbe Person kann Sekretärin eines
Vereins sein, am Desk des Schützenfests eines anderen Vereins helfen und ein
persönliches Trainingslog führen. Standblatt hat deshalb **globale
Konten**. Ein Konto meldet sich einmal an und erreicht jeden Verein, zu dem
es gehört, mit einer Rolle pro Verein, die bestimmt, was es dort darf (siehe
[Rollen](../user-guide/users.md#roles)).

Ein Verein besitzt nicht das Konto, nur die Mitgliedschaft. Wird ein
Mitglied entfernt, verliert es den Zugang zu diesem Verein; das Konto, sein
Trainingslog und seine anderen Vereine bleiben. Vereine laden Personen per
E-Mail ein, und die Person wird erst Mitglied, wenn sie den Link annimmt, so
wird niemand ohne sein Wissen einem Verein hinzugefügt.

Helfer an einem Anlass brauchen gar kein Konto: Ein
[Helferlink](../user-guide/helpers.md) gibt Zugang zu einem Anlass und einer
Aufgabe. Ein [API-Token](../user-guide/login-language.md#api-tokens) handelt
für seinen Besitzer und nie mit mehr als dessen Rolle in einem Verein.

## Schützen-Stammdaten pro Verein {#shooters}

Ein **Schütze** ist kein Konto. Er ist ein Datensatz in den Stammdaten eines
Vereins: Name, Lizenznummer, Geburtsdatum, Adresse und die Vereine, zu denen
der Schütze gehört. Jeder Verein führt seine eigenen Schützen, importiert
aus einer CSV-Datei, am Desk erfasst oder aus einer Online-Reservation
übernommen.

Es gibt kein gemeinsames Register aller lizenzierten Schützen: Die
Mitgliederdaten des Verbands sind nicht Teil von Standblatt, und da Vereine
die Daten der anderen nie sehen, bleiben die Schützen, die ein Verein
erfasst hat, für die anderen unsichtbar. Der Preis dafür ist, dass ein
Schütze, der an Anlässen von drei Vereinen schiesst, dreimal erfasst ist,
einmal in jedem Verein. Wo die Datensätze zusammenkommen müssen, für das
eigene Trainingslog eines Schützen, erledigen das die Lizenznummer und das
Geburtsdatum (siehe [Einen Anlass abschliessen](closing.md#my-training)).

## Der Verein des Schützen und der Mandantenverein {#shooter-club}

Zwei verschiedene Dinge heissen Verein:

- Der **Mandantenverein** ist der Verein, der Standblatt nutzt und den
  Anlass durchführt.
- Der **Verein des Schützen** ist ein Eintrag in der Liste **Vereine** des
  Mandanten: der Verein, zu dem ein Schütze gehört, z. B. ein Gastverein an
  einem Schützenfest. Diese Vereine haben keinen Zugang zu Standblatt; sie
  sind Namen mit einer SSV-Nummer.

Ein Schütze kann zu mehreren Vereinen gehören und meldet sich für einen
davon an. Dieser Verein wird auf die Dokumente gedruckt und zählt in den
Ranglisten und der Vereinskonkurrenz. Weil Vereine unabhängig von Gross- und
Kleinschreibung und Leerzeichen über den Namen erkannt und neue nur
absichtlich erstellt werden, macht ein Tippfehler aus einem Verein nicht
zwei; wo es trotzdem passiert ist, räumt Zusammenführen auf, sogar in
abgeschlossenen Anlässen.

## Live-Aktualisierung {#live}

Mehrere Personen arbeiten gleichzeitig am selben Anlass: Der Desk meldet an,
Warner erfassen Resultate an den Scheiben, Kiosk-Bildschirme zeigen die
Ranglisten in der Halle. Sobald etwas gespeichert wird, benachrichtigt der
Server jede offene Seite dieses Anlasses, und jede Seite holt, was sich
geändert hat. Ein an einer Scheibe erfasstes Resultat erscheint ohne
Neuladen im Dashboard und auf den Kiosk-Bildschirmen. Nach einem
Netzwerkunterbruch holt eine Seite von selbst auf, so geht nichts verloren.

## Offline {#offline}

Vereinsarbeit braucht eine Verbindung. Eine Scheibe kann nur einmal vergeben
werden, ein Schütze kann ein Programm nur einmal schiessen, und ein auf
einem Gerät gespeichertes Resultat muss auf allen anderen sein; das kann nur
der Server für alle gleichzeitig entscheiden. Ohne Verbindung zeigen
Vereinsseiten und Helferlinks einen roten Balken und speichern nichts, aber
offene Dialoge behalten, was eingetippt wurde, damit es gespeichert werden
kann, sobald die Verbindung zurück ist.

*Mein Training* ist anders, weil niemand sonst an einem persönlichen Log arbeitet. Eine ohne Verbindung gespeicherte Einheit bleibt auf dem Gerät und wird später hochgeladen, ohne sie je doppelt anzulegen, und die installierte App öffnet sich ohne Verbindung mit den letzten Daten. Siehe [Unfertige Einheiten und offline](../my-training/logging.md#unfinished-sessions-and-offline).
