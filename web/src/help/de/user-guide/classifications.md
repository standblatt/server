# Wertungen

Eine Programm-Rangliste rangiert ein Programm. Eine **Wertung** rangiert
Schützen über mehrere Programme oder über mehrere
[Durchgänge](programmes.md#attempts) eines Programms, z. B. eine
Sektionsmeisterschaft als Summe von drei Programmen oder eine
Vereinsmeisterschaft über die besten drei Durchgänge eines
Jahresprogramms. Wertungen erscheinen unter **Ranglisten**, auf
Kiosk-Anzeigen und als Bericht; wie sie mit den Programm-Ranglisten
zusammenhängen, erklärt [Wie die Ranglisten
zusammenhängen](../explanation/rankings.md#classifications).

## Eine Wertung erstellen {#create}

Ist das Modul *Wertungen* eingeschaltet (**Konfiguration** → **Module**),
öffnet **Neue Wertung** auf **Wertungen** in der Konfiguration (Admins) den
Editor:

- **Programme**: jedes Programm, das zählt, mit seinen **Zählenden
  Durchgängen**. Bei 3 zählen die besten drei Durchgänge (oder die ersten
  oder letzten drei, je nach *Zählender Durchgang* am Programm). Das Total
  eines Schützen ist die Summe aller gezählten Durchgänge.
- **Nur die besten Programme zählen**: Bei mehreren Programmen zählen nur so
  viele mit den meisten Punkten, z. B. die besten zwei von drei; 0 zählt
  alle.
- **Nur Schützen mit allen zählenden Resultaten rangieren**: Schützen, denen
  noch ein Programm oder ein Durchgang fehlt, stehen ohne Rang unter der
  Rangliste. Ohne diese Option wird jeder nach dem rangiert, was er hat.
- **Stichentscheid**: Bei gleichem Total entscheiden die angekreuzten
  Kriterien in ihrer Reihenfolge; die Pfeile ändern sie. *Passen von der
  letzten rückwärts* vergleicht alle gezählten Passen, *Mehr höchste
  Einzelschüsse* braucht [Einzelschüsse](programmes.md) (mehr Zehner, dann
  mehr Neuner …), *Alter* folgt der
  [SSV-Regel](../reference/rules.md#ranking). Standard ist die
  SSV-Reihenfolge.
- **Klassen**: getrennte Ranglisten, je für einige Alterskategorien,
  z. B. *Junioren* mit U17 bis U21 und *Elite* mit den übrigen. Eine
  Alterskategorie gehört höchstens zu einer Klasse; Schützen, deren
  Kategorie in keiner Klasse ist, erscheinen nicht. Ohne Klassen sind alle
  in einer Rangliste.
- **Auf Kiosk-Anzeigen zeigen**.

## Beispiele {#examples}

| Anlass | Programme | Einstellungen |
|---|---|---|
| Schützenfest, beste zwei Stiche | drei Stiche | *Nur die besten Programme zählen* 2 |

Schritt für Schritt: [Eine Sektionsmeisterschaft über mehrere Programme
rangieren](../how-to/section-championship.md) und [Ein Jahresprogramm über
die Saison durchführen](../how-to/annual-programme.md#classification).

**Für nächstes Jahr kopieren** kopiert die Wertungen mit den Programmen. Wird ein Programm gelöscht, verschwindet es aus den Wertungen.
