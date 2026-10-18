<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import AuthShell from '$lib/AuthShell.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { roleLabel } from '$lib/labels';
	import { meKey, meQuery } from '$lib/me';

	const token = $derived(page.params.token ?? '');
	const queryClient = useQueryClient();
	const invitation = createQuery(() => ({
		queryKey: ['invitation', token],
		queryFn: () => call(api.GET('/api/invitations/{token}', { params: { path: { token } } })),
		retry: false
	}));
	const me = meQuery();
	const role = $derived(invitation.data ? invitation.data.roleName || $t(roleLabel(invitation.data.roleKey)) : '');
	let password = $state('');
	let error = $state('');

	/** Accepts the invitation, with a new account when nobody is logged in, and opens the start page. */
	async function accept(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			const body = me.data ? undefined : { password };
			const updated = await call(api.POST('/api/invitations/{token}/accept', { params: { path: { token } }, body }));
			queryClient.setQueryData(meKey, updated);
			await goto('/');
		} catch (err) {
			error = (err as Error).message;
		}
	}
</script>

<svelte:head><title>{$t`Invitation`} · Standblatt</title></svelte:head>

<AuthShell>
	<h1 class="text-2xl font-semibold tracking-tight">{$t`Invitation`}</h1>
	{#if invitation.error}
		<p class="mt-2 mb-6 text-sm text-muted-foreground">{invitation.error.message}</p>
		<Button href="/login">{$t`Log in`}</Button>
	{:else if invitation.data && !me.isPending}
		{@const inv = invitation.data}
		<p class="mt-2 mb-6 text-sm text-muted-foreground">{$t`${inv.club} invites you to Standblatt as ${role}.`}</p>
		{#if me.data}
			<form class="grid gap-4" onsubmit={accept}>
				<p class="text-sm">{$t`You join as ${me.data.user.name}.`}</p>
				<FormError message={error} />
				<Button type="submit" class="w-full">{$t`Accept the invitation`}</Button>
			</form>
		{:else if inv.account}
			<p class="mb-4 text-sm">{$t`There is an account with ${inv.email}. Log in to accept the invitation.`}</p>
			<Button href="/login?next={encodeURIComponent(page.url.pathname)}">{$t`Log in`}</Button>
		{:else}
			<form class="grid gap-4" onsubmit={accept}>
				<FormField label={$t`Email`}><Input value={inv.email} disabled /></FormField>
				<FormField label={$t`Choose a password`} hint={$t`At least 10 characters.`}>
					<Input name="password" type="password" autocomplete="new-password" minlength={10} bind:value={password} required />
				</FormField>
				<FormError message={error} />
				<Button type="submit" class="w-full">{$t`Accept and create the account`}</Button>
			</form>
		{/if}
	{/if}
</AuthShell>
