<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { has, moduleUI, type ModuleKey } from '$lib/modules';
	import ModuleSlot from '$lib/modules/ModuleSlot.svelte';
	import { eventQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const event = eventQuery(
		() => tenantId,
		() => Number(page.params.eventId)
	);
	const key = $derived(page.params.module as ModuleKey);
	const load = $derived(moduleUI[key]?.operationPage);
</script>

{#if event.data}
	{#if load && has(event.data, key)}
		{#key key}<ModuleSlot {load} {tenantId} event={event.data} />{/key}
	{:else}
		<p class="text-sm text-muted-foreground">{$t`This module is off for this event.`}</p>
	{/if}
{/if}
