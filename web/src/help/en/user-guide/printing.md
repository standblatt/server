# Printing

Standblatt renders all documents as PDF:

| Document | Format | When |
|---|---|---|
| Lane sheet | A5 | at registration, with QR code and programme barcodes |
| Control strip | label 90 × 29 mm | at registration, stuck on the target card |
| Result sheet | A5 | when all programmes have results |
| Wreath card | A6, printed onto pre-printed cards | after choosing a wreath card award |
| Print form | any size, printed onto pre-printed sheets | for programmes with a [print form](print-forms.md), once the results are in |
| Ranking | A4 | from the rankings page, as PDF only (not sent to the range's printer) |

Every Print button prints on the range's printer if one is set up and its print agent is online; otherwise the PDF opens in the browser and can be printed from there. On a phone or tablet the PDF goes to the share sheet instead (print, save to files, send), so the app stays open; if the phone asks for another tap first, a note with a **Share** button appears.

## Setting up printers

**Ranges → Printing** (admins):

![Printing](../images/screenshots/printing.png){ .screenshot }

1. **Add print agent**, give it a name and copy the token. Install the agent on a computer next to the printers (see [Print agent](../admin-guide/print-agent.md)). One computer can serve several ranges, e.g. *Pistol 25 m* and *Pistol 50 m*: add an agent on each range and give the agent both tokens, separated by a comma.
2. Once the agent is online, it reports its printers. Choose the printer for each document (*Certificates* are the certificates from the rankings, *Print forms* every other form of the club) and its options (tray, paper size, media type, further `lp` options such as `fit-to-page`).
3. **Test page** prints a small page (90 × 29 mm, the size of the control strip) on the printer of that row with the options chosen there, also before saving: club, range, printer, options and the time. A label printer prints it as a label, other printers in the top-left corner of the sheet. It needs an online agent.
4. **Recent jobs** shows what was printed and whether it worked.

A job that the agent does not confirm within two minutes is retried; after three attempts it fails. **revoke** removes an agent immediately.
