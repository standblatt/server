<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call, type ProgrammeDefinition } from '$lib/api';
	import PageHeader from '$lib/PageHeader.svelte';
	import ProgrammeEditor from '$lib/ProgrammeEditor.svelte';
	import { ProgrammeDraft } from '$lib/programme/draft.svelte';
	import { eventQuery, keys, programmesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const queryClient = useQueryClient();
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const programmes = programmesQuery(
		() => tenantId,
		() => eventId
	);
	const draft = new ProgrammeDraft();

	/** Creates the programme and opens its page. */
	async function create(d: ProgrammeDefinition) {
		const p = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/programmes', { params: { path: { tenantId, eventId } }, body: d }));
		await queryClient.invalidateQueries({ queryKey: keys.programmes(tenantId, eventId) });
		toast.success($t`Programme "${d.name}" created`);
		await goto(`/t/${tenantId}/events/${eventId}/config/programmes/${p.id}`, { replaceState: true });
	}
</script>

<PageHeader help="user-guide/programmes" title={$t`New programme`} description={$t`Series, prices per age category and award limits.`} />

{#if event.data}
	<div class="max-w-5xl"><ProgrammeEditor {tenantId} event={event.data} {draft} others={programmes.data ?? []} submitLabel={$t`Create`} onsubmit={create} /></div>
{/if}
