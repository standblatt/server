<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import FormField from '$lib/FormField.svelte';
	import { clubsQuery, updateEvent } from '$lib/tenant';
	import type { EventProps } from '..';

	let { tenantId, event }: EventProps = $props();
	const queryClient = useQueryClient();
	const clubs = clubsQuery(() => tenantId);

	/** Saves the club that takes part out of competition (0 = none). */
	async function save(id: number) {
		await updateEvent(queryClient, tenantId, event, { settings: { ...event.settings, excludedClubId: id || undefined } });
		toast.success($t`Event saved`);
	}
</script>

<FormField label={$t`Club excluded from the club competition`}>
	<NativeSelect.Root class="w-full" value={event.settings.excludedClubId ?? 0} onchange={(e) => save(Number(e.currentTarget.value))}>
		<NativeSelect.Option value={0}>{$t`none`}</NativeSelect.Option>
		{#each clubs.data ?? [] as c (c.id)}<NativeSelect.Option value={c.id}>{c.name}</NativeSelect.Option>{/each}
	</NativeSelect.Root>
</FormField>
