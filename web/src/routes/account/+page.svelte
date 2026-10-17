<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import AccountSettings from '$lib/AccountSettings.svelte';
	import Wordmark from '$lib/Wordmark.svelte';
</script>

<svelte:head><title>{$t`Account`} · Standblatt</title></svelte:head>

<main class="mx-auto grid max-w-3xl gap-4 p-4 *:min-w-0 md:p-8">
	<header class="flex items-center justify-between gap-3">
		<a href="/" class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"><ArrowLeft class="size-4" />{$t`Back`}</a>
		<Wordmark class="text-xl" />
	</header>
	<h1 class="text-lg font-semibold tracking-tight">{$t`Account`}</h1>
	<AccountSettings />
</main>
