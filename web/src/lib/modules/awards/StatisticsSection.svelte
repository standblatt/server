<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import * as Table from '$lib/components/ui/table';
	import { chf } from '$lib/format';
	import PanelHead from '$lib/PanelHead.svelte';
	import { keys } from '$lib/tenant';
	import type { EventProps } from '..';

	let { tenantId, event }: EventProps = $props();
	const stats = createQuery(() => ({
		queryKey: [...keys.reports(tenantId, event.id), 'statistics'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/statistics', { params: { path: { tenantId, eventId: event.id } } }))
	}));
	/** Translated label of an award kind. */
	const kind = (k: string) => (k === 'card' ? $t`Wreath card` : k === 'cash' ? $t`Cash prize` : $t`Prize in kind`);
</script>

{#if stats.data}
	<section class="overflow-hidden rounded-lg border bg-card">
		<PanelHead title={$t`Awards`} />
		<Table.Root>
			<Table.Header><Table.Row><Table.Head class="hidden sm:table-cell">{$t`Programme`}</Table.Head><Table.Head>{$t`Award`}</Table.Head><Table.Head class="hidden sm:table-cell">{$t`Kind`}</Table.Head><Table.Head class="text-right">{$t`Count`}</Table.Head><Table.Head class="text-right">{$t`Value`}</Table.Head></Table.Row></Table.Header>
			<Table.Body>
				{#each stats.data.awards as a, i (i)}
					<Table.Row><Table.Cell class="hidden sm:table-cell">{a.programme}</Table.Cell><Table.Cell class="whitespace-normal"><div class="font-medium">{a.name}</div><div class="text-xs text-muted-foreground sm:hidden">{a.programme} · {kind(a.kind)}</div></Table.Cell><Table.Cell class="hidden sm:table-cell">{kind(a.kind)}</Table.Cell><Table.Cell class="text-right tabular-nums">{a.count}</Table.Cell><Table.Cell class="text-right tabular-nums">{chf(a.total)}</Table.Cell></Table.Row>
				{:else}
					<Table.Row><Table.Cell colspan={5} class="py-6 text-center text-muted-foreground">{$t`No awards given yet.`}</Table.Cell></Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</section>
{/if}
