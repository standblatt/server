<script lang="ts">
	import '../app.css';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { ModeWatcher } from 'mode-watcher';
	import ConfirmHost from '$lib/ConfirmHost.svelte';
	import UpdateNotice from '$lib/UpdateNotice.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { setLinguiContext } from 'lingui-for-svelte';
	import { i18n } from '$lib/locale.svelte';
	import { api, ApiError } from '$lib/api';
	import { startMatomo, trackPage } from '$lib/matomo';
	import { afterNavigate } from '$app/navigation';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
	/** Deployment label (STAGING …) from the server; empty on production. */
	let environment = $state('');
	onMount(() => {
		api.GET('/api/version').then(({ data }) => {
			if (data?.matomo) startMatomo(data.matomo);
			const env = (data?.environment ?? '').trim();
			environment = env && env !== 'production' ? env.toUpperCase() : '';
		});
	});
	/** Appends the environment name to the page title once. */
	function markTitle() {
		if (!environment) return;
		const suffix = ` – ${environment}`;
		if (!document.title.endsWith(suffix)) document.title += suffix;
	}
	$effect(() => {
		if (!environment) return;
		markTitle();
		const title = document.querySelector('title');
		const observer = new MutationObserver(markTitle);
		if (title) observer.observe(title, { childList: true, characterData: true, subtree: true });
		return () => observer.disconnect();
	});
	afterNavigate(({ from }) => {
		if (from) trackPage();
		markTitle();
	});
	const queryClient = new QueryClient({
		defaultOptions: { queries: { retry: (n, err) => n < 3 && !(err instanceof ApiError && err.status < 500) } }
	});
	setLinguiContext(i18n);

	/** Shows API errors that no handler caught as a toast instead of losing them in the console; 401 already redirects. */
	function apiErrorToast(e: PromiseRejectionEvent) {
		if (!(e.reason instanceof ApiError) || e.reason.status === 401) return;
		e.preventDefault();
		toast.error(e.reason.message);
	}
</script>

<svelte:window onunhandledrejection={apiErrorToast} />

<svelte:head>
	<title>Standblatt</title>
</svelte:head>

<ModeWatcher />
{#if environment}
	<div class="pointer-events-none fixed top-0 left-1/2 z-[60] -translate-x-1/2 rounded-b-md bg-red-600 px-3 py-0.5 text-xs font-bold tracking-widest text-white shadow" role="status" data-testid="environment">{environment}</div>
{/if}
<Toaster position="bottom-right" richColors />
<ConfirmHost />
<UpdateNotice />
<QueryClientProvider client={queryClient}>
	<Tooltip.Provider delayDuration={300}>
		{@render children()}
	</Tooltip.Provider>
</QueryClientProvider>
