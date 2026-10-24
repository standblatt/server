# Begriffe

Diese Seite erklärt die Begriffe, die Standblatt verwendet, und wie sie
zusammenhängen. Die deutschen, französischen und italienischen Begriffe
stehen im [Glossar](../reference/glossary.md); die Schritte für jede Aufgabe
stehen im Benutzerhandbuch.

Kurz: Ein **Verein** hat **Schiessstände**, **Schützen** und
**Benutzer**. Er führt **Anlässe** durch; ein Anlass bietet **Programme** an
und schaltet die **Module** ein, die er braucht. Eine **Anmeldung** bringt
einen Schützen mit einem oder mehreren Programmen an den Anlass, jedes
Programm wird ein **Stich**, der Stich erhält **Resultate**, und die
Resultate ergeben **Auszeichnungen** und einen Platz in den
**Ranglisten**. Ist der Anlass **abgeschlossen**, sind die Ranglisten
endgültig und gelangen in die persönlichen Trainingslogs der Schützen.

## Vereine, Benutzer und Rollen {#clubs-users-and-roles}

Ein **Verein** ist ein Mandant: Alles in Standblatt gehört zu genau einem
Verein, und Vereine sehen die Daten der anderen nie. Nur der
Plattformbetreiber legt Vereine an, im Admin-Panel. Konten entstehen aus
einer Einladung eines Vereins, aus der Registrierung bei [Mein
Training](../my-training/account.md), über **Benutzer hinzufügen** auf der
Benutzerseite eines Vereins oder durch den Betreiber.

**Benutzer** sind persönliche Konten. Ein Konto kann zu mehreren Vereinen gehören, mit einer Rolle pro Verein. Eine Rolle ist eine Sammlung von Rechten wie *Anlässe durchführen* oder *Anlässe konfigurieren*; jeder Verein startet mit den Rollen *Admin*, *Vorstand*, *Helfer* und *Mitglied* und kann sie ändern oder eigene hinzufügen (siehe [Rollen](../user-guide/users.md#roles)). Vereine laden Personen per E-Mail ein; die eingeladene Person tritt bei, indem sie den Link annimmt.

Warum Konten global sind, während alle Vereinsdaten im Verein bleiben,
erklärt [Vereine, Konten und Daten](../explanation/data.md).

## Schiessstände und Scheiben {#ranges-and-lanes}

Ein **Schiessstand** (englisch *range*) ist ein Stand des Vereins mit einer
Anzahl **Scheiben** (*lanes*). Scheiben lassen sich als ausser Betrieb
markieren. Drucker und der Druck-Agent gehören zu einem Schiessstand.

Jede Distanz ist ein eigener Schiessstand: Eine Halle mit 25-m- und
50-m-Scheiben sind zwei Schiessstände, jeder mit Scheiben ab 1. Ein Anlass
kann auf mehreren Schiessständen geschossen werden, auf einem oder auf gar
keinem (siehe [Anlässe mit und ohne
Scheiben](../explanation/modules.md#lanes)).

## Anlässe und Programme {#events-and-programmes}

Ein **Anlass** (englisch *event*) ist ein Wettkampf auf einem oder mehreren
Schiessständen an einem oder mehreren Tagen (das Datum ist der Haupttag;
weitere Tage decken Vorschiessen und mehrtägige Anlässe ab). Er hat eine
Sprache, die für gedruckte Dokumente und Kiosk-Bildschirme gilt, und
[Module](../user-guide/events.md#modules), die einschalten, was er braucht,
etwa die Scheibenzuteilung, Auszeichnungen oder die Vereinskonkurrenz.

Es gibt keine Anlasstypen: Ein Schützenfest, ein Feldschiessen und ein
Jahresprogramm sind dieselbe Art Anlass mit verschiedenen Modulen (siehe
[Module und Arten von Anlässen](../explanation/modules.md)).

Ein Anlass bietet **Programme** (*Stiche*) an: Jedes Programm hat **Passen**
mit einer Anzahl Schüsse und Höchstpunktzahl, **Preise** pro Alterskategorie
und **Auszeichnungsstufen** (Kranzkarten, Naturalgaben oder Geldpreise) mit
einem Mindesttotal pro Alterskategorie. Programme mit gleichem Namen sind
**Varianten** (zum Beispiel *Sport* und *Aufgelegt*); ein Schütze schiesst
eine Variante. Ein Programm kann mit einem anderen **kombinierbar** sein: Es
übernimmt dann einige Passen aus diesem Programm, statt separat geschossen
zu werden.

Ein Programm kann mehrere **Durchgänge** erlauben, von denen einer
zählt. Eine **Wertung** rangiert Schützen über mehrere Programme oder
Durchgänge. Wie Programm-Ranglisten, Durchgänge, Wertungen, die
Vereinskonkurrenz und Sonderpreise zusammenhängen, erklärt [Wie die
Ranglisten zusammenhängen](../explanation/rankings.md).

## Schützen, Anmeldungen und Stiche {#shooters-registrations-and-entries}

**Schützen** sind die Personen, die schiessen. Jeder Verein führt seine eigenen Schützen-Stammdaten, zusammen mit den **Vereinen**, zu denen die Schützen gehören (der eigene Verein des Schützen, nicht zu verwechseln mit dem Verein, der Standblatt nutzt).

Eine **Anmeldung** bringt einen Schützen an den Anlass, wahlweise auf eine
Scheibe, mit einem oder mehreren Programmen. Jedes Programm einer Anmeldung
ist ein **Stich** (*entry*). Eine Anmeldung ist *angemeldet* (wartend), *am
Schiessen* (auf einer Scheibe) oder *abgeschlossen* (alle Programme haben
Resultate; die Scheibe ist wieder frei).

## Resultate, Auszeichnungen und Ranglisten {#results-awards-and-rankings}

**Resultate** sind die Punkte pro Passe, wahlweise mit Mouchen (Innenzehner). Sobald alle Passen erfasst sind, kann der Schütze eine der **Auszeichnungen** erhalten, die sein Total in seiner Alterskategorie erreicht. **Ranglisten** ordnen alle Resultate eines Programms nach den [Regeln des SSV](../reference/rules.md#ranking).

Eine Auszeichnung hängt nur vom eigenen Total des Schützen ab; ein Rang
hängt von allen anderen ab. Ranglisten sind live, solange der Anlass läuft,
und werden endgültig, wenn er [abgeschlossen](../explanation/closing.md)
ist.

## Alterskategorien {#age-categories}

Die Alterskategorie folgt der Regel des SSV: Jahr des Anlasses minus
Jahrgang. Der Geburtstag spielt keine Rolle, und das Jahr ist das des
Haupttags des Anlasses, so zählt ein Vorschiesstag im Dezember für einen
Anlass im Januar. Die Kategorie wird bei der Anmeldung festgelegt und
bestimmt die Preise, Auszeichnungen und Sonderpreise, die der Schütze
erhalten kann. Ein Schütze ohne Geburtsdatum zählt als Elite. Die Kategorien
und ihre Alter stehen in den
[Wettkampfregeln](../reference/rules.md#age-categories).

## Dokumente und Drucken {#documents-and-printing}

Jedes Dokument ist ein PDF, das der Server in der Sprache des Anlasses
erzeugt. Am Schiessstand schickt ein kleiner Druck-Agent es an den richtigen
Drucker; ohne Agent öffnet sich das PDF im Browser. Siehe [Wie das Drucken
funktioniert](../explanation/printing.md).
