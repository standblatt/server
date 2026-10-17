<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { setMode, userPrefersMode } from 'mode-watcher';
	import Monitor from '@lucide/svelte/icons/monitor';
	import Moon from '@lucide/svelte/icons/moon';
	import Sun from '@lucide/svelte/icons/sun';
	import { Button } from '$lib/components/ui/button';

	const next = { system: 'light', light: 'dark', dark: 'system' } as const;
	const label = $derived({ system: $t`System theme`, light: $t`Light theme`, dark: $t`Dark theme` });
</script>

<Button variant="ghost" size="icon-sm" title={label[userPrefersMode.current]} aria-label={label[userPrefersMode.current]} onclick={() => setMode(next[userPrefersMode.current])}>
	{#if userPrefersMode.current === 'light'}<Sun />{:else if userPrefersMode.current === 'dark'}<Moon />{:else}<Monitor />{/if}
</Button>
