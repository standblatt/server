<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import AuthShell from '$lib/AuthShell.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { currentLocale } from '$lib/locale.svelte';
	import { authOptionsQuery, meKey } from '$lib/me';

	const queryClient = useQueryClient();
	const options = authOptionsQuery();
	let firstname = $state('');
	let lastname = $state('');
	let email = $state(page.url.searchParams.get('email') ?? '');
	let password = $state('');
	let error = $state('');
	let busy = $state(false);

	/** Creates the account and opens the training log. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		busy = true;
		error = '';
		try {
			const me = await call(api.POST('/api/auth/signup', { body: { firstname, lastname, email, password, locale: currentLocale() } }));
			queryClient.setQueryData(meKey, me);
			await goto('/my');
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

<svelte:head><title>{$t`Create an account`} · Standblatt</title></svelte:head>

<AuthShell>
	<h1 class="text-2xl font-semibold tracking-tight">{$t`Create an account`}</h1>
	{#if options.data && !options.data.signup}
		<p class="mt-2 text-sm text-muted-foreground">{$t`Sign-up is not open on this server.`} <a class="underline" href="/login">{$t`Log in`}</a></p>
	{:else}
		<p class="mb-6 text-sm text-muted-foreground">{$t`Your personal training log, ready in a minute.`}</p>
		<form class="grid gap-4" onsubmit={submit}>
			<div class="grid grid-cols-2 gap-3">
				<FormField label={$t`First name`}><Input name="firstname" autocomplete="given-name" bind:value={firstname} required /></FormField>
				<FormField label={$t`Last name`}><Input name="lastname" autocomplete="family-name" bind:value={lastname} required /></FormField>
			</div>
			<FormField label={$t`Email`}><Input name="email" type="email" autocomplete="email" bind:value={email} required /></FormField>
			<FormField label={$t`Password`} hint={$t`At least 10 characters.`}>
				<Input name="password" type="password" autocomplete="new-password" minlength={10} bind:value={password} required />
			</FormField>
			<FormError message={error} />
			<Button type="submit" class="w-full" disabled={busy}>{$t`Create account`}</Button>
		</form>
		<p class="mt-6 text-sm text-muted-foreground">{$t`Already have an account?`} <a class="font-medium text-foreground underline-offset-4 hover:underline" href="/login">{$t`Log in`}</a></p>
	{/if}
</AuthShell>
