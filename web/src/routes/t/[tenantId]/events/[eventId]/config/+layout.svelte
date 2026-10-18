<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { meQuery } from '$lib/me';

	let { children } = $props();
	const me = meQuery();
	const tenantId = $derived(Number(page.params.tenantId));
	const role = $derived(me.data?.tenants.find((t) => t.tenantId === tenantId)?.role);

	$effect(() => {
		if (role && role !== 'admin') goto(`/t/${tenantId}/events/${page.params.eventId}`, { replaceState: true });
	});
</script>

{#if role === 'admin'}{@render children()}{/if}
