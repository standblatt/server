<script lang="ts" generics="P extends Record<string, any>">
	import type { Lazy } from '.';

	let { load, props }: { load: Lazy<P>; props: P } = $props();
</script>

{#await load() then m}<m.default {...props} />{/await}
