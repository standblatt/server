<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery } from '@tanstack/svelte-query';
	import FileSpreadsheet from '@lucide/svelte/icons/file-spreadsheet';
	import FileText from '@lucide/svelte/icons/file-text';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { pdfClick } from '$lib/print';
	import { keys } from '$lib/tenant';
	import { has, type EventProps } from '..';
	import ClassificationTable from './ClassificationTable.svelte';
	import CertificatesButton from '$lib/forms/CertificatesButton.svelte';

	let { tenantId, event }: EventProps = $props();
	const classifications = createQuery(() => ({
		queryKey: keys.classificationRankings(tenantId, event.id),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/classification-rankings', { params: { path: { tenantId, eventId: event.id } } }))
	}));
	const reportUrl = $derived(`/api/tenants/${tenantId}/events/${event.id}/reports/classifications`);
</script>

{#each classifications.data ?? [] as c (c.id)}
	<ClassificationTable ranking={c}>
		{#snippet actions()}
			<div class="flex items-center gap-2">
				{#if has(event, 'kiosk') && !c.public}<span class="rounded-full border px-2 py-0.5 text-xs text-muted-foreground" title={$t`Not shown on kiosk displays`}>{$t`not on kiosk`}</span>{/if}
				<CertificatesButton {tenantId} {event} title={c.name} classificationId={c.id} classes={c.classes.map((x) => x.name).filter(Boolean)} />
				<Button size="sm" variant="outline" href="{reportUrl}?format=pdf" target="_blank" onclick={pdfClick}><FileText />PDF</Button>
				<Button size="sm" variant="outline" href="{reportUrl}?format=csv"><FileSpreadsheet />CSV</Button>
			</div>
		{/snippet}
	</ClassificationTable>
{/each}
