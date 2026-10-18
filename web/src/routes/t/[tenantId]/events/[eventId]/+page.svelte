<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { staffRights } from '$lib/desk';
	import EventDesk from '$lib/EventDesk.svelte';
	import SetupChecklist from '$lib/SetupChecklist.svelte';
	import { can, membershipOf, meQuery } from '$lib/me';
	import { eventQuery, programmesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const base = $derived(`/t/${tenantId}/events/${eventId}`);
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const programmes = programmesQuery(
		() => tenantId,
		() => eventId
	);
	const me = meQuery();
	const membership = $derived(membershipOf(me.data, tenantId));
	const isAdmin = $derived(can(membership, 'events.configure'));
	const operates = $derived(can(membership, 'events.operate'));
	$effect(() => {
		if (membership && !operates) goto(`${base}/rankings`, { replaceState: true });
	});
	const unready = $derived(!!programmes.data && !programmes.data.some((p) => p.active));
</script>

{#if operates && unready && !event.data?.archivedAt}
	<div class="max-w-3xl"><SetupChecklist {tenantId} {eventId} {isAdmin} /></div>
{:else if operates}
	<EventDesk {tenantId} {eventId} {base} can={staffRights} />
{/if}
