<script lang="ts">
	import type { Snippet } from 'svelte';

	let { title, hint, children }: { title: string; hint?: string; children: Snippet } = $props();
</script>

<section class="rounded-lg border bg-card p-4">
	<div class="mb-3">
		<h2 class="text-sm font-semibold">{title}</h2>
		{#if hint}<p class="text-xs text-muted-foreground">{hint}</p>{/if}
	</div>
	{@render children()}
</section>
