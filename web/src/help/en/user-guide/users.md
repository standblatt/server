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

Every permission and what it allows: [Permissions](../reference/permissions.md#permissions). Every club starts with the roles **Admin**, **Board**, **Helper** and **Member** (see [Default roles](../reference/permissions.md#default-roles)).

**New role** adds a role of your own, which may only see the events until you tick more. The pencil renames a role (an empty name keeps the default name), the bin deletes it once nobody has it any more. The app shows each member only the pages and buttons their role allows: a member without *Run events* opens an event on its rankings.

Members with *Users and club* but without the admin role stay within their own permissions: they give, change, remove and invite only with roles whose permissions they have themselves, and make or change roles only within them. So only admins make, change or remove admins.

Accounts are global (see [One account, a role per club](../explanation/data.md#accounts)): deleting an account is done by the platform operator in the admin panel. For a forgotten password see [Forgotten password](login-language.md#password).
