# Show the rankings on a screen in the hall

You need this when shooters and spectators should follow the live rankings on a screen at the range. The screen needs a browser and internet, no login. Setting it up needs the permission *Configure events* (see [Roles](../user-guide/users.md#roles)).

1. Under **Configuration** → **Modules** of the event, switch on *Kiosk displays*.
2. Open **Kiosk displays** in the configuration.
3. Enter a name for the screen (e.g. *Hall 1*), choose **Dark** or **Light** and, if the screen should not use the event's language, a language. On an event with several ranges, choose the range the screen shows, or *All ranges*.
4. Click **Create kiosk link** and **copy** the link. It is shown only once.
5. Open the link in a browser on the screen and click the **Fullscreen** button at the top right.

![Kiosk screen](../images/screenshots/kiosk.png){ .screenshot }

The screen updates live and reconnects by itself after a network interruption.

## Choose what the screen shows {#content}

- A programme shows only while **Show on kiosk displays** is on on its programme page (it is on for new programmes). Staff pages, PDF and CSV rank it anyway.
- A classification shows with **Show on kiosk displays** in its editor (see [Classifications](../user-guide/classifications.md)).
- With lane assignment, **Show lane status on the rankings** on **Kiosk displays** adds the lane status above the rankings.

## Stop a screen {#revoke}

Click **revoke** next to the screen in the list on **Kiosk displays**: the link stops working immediately. *last seen* tells when a screen last fetched data. The parts of the link are listed in [Links](../reference/links.md).
