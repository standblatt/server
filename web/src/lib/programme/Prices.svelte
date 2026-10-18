<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { ageCategories } from '$lib/categories';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { ageCategoryLabel } from '$lib/labels';
	import type { ProgrammeDraft } from './draft.svelte';

	let { draft, locked }: { draft: ProgrammeDraft; locked: boolean } = $props();
</script>

<div class="grid grid-cols-2 gap-3 *:min-w-0 sm:grid-cols-3 lg:grid-cols-5">
	{#each ageCategories as c (c)}
		{@const p = draft.prices[c]}
		<div class="grid min-w-0 gap-1.5 rounded-md border p-2.5 {p.enabled ? '' : 'bg-muted/40'}">
			<label class="flex items-center gap-2 text-sm font-medium"><Checkbox bind:checked={p.enabled} disabled={locked && draft.fixedPrices.has(c)} /><span class="min-w-0 break-words hyphens-auto">{$t(ageCategoryLabel(c))}</span></label>
			<Input class="h-8" type="number" min="0" step="0.05" bind:value={p.price} disabled={!p.enabled} aria-label={$t`Price ${$t(ageCategoryLabel(c))}`} />
			{#if draft.repeatable}
				<Input
					class="h-8"
					type="number"
					min="0"
					step="0.05"
					bind:value={p.repeat}
					disabled={!p.enabled}
					placeholder={$t`Further: same`}
					title={$t`Price of every further attempt`}
					aria-label={$t`Further attempts ${$t(ageCategoryLabel(c))}`}
				/>
			{/if}
		</div>
	{/each}
</div>
