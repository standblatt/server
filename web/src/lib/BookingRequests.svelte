<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Check from '@lucide/svelte/icons/check';
	import UserPlus from '@lucide/svelte/icons/user-plus';
	import X from '@lucide/svelte/icons/x';
	import { toast } from 'svelte-sonner';
	import { api, call, type BookingRequest, type Event } from './api';
	import { slotDay, slotTime } from './booking';
	import { formatLanes } from './format';
	import BookingStatus from './booking/BookingStatus.svelte';
	import { Button } from './components/ui/button';
	import * as Dialog from './components/ui/dialog';
	import { Textarea } from './components/ui/textarea';
	import FormError from './FormError.svelte';
	import { bookingRequestsQuery, keys } from './tenant';

	let { tenantId, event, isAdmin }: { tenantId: number; event: Event; isAdmin: boolean } = $props();

	const queryClient = useQueryClient();
	const requests = bookingRequestsQuery(
		() => tenantId,
		() => event.id
	);
	const order = { pending: 0, confirmed: 1, declined: 2, cancelled: 3 };
	const sorted = $derived([...(requests.data ?? [])].sort((a, b) => order[a.status] - order[b.status]));
	/** Name of a range of the event. */
	const rangeName = (id: number) => event.ranges.find((r) => r.id === id)?.name ?? '';
	const tz = $derived(event.timezone);

	/** Reloads requests, reservations and lanes after a decision. */
	async function reload() {
		await queryClient.invalidateQueries({ queryKey: keys.booking(tenantId, event.id) });
		await queryClient.invalidateQueries({ queryKey: keys.reservations(tenantId, event.id) });
		await queryClient.invalidateQueries({ queryKey: keys.lanes(tenantId, event.id) });
	}

	/** Runs a decision on a request and shows its result or error as a toast. */
	async function decide(r: BookingRequest, action: 'confirm' | 'shooters', done: string) {
		try {
			const params = { params: { path: { tenantId, requestId: r.id } } };
			if (action === 'confirm') await call(api.POST('/api/tenants/{tenantId}/booking-requests/{requestId}/confirm', params));
			else await call(api.POST('/api/tenants/{tenantId}/booking-requests/{requestId}/shooters', params));
			toast.success(done);
			await reload();
		} catch (err) {
			toast.error((err as Error).message);
		}
	}

	let declining = $state<BookingRequest | null>(null);
	let reason = $state('');
	let error = $state('');

	/** Declines the request in the dialog with the reason typed there. */
	async function decline(e: SubmitEvent) {
		e.preventDefault();
		if (!declining) return;
		error = '';
		try {
			await call(api.POST('/api/tenants/{tenantId}/booking-requests/{requestId}/decline', { params: { path: { tenantId, requestId: declining.id } }, body: { reason } }));
			toast.success($t`Request declined`);
			declining = null;
			reason = '';
			await reload();
		} catch (err) {
			error = (err as Error).message;
		}
	}
</script>

{#if sorted.length > 0}
	<section class="grid gap-3" aria-label={$t`Booking requests`}>
		<h2 class="text-base font-semibold">{$t`Booking requests`}</h2>
		{#each sorted as r (r.id)}
			{@const waiting = r.shooters.filter((s) => !s.shooterId).length}
			<article class="grid gap-3 rounded-lg border bg-card p-4 text-sm" class:opacity-60={r.status === 'declined' || r.status === 'cancelled'} aria-label={$t`Request of ${r.club}`}>
				<div class="flex flex-wrap items-center gap-3">
					<span class="text-base font-semibold">{r.club}</span>
					<BookingStatus status={r.status} />
					<span class="text-muted-foreground">{r.contact} · <a class="underline" href="mailto:{r.email}">{r.email}</a>{r.phone ? ` · ${r.phone}` : ''}</span>
					{#if isAdmin}
						<div class="ml-auto flex flex-wrap gap-2">
							{#if r.status === 'pending'}
								<Button size="sm" onclick={() => decide(r, 'confirm', $t`Request confirmed`)}><Check />{$t`Confirm`}</Button>
							{/if}
							{#if r.status === 'confirmed' && waiting > 0}
								<Button size="sm" variant="outline" onclick={() => decide(r, 'shooters', $t`Shooters taken over`)}><UserPlus />{$plural(waiting, { one: 'Take over # shooter', other: 'Take over # shooters' })}</Button>
							{/if}
							{#if r.status === 'pending' || r.status === 'confirmed'}
								<Button size="sm" variant="outline" onclick={() => ((declining = r), (error = ''))}><X />{$t`Decline`}</Button>
							{/if}
						</div>
					{/if}
				</div>
				{#if r.note}<p class="rounded-md bg-muted px-3 py-2">{r.note}</p>{/if}
				<ul class="grid gap-1">
					{#each r.items as it (`${it.rangeId}-${it.start}`)}
						<li class="flex flex-wrap gap-x-3">
							<span class="w-64">{slotDay(it.start, tz)}</span>
							<span class="w-28 tabular-nums">{slotTime(it.start, tz)}–{slotTime(it.end, tz)}</span>
							{#if event.ranges.length > 1}<span>{rangeName(it.rangeId)}</span>{/if}
							<span class="font-medium">{$plural(it.lanes, { one: '# lane', other: '# lanes' })}</span>
							{#if it.assigned.length > 0}<span class="text-muted-foreground">{$t`lanes ${formatLanes(it.assigned)}`}</span>{/if}
						</li>
					{/each}
				</ul>
				{#if r.shooters.length > 0}
					<details>
						<summary class="cursor-pointer">{$plural(r.shooters.length, { one: '# shooter announced', other: '# shooters announced' })}</summary>
						<ul class="mt-2 grid gap-1 sm:grid-cols-2">
							{#each r.shooters as s, i (i)}
								<li class="flex gap-2">
									<span>{s.lastname} {s.firstname}</span>
									<span class="text-muted-foreground">{s.birthday.slice(0, 4)}{s.licenseNumber ? ` · ${s.licenseNumber}` : ''}</span>
									{#if s.shooterId}<span class="text-emerald-700 dark:text-emerald-400">{$t`in master data`}</span>
									{:else if s.known}<span class="text-muted-foreground">{$t`known`}</span>
									{:else}<span class="text-sky-700 dark:text-sky-400">{$t`new`}</span>{/if}
								</li>
							{/each}
						</ul>
					</details>
				{/if}
				{#if r.status === 'declined' && r.declineReason}<p class="text-muted-foreground">{$t`Reason: ${r.declineReason}`}</p>{/if}
			</article>
		{/each}
	</section>
{/if}

<Dialog.Root open={declining !== null} onOpenChange={(o) => !o && (declining = null)}>
	<Dialog.Content class="sm:max-w-md">
		{#if declining}
			<Dialog.Header>
				<Dialog.Title>{$t`Decline request of ${declining.club}?`}</Dialog.Title>
				<Dialog.Description>{$t`The club gets a mail with your reason. Lanes already reserved for it are free again.`}</Dialog.Description>
			</Dialog.Header>
			<form class="grid gap-4" onsubmit={decline}>
				<Textarea bind:value={reason} rows={3} maxlength={1000} aria-label={$t`Reason`} placeholder={$t`Reason (optional)`} />
				<FormError message={error} />
				<Dialog.Footer><Button type="submit" variant="destructive">{$t`Decline`}</Button></Dialog.Footer>
			</form>
		{/if}
	</Dialog.Content>
</Dialog.Root>
