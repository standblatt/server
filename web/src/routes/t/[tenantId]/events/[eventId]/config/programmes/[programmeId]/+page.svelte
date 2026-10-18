<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Download from '@lucide/svelte/icons/download';
	import { toast } from 'svelte-sonner';
	import { api, call, type ProgrammeDefinition } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { programmeTitle } from '$lib/format';
	import PageHeader from '$lib/PageHeader.svelte';
	import ProgrammeEditor from '$lib/ProgrammeEditor.svelte';
	import { ProgrammeDraft } from '$lib/programme/draft.svelte';
	import { eventQuery, keys, programmesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const programmeId = $derived(Number(page.params.programmeId));
	const queryClient = useQueryClient();
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const programmes = programmesQuery(
		() => tenantId,
		() => eventId
	);
	const programme = createQuery(() => ({
		queryKey: keys.programme(tenantId, eventId, programmeId),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/programmes/{programmeId}', { params: { path: { tenantId, programmeId } } }))
	}));
	const others = $derived((programmes.data ?? []).filter((p) => p.id !== programmeId));
	let draft = $state<ProgrammeDraft>();
	let draftOf = 0;
	$effect(() => {
		if (programme.data && draftOf !== programme.data.id) {
			draftOf = programme.data.id;
			draft = new ProgrammeDraft(programme.data);
		}
	});

	/** Saves the programme. */
	async function save(d: ProgrammeDefinition) {
		await call(api.PUT('/api/tenants/{tenantId}/programmes/{programmeId}', { params: { path: { tenantId, programmeId } }, body: d }));
		await queryClient.invalidateQueries({ queryKey: keys.programmes(tenantId, eventId) });
		toast.success($t`Programme "${d.name}" saved`);
	}
</script>

<PageHeader help="user-guide/programmes" title={programme.data ? programmeTitle(programme.data) : $t`Programme`} description={$t`Series, prices per age category and award limits.`}>
	{#snippet actions()}
		<Button variant="outline" size="sm" href="/api/tenants/{tenantId}/programmes/{programmeId}/export" download><Download />{$t`Export`}</Button>
	{/snippet}
</PageHeader>

{#if event.data && programme.data && draft}
	<div class="max-w-5xl">
		<ProgrammeEditor {tenantId} event={event.data} {draft} programme={programme.data} {others} locked={programme.data.inUse} submitLabel={$t`Save`} onsubmit={save} />
	</div>
{:else if programme.error}
	<p class="text-sm text-destructive">{programme.error.message}</p>
{/if}
