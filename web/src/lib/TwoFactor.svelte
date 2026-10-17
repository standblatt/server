<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import ShieldCheck from '@lucide/svelte/icons/shield-check';
	import { toast } from 'svelte-sonner';
	import { api, call } from './api';
	import { Button } from './components/ui/button';
	import * as Dialog from './components/ui/dialog';
	import { Input } from './components/ui/input';
	import FormError from './FormError.svelte';
	import FormField from './FormField.svelte';
	import { authOptionsQuery, meKey, meQuery } from './me';

	const queryClient = useQueryClient();
	const me = meQuery();
	const options = authOptionsQuery();
	const on = $derived(!!me.data?.totp);

	let setup = $state<{ secret: string; url: string; qr: string }>();
	let code = $state('');
	let password = $state('');
	let error = $state('');
	let offOpen = $state(false);

	/** Starts the TOTP setup and fetches the secret and QR code. */
	async function start() {
		error = code = '';
		setup = await call(api.POST('/api/me/totp/setup'));
	}
	/** Switches the second factor on with the setup secret and the entered code. */
	async function enable(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			await call(api.PUT('/api/me/totp', { body: { secret: setup!.secret, code } }));
			setup = undefined;
			await queryClient.invalidateQueries({ queryKey: meKey });
			toast.success($t`Second factor switched on`);
		} catch (err) {
			error = (err as Error).message;
		}
	}
	/** Switches the second factor off after checking the password. */
	async function disable(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			await call(api.POST('/api/me/totp/disable', { body: { password } }));
			offOpen = false;
			password = '';
			await queryClient.invalidateQueries({ queryKey: meKey });
			toast.success($t`Second factor switched off`);
		} catch (err) {
			error = (err as Error).message;
		}
	}
</script>

{#if options.data?.totp}
	<section class="rounded-xl border bg-card p-4 shadow-sm" aria-label={$t`Second factor`}>
		<div class="flex flex-wrap items-center justify-between gap-2">
			<div>
				<h2 class="flex items-center gap-2 font-semibold"><ShieldCheck class="size-4" />{$t`Second factor`}</h2>
				<p class="text-sm text-muted-foreground">
					{on ? $t`Logging in asks for a code from your authenticator app.` : $t`Protect your login with a code from an authenticator app (TOTP).`}
				</p>
			</div>
			{#if on}
				<Button variant="outline" size="sm" onclick={() => {
						error = '';
						offOpen = true;
					}}>{$t`Switch off`}</Button>
			{:else}
				<Button size="sm" onclick={start}>{$t`Set up`}</Button>
			{/if}
		</div>
	</section>
{/if}

<Dialog.Root open={!!setup} onOpenChange={(o) => !o && (setup = undefined)}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>{$t`Set up the authenticator app`}</Dialog.Title>
			<Dialog.Description>{$t`Scan the code with Google Authenticator, Aegis, 1Password or a similar app, then enter the code it shows.`}</Dialog.Description>
		</Dialog.Header>
		{#if setup}
			<form class="grid gap-3" onsubmit={enable}>
				<img src={setup.qr} alt={$t`QR code for the authenticator app`} width="240" height="240" class="mx-auto rounded-md bg-white p-2" />
				<p class="text-center text-xs break-all text-muted-foreground">{$t`Cannot scan? Enter this key:`} <code class="select-all">{setup.secret}</code></p>
				<FormField label={$t`Code from the app`}><Input bind:value={code} required name="code" inputmode="numeric" autocomplete="one-time-code" pattern={"[0-9 ]{6,8}"} /></FormField>
				<FormError message={error} />
				<Button type="submit">{$t`Switch on`}</Button>
			</form>
		{/if}
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={offOpen}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header><Dialog.Title>{$t`Switch off the second factor?`}</Dialog.Title></Dialog.Header>
		<form class="grid gap-3" onsubmit={disable}>
			<FormField label={$t`Your password`}><Input type="password" bind:value={password} required name="password" autocomplete="current-password" /></FormField>
			<FormError message={error} />
			<Button type="submit" variant="destructive">{$t`Switch off`}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
