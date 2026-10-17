<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import CircleHelp from '@lucide/svelte/icons/circle-help';
	import ExternalLink from '@lucide/svelte/icons/external-link';
	import { Button } from '../components/ui/button';
	import * as Sheet from '../components/ui/sheet';
	import { currentLocale } from '../locale.svelte';
	import HelpArticle from './HelpArticle.svelte';

	let { page }: { page: string } = $props();

	let open = $state(false);
	const slug = $derived(page.split('#')[0]);
	const anchor = $derived(page.split('#')[1] ?? '');
</script>

<Button variant="ghost" size="icon-sm" aria-label={$t`Help`} title={$t`Help`} onclick={() => (open = true)}><CircleHelp /></Button>

<Sheet.Root bind:open>
	<Sheet.Content class="gap-0 p-0 data-[side=right]:w-full data-[side=right]:sm:max-w-xl">
		<Sheet.Header class="border-b">
			<Sheet.Title>{$t`Help`}</Sheet.Title>
			<a href="/help/{page}" target="_blank" rel="noopener" class="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
				><ExternalLink class="size-3.5" />{$t`Open the full help`}</a
			>
		</Sheet.Header>
		<div class="overflow-y-auto p-4">
			{#if open}<HelpArticle {slug} {anchor} lang={currentLocale()} newTab />{/if}
		</div>
	</Sheet.Content>
</Sheet.Root>
