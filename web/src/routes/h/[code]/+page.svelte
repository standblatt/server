<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import { helperRights } from '$lib/desk';
	import EventDesk from '$lib/EventDesk.svelte';
	import OfflineBanner from '$lib/OfflineBanner.svelte';
	import { formatDays } from '$lib/format';
	import { helperRoleLabel, liveLabel } from '$lib/labels';
	import { eventLive, type Live } from '$lib/live.svelte';
	import { localeNames, setLocale, type Locale } from '$lib/locale.svelte';
	import Logo from '$lib/Logo.svelte';
	import { eventQuery } from '$lib/tenant';
	import ThemeToggle from '$lib/ThemeToggle.svelte';

	const code = $derived(page.params.code ?? '');
	const base = $derived(`/h/${code}`);
	const queryClient = useQueryClient();
	const language = page.url.searchParams.get('language');

	const helper = createQuery(() => ({
		queryKey: ['helper', code],
		queryFn: () => call(api.GET('/api/helper')),
		retry: false
	}));
	const tenantId = $derived(helper.data?.tenantId ?? 0);
	const eventId = $derived(helper.data?.eventId ?? 0);
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	$effect(() => {
		if (event.data) setLocale(language && language in localeNames ? (language as Locale) : event.data.locale, false);
	});
	const only = $derived(helper.data?.role === 'scorer' && helper.data.rangeId ? { rangeId: helper.data.rangeId, lanes: helper.data.lanes } : undefined);
	const rangeName = $derived(event.data?.ranges.find((r) => r.id === only?.rangeId)?.name ?? '');

	let revoked = $state(false);
	$effect(() => {
		/** Marks the helper link as revoked. */
		const onUnauthorized = () => (revoked = true);
		window.addEventListener('standblatt:unauthorized', onUnauthorized);
		return () => window.removeEventListener('standblatt:unauthorized', onUnauthorized);
	});
	let live = $state<Live>();
	$effect(() => {
		if (!helper.data || revoked) return;
		const l = eventLive(helper.data.tenantId, helper.data.eventId, queryClient, code);
		live = l;
		return () => l.close();
	});
</script>

<svelte:head><title>{event.data ? `${event.data.name} – ${$t(helperRoleLabel(helper.data?.role ?? ''))}` : 'Standblatt'}</title></svelte:head>

<main class="flex min-h-svh flex-col">
	{#if helper.error || revoked}
		<p class="p-10 text-xl text-destructive" role="alert">{$t`This helper link is not valid (anymore).`}</p>
	{:else if helper.data && event.data}
		<header class="sticky top-0 z-20 flex items-center gap-3 bg-sidebar px-4 pt-[calc(0.625rem+env(safe-area-inset-top))] pb-2.5 text-sidebar-foreground">
			<Logo class="size-8 shrink-0 text-sidebar-primary" />
			<div class="min-w-0">
				<h1 class="truncate font-semibold text-white">{event.data.name}</h1>
				<p class="truncate text-xs">
					{formatDays(event.data.days, event.data.startsOn)} · {$t(helperRoleLabel(helper.data.role))}: {helper.data.name}{only ? ` · ${rangeName} ${helper.data.lanes.join(', ')}` : ''}
				</p>
			</div>
			<div class="ml-auto flex items-center gap-2">
				<span class="size-2.5 rounded-full {live?.connected ? 'bg-emerald-400' : 'bg-red-400'}" title={$t(liveLabel(live?.connected))}></span>
				<ThemeToggle />
			</div>
		</header>
		<OfflineBanner />
		<div class="flex flex-1 flex-col p-3 md:p-6">
			<EventDesk {tenantId} {eventId} {base} can={helperRights(helper.data.role)} {only} />
		</div>
	{:else}
		<p class="p-10 text-muted-foreground">{$t`Loading …`}</p>
	{/if}
</main>
