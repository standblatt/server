# Users

**Users** lists the accounts with access to the club, the open invitations and the roles of the club. It is open to roles with the permission *Users and club*, so at least to admins.

![Users](../images/screenshots/users.png){ .screenshot }

## Members and invitations {#invitations}

- **Invite by email**: paste one or many addresses (comma separated or one per line, up to 300) and choose the role. Every address gets a mail with a link valid for seven days; the list under the form says what happened to each address. Mails leave at about four per second, so a large list takes a minute.
- **Accepting**: the link opens a page that says which club invites with which role. Without an account the person chooses a password and is logged in at once; the email address is the login. Whoever has an account with the address logs in first and then accepts. Only then does the person appear among the members.
- **Invitations** below the members lists the open and the expired ones: **Send again** mails a new link (the old one stops working), **Withdraw** cancels the invitation.
- **Add user**: for an account without email, such as a shared account at the range: an existing account is added by its user name; for a new account also enter a password (at least 10 characters).
- **Role**: choose the role of a member in the list. A club always keeps at least one admin.
- **Remove** takes away the access to this club; the account itself stays and keeps its other clubs.

## Roles {#roles}

What a member may do depends on the role. The matrix under **Roles** has a column per role and a row per permission; a tick gives the permission at once.

![Roles](../images/screenshots/roles.png){ .screenshot }

| Permission | Allows |
|---|---|
| **See events and rankings** | events, programmes, rankings and classifications |
| **Run events** | registrations, results, awards, reservations, documents and printing |
| **Configure events** | create, configure, copy, close and delete events with programmes, classifications, kiosk and helper links and online booking |
| **Statistics and reports** | statistics, reports with amounts and ranking CSV files |
| **See shooters and clubs** | the personal data of the master data |
| **Edit shooters and clubs** | add and change shooters and clubs |
| **Manage master data** | import shooters, rename, merge and delete clubs, delete shooters |
| **Ranges and printers** | ranges, print agents and printer settings |
| **Print forms** | print forms and images (only while the club module is on) |
| **Users and club** | users, roles, invitations and club modules |

Every club starts with four roles:

| Role | Permissions |
|---|---|
| **Admin** | all; the column cannot be changed and the role cannot be deleted |
| **Board** | all but *Ranges and printers* and *Users and club* |
| **Helper** | see and run events, statistics, see and edit shooters and clubs |
| **Member** | see events and rankings |

**New role** adds a role of your own, which may only see the events until you tick more. The pencil renames a role (an empty name keeps the default name), the bin deletes it once nobody has it any more. The app shows each member only the pages and buttons their role allows: a member without *Run events* opens an event on its rankings.

Accounts are global: deleting an account is done by the platform operator in the [admin panel](../admin-guide/admin-panel.md). A forgotten password is reset with **Forgot your password?** on the login page if the account has a confirmed email address, otherwise by the operator.
