# Ein Feldschiessen über einen Zeitraum einrichten

Das brauchst du, wenn das Feldschiessen über Wochen geschossen wird,
z. B. von der Eröffnung des Schiessstands bis zum 31. August plus das
offizielle Wochenende, und die Resultate vom Desk, von Warnern an der
Scheibe oder auf Papier kommen. Es braucht das Recht *Anlässe konfigurieren*
(siehe [Rollen](../user-guide/users.md#roles)).

## Den Anlass erstellen {#event}

1. Klicke unter **Anlässe** auf **Neuer Anlass** und wähle
   **Feldschiessen**. Das schaltet Scheibenzuteilung, Auszeichnungen,
   Vereinskonkurrenz, Kiosk-Anzeigen und Helferlinks ein.
2. Gib den **Name** ein und als **Datum** den Haupttag, z. B. den Samstag
   des offiziellen Wochenendes. Sein Jahr bestimmt die Alterskategorien.
3. Kreuze unter **Schiessstände** die Schiessstände an, z. B. *Gewehr 300 m*
   und *Pistole 25 m*.
4. Füge unter **Schiesstage** die Tage des offiziellen Wochenendes mit ihren
   Schiesszeiten hinzu.
5. Schalte **Zeitraum (Saison)** ein, gib **Erster Tag** und **Letzter Tag**
   ein und kreuze unter **Schiesszeiten** die Wochentage an, an denen der
   Schiessstand offen ist, mit ihren Zeiten. Die einzelnen Tage aus Schritt
   4 ergänzen den Zeitraum oder ersetzen an diesem Datum seine Zeiten.
6. Klicke auf **Erstellen**.

## Die Programme hinzufügen {#programmes}

1. Öffne **Konfiguration** → **Programme** und füge die Programme hinzu, im
   Editor oder mit **Übernehmen aus …** (Bibliothek, Programmdatei oder
   Anlass des letzten Jahres). Ein Programm pro Distanz; Waffen mit eigenen
   Limiten für Auszeichnungen, z. B. *Stgw 90* und *Stgw 57*, sind Varianten
   mit gleichem Namen.
2. Kreuze unter **Preise (CHF)** jede Alterskategorie an, die schiessen
   darf. Für ein Feldschiessen ohne Kosten gib 0 ein.
3. Klicke unter **Auszeichnungen** für jede Auszeichnung auf **Auszeichnung
   hinzufügen**, mit ihrer **Art** und dem Mindesttotal pro Alterskategorie.
4. Klicke auf **Erstellen**. Bei einem Anlass mit mehreren Schiessständen
   wähle dann unter **Einsatz** den **Schiessstand** des Programms; er wird
   sofort gespeichert.

## Wählen, wie die Resultate hereinkommen {#results}

Die Wege lassen sich innerhalb eines Anlasses mischen.

- **Am Desk:** Behalte das Modul *Scheibenzuteilung*. Melde jeden Schützen
  im [Dashboard](../user-guide/dashboard.md) auf einer Scheibe an und
  erfasse die Passen unter **Resultate**, wenn die Scheiben zurückkommen
  (siehe [Resultate](../user-guide/results.md)). Schalte unter **Module**
  *Standblatt* und *Kontrollstreifen* ein, wenn sie bei der Anmeldung
  gedruckt werden sollen.
- **An der Scheibe:** Erstelle in der Konfiguration unter **Helferlinks**
  für jede Gruppe von Scheiben einen Warner-Link (siehe
  [Helfer](../user-guide/helpers.md#create)); die Warner erfassen die Passen
  auf einem Handy oder Tablet.
- **Auf Papier:** Schalte unter **Module** die *Scheibenzuteilung* aus, wenn
  niemand auf einer Scheibe angemeldet wird. Tippe die Blätter unter
  **Resultate erfassen** ein (siehe [Resultate
  erfassen](../user-guide/results-entry.md#quick)) oder importiere die Datei
  der Trefferanzeige (siehe [Resultate aus der Trefferanzeige
  importieren](import-results.md)).

Schalte unter **Konfiguration** → **Module** aus, was das Feldschiessen
nicht braucht, z. B. *Kiosk-Anzeigen* oder *Vereinskonkurrenz*.
