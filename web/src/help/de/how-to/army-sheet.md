# Resultate auf die Standblätter der Armee drucken

Das brauchst du, wenn die Resultate der Schützen auf vorgedruckten Blättern
landen müssen, z. B. auf den Standblättern, die die Armee für das
Obligatorische Programm liefert, und du sie lieber druckst, statt sie von
Hand auszufüllen. Es braucht die Rechte *Druckvorlagen*, *Anlässe
konfigurieren* und *Stände und Drucker* (siehe
[Rollen](../user-guide/users.md#roles)).

Scanne vor dem Start ein leeres Blatt als PNG oder JPEG und halte ein paar
Ersatzblätter für Testdrucke bereit.

## 1. Das Programm vorbereiten {#programme}

1. Prüfe, dass unter **Vereinsmodule** im Vereinsteil der Seitenleiste
   *Druckvorlagen* eingeschaltet ist (es ist an, ausser jemand hat es
   ausgeschaltet).
2. Öffne unter **Konfiguration** → **Programme** des Anlasses das Programm
   und kreuze **Einzelschüsse** für jede Passe an, für deren Schüsse das
   Blatt Kästchen hat. Mach das vor der ersten Anmeldung: Sobald Schützen
   angemeldet sind, sind die Passen festgelegt.
3. Speichere das Programm mit **Speichern**.

## 2. Die Vorlage erstellen {#form}

1. Öffne **Druckvorlagen** im Vereinsteil der Seitenleiste und klicke auf
   **Neue Druckvorlage**.
2. Gib einen Namen ein (z. B. *Obligatorisches Programm 300 m*), wähle bei
   **Druckt** *Resultate eines Programms* und das Papierformat, das dem
   Blatt am nächsten kommt, und klicke auf **Erstellen**. Der Editor öffnet
   sich.
3. Setze rechts unter **Vorlage** **Breite (mm)** und **Höhe (mm)** auf die
   Grösse des Blatts und klicke dann unter **Scan des leeren Formulars** auf
   **Scan hochladen**. Der Scan erscheint beim Bearbeiten hinter dem Blatt;
   gedruckt wird er nie.
4. Ziehe aus **Felder** links die Felder des Schützen (z. B. **Nachname**,
   **Vorname**, **Jahrgang**, **Lizenznummer**) auf ihre Kästchen auf dem
   Scan.
5. Ziehe aus der Gruppe **Einzelschüsse** **Schüsse Passe 1** auf das erste
   Schusskästchen dieser Passe. Es ist eine Reihe von Kästchen im Abstand
   von 8 mm; ziehe den rosa Punkt auf ihrem letzten Kästchen, bis jeder Wert
   in seinem Kästchen sitzt, oder setze **dx (mm)** unter
   **Eigenschaften**. Wiederhole das für die anderen Passen.
6. Füge aus **Resultat** die Passentotale (**Passe 1** …) und das **Total**
   hinzu, wo das Blatt Kästchen dafür hat.
7. **Speichern**.

![Druckvorlagen](../images/screenshots/print-forms.png){ .screenshot }

Wie man Elemente verschiebt, ausrichtet und einrastet, steht in
[Druckvorlagen](../user-guide/print-forms.md#arrange), jedes Feld in
[Platzhalter in Vorlagen](../reference/form-placeholders.md).

## 3. Den Drucker kalibrieren {#calibrate}

1. Lege am Schiessstand ein leeres Blatt in den Drucker.
2. Klicke auf **Testdruck**: Die gespeicherte Vorlage öffnet sich als PDF
   mit erfundenen Werten und ohne den Scan. Drucke sie in Originalgrösse
   (nicht an die Seite angepasst).
3. Miss, wie weit die Werte neben ihren Kästchen liegen, und trage es unter
   **Vorlage** als **Druckerversatz nach rechts (mm)** und **Druckerversatz
   nach unten (mm)** ein (negativ, um nach links oder oben zu verschieben).
4. **Speichern** und ab Schritt 1 wiederholen, bis die Werte in den Kästchen
   sitzen.

## 4. Die Vorlage und den Drucker zuweisen {#assign}

1. Öffne unter **Konfiguration** → **Programme** des Anlasses das Programm
   und wähle die Vorlage im Abschnitt **Druckvorlage**. Sie wird sofort
   gespeichert; das geht auch, wenn Schützen angemeldet sind.
2. Wähle unter **Schiessstände → Drucker einrichten** des Schiessstands den
   Drucker für **Druckvorlagen (vorgedruckt)**, z. B. mit seinem manuellen
   Einzug als Fach, und klicke auf **Speichern** (siehe [Einen Drucker am
   Schiessstand verbinden](connect-printer.md)).

## 5. Die Blätter drucken {#print}

Sobald alle Passen des Programms erfasst sind, zeigen seine Resultate **Auf
Vorlage drucken**: in den Resultaten am Desk und auf Warner-Links an der
Scheibe (siehe [Helfer](../user-guide/helpers.md)). Unter **Resultate
erfassen** hat jede Zeile unter *Soeben erfasst* einen eigenen
Druckknopf. Lege ein Blatt in den Drucker und klicke darauf. Ohne
Druck-Agent online öffnet sich das PDF und wird aus dem Browser gedruckt.
