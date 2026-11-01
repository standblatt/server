<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Copy from '@lucide/svelte/icons/copy';
	import Power from '@lucide/svelte/icons/power';
	import RefreshCw from '@lucide/svelte/icons/refresh-cw';
	import ExternalLink from '@lucide/svelte/icons/external-link';
	import { toast } from 'svelte-sonner';
	import { api, call, type BookingConfig, type Event } from '$lib/api';
	import { byDay, deadlineDay, slotDay, slotTime } from './booking';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { ask } from '$lib/confirm.svelte';
	import BookingHoursForm from './parts/BookingHoursForm.svelte';
	import { keys } from '$lib/tenant';

	let { tenantId, event }: { tenantId: number; event: Event } = $props();

	const queryClient = useQueryClient();
	const key = $derived([...keys.booking(tenantId, event.id), 'config']);
	const params = $derived({ params: { path: { tenantId, eventId: event.id } } });
	const config = createQuery(() => ({
		queryKey: key,
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/booking', params))
	}));

	/** Puts a fresh answer of the server into the cache. */
	const store = (c: BookingConfig) => queryClient.setQueryData(key, c);

	/** Switches the booking on, or gives it a new link after confirmation. */
	async function newCode(renew: boolean) {
		if (renew && !(await ask({ title: $t`New link?`, description: $t`The current link stops working. Clubs that already asked keep their personal links.`, action: $t`New link` }))) return;
		store(await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/booking/code', params)));
		toast.success(renew ? $t`New link created` : $t`Online booking switched on`);
	}

	/** Switches the booking off after confirmation; requests and reservations stay. */
	async function off() {
		if (!(await ask({ title: $t`Switch online booking off?`, description: $t`The link stops working. Requests and confirmed reservations stay.`, action: $t`Switch off` }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/events/{eventId}/booking/code', params));
		await queryClient.invalidateQueries({ queryKey: key });
		toast.success($t`Online booking switched off`);
	}

	/** Copies the public link. */
	function copy(link: string) {
		navigator.clipboard.writeText(link);
		toast.success($t`Link copied`);
	}

	/** One line per day: first start, last end, number of slots and the slot starts with a gap mark before each break. */
	const preview = $derived(
		byDay(config.data?.slots ?? [], event.timezone).map(([day, slots]) => ({
			day,
			from: slotTime(slots[0].start, event.timezone),
			to: slotTime(slots[slots.length - 1].end, event.timezone),
			count: slots.length,
			starts: slots.map((sl, i) => ({ time: slotTime(sl.start, event.timezone), gap: i > 0 && new Date(sl.start).getTime() !== new Date(slots[i - 1].end).getTime() }))
		}))
	);

	/** Where the booking window stands today, as text and badge colours. */
	function windowState(c: BookingConfig) {
		if (c.opens && new Date(c.opens) > new Date())
			return { text: $t`Opens on ${slotDay(c.opens, event.timezone)}`, tone: 'border-amber-500/50 bg-amber-500/10 text-amber-700 dark:text-amber-400' };
		if (new Date(c.deadline) > new Date())
			return { text: $t`Open until ${deadlineDay(c.deadline, event.timezone)}`, tone: 'border-emerald-500/50 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400' };
		return { text: $t`Deadline passed`, tone: 'text-muted-foreground' };
	}
</script>

{#if config.data}
	{@const c = config.data}
	<div class="grid gap-4">
		{#if !c.enabled}
			<section class="flex flex-wrap items-center gap-3 rounded-lg border bg-card p-4">
				<div class="min-w-48 flex-1 text-sm">
					<div class="font-medium">{$t`Online booking`}</div>
					<div class="text-muted-foreground">{$t`External clubs ask for lanes through a link in your invitation and announce their shooters. You confirm each request.`}</div>
				</div>
				<Button onclick={() => newCode(false)}><Power />{$t`Switch on`}</Button>
			</section>
		{:else}
			{@const win = windowState(c)}
			<section class="grid gap-3 rounded-lg border bg-card p-4">
				<div class="flex flex-wrap items-center gap-2 text-sm font-medium">{$t`Link for the invitation`}<Badge variant="outline" class={win.tone}>{win.text}</Badge></div>
				<div class="flex flex-wrap items-start gap-4">
					<img src={c.qr} alt={$t`QR code of the booking link`} class="size-28 rounded border bg-white p-1" />
					<div class="grid min-w-0 flex-1 gap-2">
						<div class="flex gap-2">
							<Input readonly value={c.link} aria-label={$t`Booking link`} />
							<Button variant="outline" size="icon" aria-label={$t`Copy link`} onclick={() => copy(c.link)}><Copy /></Button>
							<Button variant="outline" size="icon" aria-label={$t`Open link`} onclick={() => window.open(c.link, '_blank')}><ExternalLink /></Button>
						</div>
						<div class="flex flex-wrap gap-2">
							<Button variant="outline" size="sm" onclick={() => newCode(true)}><RefreshCw />{$t`New link`}</Button>
							<Button variant="outline" size="sm" onclick={off}><Power />{$t`Switch off`}</Button>
						</div>
					</div>
				</div>
			</section>
		{/if}

		<section class="grid gap-2 rounded-lg border bg-card p-4 text-sm">
			<div class="font-medium">{$t`Slots`}</div>
			{#each preview as p (p.day)}
				<div class="grid gap-1.5">
					<div class="flex flex-wrap gap-x-3"><span class="w-56">{p.day}</span><span class="tabular-nums">{p.from}–{p.to}</span><span class="text-muted-foreground">{$plural(p.count, { one: '# slot', other: '# slots' })}</span></div>
					<div class="flex flex-wrap gap-1">
						{#each p.starts as s (s.time)}
							{#if s.gap}<span class="px-1 text-xs text-muted-foreground">{$t`Break`}</span>{/if}
							<span class="rounded border px-1.5 text-xs tabular-nums">{s.time}</span>
						{/each}
					</div>
				</div>
			{:else}
				<p class="text-muted-foreground">{$t`No slots: check the hours.`}</p>
			{/each}
		</section>

		{#key c.settings}<BookingHoursForm {tenantId} {event} settings={c.settings} onsaved={store} />{/key}
	</div>
{/if}
