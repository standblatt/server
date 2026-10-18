<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Download from '@lucide/svelte/icons/download';
	import FileUp from '@lucide/svelte/icons/file-up';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call, type Layout, type LayoutInput } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import * as Table from '$lib/components/ui/table';
	import { ask } from '$lib/confirm.svelte';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import ImageLibrary from '$lib/forms/ImageLibrary.svelte';
	import type { Source } from '$lib/forms/editor.svelte';
	import { sourceKeys, sourceLabels } from '$lib/forms/sources';
	import Copy from '@lucide/svelte/icons/copy';
	import { exportImages, imagesKey, imagesQuery, importImages, type ExportedImage } from '$lib/forms/images';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const key = $derived(['tenants', tenantId, 'layouts']);
	const layouts = createQuery(() => ({
		queryKey: key,
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/layouts', { params: { path: { tenantId } } }))
	}));

	const images = imagesQuery(() => tenantId);
	const sizes: Record<string, [number, number]> = { a5: [148, 210], a4: [210, 297], a6: [105, 148], a5l: [210, 148] };
	let creating = $state(false);
	let name = $state('');
	let size = $state('a5');
	let source = $state<Source>('results');
	const builtins = createQuery(() => ({
		queryKey: ['tenants', tenantId, 'layouts', 'builtin'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/layouts/builtin', { params: { path: { tenantId } } }))
	}));
	const builtinNames = $derived<Record<string, string>>({ wreath_card: $t`Wreath card`, control_strip: $t`Control strip` });

	/** Creates an editable copy of a built-in form. */
	async function copyBuiltin(key: string, layout: LayoutInput) {
		let name = builtinNames[key];
		for (let n = 2; (layouts.data ?? []).some((l) => l.name === name); n++) name = `${builtinNames[key]} ${n}`;
		await create({ ...layout, name });
	}
	let error = $state('');
	let importInput = $state<HTMLInputElement>();

	/** Creates a form and opens it in the editor. */
	async function create(body: LayoutInput) {
		const l = await call(api.POST('/api/tenants/{tenantId}/layouts', { params: { path: { tenantId } }, body }));
		await queryClient.invalidateQueries({ queryKey: key });
		await goto(`/t/${tenantId}/forms/${l.id}`);
	}

	/** Creates an empty form of the chosen size from the dialog. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		try {
			const [width, height] = sizes[size];
			await create({ name: name.trim(), source, width, height, offsetX: 0, offsetY: 0, elements: [] });
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Creates a form from an exported JSON file. */
	async function importFile(e: Event) {
		const f = (e.currentTarget as HTMLInputElement).files?.[0];
		if (!f) return;
		try {
			const l = JSON.parse(await f.text()) as LayoutInput & { images?: ExportedImage[] };
			const body: LayoutInput = { name: l.name, source: l.source, width: l.width, height: l.height, offsetX: l.offsetX ?? 0, offsetY: l.offsetY ?? 0, elements: l.elements ?? [] };
			if (l.images?.length) {
				await importImages(tenantId, body, l.images, new Map((images.data ?? []).map((i) => [i.name, i.id])));
				await queryClient.invalidateQueries({ queryKey: imagesKey(tenantId) });
			}
			await create(body);
		} catch (err) {
			toast.error(err instanceof SyntaxError ? $t`The file is not a print form.` : (err as Error).message);
		} finally {
			if (importInput) importInput.value = '';
		}
	}

	/** Downloads a form with the images it shows as a JSON file to share it with another club. */
	async function exportForm(l: Layout) {
		const form: LayoutInput = { name: l.name, source: l.source, width: l.width, height: l.height, offsetX: l.offsetX, offsetY: l.offsetY, elements: l.elements };
		const body = { ...form, images: await exportImages(tenantId, form, new Map((images.data ?? []).map((i) => [i.id, i.name]))) };
		const a = document.createElement('a');
		a.href = URL.createObjectURL(new Blob([JSON.stringify(body, null, 2)], { type: 'application/json' }));
		a.download = `${l.name.replace(/[^A-Za-z0-9]+/g, '_')}.json`;
		a.click();
		setTimeout(() => URL.revokeObjectURL(a.href), 60_000);
	}

	/** Deletes a form after confirmation. */
	async function remove(l: Layout) {
		if (!(await ask({ title: $t`Delete print form "${l.name}"?`, description: $t`Programmes that use it no longer print a form.`, action: $t`Delete`, destructive: true }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/layouts/{layoutId}', { params: { path: { tenantId, layoutId: l.id } } }));
		await queryClient.invalidateQueries({ queryKey: key });
		toast.success($t`Print form deleted`);
	}
</script>

<PageHeader help="user-guide/print-forms" title={$t`Print forms`} description={$t`Templates that print results onto pre-printed sheets, e.g. the army's lane sheets.`}>
	{#snippet actions()}
		<Button variant="outline" onclick={() => importInput?.click()}><FileUp />{$t`Import`}</Button>
		<input bind:this={importInput} class="hidden" type="file" accept=".json,application/json" onchange={importFile} aria-label={$t`Import a print form`} />
		<Button onclick={() => ((creating = true), (name = ''), (error = ''))}><Plus />{$t`New print form`}</Button>
	{/snippet}
</PageHeader>

<div class="overflow-hidden rounded-lg border bg-card">
	<Table.Root>
		<Table.Header>
			<Table.Row><Table.Head>{$t`Name`}</Table.Head><Table.Head class="hidden sm:table-cell">{$t`Prints`}</Table.Head><Table.Head class="hidden sm:table-cell">{$t`Size`}</Table.Head><Table.Head class="hidden sm:table-cell text-right">{$t`Elements`}</Table.Head><Table.Head></Table.Head></Table.Row>
		</Table.Header>
		<Table.Body>
			{#each layouts.data ?? [] as l (l.id)}
				<Table.Row>
					<Table.Cell class="whitespace-normal">
						<a class="font-medium hover:underline" href="/t/{tenantId}/forms/{l.id}">{l.name}</a>
						<div class="text-xs text-muted-foreground sm:hidden">{l.width} × {l.height} mm</div>
					</Table.Cell>
					<Table.Cell class="hidden sm:table-cell">{$t(sourceLabels[(l.source || 'results') as Source])}</Table.Cell>
					<Table.Cell class="hidden sm:table-cell">{l.width} × {l.height} mm</Table.Cell>
					<Table.Cell class="hidden sm:table-cell text-right tabular-nums">{l.elements.length}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-1">
							<Button variant="ghost" size="icon-sm" title={$t`Export`} aria-label={$t`Export ${l.name}`} onclick={() => exportForm(l)}><Download /></Button>
							<Button variant="ghost" size="icon-sm" title={$t`Delete`} aria-label={$t`Delete ${l.name}`} onclick={() => remove(l)}><Trash /></Button>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row><Table.Cell colspan={5} class="py-10 text-center text-muted-foreground">{$t`No print forms yet.`}</Table.Cell></Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>

<section class="mt-4 overflow-hidden rounded-lg border bg-card" aria-label={$t`Built-in forms`}>
	<div class="border-b px-4 py-2.5">
		<h2 class="text-sm font-semibold">{$t`Built-in forms`}</h2>
		<p class="text-xs text-muted-foreground">{$t`Standblatt prints these unless an event chooses a form of the club; copy one to change it.`}</p>
	</div>
	<ul class="divide-y">
		{#each builtins.data ?? [] as b (b.key)}
			<li class="flex items-center gap-3 px-4 py-2 text-sm">
				<span class="mr-auto font-medium">{builtinNames[b.key]}</span>
				<span class="hidden text-muted-foreground sm:inline">{$t(sourceLabels[(b.layout.source || 'results') as Source])} · {b.layout.width} × {b.layout.height} mm</span>
				<Button size="sm" variant="outline" aria-label={$t`Copy ${builtinNames[b.key]}`} onclick={() => copyBuiltin(b.key, b.layout)}><Copy />{$t`Copy`}</Button>
			</li>
		{/each}
	</ul>
</section>

<div class="mt-4"><ImageLibrary {tenantId} /></div>

<Dialog.Root bind:open={creating}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header><Dialog.Title>{$t`New print form`}</Dialog.Title></Dialog.Header>
		<form class="grid gap-4" onsubmit={submit}>
			<FormField label={$t`Name`}><Input bind:value={name} required maxlength={100} /></FormField>
			<FormField label={$t`Prints`} hint={$t`Decides the fields and where the form is printed.`}>
				<NativeSelect.Root class="w-full" bind:value={source}>
					{#each sourceKeys as value (value)}<NativeSelect.Option {value}>{$t(sourceLabels[value])}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
			</FormField>
			<FormField label={$t`Paper size`} hint={$t`Can be changed to any size in millimetres later.`}>
				<NativeSelect.Root class="w-full" bind:value={size}>
					<NativeSelect.Option value="a5">{$t`A5 portrait (148 × 210 mm)`}</NativeSelect.Option>
					<NativeSelect.Option value="a5l">{$t`A5 landscape (210 × 148 mm)`}</NativeSelect.Option>
					<NativeSelect.Option value="a4">{$t`A4 portrait (210 × 297 mm)`}</NativeSelect.Option>
					<NativeSelect.Option value="a6">{$t`A6 portrait (105 × 148 mm)`}</NativeSelect.Option>
				</NativeSelect.Root>
			</FormField>
			<FormError message={error} />
			<Dialog.Footer><Button type="submit">{$t`Create`}</Button></Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
