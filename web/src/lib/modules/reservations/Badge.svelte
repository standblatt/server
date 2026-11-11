<script lang="ts">
	import { plural } from 'lingui-for-svelte/macro';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { bookingRequestsQuery } from '$lib/tenant';
	import type { EventProps } from '..';

	let { tenantId, event }: EventProps = $props();
	const requests = bookingRequestsQuery(
		() => tenantId,
		() => event.id
	);
	const pending = $derived((requests.data ?? []).filter((r) => r.status === 'pending').length);
</script>

{#if pending}<Sidebar.MenuBadge class="rounded-full bg-red-600 text-white peer-hover/menu-button:text-white peer-data-active/menu-button:text-white" aria-label={$plural(pending, { one: '# open request', other: '# open requests' })}>{pending}</Sidebar.MenuBadge>{/if}
