<script lang="ts">
	let {
		items,
		from = 0,
		highlight,
		label
	}: { items: { label: string; value: number; text: string }[]; from?: number; highlight?: number; label: string } = $props();

	const top = $derived(Math.max(...items.map((i) => i.value), from + 1e-9));
</script>

<div class="flex h-36 items-end gap-1.5" role="img" aria-label={label}>
	{#each items as item, i (i)}
		<div class="flex h-full min-w-0 flex-1 flex-col items-center justify-end gap-1">
			<span class="text-[10px] text-muted-foreground tabular-nums">{item.text}</span>
			<div class="flex w-full flex-1 items-end">
				<div
					class="w-full rounded-t-sm {i === highlight ? 'bg-amber-500' : 'bg-chart-3/80'}"
					style:height="{Math.max(2, ((item.value - from) / (top - from)) * 100)}%"
				></div>
			</div>
			<span class="text-xs font-medium tabular-nums">{item.label}</span>
		</div>
	{/each}
</div>
