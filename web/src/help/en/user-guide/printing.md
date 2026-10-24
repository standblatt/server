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

**Ranges → Set up the printers** (admins):

![Printing](../images/screenshots/printing.png){ .screenshot }

1. Connect a computer next to the printers: a Windows PC or a Mac with the [Standblatt Print Agent app](#app) by [code](#connect), or a Linux computer such as a Raspberry Pi with the [headless agent](https://github.com/standblatt/print-agent): **Add print agent**, give it a name and copy the token. One computer can serve several ranges, e.g. *Pistol 25 m* and *Pistol 50 m*: connect it to each range (the app with **Connect another range**, the headless agent with both tokens, separated by a comma).
2. Once the agent is online, it reports its printers and the list shows its version and operating system. Choose the printer for each document (*Certificates* are the certificates from the rankings, *Print forms* every other form of the club) and its options: the options the agent reports for that printer (on Linux tray, paper size and media type, plus further `lp` options such as `fit-to-page` in the text field; on Windows tray, paper, two-sided, scaling and colour).
3. **Test page** prints a small page (90 × 29 mm, the size of the control strip) on the printer of that row with the options chosen there, also before saving: club, range, printer, options and the time. A label printer prints it as a label, other printers in the top-left corner of the sheet. It needs an online agent.
4. **Recent jobs** shows what was printed and whether it worked.

A job that the agent does not confirm within two minutes is retried; after three attempts it fails. **revoke** removes an agent immediately.

## The Standblatt Print Agent app {#app}

The app is the simple way for a laptop or PC next to the printers: no token and no settings to type.

1. Download it from the [latest release](https://github.com/standblatt/print-agent/releases/latest): `Standblatt-Print-Agent-Setup.exe` for Windows, `standblatt-print-agent-app-darwin-arm64.dmg` for a Mac with Apple silicon (`-amd64.dmg` for an Intel Mac), the `.deb` or `.rpm` for Linux.
2. Install it. The app is not signed yet: Windows asks once, **More info → Run anyway**; on a Mac drag it into *Applications* and open it the first time with a right click → **Open**. On Windows it installs for the current user without admin rights, with a shortcut in the start menu and on the desktop, and starts with Windows.
3. Click **Connect a range** and confirm the code in Standblatt (see [below](#connect)).

![The Standblatt Print Agent app](../images/app/print-agent.png){ .screenshot }

The window shows for each range whether it prints, with its last jobs, and below the printers of the computer, each with a **Test page** that prints on that computer without Standblatt. **Advanced** has *Start with the computer*, the address of Standblatt (change it only for a test server) and removing a range. Closing the window keeps the app running in the system tray: the icon is green when every range prints, red otherwise.

On Windows and on a Mac the app tells when a new version is out; **Update now** installs it and restarts the app. To give up a computer, remove its range under **Advanced**, revoke its agent on **Ranges → Set up the printers** and uninstall the app (on Windows under *Settings → Apps*).

## Connecting a computer by code {#connect}

![A code in the app](../images/app/print-agent-code.png){ .screenshot }

The Standblatt app on a computer next to the printers needs no token: it shows a code such as `K7QD-4MZP` and opens **Connect a print agent** in the browser. There, log in if asked, check the computer named on the page, choose the club and range and the agent's name, and **Connect**. Only members who may manage ranges can connect a computer; a code is valid for ten minutes. The app then connects by itself, and the agent appears under **Ranges → Set up the printers** with its version and operating system.

If the browser on that computer is not at hand, type the code on **Ranges → Set up the printers** of the range in **Connect by code** instead: the computer is connected to that range under its own name.
