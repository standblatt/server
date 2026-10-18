<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { Snippet } from 'svelte';
	import { Switch } from '$lib/components/ui/switch';
	import HelpButton from '$lib/help/HelpButton.svelte';
	import type { ModuleUI } from './types';

	/** A module as a card: whether it is on, why its switch is locked (blocked) or what switching does besides (note). */
	type Card = { key: string; ui: Pick<ModuleUI, 'label' | 'description' | 'help' | 'icon'>; on: boolean; blocked: string; note: string };

	let { cards, ontoggle, settings }: { cards: Card[]; ontoggle: (key: string, on: boolean) => void; settings?: Snippet<[string]> } = $props();
</script>

<ul class="grid max-w-5xl gap-3 md:grid-cols-2">
	{#each cards as c (c.key)}
		<li class="flex flex-col gap-3 rounded-lg border bg-card p-4 {c.on ? '' : 'opacity-80'}" data-testid="module-{c.key}">
			<div class="flex items-start gap-3">
				<c.ui.icon class="mt-0.5 size-5 shrink-0 {c.on ? 'text-primary' : 'text-muted-foreground'}" />
				<div class="min-w-0 flex-1 text-sm">
					<div class="flex items-center gap-1 font-medium">{$t(c.ui.label)}<HelpButton page={c.ui.help} /></div>
					<p class="text-muted-foreground">{$t(c.ui.description)}</p>
					{#if c.blocked || c.note}<p class="mt-1 text-xs text-muted-foreground">{c.blocked || c.note}</p>{/if}
				</div>
				<Switch checked={c.on} disabled={!!c.blocked} aria-label={$t(c.ui.label)} onCheckedChange={(v) => ontoggle(c.key, v)} />
			</div>
			{#if c.on && settings}{@render settings(c.key)}{/if}
		</li>
	{/each}
</ul>
