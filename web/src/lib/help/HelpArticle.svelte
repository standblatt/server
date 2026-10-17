<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { loadPage } from './markdown';

	let { slug, lang, anchor = '', newTab = false }: { slug: string; lang: string; anchor?: string; newTab?: boolean } = $props();

	const page = $derived(loadPage(lang, slug, newTab));

	/** Scrolls to the anchor once the rendered page is in the DOM. */
	function scrollToAnchor(node: HTMLElement) {
		if (anchor) node.querySelector(`#${CSS.escape(anchor)}`)?.scrollIntoView();
	}
</script>

<svelte:head>
	{#await page then p}<title>{p?.title ?? ''} – {$t`Help`}</title>{/await}
</svelte:head>

{#await page}
	<p class="text-sm text-muted-foreground">{$t`Loading…`}</p>
{:then p}
	{#if p}
		{#if p.lang !== lang}<p class="mb-4 rounded-md bg-muted px-3 py-2 text-sm text-muted-foreground">{$t`This page is not translated yet.`}</p>{/if}
		<article class="help-prose" use:scrollToAnchor>
			{@html p.html}
		</article>
	{:else}
		<p class="rounded-lg border border-dashed p-6 text-center text-muted-foreground">{$t`This help page does not exist.`}</p>
	{/if}
{/await}
