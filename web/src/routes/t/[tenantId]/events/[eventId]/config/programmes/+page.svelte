<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import PageHeader from '$lib/PageHeader.svelte';
	import ProgrammeList from '$lib/ProgrammeList.svelte';
	import { eventQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
</script>

<PageHeader help="user-guide/programmes" title={$t`Programmes`} description={$t`What can be shot: series, prices per age category and awards.`} />

{#if event.data}
	<section class="overflow-hidden rounded-lg border bg-card"><ProgrammeList {tenantId} {eventId} isAdmin={!event.data.archivedAt} /></section>
{/if}
