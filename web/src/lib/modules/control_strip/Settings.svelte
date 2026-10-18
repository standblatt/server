<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import EventFormChoice from '$lib/forms/EventFormChoice.svelte';
	import type { EventProps } from '..';

	let { tenantId, event }: EventProps = $props();
</script>

<EventFormChoice {tenantId} {event} setting="controlStripFormId" source="registration" label={$t`Print form of the control strips`} />
