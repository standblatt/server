<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { has, moduleUI, type ModuleKey } from '$lib/modules';
	import ModuleSlot from '$lib/modules/ModuleSlot.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { eventQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const event = eventQuery(
		() => tenantId,
		() => Number(page.params.eventId)
	);
	const key = $derived(page.params.module as ModuleKey);
	const ui = $derived(moduleUI[key]);
</script>

{#if ui}
	<PageHeader help={ui.help} title={$t(ui.label)} description={$t(ui.description)} />
{/if}
{#if event.data}
	{#if ui?.configPage && has(event.data, key)}
		<div class="max-w-3xl">{#key key}<ModuleSlot load={ui.configPage} props={{ tenantId, event: event.data }} />{/key}</div>
	{:else}
		<p class="text-sm text-muted-foreground">
			{$t`This module is off for this event.`} <a class="underline" href="/t/{tenantId}/events/{event.data.id}/config/modules">{$t`Modules`}</a>
		</p>
	{/if}
{/if}
