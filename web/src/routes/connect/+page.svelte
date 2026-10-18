<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { api, call, type PrintAgent } from '$lib/api';
	import AuthShell from '$lib/AuthShell.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { can, meQuery } from '$lib/me';
	import { rangesQuery } from '$lib/tenant';

	const me = meQuery();
	const clubs = $derived((me.data?.tenants ?? []).filter((m) => can(m, 'ranges.manage')));
	let code = $state(page.url.searchParams.get('code') ?? '');
	const typed = $derived(code.toUpperCase().replace(/[^A-Z0-9]/g, ''));
	const pairing = createQuery(() => ({
		queryKey: ['agent-pairing', typed],
		queryFn: () => call(api.GET('/api/agent-pairings/{code}', { params: { path: { code: typed } } })),
		enabled: !!me.data && typed.length === 8,
		retry: false
	}));
	let tenantId = $state(0);
	$effect(() => {
		if (!tenantId && clubs.length) tenantId = clubs[0].tenantId;
	});
	const ranges = rangesQuery(() => tenantId);
	let rangeId = $state(0);
	$effect(() => {
		const list = ranges.data ?? [];
		if (!list.some((r) => r.id === rangeId)) rangeId = list[0]?.id ?? 0;
	});
	let name = $state('');
	$effect(() => {
		if (pairing.data && !name) name = pairing.data.name;
	});
	let error = $state('');
	let connected = $state<PrintAgent | null>(null);
	/** The display name of an operating system the agent reports. */
	const osLabel = (os: string) => ({ windows: 'Windows', darwin: 'macOS', linux: 'Linux' })[os] ?? os;

	/** Confirms the code: the computer becomes a print agent of the chosen range. */
	async function connect(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			connected = await call(
				api.POST('/api/tenants/{tenantId}/ranges/{rangeId}/print-agents/pair', { params: { path: { tenantId, rangeId } }, body: { code: typed, name } })
			);
		} catch (err) {
			error = (err as Error).message;
		}
	}
</script>

<svelte:head><title>{$t`Connect a print agent`} · Standblatt</title></svelte:head>

<AuthShell>
	<h1 class="text-2xl font-semibold tracking-tight">{$t`Connect a print agent`}</h1>
	{#if me.isPending}
		<p class="mt-2 text-sm text-muted-foreground">…</p>
	{:else if !me.data}
		<p class="mt-2 mb-6 text-sm text-muted-foreground">{$t`Log in to connect the computer with the code to a range of your club.`}</p>
		<Button href="/login?next={encodeURIComponent(page.url.pathname + page.url.search)}">{$t`Log in`}</Button>
	{:else if connected}
		{@const range = ranges.data?.find((r) => r.id === rangeId)}
		<p class="mt-2 mb-6 text-sm" role="status">{$t`Connected. ${connected.name} now prints for ${range?.name ?? ''}; the app shows it in a moment.`}</p>
		<Button href="/t/{tenantId}/ranges/{rangeId}/printing">{$t`Set up the printers`}</Button>
	{:else if clubs.length === 0}
		<p class="mt-2 text-sm text-muted-foreground">{$t`You need the permission to manage ranges in a club to connect a print agent.`}</p>
	{:else}
		<p class="mt-2 mb-6 text-sm text-muted-foreground">{$t`Enter the code the Standblatt app shows on the computer next to the printers.`}</p>
		<form class="grid gap-4" onsubmit={connect}>
			<FormField label={$t`Code`}>
				<Input name="code" bind:value={code} autocomplete="off" autocapitalize="characters" placeholder="K7QD-4MZP" required class="font-mono tracking-widest uppercase" />
			</FormField>
			{#if pairing.data}
				<p class="text-sm" role="status">
					{$t`${pairing.data.name} (${osLabel(pairing.data.os)}, ${pairing.data.version}) wants to connect.`}
				</p>
			{:else if pairing.error}
				<FormError message={pairing.error.message} />
			{/if}
			<FormField label={$t`Club`}>
				<NativeSelect.Root bind:value={tenantId} name="club">
					{#each clubs as c (c.tenantId)}<NativeSelect.Option value={c.tenantId}>{c.tenantName}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
			</FormField>
			<FormField label={$t`Range`}>
				<NativeSelect.Root bind:value={rangeId} name="range">
					{#each ranges.data ?? [] as r (r.id)}<NativeSelect.Option value={r.id}>{r.name}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
			</FormField>
			<FormField label={$t`Name`}>
				<Input name="name" bind:value={name} maxlength={100} required />
			</FormField>
			<FormError message={error} />
			<Button type="submit" class="w-full" disabled={!pairing.data || !rangeId}>{$t`Connect`}</Button>
		</form>
	{/if}
</AuthShell>
