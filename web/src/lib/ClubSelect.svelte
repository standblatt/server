<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { Club } from './api';
	import * as NativeSelect from './components/ui/native-select';
	import { clubsQuery } from './tenant';

	let {
		tenantId,
		name,
		suggested,
		value = $bindable(''),
		allowNew = true,
		exclude,
		label
	}: { tenantId: number; name: string; suggested: Club[]; value?: string; allowNew?: boolean; exclude?: number; label: string } = $props();

	const clubs = clubsQuery(() => tenantId);
	const top = $derived(suggested.filter((c) => c.id !== exclude));
	const rest = $derived((clubs.data ?? []).filter((c) => c.id !== exclude && !top.some((x) => x.id === c.id)));
</script>

<NativeSelect.Root class="w-full" bind:value aria-label={label} required>
	<NativeSelect.Option value="" disabled>{$t`Choose a club`}</NativeSelect.Option>
	{#if top.length}
		<NativeSelect.OptGroup label={$t`Suggested`}>
			{#each top as c (c.id)}<NativeSelect.Option value={String(c.id)}>{c.name}</NativeSelect.Option>{/each}
		</NativeSelect.OptGroup>
	{/if}
	{#if allowNew}<NativeSelect.Option value="new">{$t`New club “${name}”`}</NativeSelect.Option>{/if}
	{#if rest.length}
		<NativeSelect.OptGroup label={$t`Other clubs`}>
			{#each rest as c (c.id)}<NativeSelect.Option value={String(c.id)}>{c.name}</NativeSelect.Option>{/each}
		</NativeSelect.OptGroup>
	{/if}
</NativeSelect.Root>
