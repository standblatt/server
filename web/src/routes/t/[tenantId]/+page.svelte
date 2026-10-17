<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Plus from '@lucide/svelte/icons/plus';
	import { toast } from 'svelte-sonner';
	import { api, call, type Event, type EventInput } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Table from '$lib/components/ui/table';
	import EventForm from '$lib/EventForm.svelte';
	import { formatDays, rangeNames } from '$lib/format';
	import { meQuery } from '$lib/me';
	import PageHeader from '$lib/PageHeader.svelte';
	import SetupChecklist from '$lib/SetupChecklist.svelte';
	import { clubsQuery, eventsQuery, keys, rangesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const me = meQuery();
	const events = eventsQuery(() => tenantId);
	const ranges = rangesQuery(() => tenantId);
	const clubs = clubsQuery(() => tenantId);
	const isAdmin = $derived(me.data?.tenants.find((t) => t.tenantId === tenantId)?.role === 'admin');
	/** Today's date (YYYY-MM-DD) in the given time zone, UTC if it is invalid. */
	function todayIn(timeZone: string): string {
		try {
			return new Intl.DateTimeFormat('en-CA', { timeZone }).format(new Date());
		} catch {
			return new Date().toISOString().slice(0, 10);
		}
	}
	const groups = $derived.by(() => {
		const g = { today: [] as Event[], upcoming: [] as Event[], past: [] as Event[], archived: [] as Event[] };
		for (const ev of events.data ?? []) {
			const today = todayIn(ev.timezone);
			if (ev.archivedAt) g.archived.push(ev);
			else if (ev.days[0] > today) g.upcoming.push(ev);
			else if (ev.days[ev.days.length - 1] >= today) g.today.push(ev);
			else g.past.push(ev);
		}
		g.upcoming.sort((a, b) => a.startsOn.localeCompare(b.startsOn));
		return g;
	});
	const sections = $derived([
		{ label: $t`Today`, list: groups.today, running: true, muted: false },
		{ label: $t`Upcoming`, list: groups.upcoming, running: false, muted: false },
		{ label: $t`Past`, list: groups.past, running: false, muted: true }
	]);
	let creating = $state(false);

	/** Creates the event and opens its dashboard. */
	async function create(input: EventInput) {
		const ev = await call(api.POST('/api/tenants/{tenantId}/events', { params: { path: { tenantId } }, body: input }));
		await queryClient.invalidateQueries({ queryKey: keys.events(tenantId) });
		creating = false;
		toast.success($t`Event "${ev.name}" created`);
		await goto(`/t/${tenantId}/events/${ev.id}`);
	}
</script>

<PageHeader help="user-guide/events" title={$t`Events`} description={$t`Competitions of this club; choose one to open its dashboard.`}>
	{#snippet actions()}
		{#if isAdmin}<Button onclick={() => (creating = true)}><Plus />{$t`New event`}</Button>{/if}
	{/snippet}
</PageHeader>

{#if events.data?.length === 0}
	<SetupChecklist {tenantId} {isAdmin} onnewevent={() => (creating = true)} />
{:else if events.isPending}
	<p class="py-10 text-center text-muted-foreground">{$t`Loading …`}</p>
{:else if events.error}
	<p class="py-10 text-center text-muted-foreground">{events.error.message}</p>
{:else}
	<div class="space-y-6">
		{#each sections as sec (sec.label)}
			{#if sec.list.length}
				<section>
					<h2 class="mb-2 flex items-center gap-2 text-sm font-medium text-muted-foreground">
						{#if sec.running}<span class="size-2 rounded-full bg-emerald-500"></span>{/if}{sec.label}
					</h2>
					{@render table(sec.list, sec.muted)}
				</section>
			{/if}
		{/each}
		{#if groups.archived.length}
			<details class="group">
				<summary class="mb-2 cursor-pointer text-sm font-medium text-muted-foreground">{$t`Closed`} ({groups.archived.length})</summary>
				{@render table(groups.archived, true)}
			</details>
		{/if}
	</div>
{/if}

{#snippet table(list: Event[], muted: boolean)}
	<div class="overflow-hidden rounded-lg border bg-card">
		<Table.Root>
			<Table.Header>
				<Table.Row><Table.Head class="w-32">{$t`Date`}</Table.Head><Table.Head>{$t`Name`}</Table.Head><Table.Head>{$t`Range`}</Table.Head><Table.Head>{$t`Organizer`}</Table.Head></Table.Row>
			</Table.Header>
			<Table.Body>
				{#each list as ev (ev.id)}
					<Table.Row class="cursor-pointer" onclick={() => goto(`/t/${tenantId}/events/${ev.id}`)}>
						<Table.Cell class="tabular-nums {muted ? 'text-muted-foreground' : ''}">{formatDays(ev.days, ev.startsOn)}</Table.Cell>
						<Table.Cell><a class="font-medium hover:underline" href="/t/{tenantId}/events/{ev.id}">{ev.name}</a></Table.Cell>
						<Table.Cell>{rangeNames(ev.ranges)}</Table.Cell>
						<Table.Cell class="text-muted-foreground">{ev.organizer || '–'}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
{/snippet}

<Dialog.Root bind:open={creating}>
	<Dialog.Content class="max-h-[92svh] overflow-y-auto sm:max-w-xl">
		<Dialog.Header><Dialog.Title>{$t`New event`}</Dialog.Title></Dialog.Header>
		{#if creating && ranges.data}<EventForm ranges={ranges.data} clubs={clubs.data ?? []} submitLabel={$t`Create`} onsubmit={create} />{/if}
	</Dialog.Content>
</Dialog.Root>
