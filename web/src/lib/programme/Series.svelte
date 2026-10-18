<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import type { ProgrammeDraft } from './draft.svelte';

	let { draft, locked }: { draft: ProgrammeDraft; locked: boolean } = $props();

	/** Appends a series that copies the shots, maximum and mouches setting of the last one. */
	function addSeries() {
		const last = draft.series[draft.series.length - 1];
		draft.series.push({ shots: last?.shots ?? 10, maxPoints: last?.maxPoints ?? 100, name: '', trackMouches: last?.trackMouches ?? false, trackShots: last?.trackShots ?? false });
	}
</script>

<div class="overflow-hidden rounded-md border">
	<Table.Root>
		<Table.Header>
			<Table.Row><Table.Head class="hidden w-10 sm:table-cell">#</Table.Head><Table.Head>{$t`Shots`}</Table.Head><Table.Head class="whitespace-normal">{$t`Max. points`}</Table.Head><Table.Head>{$t`Name`}</Table.Head><Table.Head><span class="sm:hidden" title={$t`Mouches`}>{$t({ message: "M", comment: "Column header: mouches" })}</span><span class="hidden sm:inline">{$t`Mouches`}</span></Table.Head><Table.Head><span class="sm:hidden" title={$t`Single shots`}>{$t({ message: "SH", comment: "Column header: single shots" })}</span><span class="hidden sm:inline">{$t`Single shots`}</span></Table.Head><Table.Head></Table.Head></Table.Row>
		</Table.Header>
		<Table.Body>
			{#each draft.series as s, i (i)}
				<Table.Row>
					<Table.Cell class="hidden text-muted-foreground sm:table-cell">{i + 1}</Table.Cell>
					<Table.Cell><Input class="h-8 w-12 sm:w-20" type="number" min="1" max="100" bind:value={s.shots} required disabled={locked} /></Table.Cell>
					<Table.Cell><Input class="h-8 w-14 sm:w-24" type="number" min="1" bind:value={s.maxPoints} required disabled={locked} /></Table.Cell>
					<Table.Cell class="w-full"><Input class="h-8 min-w-12" bind:value={s.name} /></Table.Cell>
					<Table.Cell><Checkbox bind:checked={s.trackMouches} disabled={locked} aria-label={$t`Track mouches in series ${i + 1}`} /></Table.Cell>
					<Table.Cell><Checkbox bind:checked={s.trackShots} disabled={locked} aria-label={$t`Record single shots in series ${i + 1}`} /></Table.Cell>
					<Table.Cell class="text-right">
						<Button variant="ghost" size="icon-sm" title={$t`Remove series`} aria-label={$t`Remove series ${i + 1}`} onclick={() => draft.series.splice(i, 1)} disabled={locked || draft.series.length === 1}><Trash /></Button>
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
<Button variant="outline" size="sm" class="mt-2" onclick={addSeries} disabled={locked}><Plus />{$t`Add series`}</Button>
