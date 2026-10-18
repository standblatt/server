<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery } from '@tanstack/svelte-query';
	import CircleCheck from '@lucide/svelte/icons/circle-check';
	import Circle from '@lucide/svelte/icons/circle';
	import { api, call } from './api';
	import { Button } from './components/ui/button';
	import { eventsQuery, keys, rangesQuery } from './tenant';

	let { tenantId, eventId, isAdmin, onnewevent }: { tenantId: number; eventId?: number; isAdmin: boolean; onnewevent?: () => void } = $props();

	const ranges = rangesQuery(() => tenantId);
	const events = eventsQuery(() => tenantId);
	const programmes = createQuery(() => ({
		queryKey: keys.programmes(tenantId, eventId ?? 0),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/programmes', { params: { path: { tenantId, eventId: eventId! } } })),
		enabled: !!eventId
	}));
	const shooters = createQuery(() => ({
		queryKey: [...keys.shooters(tenantId), 'any'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/shooters', { params: { path: { tenantId }, query: { q: '', limit: 1 } } }))
	}));

	const base = $derived(`/t/${tenantId}`);
	const steps = $derived([
		{
			title: $t`Add a range`,
			text: $t`The shooting range with its number of lanes.`,
			done: !!ranges.data?.length,
			action: { label: $t`Ranges`, href: `${base}/ranges` }
		},
		{
			title: $t`Create an event`,
			text: $t`The competition with its date, on one of your ranges.`,
			done: !!events.data?.length,
			action: onnewevent ? { label: $t`New event`, onclick: onnewevent } : undefined
		},
		{
			title: $t`Set up programmes`,
			text: $t`What shooters can register for: series, prices and awards. Copy them from an earlier event or the library.`,
			done: !!programmes.data?.some((p) => p.active),
			action: eventId ? { label: $t`Programmes`, href: `${base}/events/${eventId}/config/programmes` } : undefined
		},
		{
			title: $t`Import shooters (optional)`,
			text: $t`Import your members from a CSV file. Shooters can also be added during registration.`,
			done: !!shooters.data?.length,
			action: { label: $t`Shooters`, href: `${base}/shooters` }
		}
	]);
	const next = $derived(steps.findIndex((s) => !s.done));
</script>

<section class="rounded-lg border bg-card" aria-label={$t`Getting started`}>
	<div class="border-b px-4 py-3">
		<h2 class="text-sm font-semibold">{$t`Getting started`}</h2>
		<p class="text-sm text-muted-foreground">
			{isAdmin ? $t`A few steps before the first shooter can register.` : $t`An administrator of your club still has to finish these steps.`}
		</p>
	</div>
	<ol class="divide-y">
		{#each steps as s, i (s.title)}
			<li class="flex flex-wrap items-center gap-3 px-4 py-3 {s.done ? 'text-muted-foreground' : ''}" data-done={s.done}>
				{#if s.done}<CircleCheck class="size-5 shrink-0 text-emerald-600" />{:else}<Circle class="size-5 shrink-0 text-muted-foreground" />{/if}
				<div class="min-w-48 flex-1">
					<div class="font-medium {s.done ? 'line-through' : ''}">{i + 1}. {s.title}</div>
					{#if !s.done}<div class="text-sm text-muted-foreground">{s.text}</div>{/if}
				</div>
				{#if isAdmin && !s.done && s.action}
					{#if 'href' in s.action}
						<Button size="sm" variant={i === next ? 'default' : 'outline'} href={s.action.href}>{s.action.label}</Button>
					{:else}
						<Button size="sm" variant={i === next ? 'default' : 'outline'} onclick={s.action.onclick}>{s.action.label}</Button>
					{/if}
				{/if}
			</li>
		{/each}
	</ol>
</section>
