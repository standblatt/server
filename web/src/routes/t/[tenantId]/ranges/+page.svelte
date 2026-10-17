<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { api, call, type Range } from '$lib/api';
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
	import Printer from '@lucide/svelte/icons/printer';
	import { keys, rangesQuery } from '$lib/tenant';
	import { formatLanes, parseLanes } from '$lib/format';
	import X from '@lucide/svelte/icons/x';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const ranges = rangesQuery(() => tenantId);

	let editId = $state<number | null>(null);
	let name = $state('');
	let laneCount = $state(10);
	let disabled = $state('');
	let groups = $state<{ id?: number; name: string; lanes: string }[]>([]);
	let error = $state('');

	/** Fills the form with a range, or empties it for a new one. */
	function edit(s: Range | null) {
		editId = s?.id ?? 0;
		name = s?.name ?? '';
		laneCount = s?.laneCount ?? 10;
		disabled = formatLanes(s?.disabledLanes ?? []);
		groups = (s?.groups ?? []).map((g) => ({ id: g.id, name: g.name, lanes: formatLanes(g.lanes) }));
		error = '';
	}

	/** Creates or updates the range from the form. */
	async function save(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		const body = { name, laneCount, disabledLanes: parseLanes(disabled), groups: groups.map((g) => ({ id: g.id, name: g.name.trim(), lanes: parseLanes(g.lanes) })) };
		if ([body.disabledLanes, ...body.groups.map((g) => g.lanes)].some((l) => l.some(Number.isNaN))) {
			error = $t`Write lanes like 1-10, 12`;
			return;
		}
		try {
			if (editId) {
				await call(api.PUT('/api/tenants/{tenantId}/ranges/{rangeId}', { params: { path: { tenantId, rangeId: editId } }, body }));
			} else {
				await call(api.POST('/api/tenants/{tenantId}/ranges', { params: { path: { tenantId } }, body }));
			}
			await queryClient.invalidateQueries({ queryKey: keys.ranges(tenantId) });
			toast.success(editId ? $t`Range saved` : $t`Range "${name}" created`);
			editId = null;
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Deletes a range after confirmation. */
	async function remove(s: Range) {
		if (!(await ask({ title: $t`Delete range "${s.name}"?`, description: $t`Ranges with events cannot be deleted.`, action: $t`Delete`, destructive: true }))) return;
		try {
			await call(api.DELETE('/api/tenants/{tenantId}/ranges/{rangeId}', { params: { path: { tenantId, rangeId: s.id } } }));
			await queryClient.invalidateQueries({ queryKey: keys.ranges(tenantId) });
			toast.success($t`Range deleted`);
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
</script>

<PageHeader help="user-guide/ranges" title={$t`Ranges`} description={$t`Shooting ranges of this club with their lanes and printers.`}>
	{#snippet actions()}<Button onclick={() => edit(null)}><Plus />{$t`New range`}</Button>{/snippet}
</PageHeader>

<div class="overflow-hidden rounded-lg border bg-card">
	<Table.Root>
		<Table.Header>
			<Table.Row><Table.Head>{$t`Name`}</Table.Head><Table.Head class="text-right">{$t`Lanes`}</Table.Head><Table.Head>{$t`Lane groups`}</Table.Head><Table.Head>{$t`Out of service`}</Table.Head><Table.Head></Table.Head></Table.Row>
		</Table.Header>
		<Table.Body>
			{#each ranges.data ?? [] as s (s.id)}
				<Table.Row>
					<Table.Cell class="font-medium">{s.name}</Table.Cell>
					<Table.Cell class="text-right tabular-nums">{s.laneCount}</Table.Cell>
					<Table.Cell>{s.groups.map((g) => `${g.name}: ${formatLanes(g.lanes)}`).join(' · ') || '–'}</Table.Cell>
					<Table.Cell>{formatLanes(s.disabledLanes) || '–'}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-1">
							<Button size="sm" variant="outline" href="/t/{tenantId}/ranges/{s.id}/printing"><Printer />{$t`Printing`}</Button>
							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									{#snippet child({ props })}<Button variant="ghost" size="icon-sm" aria-label={$t`Actions for ${s.name}`} {...props}><MoreHorizontal /></Button>{/snippet}
								</DropdownMenu.Trigger>
								<DropdownMenu.Content align="end">
									<DropdownMenu.Item onSelect={() => edit(s)}><Pencil />{$t`Edit`}</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item variant="destructive" onSelect={() => remove(s)}><Trash />{$t`Delete`}</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row><Table.Cell colspan={5} class="py-10 text-center text-muted-foreground">{$t`No ranges yet.`}</Table.Cell></Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>

<Dialog.Root open={editId !== null} onOpenChange={(o) => !o && (editId = null)}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header><Dialog.Title>{editId ? $t`Edit range` : $t`New range`}</Dialog.Title></Dialog.Header>
		<form class="grid gap-4" onsubmit={save}>
			<FormField label={$t`Name`}><Input bind:value={name} required /></FormField>
			<FormField label={$t`Number of lanes`}><Input class="w-28" type="number" min="1" max="99" bind:value={laneCount} required /></FormField>
			<FormField label={$t`Out of service`} hint={$t`Lane numbers, e.g. 3, 7`}><Input bind:value={disabled} /></FormField>
			<FormField label={$t`Lane groups`} hint={$t`Optional: lanes with another distance on the same range, e.g. 25 m on 1-10 and 50 m on 11-15. Programmes can then be tied to a group.`}>
				<div class="grid gap-2">
					{#each groups as g, i (i)}
						<div class="flex items-center gap-2">
							<Input class="w-28" bind:value={g.name} required placeholder="25 m" aria-label={$t`Group name`} />
							<Input bind:value={g.lanes} required placeholder="1-10" aria-label={$t`Lanes of ${g.name}`} />
							<Button variant="ghost" size="icon-sm" aria-label={$t`Remove group`} onclick={() => groups.splice(i, 1)}><X /></Button>
						</div>
					{/each}
					<div><Button variant="outline" size="sm" onclick={() => groups.push({ name: '', lanes: '' })}><Plus />{$t`Add group`}</Button></div>
				</div>
			</FormField>
			<FormError message={error} />
			<Dialog.Footer><Button type="submit">{$t`Save`}</Button></Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
