<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import AuthShell from '$lib/AuthShell.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { ApiError } from '$lib/api';
	import { authOptionsQuery, meKey } from '$lib/me';
	import { forgetMy, myKey } from '$lib/my.svelte';

	const queryClient = useQueryClient();
	const options = authOptionsQuery();
	let name = $state('');
	let password = $state('');
	let code = $state('');
	let remember = $state(false);
	let needCode = $state(false);
	let error = $state('');
	let busy = $state(false);

	/** Logs in; asks for the two-factor code when the server requires it. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		busy = true;
		error = '';
		try {
			const me = await call(api.POST('/api/auth/login', { body: { name, password, code: needCode ? code : undefined, remember } }));
			forgetMy();
			queryClient.removeQueries({ queryKey: myKey });
			queryClient.setQueryData(meKey, me);
			await goto('/');
		} catch (err) {
			if (err instanceof ApiError && err.status === 428) {
				needCode = true;
				return;
			}
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

<AuthShell>
	<h1 class="text-2xl font-semibold tracking-tight">{$t`Log in`}</h1>
	<p class="mb-6 text-sm text-muted-foreground">{$t`With your email address, or the user name your club admin gave you.`}</p>
	<form class="grid gap-4" onsubmit={submit}>
		<FormField label={$t`Email or user name`}><Input name="name" autocomplete="username" bind:value={name} required /></FormField>
		<FormField label={$t`Password`}><Input name="password" type="password" autocomplete="current-password" bind:value={password} required /></FormField>
		{#if needCode}
			<FormField label={$t`Code from your authenticator app`}>
				<!-- svelte-ignore a11y_autofocus -->
				<Input name="code" inputmode="numeric" autocomplete="one-time-code" pattern={"[0-9 ]{6,8}"} bind:value={code} required autofocus />
			</FormField>
		{/if}
		<label class="flex items-center gap-2 text-sm"><Checkbox bind:checked={remember} name="remember" />{$t`Stay logged in for 90 days`}</label>
		<FormError message={error} />
		<Button type="submit" class="w-full" disabled={busy}>{$t`Log in`}</Button>
	</form>
	<div class="mt-6 flex flex-wrap justify-between gap-2 text-sm">
		<a href="/forgot-password" class="text-muted-foreground hover:text-foreground">{$t`Forgot your password?`}</a>
		{#if options.data?.signup}
			<a href="/signup" class="font-medium underline-offset-4 hover:underline">{$t`Create an account`}</a>
		{/if}
	</div>
</AuthShell>
