<script lang="ts">
	import { podium, points } from '$lib/format';
	import { t } from 'lingui-for-svelte/macro';
	import type { Snippet } from 'svelte';
	import type { ClassificationRanking } from '$lib/api';
	import { ageCategoryShort } from '$lib/categories';
	import * as Table from '$lib/components/ui/table';

	let { ranking, large = false, actions }: { ranking: ClassificationRanking; large?: boolean; actions?: Snippet } = $props();

</script>

{#each ranking.classes as cl (cl.name)}
	{@const title = cl.name ? `${ranking.name} – ${cl.name}` : ranking.name}
	<section class="overflow-hidden rounded-lg border bg-card" aria-label={$t`Classification ${title}`}>
		<div class="flex flex-wrap items-center gap-x-3 gap-y-2 border-b px-4 py-2.5">
			<h2 class="font-semibold {large ? 'text-2xl' : 'text-sm'}">{title}</h2>
			<span class="text-muted-foreground {large ? 'text-base' : 'text-xs'}">{ranking.parts.map((p) => (p.count > 1 ? `${p.title} ×${p.count}` : p.title)).join(' + ')}</span>
			{#if actions}<div class="ml-auto">{@render actions()}</div>{/if}
		</div>
		<Table.Root class="tabular-nums {large ? 'text-lg [&_td]:py-3' : ''}">
			<Table.Header>
				<Table.Row>
					<Table.Head class="w-14">{$t`Rank`}</Table.Head>
					<Table.Head>{$t`Name`}</Table.Head>
					<Table.Head class="hidden sm:table-cell">{$t`Year`}</Table.Head>
					<Table.Head class="hidden sm:table-cell">{$t`Cat.`}</Table.Head>
					<Table.Head class="hidden sm:table-cell">{$t`Club`}</Table.Head>
					{#each ranking.parts as p (p.programmeId)}<Table.Head class="hidden text-right sm:table-cell">{p.title}</Table.Head>{/each}
					<Table.Head class="text-right">{$t`Total`}</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each [...cl.lines, ...cl.unranked] as l (l.shooterId)}
					<Table.Row class={l.rank ? '' : 'text-muted-foreground'}>
						<Table.Cell>
							{#if l.rank}<span class="inline-flex size-7 items-center justify-center rounded-full font-semibold {podium[l.rank - 1] ?? 'text-muted-foreground'}">{l.rank}</span>
							{:else}<span title={$t`Not all counting results yet`}>–</span>{/if}
						</Table.Cell>
						<Table.Cell class="max-sm:whitespace-normal sm:whitespace-nowrap">
							<span class="font-semibold">{l.lastname}</span> {l.firstname}
							<div class="text-xs text-muted-foreground sm:hidden">{[l.birthYear, ageCategoryShort[l.category], l.club].filter(Boolean).join(' · ')}</div>
						</Table.Cell>
						<Table.Cell class="hidden sm:table-cell">{l.birthYear || '–'}</Table.Cell>
						<Table.Cell class="hidden sm:table-cell">{ageCategoryShort[l.category]}</Table.Cell>
						<Table.Cell class="hidden max-w-48 truncate sm:table-cell">{l.club}</Table.Cell>
						{#each l.parts as p, i (i)}
							<Table.Cell class="hidden text-right sm:table-cell {p.counted ? '' : 'text-muted-foreground line-through'}">{p.points !== undefined ? points(p.points) : '–'}</Table.Cell>
						{/each}
						<Table.Cell class="text-right font-semibold">{points(l.total)}</Table.Cell>
					</Table.Row>
				{:else}
					<Table.Row><Table.Cell colspan={6 + ranking.parts.length} class="py-8 text-center text-muted-foreground">{$t`No results yet.`}</Table.Cell></Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</section>
{/each}
