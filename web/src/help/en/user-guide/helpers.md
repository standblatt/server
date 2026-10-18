# Helpers

Helpers at the event (scorers, the ammunition desk, the results office) work with a **helper link** instead of a login. A helper link opens a reduced desk of one event on any phone or laptop; admins switch on the module *Helper links* under **Configuration** → **Modules**, create the links on **Helper links** in the configuration and revoke them when the event is over.

![Scorer's helper page on a phone](../images/screenshots/helper-scorer.png){ .screenshot .phone }

## Roles {#roles}

| Role | Sees | Can |
|---|---|---|
| **Scorer** | the lane grid of their own lanes, the shooters on them | enter results and awards, **Finish** a shooter |
| **Ammunition desk** | the search box, the waiting list, the shooters on all lanes | mark the ammunition as issued (**Ammo**), assign a lane to a waiting shooter, print the lane sheet |
| **Results clerk** | the search box, the shooters on all lanes, the completed list | enter results and awards on every lane |

Everything else stays with the club's users: registering shooters, editing programmes, rankings, statistics, deleting.

## Creating a link {#create}

1. Enter a name (who gets the link, e.g. *Scorer range A*), choose the role and, for a scorer, the range and the lanes (`1-5, 8`). One scorer usually looks after several lanes.
2. **Create helper link**; copy the link and send it to the helper (e.g. by messenger). It is shown only once.
3. The helper opens the link; it stays valid until it is revoked or the event is deleted. `?language=de` (`fr`, `it`, `en`) overrides the event's language.

**revoke** in the list stops a link immediately; *last seen* tells when it was last used. A scorer who opens a shooter on another lane gets *This lane is not yours*.

## Ammunition {#ammunition}

**Ammo** on a waiting or shooting registration records the time the ammunition was handed out; the button exists only when the module *Ammunition desk* is on (air pistol and air rifle shooters usually bring their own). The dashboard shows a green *ammo* badge on such registrations, so the desk and the range officer see at a glance who has been served.
