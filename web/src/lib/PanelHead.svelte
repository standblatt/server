<script lang="ts">
	let { title, hint }: { title: string; hint?: string } = $props();
</script>

<div class="border-b px-4 py-2.5">
	<h2 class="text-sm font-semibold">{title}</h2>
	{#if hint}<p class="text-xs text-muted-foreground">{hint}</p>{/if}
</div>
