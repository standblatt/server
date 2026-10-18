<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import Copy from '@lucide/svelte/icons/copy';
	import KeyRound from '@lucide/svelte/icons/key-round';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { api, call, type Permission } from './api';
	import { Button } from './components/ui/button';
	import { Checkbox } from './components/ui/checkbox';
	import * as Dialog from './components/ui/dialog';
	import { Input } from './components/ui/input';
	import { ask } from './confirm.svelte';
	import { formatDate } from './format';
	import FormError from './FormError.svelte';
	import FormField from './FormField.svelte';
	import { meQuery } from './me';
	import { permissionKeys, permissionUI } from './permissions';

	type Training = 'training:read' | 'training:write';
	type Scope = Permission | Training;
	const key = ['me', 'tokens'];
	const queryClient = useQueryClient();
	const me = meQuery();
	const tokens = createQuery(() => ({ queryKey: key, queryFn: () => call(api.GET('/api/me/tokens')) }));
	const memberships = $derived(me.data?.tenants ?? []);
	const training = $derived<Record<Training, string>>({ 'training:read': $t`Read my training log`, 'training:write': $t`Log and change my training sessions` });

	let open = $state(false);
	let name = $state('');
	let clubs = $state<number[]>([]);
	let scopes = $state<Scope[]>([]);
	let days = $state<0 | 30 | 90 | 365>(90);
	let error = $state('');
	let secret = $state('');

	/** The club permissions the own roles have in the chosen clubs; only those can go into a token. */
	const available = $derived(permissionKeys.filter((p) => memberships.some((m) => clubs.includes(m.tenantId) && m.permissions.includes(p))));

	/** Named selections of scopes for common uses; the server checks them like any other selection. */
	const presets = $derived<{ label: string; scopes: Scope[] }[]>([
		{ label: $t`Rankings for a website`, scopes: ['events.view'] },
		{ label: $t`Training app`, scopes: ['training:read', 'training:write'] },
		{ label: $t`Read only`, scopes: ['events.view', 'reports.view', 'shooters.view', 'training:read'] }
	]);

	/** Label of a scope. */
	const scopeLabel = (s: string) => (s in training ? training[s as Training] : s in permissionUI ? $t(permissionUI[s as Permission].label) : s);
	/** Name of a club of the user by ID. */
	const clubName = (id: number) => memberships.find((m) => m.tenantId === id)?.tenantName ?? `#${id}`;

	/** Resets the form and opens the new token dialog. */
	function start() {
		name = '';
		clubs = memberships.length === 1 ? [memberships[0].tenantId] : [];
		scopes = [];
		days = 90;
		error = secret = '';
		open = true;
	}

	/** Takes the scopes of a preset that the chosen clubs allow. */
	function preset(list: Scope[]) {
		scopes = list.filter((s) => s in training || available.includes(s as Permission));
	}

	/** Adds or removes a club and drops scopes no chosen club allows any more. */
	function toggleClub(id: number, on: boolean) {
		clubs = on ? [...clubs, id] : clubs.filter((x) => x !== id);
		scopes = scopes.filter((s) => s in training || available.includes(s as Permission));
	}

	/** Adds or removes a scope. */
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
			const created = await call(api.POST('/api/me/tokens', { body: { name, scopes, tenants: clubs, expiresInDays: days } }));
			secret = created.token;
			await queryClient.invalidateQueries({ queryKey: key });
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Asks for confirmation, then revokes a token. */
	async function revoke(id: number, tokenName: string) {
		if (!(await ask({ title: $t`Revoke "${tokenName}"?`, description: $t`Scripts using this token stop working at once.`, action: $t`Revoke`, destructive: true }))) return;
		await call(api.DELETE('/api/me/tokens/{id}', { params: { path: { id } } }));
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
							{#if tok.tenants.length}{tok.tenants.map(clubName).join(', ')} · {/if}{tok.scopes.map(scopeLabel).join(', ')}
						</div>
						<div class="text-xs text-muted-foreground">
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
		<p class="text-xs text-muted-foreground">{$t`For your own scripts and tools: a token lets them use the API with the permissions you choose, without your password.`}</p>
	{/if}
</section>

<Dialog.Root bind:open>
	<Dialog.Content class="max-h-[90vh] overflow-y-auto sm:max-w-lg">
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
			<Dialog.Header>
				<Dialog.Title>{$t`New API token`}</Dialog.Title>
				<Dialog.Description>{$t`The token may do what you choose here, and never more than your role in the club.`}</Dialog.Description>
			</Dialog.Header>
			<form class="grid gap-4" onsubmit={create}>
				<FormField label={$t`Name`}><Input bind:value={name} required maxlength={100} name="tokenName" placeholder={$t`e.g. rankings on our website`} /></FormField>
				<div class="flex flex-wrap gap-2">
					{#each presets as p (p.label)}<Button size="sm" variant="outline" onclick={() => preset(p.scopes)}>{p.label}</Button>{/each}
				</div>
				{#if memberships.length}
					<fieldset class="grid gap-2 text-sm">
						<legend class="mb-1 font-medium">{$t`Clubs`}</legend>
						{#each memberships as m (m.tenantId)}
							<label class="flex items-center gap-2"><Checkbox checked={clubs.includes(m.tenantId)} onCheckedChange={(v) => toggleClub(m.tenantId, v)} />{m.tenantName}</label>
						{/each}
					</fieldset>
				{/if}
				<fieldset class="grid gap-2 text-sm">
					<legend class="mb-1 font-medium">{$t`Permissions`}</legend>
					{#each Object.keys(training) as s (s)}
						<label class="flex items-center gap-2"><Checkbox checked={scopes.includes(s as Scope)} onCheckedChange={(v) => toggle(s as Scope, v)} />{scopeLabel(s)}</label>
					{/each}
					{#each available as p (p)}
						<label class="flex items-start gap-2">
							<Checkbox class="mt-0.5" checked={scopes.includes(p)} onCheckedChange={(v) => toggle(p, v)} />
							<span>{$t(permissionUI[p].label)}<span class="block text-xs text-muted-foreground">{$t(permissionUI[p].description)}</span></span>
						</label>
					{/each}
					{#if memberships.length && !clubs.length}<p class="text-xs text-muted-foreground">{$t`Choose a club to give the token permissions there.`}</p>{/if}
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
