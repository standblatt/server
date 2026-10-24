# Segnaposto dei modelli

Un [modello di stampa](../user-guide/print-forms.md) stampa al posto di un
segnaposto come `{{shooter.name}}` il suo valore. Un segnaposto senza valore
non stampa nulla; lo stesso vale per uno che la [fonte di
dati](../user-guide/print-forms.md#source) del modello non conosce. In un
elemento ripetuto, `#` conta 1, 2, 3, …: `{{s1.shot#}}` stampa
`{{s1.shot1}}`, `{{s1.shot2}}`, ….

## Ogni fonte di dati {#common}

| Segnaposto | Campo | Esempio |
|---|---|---|
| `{{shooter.name}}` | Nome (prima il cognome) | Example Anna |
| `{{shooter.lastname}}` | Cognome | Example |
| `{{shooter.firstname}}` | Nome | Anna |
| `{{shooter.license}}` | Numero di licenza | 123456 |
| `{{shooter.birthday}}` | Data di nascita | 02.06.1985 |
| `{{shooter.birthyear}}` | Anno di nascita | 1985 |
| `{{shooter.postalcode}}` | NPA | 9999 |
| `{{shooter.city}}` | Località | Hilltown |
| `{{category}}` | Categoria d’età | Élite |
| `{{club}}` | Società | SG Hilltown |
| `{{club.ssv}}` | Numero FST | 1.23.4.56 |
| `{{event.name}}` | Manifestazione | Programma annuale 2027 |
| `{{event.date}}` | Data (giorno principale) | 12.05.2027 |
| `{{event.organizer}}` | Organizzatore | SG Hilltown |
| `{{today}}` | Oggi, nel fuso orario della manifestazione | 12.05.2027 |

## Risultati di un programma {#results}

Modelli dei programmi e carte corona.

| Segnaposto | Campo | Esempio |
|---|---|---|
| `{{programme}}` | Programma, con la sua variante | Obbligatorio |
| `{{lane}}` | Bersaglio | 7 |
| `{{registration}}` | Numero d’iscrizione | 1042 |
| `{{total}}` | Totale | 85 |
| `{{mouches}}` | Mouches | 3 |
| `{{s1}}`, `{{s2}}`, … | Serie 1, 2, … | 47 |
| `{{s1.mouches}}`, … | Mouches serie 1, … | 1 |
| `{{s1.shot1}}`, `{{s1.shot#}}`, … | colpi singoli della serie 1 (serie con **Colpi singoli**) | 5 |
| `{{award}}` | Distinzione | Carta |
| `{{award.value}}` | Valore della distinzione | CHF 12.00 |
| `{{award.value.words}}` | Valore in lettere, nella lingua della manifestazione | dodici franchi |

`{{award.value}}` e `{{award.value.words}}` restano vuoti a meno che *Stampa
il valore della distinzione sulle carte corona* sia attivo e la distinzione
abbia un valore.

## Un’iscrizione {#registration}

Strisce di controllo e documenti d’iscrizione.

| Segnaposto | Campo | Esempio |
|---|---|---|
| `{{category.short}}` | Categoria d’età, abbreviata | E |
| `{{lane}}` | Bersaglio | 7 |
| `{{registration}}` | Numero d’iscrizione | 1042 |
| `{{registration.code}}` | Codice QR dell’iscrizione (il [codice d’iscrizione](codes.md)) | REGCODEAAAA1042AAAAL123456 |
| `{{programmes}}` | Programmi, separati da virgola | Obbligatorio, Tiro in campagna |
| `{{passes}}` | Serie: `P1 P2` per programma, `P0` per un programma che prende le sue serie da un altro, `*` sulle serie che vengono prese | P1 P2 + P0 |

## Un rango {#rank}

Diplomi.

| Segnaposto | Campo | Esempio |
|---|---|---|
| `{{rank}}` | Rango | 1 |
| `{{total}}` | Totale | 98 |
| `{{ranking}}` | Classifica: il programma o la classifica generale | Campionato di sezione |
| `{{class}}` | Classe della classifica generale; vuoto per un programma | Juniori |

Un diploma di un programma compila anche ogni segnaposto di [Risultati di un
programma](#results) dal tentativo che conta.

## Elementi {#elements}

| Elemento | Stampa |
|---|---|
| *Testo libero* | il suo testo, con eventuali segnaposto |
| *Codice QR* | il suo valore come codice QR |
| *Codice a barre* | il suo valore come codice a barre Code 128 |
| *Linea* | una linea |
| immagine | un’immagine della [biblioteca di immagini](../user-guide/print-forms.md#images) della società |
