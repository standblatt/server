<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { api, call, type Member, type Role } from '$lib/api';
	import { roleLabel } from '$lib/labels';
	import Mail from '@lucide/svelte/icons/mail';
	import Plus from '@lucide/svelte/icons/plus';
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import * as Table from '$lib/components/ui/table';
	import { ask } from '$lib/confirm.svelte';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import RoleMatrix from '$lib/RoleMatrix.svelte';
	import { meKey } from '$lib/me';
	import { keys, rolesQuery, usersQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const users = usersQuery(() => tenantId);
	const roles = rolesQuery(() => tenantId);
	/** Name of a role: its own or the translated default name. */
	const roleName = (r?: Role) => r?.name || $t(roleLabel(r?.key));
	/** The helper role, preselected for new users. */
	const helperRole = $derived(roles.data?.find((r) => r.key === 'staff')?.id ?? 0);

	let name = $state('');
	let role = $state(0);
	let password = $state('');
	let error = $state('');
	let adding = $state(false);
	let inviting = $state(false);
	let inviteText = $state('');
	let inviteRole = $state(0);
	let inviteError = $state('');
	let invited = $state<{ email: string; status: string }[]>([]);
	/** Splits the pasted list into addresses (commas, semicolons, whitespace and newlines all separate). */
	const inviteEmails = $derived(inviteText.split(/[\s,;]+/).map((x) => x.trim()).filter(Boolean));

	/** Sends the invitations and shows what happened to every address. */
	async function invite(e: SubmitEvent) {
		e.preventDefault();
		inviteError = '';
		try {
			invited = await call(api.POST('/api/tenants/{tenantId}/invitations', { params: { path: { tenantId } }, body: { emails: inviteEmails, roleId: inviteRole || helperRole } }));
			inviteText = '';
			await refresh();
		} catch (err) {
			inviteError = (err as Error).message;
		}
	}
	/** Label of an invitation result. */
	const inviteStatus = (s: string) => (s === 'invited' ? $t`invited by mail` : s === 'added' ? $t`added to the club` : s === 'member' ? $t`already a member` : $t`not an email address`);

	/** Reloads the user list and the own memberships. */
	async function refresh() {
		await queryClient.invalidateQueries({ queryKey: keys.users(tenantId) });
		await queryClient.invalidateQueries({ queryKey: keys.roles(tenantId) });
		await queryClient.invalidateQueries({ queryKey: meKey });
	}

	/** Adds a user to the club. */
	async function add(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			await call(
				api.POST('/api/tenants/{tenantId}/users', {
					params: { path: { tenantId } },
					body: { name, roleId: role || helperRole, password: password || undefined }
				})
			);
			toast.success($t`${name} added`);
			name = password = '';
			adding = false;
			await refresh();
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Changes a user's role. */
	async function setRole(m: Member, roleId: number) {
		try {
			await call(api.PUT('/api/tenants/{tenantId}/users/{userId}', { params: { path: { tenantId, userId: m.userId } }, body: { roleId } }));
			toast.success($t`${m.name} is now ${roleName(roles.data?.find((r) => r.id === roleId))}`);
		} catch (err) {
			toast.error((err as Error).message);
		}
		await refresh();
	}

	/** Removes a user from the club after confirmation. */
	async function remove(m: Member) {
		if (!(await ask({ title: $t`Remove ${m.name}?`, description: $t`The user loses access to this club. The account itself stays.`, action: $t`Remove`, destructive: true }))) return;
		try {
			await call(api.DELETE('/api/tenants/{tenantId}/users/{userId}', { params: { path: { tenantId, userId: m.userId } } }));
			await refresh();
			toast.success($t`${m.name} removed`);
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
</script>

<PageHeader help="user-guide/users" title={$t`Users`} description={$t`Accounts with access to this club. What each role may do is set under Roles.`}>
	{#snippet actions()}
		<Button
			onclick={() => {
				inviting = true;
				inviteError = '';
				invited = [];
			}}><Mail />{$t`Invite by email`}</Button
		>
		<Button
			variant="outline"
			onclick={() => {
				adding = true;
				error = '';
			}}><Plus />{$t`Add user`}</Button
		>
	{/snippet}
</PageHeader>

<div class="overflow-hidden rounded-lg border bg-card">
	<Table.Root>
		<Table.Header><Table.Row><Table.Head>{$t`User`}</Table.Head><Table.Head>{$t`Role`}</Table.Head><Table.Head></Table.Head></Table.Row></Table.Header>
		<Table.Body>
			{#each users.data ?? [] as m (m.userId)}
				<Table.Row>
					<Table.Cell>
						<div class="font-medium">{m.name}</div>
						{#if m.email && m.email !== m.name}<div class="text-xs text-muted-foreground">{m.email}</div>{/if}
					</Table.Cell>
					<Table.Cell>
						<NativeSelect.Root size="sm" value={m.roleId} onchange={(e) => setRole(m, Number(e.currentTarget.value))} aria-label={$t`Role of ${m.name}`}>
							{#each roles.data ?? [] as r (r.id)}<NativeSelect.Option value={r.id}>{roleName(r)}</NativeSelect.Option>{/each}
						</NativeSelect.Root>
					</Table.Cell>
					<Table.Cell class="text-right"><Button variant="ghost" size="sm" class="text-destructive" onclick={() => remove(m)}>{$t`Remove`}</Button></Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>

<div class="mt-6"><RoleMatrix {tenantId} /></div>

<Dialog.Root bind:open={inviting}>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>{$t`Invite by email`}</Dialog.Title>
			<Dialog.Description>{$t`New people get a mail with a link to set their password; existing accounts are added to the club and told by mail.`}</Dialog.Description>
		</Dialog.Header>
		<form class="grid gap-4" onsubmit={invite}>
			<FormField label={$t`Email addresses`} hint={$t`Several addresses separated by commas or one per line.`}>
				<textarea class="min-h-28 rounded-md border bg-background px-3 py-2 text-sm" bind:value={inviteText} name="emails" required placeholder="anna@example.org, bert@example.org"></textarea>
			</FormField>
			<FormField label={$t`Role`}>
				<NativeSelect.Root class="w-full" value={inviteRole || helperRole} onchange={(e) => (inviteRole = Number(e.currentTarget.value))}>
					{#each roles.data ?? [] as r (r.id)}<NativeSelect.Option value={r.id}>{roleName(r)}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
			</FormField>
			{#if invited.length}
				<ul class="divide-y rounded-md border text-sm" data-testid="invite-results">
					{#each invited as r (r.email)}
						<li class="flex items-center justify-between gap-3 px-3 py-1.5"><span class="truncate">{r.email}</span><span class="text-xs text-muted-foreground">{inviteStatus(r.status)}</span></li>
					{/each}
				</ul>
			{/if}
			<FormError message={inviteError} />
			<Dialog.Footer><Button type="submit" disabled={!inviteEmails.length}>{$t`Send invitations`}</Button></Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={adding}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>{$t`Add user`}</Dialog.Title>
			<Dialog.Description>{$t`Existing accounts are added by name; new accounts need a password.`}</Dialog.Description>
		</Dialog.Header>
		<form class="grid gap-4" onsubmit={add}>
			<FormField label={$t`User name`}><Input bind:value={name} autocomplete="off" required /></FormField>
			<FormField label={$t`Role`}>
				<NativeSelect.Root class="w-full" value={role || helperRole} onchange={(e) => (role = Number(e.currentTarget.value))}>
					{#each roles.data ?? [] as r (r.id)}<NativeSelect.Option value={r.id}>{roleName(r)}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
			</FormField>
			<FormField label={$t`Password`} hint={$t`New users only, at least 10 characters`}><Input type="password" bind:value={password} autocomplete="new-password" /></FormField>
			<FormError message={error} />
			<Dialog.Footer><Button type="submit">{$t`Add`}</Button></Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
