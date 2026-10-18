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
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { Input } from '$lib/components/ui/input';
	import SetupChecklist from '$lib/SetupChecklist.svelte';
	import { presets } from '$lib/modules/presets';
	import { eventsQuery, keys, rangesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const me = meQuery();
	const events = eventsQuery(() => tenantId);
	const ranges = rangesQuery(() => tenantId);
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
			if (ev.template) continue;
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
	let preset = $state(presets[0].id);
	const templates = $derived((events.data ?? []).filter((ev) => ev.template));
	const template = $derived(templates.find((ev) => `template-${ev.id}` === preset));
	let copyName = $state('');
	let copyDate = $state(new Date().toISOString().slice(0, 10));
	let copyError = $state('');

	/** Creates an event from a template of the club. */
	async function createFromTemplate(e: SubmitEvent) {
		e.preventDefault();
		copyError = '';
		try {
			const ev = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/copy', { params: { path: { tenantId, eventId: template!.id } }, body: { name: copyName.trim(), startsOn: copyDate } }));
			await queryClient.invalidateQueries({ queryKey: keys.events(tenantId) });
			creating = false;
			toast.success($t`Event "${ev.name}" created`);
			await goto(`/t/${tenantId}/events/${ev.id}`);
		} catch (err) {
			copyError = (err as Error).message;
		}
	}

	/** Creates the event and opens its dashboard. */
	async function create(input: EventInput) {
		const modules = presets.find((p) => p.id === preset)?.modules;
		const ev = await call(api.POST('/api/tenants/{tenantId}/events', { params: { path: { tenantId } }, body: { ...input, modules } }));
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
		{#if templates.length}
			<details class="group">
				<summary class="mb-2 cursor-pointer text-sm font-medium text-muted-foreground">{$t`Templates`} ({templates.length})</summary>
				<div class="overflow-hidden rounded-lg border bg-card">
					<ul class="divide-y">
						{#each templates as tpl (tpl.id)}
							<li class="px-4 py-2.5 text-sm"><a class="font-medium hover:underline" href="/t/{tenantId}/events/{tpl.id}/config">{tpl.name}</a></li>
						{/each}
					</ul>
				</div>
			</details>
		{/if}
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
				<Table.Row><Table.Head class="sm:w-32">{$t`Date`}</Table.Head><Table.Head>{$t`Name`}</Table.Head><Table.Head class="hidden sm:table-cell">{$t`Range`}</Table.Head><Table.Head class="hidden sm:table-cell">{$t`Organizer`}</Table.Head></Table.Row>
			</Table.Header>
			<Table.Body>
				{#each list as ev (ev.id)}
					<Table.Row class="cursor-pointer" onclick={() => goto(`/t/${tenantId}/events/${ev.id}`)}>
						<Table.Cell class="tabular-nums {muted ? 'text-muted-foreground' : ''}">{formatDays(ev.days, ev.startsOn)}</Table.Cell>
						<Table.Cell class="whitespace-normal"><a class="font-medium hover:underline" href="/t/{tenantId}/events/{ev.id}">{ev.name}</a><div class="text-xs text-muted-foreground sm:hidden">{rangeNames(ev.ranges)}</div></Table.Cell>
						<Table.Cell class="hidden sm:table-cell">{rangeNames(ev.ranges)}</Table.Cell>
						<Table.Cell class="hidden sm:table-cell text-muted-foreground">{ev.organizer || '–'}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
{/snippet}

<Dialog.Root bind:open={creating}>
	<Dialog.Content class="max-h-[92svh] overflow-y-auto sm:max-w-2xl">
		<Dialog.Header><Dialog.Title>{$t`New event`}</Dialog.Title></Dialog.Header>
		<fieldset class="grid gap-2">
			<legend class="mb-2 text-sm font-medium">{$t`What kind of event?`}</legend>
			{#each presets as p (p.id)}
				<label class="flex cursor-pointer items-start gap-3 rounded-md border p-3 text-sm has-[:checked]:border-primary has-[:checked]:bg-primary/5">
					<input type="radio" name="preset" class="mt-1 accent-primary" value={p.id} bind:group={preset} />
					<span><span class="font-medium">{$t(p.label)}</span><span class="block text-muted-foreground">{$t(p.description)}</span></span>
				</label>
			{/each}
			<p class="text-xs text-muted-foreground">{$t`Only sets the modules; they can be changed any time under Configuration → Modules.`}</p>
			{#each templates as tpl (tpl.id)}
				<label class="flex cursor-pointer items-start gap-3 rounded-md border p-3 text-sm has-[:checked]:border-primary has-[:checked]:bg-primary/5">
					<input type="radio" name="preset" class="mt-1 accent-primary" value="template-{tpl.id}" bind:group={preset} />
					<span><span class="font-medium">{tpl.name}</span><span class="block text-muted-foreground">{$t`Template of the club: modules, settings, ranges, programmes and classifications`}</span></span>
				</label>
			{/each}
		</fieldset>
		{#if template}
			<form class="grid gap-4 sm:grid-cols-2" onsubmit={createFromTemplate}>
				<FormField label={$t`Name`} class="sm:col-span-2"><Input bind:value={copyName} required maxlength={200} /></FormField>
				<FormField label={$t`Date`}><Input type="date" bind:value={copyDate} required /></FormField>
				<FormError message={copyError} class="sm:col-span-2" />
				<div class="flex justify-end sm:col-span-2"><Button type="submit">{$t`Create`}</Button></div>
			</form>
		{:else if creating && ranges.data}<EventForm ranges={ranges.data} submitLabel={$t`Create`} onsubmit={create} />{/if}
	</Dialog.Content>
</Dialog.Root>
