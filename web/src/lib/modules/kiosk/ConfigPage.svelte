<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { Switch } from '$lib/components/ui/switch';
	import { updateEvent } from '$lib/tenant';
	import type { EventProps } from '..';
	import KioskTokens from './KioskTokens.svelte';

	let { tenantId, event }: EventProps = $props();
	const queryClient = useQueryClient();

	/** Saves whether kiosk displays show the lane status. */
	async function save(showLaneStatus: boolean) {
		await updateEvent(queryClient, tenantId, event, { settings: { ...event.settings, showLaneStatus } });
		toast.success($t`Event saved`);
	}
</script>

<div class="space-y-4">
	{#if event.ranges.length}
		<label class="flex items-center justify-between gap-3 rounded-lg border bg-card p-4 text-sm">
			{$t`Show lane status on the rankings`}
			<Switch checked={event.settings.showLaneStatus ?? false} onCheckedChange={save} />
		</label>
	{/if}
	<section class="overflow-hidden rounded-lg border bg-card"><KioskTokens {tenantId} eventId={event.id} /></section>
</div>
