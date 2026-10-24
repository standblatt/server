# Permissions and roles

A member's role in a club is a set of permissions. Roles are edited under **Users** → **Roles** (see [Roles](../user-guide/users.md#roles)); API tokens carry a selection of the same permissions.

## Permissions {#permissions}

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

- One of the permissions an action names is enough.
- Without *Run events* an event opens on its rankings.
- The app shows each member only the pages and buttons their role allows.

## Default roles {#default-roles}

| Role | Permissions |
|---|---|
| **Admin** | all; the column cannot be changed and the role cannot be deleted |
| **Board** | all but *Ranges and printers* and *Users and club* |
| **Helper** | see and run events, statistics, see and edit shooters and clubs |
| **Member** | see events and rankings |

- **New role** starts with *See events and rankings* only.
- A role can be deleted once nobody has it; the admin role never.
- A club always keeps at least one member with the admin role.
- Members with *Users and club* but without the admin role give, change and remove only roles whose permissions they have themselves.

## Helper links {#helper-links}

Helper links are not roles of the club; each link has one of these roles (see [Creating a link](../user-guide/helpers.md#create)):

| Role | Sees | Can |
|---|---|---|
| **Scorer** | the lane grid of their own lanes, the shooters on them | enter results and awards, **Finish** a shooter |
| **Ammunition desk** | the search box, the waiting list, the shooters on all lanes | mark the ammunition as issued (**Ammo**), assign a lane to a waiting shooter, print the lane sheet |
| **Results clerk** | the search box, the shooters on all lanes, the completed list | enter results and awards on every lane |

Scorers and results clerks may also print [print forms](../user-guide/print-forms.md#use) of their lanes; a results clerk may mark awards as handed out.

## API token scopes {#scopes}

A token (**Account** → **API tokens**, see [API tokens](../user-guide/login-language.md#api-tokens)) has a list of clubs and of scopes:

| Scope | Key | Allows |
|---|---|---|
| *Read my training log* | `training:read` | read the own training log |
| *Log and change my training sessions* | `training:write` | add, change and delete the own sessions, disciplines, weapons and goals |
| a permission above | `events.view`, `events.operate`, `events.configure`, `reports.view`, `shooters.view`, `shooters.edit`, `shooters.manage`, `ranges.manage`, `forms.manage`, `club.manage` | the same as in a role, in the token's clubs |

- A club permission can only be chosen if the owner's role has it in one of the token's clubs.
- A token never has more than the owner's role: a permission the role loses, or a club the owner leaves, is gone for the token at once.
- A token without clubs reaches only the training log.

### Presets {#presets}

| Preset | Scopes |
|---|---|
| **Rankings for a website** | *See events and rankings* |
| **Training app** | *Read my training log*, *Log and change my training sessions* |
| **Read only** | *See events and rankings*, *Statistics and reports*, *See shooters and clubs*, *Read my training log* |

A preset only fills in the permissions the chosen clubs allow.

### Using a token {#token-use}

| What | Value |
|---|---|
| Format | `sb_` and 40 hexadecimal digits, shown once |
| Header | `Authorization: Bearer sb_…` or `X-API-Key: sb_…` |
| Expires | in 30 or 90 days, in a year, or never |
| Endpoints | `/api/tenants/…` (club data, by permission), `/api/my/…` (training log, by scope), `GET /api/token`, `/api/version`, `/api/auth/options`; with *Configure events* also `/api/modules` and `/api/programme-templates`; everything else answers 403 |
| `GET /api/token` | the token's name, scopes and clubs |
| Rate limit | 600 requests per token in five minutes (set by the server operator), then 429 |
| Unknown tokens | after 20 unknown or revoked tokens from one address in 15 minutes, 429 |
| API description | `/api/docs` |
| Never reachable | password, email, account settings, other tokens |
