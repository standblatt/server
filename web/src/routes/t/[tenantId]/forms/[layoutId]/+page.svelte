<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { beforeNavigate, goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import FileText from '@lucide/svelte/icons/file-text';
	import Minus from '@lucide/svelte/icons/minus';
	import Plus from '@lucide/svelte/icons/plus';
	import Printer from '@lucide/svelte/icons/printer';
	import Redo from '@lucide/svelte/icons/redo-2';
	import Save from '@lucide/svelte/icons/save';
	import Trash from '@lucide/svelte/icons/trash-2';
	import Undo from '@lucide/svelte/icons/undo-2';
	import { toast } from 'svelte-sonner';
	import { api, call, type LayoutElement } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Switch } from '$lib/components/ui/switch';
	import { ask } from '$lib/confirm.svelte';
	import { FormEditor, num } from '$lib/forms/editor.svelte';
	import Palette from '$lib/forms/Palette.svelte';
	import Properties from '$lib/forms/Properties.svelte';
	import Stage from '$lib/forms/Stage.svelte';
	import { imagesQuery, imageUrl } from '$lib/forms/images';
	import { sourceKeys, sourceLabels } from '$lib/forms/sources';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { pdfClick } from '$lib/print';

	const tenantId = $derived(Number(page.params.tenantId));
	const layoutId = $derived(Number(page.params.layoutId));
	const queryClient = useQueryClient();
	const path = $derived({ params: { path: { tenantId, layoutId } } });
	const layout = createQuery(() => ({
		queryKey: ['tenants', tenantId, 'layouts', layoutId],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/layouts/{layoutId}', path))
	}));
	const editor = new FormEditor();
	const fields = createQuery(() => ({
		queryKey: ['tenants', tenantId, 'layouts', 'fields', editor.source],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/layouts/fields', { params: { path: { tenantId }, query: { source: editor.source } } }))
	}));
	const images = imagesQuery(() => tenantId);
	const samples = $derived(Object.fromEntries((fields.data ?? []).map((f) => [f.key, f.example])));

	let loaded = $state(0);
	let version = $state(0);
	let names = $state(false);
	let error = $state('');
	let busy = $state(false);
	let leaving = false;
	let zoom = $state(1);
	let autosave = $state(true);
	try {
		zoom = Number(localStorage.getItem('standblatt-form-zoom')) || 1;
		autosave = localStorage.getItem('standblatt-form-autosave') !== 'off';
	} catch {
		/* no storage: keep 100 % and autosave */
	}

	$effect(() => {
		const l = layout.data;
		if (!l || loaded === l.id) return;
		loaded = l.id;
		editor.load(l);
	});

	const base = $derived(`/api/tenants/${tenantId}/layouts/${layoutId}`);

	/** Sets the zoom between 50 and 200 % and remembers it. */
	function setZoom(z: number) {
		zoom = Math.min(2, Math.max(0.5, Math.round(z * 10) / 10));
		try {
			localStorage.setItem('standblatt-form-zoom', String(zoom));
		} catch {
			/* no storage */
		}
	}

	/** Switches saving every minute on or off and remembers it. */
	function setAutosave(on: boolean) {
		autosave = on;
		try {
			localStorage.setItem('standblatt-form-autosave', on ? 'on' : 'off');
		} catch {
			/* no storage */
		}
	}

	$effect(() => {
		if (!autosave) return;
		const timer = setInterval(() => {
			if (editor.dirty && !busy) save(true);
		}, 60_000);
		return () => clearInterval(timer);
	});

	/** Adds an element from the palette below the lowest one. */
	function add(el: Partial<LayoutElement>) {
		const lowest = Math.max(0, ...editor.items.map((i) => num(i.y)));
		editor.add({ ...el, x: 10, y: editor.items.length ? Math.min(lowest + 8, num(editor.height) - 10) : 10 });
	}

	/** Stores the form; quietly for the autosave. Changes made while it is stored stay unsaved. */
	async function save(quiet = false) {
		error = '';
		busy = true;
		const state = editor.begin();
		try {
			await call(api.PUT('/api/tenants/{tenantId}/layouts/{layoutId}', { ...path, body: editor.input() }));
			editor.markClean(state);
			version++;
			await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId, 'layouts'] });
			if (!quiet) toast.success($t`Print form saved`);
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}

	/** Uploads a scan of the empty form for the sheet. */
	async function upload(ev: Event) {
		const input = ev.currentTarget as HTMLInputElement;
		const f = input.files?.[0];
		if (!f) return;
		try {
			const form = new FormData();
			form.append('file', f);
			await call(api.PUT('/api/tenants/{tenantId}/layouts/{layoutId}/background', { ...path, body: {} as never, bodySerializer: () => form }));
			version++;
			await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId, 'layouts', layoutId] });
		} catch (err) {
			toast.error((err as Error).message);
		} finally {
			input.value = '';
		}
	}

	/** Removes the scan of the empty form. */
	async function removeBackground() {
		await call(api.DELETE('/api/tenants/{tenantId}/layouts/{layoutId}/background', path));
		version++;
		await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId, 'layouts', layoutId] });
	}

	/** Ctrl+S saves. */
	function key(e: KeyboardEvent) {
		if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
			e.preventDefault();
			save();
		}
	}

	beforeNavigate(({ cancel, to, type }) => {
		if (leaving || !editor.dirty) return;
		cancel();
		if (type === 'leave' || !to) return;
		ask({ title: $t`Leave without saving?`, description: $t`The print form has changes that are not saved yet.`, action: $t`Leave` }).then((ok) => {
			if (!ok) return;
			leaving = true;
			goto(to.url);
		});
	});
</script>

<svelte:window onkeydown={key} />

<PageHeader help="user-guide/print-forms" title={layout.data?.name ?? $t`Print form`} description={$t`Drag fields onto the sheet; positions and sizes are in millimetres from the top left corner.`}>
	{#snippet actions()}
		<Button variant="outline" size="icon" title={$t`Undo (Ctrl+Z)`} aria-label={$t`Undo`} disabled={!editor.canUndo} onclick={() => editor.undo()}><Undo /></Button>
		<Button variant="outline" size="icon" title={$t`Redo (Ctrl+Y)`} aria-label={$t`Redo`} disabled={!editor.canRedo} onclick={() => editor.redo()}><Redo /></Button>
		<Button variant="outline" href="{base}/preview?v={version}" target="_blank" onclick={pdfClick}><FileText />{$t`PDF preview`}</Button>
		<Button variant="outline" href="{base}/preview?print=true&v={version}" target="_blank" onclick={pdfClick}><Printer />{$t`Test print`}</Button>
		<Button onclick={() => save()} disabled={busy}><Save />{$t`Save`}{editor.dirty ? ' •' : ''}</Button>
	{/snippet}
</PageHeader>

{#if layout.data}
	<FormError message={error} />
	<div class="grid gap-4 *:min-w-0 lg:grid-cols-[15rem_minmax(0,1fr)_17rem]">
		<aside class="rounded-lg border bg-card p-3" aria-label={$t`Fields`}>
			<h2 class="mb-2 text-sm font-semibold">{$t`Fields`}</h2>
			<Palette {samples} images={images.data ?? []} library="/t/{tenantId}/forms" onadd={add} />
		</aside>

		<section class="grid content-start gap-2" aria-label={$t`Sheet`}>
			<div class="flex flex-wrap items-center gap-3 text-sm">
				<div class="flex items-center gap-1">
					<Button variant="outline" size="icon-sm" aria-label={$t`Zoom out`} onclick={() => setZoom(zoom - 0.1)}><Minus /></Button>
					<button type="button" class="w-14 text-center tabular-nums" title={$t`Back to 100 %`} onclick={() => setZoom(1)}>{Math.round(zoom * 100)} %</button>
					<Button variant="outline" size="icon-sm" aria-label={$t`Zoom in`} onclick={() => setZoom(zoom + 0.1)}><Plus /></Button>
				</div>
				<label class="flex items-center gap-2"><Switch bind:checked={names} />{$t`Show field names`}</label>
				<label class="flex items-center gap-2" title={$t`Saves unsaved changes every minute`}><Checkbox checked={autosave} onCheckedChange={setAutosave} />{$t`Save automatically`}</label>
				<span class="ml-auto text-xs text-muted-foreground">{$t`Alt: no snapping · Shift: add to the selection`}</span>
			</div>
			<Stage {editor} {samples} {names} {zoom} onzoom={setZoom} imageUrl={(id) => imageUrl(tenantId, id)} background={layout.data.hasBackground ? `${base}/background?v=${version}` : undefined} />
		</section>

		<div class="grid content-start gap-4">
			<aside class="rounded-lg border bg-card p-3" aria-label={$t`Properties`}>
				<h2 class="mb-2 text-sm font-semibold">{$t`Properties`}</h2>
				<Properties {editor} images={images.data ?? []} />
			</aside>
			<aside class="grid gap-3 rounded-lg border bg-card p-3" aria-label={$t`Form`}>
				<h2 class="text-sm font-semibold">{$t`Form`}</h2>
				<FormField label={$t`Name`}><Input bind:value={editor.name} required maxlength={100} /></FormField>
				<FormField label={$t`Prints`} hint={$t`Decides the fields and where the form is printed.`}>
					<NativeSelect.Root class="w-full" bind:value={editor.source}>
						{#each sourceKeys as value (value)}<NativeSelect.Option {value}>{$t(sourceLabels[value])}</NativeSelect.Option>{/each}
					</NativeSelect.Root>
				</FormField>
				<div class="grid grid-cols-2 gap-2">
					<FormField label={$t`Width (mm)`}><Input type="number" min="20" max="500" step="0.5" bind:value={editor.width} /></FormField>
					<FormField label={$t`Height (mm)`}><Input type="number" min="20" max="500" step="0.5" bind:value={editor.height} /></FormField>
					<FormField label={$t`Printer offset right (mm)`}><Input type="number" min="-50" max="50" step="0.1" bind:value={editor.offsetX} /></FormField>
					<FormField label={$t`Printer offset down (mm)`}><Input type="number" min="-50" max="50" step="0.1" bind:value={editor.offsetY} /></FormField>
				</div>
				<p class="text-xs text-muted-foreground">{$t`Moves everything when the printer prints off target; check it with a test print.`}</p>
				<FormField label={$t`Scan of the empty form`} hint={$t`PNG or JPEG of the whole sheet; only shown while editing, never printed.`}>
					<div class="flex items-center gap-2">
						<Input type="file" accept="image/png,image/jpeg" onchange={upload} aria-label={$t`Upload a scan`} />
						{#if layout.data.hasBackground}<Button variant="ghost" size="icon-sm" title={$t`Remove the scan`} aria-label={$t`Remove the scan`} onclick={removeBackground}><Trash /></Button>{/if}
					</div>
				</FormField>
			</aside>
		</div>
	</div>
{/if}
