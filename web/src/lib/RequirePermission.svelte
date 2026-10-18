<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Snippet } from 'svelte';
	import type { Permission } from './api';
	import { can, membershipOf, meQuery } from './me';

	let { tenantId, permission, children }: { tenantId: number; permission: Permission; children: Snippet } = $props();
	const me = meQuery();
	const membership = $derived(membershipOf(me.data, tenantId));
	const allowed = $derived(can(membership, permission));

	$effect(() => {
		if (membership && !allowed) goto(`/t/${tenantId}`, { replaceState: true });
	});
</script>

{#if allowed}{@render children()}{/if}
