<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { Snippet } from 'svelte';
	import Wordmark from './Wordmark.svelte';
	import ThemeToggle from './ThemeToggle.svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div class="grid min-h-svh lg:grid-cols-[minmax(0,1fr)_minmax(0,1.1fr)]">
	<div class="relative hidden flex-col justify-between bg-sidebar p-10 text-sidebar-foreground lg:flex">
		<div class="flex items-center gap-3 text-white">
			<Wordmark class="text-3xl" />
		</div>
		<div class="grid max-w-md gap-3">
			<p class="text-2xl font-semibold text-white">{$t`Competitions for clubs, a training log for shooters.`}</p>
			<p class="text-sm">{$t`Registrations, results and rankings live at the range; your own sessions and progress on the phone.`}</p>
		</div>
		<p class="text-xs opacity-60">{$t`Swiss SSV rules`}</p>
	</div>
	<div class="relative flex items-center justify-center p-6 pt-[calc(1.5rem+env(safe-area-inset-top))]">
		<div class="absolute top-[calc(1rem+env(safe-area-inset-top))] right-4"><ThemeToggle /></div>
		<div class="w-full max-w-sm">
			<div class="mb-8 flex items-center gap-3 lg:hidden">
				<Wordmark class="text-3xl" />
			</div>
			{@render children()}
		</div>
	</div>
</div>
