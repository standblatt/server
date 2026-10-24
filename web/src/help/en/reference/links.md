# Links without login

These links open a page without logging in. Whoever has the link can use it, so send it only to the people it is meant for.

| Link | Made by | Shows | Valid |
|---|---|---|---|
| Kiosk link `/kiosk/<secret>` | **Create kiosk link** on **Kiosk displays** in the configuration | the live rankings of the event: programmes and classifications with *Show on kiosk displays*, and the lane status with *Show lane status on the rankings* | until **revoke** or the event is deleted |
| Helper link `/h/<secret>` | **Create helper link** on **Helper links** in the configuration | a reduced desk of the event for one [helper role](permissions.md#helper-links) | until **revoke** or the event is deleted |
| Booking link `/b/<code>` | **Switch on** on **Reservations** in the configuration | the slots of the event with the free lanes per range, and the request form | until **New link** or **Switch off**; requests only until the deadline |
| Booking manage link `/b/manage/<secret>` | mailed to the club after **Send request** | the status of the club's request; change or cancel it | changes and cancelling until the deadline |
| Invitation link `/invite/<secret>` | **Invite by email** on **Users** | the club and the role; accept with an account or a new password | seven days; ends when accepted, withdrawn or sent again |
| Connect page `/connect?code=<code>` | opened by the Standblatt Print Agent app | the computer of the code; club, range and name of the agent | the code: ten minutes; needs a login with *Ranges and printers* |
| Email confirmation `/verify?token=…` | adding or changing the email address, sign-up | confirms the address | seven days |
| Password link `/reset-password?token=…` | **Forgot your password?**; the operator in the admin panel | sets a new password | one hour; from the operator seven days |

The secret part of kiosk, helper, manage and invitation links is shown or mailed once and stored only as a hash: a lost link cannot be shown again, only replaced.

## URL parameters {#parameters}

| Parameter | Link | Effect |
|---|---|---|
| `?theme=dark`, `?theme=light` | kiosk | dark or light display |
| `?language=de` (`fr`, `it`, `en`) | kiosk, helper | language of the page instead of the event's language |
| `?range=<id>` | kiosk | on an event with several ranges: only that range's lanes, its programmes and the programmes without a range |
| `?code=` | connect | the code to confirm, filled in by the app |

**Create kiosk link** writes the theme, the language and the range chosen there into the link.
