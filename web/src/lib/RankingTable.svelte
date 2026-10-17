<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import type { Snippet } from 'svelte';
	import type { ProgrammeRanking } from './api';
	import { ageCategoryShort } from './categories';
	import * as Table from './components/ui/table';

	let { ranking, large = false, actions }: { ranking: ProgrammeRanking; large?: boolean; actions?: Snippet } = $props();

	const podium = ['bg-amber-400 text-amber-950', 'bg-slate-300 text-slate-900', 'bg-orange-300 text-orange-950'];
	/** Formats a value without decimals when whole, else with one. */
	const fmt = (v: number) => (Number.isInteger(v) ? String(v) : v.toFixed(1));
</script>

<section class="overflow-hidden rounded-lg border bg-card" aria-label={$t`Ranking ${ranking.title}`}>
	<div class="flex flex-wrap items-center gap-x-3 gap-y-2 border-b px-4 py-2.5">
		<h2 class="font-semibold {large ? 'text-2xl' : 'text-sm'}">{ranking.title}</h2>
		<span class="text-muted-foreground {large ? 'text-base' : 'text-xs'}">{$t`${ranking.shots} shots`} · {$plural(ranking.lines.length, { one: '# shooter', other: '# shooters' })}</span>
		{#if actions}<div class="ml-auto">{@render actions()}</div>{/if}
	</div>
	<Table.Root class="tabular-nums {large ? 'text-lg [&_td]:py-3' : ''}">
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-14">{$t`Rank`}</Table.Head>
				<Table.Head>{$t`Name`}</Table.Head>
				<Table.Head>{$t`Year`}</Table.Head>
				<Table.Head>{$t`Cat.`}</Table.Head>
				<Table.Head>{$t`Club`}</Table.Head>
				{#each Array(ranking.seriesCount) as _, i (i)}<Table.Head class="text-right">{$t`S${i + 1}`}</Table.Head>{/each}
				<Table.Head class="text-right">{$t`Total`}</Table.Head>
				{#if ranking.tracksMouches}<Table.Head class="text-right">{$t({ message: "M", comment: "Column header: mouches" })}</Table.Head>{/if}
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each ranking.lines as l (l.entryId)}
				<Table.Row>
					<Table.Cell>
						<span class="inline-flex size-7 items-center justify-center rounded-full font-semibold {podium[l.rank - 1] ?? 'text-muted-foreground'}">{l.rank}</span>
					</Table.Cell>
					<Table.Cell class="whitespace-nowrap"><span class="font-semibold">{l.lastname}</span> {l.firstname}</Table.Cell>
					<Table.Cell>{l.birthYear || '–'}</Table.Cell>
					<Table.Cell>{ageCategoryShort[l.category]}</Table.Cell>
					<Table.Cell class="max-w-48 truncate">{l.club}</Table.Cell>
					{#each Array(ranking.seriesCount) as _, i (i)}<Table.Cell class="text-right text-muted-foreground">{l.series[i] !== undefined ? fmt(l.series[i]) : '–'}</Table.Cell>{/each}
					<Table.Cell class="text-right font-semibold">{fmt(l.total)}</Table.Cell>
					{#if ranking.tracksMouches}<Table.Cell class="text-right">{l.mouches}</Table.Cell>{/if}
				</Table.Row>
			{:else}
				<Table.Row><Table.Cell colspan={6 + ranking.seriesCount} class="py-8 text-center text-muted-foreground">{$t`No results yet.`}</Table.Cell></Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</section>
