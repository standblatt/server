<script lang="ts">
	import type { EventProps, Lazy } from '.';

	let { load, tenantId, event }: { load: Lazy } & EventProps = $props();
</script>

{#await load() then m}<m.default {tenantId} {event} />{/await}
