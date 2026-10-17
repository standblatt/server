# Print agent

The print agent is a small program on a computer next to the printers of a range, typically a Raspberry Pi with CUPS. It only makes outgoing HTTPS requests to Standblatt: it reports the CUPS printers and their options, waits for jobs (long polling) and prints them with `lp`. No port needs to be opened at the range.

## Installation

1. In the app, **Ranges → Printing → Add print agent**; copy the token (shown only once).
2. Download the binary for the computer (`standblatt-print-agent-linux-amd64` or `-arm64`, from the [GitHub release](https://github.com/standblatt/server/releases) of the version, or built with `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o standblatt-print-agent-linux-amd64 ./cmd/print-agent`) and install it:

    ```sh
    install standblatt-print-agent-linux-arm64 /usr/local/bin/standblatt-print-agent
    printf 'STANDBLATT_URL=https://app.standblatt.ch\nSTANDBLATT_AGENT_TOKEN=<token>\n' > /etc/standblatt-print-agent.env
    chmod 600 /etc/standblatt-print-agent.env
    cat > /etc/systemd/system/standblatt-print-agent.service <<'EOF'
    [Unit]
    Description=Standblatt print agent
    After=network-online.target cups.service
    Wants=network-online.target

    [Service]
    EnvironmentFile=/etc/standblatt-print-agent.env
    ExecStart=/usr/local/bin/standblatt-print-agent
    Restart=always
    RestartSec=5
    DynamicUser=yes

    [Install]
    WantedBy=multi-user.target
    EOF
    systemctl enable --now standblatt-print-agent
    ```

3. Back in the app, the agent shows as *online* with its printers. Assign a printer and its options to each document (see [Printing](../../web/src/help/en/user-guide/printing.md)).

## Configuration

| Variable | Default | Description |
|---|---|---|
| `STANDBLATT_URL` | required | base URL of Standblatt |
| `STANDBLATT_AGENT_TOKEN` | required | token from *Add print agent*; several tokens separated by commas serve several ranges from one computer |
| `STANDBLATT_LP` | `lp` | print command |
| `STANDBLATT_LPSTAT` | `lpstat` | used to list printers |
| `STANDBLATT_LPOPTIONS` | `lpoptions` | used to read printer options (trays, paper sizes, media types) |

The systemd unit runs the agent as a dynamic user and restarts it on failure. `journalctl -u standblatt-print-agent` shows its log.

## Behaviour

- Documents are printed with `lp -d <printer> -t <title> -o <options>`; the options come from the printing setup of the range.
- A job that is not confirmed within two minutes goes back to the queue; after three attempts it fails. Jobs are deleted after a week.
- The app treats an agent as online when it was seen within the last minute. Without an online agent, documents open as PDF in the browser.
- Revoking the agent in the app invalidates its token immediately.
