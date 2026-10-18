<script lang="ts">
	import { page } from '$app/state';
	import RequirePermission from '$lib/RequirePermission.svelte';

	let { children } = $props();
</script>

<RequirePermission tenantId={Number(page.params.tenantId)} permission="shooters.view">{@render children()}</RequirePermission>
