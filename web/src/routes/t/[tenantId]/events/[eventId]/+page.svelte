<script lang="ts">
	import { page } from '$app/state';
	import { staffRights } from '$lib/desk';
	import EventDesk from '$lib/EventDesk.svelte';
	import SetupChecklist from '$lib/SetupChecklist.svelte';
	import { meQuery } from '$lib/me';
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
	const isAdmin = $derived(me.data?.tenants.find((t) => t.tenantId === tenantId)?.role === 'admin');
	const unready = $derived(!!programmes.data && !programmes.data.some((p) => p.active));
</script>

{#if unready && !event.data?.archivedAt}
	<div class="max-w-3xl"><SetupChecklist {tenantId} {eventId} {isAdmin} /></div>
{:else}
	<EventDesk {tenantId} {eventId} {base} can={staffRights} />
{/if}
