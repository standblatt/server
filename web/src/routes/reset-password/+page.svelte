<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { api, call } from '$lib/api';
	import AuthShell from '$lib/AuthShell.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';

	let password = $state('');
	let done = $state(false);
	let error = $state('');

	/** Sets the new password with the token from the link. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			await call(api.POST('/api/auth/reset', { body: { token: page.url.searchParams.get('token') ?? '', password } }));
			done = true;
		} catch (err) {
			error = (err as Error).message;
		}
	}
</script>

<svelte:head><title>{$t`New password`} · Standblatt</title></svelte:head>

<AuthShell>
	<h1 class="text-2xl font-semibold tracking-tight">{$t`New password`}</h1>
	{#if done}
		<p class="mt-2 mb-6 text-sm text-muted-foreground">{$t`Your password is changed. Log in with the new one.`}</p>
		<Button href="/login">{$t`Log in`}</Button>
	{:else}
		<form class="mt-6 grid gap-4" onsubmit={submit}>
			<FormField label={$t`New password`} hint={$t`At least 10 characters.`}>
				<Input name="password" type="password" autocomplete="new-password" minlength={10} bind:value={password} required />
			</FormField>
			<FormError message={error} />
			<Button type="submit" class="w-full">{$t`Save password`}</Button>
		</form>
	{/if}
</AuthShell>
