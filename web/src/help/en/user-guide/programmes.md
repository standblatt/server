# Programmes

Programmes belong to an event and are managed on the *Programmes* tab of its **Setup** (admins only).

![Setup with the programme list](../images/screenshots/setup.png){ .screenshot }

## Adding programmes

- **New programme** opens the editor.
- **Add from … → Programme library**: a template of the global library, maintained by the platform operator.
- **Add from … → Earlier event**: copies all programmes of another event of the club, including combinations and special prizes. Programmes whose name and variant already exist are skipped.
- **Add from … → Programme file**: imports a [programme file](../reference/programme-file.md) exported from another event or club.

## The editor

![Programme editor](../images/screenshots/programme-editor.png){ .screenshot }

- **Name** and **Variant**: programmes with the same name are variants of each other; a shooter registers for one of them. The name field suggests the names of the event's programmes, and as soon as the name matches one, the editor says that you are creating a variant and warns when the variant name is empty or already taken.
- **Barcode**: printed on the lane sheet for electronic scoring (see [Codes](../reference/codes.md)).
- **Counts for the club competition** and **Show on kiosk displays**.
- **Series**: shots, maximum points and whether mouches are tracked, per series.
- **Prices (CHF)**: only the ticked age categories may shoot the programme.
- **Awards**: name, kind (wreath card, prize in kind or cash prize paid out in money), value and the minimum total per age category; leave a category empty where the award is not available.

Once shooters are registered for a programme, the **series** are fixed so results keep their meaning; the editor says so and greys them out. Everything else can still be changed: name, variant, description, barcode, the flags, prices, series names, award names, values and limits. Age categories and awards can be **added** but not removed. Also available on a programme in use:

- **Active**: inactive programmes cannot be chosen for new registrations; existing ones keep them.
- **Combinable with**: the programme whose series this one may take at registration.
- **Range** (events on several ranges): the range the programme is shot on, e.g. *Pistol 50 m*; the registration then offers only its lanes and kiosk displays group the rankings by range. *any* leaves it open. Programmes of different ranges need separate registrations.
- **Special prizes** (see below).

**Export** downloads the programme as a file; **Delete** is only possible while nobody is registered.

## Special prizes {#special-prizes}

**Special prizes** in the programme menu awards the best ranks of every age category, for example ranks 1–3 in every category with at least 10 participants, counting only shooters whose club has at least 4 participants in this programme. Skipped shooters make room for the next one. The exact rule is in [Competition rules](../reference/rules.md#special-prizes); the winners are in the report *Special prizes* under [Statistics](statistics-reports.md).

![Special prizes dialog](../images/screenshots/special-prizes.png){ .screenshot }

Ranks 0 turns special prizes off. The settings can be changed at any time, also after the event.
