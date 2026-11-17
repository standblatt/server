<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { intlLocale } from '#lib/locale.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Copy from '@lucide/svelte/icons/copy';
	import MoreHorizontal from '@lucide/svelte/icons/ellipsis';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { formatLanes, fromZonedInput, zonedInput } from '#lib/format';
	import { api, call, type BookingRequest, type ReservationView } from '#lib/api';
	import { Badge } from '#lib/components/ui/badge';
	import { Button } from '#lib/components/ui/button';
	import * as Dialog from '#lib/components/ui/dialog';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu';
	import { Input } from '#lib/components/ui/input';
	import * as NativeSelect from '#lib/components/ui/native-select';
	import { Switch } from '#lib/components/ui/switch';
	import { ask } from '#lib/confirm.svelte';
	import FormError from '#lib/FormError.svelte';
	import FormField from '#lib/FormField.svelte';
	import PageHeader from '#lib/PageHeader.svelte';
	import { bookingRequestsQuery, keys } from '#lib/tenant';
	import Info from '@lucide/svelte/icons/info';
	import Check from '@lucide/svelte/icons/check';
	import X from '@lucide/svelte/icons/x';
	import BookingStatus from './booking/parts/BookingStatus.svelte';
	import RequestSheet from './booking/RequestSheet.svelte';
	import { byDay } from './booking/booking';
	import LaneStrip from './LaneStrip.svelte';
	import type { EventProps } from '..';
	import { can, membershipOf, meQuery } from '#lib/me';

	let { tenantId, event }: EventProps = $props();
	const eventId = $derived(event.id);
	const queryClient = useQueryClient();
	const ranges = $derived(event.ranges);
	const tz = $derived(event.timezone);
	const me = meQuery();
	const isAdmin = $derived(can(membershipOf(me.data, tenantId), 'events.configure'));
	let showPast = $state(false);
	const key = $derived([...keys.reservations(tenantId, eventId), showPast, ranges.map((r) => r.id).join(',')]);
	const reservations = createQuery(() => ({
		queryKey: key,
		queryFn: async () => {
			const lists = await Promise.all(
				ranges.map((r) => call(api.GET('/api/tenants/{tenantId}/ranges/{rangeId}/reservations', { params: { path: { tenantId, rangeId: r.id }, query: { all: showPast } } })))
			);
			return lists.flatMap((list, i) => list.map((r) => ({ ...r, rangeName: ranges[i].name })));
		},
		enabled: ranges.length > 0
	}));
	type Row = ReservationView & { rangeName: string };

	/** The next full quarter hour from now. */
	function nextQuarter() {
		const d = new Date();
		d.setSeconds(0, 0);
		d.setMinutes(Math.ceil((d.getMinutes() + 1) / 15) * 15);
		return d;
	}

	type Form = { id: number; rangeId: number; lanes: string; start: string; end: string; lead: number; description: string };
	let form = $state<Form | null>(null);
	let error = $state('');

	/** Opens the form to edit or copy a reservation, or for a new one starting at the next quarter hour. */
	function open(r?: ReservationView, copy = false) {
		error = '';
		if (r) {
			form = { id: copy ? 0 : r.id, rangeId: r.rangeId, lanes: r.lanesText, start: zonedInput(new Date(r.startsAt), tz), end: zonedInput(new Date(r.endsAt), tz), lead: r.leadMinutes, description: r.description };
		} else {
			const start = nextQuarter();
			form = { id: 0, rangeId: ranges[0]?.id ?? 0, lanes: '', start: zonedInput(start, tz), end: zonedInput(new Date(start.getTime() + 3600000), tz), lead: 30, description: '' };
		}
	}

	/** Sets the end one hour after the new start. */
	function startChanged() {
		if (form) form.end = zonedInput(new Date(fromZonedInput(form.start, tz).getTime() + 3600000), tz);
	}

	/** Creates or updates the reservation from the form. */
	async function save(e: SubmitEvent) {
		e.preventDefault();
		if (!form) return;
		error = '';
		const body = {
			lanes: form.lanes,
			startsAt: fromZonedInput(form.start, tz).toISOString(),
			endsAt: fromZonedInput(form.end, tz).toISOString(),
			leadMinutes: Number(form.lead),
			description: form.description
		};
		try {
			if (form.id) await call(api.PUT('/api/tenants/{tenantId}/reservations/{reservationId}', { params: { path: { tenantId, reservationId: form.id } }, body }));
			else await call(api.POST('/api/tenants/{tenantId}/ranges/{rangeId}/reservations', { params: { path: { tenantId, rangeId: form.rangeId } }, body }));
			toast.success(form.id ? $t`Reservation saved` : $t`Reservation created`);
			form = null;
			await queryClient.invalidateQueries({ queryKey: keys.reservations(tenantId, eventId) });
			await queryClient.invalidateQueries({ queryKey: keys.lanes(tenantId, eventId) });
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Deletes a reservation after confirmation. */
	async function remove(r: ReservationView) {
		if (!(await ask({ title: $t`Delete reservation?`, description: $t`"${r.description}" is removed and its lanes are free again.`, action: $t`Delete`, destructive: true }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/reservations/{reservationId}', { params: { path: { tenantId, reservationId: r.id } } }));
		await queryClient.invalidateQueries({ queryKey: keys.reservations(tenantId, eventId) });
		await queryClient.invalidateQueries({ queryKey: keys.lanes(tenantId, eventId) });
	}

	/** Whether a reservation is planned, active (lead time reached) or finished. */
	function status(r: ReservationView) {
		const now = Date.now();
		if (new Date(r.endsAt).getTime() <= now) return 'finished';
		if (new Date(r.startsAt).getTime() - r.leadMinutes * 60000 <= now) return 'active';
		return 'planned';
	}
	/** Translated label of a reservation status. */
	const statusLabel = (st: string) => (st === 'finished' ? $t`finished` : st === 'active' ? $t`active` : $t`planned`);
	/** Formats a timestamp as hours and minutes in the event's time zone. */
	const time = (iso: string) => new Date(iso).toLocaleTimeString(intlLocale(), { hour: '2-digit', minute: '2-digit', timeZone: tz });
	const requests = bookingRequestsQuery(
		() => tenantId,
		() => eventId
	);
	type BookingItem = BookingRequest['items'][number];
	type Line = { key: string; start: string; end: string; rangeId: number; reservation?: Row; request?: BookingRequest; item?: BookingItem };
	/** Reservations and the slots of open requests, grouped by day; confirmed requests show through their reservations. */
	const days = $derived(
		byDay<Line>(
			[
				...(reservations.data ?? []).map((r) => ({ key: `r${r.id}`, start: r.startsAt, end: r.endsAt, rangeId: r.rangeId, reservation: r })),
				...(requests.data ?? [])
					.filter((q) => q.status === 'pending')
					.flatMap((q) => q.items.map((it, i) => ({ key: `q${q.id}-${i}`, start: it.start, end: it.end, rangeId: it.rangeId, request: q, item: it })))
			],
			tz
		)
	);
	const closed = $derived((requests.data ?? []).filter((q) => q.status === 'declined' || q.status === 'cancelled'));
	/** The range of the event with that ID. */
	const rangeOf = (id: number) => ranges.find((r) => r.id === id);
	let requestId = $state<number | null>(null);
	let sheet: RequestSheet;
</script>

<PageHeader help="user-guide/reservations" title={$t`Reservations`} description={$t`Lanes reserved for trainings or groups; they show as reserved from the lead time on.`}>
	{#snippet actions()}
		<label class="flex items-center gap-2 text-sm text-muted-foreground"><Switch bind:checked={showPast} />{$t`show past`}</label>
		<Button onclick={() => open()}><Plus />{$t`New reservation`}</Button>
	{/snippet}
</PageHeader>

<div class="space-y-4">
	{#each days as [day, list] (day)}
		<section class="overflow-hidden rounded-lg border bg-card">
			<div class="border-b px-4 py-2.5"><h2 class="text-sm font-semibold">{day}</h2></div>
			<ul class="divide-y">
				{#each list as l (l.key)}
					{@const range = rangeOf(l.rangeId)}
					{#if l.reservation}
						{@const r = l.reservation}
						{@const st = status(r)}
						<li class="flex flex-wrap items-center gap-x-4 gap-y-1 px-4 py-2.5 text-sm" class:opacity-50={st === 'finished'} aria-label={$t`Reservation ${r.description}`}>
							<Badge
								variant="outline"
								class="min-w-20 justify-center {st === 'active' ? 'border-amber-500/50 bg-amber-500/10 text-amber-700 dark:text-amber-400' : st === 'planned' ? 'border-sky-500/50 text-sky-700 dark:text-sky-400' : ''}">{statusLabel(st)}</Badge
							>
							<span class="w-28 font-medium tabular-nums">{time(r.startsAt)}–{time(r.endsAt)}</span>
							<span class="flex min-w-0 items-baseline gap-2 sm:w-64">
								{#if r.bookingRequestId}
									{@const qid = r.bookingRequestId}
									<button type="button" class="truncate font-semibold hover:underline" onclick={() => (requestId = qid)}>{r.description}</button>
									<span class="shrink-0 text-xs text-muted-foreground">{$t`online booking`}</span>
								{:else}
									<span class="truncate font-semibold">{r.description}</span>
								{/if}
							</span>
							{#if range}<LaneStrip count={range.laneCount} lanes={r.lanes} disabled={range.disabledLanes} />{/if}
							<span class="text-muted-foreground">{ranges.length > 1 ? `${r.rangeName} · ` : ''}{$t`lanes ${r.lanesText}`}</span>
							{#if r.leadMinutes}<span class="text-xs text-muted-foreground">{$t`shown ${r.leadMinutes} min before`}</span>{/if}
							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									{#snippet child({ props })}<Button variant="ghost" size="icon-sm" class="ml-auto" aria-label={$t`Actions for ${r.description}`} {...props}><MoreHorizontal /></Button>{/snippet}
								</DropdownMenu.Trigger>
								<DropdownMenu.Content align="end">
									{#if r.bookingRequestId}
										{@const qid = r.bookingRequestId}
										<DropdownMenu.Item onSelect={() => (requestId = qid)}><Info />{$t`Request details`}</DropdownMenu.Item>
										<DropdownMenu.Separator />
									{/if}
									<DropdownMenu.Item onSelect={() => open(r)}><Pencil />{$t`Edit`}</DropdownMenu.Item>
									<DropdownMenu.Item onSelect={() => open(r, true)}><Copy />{$t`Copy`}</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item variant="destructive" onSelect={() => remove(r)}><Trash />{$t`Delete`}</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</li>
					{:else if l.request && l.item}
						{@const q = l.request}
						{@const it = l.item}
						<li class="flex flex-wrap items-center gap-x-4 gap-y-1 bg-amber-500/5 px-4 py-2.5 text-sm" aria-label={$t`Request of ${q.club}`}>
							<Badge variant="outline" class="min-w-20 justify-center border-amber-500/50 bg-amber-500/10 text-amber-700 dark:text-amber-400">{$t`requested`}</Badge>
							<span class="w-28 font-medium tabular-nums">{time(it.start)}–{time(it.end)}</span>
							<span class="flex min-w-0 items-baseline gap-2 sm:w-64"><button type="button" class="truncate font-semibold hover:underline" onclick={() => (requestId = q.id)}>{q.club}</button></span>
							{#if range}<LaneStrip count={range.laneCount} lanes={it.wanted ?? []} disabled={range.disabledLanes} pending />{/if}
							<span class="text-muted-foreground">
								{ranges.length > 1 ? `${range?.name ?? ''} · ` : ''}{$plural(it.lanes, { one: '# lane', other: '# lanes' })}{#if it.wanted?.length}{' · '}{$t`wanted ${formatLanes(it.wanted)}`}{/if}
							</span>
							{#if isAdmin}
								<div class="ml-auto flex gap-2">
									<Button size="sm" onclick={() => sheet.openConfirm(q)}><Check />{$t`Confirm`}</Button>
									<Button size="sm" variant="outline" onclick={() => sheet.openDecline(q)}><X />{$t`Decline`}</Button>
								</div>
							{/if}
						</li>
					{/if}
				{/each}
			</ul>
		</section>
	{:else}
		<div class="rounded-lg border border-dashed py-12 text-center text-sm text-muted-foreground">{$t`No reservations.`}</div>
	{/each}
	{#if closed.length > 0}
		<details class="text-sm">
			<summary class="cursor-pointer text-muted-foreground">{$plural(closed.length, { one: '# declined or cancelled request', other: '# declined or cancelled requests' })}</summary>
			<ul class="mt-2 divide-y rounded-lg border bg-card">
				{#each closed as q (q.id)}
					<li class="flex flex-wrap items-center gap-x-4 gap-y-1 px-4 py-2">
						<BookingStatus status={q.status} />
						<button type="button" class="font-semibold hover:underline" onclick={() => (requestId = q.id)}>{q.club}</button>
						{#if q.declineReason}<span class="text-muted-foreground">{q.declineReason}</span>{/if}
					</li>
				{/each}
			</ul>
		</details>
	{/if}
</div>

<RequestSheet bind:this={sheet} bind:requestId {tenantId} {event} {isAdmin} />

<Dialog.Root open={form !== null} onOpenChange={(o) => !o && (form = null)}>
	<Dialog.Content class="sm:max-w-lg">
		{#if form}
			<Dialog.Header><Dialog.Title>{form.id ? $t`Edit reservation` : $t`New reservation`}</Dialog.Title></Dialog.Header>
			<form class="grid gap-4 sm:grid-cols-2" onsubmit={save}>
				<FormField label={$t`Description`} class="sm:col-span-2"><Input bind:value={form.description} required /></FormField>
				{#if ranges.length > 1}
					<FormField label={$t`Range`} class="sm:col-span-2">
						<NativeSelect.Root class="w-full" bind:value={form.rangeId} disabled={form.id > 0}>
							{#each ranges as r (r.id)}<NativeSelect.Option value={r.id}>{r.name}</NativeSelect.Option>{/each}
						</NativeSelect.Root>
					</FormField>
				{/if}
				<FormField label={$t`Lanes`} class="sm:col-span-2" hint={$t`e.g. 5 or 1, 5, 7-10`}><Input bind:value={form.lanes} placeholder={$t`e.g. 5 or 1, 5, 7-10`} required /></FormField>
				<FormField label={$t`From`}><Input type="datetime-local" bind:value={form.start} onchange={startChanged} required /></FormField>
				<FormField label={$t`To`}><Input type="datetime-local" bind:value={form.end} required /></FormField>
				<FormField label={$t`Show as reserved (minutes before)`} class="sm:col-span-2"><Input class="w-28" type="number" min="0" max="1440" bind:value={form.lead} /></FormField>
				<FormError message={error} class="sm:col-span-2" />
				<Dialog.Footer class="sm:col-span-2"><Button type="submit">{$t`Save`}</Button></Dialog.Footer>
			</form>
		{/if}
	</Dialog.Content>
</Dialog.Root>
