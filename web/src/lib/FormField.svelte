<script lang="ts">
	import type { Snippet } from 'svelte';

	let { label, hint, class: className = '', children }: { label: string; hint?: string; class?: string; children: Snippet } = $props();
</script>

<label class="grid min-w-0 content-start gap-1.5 text-sm {className}">
	<span class="font-medium break-words hyphens-auto">{label}</span>
	{@render children()}
	{#if hint}<span class="text-xs text-muted-foreground">{hint}</span>{/if}
</label>
