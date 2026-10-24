# Module und Arten von Anlässen

Vereine führen sehr unterschiedliche Anlässe durch: ein Schützenfest mit
Desk, Scheibenzuteilung und Kranzkarten, ein Feldschiessen am offiziellen
Wochenende, ein Jahresprogramm über die Saison, eine
Vereinsmeisterschaft. Standblatt hat für keinen davon eigenen Code. Jeder
Anlass ist dieselbe Art Ding; Anlässe unterscheiden sich nur in ihren
Programmen, ihren Schiessständen und Tagen und den **Modulen**, die sie
einschalten.

## Der Kern und die Module {#core}

Jeder Anlass hat denselben Kern: Name, Tage, Schiessstände und Sprache,
Programme mit ihren Passen und Preisen, Anmeldungen, Resultate, Ranglisten,
Statistik und Abschliessen. Alles darüber hinaus ist ein Modul:
Scheibenzuteilung, Reservationen mit der Online-Reservation, das Standblatt
und der Kontrollstreifen, die Munitionsausgabe, Durchgänge, Auszeichnungen,
Sonderpreise, die Vereinskonkurrenz, Wertungen, Kiosk-Anzeigen und
Helferlinks.

Ein ausgeschaltetes Modul ist weg, nicht nur versteckt: Es verschwindet aus
der Seitenleiste, dem Dashboard, der Programmseite, den Ranglisten und den
Berichten, und auch der Server verweigert seine Funktionen, sodass weder ein
Helferlink noch ein API-Token es nutzen kann. Ein Jahresprogramm ohne
Auszeichnungen hat deshalb keinen Auszeichnungsabschnitt in seinen
Programmen und keine Auszeichnungsberichte; niemand muss überspringen, was
der Anlass nicht braucht.

Manche Module brauchen andere: *Sonderpreise* schalten *Auszeichnungen* ein,
weil ein Sonderpreis eine Auszeichnung ist, und *Reservationen* brauchen die
*Scheibenzuteilung*, weil eine Reservation Scheiben sperrt. Ein Modul bleibt
an, solange ein Modul an ist, das es braucht.

## Arten von Anlässen sind nur ein Ausgangspunkt {#kinds}

**Neuer Anlass** fragt, was für ein Anlass es ist: Schützenfest, Jahresprogramm oder Training, Feldschiessen, Vereinsmeisterschaft oder eine eigene Auswahl. Die Art schaltet nur eine Reihe von Modulen ein (die Liste steht unter [Arten von Anlässen](../reference/modules.md#kinds)); der Anlass merkt sie sich nicht. Danach lassen sich die Module wie bei jedem anderen Anlass ändern: Ein Feldschiessen, das auch Wertungen will, schaltet sie ein, und ein Schützenfest ohne Desk schaltet die Scheibenzuteilung aus.

Deshalb braucht eine neue Art von Anlass keine neue Funktion: Wenn eine
Kombination von Modulen sie abdeckt, lässt sie sich heute durchführen. Ein
Verein, der jedes Jahr einen Anlass in derselben Form durchführt, speichert
ihn stattdessen als [Vorlage](../user-guide/events.md#templates); eine
Vorlage ist ein Anlass mit Modulen, Programmen und Wertungen, keine Art.

## Warum ein Modul mit Daten an bleibt {#in-use}

Ein Modul mit Daten auszuschalten würde Daten verstecken, die noch zählen:
Auszeichnungsstufen, die Schützen erreicht haben, Kiosk-Links, die noch die
Ranglisten in der Halle zeigen, Programme, die mehrere Durchgänge
erlauben. Standblatt lässt ein solches Modul deshalb an, bis seine Daten weg
sind, und die Modulkarte sagt, warum ihr Schalter gesperrt ist:

| Modul | Bleibt an, solange |
|---|---|
| Scheibenzuteilung | eine Anmeldung eine Scheibe hat |
| Reservationen | der Link der Online-Reservation aktiv ist |
| Durchgänge | ein Programm mehr als einen Durchgang erlaubt |
| Auszeichnungen | ein Programm Auszeichnungsstufen hat |
| Sonderpreise | ein Programm Sonderpreise vergibt |
| Wertungen | der Anlass eine Wertung hat |
| Kiosk-Anzeigen | der Anlass einen Kiosk-Link hat |
| Helferlinks | der Anlass einen Helferlink hat |

Das Standblatt, der Kontrollstreifen, die Munitionsausgabe und die
Vereinskonkurrenz haben keine eigenen Daten und lassen sich jederzeit
umschalten.

## Konfiguration und Betrieb {#configuration}

Ein Anlass wird einmal vorbereitet und oft durchgeführt: Programme, Preise
und Module werden eingerichtet, bevor der erste Schütze kommt, während der
Desk am Tag selbst arbeitet. Standblatt hält beides auseinander. Der
**Betrieb** (Dashboard, Resultate, Ranglisten, Statistik und die Seiten von
Modulen wie Reservationen) braucht das Recht *Anlässe durchführen*; die
**Konfiguration** (allgemeine Einstellungen, Module, Programme, die
Einrichtungsseiten der Module, Kopieren und Abschliessen) braucht *Anlässe
konfigurieren*. Ein Helfer am Desk kann nicht versehentlich einen Preis
ändern, und die Konfiguration steht dem Desk nicht im Weg.

## Vereinsmodule und Anlassmodule {#club-modules}

Die meisten Module gehören zu einem Anlass, weil jeder Anlass entscheidet,
ob er Auszeichnungen abgibt oder Kiosk-Bildschirme nutzt. Manche Bausteine
nutzt jeder Anlass eines Vereins, deshalb gehören sie dem Verein: Die
**Druckvorlagen** mit ihrer Bildbibliothek werden einmal gezeichnet und von
den Programmen jedes Anlasses genutzt. Solche **Vereinsmodule** werden unter
[Vereinsmodule](../user-guide/events.md#club-modules) geschaltet. Ein neuer
Verein hat die Druckvorlagen eingeschaltet; wie bei Anlassmodulen bleibt das
Modul an, solange der Verein Vorlagen oder Bilder hat.

## Anlässe mit und ohne Scheiben {#lanes}

Scheiben spielen nur dort eine Rolle, wo ein Desk sie zuteilt. Standblatt
trennt deshalb den Schiessstand von den Scheiben:

| Der Anlass hat | Was daraus folgt |
|---|---|
| einen Schiessstand und *Scheibenzuteilung* | den Desk eines Schützenfests: Scheibenraster, Reservationen, Online-Reservation, Warner-Links und Scheibenstatus auf Kiosk-Bildschirmen |
| einen Schiessstand, aber keine *Scheibenzuteilung* | Schützen werden ohne Scheibe angemeldet; Dokumente gehen weiterhin an den Druck-Agent des Schiessstands, und Programme lassen sich weiterhin an den Schiessstand binden |
| keinen Schiessstand | keine Scheiben und nichts, was sie braucht; Dokumente öffnen sich als PDF, weil Druck-Agents zu einem Schiessstand gehören |

Ein Jahresprogramm oder ein Feldschiessen, dessen Resultate auf Papier oder
aus der Trefferanzeige kommen, braucht keine Scheiben, aber der Verein will
vielleicht trotzdem seinen Drucker am Schiessstand; das ist die mittlere
Zeile. Module, die Scheiben brauchen, schalten sich von selbst aus, wenn der
Anlass keinen Schiessstand hat.

Resultate kommen in einen Anlass ohne Scheiben auf demselben Weg wie in
jeden anderen: über **Resultate erfassen**, eine CSV-Datei oder die
Resultate einer Anmeldung (siehe [Resultate
erfassen](../user-guide/results-entry.md)).

## Kopieren behält die Form {#copy}

**Für nächstes Jahr kopieren** und Vorlagen übernehmen die eingeschalteten Module und jene, die Daten haben, zusammen mit den Programmen und Wertungen. Ein vom letzten Jahr kopierter Anlass startet deshalb mit derselben Form, ohne seine Anmeldungen und Resultate.
