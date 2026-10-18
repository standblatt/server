<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { intlLocale } from '$lib/locale.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Copy from '@lucide/svelte/icons/copy';
	import MoreHorizontal from '@lucide/svelte/icons/ellipsis';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { localInput } from '$lib/format';
	import { api, call, type ReservationView } from '$lib/api';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import { Switch } from '$lib/components/ui/switch';
	import { ask } from '$lib/confirm.svelte';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { keys } from '$lib/tenant';
	import BookingRequests from './booking/BookingRequests.svelte';
	import type { EventProps } from '..';
	import { membershipOf, meQuery } from '$lib/me';

	let { tenantId, event }: EventProps = $props();
	const eventId = $derived(event.id);
	const queryClient = useQueryClient();
	const ranges = $derived(event.ranges);
	const me = meQuery();
	const isAdmin = $derived(membershipOf(me.data, tenantId)?.role === 'admin');
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
			form = { id: copy ? 0 : r.id, rangeId: r.rangeId, lanes: r.lanesText, start: localInput(new Date(r.startsAt)), end: localInput(new Date(r.endsAt)), lead: r.leadMinutes, description: r.description };
		} else {
			const start = nextQuarter();
			form = { id: 0, rangeId: ranges[0]?.id ?? 0, lanes: '', start: localInput(start), end: localInput(new Date(start.getTime() + 3600000)), lead: 30, description: '' };
		}
	}

	/** Sets the end one hour after the new start. */
	function startChanged() {
		if (form) form.end = localInput(new Date(new Date(form.start).getTime() + 3600000));
	}

	/** Creates or updates the reservation from the form. */
	async function save(e: SubmitEvent) {
		e.preventDefault();
		if (!form) return;
		error = '';
		const body = {
			lanes: form.lanes,
			startsAt: new Date(form.start).toISOString(),
			endsAt: new Date(form.end).toISOString(),
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
	/** Formats a timestamp as local hours and minutes. */
	const time = (iso: string) => new Date(iso).toLocaleTimeString(intlLocale(), { hour: '2-digit', minute: '2-digit' });
	const byDay = $derived.by(() => {
		const days = new Map<string, Row[]>();
		for (const r of [...(reservations.data ?? [])].sort((a, b) => a.startsAt.localeCompare(b.startsAt))) {
			const day = new Date(r.startsAt).toLocaleDateString(intlLocale(), { weekday: 'long', day: '2-digit', month: 'long', year: 'numeric' });
			days.set(day, [...(days.get(day) ?? []), r]);
		}
		return [...days.entries()];
	});
</script>

<PageHeader help="user-guide/reservations" title={$t`Reservations`} description={$t`Lanes reserved for trainings or groups; they show as reserved from the lead time on.`}>
	{#snippet actions()}
		<label class="flex items-center gap-2 text-sm text-muted-foreground"><Switch bind:checked={showPast} />{$t`show past`}</label>
		<Button onclick={() => open()}><Plus />{$t`New reservation`}</Button>
	{/snippet}
</PageHeader>

<div class="space-y-4">
	<BookingRequests {tenantId} {event} {isAdmin} />
	{#each byDay as [day, list] (day)}
		<section class="overflow-hidden rounded-lg border bg-card">
			<div class="border-b px-4 py-2.5"><h2 class="text-sm font-semibold">{day}</h2></div>
			<ul class="divide-y">
				{#each list as r (r.id)}
					{@const st = status(r)}
					<li class="flex flex-wrap items-center gap-x-4 gap-y-1 px-4 py-2.5 text-sm" class:opacity-50={st === 'finished'} aria-label={$t`Reservation ${r.description}`}>
						<Badge
							variant="outline"
							class="w-20 justify-center {st === 'active' ? 'border-amber-500/50 bg-amber-500/10 text-amber-700 dark:text-amber-400' : st === 'planned' ? 'border-sky-500/50 text-sky-700 dark:text-sky-400' : ''}">{statusLabel(st)}</Badge
						>
						<span class="w-28 font-medium tabular-nums">{time(r.startsAt)}–{time(r.endsAt)}</span>
						<span class="font-semibold">{r.description}</span>
						<span class="text-muted-foreground">{ranges.length > 1 ? `${r.rangeName} · ` : ''}{$t`lanes ${r.lanesText}`}</span>
						{#if r.leadMinutes}<span class="text-xs text-muted-foreground">{$t`shown ${r.leadMinutes} min before`}</span>{/if}
						<DropdownMenu.Root>
							<DropdownMenu.Trigger>
								{#snippet child({ props })}<Button variant="ghost" size="icon-sm" class="ml-auto" aria-label={$t`Actions for ${r.description}`} {...props}><MoreHorizontal /></Button>{/snippet}
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end">
								<DropdownMenu.Item onSelect={() => open(r)}><Pencil />{$t`Edit`}</DropdownMenu.Item>
								<DropdownMenu.Item onSelect={() => open(r, true)}><Copy />{$t`Copy`}</DropdownMenu.Item>
								<DropdownMenu.Separator />
								<DropdownMenu.Item variant="destructive" onSelect={() => remove(r)}><Trash />{$t`Delete`}</DropdownMenu.Item>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</li>
				{/each}
			</ul>
		</section>
	{:else}
		<div class="rounded-lg border border-dashed py-12 text-center text-sm text-muted-foreground">{$t`No reservations.`}</div>
	{/each}
</div>

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
