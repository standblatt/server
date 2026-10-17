<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { api, call } from '$lib/api';
	import AuthShell from '$lib/AuthShell.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';

	let email = $state('');
	let sent = $state(false);
	let error = $state('');

	/** Requests a password reset mail and shows the confirmation. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			await call(api.POST('/api/auth/forgot', { body: { email } }));
			sent = true;
		} catch (err) {
			error = (err as Error).message;
		}
	}
</script>

<svelte:head><title>{$t`Forgot your password?`} · Standblatt</title></svelte:head>

<AuthShell>
	<h1 class="text-2xl font-semibold tracking-tight">{$t`Forgot your password?`}</h1>
	{#if sent}
		<p class="mt-2 text-sm text-muted-foreground">{$t`If an account uses ${email}, we sent a link to choose a new password. Check your inbox.`}</p>
	{:else}
		<p class="mb-6 text-sm text-muted-foreground">{$t`Enter your email address and we send you a link to choose a new one. Club accounts without email: ask your club admin.`}</p>
		<form class="grid gap-4" onsubmit={submit}>
			<FormField label={$t`Email`}><Input name="email" type="email" autocomplete="email" bind:value={email} required /></FormField>
			<FormError message={error} />
			<Button type="submit" class="w-full">{$t`Send link`}</Button>
		</form>
	{/if}
	<p class="mt-6 text-sm"><a class="text-muted-foreground hover:text-foreground" href="/login">{$t`Back to log in`}</a></p>
</AuthShell>
