# Print agent

The print agent is a small program on a computer next to the printers of a range. It only makes outgoing HTTPS requests to Standblatt: it reports the printers and their options, waits for jobs (long polling) and prints them. No port needs to be opened at the range.

It lives in its own repository, [standblatt/print-agent](https://github.com/standblatt/print-agent), with its own [releases](https://github.com/standblatt/print-agent/releases/latest):

- **Desktop app** *Standblatt Print Agent* for Windows, macOS and Linux: install it on the computer next to the printers, click *Connect a range* and confirm the code it shows on **Connect a print agent** in Standblatt (see [Printing](../../web/src/help/en/user-guide/printing.md#connect)). No token to copy.
- **Headless agent** for Linux with CUPS (e.g. a Raspberry Pi), run as a systemd service with the token from **Ranges → Set up the printers → Add print agent**. Installation and configuration are in the [README](https://github.com/standblatt/print-agent#readme).

## Behaviour

- The agent reports its printers with the options it can set; the printing setup of the range chooses a printer and options per document.
- A job that is not confirmed within two minutes goes back to the queue; after three attempts it fails. Jobs are deleted after a week.
- The app treats an agent as online when it was seen within the last minute. Without an online agent, documents open as PDF in the browser.
- Revoking the agent in the app invalidates its token immediately.
