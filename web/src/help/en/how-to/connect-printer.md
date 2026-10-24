# Connect a printer at the range

You need this when lane sheets, control strips and wreath cards should come out of the printer at the range instead of opening as PDF. It needs the permission *Ranges and printers* (see [Roles](../user-guide/users.md#roles)).

A computer next to the printers runs the print agent: a Windows PC or a Mac with the Standblatt Print Agent app, or a Linux computer such as a Raspberry Pi with the headless agent. Why it works this way is in [Printing](../explanation/printing.md).

## With the app on Windows or a Mac {#app}

1. Download it from the [latest release](https://github.com/standblatt/print-agent/releases/latest): `Standblatt-Print-Agent-Setup.exe` for Windows, `standblatt-print-agent-app-darwin-arm64.dmg` for a Mac with Apple silicon (`-amd64.dmg` for an Intel Mac), the `.deb` or `.rpm` for Linux.
2. Install it. The app is not signed yet: Windows asks once, **More info → Run anyway**; on a Mac drag it into *Applications* and open it the first time with a right click → **Open**. On Windows it installs for the current user without admin rights, with a shortcut in the start menu and on the desktop, and starts with Windows.
3. In the app, click **Connect a range**. It shows a code such as `K7QD-4MZP` and opens **Connect a print agent** in the browser.
4. Log in if asked, check the computer named on the page, choose the **Club**, the **Range** and the agent's **Name**, and **Connect**. The code is valid for ten minutes.

![A code in the app](../images/app/print-agent-code.png){ .screenshot }

If the browser on that computer is not at hand, open **Ranges → Set up the printers** of the range on any other computer, type the code into the field next to **Connect by code** and click it: the computer is connected to that range under its own name.

For a second range on the same computer, e.g. *Pistol 25 m* and *Pistol 50 m*, click **Connect another range** in the app and repeat steps 3 and 4.

## With the headless agent on Linux {#headless}

1. Open **Ranges** and click **Set up the printers** of the range.
2. Under **Print agents**, enter a name (e.g. *Pi hall 1*) and click **Add print agent**.
3. Click **Copy** to copy the two lines `STANDBLATT_URL` and `STANDBLATT_AGENT_TOKEN`. The token is shown only now.
4. Install and start the agent with these settings as described in the [README of the print agent](https://github.com/standblatt/print-agent#readme).

For a second range on the same computer, add an agent there too and give the agent both tokens, separated by a comma.

## Choose the printers {#printers}

1. Open **Ranges → Set up the printers** of the range. The agent shows *online* with its version and operating system; its printers appear once it has connected.
2. For each document, choose the printer, then its options (tray, paper, two-sided …, as the agent reports them for that printer; on Linux further `lp` options such as `fit-to-page` go in the text field).
3. Click **Test page** in the row: a small page (90 × 29 mm) with club, range, printer, options and the time prints with the options of that row. A label printer prints it as a label, other printers in the top-left corner of the sheet.
4. Click **Save** in the row. A document left on *browser (PDF)* keeps opening as PDF.

![Printing](../images/screenshots/printing.png){ .screenshot }

*Certificates* are the certificates from the rankings, *Print forms (pre-printed)* every other [print form](../user-guide/print-forms.md) of the club. Rankings are never sent to a printer.

## Check that it prints {#jobs}

Print something, e.g. a lane sheet at the next registration, and look at **Recent jobs** on the same page: *done* means it printed, *failed* shows the error next to it. A job that the agent does not confirm within two minutes is retried; after three attempts it fails.

To give up a computer, click **revoke** next to its agent; with the app, also remove the range under **Advanced** and uninstall the app. All screens of the printing page are described in [Printing](../user-guide/printing.md).
