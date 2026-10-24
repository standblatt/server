# Wie die Ranglisten zusammenhängen

Ein Anlass kann seine Schützen auf mehrere Arten gleichzeitig rangieren: pro
Programm, über mehrere Programme oder Durchgänge, nach Verein und nach
Alterskategorie für Sonderpreise. Das sind keine separaten Wettkämpfe mit
separaten Daten. Alle gehen von denselben Resultaten aus, und die meisten
von derselben Programm-Rangliste.

## Die Programm-Rangliste {#programme}

Die Programm-Rangliste ist die Grundlage. Jedes Programm hat eine Rangliste
aller, die es geschossen haben, geordnet nach den [Regeln des
SSV](../reference/rules.md#ranking): Total, dann die Passen von der letzten
rückwärts, dann Mouchen, dann Alter. Es gibt nur eine Umsetzung dieser
Reihenfolge, und sie speist die Ranglistenseite, die Kiosk-Bildschirme, das
Ranglisten-PDF und -CSV, die Sonderpreise und den Rang, den ein Schütze in
*Mein Training* sieht. Der Platz eines Schützen ist deshalb überall
derselbe.

Die Ränge sind fortlaufend (1, 2, 3, … auch bei gleichen Resultaten): Der
Stichentscheid entscheidet fast jeden Gleichstand, so teilen sich zwei
Schützen nie einen Rang.

**Varianten** sind separate Programme mit gleichem Namen, zum Beispiel *Sport* und *Aufgelegt*, und jede hat ihre eigene Rangliste: Aufgelegt und frei geschossen wird nicht verglichen.

Ein **kombiniertes** Programm übernimmt Passen aus einem anderen Programm,
statt separat geschossen zu werden. Seine Resultate sind Kopien der
gewählten Passen seiner Quelle, deshalb wird es wie jedes andere Programm
rangiert; ändert sich die Quelle, folgt die Kopie und mit ihr die Rangliste.

## Durchgänge und der zählende Durchgang {#attempts}

Ein Programm, das mehrere Durchgänge erlaubt, behält jeden Durchgang als
Stich mit eigenen Resultaten. Eine Rangliste führt jeden Schützen aber nur
einmal auf. Der **zählende Durchgang** des Programms (der beste, der erste
oder der letzte) entscheidet, welcher Durchgang für den Schützen steht, und
die Rangliste zeigt, welcher es ist, z. B. *2/3*.

Derselbe zählende Durchgang gilt überall, wo ein Resultat pro Schütze
gebraucht wird: in der Rangliste, der Vereinskonkurrenz, den Sonderpreisen
und beim Rang in *Mein Training*. Bei Auszeichnungen ist es anders, weil
eine Auszeichnung zu einem Resultat gehört und nicht zu einer Rangliste:
Jeder Durchgang kann seine eigene Auszeichnung erhalten, ausser *Nur der
zählende Durchgang erhält eine Auszeichnung* ist an; dann wandert die
Auszeichnung zu dem Durchgang, der gerade zählt.

Kombinationen übernehmen die Passen eines Resultats; bei mehreren
Durchgängen wäre unklar, welches, deshalb können Programme mit mehreren
Durchgängen nicht kombiniert werden.

## Wertungen {#classifications}

Eine Wertung rangiert Schützen über mehrere Programme oder über mehrere
Durchgänge eines Programms: die Summe von drei Programmen für eine
Sektionsmeisterschaft, die besten drei Durchgänge eines Jahresprogramms für
eine Vereinsmeisterschaft. Sie nimmt die zählenden Durchgänge jedes
Programms, zählt zusammen, was zählt, kann nur die besten Programme zählen
und rangiert mit ihrem eigenen Stichentscheid und ihren eigenen Klassen von
Alterskategorien.

Eine Wertung liest nur Resultate. Sie ändert nie ein Resultat, eine
Auszeichnung oder eine Programm-Rangliste, und ein Anlass kann neben seinen
Programm-Ranglisten mehrere Wertungen haben. Die Programm-Rangliste selbst
lässt sich als einfachste Wertung sehen: ein Programm, ein Durchgang, der
Stichentscheid des SSV und keine Klassen. Die Einstellungen stehen unter
[Wertungen](../user-guide/classifications.md).

## Vereinskonkurrenz {#club-competition}

Die Vereinskonkurrenz rangiert Vereine, nicht Schützen. Sie nutzt nur
Programme mit *Zählt für die Vereinskonkurrenz* und nur abgeschlossene
Resultate von Schützen, die sich für einen Verein angemeldet haben; der
Verein ist der bei der Anmeldung gewählte, nicht der Mandantenverein. Jedes
Resultat zählt mit seinen ersten zwei Passen, und das Resultat eines Vereins
ist der Durchschnitt seiner besten Resultate, mit einem kleinen Bonus für
jeden weiteren Teilnehmer (siehe
[Vereinskonkurrenz](../reference/rules.md#club-competition)). Weil sie ihre
eigene Punktzahl hat, folgt die Vereinskonkurrenz nicht der
Programm-Rangliste.

## Sonderpreise {#special-prizes}

Sonderpreise gehen an die besten Ränge jeder Alterskategorie. Sie folgen
sehr wohl der Programm-Rangliste: Sie gehen sie pro Alterskategorie von oben
nach unten durch und überspringen Schützen, deren Verein oder Kategorie zu
wenige Teilnehmer hat (siehe
[Sonderpreise](../reference/rules.md#special-prizes)). Ihre Einstellungen
lassen sich ändern, während der Anlass läuft, weil sie nur die Rangliste
lesen.

## Welche Rangliste was verwendet {#overview}

| Rangliste | Rangiert | Aus | Zählender Durchgang | Reihenfolge |
|---|---|---|---|---|
| Programm-Rangliste | Schützen | einem Programm | ja | Regeln des SSV |
| Wertung | Schützen, pro Klasse | mehreren Programmen oder Durchgängen | eigene Anzahl pro Programm | eigener Stichentscheid, standardmässig SSV |
| Vereinskonkurrenz | Vereine | Programmen, die für die Vereinskonkurrenz zählen | ja | Vereinsresultat |
| Sonderpreise | Schützen, pro Alterskategorie | der Programm-Rangliste | ja | die Programm-Rangliste |

Auszeichnungen sind keine Rangliste: Ein Schütze erreicht eine Auszeichnung
mit seinem eigenen Total in seiner Alterskategorie, egal wie die anderen
schiessen.

## Warum die Passen gesperrt sind {#locked}

Resultate werden pro Passe gespeichert. Sobald Schützen für ein Programm
angemeldet sind, würde eine Änderung an der Anzahl Passen, ihren Schüssen,
ihrem Maximum oder daran, ob sie Mouchen oder Einzelschüsse erfassen, die
Bedeutung der gespeicherten Resultate ändern: Ein Total von 96 von 100 ist
kein Total von 200, und jede Rangliste, die sie vergleicht, wäre
falsch. Standblatt legt die Passen deshalb fest, sobald ein Schütze
angemeldet ist.

Alles, was die Bedeutung eines Resultats nicht ändert, bleibt offen: Namen,
Preise, Werte und Limiten der Auszeichnungen, die Optionen. Alterskategorien
und Auszeichnungsstufen lassen sich hinzufügen, aber nicht entfernen, weil
Anmeldungen danach berechnet und Auszeichnungen danach vergeben wurden.

## Live und endgültig {#live}

Alle Ranglisten werden jedes Mal aus den Resultaten berechnet, wenn sie
angezeigt werden, deshalb erscheint ein gespeichertes Resultat sofort auf
der Ranglistenseite und den Kiosk-Bildschirmen. Solange der Anlass läuft,
können sie sich noch ändern. Das [Abschliessen](closing.md) des Anlasses
macht sie endgültig.
