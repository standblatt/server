<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call, type ProgrammeDefinition } from './api';
	import * as Sheet from './components/ui/sheet';
	import ProgrammeEditor from './ProgrammeEditor.svelte';
	import { keys, programmesQuery } from './tenant';

	let {
		tenantId,
		eventId,
		programmeId,
		onclose
	}: { tenantId: number; eventId: number; programmeId?: number; onclose: () => void } = $props();

	const queryClient = useQueryClient();
	const programme = createQuery(() => ({
		queryKey: ['tenants', tenantId, 'programme', programmeId],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/programmes/{programmeId}', { params: { path: { tenantId, programmeId: programmeId! } } })),
		enabled: !!programmeId
	}));

	const all = programmesQuery(
		() => tenantId,
		() => eventId
	);
	const others = $derived((all.data ?? []).filter((p) => p.id !== programmeId));

	/** Refreshes the programme queries, shows the message and closes the sheet. */
	async function done(message: string) {
		await queryClient.invalidateQueries({ queryKey: keys.programmes(tenantId, eventId) });
		if (programmeId) await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId, 'programme', programmeId] });
		toast.success(message);
		onclose();
	}

	/** Creates a new programme in the event. */
	async function create(d: ProgrammeDefinition) {
		await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/programmes', { params: { path: { tenantId, eventId } }, body: d }));
		await done($t`Programme "${d.name}" created`);
	}

	/** Saves the edited programme. */
	async function save(d: ProgrammeDefinition) {
		await call(api.PUT('/api/tenants/{tenantId}/programmes/{programmeId}', { params: { path: { tenantId, programmeId: programmeId! } }, body: d }));
		await done($t`Programme "${d.name}" saved`);
	}
</script>

<Sheet.Root open onOpenChange={(o) => !o && onclose()}>
	<Sheet.Content class="gap-0 p-0 data-[side=right]:w-full data-[side=right]:sm:max-w-3xl">
		<Sheet.Header class="border-b px-6 py-4">
			<Sheet.Title>{programmeId ? $t`Edit programme` : $t`New programme`}</Sheet.Title>
			<Sheet.Description>{$t`Series, prices per age category and award limits.`}</Sheet.Description>
		</Sheet.Header>
		{#if !programmeId}
			<ProgrammeEditor {others} submitLabel={$t`Create`} onsubmit={create} />
		{:else if programme.data}
			<ProgrammeEditor initial={programme.data} {others} locked={programme.data.inUse} submitLabel={$t`Save`} onsubmit={save} />
		{:else if programme.error}
			<p class="p-6 text-sm text-destructive">{programme.error.message}</p>
		{/if}
	</Sheet.Content>
</Sheet.Root>
