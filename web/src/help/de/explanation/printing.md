# Wie das Drucken funktioniert

An einem Schützenfest wird viel gedruckt: ein Standblatt und ein
Kontrollstreifen für jede Anmeldung, Resultatblätter, Kranzkarten,
Vorlagen. Standblatt ist eine Web-App, und das prägt, wie das Drucken
funktioniert.

## Jedes Dokument ist ein PDF {#pdf}

Der Server erzeugt jedes Dokument als PDF, in der Sprache des Anlasses:
Standblätter, Kontrollstreifen, Resultatblätter, Kranzkarten, Druckvorlagen,
Urkunden und Ranglisten. Dasselbe PDF wird am Schiessstand gedruckt, im
Browser geöffnet oder vom Handy geteilt, so sieht ein Dokument überall
gleich aus. Ranglisten werden nur heruntergeladen, als PDF oder CSV, und nie
an einen Drucker des Schiessstands geschickt.

## Warum ein Druck-Agent {#agent}

Ein Browser kann nicht von selbst drucken: Jeder Druck verlangt einen
Dialog, einen Drucker und einen Klick. An einem vollen Desk ist das zu
langsam, und ein Standblatt soll aus dem Drucker neben dem Desk kommen,
während der Kontrollstreifen aus dem Etikettendrucker kommt.

Der **Druck-Agent** löst das. Er ist ein kleines Programm auf einem Computer
neben den Druckern: die App Standblatt Print Agent unter Windows oder auf
einem Mac oder ein Agent ohne Oberfläche auf einem Linux-Computer wie einem
Raspberry Pi. Drückt jemand einen Druckknopf, stellt der Server das PDF in
eine Warteschlange; der Agent holt es und druckt es ohne Nachfrage auf dem
gewählten Drucker. Der Agent verbindet sich mit Standblatt, nicht umgekehrt,
so muss am Schiessstand nichts aus dem Internet erreichbar sein.

Jeden Druckauftrag bestätigt der Agent. Ein Auftrag, der nicht innert zwei
Minuten bestätigt ist, wird erneut ausgegeben, bis zu drei Versuche;
**Letzte Aufträge** zeigt, was passiert ist.

## Agents und Drucker gehören zu einem Schiessstand {#range}

Drucker stehen an einem Schiessstand, nicht in einem Anlass, deshalb gehören
der Druck-Agent und die Druckereinstellungen zum
[Schiessstand](../getting-started/concepts.md#ranges-and-lanes). Jeder
Anlass auf diesem Schiessstand nutzt sie ohne eigene Einrichtung. Die
Dokumente einer Anmeldung werden auf dem Schiessstand ihrer Scheibe
gedruckt; ein Anlass auf mehreren Schiessständen druckt deshalb dort, wo der
Schütze schiesst. Ein Computer kann mehrere Schiessstände bedienen,
z. B. *Pistole 25 m* und *Pistole 50 m* in derselben Halle.

Die Einstellungen wählen einen Drucker **pro Dokument**, weil die Dokumente
verschiedene Drucker brauchen: Der Kontrollstreifen geht an einen
Etikettendrucker, das Standblatt an ein A5-Fach, Kranzkarten und
Druckvorlagen an einen manuellen Einzug mit den vorgedruckten Blättern. Die
Druckeroptionen, die der Agent meldet (Fach, Papierformat, beidseitig usw.),
werden mit der Wahl gespeichert.

Ein Agent wird entweder mit einem Token verbunden, das in seine
Einstellungen kopiert wird, oder, mit der App, über einen Code, der auf dem
Computer erscheint und in Standblatt bestätigt wird (siehe
[Drucken](../user-guide/printing.md#connect)). So oder so speichert der
Server nur einen Fingerabdruck des Tokens, nie das Token selbst.

## Wenn kein Drucker da ist {#fallback}

Ein Dokument geht nur dann an den Drucker, wenn sein Schiessstand für dieses
Dokument einen Drucker gewählt hat und ein Agent des Schiessstands in der
letzten Minute online war. Sonst öffnet sich stattdessen das PDF: auf einem
Computer in einem neuen Tab, auf einem Handy oder Tablet im Teilen-Menü
(drucken, in Dateien sichern, senden), so bleibt die App offen. Ist der
Agent aus, geht nichts verloren; der Desk druckt einfach von Hand.

Dasselbe gilt für Anlässe ohne Schiessstand: Es gibt keinen Agent, an den
gesendet werden könnte, also öffnet sich jedes Dokument als PDF.

## Druckvorlagen auf vorgedruckten Blättern {#forms}

Manche Dokumente gestaltet nicht Standblatt, sie kommen auf Papier:
Kranzkarten, die Standblätter, welche die Armee für das Obligatorische
Programm liefert, die eigenen Urkunden eines Vereins. Eine **Druckvorlage**
beschreibt ein solches Blatt: wohin auf dem Papier jeder Wert kommt, in
Millimetern ab der oberen linken Ecke. Beim Drucken einer Vorlage kommen nur
die Werte auf das Blatt, das schon im Drucker liegt.

Eine Vorlage hat eine **Datenquelle**: die Resultate eines Programms, eine
Anmeldung oder einen Rang. Die Quelle bestimmt, welche Werte die Vorlage
drucken kann und wo sie verwendet wird: Vorlagen von Resultaten werden nach
Abschluss eines Programms oder als Kranzkarten gedruckt, Vorlagen einer
Anmeldung als Kontrollstreifen oder Teilnehmerkarten, Vorlagen eines Rangs
als Urkunden aus den Ranglisten. Die Kranzkarte und der Kontrollstreifen
sind selbst eingebaute Vorlagen, die ein Verein kopieren und anpassen kann.

Ein Scan des leeren Blatts kann im Editor hinter der Vorlage liegen, damit
sich die Felder auf seine Kästchen setzen lassen. Der Scan dient nur zum
Platzieren und wird nie gedruckt.

## Kalibrieren im Grundsatz {#calibration}

Drucker drucken selten genau dort, wo sie sollen: Papiereinzüge und Fächer
verschieben den Druck um einen oder zwei Millimeter. Eine Vorlage, die am
Bildschirm stimmt, kann deshalb auf Papier neben den Kästchen landen. Der
**Druckerversatz** der Vorlage verschiebt jedes Element um denselben Betrag,
um das auszugleichen: Drucke ein Testblatt, miss, wie weit die Werte daneben
liegen, und trage die Differenz ein. Der Versatz ist Teil der Vorlage,
deshalb braucht eine Vorlage, die auf einem anderen, anders verschiebenden
Drucker gedruckt wird, eine Kopie mit eigenem Versatz. Die Vorschau
ignoriert den Versatz und zeigt die Vorlage wie entworfen. Die Schritte
stehen in [Resultate auf die Standblätter der Armee
drucken](../how-to/army-sheet.md#calibrate).
