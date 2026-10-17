<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import createClient from 'openapi-fetch';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { call } from '$lib/api';
	import type { LaneState } from '$lib/api';
	import type { paths } from '$lib/api.gen';
	import { formatDays, laneSections } from '$lib/format';
	import { intlLocale, localeNames, setLocale, type Locale } from '$lib/locale.svelte';
	import { liveLabel } from '$lib/labels';
	import { Live, refresh } from '$lib/live.svelte';
	import Logo from '$lib/Logo.svelte';
	import RankingTable from '$lib/RankingTable.svelte';
	import Maximize from '@lucide/svelte/icons/maximize';
	import { Button } from '$lib/components/ui/button';
	import ThemeToggle from '$lib/ThemeToggle.svelte';
	import { setMode } from 'mode-watcher';

	const token = $derived(page.params.token ?? '');
	const kiosk = $derived(createClient<paths>({ baseUrl: '', headers: { Authorization: `Bearer ${token}` } }));
	const queryClient = useQueryClient();

	const theme = page.url.searchParams.get('theme');
	if (theme === 'dark' || theme === 'light') setMode(theme);
	const language = page.url.searchParams.get('language');
	const rangeFilter = Number(page.url.searchParams.get('range')) || 0;

	const event = createQuery(() => ({
		queryKey: ['kiosk', token, 'event'],
		queryFn: async () => {
			const ev = await call(kiosk.GET('/api/kiosk/event'));
			await setLocale(language && language in localeNames ? (language as Locale) : ev.locale, false);
			return ev;
		},
		retry: false
	}));
	const shownRanges = $derived((event.data?.ranges ?? []).filter((r) => !rangeFilter || r.id === rangeFilter));
	const rankings = createQuery(() => ({
		queryKey: ['kiosk', token, 'rankings'],
		queryFn: () => call(kiosk.GET('/api/kiosk/rankings')),
		enabled: !!event.data
	}));
	const lanes = createQuery(() => ({
		queryKey: ['kiosk', token, 'lanes'],
		queryFn: () => call(kiosk.GET('/api/kiosk/lanes')),
		enabled: !!event.data?.showLaneStatus
	}));

	let live = $state<Live>();
	$effect(() => {
		if (!event.data) return;
		const l = new Live(`/api/kiosk/live?token=${encodeURIComponent(token)}`, (topic) => {
			if (topic === 'lanes') refresh(queryClient, ['kiosk', token, 'lanes']);
			else refresh(queryClient, ['kiosk', token]);
		});
		live = l;
		return () => l.close();
	});

	let now = $state(new Date());
	$effect(() => {
		const timer = setInterval(() => (now = new Date()), 1000);
		return () => clearInterval(timer);
	});

	/** Toggles fullscreen. */
	function fullscreen() {
		if (document.fullscreenElement) document.exitFullscreen();
		else document.documentElement.requestFullscreen();
	}

	/** Text shown on a lane: shooter, reservation, out of service or free. */
	const laneText = (l: LaneState) =>
		l.state === 'occupied' ? (l.shooterName ?? '') : l.state === 'disabled' ? $t`out of service` : l.state === 'reserved' ? (l.reservation ?? '') : $t`free`;

	const laneClass: Record<string, string> = { free: 'lane-free', occupied: 'lane-occupied', reserved: 'lane-reserved', disabled: 'lane-disabled' };

	/** Rankings grouped by range when the event has several; programmes without a range come last. */
	const rankingGroups = $derived.by(() => {
		const list = (rankings.data ?? []).filter((r) => !rangeFilter || !r.rangeId || r.rangeId === rangeFilter);
		if (shownRanges.length < 2) return [{ name: '', items: list }];
		const ids = [...shownRanges.map((r) => r.id), 0];
		return ids
			.map((id) => ({ name: id ? (shownRanges.find((r) => r.id === id)?.name ?? '') : $t`All ranges`, items: list.filter((r) => (r.rangeId ?? 0) === id) }))
			.filter((g) => g.items.length > 0);
	});
</script>

<svelte:head><title>{event.data ? $t`${event.data.name} – Rankings` : $t`Rankings`}</title></svelte:head>

<main class="min-h-svh">
	{#if event.error}
		<p class="p-10 text-xl text-destructive" role="alert">{$t`This kiosk link is not valid (anymore).`}</p>
	{:else if event.data}
		<header class="sticky top-0 z-10 flex flex-wrap items-center gap-x-6 gap-y-2 bg-sidebar px-6 pt-[calc(1rem+env(safe-area-inset-top))] pb-4 text-sidebar-foreground">
			<Logo class="size-11 text-sidebar-primary" />
			<div class="min-w-0">
				<h1 class="truncate text-3xl font-semibold text-white">{event.data.name}</h1>
				<p class="text-lg">{formatDays(event.data.days, event.data.startsOn)} · {shownRanges.map((r) => r.name).join(', ')}{event.data.organizer ? ` · ${event.data.organizer}` : ''}</p>
			</div>
			<div class="ml-auto flex items-center gap-4">
				<span class="text-4xl font-semibold text-white tabular-nums">{now.toLocaleTimeString(intlLocale())}</span>
				<span class="size-3 rounded-full {live?.connected ? 'bg-emerald-400' : 'bg-red-400'}" title={$t(liveLabel(live?.connected))}></span>
				<div class="flex gap-1 opacity-40 transition-opacity hover:opacity-100">
					<ThemeToggle />
					<Button variant="ghost" size="icon-sm" title={$t`Fullscreen`} aria-label={$t`Fullscreen`} onclick={fullscreen}><Maximize /></Button>
				</div>
			</div>
		</header>

		<div class="space-y-6 p-6">
			{#if event.data.showLaneStatus && lanes.data}
				{#each lanes.data.filter((rl) => !rangeFilter || rl.rangeId === rangeFilter) as rl (rl.rangeId)}
					<section aria-label={$t`Lane status`}>
						{#if shownRanges.length > 1}<h2 class="mb-2 text-xl font-semibold">{rl.rangeName}</h2>{/if}
						{#each laneSections(rl.lanes, rl.groups) as sec (sec.name)}
						{#if sec.name}<h3 class="mt-2 mb-1 text-lg font-medium text-muted-foreground">{sec.name}</h3>{/if}
						<div class="grid grid-cols-[repeat(auto-fill,minmax(8rem,1fr))] gap-2">
							{#each sec.lanes as l (l.lane)}
								<div class="lane {laneClass[l.state]} p-2.5" aria-label={$t`Lane ${l.lane}: ${laneText(l)}`}>
									<span class="text-2xl font-bold">{l.lane}</span>
									<span class="truncate text-sm">
										{laneText(l)}
									</span>
								</div>
							{/each}
						</div>
						{/each}
					</section>
				{/each}
			{/if}

			{#each rankingGroups as g (g.name)}
				{#if g.name}<h2 class="text-xl font-semibold">{g.name}</h2>{/if}
				<div class="grid gap-6 *:min-w-0 2xl:grid-cols-2">
					{#each g.items as r (r.programmeId)}
						<RankingTable ranking={r} large />
					{/each}
				</div>
			{/each}
		</div>
	{:else}
		<p class="p-10 text-muted-foreground">{$t`Loading …`}</p>
	{/if}
</main>
