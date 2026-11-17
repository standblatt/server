<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Check from '@lucide/svelte/icons/check';
	import Pencil from '@lucide/svelte/icons/pencil';
	import UserPlus from '@lucide/svelte/icons/user-plus';
	import X from '@lucide/svelte/icons/x';
	import { toast } from 'svelte-sonner';
	import { api, call, type BookingRequest, type BookingShooter, type Club, type Event } from '#lib/api';
	import { slotDay, slotTime } from './booking';
	import { formatLanes } from '#lib/format';
	import { intlLocale } from '#lib/locale.svelte';
	import BookingStatus from './parts/BookingStatus.svelte';
	import ShooterList from './parts/ShooterList.svelte';
	import ClubSelect from '#lib/ClubSelect.svelte';
	import { Button } from '#lib/components/ui/button';
	import * as Dialog from '#lib/components/ui/dialog';
	import * as Sheet from '#lib/components/ui/sheet';
	import { Textarea } from '#lib/components/ui/textarea';
	import FormError from '#lib/FormError.svelte';
	import { bookingRequestsQuery, clubsQuery, keys } from '#lib/tenant';

	let { tenantId, event, isAdmin, requestId = $bindable(null) }: { tenantId: number; event: Event; isAdmin: boolean; requestId?: number | null } = $props();

	const queryClient = useQueryClient();
	const requests = bookingRequestsQuery(
		() => tenantId,
		() => event.id
	);
	const r = $derived(requests.data?.find((x) => x.id === requestId));
	/** Name of a range of the event. */
	const rangeName = (id: number) => event.ranges.find((x) => x.id === id)?.name ?? '';
	const tz = $derived(event.timezone);
	/** Date and time to the second in the event's time zone, so the order of requests is plain. */
	const moment = (iso: string) => new Date(iso).toLocaleString(intlLocale(), { dateStyle: 'medium', timeStyle: 'medium', timeZone: tz });
	const clubs = clubsQuery(() => tenantId);
	/** Name of a club of the master data. */
	const clubName = (id: number) => clubs.data?.find((c) => c.id === id)?.name ?? '';

	/** Reloads requests, reservations and lanes after a decision. */
	async function reload() {
		await queryClient.invalidateQueries({ queryKey: keys.booking(tenantId, event.id) });
		await queryClient.invalidateQueries({ queryKey: keys.reservations(tenantId, event.id) });
		await queryClient.invalidateQueries({ queryKey: keys.lanes(tenantId, event.id) });
	}

	/** Takes the shooters added after the confirmation into the master data. */
	async function adopt(r: BookingRequest) {
		try {
			await call(api.POST('/api/tenants/{tenantId}/booking-requests/{requestId}/shooters', { params: { path: { tenantId, requestId: r.id } } }));
			toast.success($t`Shooters taken over`);
			await reload();
		} catch (err) {
			toast.error((err as Error).message);
		}
	}

	let draft = $state<BookingShooter[] | null>(null);
	let shooterError = $state('');

	/** Opens the announced shooters of the request for editing. */
	function editShooters(r: BookingRequest) {
		shooterError = '';
		draft = r.shooters.map((s) => ({ firstname: s.firstname, lastname: s.lastname, birthday: s.birthday, licenseNumber: s.licenseNumber ?? '' }));
	}

	/** Saves the edited shooter list; changed shooters wait for adoption again. */
	async function saveShooters(e: SubmitEvent) {
		e.preventDefault();
		if (!draft || requestId === null) return;
		shooterError = '';
		try {
			await call(api.PUT('/api/tenants/{tenantId}/booking-requests/{requestId}/shooters', { params: { path: { tenantId, requestId } }, body: { shooters: draft } }));
			toast.success($t`Shooters saved`);
			draft = null;
			await reload();
		} catch (err) {
			shooterError = (err as Error).message;
		}
	}

	let confirming = $state<BookingRequest | null>(null);
	let suggested = $state<Club[]>([]);
	let clubChoice = $state('');

	/** Opens the confirmation with the club of the master data that matches the request, or the similar ones. */
	export async function openConfirm(r: BookingRequest) {
		error = '';
		const m = await call(api.GET('/api/tenants/{tenantId}/clubs/similar', { params: { path: { tenantId }, query: { name: r.club } } }));
		suggested = [...(m.match ? [m.match] : []), ...m.similar];
		clubChoice = m.match ? String(m.match.id) : m.similar.length ? '' : 'new';
		confirming = r;
	}

	/** Opens the dialog that declines the request. */
	export function openDecline(r: BookingRequest) {
		error = '';
		declining = r;
	}

	/** Confirms the request in the dialog with the chosen club. */
	async function confirm(e: SubmitEvent) {
		e.preventDefault();
		if (!confirming) return;
		error = '';
		try {
			const body = clubChoice === 'new' ? {} : { clubId: Number(clubChoice) };
			await call(api.POST('/api/tenants/{tenantId}/booking-requests/{requestId}/confirm', { params: { path: { tenantId, requestId: confirming.id } }, body }));
			toast.success($t`Request confirmed`);
			confirming = null;
			await reload();
			await queryClient.invalidateQueries({ queryKey: keys.clubs(tenantId) });
		} catch (err) {
			error = (err as Error).message;
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

<Sheet.Root open={!!r} onOpenChange={(o) => !o && ((requestId = null), (draft = null))}>
	<Sheet.Content class="gap-0 overflow-y-auto p-0 data-[side=right]:w-full data-[side=right]:sm:max-w-xl">
		{#if r}
			{@const waiting = r.shooters.filter((s) => !s.shooterId).length}
			{@const editable = isAdmin && (r.status === 'pending' || r.status === 'confirmed')}
			<Sheet.Header class="border-b px-6 py-4">
				<Sheet.Title class="text-lg">{r.club}</Sheet.Title>
				<Sheet.Description>{r.contact} · <a class="underline" href="mailto:{r.email}">{r.email}</a>{r.phone ? ` · ${r.phone}` : ''}</Sheet.Description>
			</Sheet.Header>
			<div class="grid gap-4 px-6 py-4 text-sm">
				<div class="flex flex-wrap items-center gap-2">
					<BookingStatus status={r.status} />
					<span class="text-muted-foreground">
						{$t`Received ${moment(r.createdAt)}`}{#if r.status === 'pending' && r.updatedAt !== r.createdAt}{' · '}{$t`changed ${moment(r.updatedAt)}`}{/if}
					</span>
				</div>
				{#if r.clubId && clubName(r.clubId) !== r.club}<p class="text-muted-foreground">{$t`Club in the master data: ${clubName(r.clubId)}`}</p>{/if}
				{#if r.note}<p class="rounded-md bg-muted px-3 py-2">{r.note}</p>{/if}
				{#if r.status === 'declined' && r.declineReason}<p class="text-muted-foreground">{$t`Reason: ${r.declineReason}`}</p>{/if}
				<ul class="grid gap-1">
					{#each r.items as it (`${it.rangeId}-${it.start}`)}
						<li class="flex flex-wrap gap-x-3">
							<span>{slotDay(it.start, tz)}</span>
							<span class="tabular-nums">{slotTime(it.start, tz)}–{slotTime(it.end, tz)}</span>
							{#if event.ranges.length > 1}<span>{rangeName(it.rangeId)}</span>{/if}
							<span class="font-medium">{$plural(it.lanes, { one: '# lane', other: '# lanes' })}</span>
							{#if it.assigned.length > 0}<span class="text-muted-foreground">{$t`lanes ${formatLanes(it.assigned)}`}</span>{:else if it.wanted?.length}<span class="text-muted-foreground">{$t`wanted ${formatLanes(it.wanted)}`}</span>{/if}
						</li>
					{/each}
				</ul>
				{#if draft}
					<form class="grid gap-3" onsubmit={saveShooters}>
						<h3 class="font-semibold">{$t`Announced shooters`}</h3>
						<ShooterList bind:shooters={draft} />
						<FormError message={shooterError} />
						<div class="flex gap-2">
							<Button type="submit" size="sm">{$t`Save`}</Button>
							<Button size="sm" variant="outline" onclick={() => (draft = null)}>{$t`Cancel`}</Button>
						</div>
					</form>
				{:else if r.shooters.length > 0 || editable}
					<section class="grid gap-2">
						<div class="flex items-center justify-between gap-2">
							<h3 class="font-semibold">{$plural(r.shooters.length, { one: '# shooter announced', other: '# shooters announced' })}</h3>
							{#if editable}<Button size="sm" variant="ghost" onclick={() => editShooters(r)}><Pencil />{$t`Edit`}</Button>{/if}
						</div>
						<ul class="grid gap-1">
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
					</section>
				{/if}
				{#if editable && !draft}
					<div class="flex flex-wrap gap-2 border-t pt-4">
						{#if r.status === 'pending'}
							<Button size="sm" onclick={() => openConfirm(r)}><Check />{$t`Confirm`}</Button>
						{/if}
						{#if r.status === 'confirmed' && waiting > 0}
							<Button size="sm" variant="outline" onclick={() => adopt(r)}><UserPlus />{$plural(waiting, { one: 'Take over # shooter', other: 'Take over # shooters' })}</Button>
						{/if}
						<Button size="sm" variant="outline" onclick={() => openDecline(r)}><X />{$t`Decline`}</Button>
					</div>
				{/if}
			</div>
		{/if}
	</Sheet.Content>
</Sheet.Root>

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

<Dialog.Root open={confirming !== null} onOpenChange={(o) => !o && (confirming = null)}>
	<Dialog.Content class="sm:max-w-md">
		{#if confirming}
			<Dialog.Header>
				<Dialog.Title>{$t`Confirm request of ${confirming.club}`}</Dialog.Title>
				<Dialog.Description>{$t`Free lanes are reserved for every slot. The announced shooters join this club of the master data:`}</Dialog.Description>
			</Dialog.Header>
			<form class="grid gap-4" onsubmit={confirm}>
				<ClubSelect {tenantId} name={confirming.club} {suggested} bind:value={clubChoice} label={$t`Club`} />
				<FormError message={error} />
				<Dialog.Footer><Button type="submit"><Check />{$t`Confirm`}</Button></Dialog.Footer>
			</form>
		{/if}
	</Dialog.Content>
</Dialog.Root>
