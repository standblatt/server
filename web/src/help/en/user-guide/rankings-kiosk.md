# Rankings and kiosk

## Rankings

**Rankings** shows a live ranking per programme, ordered by the [SSV rules](../reference/rules.md#ranking). The first three ranks are highlighted.

How the programme rankings, classifications, club competition and special prizes relate: [How rankings fit together](../explanation/rankings.md).

![Rankings](../images/screenshots/rankings.png){ .screenshot }

- **PDF** opens the ranking (A4) for printing or sharing, **CSV** downloads it for Excel. Rankings do not go to the range's printer.
- Programmes switched off for kiosk displays are still ranked here and marked *not on kiosk*.
- [Classifications](classifications.md) follow the programmes, one table per class; their **PDF** and **CSV** are the report *Classifications*.
- **Certificates** prints certificates for the best ranks onto a [print form](print-forms.md#certificates) of the club (see [Print certificates for the best ranks](../how-to/certificates.md)); the button appears once the club has a form that prints a rank.

## Kiosk screens

A kiosk link shows the rankings of one event on a screen in the hall, without login; set one up with [Show the rankings on a screen in the hall](../how-to/kiosk.md).

**Kiosk displays** in the configuration (module *Kiosk displays*) lists the links. A new one takes a name, **Dark** or **Light**, a language and, on an event with several ranges, a range; **Create kiosk link** shows the link only once.

![Kiosk screen](../images/screenshots/kiosk.png){ .screenshot }

- The kiosk shows only programmes and classifications with *Show on kiosk displays* switched on, and, with the module *Lane assignment*, the lane status if *Show lane status on the rankings* is on **Kiosk displays**.
- The link sets theme, language and range (see [URL parameters](../reference/links.md#parameters)).
- The **Fullscreen** button at the top right switches the screen to full screen.
- The screen updates live and reconnects by itself after a network interruption. When a new version of Standblatt is released, it reloads by itself.
- **revoke** in the list stops a link immediately. *last seen* tells when a screen last fetched data.
