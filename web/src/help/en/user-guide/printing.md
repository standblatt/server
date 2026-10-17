# Printing

Standblatt renders all documents as PDF:

| Document | Format | When |
|---|---|---|
| Lane sheet | A5 | at registration, with QR code and programme barcodes |
| Control strip | label 90 × 29 mm | at registration, stuck on the target card |
| Result sheet | A5 | when all programmes have results |
| Wreath card | A6, printed onto pre-printed cards | after choosing a wreath card award |
| Ranking | A4 | from the rankings page, as PDF only (not sent to the range's printer) |

Every Print button prints on the range's printer if one is set up and its print agent is online; otherwise the PDF opens in the browser and can be printed from there. On a phone or tablet the PDF goes to the share sheet instead (print, save to files, send), so the app stays open; if the phone asks for another tap first, a note with a **Share** button appears.

## Setting up printers

**Ranges → Printing** (admins):

![Printing](../images/screenshots/printing.png){ .screenshot }

1. **Add print agent**, give it a name and copy the token. Install the agent on a computer next to the printers (see [Print agent](../admin-guide/print-agent.md)).
2. Once the agent is online, it reports its printers. Choose the printer for each document and its options (tray, paper size, media type, further `lp` options such as `fit-to-page`).
3. **Recent jobs** shows what was printed and whether it worked.

A job that the agent does not confirm within two minutes is retried; after three attempts it fails. **revoke** removes an agent immediately.
