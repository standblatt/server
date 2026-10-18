# Modules

Standblatt has a **core** that every event has (general settings, programmes with series and prices, registration, results, rankings, statistics, closing) and **modules**: building blocks a club switches on per event, such as reservations, kiosk displays or attempts. A new feature that not every event needs is a module. It plugs into fixed extension points and leaves the rest of the code alone.

## The registry

`internal/modules` is the only list of modules: key, level (`event`, later `club`), the modules it requires, whether it needs a range and whether it is off by default for API clients. `modules.Effective` turns the chosen keys into the ones that are on: what was chosen, what holds data (see below) and what those require, without the modules that need a range when the event has none.

Events store the chosen keys in `events.modules`; `store.Event.Modules` is the effective list and `InUse` the modules that hold data.

## Backend extension points

| Point | Where | What a module adds |
|---|---|---|
| Key | `internal/modules/modules.go` | an entry in `All` plus its key in the enum tags on `store.Event`/`EventInput` (club modules: `store.Membership` and `set-club-modules`); `TestModuleEnums` compares every module enum of the API with `All` |
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
| `programmeSections` | the programme page (`ProgrammeEditor`), `after` one of `basics`, `series`, `prices`, `usage` | sections that edit the `ProgrammeDraft` (`$lib/programme/draft.svelte.ts`) saved with the programme, or, after `usage`, save at once through their own endpoint (only for saved programmes); props `ProgrammeProps` |
| `reports` | the statistics page's export list | report keys with their labels (the server checks `reportModules`) |

`has(event, key)` tells whether a module is on; check it where the feature lives, not in the page around it. Static routes win over `[module]`, so a module key must not be `entry`, `rankings`, `statistics`, `config`, `modules`, `programmes` or `close`.

## Club modules

Modules with `Level: Club` belong to a club: `tenants.modules` holds the chosen keys, `clubUsage` in `internal/store/modules.go` tells when a club holds data, `store.Membership` (in `/api/me`) carries `modules` and `modulesInUse`, and `PUT /api/tenants/{tenantId}/modules` chooses them. The gate checks the club of `tenantId` for operations tagged with a club module. In the UI, `clubModuleUI` in `$lib/modules/index.ts` (`ClubModuleUI`: label, description, help, icon, `clubPage` = a page of the club in the sidebar with the permission it needs, `programmeSections`) and `hasClub(membership, key)`; the *Club modules* page at `/t/<id>/modules` uses the same `ModuleCards` as the event's *Modules* page.

## Permissions

What a member may do comes from the permissions of their role, not from the module. `internal/permissions` is the only list (`All` with the key and, for a permission that belongs to a club module, its `Module`; `Enum` and the enum tags on `store.Role`/`RoleInput`/`Membership` must match, `TestModuleEnums` checks the spec, `TestPermissionsUsed` that every permission is required somewhere). Every operation under `/api/tenants/{tenantId}` declares its permissions in its metadata, `Metadata: needs(moduleKey, permissions.X, …)` (or `withNeeds(op, …)`; one of them is enough); the huma middleware `permissionGate` refuses the request before the input is parsed (404 for non-members, 403 without a permission), and handlers call `s.access(ctx, tenantID)` only for the user. `TestPermissionTable` fails on an operation without permissions and compares the whole table with `internal/web/testdata/permissions.txt` (after a deliberate change: `go test ./internal/web -run TestPermissionTable -args -update-permissions`, then review the diff); `TestPermissionLeaks` calls every operation as a member without its permissions and with every other one (403) and with each of its permissions (no 403), and as every helper link role outside `helperOps` (403). A module that needs a right of its own adds a key there with its `Module`, names it in `needs` and gives it a text in `permissionUI` (`$lib/permissions.ts`, a `Record<Permission, …>`, so svelte-check catches a missing one); the role matrix shows it only while the module is on. The UI asks `can(membership, permission)` from `$lib/me.ts`; pages that need a permission are wrapped in `$lib/RequirePermission.svelte`.

## Adding a module

1. Add the key to `internal/modules` (and to the enum tags), with `Requires` and `NeedsRange` as needed.
2. Add its usage SQL and, if it can hold data, its message in `moduleUsed`.
3. Write its migration, store file and web file as usual; tag every operation with `Metadata`.
4. Add its documents, reports and copier to the lists above.
5. Tests: the module's own tests, plus a case in `TestModules` when it changes the switch-off rules.
6. UI: a folder `web/src/lib/modules/<key>/` with `index.ts` and its components, and the entry in `moduleUI`.
7. Help: a section or page in `web/src/help/en/user-guide/` and the module table in `events.md`.
