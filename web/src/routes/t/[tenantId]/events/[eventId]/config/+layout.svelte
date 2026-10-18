<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { can, membershipOf, meQuery } from '$lib/me';

	let { children } = $props();
	const me = meQuery();
	const tenantId = $derived(Number(page.params.tenantId));
	const membership = $derived(membershipOf(me.data, tenantId));
	const allowed = $derived(can(membership, 'events.configure'));

	$effect(() => {
		if (membership && !allowed) goto(`/t/${tenantId}/events/${page.params.eventId}`, { replaceState: true });
	});
</script>

{#if allowed}{@render children()}{/if}
