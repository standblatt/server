# Module

Module sind die Bausteine, die ein Anlass oder ein Verein einschaltet. Was
aus ist, verschwindet aus der Seitenleiste, dem Dashboard, der
Programmseite, den Ranglisten und den Berichten. Wie man sie schaltet:
[Module](../user-guide/events.md#modules) und
[Vereinsmodule](../user-guide/events.md#club-modules).

## Anlassmodule {#event-modules}

Unter **Konfiguration** → **Module** eines Anlasses.

| Modul | Was es bringt |
|---|---|
| Scheibenzuteilung | Schützen werden auf Scheiben angemeldet: das Scheibenraster im Dashboard, der Scheibenstatus auf Kiosk-Anzeigen und Warner-Links (braucht einen Schiessstand; siehe unten) |
| Reservationen | für Gruppen und Vereine reservierte Scheiben; seine Seite in der Konfiguration schaltet die [Online-Reservation](../user-guide/online-booking.md) ein, über die Vereine mit einem öffentlichen Link Scheiben anfragen (braucht einen Schiessstand) |
| Standblatt | das bei der Anmeldung gedruckte Standblatt: Schütze, Scheibe, Programme und Barcodes |
| Kontrollstreifen | der bei der Anmeldung gedruckte Kontrollstreifen: die Passen, die der Schütze schiessen darf |
| Munitionsausgabe | der Knopf **Munition** im Dashboard und für die Munitionsausgabe (siehe [Helfer](../user-guide/helpers.md#ammunition)); aus lassen, wo die Schützen ihre eigene mitbringen, z. B. Luftpistole und Luftgewehr |
| Durchgänge | Programme, die mehrmals geschossen werden können (siehe [Programme](../user-guide/programmes.md#attempts)) |
| Auszeichnungen | Kranzkarten, Naturalgaben und Geldpreise nach Punkten pro Alterskategorie |
| Sonderpreise | die besten Ränge jeder Alterskategorie (schaltet Auszeichnungen ein, siehe [Programme](../user-guide/programmes.md#special-prizes)) |
| Vereinskonkurrenz | Vereine nach den Resultaten ihrer Schützen rangiert |
| Wertungen | Ranglisten über mehrere Programme oder Durchgänge (siehe [Wertungen](../user-guide/classifications.md)) |
| Kiosk-Anzeigen | Ranglisten auf Bildschirmen in der Halle (siehe [Ranglisten und Kiosk](../user-guide/rankings-kiosk.md)) |
| Helferlinks | Warner, Munitionsausgabe und Resultaterfassung ohne Konto (siehe [Helfer](../user-guide/helpers.md)) |

### Wo ein Modul erscheint {#places}

| Modul | Seiten und Abschnitte |
|---|---|
| Scheibenzuteilung | Scheibenraster im Dashboard, Scheiben in der Anmeldung, Scheibenstatus auf Kiosk-Anzeigen, Rolle *Warner* der Helferlinks |
| Reservationen | Seite **Reservationen** des Anlasses, mit den Reservationsanfragen und ihrer Anzahl in der Seitenleiste; Seite **Reservationen** in der Konfiguration (Online-Reservation); Bericht *Scheibenreservationen* |
| Standblatt | *Standblatt drucken* in der Anmeldung |
| Kontrollstreifen | *Kontrollstreifen drucken* in der Anmeldung |
| Munitionsausgabe | Knopf **Munition** und Abzeichen *Munition* im Dashboard und auf Helferlinks |
| Durchgänge | Abschnitt **Durchgänge** der Programmseite; Zähler pro Programm in der Anmeldung |
| Auszeichnungen | Abschnitt **Auszeichnungen** der Programmseite; Auszeichnungen in den Resultaten; Kennzahlen unter **Statistik**; Berichte *Auszeichnungen* und *Abzugebende Auszeichnungen* |
| Sonderpreise | Abschnitt **Sonderpreise** der Programmseite; Bericht *Sonderpreise* |
| Vereinskonkurrenz | *Zählt für die Vereinskonkurrenz* auf der Programmseite; Vereinskonkurrenz unter **Statistik**; Bericht *Vereinskonkurrenz* |
| Wertungen | Seite **Wertungen** in der Konfiguration; Tabellen unter **Ranglisten** und auf Kiosk-Anzeigen; Bericht *Wertungen* |
| Kiosk-Anzeigen | Seite **Kiosk-Anzeigen** in der Konfiguration; *Auf Kiosk-Anzeigen zeigen* auf der Programmseite |
| Helferlinks | Seite **Helferlinks** in der Konfiguration |

### Regeln {#rules}

| Modul | Braucht | Braucht einen Schiessstand | Einstellungen | Lässt sich nicht ausschalten, solange |
|---|---|---|---|---|
| Scheibenzuteilung | – | ja | – | eine Anmeldung eine Scheibe hat |
| Reservationen | Scheibenzuteilung | ja | die Online-Reservation auf ihrer Konfigurationsseite | der Link der Online-Reservation aktiv ist |
| Standblatt | – | – | – | – |
| Kontrollstreifen | – | – | *Vorlage der Kontrollstreifen* | – |
| Munitionsausgabe | – | – | – | – |
| Durchgänge | – | – | – | ein Programm mehr als einen Durchgang erlaubt |
| Auszeichnungen | – | – | *Wert der Auszeichnung auf Kranzkarten drucken* (bei einem neuen Anlass an), *Vorlage der Kranzkarten* | ein Programm Auszeichnungsstufen hat |
| Sonderpreise | Auszeichnungen | – | – | ein Programm Sonderpreise vergibt |
| Vereinskonkurrenz | – | – | *Von der Vereinskonkurrenz ausgeschlossener Verein* | – |
| Wertungen | – | – | – | der Anlass eine Wertung hat |
| Kiosk-Anzeigen | – | – | *Scheibenstatus auf den Ranglisten zeigen* (mit Scheibenzuteilung) | der Anlass einen Kiosk-Link hat |
| Helferlinks | – | – | – | der Anlass einen Helferlink hat |

- Ein Modul, das ein anderes braucht, bleibt an, solange das andere an ist.
- Module, die einen Schiessstand brauchen, sind in einem Anlass ohne
  Schiessstand aus, egal was gewählt ist (siehe [Anlässe ohne
  Scheiben](../user-guide/events.md#no-lanes)).
- Die Einstellungen stehen auf der Karte des Moduls unter **Module**;
  *Scheibenstatus auf den Ranglisten zeigen* steht auf der Seite
  **Kiosk-Anzeigen**.
- **Für nächstes Jahr kopieren** und Vorlagen übernehmen die gewählten
  Module und die Module, die in Gebrauch sind.
- Ein über die API ohne Liste von Modulen erstellter Anlass erhält jedes
  Anlassmodul ausser *Munitionsausgabe*.

## Vereinsmodule {#club-modules}

Unter **Vereinsmodule** im Vereinsteil der Seitenleiste (Recht *Benutzer und
Verein*).

| Modul | Was es bringt | Lässt sich nicht ausschalten, solange |
|---|---|---|
| Druckvorlagen | Vorlagen für vorgedruckte Blätter und der Abschnitt *Druckvorlage* jedes Programms (siehe [Druckvorlagen](../user-guide/print-forms.md)); Seite **Druckvorlagen** und das Recht *Druckvorlagen* | der Verein eine Druckvorlage oder ein Bild hat |

Ein neuer Verein startet mit eingeschalteten *Druckvorlagen*. Sind sie aus,
drucken Anlässe die eingebaute Kranzkarte und den eingebauten
Kontrollstreifen.

## Arten von Anlässen {#kinds}

**Neuer Anlass** fragt nach der Art; sie schaltet nur Module ein, alles lässt sich später ändern.

| Art | Module |
|---|---|
| Schützenfest | Scheibenzuteilung, Reservationen, Standblatt, Kontrollstreifen, Auszeichnungen, Sonderpreise, Vereinskonkurrenz, Kiosk-Anzeigen, Helferlinks |
| Jahresprogramm oder Training | Durchgänge, Wertungen |
| Feldschiessen | Scheibenzuteilung, Auszeichnungen, Vereinskonkurrenz, Kiosk-Anzeigen, Helferlinks |
| Vereinsmeisterschaft | Durchgänge, Wertungen, Kiosk-Anzeigen |
| Eigene Auswahl | keine |

Eine [Vorlage](../user-guide/events.md#templates) des Vereins bringt
stattdessen ihre eigenen Module mit.
