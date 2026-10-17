<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import AuthShell from '$lib/AuthShell.svelte';
	import { Button } from '$lib/components/ui/button';
	import { myKey } from '$lib/my.svelte';

	const queryClient = useQueryClient();
	let state = $state<'busy' | 'done' | string>('busy');

	onMount(async () => {
		try {
			await call(api.POST('/api/auth/verify', { body: { token: page.url.searchParams.get('token') ?? '' } }));
			state = 'done';
			await queryClient.invalidateQueries({ queryKey: myKey });
		} catch (err) {
			state = (err as Error).message;
		}
	});
</script>

<AuthShell>
	{#if state === 'busy'}
		<p class="text-sm text-muted-foreground">{$t`Confirming …`}</p>
	{:else if state === 'done'}
		<h1 class="text-2xl font-semibold tracking-tight">{$t`Email address confirmed`}</h1>
		<p class="mt-2 mb-6 text-sm text-muted-foreground">{$t`Thank you. Time for a training!`}</p>
		<Button href="/my">{$t`Open my log`}</Button>
	{:else}
		<h1 class="text-2xl font-semibold tracking-tight">{$t`That did not work`}</h1>
		<p class="mt-2 mb-6 text-sm text-muted-foreground">{state}</p>
		<Button href="/my/settings" variant="outline">{$t`Settings`}</Button>
	{/if}
</AuthShell>
