<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { online } from 'svelte/reactivity/window';
	import CloudOff from '@lucide/svelte/icons/cloud-off';
</script>

{#if online.current === false}
	<div class="flex items-center gap-2 bg-destructive px-4 py-2 text-sm font-medium text-white" role="alert" data-testid="offline">
		<CloudOff class="size-4 shrink-0" />{$t`No connection. Changes are only saved when you are back online.`}
	</div>
{/if}
