<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import ChevronLeft from '@lucide/svelte/icons/chevron-left';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import { Button } from './components/ui/button';
	import { lastPage } from './format';

	let { page = $bindable(), total, size }: { page: number; total: number; size: number } = $props();

	const current = $derived(Math.min(page, lastPage(total, size)));
	const from = $derived(Math.min(total, current * size + 1));
	const to = $derived(Math.min(total, (current + 1) * size));
</script>

{#if total > size}
	<nav class="flex items-center justify-end gap-2 border-t px-3 py-2 text-sm text-muted-foreground" aria-label={$t`Pages`}>
		<span class="tabular-nums">{$t`${from}–${to} of ${total}`}</span>
		<Button variant="ghost" size="icon-sm" aria-label={$t`Previous page`} disabled={current === 0} onclick={() => (page = current - 1)}><ChevronLeft /></Button>
		<Button variant="ghost" size="icon-sm" aria-label={$t`Next page`} disabled={to >= total} onclick={() => (page = current + 1)}><ChevronRight /></Button>
	</nav>
{/if}
