<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call } from './api';
	import ApiTokens from './ApiTokens.svelte';
	import { Button } from './components/ui/button';
	import * as Dialog from './components/ui/dialog';
	import { Input } from './components/ui/input';
	import FormError from './FormError.svelte';
	import FormField from './FormField.svelte';
	import { currentLocale, localeNames } from './locale.svelte';
	import { meKey, meQuery, setMyLocale } from './me';
	import { myKey } from './my.svelte';
	import ThemeToggle from './ThemeToggle.svelte';
	import TwoFactor from './TwoFactor.svelte';

	const queryClient = useQueryClient();
	const me = meQuery();
	const inClub = $derived((me.data?.tenants.length ?? 0) > 0);

	let emailOpen = $state(false);
	let email = $state('');
	let emailPassword = $state('');
	let emailError = $state('');
	let sent = $state(false);
	/** Changes the email address and sends a confirmation link to the new one. */
	async function changeEmail(e: SubmitEvent) {
		e.preventDefault();
		emailError = '';
		try {
			await call(api.PUT('/api/my/email', { body: { email, password: emailPassword } }));
			emailOpen = false;
			emailPassword = '';
			await Promise.all([queryClient.invalidateQueries({ queryKey: meKey }), queryClient.invalidateQueries({ queryKey: myKey })]);
			toast.success($t`We sent a confirmation link to ${email}.`);
		} catch (err) {
			emailError = (err as Error).message;
		}
	}
	/** Sends the email verification link again. */
	async function resend() {
		await call(api.POST('/api/my/email/verify'));
		sent = true;
		toast.success($t`We sent you a new link.`);
	}

	let currentPassword = $state('');
	let newPassword = $state('');
	let passwordError = $state('');
	/** Changes the password; other sessions of the user end. */
	async function changePassword(e: SubmitEvent) {
		e.preventDefault();
		passwordError = '';
		try {
			await call(api.PUT('/api/me/password', { body: { currentPassword, newPassword } }));
			currentPassword = newPassword = '';
			toast.success($t`Password changed. Other devices were logged out.`);
		} catch (err) {
			passwordError = (err as Error).message;
		}
	}

	let deleteOpen = $state(false);
	let deletePassword = $state('');
	let deleteError = $state('');
	/** Deletes the account after the password check and goes to the login. */
	async function deleteAccount(e: SubmitEvent) {
		e.preventDefault();
		deleteError = '';
		try {
			await call(api.DELETE('/api/my/account', { body: { password: deletePassword } }));
			queryClient.clear();
			await goto('/login');
		} catch (err) {
			deleteError = (err as Error).message;
		}
	}
</script>

{#if me.data}
	<section class="rounded-xl border bg-card p-4 shadow-sm">
		<h2 class="font-semibold">{$t`Login`}</h2>
		<dl class="mt-2 grid gap-2 text-sm sm:grid-cols-[8rem_1fr]">
			<dt class="text-muted-foreground">{$t`User name`}</dt>
			<dd>{me.data.user.name}</dd>
			<dt class="text-muted-foreground">{$t`Email`}</dt>
			<dd class="flex flex-wrap items-center gap-2">
				<span class="[overflow-wrap:anywhere]">{me.data.email || $t`No email address yet`}</span>
				{#if me.data.email}
					<span class="text-xs text-muted-foreground">· {me.data.emailVerified ? $t`confirmed` : $t`not confirmed`}</span>
					{#if !me.data.emailVerified}<Button size="xs" variant="ghost" disabled={sent} onclick={resend}>{$t`Send again`}</Button>{/if}
				{/if}
				<Button
					variant="outline"
					size="xs"
					onclick={() => {
						email = me.data?.email ?? '';
						emailOpen = true;
					}}>{me.data.email ? $t`Change` : $t`Add`}</Button
				>
			</dd>
		</dl>
		<p class="mt-2 text-xs text-muted-foreground">{$t`With a confirmed email address you can reset a forgotten password yourself.`}</p>
	</section>

	<section class="rounded-xl border bg-card p-4 shadow-sm">
		<h2 class="mb-3 font-semibold">{$t`Password`}</h2>
		<form class="grid gap-3 sm:grid-cols-2" onsubmit={changePassword}>
			<FormField label={$t`Current password`}><Input type="password" bind:value={currentPassword} required name="currentPassword" autocomplete="current-password" /></FormField>
			<FormField label={$t`New password`} hint={$t`At least 10 characters.`}><Input type="password" bind:value={newPassword} required minlength={10} name="newPassword" autocomplete="new-password" /></FormField>
			<FormError message={passwordError} class="sm:col-span-2" />
			<Button type="submit" class="w-fit">{$t`Change password`}</Button>
		</form>
	</section>

	<section class="grid gap-3 rounded-xl border bg-card p-4 text-sm shadow-sm">
		<h2 class="font-semibold">{$t`Language and theme`}</h2>
		<div class="flex flex-wrap items-center gap-3">
			<span class="flex-1">{$t`Language`}</span>
			<select class="rounded-md border bg-background px-2 py-1.5" value={currentLocale()} onchange={(e) => setMyLocale(queryClient, e.currentTarget.value)} aria-label={$t`Language`}>
				{#each Object.entries(localeNames) as [value, label] (value)}<option {value}>{label}</option>{/each}
			</select>
		</div>
		<div class="flex items-center gap-3"><span class="flex-1">{$t`Theme`}</span><ThemeToggle /></div>
	</section>

	<TwoFactor />
	<ApiTokens />

	<section class="grid gap-3 rounded-xl border bg-card p-4 text-sm shadow-sm">
		<h2 class="font-semibold">{$t`Delete account`}</h2>
		{#if inClub}
			<p class="text-muted-foreground">{$t`Your club manages this account. Ask your club admin to remove it.`}</p>
		{:else}
			<p class="text-muted-foreground">{$t`Your training log, goals and settings are deleted for good. Export your sessions first if you want to keep them.`}</p>
			<Button variant="destructive" size="sm" class="w-fit" onclick={() => (deleteOpen = true)}><Trash />{$t`Delete account`}</Button>
		{/if}
	</section>
{/if}

<Dialog.Root bind:open={emailOpen}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header><Dialog.Title>{$t`Email address`}</Dialog.Title></Dialog.Header>
		<form class="grid gap-3" onsubmit={changeEmail}>
			<FormField label={$t`Email`}><Input type="email" bind:value={email} required name="email" /></FormField>
			<FormField label={$t`Your password`}><Input type="password" bind:value={emailPassword} required name="password" autocomplete="current-password" /></FormField>
			<FormError message={emailError} />
			<Button type="submit">{$t`Send confirmation link`}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={deleteOpen}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>{$t`Delete your account?`}</Dialog.Title>
			<Dialog.Description>{$t`All your sessions, goals and settings are deleted for good. Export them first if you want to keep them.`}</Dialog.Description>
		</Dialog.Header>
		<form class="grid gap-3" onsubmit={deleteAccount}>
			<FormField label={$t`Your password`}><Input type="password" bind:value={deletePassword} required name="password" autocomplete="current-password" /></FormField>
			<FormError message={deleteError} />
			<Button type="submit" variant="destructive">{$t`Delete account`}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
