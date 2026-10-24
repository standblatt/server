# Ein Jahresprogramm über die Saison durchführen

Das brauchst du, wenn Mitglieder ein Programm über die Saison so oft
schiessen, wie sie wollen, und eine Vereinsmeisterschaft die besten drei
Resultate jedes Schützen rangiert. Es braucht das Recht *Anlässe
konfigurieren* (siehe [Rollen](../user-guide/users.md#roles)); für das
Erfassen der Resultate braucht es *Anlässe durchführen*.

## Den Anlass erstellen {#event}

1. Klicke unter **Anlässe** auf **Neuer Anlass** und wähle **Jahresprogramm
   oder Training**. Das schaltet die Module *Durchgänge* und *Wertungen*
   ein; wähle stattdessen **Vereinsmeisterschaft**, wenn die Rangliste auch
   auf einer Kiosk-Anzeige laufen soll.
2. Gib den **Name** ein, z. B. *Jahresprogramm 2027*, und als **Datum** den
   ersten Tag der Saison. Sein Jahr bestimmt die Alterskategorien.
3. Kreuze unter **Schiessstände** den Schiessstand an, oder keinen, wenn die
   Resultate nur auf Papier kommen (siehe [Anlässe ohne
   Scheiben](../user-guide/events.md#no-lanes)).
4. Schalte **Zeitraum (Saison)** ein, gib **Erster Tag** und **Letzter Tag**
   ein und kreuze unter **Schiesszeiten** die Wochentage des
   Vereinstrainings mit ihren Zeiten an. Kein angekreuzter Wochentag heisst
   jeder Tag des Zeitraums.
5. Klicke auf **Erstellen**.

## Das Programm hinzufügen {#programme}

1. Öffne **Konfiguration** → **Programme** und klicke auf **Neues Programm**
   (oder auf **Übernehmen aus …**, um es aus der Bibliothek, einer Datei
   oder einem früheren Anlass zu übernehmen).
2. Gib den Namen und die Passen ein.
3. Kreuze unter **Durchgänge** **Unbegrenzt** an und setze **Zählender
   Durchgang** auf **Der beste**.
4. Kreuze unter **Preise (CHF)** jede Alterskategorie an, die schiessen
   darf, und gib den Preis ein. Das zweite Feld ist der Preis jedes weiteren
   Durchgangs (*Nachdoppel*): leer heisst gleicher Preis, 0 heisst, das
   Programm wird einmal bezahlt.
5. Klicke auf **Erstellen**.

## Die besten drei Durchgänge rangieren {#classification}

1. Öffne **Konfiguration** → **Wertungen** und klicke auf **Neue Wertung**.
2. Gib einen **Name** ein, z. B. *Vereinsmeisterschaft*, wähle das
   Jahresprogramm als **Programm 1** und setze **Zählende Durchgänge** auf
   3.
3. Schalte **Nur Schützen mit allen zählenden Resultaten rangieren** ein,
   wenn nur Schützen mit drei Resultaten einen Rang erhalten sollen; die
   anderen stehen dann unter der Rangliste.
4. Füge **Klassen** hinzu, wenn Junioren und Erwachsene getrennt rangiert
   werden, z. B. *Junioren* mit U17 bis U21 und *Elite* mit den anderen.
5. Klicke auf **Erstellen**.

Alle Einstellungen des Editors beschreibt
[Wertungen](../user-guide/classifications.md).

## Die Resultate während der Saison erfassen {#results}

1. Öffne den Anlass und klicke in der Seitenleiste auf **Resultate
   erfassen**.
2. Wähle das **Programm**, tippe eine Lizenznummer oder einen Namen in
   **Schütze** und drücke ++enter++.
3. Tippe die Passen und drücke ++enter++. Jedes Resultat eines Schützen ist
   ein neuer Durchgang.

Resultate aus der Trefferanzeige lassen sich auch als Datei importieren
(siehe [Resultate aus der Trefferanzeige importieren](import-results.md)).

**Ranglisten** zeigt das Programm mit dem besten Durchgang jedes Schützen (*2/5*: der zweite von fünf Durchgängen zählt) und die Vereinsmeisterschaft mit der Summe der besten drei.

## Die Saison abschliessen {#close}

Wenn die letzten Resultate erfasst und geprüft sind, öffne **Konfiguration**
→ **Kopieren & abschliessen** und klicke auf **Anlass abschliessen**. Was
das Abschliessen ändert, erklärt [Einen Anlass
abschliessen](../explanation/closing.md); die nächste Saison beginnt mit
[Den Anlass des letzten Jahres übernehmen](next-year.md).
