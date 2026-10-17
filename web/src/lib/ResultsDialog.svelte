<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Award from '@lucide/svelte/icons/award';
	import Pencil from '@lucide/svelte/icons/pencil';
	import { toast } from 'svelte-sonner';
	import { Badge } from './components/ui/badge';
	import { api, call, type SeriesResult } from './api';
	import { ageCategoryLabel, registrationStatus } from './labels';
	import { ask } from './confirm.svelte';
	import { intlLocale } from './locale.svelte';
	import { Button } from './components/ui/button';
	import * as Dialog from './components/ui/dialog';
	import EntryCard from './EntryCard.svelte';
	import PrintButton from './PrintButton.svelte';
	import { eventQuery, keys } from './tenant';

	let {
		tenantId,
		eventId,
		registrationId,
		canEdit = true,
		canPrint = true,
		canPrintForm = true,
		canIssue = true,
		canRevoke = true,
		onclose,
		onedit
	}: {
		tenantId: number;
		eventId: number;
		registrationId: number;
		/** Helper links cannot edit the registration or print. */
		canEdit?: boolean;
		canPrint?: boolean;
		canPrintForm?: boolean;
		canIssue?: boolean;
		canRevoke?: boolean;
		onclose: () => void;
		onedit: (registrationId: number) => void;
	} = $props();

	const queryClient = useQueryClient();
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const key = $derived([...keys.registrations(tenantId, eventId), registrationId, 'results']);
	const sheet = createQuery(() => ({
		queryKey: key,
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/registrations/{registrationId}/results', { params: { path: { tenantId, registrationId } } }))
	}));

	let body = $state<HTMLElement>();
	let focused = false;
	$effect(() => {
		if (!sheet.data || !body || focused) return;
		focused = true;
		const inputs = [...body.querySelectorAll<HTMLInputElement>('input[data-score]')];
		(inputs.find((i) => i.value === '') ?? inputs[0])?.focus();
	});

	/** Stores the returned result sheet in the cache and refreshes the registrations list. */
	async function changed(s: NonNullable<typeof sheet.data>) {
		queryClient.setQueryData(key, s);
		await queryClient.invalidateQueries({ queryKey: keys.registrations(tenantId, eventId), refetchType: 'active', exact: true });
	}

	/** Saves the series results of an entry. */
	async function save(entryId: number, series: SeriesResult[]) {
		await changed(await call(api.PUT('/api/tenants/{tenantId}/entries/{entryId}/results', { params: { path: { tenantId, entryId } }, body: { series } })));
	}

	/** Marks the awards of the registration as handed out. */
	async function issueAwards() {
		await call(api.POST('/api/tenants/{tenantId}/registrations/{registrationId}/awards-issued', { params: { path: { tenantId, registrationId } } }));
		await queryClient.invalidateQueries({ queryKey: key });
		await queryClient.invalidateQueries({ queryKey: keys.registrations(tenantId, eventId), refetchType: 'active', exact: true });
		toast.success($t`Awards handed out`);
	}

	/** Asks for confirmation, then takes the awards handout back. */
	async function revokeAwards() {
		if (!(await ask({ title: $t`Take the awards handout back?`, description: $t`The awards count as open again.`, action: $t`Take back` }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/registrations/{registrationId}/awards-issued', { params: { path: { tenantId, registrationId } } }));
		await queryClient.invalidateQueries({ queryKey: key });
		await queryClient.invalidateQueries({ queryKey: keys.registrations(tenantId, eventId), refetchType: 'active', exact: true });
	}

	/** Sets the award tier of an entry (none when undefined). */
	async function award(entryId: number, tierId: number | undefined) {
		await changed(await call(api.PUT('/api/tenants/{tenantId}/entries/{entryId}/award', { params: { path: { tenantId, entryId } }, body: { tierId } })));
	}
</script>

<Dialog.Root open onOpenChange={(o) => !o && onclose()}>
	<Dialog.Content class="flex max-h-[92svh] flex-col gap-0 p-0 sm:max-w-4xl" onOpenAutoFocus={(e) => e.preventDefault()}>
		<Dialog.Header class="border-b px-6 py-4 pr-12">
			{#if sheet.data}
				{@const r = sheet.data.registration}
				<Dialog.Title class="text-base">{$t`Results: ${r.shooter.lastname} ${r.shooter.firstname}`}</Dialog.Title>
				<Dialog.Description>{r.clubName || $t`no club`} · {$t(ageCategoryLabel(r.ageCategory))} · {r.lane ? $t`lane ${r.lane}` : $t(registrationStatus(r.status))}</Dialog.Description>
				<div class="mt-2 flex flex-wrap items-center gap-2">
					{#if canEdit}<Button size="sm" variant="outline" onclick={() => onedit(r.id)}><Pencil />{$t`Edit registration`}</Button>{/if}
					{#if sheet.data.entries.some((e) => e.award)}
						{#if r.awardsIssuedAt}
							<Badge variant="outline" class="border-emerald-500/50 text-emerald-700 dark:text-emerald-400">{$t`awards handed out`} · {new Date(r.awardsIssuedAt).toLocaleTimeString(intlLocale(), { hour: '2-digit', minute: '2-digit' })}</Badge>
							{#if canRevoke && !event.data?.archivedAt}<Button size="sm" variant="ghost" onclick={revokeAwards}>{$t`Take back`}</Button>{/if}
						{:else if canIssue && !event.data?.archivedAt}
							<Button size="sm" variant="outline" onclick={issueAwards}><Award />{$t`Awards handed out`}</Button>
						{/if}
					{/if}
					{#if r.status === 'completed' && sheet.data.entries.every((e) => e.status === 'completed')}
						<span class="ml-auto flex flex-wrap items-center gap-2 text-sm text-emerald-700 dark:text-emerald-400" role="status">
							{$t`All programmes have results.`}
							{#if canPrint}<PrintButton {tenantId} target={{ document: 'result_sheet', registrationId: r.id }} />{/if}
						</span>
					{/if}
				</div>
			{:else}
				<Dialog.Title class="text-base">{$t`Results`}</Dialog.Title>
			{/if}
		</Dialog.Header>
		<div bind:this={body} class="grid min-h-0 flex-1 gap-4 overflow-y-auto bg-muted/30 p-4 sm:p-6">
			{#if sheet.data}
				{#each sheet.data.entries as entry (entry.entryId)}
					<EntryCard {entry} {tenantId} readonly={!!event.data?.archivedAt} print={canPrint} printForm={canPrintForm} onsave={(s) => save(entry.entryId, s)} onaward={(t) => award(entry.entryId, t)} />
				{/each}
			{:else if sheet.error}
				<p class="text-sm text-destructive">{sheet.error.message}</p>
			{:else}
				<p class="text-sm text-muted-foreground">{$t`Loading …`}</p>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
