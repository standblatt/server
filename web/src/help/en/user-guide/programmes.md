# Programmes

Programmes belong to an event and are managed under **Configuration** → **Programmes** (admins only).

![The programme list](../images/screenshots/setup.png){ .screenshot }

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
- **Series**: shots, maximum points, whether mouches are tracked and whether **Single shots** are recorded, per series. With single shots, every shot can be entered and the series is their sum, e.g. to print them onto a form; the series total alone can still be entered.
- **Attempts** (see below): how often a shooter may shoot the programme at the event.
- **Prices (CHF)**: only the ticked age categories may shoot the programme. With several attempts, a second field per category is the price of every further attempt (*Nachdoppel*); empty means the same price, 0 means the programme is paid once.
- **Awards**: name, kind (wreath card, prize in kind or cash prize paid out in money), value and the minimum total per age category; leave a category empty where the award is not available.

Once shooters are registered for a programme, the **series** are fixed so results keep their meaning; the editor says so and greys them out. Everything else can still be changed: name, variant, description, barcode, the flags, prices, series names, award names, values and limits. Age categories and awards can be **added** but not removed. Also available on a programme in use:

- **Active**: inactive programmes cannot be chosen for new registrations; existing ones keep them.
- **Combinable with**: the programme whose series this one may take at registration.
- **Range** (events on several ranges): the range the programme is shot on, e.g. *Pistol 50 m*; the registration then offers only its lanes and kiosk displays group the rankings by range. *any* leaves it open. Programmes of different ranges need separate registrations.
- **Print form** (when the club has [print forms](print-forms.md)): the form the results of the programme are printed onto; *none* prints no form.
- **Special prizes** (see below).

**Export** downloads the programme as a file; **Delete** is only possible while nobody is registered.

## Attempts {#attempts}

By default a shooter shoots a programme once per event. **Attempts per shooter** allows more, e.g. 3 for a lucky target with further attempts at a shooting festival, or **Unlimited** for an annual programme shot over the season.

- **Counting attempt** decides which attempt the ranking, the club competition, special prizes and *My training* count: **The best** (by the ranking rules), **The first** or **The last**. The ranking shows it as *2/3*: the second of three attempts counts.
- **Only the counting attempt receives an award**: off, every attempt can earn its award (a wreath card per further attempt); on, only the counting attempt gets one, and the award moves when a better attempt is shot.
- At registration, a programme with several attempts shows a counter: two attempts at once are two result cards, *Attempt 1* and *Attempt 2*. Further attempts can also be added later in a new registration, or by entering another result on **Enter results**, until the shooter has no attempts left.

The number of attempts can be changed while shooters are registered, but not below the attempts a shooter already has. Programmes with several attempts cannot be combined.

## Special prizes {#special-prizes}

**Special prizes** in the programme menu awards the best ranks of every age category, for example ranks 1–3 in every category with at least 10 participants, counting only shooters whose club has at least 4 participants in this programme. Skipped shooters make room for the next one. The exact rule is in [Competition rules](../reference/rules.md#special-prizes); the winners are in the report *Special prizes* under [Statistics](statistics-reports.md).

![Special prizes dialog](../images/screenshots/special-prizes.png){ .screenshot }

Ranks 0 turns special prizes off. The settings can be changed at any time, also after the event.
