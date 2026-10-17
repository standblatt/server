<script lang="ts">
	import { Trans, t } from 'lingui-for-svelte/macro';
	import { createQuery } from '@tanstack/svelte-query';
	import { api, call } from './api';
	import { Button } from './components/ui/button';
	import { Input } from './components/ui/input';
	import { clubsQuery, keys } from './tenant';

	let {
		tenantId,
		value = $bindable(''),
		ready = $bindable(true),
		id,
		required = false
	}: { tenantId: number; value?: string; ready?: boolean; id: string; required?: boolean } = $props();

	/** The name clubs are matched by, like domain.ClubKey: lower case with single spaces. */
	const clubKey = (s: string) => s.toLowerCase().split(/\s+/).filter(Boolean).join(' ');

	const clubs = clubsQuery(() => tenantId);
	const name = $derived(value.trim());
	const known = $derived(!!clubs.data?.some((c) => clubKey(c.name) === clubKey(name)));
	let typed = $state('');
	let confirmed = $state('');
	$effect(() => {
		const v = name;
		const timer = setTimeout(() => (typed = v), 300);
		return () => clearTimeout(timer);
	});
	const similar = createQuery(() => ({
		queryKey: [...keys.clubs(tenantId), 'similar', typed],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/clubs/similar', { params: { path: { tenantId }, query: { name: typed } } })),
		enabled: typed !== ''
	}));
	$effect(() => {
		ready = name === '' || known || confirmed === name;
	});
</script>

<Input list="{id}-names" bind:value placeholder={$t`Club`} aria-label={$t`Club`} {required} />
<datalist id="{id}-names">{#each clubs.data ?? [] as c (c.id)}<option value={c.name}></option>{/each}</datalist>
{#if name && !known && clubs.data}
	<div class="mt-2 grid gap-2 rounded-md border bg-muted/40 p-3 text-sm" role="status">
		{#if confirmed === name}
			<p><Trans>“{name}” will be added as a new club.</Trans></p>
		{:else}
			<p><Trans>“{name}” is not in the club list yet.</Trans></p>
			{#if typed === name && similar.data?.similar.length}
				<div class="flex flex-wrap items-center gap-1.5">
					<span class="text-muted-foreground">{$t`Did you mean`}</span>
					{#each similar.data.similar as c (c.id)}<Button size="sm" variant="outline" onclick={() => (value = c.name)}>{c.name}</Button>{/each}
				</div>
			{/if}
			<div><Button size="sm" variant="secondary" onclick={() => (confirmed = name)}>{$t`Create new club “${name}”`}</Button></div>
		{/if}
	</div>
{/if}
