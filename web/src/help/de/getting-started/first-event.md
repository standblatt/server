# Dein erster Anlass

Diese Anleitung führt einen Verein vom leeren Konto bis zum laufenden
Wettkampf. Sie setzt voraus, dass der Plattformbetreiber den Verein und dein
Konto als dessen Admin angelegt hat (siehe
[Admin-Panel](../admin-guide/admin-panel.md)).

## 1. Anmelden

Öffne die App, melde dich mit dem erhaltenen Konto an und wähle deine
Sprache im Benutzermenü unten in der Seitenleiste.

Solange der Verein keinen Anlass hat, zeigt die Seite **Anlässe** eine
Checkliste **Erste Schritte** mit den Schritten unten; jeder offene Schritt
hat einen Knopf, der dich hinführt. Ein Anlass ohne aktive Programme zeigt
dieselbe Checkliste statt seines Dashboards.

![Login](../images/screenshots/login.png){ .screenshot }

## 2. Schiessstand einrichten

Lege unter **Schiessstände** den Stand mit seiner Anzahl Scheiben an und
markiere Scheiben, die ausser Betrieb sind. Wenn du am Stand druckst, füge
später einen Druck-Agenten hinzu (siehe
[Drucken](../user-guide/printing.md)).

## 3. Schützen laden

Importiere unter **Schützen** deine Schützenliste als CSV-Datei (siehe
[CSV-Import](../reference/csv-import.md)) oder lege Schützen einzeln
an. Schützen lassen sich auch während des Anlasses am Empfang erfassen.

## 4. Anlass erstellen

Erstelle unter **Anlässe** den Anlass: Wähle, was für ein Anlass es ist (das
schaltet die passenden [Module](../user-guide/events.md#modules) ein), dann
Name, Datum, Schiessstand, Veranstalter und die Sprache von Dokumenten und
Kiosk-Bildschirmen. Hast du den Anlass letztes Jahr schon durchgeführt,
öffne den Anlass vom letzten Jahr und nutze stattdessen **Für nächstes Jahr
kopieren**: Das übernimmt die Einstellungen und alle Programme. Eine
[Vorlage des Vereins](../user-guide/events.md#templates) tut dasselbe für
Anlässe, die in verschiedener Form wiederkehren.

![Die Anlässe eines Vereins](../images/screenshots/events.png){ .screenshot
}

## 5. Programme hinzufügen

Öffne **Konfiguration** in der Seitenleiste des Anlasses. Schalte unter
**Module** ein, was der Anlass verwendet (Kiosk-Anzeigen, Helferlinks,
Auszeichnungen …, siehe [Module](../user-guide/events.md#modules)); füge
dann unter **Programme** die Programme hinzu: mit dem Editor, aus der
Programmbibliothek, aus einer Programmdatei oder als Kopie aus einem
früheren Anlass. Prüfe Preise, Auszeichnungsstufen, Varianten und
Kombinationen, und richte
[Sonderpreise](../user-guide/programmes.md#special-prizes) ein, wenn dein
Anlass welche vergibt.

![Die Programme des Anlasses](../images/screenshots/setup.png){ .screenshot
}

## 6. Kiosk-Bildschirme und Drucker

Erstelle, immer noch in der Konfiguration, auf **Kiosk-Anzeigen** einen
Kiosk-Link für jeden Bildschirm in der Halle (siehe [Ranglisten und
Kiosk](../user-guide/rankings-kiosk.md)). Weise unter Schiessstände →
Drucken jedem Dokument einen Drucker zu (siehe
[Drucken](../user-guide/printing.md)).

## 7. Den Anlass durchführen

Drücke im **Dashboard** ++f4++, suche den Schützen, melde ihn auf einer
Scheibe an und lass das Standblatt drucken. Erfasse die Resultate, wenn der
Schütze die Scheiben abgibt; Ranglisten und Kiosk-Bildschirme aktualisieren
sich sofort. Siehe [Dashboard](../user-guide/dashboard.md),
[Anmeldungen](../user-guide/registrations.md) und
[Resultate](../user-guide/results.md).

## 8. Nach dem Anlass

Unter **Statistik** lädst du die Berichte herunter (Finanzbericht,
Auszeichnungen, Vereinskonkurrenz, Sonderpreise …). Wenn alles erledigt ist,
**schliesst** du den Anlass unter **Konfiguration** → **Kopieren &
abschliessen** **ab**: Die Ranglisten werden für das persönliche
Trainingslog der Schützen freigegeben, und der Anlass bleibt lesbar, kann
aber nicht mehr geändert werden.
