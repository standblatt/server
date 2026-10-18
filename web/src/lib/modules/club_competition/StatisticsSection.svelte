<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import * as Table from '$lib/components/ui/table';
	import PanelHead from '$lib/PanelHead.svelte';
	import { keys } from '$lib/tenant';
	import type { EventProps } from '..';

	let { tenantId, event }: EventProps = $props();
	const club = createQuery(() => ({
		queryKey: [...keys.reports(tenantId, event.id), 'club-competition'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/club-competition', { params: { path: { tenantId, eventId: event.id } } }))
	}));
</script>

{#if club.data}
	<section class="overflow-hidden rounded-lg border bg-card" aria-label={$t`Club competition`}>
		<PanelHead title={$t`Club competition`} hint={club.data.programmes.length ? $t`Counted: ${club.data.programmes.join(', ')}` : $t`No programme counts for the club competition.`} />
		<Table.Root>
			<Table.Header><Table.Row><Table.Head>{$t`Rank`}</Table.Head><Table.Head>{$t`Club`}</Table.Head><Table.Head class="text-right">{$t`Result`}</Table.Head><Table.Head class="text-right">{$t`Participants`}</Table.Head></Table.Row></Table.Header>
			<Table.Body>
				{#each club.data.ranked as c (c.clubId)}
					<Table.Row class={c.excluded ? 'opacity-50' : ''}>
						<Table.Cell class="whitespace-normal">{c.excluded ? $t`out of competition` : c.rank}</Table.Cell>
						<Table.Cell class="font-medium whitespace-normal">{c.name}</Table.Cell>
						<Table.Cell class="text-right tabular-nums">{c.result.toFixed(3)}</Table.Cell>
						<Table.Cell class="text-right tabular-nums">{c.participants}</Table.Cell>
					</Table.Row>
				{/each}
				{#each club.data.unranked as c (c.clubId)}
					<Table.Row class="text-muted-foreground"><Table.Cell class="whitespace-normal">{$t`not ranked`}</Table.Cell><Table.Cell class="whitespace-normal">{c.name}</Table.Cell><Table.Cell></Table.Cell><Table.Cell class="text-right tabular-nums">{c.participants}</Table.Cell></Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</section>
{/if}
