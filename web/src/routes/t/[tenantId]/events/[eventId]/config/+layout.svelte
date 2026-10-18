<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { membershipOf, meQuery } from '$lib/me';

	let { children } = $props();
	const me = meQuery();
	const tenantId = $derived(Number(page.params.tenantId));
	const role = $derived(membershipOf(me.data, tenantId)?.role);

	$effect(() => {
		if (role && role !== 'admin') goto(`/t/${tenantId}/events/${page.params.eventId}`, { replaceState: true });
	});
</script>

{#if role === 'admin'}{@render children()}{/if}
