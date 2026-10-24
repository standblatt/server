# Printing

Standblatt renders all documents as PDF. Every document with its format, content and printer setting: [Documents](../reference/documents.md#documents).

Every Print button prints on the range's printer if one is set up and its print agent is online; otherwise the PDF opens (see [When there is no printer](../explanation/printing.md#fallback)).

Why a print agent and how it works: [How printing works](../explanation/printing.md).

## Setting up printers

**Ranges → Set up the printers** (permission *Ranges and printers*) has three parts:

![Printing](../images/screenshots/printing.png){ .screenshot }

- **Print agents**: the computers that print for this range, each with its version and operating system, *online* or when it was last seen, and its number of printers. **Add print agent** adds a headless agent and shows its token once, **Connect by code** connects the [app](#app) by its [code](#connect), **revoke** removes an agent immediately.
- **Documents**: per document the printer (*browser (PDF)* keeps opening the PDF) and the options the agent reports for that printer; on Linux further `lp` options such as `fit-to-page` go in the text field. *Certificates* are the certificates from the rankings, *Print forms (pre-printed)* every other form of the club. **Save** stores the row. **Test page** prints a small page (90 × 29 mm) with club, range, printer, options and the time on the printer of that row, also before saving; it needs an online agent.
- **Recent jobs** shows what was printed and whether it worked.

A job that the agent does not confirm within two minutes is retried; after three attempts it fails.

Step by step: [Connect a printer at the range](../how-to/connect-printer.md).

## The Standblatt Print Agent app {#app}

The app is the simple way for a laptop or PC next to the printers: no token and no settings to type.

Download, install and connect it: [Connect a printer at the range](../how-to/connect-printer.md#app).

![The Standblatt Print Agent app](../images/app/print-agent.png){ .screenshot }

The window shows for each range whether it prints, with its last jobs, and below the printers of the computer, each with a **Test page** that prints on that computer without Standblatt. **Advanced** has *Start with the computer*, the address of Standblatt (change it only for a test server) and removing a range. Closing the window keeps the app running in the system tray: the icon is green when every range prints, red otherwise.

On Windows and on a Mac the app tells when a new version is out; **Update now** installs it and restarts the app. To give up a computer, remove its range under **Advanced**, revoke its agent on **Ranges → Set up the printers** and uninstall the app (on Windows under *Settings → Apps*).

## Connecting a computer by code {#connect}

![A code in the app](../images/app/print-agent-code.png){ .screenshot }

**Connect a print agent** (opened by the app with a code such as `K7QD-4MZP`) shows the computer and asks for club, range and name; only members who may manage ranges can connect a computer, and a code is valid for ten minutes. **Connect by code** on **Ranges → Set up the printers** connects it to that range under its own name. Steps: [Connect a printer at the range](../how-to/connect-printer.md#app).
