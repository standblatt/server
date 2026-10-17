<script lang="ts">
	import { Trans, plural, t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import FileUp from '@lucide/svelte/icons/file-up';
	import Search from '@lucide/svelte/icons/search';
	import UserPlus from '@lucide/svelte/icons/user-plus';
	import { api, call, type ImportReport } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import * as Table from '$lib/components/ui/table';
	import { formatDate } from '$lib/format';
	import FormError from '$lib/FormError.svelte';
	import { meQuery } from '$lib/me';
	import NewShooterDialog from '$lib/NewShooterDialog.svelte';
	import Pager from '$lib/Pager.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import ClubSelect from '$lib/ClubSelect.svelte';
	import ShooterSheet from '$lib/ShooterSheet.svelte';
	import { clubsQuery, keys } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const me = meQuery();
	const clubs = clubsQuery(() => tenantId);
	const isAdmin = $derived(me.data?.tenants.find((t) => t.tenantId === tenantId)?.role === 'admin');

	let input = $state('');
	let q = $state('');
	let clubId = $state(0);
	let timer: ReturnType<typeof setTimeout>;
	/** Debounces the search input. */
	function onInput() {
		clearTimeout(timer);
		timer = setTimeout(() => {
			q = input.trim();
			pageNo = 0;
		}, 250);
	}

	const pageSize = 50;
	let pageNo = $state(0);
	const results = createQuery(() => ({
		queryKey: [...keys.shooters(tenantId), q, clubId, pageNo],
		queryFn: async () => {
			const res = api.GET('/api/tenants/{tenantId}/shooters', {
				params: { path: { tenantId }, query: { q, clubId: clubId || undefined, limit: pageSize, offset: pageNo * pageSize } }
			});
			const items = await call(res);
			return { items, total: Number((await res).response.headers.get('X-Total-Count') ?? items.length) };
		},
		placeholderData: (prev) => prev
	}));

	let creating = $state(false);
	let importing = $state(false);
	const selected = $derived(Number(page.url.searchParams.get('shooter')) || 0);
	/** Opens the shooter sheet. */
	const openShooter = (id: number) => goto(`?shooter=${id}`, { noScroll: true, keepFocus: true });
	/** Closes the shooter sheet. */
	const closeShooter = () => goto(page.url.pathname, { noScroll: true, keepFocus: true, replaceState: true });

	let file = $state<FileList>();
	let report = $state<ImportReport | null>(null);
	let importError = $state('');
	let busy = $state(false);
	/** Choice per unknown club of the file: an ID, "new" or "" while undecided. */
	let clubMap = $state<Record<string, string>>({});
	/** Uploads the CSV file with the clubs mapped so far and shows the import report or the clubs still to map. */
	async function importCSV(e: SubmitEvent) {
		e.preventDefault();
		if (!file?.[0]) return;
		busy = true;
		importError = '';
		const mapped = Object.fromEntries(Object.entries(clubMap).map(([name, v]) => [name, v === 'new' ? 0 : Number(v)]));
		report = null;
		try {
			const form = new FormData();
			form.append('file', file[0]);
			form.append('clubs', JSON.stringify(mapped));
			report = await call(
				api.POST('/api/tenants/{tenantId}/shooters/import', {
					params: { path: { tenantId } },
					body: {} as never,
					bodySerializer: () => form
				})
			);
			for (const u of report.unknownClubs) clubMap[u.name] ??= u.similar.length ? '' : 'new';
			await queryClient.invalidateQueries({ queryKey: keys.clubs(tenantId) });
			await queryClient.invalidateQueries({ queryKey: keys.shooters(tenantId) });
		} catch (err) {
			importError = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

<PageHeader help="user-guide/shooters-clubs" title={$t`Shooters`} description={$t`Master data of the people who shoot for this club.`}>
	{#snippet actions()}
		{#if isAdmin}<Button variant="outline" onclick={() => ((importing = true), (report = null), (importError = ''), (clubMap = {}))}><FileUp />{$t`Import CSV`}</Button>{/if}
		<Button onclick={() => (creating = true)}><UserPlus />{$t`New shooter`}</Button>
	{/snippet}
</PageHeader>

<div class="overflow-hidden rounded-lg border bg-card">
	<div class="flex flex-wrap items-center gap-2 border-b p-3">
		<div class="relative min-w-60 flex-1 sm:max-w-sm">
			<Search class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
			<Input class="pl-8" type="search" placeholder={$t`License number or name`} bind:value={input} oninput={onInput} />
		</div>
		<NativeSelect.Root bind:value={clubId} onchange={() => (pageNo = 0)} aria-label={$t`Club`}>
			<NativeSelect.Option value={0}>{$t`All clubs`}</NativeSelect.Option>
			{#each clubs.data ?? [] as c (c.id)}<NativeSelect.Option value={c.id}>{c.name}</NativeSelect.Option>{/each}
		</NativeSelect.Root>
		<span class="ml-auto text-xs text-muted-foreground">{results.data ? $plural(results.data.total, { one: '# shooter', other: '# shooters' }) : ''}</span>
	</div>
	<Table.Root>
		<Table.Header>
			<Table.Row><Table.Head class="hidden sm:table-cell w-28">{$t`License`}</Table.Head><Table.Head>{$t`Name`}</Table.Head><Table.Head class="hidden sm:table-cell">{$t`Birthday`}</Table.Head><Table.Head class="hidden md:table-cell">{$t`City`}</Table.Head><Table.Head class="hidden sm:table-cell">{$t`Clubs`}</Table.Head></Table.Row>
		</Table.Header>
		<Table.Body>
			{#each results.data?.items ?? [] as sh (sh.id)}
				<Table.Row class="cursor-pointer" onclick={() => openShooter(sh.id)}>
					<Table.Cell class="hidden sm:table-cell font-mono text-xs">{sh.licenseNumber || '–'}</Table.Cell>
					<Table.Cell class="whitespace-normal">
						<a class="font-medium hover:underline" href="?shooter={sh.id}" data-sveltekit-noscroll>{sh.lastname} {sh.firstname}</a>
						<div class="text-xs text-muted-foreground sm:hidden">{[sh.licenseNumber, sh.birthday && formatDate(sh.birthday), ...sh.clubs.map((c) => c.name)].filter(Boolean).join(' · ')}</div>
					</Table.Cell>
					<Table.Cell class="hidden sm:table-cell tabular-nums">{sh.birthday ? formatDate(sh.birthday) : '–'}</Table.Cell>
					<Table.Cell class="hidden md:table-cell">{[sh.postalcode, sh.city].filter(Boolean).join(' ') || '–'}</Table.Cell>
					<Table.Cell class="hidden sm:table-cell text-muted-foreground">{sh.clubs.map((c) => c.name).join(', ') || '–'}</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={5} class="py-10 text-center text-muted-foreground">{results.isPending ? $t`Loading …` : results.error ? results.error.message : $t`No shooters found.`}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
	<Pager bind:page={pageNo} total={results.data?.total ?? 0} size={pageSize} />
</div>

<NewShooterDialog {tenantId} bind:open={creating} oncreated={(sh) => openShooter(sh.id)} />

{#if selected}
	{#key selected}<ShooterSheet {tenantId} shooterId={selected} {isAdmin} onclose={closeShooter} />{/key}
{/if}

<Dialog.Root bind:open={importing}>
	<Dialog.Content class="sm:max-w-xl">
		<Dialog.Header>
			<Dialog.Title>{$t`Import shooters from CSV`}</Dialog.Title>
			<Dialog.Description>
				<Trans
					>Header row with the columns <code>firstname</code> and <code>lastname</code>, optionally <code>license</code>, <code>birthday</code>,
					<code>postalcode</code>, <code>city</code> and <code>club</code>; separated by <code>;</code> or <code>,</code>. One line per shooter and
					club. Shooters are matched by license number (or name and birthday without one); nothing is deleted.</Trans
				>
			</Dialog.Description>
		</Dialog.Header>
		<form class="grid gap-3" onsubmit={importCSV}>
			<div class="flex flex-wrap items-center gap-2">
				<Input class="flex-1" type="file" name="file" accept=".csv,text/csv" bind:files={file} onchange={() => ((report = null), (clubMap = {}))} required />
				<Button type="submit" disabled={busy}>{$t`Import`}</Button>
			</div>
			{#if report?.unknownClubs.length}
				<div class="grid gap-2 rounded-md border bg-muted/40 p-3 text-sm" role="status">
					<p class="font-medium">{$t`Nothing imported yet: these clubs are not in the club list. Choose an existing club or a new one for each, then import again.`}</p>
					{#each report.unknownClubs as u (u.name)}
						<div class="grid items-center gap-1 sm:grid-cols-2 sm:gap-3 *:min-w-0">
							<span class="font-medium break-words">{u.name}</span>
							<ClubSelect {tenantId} name={u.name} suggested={u.similar} bind:value={clubMap[u.name]} label={$t`Club for ${u.name}`} />
						</div>
					{/each}
				</div>
			{/if}
		</form>
		<FormError message={importError} />
		{#if report && !report.unknownClubs.length}
			{#if report.errors.length > 0}
				<div class="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive" role="alert">
					<p class="font-medium">{$t`Nothing imported, please fix these lines:`}</p>
					<ul class="mt-1 list-disc pl-5">
						{#each report.errors as e (e.line)}<li>{$t`Line ${e.line}: ${e.message}`}</li>{/each}
					</ul>
				</div>
			{:else}
				<p class="rounded-md border border-emerald-500/40 bg-emerald-500/10 px-3 py-2 text-sm" role="status">
					{$t`Imported: ${report.created} new, ${report.updated} updated, ${report.clubsCreated} new clubs, ${report.links} club assignments.`}
				</p>
			{/if}
		{/if}
	</Dialog.Content>
</Dialog.Root>
