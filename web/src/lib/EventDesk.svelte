<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Flag from '@lucide/svelte/icons/flag';
	import ListChecks from '@lucide/svelte/icons/list-checks';
	import Package from '@lucide/svelte/icons/package-check';
	import Pencil from '@lucide/svelte/icons/pencil';
	import ScanLine from '@lucide/svelte/icons/scan-line';
	import UserPlus from '@lucide/svelte/icons/user-plus';
	import { toast } from 'svelte-sonner';
	import { api, call, type Registration } from './api';
	import { ageCategoryLabel } from './labels';
	import { intlLocale } from './locale.svelte';
	import { Badge } from './components/ui/badge';
	import { Button } from './components/ui/button';
	import { Kbd } from './components/ui/kbd';
	import * as Table from './components/ui/table';
	import { ask } from './confirm.svelte';
	import type { DeskRights } from './desk';
	import { pad2, pageOf, programmeTitle, withRange } from './format';
	import Pager from './Pager.svelte';
	import NewShooterDialog from './NewShooterDialog.svelte';
	import PrintButton from './PrintButton.svelte';
	import RegistrationDialog from './RegistrationDialog.svelte';
	import ResultsDialog from './ResultsDialog.svelte';
	import { eventQuery, keys, lanesQuery, programmesQuery, registrationsQuery } from './tenant';

	let {
		tenantId,
		eventId,
		base,
		can,
		only
	}: {
		tenantId: number;
		eventId: number;
		/** Route the dialogs hang off with query parameters (?results=, ?registration=, ?register=). */
		base: string;
		can: DeskRights;
		/** Scorers: show only these lanes of this range. */
		only?: { rangeId: number; lanes: number[] };
	} = $props();

	const queryClient = useQueryClient();
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const archived = $derived(!!event.data?.archivedAt);
	const ammoIssue = $derived(!!event.data?.modules.includes('ammo'));
	/** The event assigns lanes; without, registrations stay open until their results are complete. */
	const lanesOn = $derived(!!event.data?.modules.includes('lanes'));
	const laneSheet = $derived(!!event.data?.modules.includes('lane_sheet'));
	const registrations = registrationsQuery(
		() => tenantId,
		() => eventId
	);
	const lanes = lanesQuery(
		() => tenantId,
		() => eventId
	);
	const programmes = programmesQuery(
		() => tenantId,
		() => eventId
	);
	const stats = createQuery(() => ({
		queryKey: keys.stats(tenantId, eventId),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/stats', { params: { path: { tenantId, eventId } } })),
		enabled: !only
	}));

	/** Whether a registration is on one of the helper's lanes (always true without a lane restriction). */
	const mine = (r: { rangeId?: number; lane?: number }) => !only || (r.rangeId === only.rangeId && !!r.lane && only.lanes.includes(r.lane));
	const waiting = $derived((registrations.data ?? []).filter((r) => r.status === 'registered'));
	const shooting = $derived((registrations.data ?? []).filter((r) => r.status === 'shooting' && mine(r)));
	const completed = $derived(
		(registrations.data ?? [])
			.filter((r) => r.status === 'completed')
			.sort((a, b) => (b.completedAt ?? '').localeCompare(a.completedAt ?? ''))
	);
	const pageSize = 50;
	let waitingPage = $state(0);
	let completedPage = $state(0);
	const shownLanes = $derived(
		(lanes.data ?? [])
			.filter((rl) => !only || rl.rangeId === only.rangeId)
			.map((rl) => (only ? { ...rl, lanes: rl.lanes.filter((l) => only.lanes.includes(l.lane)) } : rl))
	);
	const activeByShooter = $derived(new Map((registrations.data ?? []).filter((r) => r.status !== 'completed').map((r) => [r.shooter.id, r.id])));
	/** Whether any entry of the registration still lacks results. */
	const missingResults = (r: Registration) => r.entries.some((e) => e.status !== 'completed');
	/** Whether any entry of the registration has an award. */
	const hasAwards = (r: Registration) => r.entries.some((e) => e.awardTierId);
	const shootingColumns = $derived(2 + (can.results || can.assign || can.finish || can.ammo ? 1 : 0) + 2);

	let now = $state(Date.now());
	$effect(() => {
		const t = setInterval(() => (now = Date.now()), 1000);
		return () => clearInterval(t);
	});
	/** Time since a timestamp in minutes, or hours and minutes from an hour on. */
	function since(iso?: string) {
		if (!iso) return '';
		const min = Math.max(0, Math.floor((now - new Date(iso).getTime()) / 60000));
		return min < 60 ? $t`${min} min` : $t`${Math.floor(min / 60)} h ${min % 60} min`;
	}

	/** Waiting time as m:ss (h:mm:ss from an hour on): at the desk, seconds show that the list moves. */
	function waitingFor(iso: string) {
		const s = Math.max(0, Math.floor((now - new Date(iso).getTime()) / 1000));
		return s >= 3600 ? `${Math.floor(s / 3600)}:${pad2(Math.floor((s % 3600) / 60))}:${pad2(s % 60)}` : `${Math.floor(s / 60)}:${pad2(s % 60)}`;
	}

	/** Comma-separated programme names (with variant) of a registration's entries. */
	function programmeNames(r: Registration) {
		return r.entries
			.map((e) => programmes.data?.find((p) => p.id === e.programmeId))
			.map((p) => (p ? programmeTitle(p) : '?'))
			.join(', ');
	}

	let searchInput = $state<HTMLInputElement>();
	let input = $state('');
	let q = $state('');
	let scanError = $state('');
	let timer: ReturnType<typeof setTimeout>;
	const results = createQuery(() => ({
		queryKey: [...keys.shooters(tenantId), q, 0],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/shooters', { params: { path: { tenantId }, query: { q, limit: 20 } } })),
		enabled: can.search && q !== '' && !q.startsWith('REGCODE')
	}));

	/** Debounces the search input into the shooter query. */
	function onInput() {
		scanError = '';
		clearTimeout(timer);
		timer = setTimeout(() => (q = input.trim()), 250);
	}

	/** Escape clears the search; Enter on a scanned registration code checks event and license and opens the registration. */
	async function onKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			input = q = scanError = '';
		}
		if (e.key !== 'Enter') return;
		const m = input.trim().match(/^REGCODEAAAA(\d+)AAAAL(\d+)$/);
		if (!m) return;
		e.preventDefault();
		try {
			const r = await call(api.GET('/api/tenants/{tenantId}/registrations/{registrationId}', { params: { path: { tenantId, registrationId: Number(m[1]) } } }));
			const license = m[2] === '0' ? '' : m[2];
			if (r.eventId !== eventId) {
				scanError = $t`The code belongs to another event`;
				return;
			}
			if (license && r.shooter.licenseNumber && r.shooter.licenseNumber !== license) {
				scanError = $t`The code does not match the registration`;
				return;
			}
			input = q = '';
			open(r.id);
		} catch (err) {
			scanError = (err as Error).message;
		}
	}

	/** F4 clears and focuses the shooter search. */
	function onGlobalKey(e: KeyboardEvent) {
		if (e.key === 'F4' && can.search) {
			e.preventDefault();
			input = q = '';
			searchInput?.focus();
		}
	}

	/** Lane label (with range name when there are several) or the waiting list. */
	const whereOf = (r: Registration) => (r.lane ? withRange(event.data?.ranges ?? [], r.rangeId, $t`Lane ${r.lane}`) : $t`the waiting list`);

	/** Asks for confirmation, then finishes the registration and frees its lane. */
	async function finish(r: Registration) {
		const name = `${r.shooter.firstname} ${r.shooter.lastname}`;
		const where = whereOf(r);
		const description = r.lane || lanesOn ? $t`${name} leaves ${where}.` : $t`${name} is marked as finished.`;
		if (!(await ask({ title: $t`Finish shooting?`, description, action: $t`Finish` }))) return;
		await call(api.POST('/api/tenants/{tenantId}/registrations/{registrationId}/finish', { params: { path: { tenantId, registrationId: r.id } } }));
		await queryClient.invalidateQueries({ queryKey: keys.registrations(tenantId, eventId) });
		toast.success($t`${name} finished`);
	}

	/** Asks for confirmation, then takes the ammunition issue or the awards handout back. */
	async function revoke(r: Registration, what: 'ammo' | 'awards') {
		const name = `${r.shooter.firstname} ${r.shooter.lastname}`;
		const ok = await ask(
			what === 'ammo'
				? { title: $t`Take the ammunition issue back?`, description: $t`${name} counts as not served again.`, action: $t`Take back` }
				: { title: $t`Take the awards handout back?`, description: $t`The awards of ${name} count as open again.`, action: $t`Take back` }
		);
		if (!ok) return;
		if (what === 'ammo') await call(api.DELETE('/api/tenants/{tenantId}/registrations/{registrationId}/ammo', { params: { path: { tenantId, registrationId: r.id } } }));
		else await call(api.DELETE('/api/tenants/{tenantId}/registrations/{registrationId}/awards-issued', { params: { path: { tenantId, registrationId: r.id } } }));
		await queryClient.invalidateQueries({ queryKey: keys.registrations(tenantId, eventId) });
	}

	/** Marks the ammunition as issued to the shooter. */
	async function issueAmmo(r: Registration) {
		await call(api.POST('/api/tenants/{tenantId}/registrations/{registrationId}/ammo', { params: { path: { tenantId, registrationId: r.id } } }));
		await queryClient.invalidateQueries({ queryKey: keys.registrations(tenantId, eventId) });
		toast.success($t`Ammunition issued to ${r.shooter.firstname} ${r.shooter.lastname}`);
	}

	const params = $derived(page.url.searchParams);
	const dialog = $derived.by(() => {
		/** Numeric query parameter, undefined when missing or 0. */
		const n = (k: string) => Number(params.get(k)) || undefined;
		if (n('results')) return { kind: 'results' as const, id: n('results')! };
		if (n('registration')) return { kind: 'edit' as const, id: n('registration')! };
		if (n('register')) return { kind: 'new' as const, shooterId: n('register')!, clubId: n('club') };
		return null;
	});
	/** Navigates to the desk with the given query string, keeping scroll and focus. */
	const show = (query: string) => goto(`${base}${query}`, { noScroll: true, keepFocus: true });
	/** Opens the results dialog of a registration through the URL query parameter. */
	const openResults = (id: number) => show(`?results=${id}`);
	/** Opens the registration dialog through the URL query parameter. */
	const openRegistration = (id: number) => show(`?registration=${id}`);
	/** Opens what the viewer may work on: results, else the registration. */
	const open = (id: number) => (can.results ? openResults(id) : can.assign ? openRegistration(id) : undefined);
	/** Clears the search and opens the registration dialog for a shooter and optional club. */
	const register = (shooterId: number, clubId?: number) => {
		input = q = '';
		show(`?register=${shooterId}${clubId ? `&club=${clubId}` : ''}`);
	};
	/** Closes any dialog by dropping the query parameters. */
	const close = () => goto(base, { noScroll: true, keepFocus: true, replaceState: true });
	let newShooter = $state(false);
</script>

<svelte:window onkeydown={onGlobalKey} />

{#snippet shooter(r: Registration)}
	<div class="font-medium">{r.shooter.lastname} {r.shooter.firstname}</div>
	<div class="text-xs text-muted-foreground">{r.clubName}</div>
{/snippet}

{#snippet panel(title: string, count: number)}
	<div class="flex items-center gap-2 border-b px-4 py-2.5">
		<h2 class="text-sm font-semibold">{title} ({count})</h2>
	</div>
{/snippet}

{#snippet awards(r: Registration)}
	{#if hasAwards(r)}
		{#if r.awardsIssuedAt}
			{@const at = new Date(r.awardsIssuedAt).toLocaleTimeString(intlLocale(), { hour: '2-digit', minute: '2-digit' })}
			{#if can.revoke && !archived}
				<button type="button" class="cursor-pointer" title={$t`${at} · click to take back`} onclick={() => revoke(r, 'awards')}>
					<Badge variant="outline" class="border-emerald-500/50 text-emerald-700 hover:bg-emerald-500/10 dark:text-emerald-400">{$t`awards handed out`}</Badge>
				</button>
			{:else}
				<Badge variant="outline" class="border-emerald-500/50 text-emerald-700 dark:text-emerald-400" title={at}>{$t`awards handed out`}</Badge>
			{/if}
		{:else if can.results}
			<button type="button" class="cursor-pointer" title={$t`Open the results to see and hand out the awards`} onclick={() => openResults(r.id)}>
				<Badge variant="outline" class="border-amber-500/50 text-amber-700 hover:bg-amber-500/10 dark:text-amber-400">{$t`awards to hand out`}</Badge>
			</button>
		{:else}
			<Badge variant="outline" class="border-amber-500/50 text-amber-700 dark:text-amber-400">{$t`awards to hand out`}</Badge>
		{/if}
	{/if}
{/snippet}

{#snippet ammo(r: Registration)}
	{#if r.ammoIssuedAt}
		{@const at = new Date(r.ammoIssuedAt).toLocaleTimeString(intlLocale(), { hour: '2-digit', minute: '2-digit' })}
		{#if can.revoke && !archived}
			<button type="button" class="cursor-pointer" title={$t`${at} · click to take back`} onclick={() => revoke(r, 'ammo')}>
				<Badge variant="outline" class="border-emerald-500/50 text-emerald-700 hover:bg-emerald-500/10 dark:text-emerald-400">{$t`ammo`}</Badge>
			</button>
		{:else}
			<Badge variant="outline" class="border-emerald-500/50 text-emerald-700 dark:text-emerald-400" title={at}>{$t`ammo`}</Badge>
		{/if}
	{:else if can.ammo && ammoIssue && !archived}
		<Button size="sm" variant="outline" title={$t`Ammunition issued`} onclick={() => issueAmmo(r)}><Package /><span class="sr-only xl:not-sr-only">{$t`Ammo`}</span></Button>
	{/if}
{/snippet}

<div class="grid flex-1 gap-4 *:min-w-0 xl:grid-cols-[minmax(0,1fr)_20rem]">
	<div class="space-y-4">
		{#if can.search}
			<div class="relative flex flex-wrap items-center gap-2">
				<div class="relative min-w-64 flex-1">
					<ScanLine class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
					<input
						bind:this={searchInput}
						class="h-11 w-full rounded-md border border-input bg-card pr-12 pl-9 text-base shadow-xs outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/40 dark:bg-input/30"
						type="search"
						placeholder={$t`Name or license number, or scan the code`}
						aria-label={$t`Search shooter`}
						bind:value={input}
						oninput={onInput}
						onkeydown={onKeydown}
					/>
					<Kbd class="pointer-events-none absolute top-1/2 right-3 -translate-y-1/2">F4</Kbd>
					{#if (results.data && q) || scanError}
						<div class="absolute inset-x-0 top-full z-30 mt-1 overflow-hidden rounded-md border bg-popover shadow-lg">
							{#if scanError}<p class="px-4 py-3 text-sm text-destructive" role="alert">{scanError}</p>{/if}
							{#if results.data && q}
								<ul class="max-h-96 divide-y overflow-y-auto">
									{#each results.data as sh (sh.id)}
										<li class="flex flex-wrap items-center gap-3 px-4 py-2">
											<span class="w-16 font-mono text-xs text-muted-foreground">{sh.licenseNumber || '–'}</span>
											<span class="font-medium">{sh.lastname} {sh.firstname}</span>
											<span class="ml-auto flex flex-wrap gap-2">
												{#if activeByShooter.has(sh.id)}
													{#if can.assign}
														<Button size="sm" variant="outline" onclick={() => openRegistration(activeByShooter.get(sh.id)!)}>{$t`Open registration`}</Button>
													{:else if can.results}
														<Button size="sm" variant="outline" onclick={() => openResults(activeByShooter.get(sh.id)!)}>{$t`Results`}</Button>
													{/if}
												{/if}
												{#if can.register && !archived}
													{#each sh.clubs as c (c.id)}
														<Button size="sm" onclick={() => register(sh.id, c.id)}>{$t`Register for ${c.name}`}</Button>
													{:else}
														<Button size="sm" onclick={() => register(sh.id)}>{$t`Register`}</Button>
													{/each}
												{/if}
											</span>
										</li>
									{:else}
										<li class="flex items-center gap-3 px-4 py-3 text-sm text-muted-foreground">
											{$t`No shooter found.`}
											{#if can.register && !archived}<Button size="sm" variant="outline" class="ml-auto" onclick={() => (newShooter = true)}><UserPlus />{$t`New shooter`}</Button>{/if}
										</li>
									{/each}
								</ul>
							{/if}
						</div>
					{/if}
				</div>
				{#if can.register && !archived}<Button variant="outline" class="h-11" onclick={() => (newShooter = true)}><UserPlus />{$t`New shooter`}</Button>{/if}
			</div>
		{/if}

		{#if lanesOn}
		<section class="overflow-hidden rounded-lg border bg-card">
			{@render panel($t`Shooting`, shooting.length)}
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="w-14">{$t`Lane`}</Table.Head>
						<Table.Head>{$t`Shooter`}</Table.Head>
						<Table.Head class="hidden 2xl:table-cell">{$t`Category`}</Table.Head>
						<Table.Head class="hidden lg:table-cell">{$t`Programmes`}</Table.Head>
						<Table.Head></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each shooting as r (r.id)}
						<Table.Row class={can.results || can.assign ? 'cursor-pointer' : ''} onclick={() => open(r.id)}>
							<Table.Cell class="tabular-nums">
								<div class="text-lg font-semibold">{r.lane}</div>
								<div class="text-xs whitespace-nowrap text-muted-foreground">{since(r.startedAt)}</div>
							</Table.Cell>
							<Table.Cell class="whitespace-normal">{@render shooter(r)}</Table.Cell>
							<Table.Cell class="hidden 2xl:table-cell">{$t(ageCategoryLabel(r.ageCategory))}</Table.Cell>
							<Table.Cell class="hidden whitespace-normal lg:table-cell">{programmeNames(r)}</Table.Cell>
							<Table.Cell onclick={(e: MouseEvent) => e.stopPropagation()}>
								<div class="flex items-center justify-end gap-1.5">
									{@render ammo(r)}
									{#if can.results}<Button size="sm" href="{base}?results={r.id}" data-sveltekit-noscroll data-sveltekit-keepfocus title={$t`Results`}><ListChecks /><span class="sr-only sm:not-sr-only">{$t`Results`}</span></Button>{/if}
									{#if can.assign}<Button size="sm" variant="outline" title={$t`Edit registration`} onclick={() => openRegistration(r.id)}><Pencil /><span class="sr-only xl:not-sr-only">{$t`Edit`}</span></Button>{/if}
									{#if can.finish && !archived}<Button size="sm" variant="outline" title={$t`Finish shooting`} onclick={() => finish(r)}><Flag /><span class="sr-only xl:not-sr-only">{$t`Finish`}</span></Button>{/if}
								</div>
							</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row><Table.Cell colspan={shootingColumns} class="py-6 text-center text-muted-foreground">{$t`Nobody is shooting.`}</Table.Cell></Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</section>
		{/if}

		{#if can.waiting || (!lanesOn && can.results)}
			<section class="overflow-hidden rounded-lg border bg-card">
				{@render panel(lanesOn ? $t`Waiting` : $t`Registered`, waiting.length)}
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{$t`Shooter`}</Table.Head>
							<Table.Head class="hidden md:table-cell">{$t`Category`}</Table.Head>
							<Table.Head class="hidden lg:table-cell">{$t`Programmes`}</Table.Head>
							<Table.Head class="hidden sm:table-cell">{$t`Waiting`}</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each pageOf(waiting, waitingPage, pageSize) as r (r.id)}
							<Table.Row
								class={lanesOn ? (can.assign ? 'cursor-pointer' : '') : can.results || can.register ? 'cursor-pointer' : ''}
								onclick={() => (lanesOn ? can.assign && openRegistration(r.id) : can.results ? openResults(r.id) : can.register && openRegistration(r.id))}
							>
								<Table.Cell class="whitespace-normal">{@render shooter(r)}</Table.Cell>
								<Table.Cell class="hidden md:table-cell">{$t(ageCategoryLabel(r.ageCategory))}</Table.Cell>
								<Table.Cell class="hidden whitespace-normal lg:table-cell">{programmeNames(r)}</Table.Cell>
								<Table.Cell class="hidden text-muted-foreground tabular-nums sm:table-cell">{waitingFor(r.registeredAt)}</Table.Cell>
								<Table.Cell onclick={(e: MouseEvent) => e.stopPropagation()}>
									<div class="flex items-center justify-end gap-1.5">
										{@render ammo(r)}
										{#if laneSheet && can.print && !can.register}<PrintButton {tenantId} target={{ document: 'lane_sheet', registrationId: r.id }} />{/if}
										{#if lanesOn}
											{#if can.assign}<Button size="sm" onclick={() => openRegistration(r.id)}>{$t`Assign lane`}</Button>{/if}
										{:else}
											{#if can.results}<Button size="sm" href="{base}?results={r.id}" data-sveltekit-noscroll data-sveltekit-keepfocus title={$t`Results`}><ListChecks /><span class="sr-only sm:not-sr-only">{$t`Results`}</span></Button>{/if}
											{#if can.register}<Button size="sm" variant="outline" title={$t`Edit registration`} onclick={() => openRegistration(r.id)}><Pencil /><span class="sr-only xl:not-sr-only">{$t`Edit`}</span></Button>{/if}
											{#if can.finish && !archived}<Button size="sm" variant="outline" title={$t`Finish shooting`} onclick={() => finish(r)}><Flag /><span class="sr-only xl:not-sr-only">{$t`Finish`}</span></Button>{/if}
										{/if}
									</div>
								</Table.Cell>
							</Table.Row>
						{:else}
							<Table.Row><Table.Cell colspan={5} class="py-6 text-center text-muted-foreground">{lanesOn ? $t`Nobody is waiting.` : $t`Nobody is registered.`}</Table.Cell></Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
				<Pager bind:page={waitingPage} total={waiting.length} size={pageSize} />
			</section>
		{/if}

		{#if can.completed}
			<section class="overflow-hidden rounded-lg border bg-card">
				{@render panel($t`Completed`, completed.length)}
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{$t`Shooter`}</Table.Head>
							<Table.Head class="hidden md:table-cell">{$t`Category`}</Table.Head>
							<Table.Head class="hidden lg:table-cell">{$t`Programmes`}</Table.Head>
							<Table.Head>{$t`Results`}</Table.Head>
							<Table.Head>{$t`Awards`}</Table.Head>
							<Table.Head class="hidden sm:table-cell"></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each pageOf(completed, completedPage, pageSize) as r (r.id)}
							<Table.Row class={can.results ? 'cursor-pointer' : ''} onclick={() => can.results && openResults(r.id)}>
								<Table.Cell class="whitespace-normal">{@render shooter(r)}</Table.Cell>
								<Table.Cell class="hidden md:table-cell">{$t(ageCategoryLabel(r.ageCategory))}</Table.Cell>
								<Table.Cell class="hidden whitespace-normal lg:table-cell">{programmeNames(r)}</Table.Cell>
								<Table.Cell>
									{#if missingResults(r)}<Badge variant="outline" class="border-amber-500/50 text-amber-700 dark:text-amber-400">{$t`missing`}</Badge>{:else}<Badge variant="outline" class="border-emerald-500/50 text-emerald-700 dark:text-emerald-400">{$t`complete`}</Badge>{/if}
								</Table.Cell>
								<Table.Cell onclick={(e: MouseEvent) => e.stopPropagation()}><div class="flex flex-wrap items-center gap-1.5">{@render awards(r)}</div></Table.Cell>
								<Table.Cell class="hidden text-right sm:table-cell">{#if can.results}<Button size="sm" variant="outline" href="{base}?results={r.id}" data-sveltekit-noscroll data-sveltekit-keepfocus>{$t`Results`}</Button>{/if}</Table.Cell>
							</Table.Row>
						{:else}
							<Table.Row><Table.Cell colspan={6} class="py-6 text-center text-muted-foreground">{$t`Nobody has finished yet.`}</Table.Cell></Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
				<Pager bind:page={completedPage} total={completed.length} size={pageSize} />
			</section>
		{/if}
	</div>

	<aside class="space-y-4">
		{#if stats.data}
			<section class="grid {lanesOn ? 'grid-cols-3' : 'grid-cols-2'} divide-x overflow-hidden rounded-lg border bg-card" aria-label={$t`Statistics`}>
				{#each lanesOn ? [[stats.data.registered, $t`waiting`], [stats.data.shooting, $t`shooting`], [stats.data.completed, $t`completed`]] : [[stats.data.registered, $t`open`], [stats.data.completed, $t`completed`]] as [n, label] (label)}
					<div class="px-3 py-3 text-center"><div class="text-2xl font-semibold tabular-nums">{n}</div><div class="text-xs text-muted-foreground">{label}</div></div>
				{/each}
			</section>
		{/if}
		{#each shownLanes as rl (rl.rangeId)}
			<section class="rounded-lg border bg-card">
				<div class="flex items-center border-b px-4 py-2.5"><h2 class="text-sm font-semibold">{(lanes.data?.length ?? 0) > 1 ? rl.rangeName : $t`Lanes`}</h2></div>
				<div class="grid grid-cols-4 gap-1.5 p-3 sm:grid-cols-6 xl:grid-cols-4">
					{#each rl.lanes as l (l.lane)}
						{#if l.state === 'occupied'}
							<button class="lane lane-occupied text-left hover:border-sky-400" onclick={() => open(l.registrationId!)} aria-label={$t`Lane ${l.lane}: ${l.shooterName ?? ''}`}>
								<span class="text-base font-bold">{l.lane}</span>
								<span class="truncate">{l.shooterName}</span>
								<span class="opacity-70">{since(l.since)}</span>
							</button>
						{:else if l.state === 'reserved'}
							<div class="lane lane-reserved" aria-label={$t`Lane ${l.lane}: reserved`} title={l.reservation}>
								<span class="text-base font-bold">{l.lane}</span>
								<span class="truncate">{l.reservation}</span>
								<span class="opacity-70">{$t`until ${l.reservedUntil ? new Date(l.reservedUntil).toLocaleTimeString(intlLocale(), { hour: '2-digit', minute: '2-digit' }) : ''}`}</span>
							</div>
						{:else}
							<div class="lane {l.state === 'disabled' ? 'lane-disabled' : 'lane-free'}" aria-label={$t`Lane ${l.lane}: ${l.state === 'disabled' ? $t`out of service` : $t`free`}`}>
								<span class="text-base font-bold">{l.lane}</span>
								<span>{l.state === 'disabled' ? $t`out of service` : $t`free`}</span>
							</div>
						{/if}
					{/each}
				</div>
			</section>
		{/each}
	</aside>
</div>

{#if can.register}
	<NewShooterDialog {tenantId} bind:open={newShooter} submitLabel={$t`Create and register`} oncreated={(sh) => register(sh.id, sh.clubs[0]?.id)} />
{/if}

{#if dialog?.kind === 'results' && can.results}
	{#key dialog.id}<ResultsDialog {tenantId} {eventId} registrationId={dialog.id} canEdit={can.assign} canPrint={can.print} canPrintForm={can.forms} canIssue={can.awards} canRevoke={can.revoke} onclose={close} onedit={openRegistration} />{/key}
{:else if dialog?.kind === 'edit' && can.assign}
	{#key dialog.id}<RegistrationDialog {tenantId} {eventId} registrationId={dialog.id} lanesOnly={!can.register} onclose={close} onresults={openResults} />{/key}
{:else if dialog?.kind === 'new' && can.register}
	{#key dialog.shooterId}<RegistrationDialog {tenantId} {eventId} shooterId={dialog.shooterId} clubId={dialog.clubId} onclose={close} onresults={openResults} onopen={openRegistration} />{/key}
{/if}
