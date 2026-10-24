# Einen Anlass abschliessen

Solange ein Anlass läuft, kann sich alles darin noch ändern: Ein Resultat
wird korrigiert, eine Anmeldung wechselt die Scheibe, ein Programm erhält
einen neuen Preis. Irgendwann muss der Organisator sagen: *Das ist das
Resultat*. Genau das macht **Anlass abschliessen**: Es macht aus einem
laufenden Anlass einen fertigen.

## Was Abschliessen bedeutet {#meaning}

Abschliessen bewirkt zwei Dinge auf einmal:

- **Die Ranglisten werden endgültig und freigegeben.** Bisher waren sie ein
  Arbeitsstand; ab jetzt sind sie das Resultat des Anlasses.
- **Der Anlass wird schreibgeschützt.** Anmeldungen, Stiche, Resultate,
  Auszeichnungen, Programme mit ihren Passen, Preisen und
  Auszeichnungsstufen, Wertungen und die Einstellungen des Anlasses lassen
  sich nicht mehr ändern.

Alles bleibt sichtbar: das Dashboard, die Resultate, die Ranglisten, die
Statistik, die Berichte und die Dokumente. Ein Banner auf jeder Seite des
Anlasses sagt, dass er abgeschlossen ist.

Beides gehört zusammen. Eine Rangliste kann nur endgültig sein, wenn sich
darunter nichts mehr ändern kann, und Standblatt setzt das in der Datenbank
selbst durch, nicht nur durch ausgeblendete Knöpfe: Jede Änderung an den
Daten eines abgeschlossenen Anlasses wird abgelehnt, ob sie aus der App,
über einen Helferlink oder mit einem API-Token kommt.

## Warum abgeschlossene Anlässe schreibgeschützt sind {#read-only}

Ist ein Anlass abgeschlossen, verlassen sich andere Dinge darauf:

- Schützen übernehmen ihre Resultate und Ränge in ihr persönliches
  Trainingslog.
- Der Finanzbericht, die Auszeichnungen und die Vereinskonkurrenz sind
  abgegeben und abgerechnet.
- Der Anlass des nächsten Jahres wurde vielleicht davon kopiert.

Könnten sich Resultate noch ändern, zeigte das Log eines Schützen vielleicht
einen Rang, den die Rangliste nicht mehr hat, und die Berichte stimmten
nicht mehr mit dem überein, was bezahlt und abgegeben wurde. Deshalb muss
auch, was am Desk geschieht, etwa Auszeichnungen als abgegeben markieren,
vor dem Abschliessen erledigt sein.

## Die Verbindung zu Mein Training {#my-training}

Jeder Standblatt-Benutzer hat ein persönliches Trainingslog, [Mein
Training](../my-training/index.md). Wer im Profil eine **Lizenznummer** und
ein **Geburtsdatum** einträgt, wird mit den Schützen aller Vereine
abgeglichen: Wo immer die Stammdaten eines Vereins einen Schützen mit
derselben Lizenznummer (führende Nullen spielen keine Rolle) und demselben
Geburtsdatum haben, gehören die Resultate dieses Schützen dem Benutzer.

Angeboten werden nur Resultate **abgeschlossener** Anlässe, weil nur deren
Rangliste endgültig ist. Jedes kommt mit dem Rang in seinem Programm und der
Anzahl Teilnehmer, z. B. *Rang 12 von 87*. Nichts wird von selbst
hinzugefügt: Der Schütze wählt eine Disziplin und fügt das Resultat hinzu,
und es wird eine Wettkampf-Einheit im Log mit ihren Passen und ihrem Rang.

Der Abgleich funktioniert über Vereine hinweg, ohne dass ein Verein seine
Daten hergibt: Die Stammdaten der Vereine bleiben, wo sie sind, der Abgleich
liest sie nur für den Benutzer, dem sie gehören, und es werden nur die
eigenen Resultate des Schützen gezeigt. Wer keine Lizenznummer oder kein
Geburtsdatum im Profil hat oder wessen Verein ein anderes Geburtsdatum
erfasst hat, erhält keine Wettkampfresultate. Siehe [Profil und
Wettkampfresultate](../my-training/account.md#competitions).

## Wiedereröffnen {#reopen}

**Wiedereröffnen** macht einen abgeschlossenen Anlass wieder bearbeitbar, z. B. wenn nach dem Abschliessen ein falsches Resultat auftaucht. Es ist für Korrekturen gedacht, nicht als zweiter Arbeitsstand:

- Resultate, die Schützen schon in ihr Log übernommen haben, bleiben dort
  so, wie sie übernommen wurden; eine spätere Korrektur erreicht sie nicht.
- Noch nicht übernommene Resultate werden in *Mein Training* nicht mehr
  angeboten, bis der Anlass wieder abgeschlossen ist.

Schliesse den Anlass wieder ab, sobald die Korrektur erledigt ist.

Ein abgeschlossener Anlass lässt sich nicht löschen; eröffne ihn zuerst
wieder. Fürs nächste Jahr kopieren lässt er sich auch im abgeschlossenen
Zustand.
