<script lang="ts">
	import type { EventProps } from '..';
	import Classifications from './Classifications.svelte';

	let { tenantId, event }: EventProps = $props();
</script>

<Classifications {tenantId} eventId={event.id} />
