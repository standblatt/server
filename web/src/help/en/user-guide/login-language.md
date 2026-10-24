# Login and language

## Logging in

Log in with your email address (or user name) and password: you chose the password when you accepted the invitation of your club, or you received it from your club admin or the platform operator. Club accounts are not created by sign-up; only [My training](../my-training/account.md) has its own sign-up. Logged in, you change your password yourself under [Account](#account).

After the login you land on the start page with a card for every club you belong to and one for [My training](../my-training/index.md), which every user has; **Account** and **Log out** are below the cards. Without a club membership the start page takes you straight to My training.

![Login](../images/screenshots/login.png){ .screenshot }

![Start page](../images/screenshots/home.png){ .screenshot }

A login stays valid for seven days of inactivity and at most 30 days; with **Stay logged in for 90 days** ticked, it lasts 90 days on that device. After too many failed attempts, logins from the same address are blocked for a while.

If your account belongs to several clubs, you choose the club after logging in; **Switch club** in the user menu goes back to that choice.

## Forgotten password {#password}

If your account has a confirmed email address, use **Forgot your password?** on the login page. You get a link that is valid for one hour. Setting a new password logs you out everywhere. Without a confirmed email address, ask the operator to set a new password in the admin panel.

## Account {#account}

**Account** in the user menu (bottom of the sidebar), on the start page and in My training's sidebar opens your account settings:

- **Login**: your user name (it cannot be changed) and your email address. **Add** or **Change** asks for your password and sends a confirmation link; until you confirm, the address counts as *not confirmed* and **Send again** sends a new link. A confirmed address lets you reset a forgotten password yourself.
- **Password**: enter the current and the new password (at least 10 characters). Other devices are logged out, the one you are on stays.
- **Language and theme**: the interface language of your account and light or dark mode.
- **Second factor**: see below.
- **API tokens**: see below.
- **Delete account**: only for accounts that belong to no club; the training log, goals and settings are deleted for good. Club accounts are removed by the club admin.

![Account](../images/screenshots/account.png){ .screenshot }

## API tokens {#api-tokens}

Want to show the rankings on the club website, import sessions from another app or analyse data with your own script? Create a token under **Account → API tokens** with **New token**:

- a **name** that tells you later what uses it,
- a **preset** fills in the permissions for a common use: **Rankings for a website**, **Training app** or **Read only**,
- the **clubs** the token may use (none: only your training log),
- the **permissions**: *Read my training log*, *Log and change my training sessions* and, for the chosen clubs, the [permissions of roles](users.md#roles) your own role has there,
- when it **expires**: in 30 or 90 days, in a year, or never.

The token (`sb_…`) is shown only once; copy it right away. Scripts send it as `Authorization: Bearer sb_…` (or `X-API-Key: sb_…`). It may do what you chose and never more than your role in the club: if you lose a permission or leave the club, the token loses it at once. It never reaches your password, email, account or other tokens. `GET /api/token` tells a script its name, permissions and clubs. A token makes at most 600 requests in five minutes. The list shows when each token was last used; revoking one stops it at once. Developers find the endpoints under `/api/docs`.

## Second factor {#totp}

An authenticator app (Google Authenticator, Aegis, 1Password and the like) can protect your login: after the password, Standblatt asks for the six-digit code the app shows. Set it up under [Account](#account): **Set up**, scan the QR code, enter the current code, **Switch on**. **Switch off** asks for your password. Lost the phone? The operator can reset the second factor in the admin panel. The operator can also switch the feature off for the whole server.

## Language

The interface is available in English, German, French and Italian. Without a choice, Standblatt follows the browser's language. **Language** in the user menu (bottom of the sidebar) sets it for your account on every device.

Printed documents, CSV exports and kiosk screens do not follow your choice but the language of the event, so that every printout of an event looks the same (see [Events](events.md#settings)).

## Dark mode

The sun/moon button in the header switches between light and dark mode; the choice is kept in the browser.

![Dashboard in dark mode](../images/screenshots/dashboard-dark.png){ .screenshot }

## Version

The bottom of the sidebar shows the version of the server, which helps when reporting a problem; in My training on a phone it is under **Settings**. When a new version is released, open pages offer **Reload** within a minute and load it on the next page change at the latest.
