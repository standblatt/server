<script lang="ts">
	import { Trans, t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import FileUp from '@lucide/svelte/icons/file-up';
	import Search from '@lucide/svelte/icons/search';
	import UserPlus from '@lucide/svelte/icons/user-plus';
	import { toast } from 'svelte-sonner';
	import { api, call, type EntryResult, type ResultImportReport, type SeriesResult, type Shooter } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import * as Table from '$lib/components/ui/table';
	import EntryCard from '$lib/EntryCard.svelte';
	import { formatDate, programmeTitle } from '$lib/format';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { meQuery } from '$lib/me';
	import NewShooterDialog from '$lib/NewShooterDialog.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import PrintButton from '$lib/PrintButton.svelte';
	import { eventQuery, keys, programmesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const queryClient = useQueryClient();
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const programmes = programmesQuery(
		() => tenantId,
		() => eventId
	);
	const active = $derived((programmes.data ?? []).filter((p) => p.active));
	const archived = $derived(!!event.data?.archivedAt);
	const me = meQuery();
	const isAdmin = $derived(me.data?.tenants.find((x) => x.tenantId === tenantId)?.role === 'admin');

	let importing = $state(false);
	let file = $state<FileList>();
	let report = $state<ResultImportReport | null>(null);
	let importError = $state('');
	let busy = $state(false);
	/** Uploads the result file and shows the report; nothing is stored when a line has an error. */
	async function importCSV(e: SubmitEvent) {
		e.preventDefault();
		if (!file?.[0]) return;
		busy = true;
		importError = '';
		report = null;
		try {
			const form = new FormData();
			form.append('file', file[0]);
			report = await call(
				api.POST('/api/tenants/{tenantId}/events/{eventId}/results/import', {
					params: { path: { tenantId, eventId } },
					body: {} as never,
					bodySerializer: () => form
				})
			);
			await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId] });
		} catch (err) {
			importError = (err as Error).message;
		} finally {
			busy = false;
		}
	}

	let input = $state('');
	let q = $state('');
	let timer: ReturnType<typeof setTimeout>;
	let searchInput = $state<HTMLElement | null>(null);
	/** The search query of the shooters matching text. */
	const shooterSearch = (text: string) => ({
		queryKey: [...keys.shooters(tenantId), text, 0],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/shooters', { params: { path: { tenantId }, query: { q: text, limit: 10 } } }))
	});
	const results = createQuery(() => ({ ...shooterSearch(q), enabled: q !== '' }));

	let shooter = $state<Shooter | null>(null);
	let clubId = $state(0);
	let programmeId = $state(0);
	let creating = $state(false);
	let entered = $state<{ registrationId: number; entryId: number; form: boolean; name: string; title: string; total: number }[]>([]);

	$effect(() => {
		if (!programmeId && active.length) programmeId = active[0].id;
	});
	const programme = $derived(active.find((p) => p.id === programmeId));
	/** An empty result card for the chosen programme; EntryCard fills it. */
	const blank = $derived<EntryResult | null>(
		programme
			? { entryId: 0, programmeId: programme.id, title: programmeTitle(programme), status: 'pending', series: programme.series, maxPoints: programme.maxPoints, combinedSeries: [], results: [], total: 0, mouches: 0, eligibleAwards: [], form: false }
			: null
	);

	/** Debounces the search input into the shooter query. */
	function onInput() {
		clearTimeout(timer);
		timer = setTimeout(() => (q = input.trim()), 250);
	}

	/** Takes a shooter for the next result, with the first of the shooter's clubs. */
	function pick(s: Shooter) {
		shooter = s;
		clubId = s.clubs[0]?.id ?? 0;
		input = q = '';
	}

	/** Enter searches at once (a barcode scanner types faster than the debounce) and takes the first hit. */
	async function onKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') input = q = '';
		if (e.key !== 'Enter' || !input.trim()) return;
		e.preventDefault();
		clearTimeout(timer);
		q = input.trim();
		const list = await queryClient.fetchQuery(shooterSearch(q));
		if (list[0]) pick(list[0]);
	}

	/** Saves the result of the chosen shooter and programme, lists it and gets ready for the next shooter. */
	async function save(series: SeriesResult[]) {
		const sh = shooter!;
		const sheet = await call(
			api.POST('/api/tenants/{tenantId}/events/{eventId}/results', {
				params: { path: { tenantId, eventId } },
				body: { shooterId: sh.id, clubId: clubId || undefined, programmeId, series }
			})
		);
		const entry = sheet.entries.find((e) => e.programmeId === programmeId);
		const name = `${sh.lastname} ${sh.firstname}`;
		entered = [{ registrationId: sheet.registration.id, entryId: entry?.entryId ?? 0, form: !!entry?.form, name, title: entry?.title ?? '', total: entry?.total ?? 0 }, ...entered];
		toast.success($t`${name}: ${entry?.total ?? 0} points saved`);
		shooter = null;
		await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId, 'events', eventId] });
		searchInput?.focus();
	}
</script>

<PageHeader help="user-guide/results-entry" title={$t`Enter results`} description={$t`Results from paper or the scoring system, shooter by shooter, without the desk.`}>
	{#snippet actions()}
		{#if isAdmin}<Button variant="outline" onclick={() => ((importing = true), (report = null), (importError = ''))} disabled={archived}><FileUp />{$t`Import CSV`}</Button>{/if}
		<Button variant="outline" onclick={() => (creating = true)} disabled={archived}><UserPlus />{$t`New shooter`}</Button>
	{/snippet}
</PageHeader>

{#if archived}
	<p class="text-sm text-muted-foreground">{$t`This event is closed.`}</p>
{:else if programmes.data && active.length === 0}
	<p class="text-sm text-muted-foreground">{$t`This event has no active programmes yet.`}</p>
{:else}
	<div class="grid max-w-3xl gap-4">
		<div class="grid gap-3 rounded-lg border bg-card p-4 sm:grid-cols-2">
			<FormField label={$t`Programme`}>
				<NativeSelect.Root class="w-full" bind:value={programmeId}>
					{#each active as p (p.id)}<NativeSelect.Option value={p.id}>{programmeTitle(p)}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
			</FormField>
			<FormField label={$t`Shooter`}>
				<div class="relative">
					<Search class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
					<Input bind:ref={searchInput} class="pl-8" type="search" placeholder={$t`License number or name`} bind:value={input} oninput={onInput} onkeydown={onKeydown} autofocus />
				</div>
			</FormField>
			{#if q && results.data}
				<ul class="grid gap-1 sm:col-span-2" aria-label={$t`Shooters found`}>
					{#each results.data as s (s.id)}
						<li>
							<button type="button" class="flex w-full flex-wrap gap-x-3 rounded-md border px-3 py-2 text-left text-sm hover:bg-muted" onclick={() => pick(s)}>
								<span class="font-medium">{s.lastname} {s.firstname}</span>
								<span class="text-muted-foreground">{[s.licenseNumber, s.birthday ? formatDate(s.birthday) : '', s.clubs.map((c) => c.name).join(', ')].filter(Boolean).join(' · ')}</span>
							</button>
						</li>
					{:else}
						<li class="text-sm text-muted-foreground">{$t`No shooter found.`}</li>
					{/each}
				</ul>
			{/if}
		</div>

		{#if shooter && blank}
			<div class="flex flex-wrap items-center gap-3">
				<span class="font-semibold">{shooter.lastname} {shooter.firstname}</span>
				{#if shooter.clubs.length > 1}
					<NativeSelect.Root size="sm" bind:value={clubId} aria-label={$t`Club`}>
						{#each shooter.clubs as c (c.id)}<NativeSelect.Option value={c.id}>{c.name}</NativeSelect.Option>{/each}
					</NativeSelect.Root>
				{:else if shooter.clubs.length}
					<span class="text-sm text-muted-foreground">{shooter.clubs[0].name}</span>
				{/if}
				<Button variant="ghost" size="sm" class="ml-auto" onclick={() => ((shooter = null), searchInput?.focus())}>{$t`Other shooter`}</Button>
			</div>
			{#key `${shooter.id}-${programmeId}`}<EntryCard entry={blank} {tenantId} print={false} onsave={save} onaward={async () => {}} />{/key}
		{/if}

		{#if entered.length}
			<section class="overflow-hidden rounded-lg border bg-card" aria-label={$t`Entered now`}>
				<h2 class="border-b px-4 py-2.5 text-sm font-semibold">{$t`Entered now`}</h2>
				<Table.Root>
					<Table.Body>
						{#each entered as r, i (i)}
							<Table.Row>
								<Table.Cell class="whitespace-normal"><a class="font-medium hover:underline" href="/t/{tenantId}/events/{eventId}?results={r.registrationId}">{r.name}</a></Table.Cell>
								<Table.Cell class="hidden sm:table-cell">{r.title}</Table.Cell>
								<Table.Cell class="text-right font-semibold tabular-nums">{r.total}</Table.Cell>
								<Table.Cell class="w-0 text-right">{#if r.form}<PrintButton {tenantId} target={{ document: 'form', entryId: r.entryId }} />{/if}</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</section>
		{/if}
	</div>
{/if}

<NewShooterDialog {tenantId} bind:open={creating} submitLabel={$t`Create and enter results`} oncreated={pick} />

<Dialog.Root bind:open={importing}>
	<Dialog.Content class="sm:max-w-xl">
		<Dialog.Header>
			<Dialog.Title>{$t`Import results from CSV`}</Dialog.Title>
			<Dialog.Description>
				<Trans
					>Header row with <code>license</code> or <code>firstname</code> and <code>lastname</code> (optionally <code>birthday</code>), optionally
					<code>club</code>, then <code>programme</code> (barcode or name) and the series <code>s1</code>, <code>s2</code>, …; optionally mouches
					<code>m1</code>, … and single shots <code>s1_1</code>, <code>s1_2</code>, …. One line per shooter and programme. If any line has an error,
					nothing is imported.</Trans
				>
			</Dialog.Description>
		</Dialog.Header>
		<form class="flex flex-wrap items-center gap-2" onsubmit={importCSV}>
			<Input class="flex-1" type="file" name="file" accept=".csv,text/csv" bind:files={file} required />
			<Button type="submit" disabled={busy}>{$t`Import`}</Button>
		</form>
		<FormError message={importError} />
		{#if report}
			{#if report.errors.length > 0}
				<div class="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive" role="alert">
					<p class="font-medium">{$t`Nothing imported, please fix these lines:`}</p>
					<ul class="mt-1 list-disc pl-5">
						{#each report.errors as e, i (i)}<li>{$t`Line ${e.line}: ${e.message}`}</li>{/each}
					</ul>
				</div>
			{:else}
				<p class="rounded-md border border-emerald-500/40 bg-emerald-500/10 px-3 py-2 text-sm" role="status">
					{$t`Imported: ${report.imported} results, ${report.createdShooters} new shooters.`}
				</p>
			{/if}
		{/if}
	</Dialog.Content>
</Dialog.Root>
