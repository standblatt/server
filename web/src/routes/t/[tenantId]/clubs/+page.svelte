<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { api, call, type Club } from '$lib/api';
	import MoreHorizontal from '@lucide/svelte/icons/ellipsis';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import { ask } from '$lib/confirm.svelte';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import Pager from '$lib/Pager.svelte';
	import { pageOf } from '$lib/format';
	import Search from '@lucide/svelte/icons/search';
	import { meQuery } from '$lib/me';
	import { clubsQuery, keys } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const me = meQuery();
	const clubs = clubsQuery(() => tenantId);
	const isAdmin = $derived(me.data?.tenants.find((t) => t.tenantId === tenantId)?.role === 'admin');

	let filter = $state('');
	const shown = $derived((clubs.data ?? []).filter((c) => c.name.toLowerCase().includes(filter.trim().toLowerCase())));
	const pageSize = 50;
	let pageNo = $state(0);

	let editId = $state<number | null>(null);
	let name = $state('');
	let ssvNumber = $state('');
	let error = $state('');

	/** Fills the form with a club, or empties it for a new one. */
	function edit(c: Club | null) {
		editId = c?.id ?? 0;
		name = c?.name ?? '';
		ssvNumber = c?.ssvNumber ?? '';
		error = '';
	}

	/** Creates or updates the club from the form. */
	async function save(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		const body = { name, ssvNumber };
		try {
			if (editId) {
				await call(api.PUT('/api/tenants/{tenantId}/clubs/{clubId}', { params: { path: { tenantId, clubId: editId } }, body }));
			} else {
				await call(api.POST('/api/tenants/{tenantId}/clubs', { params: { path: { tenantId } }, body }));
			}
			await queryClient.invalidateQueries({ queryKey: keys.clubs(tenantId) });
			toast.success(editId ? $t`Club saved` : $t`Club "${name}" created`);
			editId = null;
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Deletes a club after confirmation; its shooters stay. */
	async function remove(c: Club) {
		if (!(await ask({ title: $t`Delete club "${c.name}"?`, description: $plural(c.shooterCount, { one: '# shooter assignment is removed, the shooter stays.', other: '# shooter assignments are removed, the shooters stay.' }), action: $t`Delete`, destructive: true }))) return;
		try {
			await call(api.DELETE('/api/tenants/{tenantId}/clubs/{clubId}', { params: { path: { tenantId, clubId: c.id } } }));
			await queryClient.invalidateQueries({ queryKey: keys.clubs(tenantId) });
			toast.success($t`Club deleted`);
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
</script>

<PageHeader help="user-guide/shooters-clubs" title={$t`Clubs`} description={$t`Clubs that shooters belong to; used for registrations and the club competition.`}>
	{#snippet actions()}
		{#if isAdmin}<Button onclick={() => edit(null)}><Plus />{$t`New club`}</Button>{/if}
	{/snippet}
</PageHeader>

<div class="overflow-hidden rounded-lg border bg-card">
	<div class="flex items-center gap-2 border-b p-3">
		<div class="relative w-full sm:max-w-sm">
			<Search class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
			<Input class="pl-8" type="search" placeholder={$t`Filter`} bind:value={filter} oninput={() => (pageNo = 0)} />
		</div>
	</div>
	<Table.Root>
		<Table.Header>
			<Table.Row><Table.Head>{$t`Name`}</Table.Head><Table.Head>{$t`SSV number`}</Table.Head><Table.Head class="text-right">{$t`Shooters`}</Table.Head><Table.Head class="w-10"></Table.Head></Table.Row>
		</Table.Header>
		<Table.Body>
			{#each pageOf(shown, pageNo, pageSize) as c (c.id)}
				<Table.Row>
					<Table.Cell class="font-medium">{c.name}</Table.Cell>
					<Table.Cell class="font-mono text-xs">{c.ssvNumber || '–'}</Table.Cell>
					<Table.Cell class="text-right tabular-nums">{c.shooterCount}</Table.Cell>
					<Table.Cell class="text-right">
						{#if isAdmin}
							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									{#snippet child({ props })}<Button variant="ghost" size="icon-sm" aria-label={$t`Actions for ${c.name}`} {...props}><MoreHorizontal /></Button>{/snippet}
								</DropdownMenu.Trigger>
								<DropdownMenu.Content align="end">
									<DropdownMenu.Item onSelect={() => edit(c)}><Pencil />{$t`Edit`}</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item variant="destructive" onSelect={() => remove(c)}><Trash />{$t`Delete`}</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						{/if}
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row><Table.Cell colspan={4} class="py-10 text-center text-muted-foreground">{$t`No clubs.`}</Table.Cell></Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
	<Pager bind:page={pageNo} total={shown.length} size={pageSize} />
</div>

<Dialog.Root open={editId !== null} onOpenChange={(o) => !o && (editId = null)}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header><Dialog.Title>{editId ? $t`Edit club` : $t`New club`}</Dialog.Title></Dialog.Header>
		<form class="grid gap-4" onsubmit={save}>
			<FormField label={$t`Name`}><Input bind:value={name} required /></FormField>
			<FormField label={$t`SSV number`}><Input bind:value={ssvNumber} /></FormField>
			<FormError message={error} />
			<Dialog.Footer><Button type="submit">{$t`Save`}</Button></Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
