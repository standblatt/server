<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call } from '$lib/api';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import FormField from '$lib/FormField.svelte';
	import { programmeTitle } from '$lib/format';
	import Section from '$lib/programme/Section.svelte';
	import { keys } from '$lib/tenant';
	import type { ProgrammeProps } from '..';

	let { tenantId, event, programme }: ProgrammeProps = $props();
	const queryClient = useQueryClient();
	const layouts = createQuery(() => ({
		queryKey: ['tenants', tenantId, 'layouts'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/layouts', { params: { path: { tenantId } } }))
	}));

	const results = $derived((layouts.data ?? []).filter((l) => (l.source || 'results') === 'results'));

	/** Chooses the print form of the programme (its ID) or none (""). */
	async function choose(value: string) {
		if (!programme) return;
		await call(api.PUT('/api/tenants/{tenantId}/programmes/{programmeId}/form', { params: { path: { tenantId, programmeId: programme.id } }, body: value ? { layoutId: Number(value) } : {} }));
		await queryClient.invalidateQueries({ queryKey: keys.programmes(tenantId, event.id) });
		toast.success($t`Programme "${programmeTitle(programme)}" saved`);
	}
</script>

{#if programme && results.length}
	<Section title={$t`Print form`} hint={$t`The form the results of this programme are printed onto.`}>
		<FormField label={$t`Print form`}>
			<NativeSelect.Root class="w-full sm:w-80" value={String(programme.layoutId ?? '')} onchange={(e) => choose(e.currentTarget.value)}>
				<NativeSelect.Option value="">{$t`none`}</NativeSelect.Option>
				{#each results as l (l.id)}<NativeSelect.Option value={String(l.id)}>{l.name}</NativeSelect.Option>{/each}
			</NativeSelect.Root>
		</FormField>
	</Section>
{/if}
