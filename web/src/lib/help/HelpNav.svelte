<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import Search from '@lucide/svelte/icons/search';
	import { Input } from '../components/ui/input';
	import { allPages, nav } from './nav';

	let { lang, current }: { lang: string; current: string } = $props();

	const pages = $derived(allPages(lang));
	let query = $state('');

	/** Sections of the pages whose heading or text contain the query, with a short excerpt around the hit. */
	function search(all: Awaited<ReturnType<typeof allPages>>, q: string) {
		const needle = q.trim().toLowerCase();
		if (needle.length < 2) return [];
		const hits: { href: string; title: string; heading: string; excerpt: string }[] = [];
		for (const p of all) {
			for (const s of p.sections) {
				const at = s.text.toLowerCase().indexOf(needle);
				if (at < 0 && !s.heading.toLowerCase().includes(needle)) continue;
				const from = Math.max(0, at - 40);
				const excerpt = at < 0 ? s.text.slice(0, 100) : (from > 0 ? '…' : '') + s.text.slice(from, at + needle.length + 60) + '…';
				hits.push({ href: `/help/${p.slug}${s.id ? `#${s.id}` : ''}`, title: p.title, heading: s.heading, excerpt });
				if (hits.length >= 30) return hits;
			}
		}
		return hits;
	}
</script>

<div class="grid gap-4 text-sm">
	<label class="relative block">
		<Search class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
		<Input class="pl-8" type="search" bind:value={query} placeholder={$t`Search the help`} aria-label={$t`Search the help`} />
	</label>
	{#await pages then all}
		{#if query.trim().length >= 2}
			{@const hits = search(all, query)}
			<ul class="grid gap-3" aria-label={$t`Search results`}>
				{#each hits as h (h.href)}
					<li>
						<a href={h.href} class="font-medium hover:underline">{h.heading && h.heading !== h.title ? `${h.title} › ${h.heading}` : h.title}</a>
						<p class="text-xs text-muted-foreground">{h.excerpt}</p>
					</li>
				{:else}
					<li class="text-muted-foreground">{$t`Nothing found.`}</li>
				{/each}
			</ul>
		{:else}
			{@const titles = new Map(all.map((p) => [p.slug, p.title]))}
			<nav class="grid gap-4" aria-label={$t`Help pages`}>
				{#each nav as group (group.label.id)}
					<div class="grid gap-1">
						<div class="text-xs font-semibold tracking-wide text-muted-foreground uppercase">{$t(group.label)}</div>
						{#each group.pages as slug (slug)}
							<a
								href="/help/{slug}"
								class="rounded px-2 py-1 hover:bg-muted {slug === current ? 'bg-muted font-medium' : ''}"
								aria-current={slug === current ? 'page' : undefined}>{titles.get(slug) ?? slug}</a
							>
						{/each}
					</div>
				{/each}
			</nav>
		{/if}
	{/await}
</div>
