<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Copy from '@lucide/svelte/icons/copy';
	import Power from '@lucide/svelte/icons/power';
	import RefreshCw from '@lucide/svelte/icons/refresh-cw';
	import { toast } from 'svelte-sonner';
	import { api, call, type BookingConfig, type Event } from '$lib/api';
	import { byDay, deadlineDay, slotTime } from './booking';
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

	/** One line per day: first start, last end and number of slots. */
	const preview = $derived(
		byDay(config.data?.slots ?? [], event.timezone).map(([day, slots]) => ({
			day,
			from: slotTime(slots[0].start, event.timezone),
			to: slotTime(slots[slots.length - 1].end, event.timezone),
			count: slots.length
		}))
	);
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
			<section class="grid gap-3 rounded-lg border bg-card p-4">
				<div class="text-sm font-medium">{$t`Link for the invitation`}</div>
				<div class="flex flex-wrap items-start gap-4">
					<img src={c.qr} alt={$t`QR code of the booking link`} class="size-28 rounded border bg-white p-1" />
					<div class="grid min-w-0 flex-1 gap-2">
						<div class="flex gap-2">
							<Input readonly value={c.link} aria-label={$t`Booking link`} />
							<Button variant="outline" size="icon" aria-label={$t`Copy link`} onclick={() => copy(c.link)}><Copy /></Button>
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
				<div class="flex flex-wrap gap-x-3"><span class="w-56">{p.day}</span><span class="tabular-nums">{p.from}–{p.to}</span><span class="text-muted-foreground">{$plural(p.count, { one: '# slot', other: '# slots' })}</span></div>
			{:else}
				<p class="text-muted-foreground">{$t`No slots: check the hours.`}</p>
			{/each}
			<p class="text-muted-foreground">{$t`Requests and changes until ${deadlineDay(c.deadline, event.timezone)}.`}</p>
		</section>

		{#key c.settings}<BookingHoursForm {tenantId} {event} settings={c.settings} onsaved={store} />{/key}
	</div>
{/if}
