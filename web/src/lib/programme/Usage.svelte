<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call, type Event, type Programme } from '$lib/api';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import { Switch } from '$lib/components/ui/switch';
	import FormField from '$lib/FormField.svelte';
	import { programmeTitle } from '$lib/format';
	import { keys, programmesQuery } from '$lib/tenant';
	import type { ProgrammeDraft } from './draft.svelte';

	let { tenantId, event, programme, draft }: { tenantId: number; event: Event; programme: Programme; draft: ProgrammeDraft } = $props();
	const queryClient = useQueryClient();
	const programmes = programmesQuery(
		() => tenantId,
		() => event.id
	);
	const others = $derived((programmes.data ?? []).filter((p) => p.id !== programme.id));
	const path = $derived({ params: { path: { tenantId, programmeId: programme.id } } });

	/** Runs a change that saves at once and refreshes the programme. */
	async function change(save: () => Promise<unknown>) {
		await save();
		await queryClient.invalidateQueries({ queryKey: keys.programmes(tenantId, event.id) });
		await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId, 'programme', programme.id] });
		toast.success($t`Programme "${programmeTitle(programme)}" saved`);
	}
</script>

<div class="grid gap-4 sm:grid-cols-2">
	<label class="flex items-center justify-between gap-3 text-sm sm:col-span-2">
		{$t`Active: shooters can register for it`}
		<Switch checked={programme.active} onCheckedChange={(active) => change(() => call(api.PUT('/api/tenants/{tenantId}/programmes/{programmeId}/active', { ...path, body: { active } })))} />
	</label>
	{#if event.ranges.length > 1}
		<FormField label={$t`Range`}>
			<NativeSelect.Root class="w-full" value={String(programme.rangeId ?? '')} onchange={(e) => change(() => call(api.PUT('/api/tenants/{tenantId}/programmes/{programmeId}/range', { ...path, body: e.currentTarget.value ? { rangeId: Number(e.currentTarget.value) } : {} })))}>
				<NativeSelect.Option value="">{$t`any`}</NativeSelect.Option>
				{#each event.ranges as r (r.id)}<NativeSelect.Option value={String(r.id)}>{r.name}</NativeSelect.Option>{/each}
			</NativeSelect.Root>
		</FormField>
	{/if}
	{#if others.length && !draft.repeatable}
		<FormField label={$t`Combinable with`}>
			<NativeSelect.Root class="w-full" value={programme.combinableWithId ?? 0} onchange={(e) => change(() => call(api.PUT('/api/tenants/{tenantId}/programmes/{programmeId}/combination', { ...path, body: Number(e.currentTarget.value) ? { combinableWithId: Number(e.currentTarget.value) } : {} })))}>
				<NativeSelect.Option value={0}>–</NativeSelect.Option>
				{#each others as o (o.id)}<NativeSelect.Option value={o.id}>{programmeTitle(o)}</NativeSelect.Option>{/each}
			</NativeSelect.Root>
		</FormField>
	{/if}
</div>
