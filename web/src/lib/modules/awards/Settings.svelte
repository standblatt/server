<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { Switch } from '$lib/components/ui/switch';
	import { updateEvent } from '$lib/tenant';
	import type { EventProps } from '..';

	let { tenantId, event }: EventProps = $props();
	const queryClient = useQueryClient();

	/** Saves whether wreath cards show the award value. */
	async function save(printAwardValue: boolean) {
		await updateEvent(queryClient, tenantId, event, { settings: { ...event.settings, printAwardValue } });
		toast.success($t`Event saved`);
	}
</script>

<label class="flex items-center justify-between gap-3 text-sm">
	{$t`Print the award value on wreath cards`}
	<Switch checked={event.settings.printAwardValue ?? true} onCheckedChange={save} />
</label>
