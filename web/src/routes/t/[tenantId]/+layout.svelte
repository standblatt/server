<script lang="ts">
	import { Trans, t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import AppSidebar from '$lib/AppSidebar.svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { eventLive, type Live } from '$lib/live.svelte';
	import { liveLabel } from '$lib/labels';
	import { meQuery, setMyLocale } from '$lib/me';
	import { logout as endSession } from '$lib/my.svelte';
	import { eventsQuery } from '$lib/tenant';
	import OfflineBanner from '$lib/OfflineBanner.svelte';
	import ThemeToggle from '$lib/ThemeToggle.svelte';

	let { children } = $props();
	const queryClient = useQueryClient();
	const me = meQuery();

	const tenantId = $derived(Number(page.params.tenantId));
	const routeEventId = $derived(page.params.eventId ? Number(page.params.eventId) : undefined);
	const membership = $derived(me.data?.tenants.find((t) => t.tenantId === tenantId));
	const events = eventsQuery(() => tenantId);
	const event = $derived(events.data?.find((e) => e.id === routeEventId));

	let live = $state<Live>();
	$effect(() => {
		if (!routeEventId) return;
		const l = eventLive(tenantId, routeEventId, queryClient);
		live = l;
		return () => {
			l.close();
			live = undefined;
		};
	});

	const sections: Record<string, string> = $derived({
		reservations: $t`Reservations`,
		rankings: $t`Rankings`,
		statistics: $t`Statistics`,
		setup: $t`Setup`,
		shooters: $t`Shooters`,
		clubs: $t`Clubs`,
		ranges: $t`Ranges`,
		users: $t`Users`,
		printing: $t`Printing`
	});
	const crumbs = $derived.by(() => {
		const base = `/t/${tenantId}`;
		const rest = page.url.pathname.slice(base.length).split('/').filter(Boolean);
		const list: { label: string; href?: string }[] = [];
		if (routeEventId) {
			list.push({ label: event?.name ?? '…', href: `${base}/events/${routeEventId}` });
			list.push({ label: sections[rest[2]] ?? $t`Dashboard` });
		} else if (rest.length === 0) {
			list.push({ label: $t`Events` });
		} else {
			list.push({ label: sections[rest[0]] ?? rest[0], href: rest.length > 1 ? `${base}/${rest[0]}` : undefined });
			if (rest.length > 2) list.push({ label: sections[rest[2]] ?? rest[2] });
		}
		return list;
	});

	const sidebarOpen = typeof document !== 'undefined' && !document.cookie.includes('sidebar_state=false');

	/** Stores the chosen UI language and applies it. */
	async function changeLocale(locale: string) {
		await setMyLocale(queryClient, locale);
	}
</script>

{#if me.data && !membership}
	<main class="flex min-h-svh items-center justify-center p-6">
		<p class="text-sm"><Trans>Club not found. <a class="font-medium underline" href="/">Back to the list</a></Trans></p>
	</main>
{:else if me.data && membership}
	<Sidebar.Provider open={sidebarOpen}>
		<AppSidebar me={me.data} {tenantId} {event} live={live?.connected} onlogout={() => endSession(queryClient)} onlocale={changeLocale} />
		<Sidebar.Inset class="min-w-0">
			<header class="sticky top-0 z-20 flex h-12 shrink-0 items-center gap-2 border-b bg-background/90 px-3 backdrop-blur">
				<Sidebar.Trigger />
				<Separator orientation="vertical" class="mr-1 h-4 self-center" />
				<Breadcrumb.Root>
					<Breadcrumb.List>
						<Breadcrumb.Item class="hidden md:block">
							<Breadcrumb.Link href="/t/{tenantId}">{membership.tenantName}</Breadcrumb.Link>
						</Breadcrumb.Item>
						{#each crumbs as c, i (i)}
							<Breadcrumb.Separator class={i === 0 ? 'hidden md:block' : ''} />
							<Breadcrumb.Item>
								{#if c.href}<Breadcrumb.Link href={c.href}>{c.label}</Breadcrumb.Link>{:else}<Breadcrumb.Page>{c.label}</Breadcrumb.Page>{/if}
							</Breadcrumb.Item>
						{/each}
					</Breadcrumb.List>
				</Breadcrumb.Root>
				<div class="ml-auto flex items-center gap-2">
					{#if routeEventId}
						<span
							class="flex items-center gap-1.5 text-xs {live?.connected ? 'text-emerald-600 dark:text-emerald-400' : 'text-destructive'}"
							role="status"
							aria-label={$t`Live updates`}
						>
							<span class="size-1.5 rounded-full bg-current"></span>{$t(liveLabel(live?.connected))}
						</span>
					{/if}
					<ThemeToggle />
				</div>
			</header>
			<OfflineBanner />
			<div class="flex flex-1 flex-col p-4 md:p-6">
				{@render children()}
			</div>
		</Sidebar.Inset>
	</Sidebar.Provider>
{/if}
