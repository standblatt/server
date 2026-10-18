# Helfer

Helfer am Anlass (Warner, Munitionsausgabe, Resultaterfassung) arbeiten mit
einem **Helferlink** statt mit einem Login. Ein Helferlink öffnet auf jedem
Handy oder Laptop eine reduzierte Ansicht eines Anlasses; Admins schalten
das Modul *Helferlinks* unter **Konfiguration** → **Module** ein, erstellen
die Links auf **Helferlinks** in der Konfiguration und widerrufen sie, wenn
der Anlass vorbei ist.

![Helferseite des Warners auf einem
Handy](../images/screenshots/helper-scorer.png){ .screenshot .phone }

## Rollen {#roles}

| Rolle | Sieht | Darf |
|---|---|---|
| **Warner** | das Scheibenraster der eigenen Scheiben, die Schützen darauf | Resultate und Auszeichnungen erfassen, einen Schützen **Beenden** |
| **Munitionsausgabe** | das Suchfeld, die Warteliste, die Schützen auf allen Scheiben | die Munition als ausgegeben markieren (**Munition**), einem wartenden Schützen eine Scheibe zuteilen, das Standblatt drucken |
| **Resultaterfassung** | das Suchfeld, die Schützen auf allen Scheiben, die Liste der Beendeten | Resultate und Auszeichnungen auf jeder Scheibe erfassen |

Alles andere bleibt bei den Benutzern des Vereins: Schützen anmelden,
Programme bearbeiten, Ranglisten, Statistik, Löschen.

## Einen Link erstellen {#create}

1. Gib einen Namen ein (wer den Link bekommt, z. B. *Warner Stand A*), wähle
   die Rolle und beim Warner den Stand und die Scheiben (`1-5, 8`). Ein
   Warner betreut meist mehrere Scheiben. Warner gibt es nur, solange der
   Anlass Scheiben zuteilt (Modul *Scheibenzuteilung*).
2. **Helferlink erstellen**; kopiere den Link und schicke ihn dem Helfer
   (z. B. per Messenger). Er wird nur einmal angezeigt.
3. Der Helfer öffnet den Link; er gilt, bis er widerrufen oder der Anlass
   gelöscht wird. `?language=de` (`fr`, `it`, `en`) übersteuert die Sprache
   des Anlasses.

**widerrufen** in der Liste stoppt einen Link sofort; *zuletzt gesehen* sagt, wann er zuletzt benutzt wurde. Ein Warner, der einen Schützen auf einer anderen Scheibe öffnet, bekommt *Diese Scheibe gehört nicht zu dir*.

## Munition {#ammunition}

**Munition** auf einer wartenden oder schiessenden Anmeldung hält die Zeit fest, zu der die Munition ausgegeben wurde; den Knopf gibt es nur, wenn das Modul *Munitionsausgabe* an ist (Luftpistolen- und Luftgewehrschützen bringen meist ihre eigene mit). Das Dashboard zeigt bei solchen Anmeldungen ein grünes *Munition*-Abzeichen, so sehen Ausgabe und Rangeur auf einen Blick, wer bedient ist.
