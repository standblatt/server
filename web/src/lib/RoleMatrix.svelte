<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call, type Permission, type Role } from './api';
	import { Button } from './components/ui/button';
	import { Checkbox } from './components/ui/checkbox';
	import * as Dialog from './components/ui/dialog';
	import { Input } from './components/ui/input';
	import * as Table from './components/ui/table';
	import { ask } from './confirm.svelte';
	import FormError from './FormError.svelte';
	import FormField from './FormField.svelte';
	import { roleLabel } from './labels';
	import { meKey, membershipOf, meQuery } from './me';
	import { hasClub } from './modules';
	import { permissionKeys, permissionUI } from './permissions';
	import { keys, rolesQuery } from './tenant';

	let { tenantId }: { tenantId: number } = $props();
	const queryClient = useQueryClient();
	const roles = rolesQuery(() => tenantId);
	const me = meQuery();
	const club = $derived(membershipOf(me.data, tenantId));
	const shown = $derived(permissionKeys.filter((p) => !permissionUI[p].module || hasClub(club, permissionUI[p].module!)));

	/** Name of a role: its own or the translated default name. */
	const roleName = (r: Role) => r.name || $t(roleLabel(r.key));

	let editing = $state<Role | 'new' | null>(null);
	let name = $state('');
	let error = $state('');

	/** Reloads the roles and the own permissions. */
	async function refresh() {
		await queryClient.invalidateQueries({ queryKey: keys.roles(tenantId) });
		await queryClient.invalidateQueries({ queryKey: meKey });
	}

	/** Gives a role a permission or takes it away. */
	async function toggle(r: Role, p: Permission, on: boolean) {
		const permissions = on ? [...r.permissions, p] : r.permissions.filter((x) => x !== p);
		try {
			await call(api.PUT('/api/tenants/{tenantId}/roles/{roleId}', { params: { path: { tenantId, roleId: r.id } }, body: { name: r.name, permissions } }));
		} catch (err) {
			toast.error((err as Error).message);
		}
		await refresh();
	}

	/** Opens the name dialog for a new role or for renaming one. */
	function open(r: Role | 'new') {
		editing = r;
		name = r === 'new' ? '' : r.name;
		error = '';
	}

	/** Creates the role or saves its new name. */
	async function save(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			if (editing === 'new') {
				await call(api.POST('/api/tenants/{tenantId}/roles', { params: { path: { tenantId } }, body: { name, permissions: ['events.view'] } }));
				toast.success($t`Role ${name} created`);
			} else if (editing) {
				await call(api.PUT('/api/tenants/{tenantId}/roles/{roleId}', { params: { path: { tenantId, roleId: editing.id } }, body: { name, permissions: editing.permissions } }));
			}
			editing = null;
			await refresh();
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Deletes a role after confirmation. */
	async function remove(r: Role) {
		if (!(await ask({ title: $t`Delete the role ${roleName(r)}?`, action: $t`Delete`, destructive: true }))) return;
		try {
			await call(api.DELETE('/api/tenants/{tenantId}/roles/{roleId}', { params: { path: { tenantId, roleId: r.id } } }));
			await refresh();
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
</script>

<section class="rounded-lg border bg-card" aria-label={$t`Roles`}>
	<div class="flex flex-wrap items-center justify-between gap-3 border-b px-4 py-3">
		<div>
			<h2 class="text-sm font-semibold">{$t`Roles`}</h2>
			<p class="text-xs text-muted-foreground">{$t`What the members of each role may do. Admins may always do everything.`}</p>
		</div>
		<Button size="sm" variant="outline" onclick={() => open('new')}><Plus />{$t`New role`}</Button>
	</div>
	<div class="overflow-x-auto" data-swipe>
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>{$t`Permission`}</Table.Head>
					{#each roles.data ?? [] as r (r.id)}
						<Table.Head class="text-center">
							<div class="flex items-center justify-center gap-0.5 whitespace-nowrap">
								<span>{roleName(r)}</span>
								<Button variant="ghost" size="icon" class="size-7" aria-label={$t`Rename ${roleName(r)}`} onclick={() => open(r)}><Pencil /></Button>
								{#if r.key !== 'admin'}<Button variant="ghost" size="icon" class="size-7" aria-label={$t`Delete ${roleName(r)}`} onclick={() => remove(r)}><Trash /></Button>{/if}
							</div>
							<div class="text-xs font-normal text-muted-foreground">{$plural(r.members, { one: '# member', other: '# members' })}</div>
						</Table.Head>
					{/each}
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each shown as p, i (p)}
					{#if i === 0 || permissionUI[shown[i - 1]].group !== permissionUI[p].group}
						<Table.Row class="bg-muted/40 hover:bg-muted/40"><Table.Cell colspan={(roles.data?.length ?? 0) + 1} class="py-1.5 text-xs font-semibold uppercase text-muted-foreground">{$t(permissionUI[p].group)}</Table.Cell></Table.Row>
					{/if}
					<Table.Row>
						<Table.Cell class="min-w-56 whitespace-normal">
							<div class="font-medium">{$t(permissionUI[p].label)}</div>
							<div class="text-xs text-muted-foreground">{$t(permissionUI[p].description)}</div>
						</Table.Cell>
						{#each roles.data ?? [] as r (r.id)}
							<Table.Cell class="text-center">
								<Checkbox
									checked={r.permissions.includes(p)}
									disabled={r.key === 'admin'}
									onCheckedChange={(v) => toggle(r, p, v)}
									aria-label={$t`${$t(permissionUI[p].label)} for ${roleName(r)}`}
								/>
							</Table.Cell>
						{/each}
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
</section>

<Dialog.Root open={editing !== null} onOpenChange={(o) => !o && (editing = null)}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>{editing === 'new' ? $t`New role` : $t`Rename role`}</Dialog.Title>
			<Dialog.Description>
				{editing === 'new' ? $t`A new role may only see the events until you give it more permissions.` : $t`Leave the name empty to keep the name every club starts with.`}
			</Dialog.Description>
		</Dialog.Header>
		<form class="grid gap-4" onsubmit={save}>
			<FormField label={$t`Name`}>
				<Input bind:value={name} maxlength={60} placeholder={editing && editing !== 'new' && editing.key ? $t(roleLabel(editing.key)) : ''} required={editing === 'new' || !(editing && editing.key)} />
			</FormField>
			<FormError message={error} />
			<Dialog.Footer><Button type="submit">{$t`Save`}</Button></Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
