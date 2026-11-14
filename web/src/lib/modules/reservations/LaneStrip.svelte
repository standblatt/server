<script lang="ts">
	let { count, lanes, disabled = [], pending = false }: { count: number; lanes: number[]; disabled?: number[]; pending?: boolean } = $props();
	const taken = $derived(new Set(lanes));
	const off = $derived(new Set(disabled));
</script>

<div class="hidden h-3 w-40 gap-px sm:flex" aria-hidden="true">
	{#each { length: count }, i (i)}
		<span
			class="flex-1 rounded-[1px] {taken.has(i + 1) ? (pending ? 'bg-amber-500/70' : 'bg-primary') : off.has(i + 1) ? 'bg-muted-foreground/30' : 'bg-muted'}"
		></span>
	{/each}
</div>
