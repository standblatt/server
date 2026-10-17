<script lang="ts">
	import type { Snippet } from 'svelte';
	import HelpButton from './help/HelpButton.svelte';

	let { title, description, actions, help }: { title: string; description?: string; actions?: Snippet; help?: string } = $props();
</script>

<div class="mb-4 flex flex-wrap items-center gap-x-4 gap-y-2">
	<div class="min-w-0">
		<h1 class="flex items-center gap-1 truncate text-lg font-semibold tracking-tight">{title}{#if help}<HelpButton page={help} />{/if}</h1>
		{#if description}<p class="text-sm text-muted-foreground">{description}</p>{/if}
	</div>
	{#if actions}<div class="ml-auto flex flex-wrap items-center gap-2">{@render actions()}</div>{/if}
</div>
