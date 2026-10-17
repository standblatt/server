<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { onMount } from 'svelte';
	import { dev } from '$app/environment';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Building from '@lucide/svelte/icons/building-2';
	import House from '@lucide/svelte/icons/house';
	import LogOut from '@lucide/svelte/icons/log-out';
	import List from '@lucide/svelte/icons/list';
	import Plus from '@lucide/svelte/icons/plus';
	import Settings from '@lucide/svelte/icons/settings';
	import TrendingUp from '@lucide/svelte/icons/trending-up';
	import Wordmark from '$lib/Wordmark.svelte';
	import { meQuery, versionQuery } from '$lib/me';
	import UserCog from '@lucide/svelte/icons/user-cog';
	import CircleHelp from '@lucide/svelte/icons/circle-help';
	import Notices from '$lib/my/Notices.svelte';
	import { flushOutbox, logout, myQuery } from '$lib/my.svelte';
	import ThemeToggle from '$lib/ThemeToggle.svelte';

	let { children } = $props();
	const queryClient = useQueryClient();
	const my = myQuery();
	const me = meQuery();
	const version = versionQuery();

	onMount(() => {
		navigator.serviceWorker?.register('/service-worker.js', { scope: '/my', type: dev ? 'module' : 'classic' }).catch(() => {});
	});

	const userId = $derived(my.data?.userId);
	$effect(() => {
		if (!userId || my.isPlaceholderData) return;
		const id = userId;
		/** Sends the queued offline sessions; runs again when the browser is back online or the app comes to the front. */
		const flush = () => document.visibilityState === 'visible' && flushOutbox(queryClient, id);
		flush();
		window.addEventListener('online', flush);
		window.addEventListener('pageshow', flush);
		document.addEventListener('visibilitychange', flush);
		return () => {
			window.removeEventListener('online', flush);
			window.removeEventListener('pageshow', flush);
			document.removeEventListener('visibilitychange', flush);
		};
	});

	const path = $derived(page.url.pathname);
	const fullscreen = $derived(path.startsWith('/my/train'));
	const tabs = $derived([
		{ href: '/my', label: $t`Today`, icon: House, active: path === '/my' },
		{ href: '/my/log', label: $t`Log`, icon: List, active: path.startsWith('/my/log') },
		{ href: '/my/progress', label: $t`Progress`, icon: TrendingUp, active: path.startsWith('/my/progress') },
		{ href: '/my/settings', label: $t`Settings`, icon: Settings, active: path.startsWith('/my/settings') }
	]);
	const clubs = $derived(me.data?.tenants ?? []);
</script>

{#if my.data && fullscreen}
	{@render children()}
{:else if my.data}
	<div class="min-h-svh md:grid md:grid-cols-[15rem_minmax(0,1fr)]">
		<aside class="sticky top-0 hidden h-svh flex-col gap-1 bg-sidebar p-3 text-sidebar-foreground md:flex">
			<a href="/my" class="mb-4 flex items-center gap-2 px-2 py-1 text-white">
				<Wordmark class="text-2xl" />
			</a>
			<a
				href="/my/train"
				class="mb-3 flex items-center justify-center gap-2 rounded-md bg-sidebar-primary px-3 py-2.5 font-semibold text-sidebar-primary-foreground hover:brightness-105"
			>
				<Plus class="size-4" />{$t`Start training`}
			</a>
			{#each tabs as tab (tab.href)}
				<a
					href={tab.href}
					class="flex items-center gap-3 rounded-md px-3 py-2 text-sm {tab.active ? 'bg-sidebar-accent text-sidebar-accent-foreground' : 'hover:bg-sidebar-accent/60'}"
					aria-current={tab.active ? 'page' : undefined}
				>
					<tab.icon class="size-4" />{tab.label}
				</a>
			{/each}
			<div class="mt-auto grid gap-1">
				{#each clubs as c (c.tenantId)}
					<a href="/t/{c.tenantId}" class="flex items-center gap-3 rounded-md px-3 py-2 text-sm hover:bg-sidebar-accent/60">
						<Building class="size-4" /><span class="truncate">{c.tenantName}</span>
					</a>
				{/each}
				{#if clubs.length > 1}
					<a href="/" class="flex items-center gap-3 rounded-md px-3 py-2 text-sm hover:bg-sidebar-accent/60">
						<Building class="size-4" /><span class="truncate">{$t`All clubs`}</span>
					</a>
				{/if}
				<a href="/account" class="flex items-center gap-3 rounded-md px-3 py-2 text-sm hover:bg-sidebar-accent/60">
					<UserCog class="size-4" /><span class="truncate">{$t`Account`}</span>
				</a>
				<a href="/help/my-training" target="_blank" rel="noopener" class="flex items-center gap-3 rounded-md px-3 py-2 text-sm hover:bg-sidebar-accent/60">
					<CircleHelp class="size-4" /><span class="truncate">{$t`Help`}</span>
				</a>
				<div class="flex items-center justify-between px-2 pt-2 text-xs">
					<span class="truncate">{`${my.data.profile.firstname} ${my.data.profile.lastname}`.trim() || me.data?.user.name}</span>
					<div class="flex items-center">
						<ThemeToggle />
						<button type="button" class="rounded-md p-2 hover:bg-sidebar-accent/60" onclick={() => logout(queryClient)} title={$t`Log out`} aria-label={$t`Log out`}><LogOut class="size-4" /></button>
					</div>
				</div>
				{#if version.data}<p class="truncate px-2 text-[11px] text-sidebar-foreground/60" title={$t`Version`}>Standblatt {version.data.version || 'dev'}</p>{/if}
			</div>
		</aside>
		<div class="min-w-0 pb-24 md:pb-8">
			<header class="sticky top-0 z-20 box-content flex h-12 items-center pt-[env(safe-area-inset-top)] gap-2 border-b bg-background/90 px-4 backdrop-blur md:hidden">
				<Wordmark class="text-xl" />
				<div class="ml-auto flex items-center">
					{#if clubs.length}
						<a href={clubs.length === 1 ? `/t/${clubs[0].tenantId}` : '/'} class="p-2 text-muted-foreground" aria-label={$t`Club`}><Building class="size-4" /></a>
					{/if}
					<ThemeToggle />
					<button type="button" class="p-2 text-muted-foreground" onclick={() => logout(queryClient)} aria-label={$t`Log out`}><LogOut class="size-4" /></button>
				</div>
			</header>
			<main class="mx-auto grid max-w-3xl gap-4 p-4 *:min-w-0 md:p-8">
				<Notices my={my.data} />
				{@render children()}
			</main>
		</div>
		<nav class="fixed inset-x-0 bottom-0 z-30 grid grid-cols-5 border-t bg-background/95 pb-[env(safe-area-inset-bottom)] backdrop-blur md:hidden" aria-label={$t`Main`}>
			{#each tabs.slice(0, 2) as tab (tab.href)}
				<a href={tab.href} class="flex flex-col items-center gap-0.5 py-2 text-[11px] {tab.active ? 'text-foreground' : 'text-muted-foreground'}" aria-current={tab.active ? 'page' : undefined}>
					<tab.icon class="size-5" />{tab.label}
				</a>
			{/each}
			<a href="/my/train" class="flex items-center justify-center" aria-label={$t`Start training`}>
				<span class="-mt-5 flex size-14 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-lg ring-4 ring-background">
					<Plus class="size-7" />
				</span>
			</a>
			{#each tabs.slice(2) as tab (tab.href)}
				<a href={tab.href} class="flex flex-col items-center gap-0.5 py-2 text-[11px] {tab.active ? 'text-foreground' : 'text-muted-foreground'}" aria-current={tab.active ? 'page' : undefined}>
					<tab.icon class="size-5" />{tab.label}
				</a>
			{/each}
		</nav>
	</div>
{/if}
