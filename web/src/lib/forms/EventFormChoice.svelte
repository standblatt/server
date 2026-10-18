<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call, type Event } from '../api';
	import * as NativeSelect from '../components/ui/native-select';
	import FormField from '../FormField.svelte';
	import { meQuery } from '../me';
	import { hasClub } from '../modules';
	import { updateEvent } from '../tenant';
	import type { Source } from './editor.svelte';

	let {
		tenantId,
		event,
		setting,
		source,
		label
	}: { tenantId: number; event: Event; setting: 'wreathCardFormId' | 'controlStripFormId'; source: Source; label: string } = $props();
	const queryClient = useQueryClient();
	const me = meQuery();
	const forms = $derived(hasClub(me.data?.tenants.find((m) => m.tenantId === tenantId), 'forms'));
	const layouts = createQuery(() => ({
		queryKey: ['tenants', tenantId, 'layouts'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/layouts', { params: { path: { tenantId } } })),
		enabled: forms
	}));
	const choices = $derived((layouts.data ?? []).filter((l) => (l.source || 'results') === source));

	/** Saves the chosen print form, or the built-in one for 0. */
	async function save(id: number) {
		await updateEvent(queryClient, tenantId, event, { settings: { ...event.settings, [setting]: id || undefined } });
		toast.success($t`Event saved`);
	}
</script>

{#if forms && choices.length}
	<FormField {label}>
		<NativeSelect.Root class="w-full" value={event.settings[setting] ?? 0} onchange={(e) => save(Number(e.currentTarget.value))}>
			<NativeSelect.Option value={0}>{$t`Built-in`}</NativeSelect.Option>
			{#each choices as l (l.id)}<NativeSelect.Option value={l.id}>{l.name}</NativeSelect.Option>{/each}
		</NativeSelect.Root>
	</FormField>
{/if}
