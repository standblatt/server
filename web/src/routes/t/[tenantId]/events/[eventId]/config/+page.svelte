<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import type { EventInput } from '$lib/api';
	import EventForm from '$lib/EventForm.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { eventQuery, keys, rangesQuery, updateEvent } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const queryClient = useQueryClient();
	const ranges = rangesQuery(() => tenantId);
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);

	/** Saves the general settings of the event. */
	async function save(input: EventInput) {
		await updateEvent(queryClient, tenantId, event.data!, input);
		await queryClient.invalidateQueries({ queryKey: keys.lanes(tenantId, eventId) });
		toast.success($t`Event saved`);
	}
</script>

<PageHeader help="user-guide/events#settings" title={$t`General`} description={$t`Name, days, ranges and language of the event.`} />

{#if event.data && ranges.data}
	{@const ev = event.data}
	{#if ev.archivedAt}
		<p class="text-sm text-muted-foreground">
			{$t`A closed event cannot be changed.`} <a class="underline" href="/t/{tenantId}/events/{eventId}/config/close">{$t`Reopen it`}</a>
		</p>
	{:else}
		<section class="max-w-3xl rounded-lg border bg-card p-4">
			{#key ev.id}<EventForm ranges={ranges.data} initial={{ ...ev, rangeIds: ev.ranges.map((r) => r.id) }} submitLabel={$t`Save`} onsubmit={save} />{/key}
		</section>
	{/if}
{/if}
