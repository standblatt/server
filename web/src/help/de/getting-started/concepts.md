# Begriffe

## Vereine, Benutzer und Rollen {#clubs-users-and-roles}

Ein **Verein** ist ein Mandant: Alles in Standblatt gehört zu genau einem
Verein, und Vereine sehen die Daten der anderen nie. Der Plattformbetreiber
legt Vereine und Benutzerkonten im
[Admin-Panel](../admin-guide/admin-panel.md) an.

**Benutzer** sind persönliche Konten. Ein Konto kann zu mehreren Vereinen gehören, mit einer Rolle pro Verein:

| Role | Can |
|---|---|
| **Admin** | everything: ranges, events, programmes, kiosk links, users of the club, special prizes, closing events |
| **Staff** | run events: register shooters, enter results, print, reservations; read everything else |

## Schiessstände und Scheiben

Ein **Schiessstand** (englisch *range*) ist ein Stand des Vereins mit einer
Anzahl **Scheiben** (*lanes*). Scheiben lassen sich als ausser Betrieb
markieren. Drucker und der Druck-Agent gehören zu einem Schiessstand.

## Anlässe und Programme

Ein **Anlass** (englisch *event*) ist ein Wettkampf auf einem oder mehreren
Schiessständen an einem oder mehreren Tagen (das Datum ist der Haupttag;
weitere Tage decken Vorschiessen und mehrtägige Anlässe ab). Er hat eine
Sprache, die für gedruckte Dokumente und Kiosk-Bildschirme gilt, und einige
Einstellungen wie den Verein, der von der Vereinskonkurrenz ausgeschlossen
ist.

Ein Anlass bietet **Programme** (*Stiche*) an: Jedes Programm hat **Passen**
mit einer Anzahl Schüsse und Höchstpunktzahl, **Preise** pro Alterskategorie
und **Auszeichnungsstufen** (Kranzkarten oder Naturalgaben) mit einem
Mindesttotal pro Alterskategorie. Programme mit gleichem Namen sind
**Varianten** (zum Beispiel *Sport* und *Aufgelegt*); ein Schütze schiesst
eine Variante. Ein Programm kann mit einem anderen **kombinierbar** sein: Es
übernimmt dann einige Passen aus diesem Programm, statt separat geschossen
zu werden.

## Schützen, Anmeldungen und Stiche

**Schützen** sind die Personen, die schiessen. Jeder Verein führt seine eigenen Schützen-Stammdaten, zusammen mit den **Vereinen**, zu denen die Schützen gehören (der eigene Verein des Schützen, nicht zu verwechseln mit dem Verein, der Standblatt nutzt).

Eine **Anmeldung** bringt einen Schützen an den Anlass, wahlweise auf eine
Scheibe, mit einem oder mehreren Programmen. Jedes Programm einer Anmeldung
ist ein **Stich** (*entry*). Eine Anmeldung ist *angemeldet* (wartend), *am
Schiessen* (auf einer Scheibe) oder *abgeschlossen* (alle Programme haben
Resultate; die Scheibe ist wieder frei).

## Resultate, Auszeichnungen und Ranglisten

**Resultate** sind die Punkte pro Passe, wahlweise mit Mouchen (Innenzehner). Sobald alle Passen erfasst sind, kann der Schütze eine der **Auszeichnungen** erhalten, die sein Total in seiner Alterskategorie erreicht. **Ranglisten** ordnen alle Resultate eines Programms nach den [Regeln des SSV](../reference/rules.md#ranking).

## Alterskategorien {#age-categories}

Die Alterskategorie ergibt sich aus dem Geburtsjahr und dem Jahr des
Anlasses:

| Kategorie | Kurz | Alter im Jahr des Anlasses |
|---|---|---|
| U10 | U10 | bis 9 |
| U13 | U13 | 10–12 |
| U15 | U15 | 13–14 |
| U17 | U17 | 15–16 |
| U19 | U19 | 17–18 |
| U21 | U21 | 19–20 |
| Elite | E | 21–45 |
| Senioren | S | 46–59 |
| Veteranen | V | 60–69 |
| Seniorveteranen | SV | 70 und älter |

Ein Schütze ohne Geburtsdatum zählt als Elite.

## Glossar

| Englisch (Standblatt) | Deutsch | Französisch | Italienisch |
|---|---|---|---|
| Range | Stand | Stand | Stand |
| Lane | Scheibe | Cible | Bersaglio |
| Event | Anlass | Manifestation | Manifestazione |
| Programme | Stich / Programm | Programme | Programma |
| Series | Passe | Passe | Serie |
| Award | Auszeichnung | Distinction | Distinzione |
| Wreath card | Kranzkarte | Carte-couronne | Carta corona |
| Prize in kind | Naturalgabe | Prix en nature | Premio in natura |
| Club competition | Vereinskonkurrenz | Concours de sections | Concorso di sezione |
| Lane sheet | Standblatt | Feuille de stand | Foglio di stand |
| Control strip | Kontrollstreifen | Bande de contrôle | Striscia di controllo |
