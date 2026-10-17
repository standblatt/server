<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import Copy from '@lucide/svelte/icons/copy';
	import KeyRound from '@lucide/svelte/icons/key-round';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { ask } from '$lib/confirm.svelte';
	import { formatDate } from '$lib/format';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';

	type Scope = 'training:read' | 'training:write';
	const key = ['my', 'tokens'];
	const queryClient = useQueryClient();
	const tokens = createQuery(() => ({ queryKey: key, queryFn: () => call(api.GET('/api/my/tokens')) }));
	const scopeLabels = $derived<Record<Scope, string>>({ 'training:read': $t`Read sessions`, 'training:write': $t`Log and change sessions` });

	let open = $state(false);
	let name = $state('');
	let scopes = $state<Scope[]>(['training:read']);
	let days = $state<0 | 30 | 90 | 365>(90);
	let error = $state('');
	let secret = $state('');

	/** Resets the form and opens the new token dialog. */
	function start() {
		name = '';
		scopes = ['training:read'];
		days = 90;
		error = secret = '';
		open = true;
	}

	/** Adds or removes a scope from the selected scopes. */
	function toggle(s: Scope, on: boolean) {
		scopes = on ? [...scopes, s] : scopes.filter((x) => x !== s);
	}

	/** Creates a token with at least one scope and shows its secret once. */
	async function create(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		if (!scopes.length) {
			error = $t`Choose at least one permission.`;
			return;
		}
		try {
			const created = await call(api.POST('/api/my/tokens', { body: { name, scopes, expiresInDays: days } }));
			secret = created.token;
			await queryClient.invalidateQueries({ queryKey: key });
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Asks for confirmation, then revokes a token. */
	async function revoke(id: number, tokenName: string) {
		if (!(await ask({ title: $t`Revoke "${tokenName}"?`, description: $t`Scripts using this token stop working at once.`, action: $t`Revoke`, destructive: true }))) return;
		await call(api.DELETE('/api/my/tokens/{id}', { params: { path: { id } } }));
		await queryClient.invalidateQueries({ queryKey: key });
		toast.success($t`Token revoked`);
	}

	/** Whether an expiry timestamp lies in the past. */
	const expired = (at?: string | null) => !!at && new Date(at) < new Date();
</script>

<section class="rounded-xl border bg-card p-4 shadow-sm" aria-label={$t`API tokens`}>
	<div class="mb-2 flex items-center justify-between gap-3">
		<h2 class="font-semibold">{$t`API tokens`}</h2>
		<Button variant="outline" size="sm" onclick={start}><Plus />{$t`New token`}</Button>
	</div>
	{#if tokens.data?.length}
		<ul class="divide-y text-sm">
			{#each tokens.data as tok (tok.id)}
				<li class="flex items-center gap-3 py-2" data-testid="api-token">
					<KeyRound class="size-4 shrink-0 text-muted-foreground" />
					<div class="min-w-0 flex-1">
						<div class="truncate font-medium">{tok.name}</div>
						<div class="text-xs text-muted-foreground">
							{tok.scopes.map((s) => scopeLabels[s as Scope] ?? s).join(', ')} ·
							{tok.lastUsedAt ? $t`last used ${formatDate(tok.lastUsedAt)}` : $t`never used`} ·
							{#if expired(tok.expiresAt)}<span class="text-destructive">{$t`expired`}</span>
							{:else if tok.expiresAt}{$t`expires ${formatDate(tok.expiresAt)}`}
							{:else}{$t`does not expire`}{/if}
						</div>
					</div>
					<Button variant="ghost" size="icon-xs" onclick={() => revoke(tok.id, tok.name)} aria-label={$t`Revoke ${tok.name}`}><Trash /></Button>
				</li>
			{/each}
		</ul>
	{:else}
		<p class="text-xs text-muted-foreground">{$t`For your own scripts and tools: a token lets them read or log sessions through the API, without your password.`}</p>
	{/if}
</section>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-md">
		{#if secret}
			<Dialog.Header>
				<Dialog.Title>{$t`Copy your token now`}</Dialog.Title>
				<Dialog.Description>{$t`It is shown only this once. Send it as "Authorization: Bearer <token>".`}</Dialog.Description>
			</Dialog.Header>
			<code class="block rounded-md bg-muted px-3 py-2 text-sm break-all" data-testid="api-token-secret">{secret}</code>
			<Button
				onclick={() => {
					navigator.clipboard.writeText(secret);
					toast.success($t`Copied`);
				}}><Copy />{$t`Copy`}</Button
			>
		{:else}
			<Dialog.Header><Dialog.Title>{$t`New API token`}</Dialog.Title></Dialog.Header>
			<form class="grid gap-3" onsubmit={create}>
				<FormField label={$t`Name`}><Input bind:value={name} required maxlength={100} name="tokenName" placeholder={$t`e.g. import script`} /></FormField>
				<fieldset class="grid gap-2 text-sm">
					<legend class="mb-1 font-medium">{$t`Permissions`}</legend>
					{#each Object.entries(scopeLabels) as [s, label] (s)}
						<label class="flex items-center gap-2">
							<input type="checkbox" checked={scopes.includes(s as Scope)} onchange={(e) => toggle(s as Scope, e.currentTarget.checked)} />
							{label}
						</label>
					{/each}
				</fieldset>
				<FormField label={$t`Expires`}>
					<select class="rounded-md border bg-background px-2 py-1.5 text-sm" bind:value={days} name="expiry">
						<option value={30}>{$t`in 30 days`}</option>
						<option value={90}>{$t`in 90 days`}</option>
						<option value={365}>{$t`in a year`}</option>
						<option value={0}>{$t`never`}</option>
					</select>
				</FormField>
				<FormError message={error} />
				<Button type="submit">{$t`Create token`}</Button>
			</form>
		{/if}
	</Dialog.Content>
</Dialog.Root>
