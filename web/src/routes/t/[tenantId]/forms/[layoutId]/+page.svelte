<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Copy from '@lucide/svelte/icons/copy';
	import Plus from '@lucide/svelte/icons/plus';
	import Printer from '@lucide/svelte/icons/printer';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call, type LayoutElement } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import LayoutCanvas from '$lib/LayoutCanvas.svelte';
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
	const fields = createQuery(() => ({
		queryKey: ['tenants', tenantId, 'layouts', 'fields'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/layouts/fields', { params: { path: { tenantId } } }))
	}));

	let name = $state('');
	let width = $state(148);
	let height = $state(210);
	let offsetX = $state(0);
	let offsetY = $state(0);
	let elements = $state<LayoutElement[]>([]);
	let loaded = $state(0);
	let version = $state(0);
	let selected = $state(-1);
	let view = $state<'sheet' | 'pdf'>('sheet');
	let error = $state('');
	let busy = $state(false);

	$effect(() => {
		const l = layout.data;
		if (!l || loaded === l.id) return;
		loaded = l.id;
		name = l.name;
		width = l.width;
		height = l.height;
		offsetX = l.offsetX ?? 0;
		offsetY = l.offsetY ?? 0;
		elements = l.elements.map((e) => ({ ...e }));
	});

	const base = $derived(`/api/tenants/${tenantId}/layouts/${layoutId}`);
	const preview = $derived(`${base}/preview?v=${version}`);
	const samples = $derived(Object.fromEntries((fields.data ?? []).map((f) => [f.key, f.example])));

	/** A new text element at x/y, or below the last one, and selects it. */
	function addElement(x?: number, y?: number) {
		const last = elements.at(-1);
		elements.push({ type: 'text', x: x ?? last?.x ?? 10, y: y ?? (Number(last?.y ?? 5) + 8), w: 0, h: 0, size: 10, bold: false, align: 'L', value: '', repeat: 0, dx: 0, dy: 0 });
		selected = elements.length - 1;
	}

	/** Puts a placeholder into the value of the element edited last (a shot as shot# when it repeats), or copies it. */
	async function usePlaceholder(key: string) {
		const el = elements[selected];
		const shot = key.match(/^(s\d+\.shot)\d+$/);
		const text = `{{${el?.repeat && shot ? `${shot[1]}#` : key}}}`;
		if (el) {
			el.value = (el.value ?? '') + text;
			return;
		}
		await navigator.clipboard?.writeText(text);
		toast.success($t`${text} copied`);
	}

	/** Saves name, size, offset and elements; the preview reloads. */
	async function save(e?: SubmitEvent) {
		e?.preventDefault();
		error = '';
		busy = true;
		try {
			await call(
				api.PUT('/api/tenants/{tenantId}/layouts/{layoutId}', {
					...path,
					body: { name: name.trim(), width: Number(width), height: Number(height), offsetX: Number(offsetX), offsetY: Number(offsetY), elements: $state.snapshot(elements).map(numbers) }
				})
			);
			version++;
			await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId, 'layouts'] });
			toast.success($t`Print form saved`);
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}

	/** Number inputs come back as strings when typed into; the API wants numbers. */
	function numbers(e: LayoutElement): LayoutElement {
		return { ...e, x: Number(e.x), y: Number(e.y), w: Number(e.w ?? 0), h: Number(e.h ?? 0), size: Number(e.size ?? 0), repeat: Number(e.repeat ?? 0), dx: Number(e.dx ?? 0), dy: Number(e.dy ?? 0) };
	}

	/** Uploads a scan of the empty form for the preview. */
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

	const num = 'h-8 w-16 px-1.5 tabular-nums';
</script>

<PageHeader help="user-guide/print-forms" title={layout.data?.name ?? $t`Print form`} description={$t`Positions and sizes in millimetres from the top left corner of the sheet.`}>
	{#snippet actions()}
		<Button variant="outline" href="{base}/preview?print=true&v={version}" target="_blank" onclick={pdfClick}><Printer />{$t`Test print`}</Button>
		<Button onclick={() => save()} disabled={busy}>{$t`Save`}</Button>
	{/snippet}
</PageHeader>

{#if layout.data}
	<div class="grid gap-4 *:min-w-0 xl:grid-cols-[minmax(0,1fr)_28rem]">
		<form class="grid content-start gap-4" onsubmit={save}>
			<div class="grid gap-3 rounded-lg border bg-card p-4 sm:grid-cols-3">
				<FormField label={$t`Name`} class="sm:col-span-3"><Input bind:value={name} required maxlength={100} /></FormField>
				<FormField label={$t`Width (mm)`}><Input type="number" min="20" max="500" step="0.5" bind:value={width} required /></FormField>
				<FormField label={$t`Height (mm)`}><Input type="number" min="20" max="500" step="0.5" bind:value={height} required /></FormField>
				<div></div>
				<FormField label={$t`Printer offset right (mm)`} hint={$t`Moves everything when the printer prints off target.`}><Input type="number" min="-50" max="50" step="0.5" bind:value={offsetX} /></FormField>
				<FormField label={$t`Printer offset down (mm)`}><Input type="number" min="-50" max="50" step="0.5" bind:value={offsetY} /></FormField>
				<FormField label={$t`Scan of the empty form`} hint={$t`PNG or JPEG of the whole sheet; only shown in the preview, never printed.`}>
					<div class="flex items-center gap-2">
						<Input type="file" accept="image/png,image/jpeg" onchange={upload} aria-label={$t`Upload a scan`} />
						{#if layout.data.hasBackground}<Button variant="ghost" size="icon-sm" title={$t`Remove the scan`} aria-label={$t`Remove the scan`} onclick={removeBackground}><Trash /></Button>{/if}
					</div>
				</FormField>
			</div>

			<section class="rounded-lg border bg-card" aria-label={$t`Elements`}>
				<div class="flex items-center justify-between border-b px-4 py-2.5">
					<h2 class="text-sm font-semibold">{$t`Elements`}</h2>
					<Button size="sm" variant="outline" onclick={() => addElement()}><Plus />{$t`Add element`}</Button>
				</div>
				<div class="overflow-x-auto" data-swipe>
					<table class="w-full text-sm">
						<thead class="text-left text-xs text-muted-foreground">
							<tr>
								<th class="px-2 py-2">{$t`Type`}</th><th class="px-2">{$t`Value`}</th><th class="px-2">x</th><th class="px-2">y</th><th class="px-2">{$t`Width`}</th><th class="px-2">{$t`Height`}</th>
								<th class="px-2">{$t`Size`}</th><th class="px-2">{$t`Bold`}</th><th class="px-2">{$t`Align`}</th><th class="px-2">{$t`Repeat`}</th><th class="px-2">dx</th><th class="px-2">dy</th><th></th>
							</tr>
						</thead>
						<tbody>
							{#each elements as el, i (i)}
								<tr class="border-t {selected === i ? 'bg-sky-500/10' : ''}" aria-label={$t`Element ${i + 1}`} onfocusin={() => (selected = i)}>
									<td class="px-2 py-1.5">
										<NativeSelect.Root size="sm" bind:value={el.type} aria-label={$t`Type of element ${i + 1}`}>
											<NativeSelect.Option value="text">{$t`Text`}</NativeSelect.Option>
											<NativeSelect.Option value="qr">QR</NativeSelect.Option>
											<NativeSelect.Option value="barcode">{$t`Barcode`}</NativeSelect.Option>
											<NativeSelect.Option value="line">{$t`Line`}</NativeSelect.Option>
										</NativeSelect.Root>
									</td>
									<td class="px-2"><Input class="h-8 min-w-48" bind:value={el.value} aria-label={$t`Value of element ${i + 1}`} disabled={el.type === 'line'} /></td>
									<td class="px-2"><Input class={num} type="number" step="0.5" bind:value={el.x} aria-label={$t`x of element ${i + 1}`} /></td>
									<td class="px-2"><Input class={num} type="number" step="0.5" bind:value={el.y} aria-label={$t`y of element ${i + 1}`} /></td>
									<td class="px-2"><Input class={num} type="number" step="0.5" bind:value={el.w} aria-label={$t`Width of element ${i + 1}`} /></td>
									<td class="px-2"><Input class={num} type="number" step="0.5" bind:value={el.h} aria-label={$t`Height of element ${i + 1}`} /></td>
									<td class="px-2"><Input class={num} type="number" min="4" max="72" bind:value={el.size} aria-label={$t`Font size of element ${i + 1}`} /></td>
									<td class="px-2 text-center"><Checkbox checked={!!el.bold} onCheckedChange={(v) => (el.bold = v)} aria-label={$t`Bold element ${i + 1}`} /></td>
									<td class="px-2">
										<NativeSelect.Root size="sm" bind:value={el.align} aria-label={$t`Alignment of element ${i + 1}`}>
											<NativeSelect.Option value="L">{$t`left`}</NativeSelect.Option>
											<NativeSelect.Option value="C">{$t`centre`}</NativeSelect.Option>
											<NativeSelect.Option value="R">{$t`right`}</NativeSelect.Option>
										</NativeSelect.Root>
									</td>
									<td class="px-2"><Input class={num} type="number" min="0" max="100" bind:value={el.repeat} aria-label={$t`Repeat element ${i + 1}`} /></td>
									<td class="px-2"><Input class={num} type="number" step="0.5" bind:value={el.dx} aria-label={$t`dx of element ${i + 1}`} /></td>
									<td class="px-2"><Input class={num} type="number" step="0.5" bind:value={el.dy} aria-label={$t`dy of element ${i + 1}`} /></td>
									<td class="px-2 whitespace-nowrap">
										<Button variant="ghost" size="icon-sm" title={$t`Duplicate`} aria-label={$t`Duplicate element ${i + 1}`} onclick={() => elements.splice(i + 1, 0, { ...$state.snapshot(el), y: Number(el.y) + 8 })}><Copy /></Button>
										<Button variant="ghost" size="icon-sm" title={$t`Remove`} aria-label={$t`Remove element ${i + 1}`} onclick={() => (elements.splice(i, 1), (selected = -1))}><Trash /></Button>
									</td>
								</tr>
							{:else}
								<tr><td colspan="13" class="px-4 py-6 text-center text-muted-foreground">{$t`No elements yet.`}</td></tr>
							{/each}
						</tbody>
					</table>
				</div>
			</section>

			<section class="rounded-lg border bg-card p-4" aria-label={$t`Placeholders`}>
				<h2 class="text-sm font-semibold">{$t`Placeholders`}</h2>
				<p class="mb-2 text-xs text-muted-foreground">{$t`Click to add one to the element edited last. With Repeat, # counts 1, 2, 3 …: s1.shot# prints the shots of series 1 one after the other.`}</p>
				<div class="flex flex-wrap gap-1.5">
					{#each fields.data ?? [] as f (f.key)}
						<button type="button" class="rounded-md border px-2 py-0.5 font-mono text-xs hover:bg-muted" title={f.example} onclick={() => usePlaceholder(f.key)}>{f.key}</button>
					{/each}
				</div>
			</section>
			<FormError message={error} />
		</form>

		<section class="grid content-start gap-2 xl:sticky xl:top-4" aria-label={$t`Preview`}>
			<div class="flex gap-1 border-b" role="tablist" aria-label={$t`Preview`}>
				{#each [['sheet', $t`Sheet`], ['pdf', $t`PDF preview`]] as [key, label] (key)}
					<button type="button" role="tab" aria-selected={view === key} class="-mb-px border-b-2 px-3 py-1.5 text-sm font-medium {view === key ? 'border-primary' : 'border-transparent text-muted-foreground hover:text-foreground'}" onclick={() => (view = key as 'sheet' | 'pdf')}>{label}</button>
				{/each}
			</div>
			{#if view === 'sheet'}
				<LayoutCanvas width={Number(width) || 1} height={Number(height) || 1} bind:elements bind:selected {samples} background={layout.data.hasBackground ? `${base}/background?v=${version}` : undefined} onadd={addElement} />
				<p class="text-xs text-muted-foreground">{$t`Drag elements, resize them at the corner, move the selected one with the arrow keys (Shift: 5 mm). Double-click the empty sheet to add text there.`}</p>
			{:else}
				<iframe title={$t`Preview`} src={preview} class="w-full rounded-lg border bg-white" style:aspect-ratio="{width} / {height}"></iframe>
				<p class="text-xs text-muted-foreground">{$t`With made-up values and the scan, without the printer offset. Save to refresh.`}</p>
			{/if}
		</section>
	</div>
{/if}
