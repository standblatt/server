<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { fmt } from '$lib/my.svelte';
	import { direction } from '$lib/targets';

	let { group }: { group: { x: number; y: number; radius: number; shots: number } } = $props();

	const dir = $derived(direction(group.x, group.y, Math.max(0.3, group.radius * 0.3)));
	const offset = $derived(Math.hypot(group.x, group.y));
	const horizontal = $derived(dir.horizontal === 'right' ? $t`${fmt(Math.abs(group.x), 1)} mm right` : dir.horizontal === 'left' ? $t`${fmt(Math.abs(group.x), 1)} mm left` : '');
	const vertical = $derived(dir.vertical === 'high' ? $t`${fmt(Math.abs(group.y), 1)} mm high` : dir.vertical === 'low' ? $t`${fmt(Math.abs(group.y), 1)} mm low` : '');
	const sight = $derived(
		[dir.vertical === 'low' ? $t`up` : dir.vertical === 'high' ? $t`down` : '', dir.horizontal === 'right' ? $t`left` : dir.horizontal === 'left' ? $t`right` : '']
			.filter(Boolean)
			.join(' + ')
	);
</script>

<div class="grid content-start gap-3 text-sm">
	<h2 class="font-semibold">{$t`Your group`}</h2>
	<dl class="grid grid-cols-2 gap-2">
		<div class="rounded-lg bg-muted p-3">
			<dt class="text-xs text-muted-foreground">{$t`Group size`}</dt>
			<dd class="text-xl font-semibold tabular-nums">{fmt(group.radius * 2, 1)} mm</dd>
		</div>
		<div class="rounded-lg bg-muted p-3">
			<dt class="text-xs text-muted-foreground">{$t`Off centre`}</dt>
			<dd class="text-xl font-semibold tabular-nums">{fmt(offset, 1)} mm</dd>
		</div>
	</dl>
	{#if horizontal || vertical}
		<p>{$t`The centre of your ${group.shots} hits sits ${[horizontal, vertical].filter(Boolean).join(', ')}.`}</p>
		<p class="rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-3">{$t`Tip: if this shows up again and again, adjust your sight ${sight} — or look at your hold and trigger first.`}</p>
	{:else}
		<p>{$t`Your group is well centred. Work on making it smaller.`}</p>
	{/if}
	<p class="text-xs text-muted-foreground">{$t`Group size is twice the average distance of your hits from their centre.`}</p>
</div>
