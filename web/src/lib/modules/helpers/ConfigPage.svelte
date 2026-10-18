<script lang="ts">
	import type { EventProps } from '..';
	import HelperCodes from './HelperCodes.svelte';

	let { tenantId, event }: EventProps = $props();
</script>

<section class="overflow-hidden rounded-lg border bg-card"><HelperCodes {tenantId} eventId={event.id} /></section>
