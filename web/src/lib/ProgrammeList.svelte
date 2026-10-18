<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Pencil from '@lucide/svelte/icons/pencil';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import Copy from '@lucide/svelte/icons/copy';
	import Download from '@lucide/svelte/icons/download';
	import Library from '@lucide/svelte/icons/library';
	import MoreHorizontal from '@lucide/svelte/icons/ellipsis';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import Upload from '@lucide/svelte/icons/upload';
	import { toast } from 'svelte-sonner';
	import { programmeTitle } from './format';
	import { api, call, type Programme, type ProgrammeExport } from './api';
	import { ageCategoryLabel } from './labels';
	import { Button } from './components/ui/button';
	import { Checkbox } from './components/ui/checkbox';
	import * as Dialog from './components/ui/dialog';
	import * as DropdownMenu from './components/ui/dropdown-menu';
	import { Input } from './components/ui/input';
	import * as NativeSelect from './components/ui/native-select';
	import * as Table from './components/ui/table';
	import { ask } from './confirm.svelte';
	import FormError from './FormError.svelte';
	import { eventsQuery, keys } from './tenant';

	let { tenantId, eventId, isAdmin }: { tenantId: number; eventId: number; isAdmin: boolean } = $props();

	const queryClient = useQueryClient();
	const programmes = createQuery(() => ({
		queryKey: keys.programmes(tenantId, eventId),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/programmes', { params: { path: { tenantId, eventId } } }))
	}));
	const templates = createQuery(() => ({
		queryKey: ['programme-templates'],
		queryFn: () => call(api.GET('/api/programme-templates')),
		enabled: isAdmin
	}));
	const events = eventsQuery(() => tenantId);
	const otherEvents = $derived((events.data ?? []).filter((e) => e.id !== eventId));

	let error = $state('');
	let adding = $state<'library' | 'copy' | 'import' | null>(null);
	const base = $derived(`/t/${tenantId}/events/${eventId}/config/programmes`);
	let templateId = $state(0);
	let fromEventId = $state(0);
	let file = $state<FileList>();
	const evPath = $derived({ params: { path: { tenantId, eventId } } });

	/** Runs a programme action, refreshes the list and shows its message or error. */
	async function run(action: () => Promise<string>) {
		error = '';
		try {
			const message = await action();
			await queryClient.invalidateQueries({ queryKey: keys.programmes(tenantId, eventId) });
			if (message) toast.success(message);
			adding = null;
		} catch (err) {
			error = (err as Error).message;
			if (!adding) toast.error(error);
		}
	}

	/** Adds the selected library template to the event. */
	const addTemplate = () =>
		run(async () => {
			const p = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/programmes/from-template', { ...evPath, body: { templateId } }));
			return $t`Added ${programmeTitle(p)} from the library.`;
		});

	/** Copies all programmes of the selected earlier event. */
	const copyFromEvent = () =>
		run(async () => {
			const r = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/programmes/copy', { ...evPath, body: { fromEventId } }));
			return $plural(r.copied, { one: 'Copied # programme.', other: 'Copied # programmes.' });
		});

	/** Imports a programme from the chosen JSON export file. */
	const importFile = () =>
		run(async () => {
			if (!file?.[0]) return '';
			let doc: ProgrammeExport;
			try {
				doc = JSON.parse(await file[0].text());
			} catch {
				throw new Error($t`The file is not valid JSON`);
			}
			const p = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/programmes/import', { ...evPath, body: doc }));
			return $t`Imported ${programmeTitle(p)}.`;
		});

	/** Turns a programme on or off. */
	const setActive = (p: Programme, active: boolean) =>
		run(async () => {
			await call(api.PUT('/api/tenants/{tenantId}/programmes/{programmeId}/active', { params: { path: { tenantId, programmeId: p.id } }, body: { active } }));
			return '';
		});

	/** Asks for confirmation, then deletes a programme. */
	const remove = (p: Programme) =>
		run(async () => {
			if (!(await ask({ title: $t`Delete programme?`, description: $t`${programmeTitle(p)} is deleted with its series, prices and awards.`, action: $t`Delete`, destructive: true }))) return '';
			await call(api.DELETE('/api/tenants/{tenantId}/programmes/{programmeId}', { params: { path: { tenantId, programmeId: p.id } } }));
			return $t`Deleted ${programmeTitle(p)}.`;
		});


	/** Price or price range in CHF with the age categories it is offered to. */
	function priceRange(p: Programme) {
		if (p.prices.length === 0) return $t`nobody`;
		const values = p.prices.map((x) => x.price);
		const min = Math.min(...values).toFixed(2);
		const max = Math.max(...values).toFixed(2);
		const cats = p.prices.map((x) => $t(ageCategoryLabel(x.ageCategory))).join(', ');
		return `CHF ${min === max ? min : `${min}–${max}`} (${cats})`;
	}
</script>

<div class="flex flex-wrap items-center gap-2 border-b px-4 py-3">
	{#if isAdmin}
		<div class="ml-auto flex gap-2">
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}<Button size="sm" variant="outline" {...props}>{$t`Add from …`}<ChevronDown /></Button>{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					<DropdownMenu.Item onSelect={() => (adding = 'library')}><Library />{$t`Programme library`}</DropdownMenu.Item>
					<DropdownMenu.Item onSelect={() => (adding = 'copy')}><Copy />{$t`Earlier event`}</DropdownMenu.Item>
					<DropdownMenu.Item onSelect={() => (adding = 'import')}><Upload />{$t`Programme file`}</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
			<Button size="sm" href="{base}/new"><Plus />{$t`New programme`}</Button>
		</div>
	{/if}
</div>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>{$t`Programme`}</Table.Head>
			<Table.Head class="hidden sm:table-cell">{$t`Shots`}</Table.Head>
			<Table.Head class="hidden sm:table-cell">{$t`Prices`}</Table.Head>
			<Table.Head class="hidden sm:table-cell">{$t`Awards`}</Table.Head>
			<Table.Head>{$t`Active`}</Table.Head>
			<Table.Head class="w-10"></Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each programmes.data ?? [] as p (p.id)}
			<Table.Row class={p.active ? '' : 'opacity-50'}>
				<Table.Cell class="whitespace-normal">
					<a class="font-medium hover:underline" href="{base}/{p.id}">{programmeTitle(p)}</a>
					<div class="text-xs text-muted-foreground sm:hidden">{$t`${p.shots} / ${p.maxPoints} pts`} · {priceRange(p)}</div>
				</Table.Cell>
				<Table.Cell class="hidden sm:table-cell tabular-nums">{$t`${p.shots} / ${p.maxPoints} pts`}</Table.Cell>
				<Table.Cell class="hidden sm:table-cell max-w-64 truncate" title={priceRange(p)}>{priceRange(p)}</Table.Cell>
				<Table.Cell class="hidden sm:table-cell whitespace-normal">
					{p.awards.map((a) => a.name).join(', ') || '–'}
					{#if p.specialPrizes.ranks}<div class="text-xs text-muted-foreground">{$t`Special prizes: ranks 1–${p.specialPrizes.ranks}`}</div>{/if}
				</Table.Cell>
				<Table.Cell>
					<Checkbox checked={p.active} disabled={!isAdmin} onCheckedChange={(v) => setActive(p, v)} aria-label={$t`Active`} />
				</Table.Cell>
				<Table.Cell class="text-right">
					<DropdownMenu.Root>
						<DropdownMenu.Trigger>
							{#snippet child({ props })}<Button variant="ghost" size="icon-sm" aria-label={$t`Actions for ${programmeTitle(p)}`} {...props}><MoreHorizontal /></Button>{/snippet}
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="end">
							<DropdownMenu.Item>
								{#snippet child({ props })}<a href="{base}/{p.id}" {...props}><Pencil />{isAdmin ? $t`Edit` : $t`Open`}</a>{/snippet}
							</DropdownMenu.Item>
							<DropdownMenu.Item>
								{#snippet child({ props })}<a href="/api/tenants/{tenantId}/programmes/{p.id}/export" download {...props}><Download />{$t`Export`}</a>{/snippet}
							</DropdownMenu.Item>
							{#if isAdmin}
								<DropdownMenu.Separator />
								<DropdownMenu.Item variant="destructive" onSelect={() => remove(p)}><Trash />{$t`Delete`}</DropdownMenu.Item>
							{/if}
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</Table.Cell>
			</Table.Row>
		{:else}
			<Table.Row><Table.Cell colspan={6} class="py-8 text-center text-muted-foreground">{$t`No programmes yet.`}</Table.Cell></Table.Row>
		{/each}
	</Table.Body>
</Table.Root>

<Dialog.Root open={adding !== null} onOpenChange={(o) => !o && ((adding = null), (error = ''))}>
	<Dialog.Content class="sm:max-w-md">
		{#if adding === 'library'}
			<Dialog.Header><Dialog.Title>{$t`Add from the library`}</Dialog.Title><Dialog.Description>{$t`Copies a template of the programme library into this event.`}</Dialog.Description></Dialog.Header>
			<NativeSelect.Root class="w-full" bind:value={templateId} aria-label={$t`Library template`}>
				<NativeSelect.Option value={0}>{$t`Choose a template …`}</NativeSelect.Option>
				{#each templates.data ?? [] as tpl (tpl.id)}<NativeSelect.Option value={tpl.id}>{programmeTitle(tpl.definition)}</NativeSelect.Option>{/each}
			</NativeSelect.Root>
			<FormError message={error} />
			<Dialog.Footer><Button disabled={!templateId} onclick={addTemplate}>{$t`Add from library`}</Button></Dialog.Footer>
		{:else if adding === 'copy'}
			<Dialog.Header><Dialog.Title>{$t`Copy from an earlier event`}</Dialog.Title><Dialog.Description>{$t`Copies all programmes of the chosen event.`}</Dialog.Description></Dialog.Header>
			<NativeSelect.Root class="w-full" bind:value={fromEventId} aria-label={$t`Copy from event`}>
				<NativeSelect.Option value={0}>{$t`Choose an event …`}</NativeSelect.Option>
				{#each otherEvents as e (e.id)}<NativeSelect.Option value={e.id}>{e.name} ({e.startsOn})</NativeSelect.Option>{/each}
			</NativeSelect.Root>
			<FormError message={error} />
			<Dialog.Footer><Button disabled={!fromEventId} onclick={copyFromEvent}>{$t`Copy all programmes`}</Button></Dialog.Footer>
		{:else if adding === 'import'}
			<Dialog.Header><Dialog.Title>{$t`Import a programme file`}</Dialog.Title><Dialog.Description>{$t`A JSON file exported from Standblatt (format standblatt.programme v1).`}</Dialog.Description></Dialog.Header>
			<Input type="file" accept=".json,application/json" bind:files={file} aria-label={$t`Programme file`} />
			<FormError message={error} />
			<Dialog.Footer><Button disabled={!file?.length} onclick={importFile}>{$t`Import programme file`}</Button></Dialog.Footer>
		{/if}
	</Dialog.Content>
</Dialog.Root>
