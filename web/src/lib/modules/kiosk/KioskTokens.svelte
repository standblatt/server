<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { intlLocale, localeNames } from '$lib/locale.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Copy from '@lucide/svelte/icons/copy';
	import MonitorPlay from '@lucide/svelte/icons/monitor-play';
	import { toast } from 'svelte-sonner';
	import { api, call, type KioskToken } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import { ask } from '$lib/confirm.svelte';
	import FormError from '$lib/FormError.svelte';
	import { keys } from '$lib/tenant';
	import type { EventProps } from '..';

	let { tenantId, event }: EventProps = $props();
	const eventId = $derived(event.id);
	const queryClient = useQueryClient();
	const tokens = createQuery(() => ({
		queryKey: keys.kioskTokens(tenantId, eventId),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/kiosk-tokens', { params: { path: { tenantId, eventId } } }))
	}));

	const ranges = $derived(event.ranges);
	let name = $state('');
	let theme = $state('dark');
	let language = $state('');
	let rangeId = $state(0);
	let url = $state('');
	let error = $state('');

	/** Creates a kiosk token and builds its display link with theme, language and range. */
	async function create(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			const r = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/kiosk-tokens', { params: { path: { tenantId, eventId } }, body: { name } }));
			url = `${location.origin}/kiosk/${r.secret}?theme=${theme}${language ? `&language=${language}` : ''}${rangeId ? `&range=${rangeId}` : ''}`;
			name = '';
			await queryClient.invalidateQueries({ queryKey: keys.kioskTokens(tenantId, eventId) });
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Asks for confirmation, then revokes a kiosk token. */
	async function revoke(k: KioskToken) {
		if (!(await ask({ title: $t`Revoke kiosk display?`, description: $t`"${k.name}" stops working immediately.`, action: $t`Revoke`, destructive: true }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/kiosk-tokens/{tokenId}', { params: { path: { tenantId, tokenId: k.id } } }));
		await queryClient.invalidateQueries({ queryKey: keys.kioskTokens(tenantId, eventId) });
	}
</script>

<div class="border-b px-4 py-3">
	<h2 class="text-sm font-semibold">{$t`Kiosk displays`}</h2>
	<p class="text-xs text-muted-foreground">{$t`A kiosk link shows the live rankings of this event without login. Anyone with the link can see them.`}</p>
</div>
<ul class="divide-y">
	{#each tokens.data ?? [] as k (k.id)}
		<li class="flex flex-wrap items-center gap-3 px-4 py-2.5 text-sm">
			<MonitorPlay class="size-4 text-muted-foreground" />
			<span class="font-medium">{k.name}</span>
			<span class="text-muted-foreground">{k.lastUsedAt ? $t`last seen ${new Date(k.lastUsedAt).toLocaleString(intlLocale())}` : $t`never used`}</span>
			<Button variant="ghost" size="sm" class="ml-auto text-destructive" onclick={() => revoke(k)}>{$t`revoke`}</Button>
		</li>
	{:else}
		<li class="px-4 py-3 text-sm text-muted-foreground">{$t`No kiosk displays.`}</li>
	{/each}
</ul>
<div class="space-y-3 border-t p-4">
	<form class="flex max-w-2xl flex-wrap gap-2 sm:flex-nowrap" onsubmit={create}>
		<Input placeholder={$t`Display name, e.g. Hall 1`} bind:value={name} required aria-label={$t`Display name`} />
		<NativeSelect.Root class="w-28 shrink-0" bind:value={theme} aria-label={$t`Theme`}>
			<NativeSelect.Option value="dark">{$t`Dark`}</NativeSelect.Option>
			<NativeSelect.Option value="light">{$t`Light`}</NativeSelect.Option>
		</NativeSelect.Root>
		<NativeSelect.Root class="w-40 shrink-0" bind:value={language} aria-label={$t`Language`}>
			<NativeSelect.Option value="">{$t`Event language`}</NativeSelect.Option>
			{#each Object.entries(localeNames) as [code, label] (code)}<NativeSelect.Option value={code}>{label}</NativeSelect.Option>{/each}
		</NativeSelect.Root>
		{#if ranges.length > 1}
			<NativeSelect.Root class="w-40 shrink-0" bind:value={rangeId} aria-label={$t`Range`}>
				<NativeSelect.Option value={0}>{$t`All ranges`}</NativeSelect.Option>
				{#each ranges as r (r.id)}<NativeSelect.Option value={r.id}>{r.name}</NativeSelect.Option>{/each}
			</NativeSelect.Root>
		{/if}
		<Button type="submit" variant="outline">{$t`Create kiosk link`}</Button>
	</form>
	{#if url}
		<div class="flex flex-wrap items-center gap-2 rounded-md border border-emerald-500/40 bg-emerald-500/10 px-3 py-2 text-sm" role="status">
			<span>{$t`Kiosk link (shown only now):`}</span>
			<a class="font-medium break-all underline" href={url} target="_blank" data-testid="kiosk-url">{url}</a>
			<Button
				variant="ghost"
				size="sm"
				onclick={() => {
					navigator.clipboard.writeText(url);
					toast.success($t`Link copied`);
				}}><Copy />{$t`copy`}</Button
			>
		</div>
	{/if}
	<FormError message={error} />
</div>
