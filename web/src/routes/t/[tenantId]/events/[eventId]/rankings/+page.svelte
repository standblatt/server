<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import FileSpreadsheet from '@lucide/svelte/icons/file-spreadsheet';
	import FileText from '@lucide/svelte/icons/file-text';
	import { Button } from '$lib/components/ui/button';
	import PageHeader from '$lib/PageHeader.svelte';
	import { pdfClick } from '$lib/print';
	import RankingTable from '$lib/RankingTable.svelte';
	import { keys } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const rankings = createQuery(() => ({
		queryKey: keys.rankings(tenantId, eventId),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/rankings', { params: { path: { tenantId, eventId } } }))
	}));
</script>

<PageHeader help="user-guide/rankings-kiosk" title={$t`Rankings`} description={$t`Live rankings per programme, ranked by the SSV rules.`} />

{#if rankings.data}
	<div class="grid gap-4 *:min-w-0 2xl:grid-cols-2">
		{#each rankings.data as r (r.programmeId)}
			<RankingTable ranking={r}>
				{#snippet actions()}
					<div class="flex items-center gap-2">
						{#if !r.public}<span class="rounded-full border px-2 py-0.5 text-xs text-muted-foreground" title={$t`Not shown on kiosk displays`}>{$t`not on kiosk`}</span>{/if}
						<Button size="sm" variant="outline" href="/api/tenants/{tenantId}/programmes/{r.programmeId}/ranking" target="_blank" onclick={pdfClick}><FileText />PDF</Button>
						<Button size="sm" variant="outline" href="/api/tenants/{tenantId}/programmes/{r.programmeId}/ranking.csv"><FileSpreadsheet />CSV</Button>
					</div>
				{/snippet}
			</RankingTable>
		{:else}
			<p class="text-sm text-muted-foreground">{$t`This event has no programmes yet.`}</p>
		{/each}
	</div>
{:else if rankings.error}
	<p class="text-sm text-destructive">{rankings.error.message}</p>
{/if}
