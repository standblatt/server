<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Flag from '@lucide/svelte/icons/flag';
	import ListChecks from '@lucide/svelte/icons/list-checks';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call, type ProgrammeChoice } from './api';
	import { ageCategoryLabel, registrationStatus } from './labels';
	import { Badge } from './components/ui/badge';
	import { Button } from './components/ui/button';
	import { Checkbox } from './components/ui/checkbox';
	import * as Dialog from './components/ui/dialog';
	import { ask } from './confirm.svelte';
	import { programmeTitle, withRange } from './format';
	import PrintButton from './PrintButton.svelte';
	import { printDocument } from './print';
	import RegistrationForm from './RegistrationForm.svelte';
	import { eventQuery, keys, lanesQuery, programmesQuery, registrationsQuery } from './tenant';

	let {
		tenantId,
		eventId,
		shooterId,
		clubId,
		registrationId,
		lanesOnly = false,
		onclose,
		onresults,
		onopen
	}: {
		tenantId: number;
		eventId: number;
		shooterId?: number;
		clubId?: number;
		registrationId?: number;
		/** Helper links of the ammunition desk: move lanes and print the lane sheet, nothing else. */
		lanesOnly?: boolean;
		onclose: () => void;
		onresults: (registrationId: number) => void;
		onopen?: (registrationId: number) => void;
	} = $props();

	const queryClient = useQueryClient();
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const programmes = programmesQuery(
		() => tenantId,
		() => eventId
	);
	const lanes = lanesQuery(
		() => tenantId,
		() => eventId
	);
	const shooter = createQuery(() => ({
		queryKey: keys.shooter(tenantId, shooterId ?? 0),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/shooters/{shooterId}', { params: { path: { tenantId, shooterId: shooterId! } } })),
		enabled: !!shooterId
	}));
	const registration = createQuery(() => ({
		queryKey: [...keys.registrations(tenantId, eventId), registrationId],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/registrations/{registrationId}', { params: { path: { tenantId, registrationId: registrationId! } } })),
		enabled: !!registrationId
	}));

	const eventRegistrations = registrationsQuery(
		() => tenantId,
		() => eventId
	);
	const previous = $derived(shooterId ? (eventRegistrations.data ?? []).filter((r) => r.shooter.id === shooterId) : []);
	const taken = $derived(new Set(previous.flatMap((r) => r.entries.map((e) => e.programmeId))));
	/** Name (with variant) of a programme by ID. */
	const titleOf = (id: number) => {
		const p = programmes.data?.find((x) => x.id === id);
		return p ? programmeTitle(p) : `#${id}`;
	};

	let printLaneSheet = $state(true);
	let printControlStrip = $state(true);

	const archived = $derived(!!event.data?.archivedAt);
	const person = $derived(registration.data?.shooter ?? shooter.data);
	const ready = $derived(!!person && !!event.data && !!programmes.data && !!lanes.data);

	/** Invalidates the registrations of the event. */
	async function refresh() {
		await queryClient.invalidateQueries({ queryKey: keys.registrations(tenantId, eventId) });
	}

	/** Name of a range of the event by ID. */
	const rangeName = (id?: number) => event.data?.ranges.find((r) => r.id === id)?.name ?? '';
	/** Lane label with the range name when the event has several ranges. */
	const laneLabel = (reg: { rangeId?: number; lane?: number }) => (reg.lane ? withRange(event.data?.ranges ?? [], reg.rangeId, $t`Lane ${reg.lane}`) : '');

	/** Registers the shooter, prints the chosen documents and closes the dialog. */
	async function register(r: { clubId?: number; rangeId?: number; lane?: number; programmes: ProgrammeChoice[] }) {
		const reg = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/registrations', { params: { path: { tenantId, eventId } }, body: { shooterId: shooterId!, ...r } }));
		if (printLaneSheet) await printDocument(tenantId, { document: 'lane_sheet', registrationId: reg.id });
		if (printControlStrip) await printDocument(tenantId, { document: 'control_strip', registrationId: reg.id });
		await refresh();
		toast.success(
			reg.lane && (event.data?.ranges.length ?? 0) > 1
				? $t`${reg.shooter.firstname} ${reg.shooter.lastname} registered on ${rangeName(reg.rangeId)}, lane ${reg.lane}`
				: reg.lane
					? $t`${reg.shooter.firstname} ${reg.shooter.lastname} registered on lane ${reg.lane}`
					: $t`${reg.shooter.firstname} ${reg.shooter.lastname} registered`
		);
		onclose();
	}

	/** Saves lane and programmes of the open registration and closes the dialog. */
	async function save(r: { rangeId?: number; lane?: number; programmes: ProgrammeChoice[] }) {
		await call(api.PUT('/api/tenants/{tenantId}/registrations/{registrationId}', { params: { path: { tenantId, registrationId: registrationId! } }, body: { rangeId: r.rangeId, lane: r.lane, programmes: r.programmes } }));
		await refresh();
		toast.success($t`Registration saved`);
		onclose();
	}

	/** Asks for confirmation, then finishes or deletes the open registration. */
	async function action(kind: 'finish' | 'delete') {
		const r = registration.data!;
		const name = `${r.shooter.firstname} ${r.shooter.lastname}`;
		const where = laneLabel(r) || $t`the waiting list`;
		const ok = await ask(
			kind === 'finish'
				? { title: $t`Finish shooting?`, description: $t`${name} leaves ${where}.`, action: $t`Finish` }
				: { title: $t`Delete registration?`, description: $t`The registration of ${name} and its results are deleted.`, action: $t`Delete`, destructive: true }
		);
		if (!ok) return;
		const p = { params: { path: { tenantId, registrationId: r.id } } };
		try {
			if (kind === 'finish') await call(api.POST('/api/tenants/{tenantId}/registrations/{registrationId}/finish', p));
			else await call(api.DELETE('/api/tenants/{tenantId}/registrations/{registrationId}', p));
			await refresh();
			toast.success(kind === 'finish' ? $t`${name} finished` : $t`Registration deleted`);
			onclose();
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
</script>

<Dialog.Root open onOpenChange={(o) => !o && onclose()}>
	<Dialog.Content class="flex max-h-[92svh] flex-col gap-0 p-0 sm:max-w-5xl">
		<Dialog.Header class="border-b px-6 py-4 pr-12">
			<Dialog.Title class="flex flex-wrap items-center gap-2 text-base">
				{#if person}{person.lastname} {person.firstname}{/if}
				{#if registration.data}
					<Badge variant="outline">#{registration.data.id} · {$t(registrationStatus(registration.data.status))}</Badge>
				{/if}
			</Dialog.Title>
			<Dialog.Description>
				{registrationId ? $t`Edit registration` : $t`New registration`}{registration.data ? ` · ${registration.data.clubName || $t`no club`} · ${$t(ageCategoryLabel(registration.data.ageCategory))}` : ''}
			</Dialog.Description>
			{#if registration.data}
				{@const r = registration.data}
				<div class="mt-2 flex flex-wrap items-center gap-2">
					{#if !lanesOnly}<Button size="sm" onclick={() => onresults(r.id)}><ListChecks />{$t`Results`}</Button>{/if}
					{#if r.status !== 'completed' && !archived && !lanesOnly}
						<Button size="sm" variant="outline" onclick={() => action('finish')}><Flag />{$t`Finish shooting`}</Button>
					{/if}
					<PrintButton {tenantId} target={{ document: 'lane_sheet', registrationId: r.id }} />
					{#if !lanesOnly}<PrintButton {tenantId} target={{ document: 'control_strip', registrationId: r.id }} />{/if}
					{#if !archived && !lanesOnly}<Button size="sm" variant="destructive" class="ml-auto" onclick={() => action('delete')}><Trash />{$t`Delete`}</Button>{/if}
				</div>
			{/if}
			{#if previous.length > 0}
				<div class="mt-3 space-y-1.5 rounded-md border border-amber-500/40 bg-amber-500/10 px-3 py-2 text-sm" role="status" data-testid="previous-registrations">
					<p class="font-medium">{$t`Already registered at this event. Programmes already booked cannot be chosen again.`}</p>
					{#each previous as r (r.id)}
						<div class="flex flex-wrap items-center gap-2">
							<span>#{r.id} · {$t(registrationStatus(r.status))}{r.lane ? ` · ${laneLabel(r)}` : ''} · {r.entries.map((e) => titleOf(e.programmeId)).join(', ')}</span>
							{#if onopen}<Button size="xs" variant="outline" class="ml-auto" onclick={() => onopen(r.id)}>{$t`Open registration`}</Button>{/if}
						</div>
					{/each}
				</div>
			{/if}
		</Dialog.Header>

		{#if ready && person}
			{#if registration.data}
				<RegistrationForm
					shooter={person}
					eventDate={event.data!.startsOn}
					programmes={programmes.data!}
					lanes={lanes.data!}
					initial={registration.data}
					readonly={archived}
					{lanesOnly}
					submitLabel={$t`Save`}
					onsubmit={save}
				/>
			{:else}
				<RegistrationForm
					shooter={person}
					eventDate={event.data!.startsOn}
					programmes={programmes.data!}
					lanes={lanes.data!}
					initialClubId={clubId}
					{taken}
					submitLabel={$t`Register`}
					onsubmit={register}
				>
					{#snippet footer()}
						<label class="flex items-center gap-2 text-sm"><Checkbox bind:checked={printLaneSheet} />{$t`Print lane sheet`}</label>
						<label class="flex items-center gap-2 text-sm"><Checkbox bind:checked={printControlStrip} />{$t`Print control strip`}</label>
					{/snippet}
				</RegistrationForm>
			{/if}
		{:else if shooter.error || registration.error}
			<p class="p-6 text-sm text-destructive">{(shooter.error ?? registration.error)?.message}</p>
		{:else}
			<p class="p-6 text-sm text-muted-foreground">{$t`Loading …`}</p>
		{/if}
	</Dialog.Content>
</Dialog.Root>
