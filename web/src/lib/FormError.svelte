<script lang="ts">
	let { message, class: className = '' }: { message: string; class?: string } = $props();
</script>

{#if message}<p class="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive {className}" role="alert">{message}</p>{/if}
