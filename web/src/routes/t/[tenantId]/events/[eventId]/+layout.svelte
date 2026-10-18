<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Archive from '@lucide/svelte/icons/archive';
	import LayoutTemplate from '@lucide/svelte/icons/layout-template';
	import { intlLocale } from '$lib/locale.svelte';
	import { eventQuery } from '$lib/tenant';

	let { children } = $props();
	const event = eventQuery(
		() => Number(page.params.tenantId),
		() => Number(page.params.eventId)
	);
	$effect(() => {
		const base = `/t/${page.params.tenantId}/events/${page.params.eventId}`;
		if (event.data?.template && !page.url.pathname.startsWith(`${base}/config`)) goto(`${base}/config`, { replaceState: true });
	});
</script>

{#if event.data}
	{#if event.data.template}
		<div class="mb-4 flex items-center gap-2 rounded-md border border-sky-500/40 bg-sky-500/10 px-4 py-2 text-sm" role="status" data-testid="template">
			<LayoutTemplate class="size-4 shrink-0" />
			{$t`A template of the club: new events start from it under Events → New event. It takes no registrations.`}
		</div>
	{/if}
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
