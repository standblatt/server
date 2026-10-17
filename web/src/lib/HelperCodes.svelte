<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { intlLocale } from './locale.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Copy from '@lucide/svelte/icons/copy';
	import HandHelping from '@lucide/svelte/icons/hand-helping';
	import { toast } from 'svelte-sonner';
	import { api, call, type HelperCode } from './api';
	import { Button } from './components/ui/button';
	import { Input } from './components/ui/input';
	import * as NativeSelect from './components/ui/native-select';
	import { ask } from './confirm.svelte';
	import FormError from './FormError.svelte';
	import { helperRoleLabel } from './labels';
	import { eventQuery, keys } from './tenant';

	let { tenantId, eventId }: { tenantId: number; eventId: number } = $props();
	const queryClient = useQueryClient();
	const codes = createQuery(() => ({
		queryKey: keys.helperCodes(tenantId, eventId),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/helper-codes', { params: { path: { tenantId, eventId } } }))
	}));
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const ranges = $derived(event.data?.ranges ?? []);
	/** Name of a range of the event by ID. */
	const rangeName = (id?: number) => ranges.find((r) => r.id === id)?.name ?? '';

	let name = $state('');
	let role = $state<HelperCode['role']>('scorer');
	let rangeId = $state(0);
	let lanes = $state('');
	let url = $state('');
	let error = $state('');

	/** Creates a helper code (scorers get a range and lanes) and shows its link. */
	async function create(e: SubmitEvent) {
		e.preventDefault();
		error = url = '';
		try {
			const r = await call(
				api.POST('/api/tenants/{tenantId}/events/{eventId}/helper-codes', {
					params: { path: { tenantId, eventId } },
					body: { name, role, rangeId: role === 'scorer' ? rangeId || ranges[0]?.id : undefined, lanes: role === 'scorer' ? lanes : undefined }
				})
			);
			url = `${location.origin}/h/${r.secret}`;
			name = lanes = '';
			await queryClient.invalidateQueries({ queryKey: keys.helperCodes(tenantId, eventId) });
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Asks for confirmation, then revokes a helper code. */
	async function revoke(c: HelperCode) {
		if (!(await ask({ title: $t`Revoke helper link?`, description: $t`"${c.name}" stops working immediately.`, action: $t`Revoke`, destructive: true }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/helper-codes/{codeId}', { params: { path: { tenantId, codeId: c.id } } }));
		await queryClient.invalidateQueries({ queryKey: keys.helperCodes(tenantId, eventId) });
	}
</script>

<div class="border-b px-4 py-3">
	<h2 class="text-sm font-semibold">{$t`Helpers`}</h2>
	<p class="text-xs text-muted-foreground">
		{$t`A helper link opens a reduced desk for this event without login: scorers enter the results of their lanes, the ammunition desk issues ammunition, assigns lanes and prints lane sheets, results clerks enter results on every lane.`}
	</p>
</div>
<ul class="divide-y">
	{#each codes.data ?? [] as c (c.id)}
		<li class="flex flex-wrap items-center gap-3 px-4 py-2.5 text-sm">
			<HandHelping class="size-4 text-muted-foreground" />
			<span class="font-medium">{c.name}</span>
			<span class="text-muted-foreground">{$t(helperRoleLabel(c.role))}{c.role === 'scorer' ? ` · ${ranges.length > 1 ? `${rangeName(c.rangeId)} ` : ''}${c.lanes.join(', ')}` : ''}</span>
			<span class="text-muted-foreground">{c.lastUsedAt ? $t`last seen ${new Date(c.lastUsedAt).toLocaleString(intlLocale())}` : $t`never used`}</span>
			<Button variant="ghost" size="sm" class="ml-auto text-destructive" onclick={() => revoke(c)}>{$t`revoke`}</Button>
		</li>
	{:else}
		<li class="px-4 py-3 text-sm text-muted-foreground">{$t`No helper links.`}</li>
	{/each}
</ul>
<div class="space-y-3 border-t p-4">
	<form class="flex max-w-3xl flex-wrap gap-2" onsubmit={create}>
		<Input class="min-w-48 flex-1" placeholder={$t`Helper name, e.g. Scorer range A`} bind:value={name} required aria-label={$t`Helper name`} />
		<NativeSelect.Root class="w-44 shrink-0" bind:value={role} aria-label={$t`Role`}>
			{#each ['scorer', 'ammunition', 'results'] as r (r)}<NativeSelect.Option value={r}>{$t(helperRoleLabel(r))}</NativeSelect.Option>{/each}
		</NativeSelect.Root>
		{#if role === 'scorer'}
			{#if ranges.length > 1}
				<NativeSelect.Root class="w-40 shrink-0" bind:value={rangeId} aria-label={$t`Range`}>
					{#each ranges as r (r.id)}<NativeSelect.Option value={r.id}>{r.name}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
			{/if}
			<Input class="w-32 shrink-0" placeholder={$t`Lanes, e.g. 1-5, 8`} bind:value={lanes} required aria-label={$t`Lanes`} />
		{/if}
		<Button type="submit" variant="outline">{$t`Create helper link`}</Button>
	</form>
	{#if url}
		<div class="flex flex-wrap items-center gap-2 rounded-md border border-emerald-500/40 bg-emerald-500/10 px-3 py-2 text-sm" role="status">
			<span>{$t`Helper link (shown only now):`}</span>
			<a class="font-medium break-all underline" href={url} target="_blank" data-testid="helper-url">{url}</a>
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
