# Modules

Standblatt has a **core** that every event has (general settings, programmes with series and prices, registration, results, rankings, statistics, closing) and **modules**: building blocks a club switches on per event, such as online booking, kiosk displays or attempts. A new feature that not every event needs is a module. It plugs into fixed extension points and leaves the rest of the code alone.

## The registry

`internal/modules` is the only list of modules: key, level (`event`, later `club`), the modules it requires, whether it needs a range and whether it is off by default for API clients. `modules.Effective` turns the chosen keys into the ones that are on: what was chosen, what holds data (see below) and what those require, without the modules that need a range when the event has none.

Events store the chosen keys in `events.modules`; `store.Event.Modules` is the effective list and `InUse` the modules that hold data.

## Backend extension points

| Point | Where | What a module adds |
|---|---|---|
| Key | `internal/modules/modules.go` | an entry in `All` and its key in `EventEnum` (a test compares them) plus the enum tags on `store.Event`/`EventInput` |
| Usage | `eventUsage` in `internal/store/modules.go` | SQL over the alias `e` that is true while the event holds data of the module, or `""`; a module with data stays on and switching it off answers 409 |
| Switch-off message | `moduleUsed` in `internal/web/modules.go` | a `case` with what to remove first |
| API gate | the module's `huma.Operation` | `Metadata: modules.Metadata(modules.Key)`; the middleware `moduleGate` answers 404 while the module is off. The path needs `eventId`, `programmeId`, `registrationId` or `entryId` (`TestModuleOperations`) |
| Documents and reports | `documentModules`, `reportModules` in `internal/web/modules.go` | the document or report key |
| Copy for next year | `copiers` in `internal/store/modules.go` | a function that copies the module's own data; programme fields are copied by `copyProgrammes` |

Handlers of a module need no checks of their own. Where the gate cannot see the event (a document type in the body), call `requireModule(ev, key)`.

## Frontend extension points

Every module has a folder `web/src/lib/modules/<key>/` with an `index.ts` that exports its `ModuleUI` (`$lib/modules/types.ts`): label, description, help page and icon, plus the extension points it uses. `$lib/modules/index.ts` collects them in `moduleUI`, typed `Record<ModuleKey, ModuleUI>`, so svelte-check fails when the API knows a module the UI lacks. Components are loaded lazily and get `{ tenantId, event }`; `ModuleSlot.svelte` renders one.

| Point | Rendered by | Use |
|---|---|---|
| `configPage` | `routes/…/events/[eventId]/config/[module]` with the module's label as title | its own page in the configuration, listed in the configuration sidebar |
| `operationPage` (+ `badge`) | `routes/…/events/[eventId]/[module]` | its own page in the operation sidebar, e.g. *Reservations*; renders its own `PageHeader` |
| `settings` | the *Modules* page, under the module's switch | a few fields saved with `updateEvent` from `$lib/tenant` |
| `rankingsSections`, `statisticsSections` | the rankings and statistics pages | panels after the core ones |
| `reports` | the statistics page's export list | report keys with their labels (the server checks `reportModules`) |

`has(event, key)` tells whether a module is on; check it where the feature lives, not in the page around it. Static routes win over `[module]`, so a module key must not be `entry`, `rankings`, `statistics`, `config`, `modules`, `programmes` or `close`.

## Adding a module

1. Add the key to `internal/modules` (and to the enum tags), with `Requires` and `NeedsRange` as needed.
2. Add its usage SQL and, if it can hold data, its message in `moduleUsed`.
3. Write its migration, store file and web file as usual; tag every operation with `Metadata`.
4. Add its documents, reports and copier to the lists above.
5. Tests: the module's own tests, plus a case in `TestModules` when it changes the switch-off rules.
6. UI: a folder `web/src/lib/modules/<key>/` with `index.ts` and its components, and the entry in `moduleUI`.
7. Help: a section or page in `web/src/help/en/user-guide/` and the module table in `events.md`.
