<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { Snippet } from 'svelte';
	import * as NativeSelect from './components/ui/native-select';
	import { currentLocale, localeNames, setLocale } from './locale.svelte';
	import ThemeToggle from './ThemeToggle.svelte';
	import Wordmark from './Wordmark.svelte';

	let { children, wide = false }: { children: Snippet; wide?: boolean } = $props();
</script>

<div class="min-h-svh bg-background">
	<header class="bg-sidebar text-white">
		<div class="mx-auto flex {wide ? 'max-w-6xl' : 'max-w-4xl'} items-center gap-3 px-4 py-3">
			<a href="/" aria-label="Standblatt"><Wordmark class="text-2xl" /></a>
			<div class="ml-auto flex items-center gap-2">
				<NativeSelect.Root
					size="sm"
					class="text-white [&_select]:border-white/30"
					value={currentLocale()}
					onchange={(e) => setLocale((e.currentTarget as HTMLSelectElement).value)}
					aria-label={$t`Language`}
				>
					{#each Object.entries(localeNames) as [code, name] (code)}<NativeSelect.Option value={code} class="text-foreground">{name}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
				<ThemeToggle />
			</div>
		</div>
	</header>
	<main class="mx-auto grid {wide ? 'max-w-6xl' : 'max-w-4xl'} gap-6 px-4 py-8">
		{@render children()}
	</main>
</div>
