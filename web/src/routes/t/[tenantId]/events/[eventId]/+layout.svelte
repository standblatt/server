<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import Archive from '@lucide/svelte/icons/archive';
	import { intlLocale } from '$lib/locale.svelte';
	import { eventQuery } from '$lib/tenant';

	let { children } = $props();
	const event = eventQuery(
		() => Number(page.params.tenantId),
		() => Number(page.params.eventId)
	);
</script>

{#if event.data}
	{#if event.data.archivedAt}
		<div class="mb-4 flex items-center gap-2 rounded-md border border-amber-500/40 bg-amber-500/10 px-4 py-2 text-sm" role="status" data-testid="archived">
			<Archive class="size-4 shrink-0" />
			{$t`Closed on ${new Date(event.data.archivedAt).toLocaleDateString(intlLocale())}: the rankings are released, everything stays visible, but nothing can be changed.`}
		</div>
	{/if}
	{@render children()}
{:else if event.error}
	<p class="text-sm text-destructive">{event.error.message}</p>
{/if}
