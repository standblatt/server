<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import Copy from '@lucide/svelte/icons/copy';
	import MousePointer from '@lucide/svelte/icons/mouse-pointer-2';
	import Trash from '@lucide/svelte/icons/trash-2';
	import type { LayoutElement, LibraryImage } from '../api';
	import { Button } from '../components/ui/button';
	import { Checkbox } from '../components/ui/checkbox';
	import { Input } from '../components/ui/input';
	import * as NativeSelect from '../components/ui/native-select';
	import FormField from '../FormField.svelte';
	import type { FormEditor } from './editor.svelte';

	let { editor, images = [] }: { editor: FormEditor; images?: LibraryImage[] } = $props();

	const sel = $derived(editor.selection);
	const item = $derived(sel.length === 1 ? sel[0] : undefined);
	let before: string | null = null;

	/** Remembers the state when a field gets the focus. */
	const focusin = () => (before = editor.begin());
	/** Records the edit of the field that lost the focus as one undo step. */
	function focusout() {
		if (before !== null) editor.commit(before);
		before = null;
	}

	/** Sets the height of an image element from its width and the image's aspect ratio. */
	function keepRatio(it: LayoutElement) {
		const img = images.find((i) => i.id === it.image);
		if (!img) return;
		editor.remember();
		it.h = Math.round((Number(it.w) * img.height * 10) / img.width) / 10;
	}

	/** Sets a property on every selected element as one undo step. */
	function setAll<K extends keyof LayoutElement>(key: K, value: LayoutElement[K]) {
		if (before === null) editor.remember();
		for (const i of sel) (i as LayoutElement)[key] = value;
	}
</script>

<div class="grid content-start gap-3" onfocusin={focusin} onfocusout={focusout} role="group" aria-label={$t`Properties`}>
	{#if !sel.length}
		<p class="grid justify-items-center gap-2 py-6 text-center text-sm text-muted-foreground">
			<MousePointer class="size-5" />
			{$t`Select an element on the sheet, or add one from the fields. Right-click for more.`}
		</p>
	{:else}
		{#if item?.type === 'image'}
			<FormField label={$t`Image`}>
				<NativeSelect.Root class="w-full" value={item.image} onchange={(e) => setAll('image', Number(e.currentTarget.value))}>
					{#each images as img (img.id)}<NativeSelect.Option value={img.id}>{img.name}</NativeSelect.Option>{/each}
				</NativeSelect.Root>
			</FormField>
			<div class="grid grid-cols-2 gap-2">
				<FormField label={$t`x (mm)`}><Input type="number" step="0.1" bind:value={item.x} /></FormField>
				<FormField label={$t`y (mm)`}><Input type="number" step="0.1" bind:value={item.y} /></FormField>
				<FormField label={$t`Width (mm)`}><Input type="number" step="0.1" bind:value={item.w} /></FormField>
				<FormField label={$t`Height (mm)`}><Input type="number" step="0.1" bind:value={item.h} /></FormField>
			</div>
			<Button variant="outline" size="sm" onclick={() => keepRatio(item)}>{$t`Keep the aspect ratio`}</Button>
		{:else if item}
			<FormField label={$t`Type`}>
				<NativeSelect.Root class="w-full" value={item.type} onchange={(e) => setAll('type', e.currentTarget.value as LayoutElement['type'])}>
					<NativeSelect.Option value="text">{$t`Text`}</NativeSelect.Option>
					<NativeSelect.Option value="qr">{$t`QR code`}</NativeSelect.Option>
					<NativeSelect.Option value="barcode">{$t`Barcode`}</NativeSelect.Option>
					<NativeSelect.Option value="line">{$t`Line`}</NativeSelect.Option>
				</NativeSelect.Root>
			</FormField>
			{#if item.type !== 'line'}
				<FormField label={$t`Value`} hint={$t`Placeholders in double curly braces are filled in when printing; add them from the fields.`}>
					<Input bind:value={item.value} />
				</FormField>
			{/if}
			<div class="grid grid-cols-2 gap-2">
				<FormField label={$t`x (mm)`}><Input type="number" step="0.1" bind:value={item.x} /></FormField>
				<FormField label={$t`y (mm)`}><Input type="number" step="0.1" bind:value={item.y} /></FormField>
				<FormField label={item.type === 'line' ? $t`Length x (mm)` : $t`Width (mm)`}><Input type="number" step="0.1" bind:value={item.w} /></FormField>
				{#if item.type !== 'qr'}<FormField label={item.type === 'line' ? $t`Length y (mm)` : $t`Height (mm)`}><Input type="number" step="0.1" bind:value={item.h} /></FormField>{/if}
			</div>
		{:else}
			<p class="text-sm font-medium">{$plural(sel.length, { one: '# element selected', other: '# elements selected' })}</p>
		{/if}

		{#if sel.some((i) => i.type === 'text')}
			<div class="grid grid-cols-2 items-end gap-2">
				<FormField label={$t`Font size (pt)`}>
					<Input type="number" min="4" max="72" value={sel[0].size || 10} oninput={(e) => sel.forEach((i) => (i.size = Number(e.currentTarget.value)))} />
				</FormField>
				<label class="flex h-9 items-center gap-2 text-sm">
					<Checkbox checked={sel.every((i) => !!i.bold)} onCheckedChange={(v) => setAll('bold', v)} />
					{$t`Bold`}
				</label>
				<label class="flex h-9 items-center gap-2 text-sm">
					<Checkbox checked={sel.every((i) => !!i.italic)} onCheckedChange={(v) => setAll('italic', v)} />
					{$t`Italic`}
				</label>
				<label class="flex h-9 items-center gap-2 text-sm" title={$t`Breaks the text into lines within the width; the height is the line height.`}>
					<Checkbox checked={sel.every((i) => !!i.wrap)} onCheckedChange={(v) => setAll('wrap', v)} />
					{$t`Wrap text`}
				</label>
			</div>
			<FormField label={$t`Alignment`}>
				<NativeSelect.Root class="w-full" value={sel[0].align || 'L'} onchange={(e) => setAll('align', e.currentTarget.value as LayoutElement['align'])}>
					<NativeSelect.Option value="L">{$t`left`}</NativeSelect.Option>
					<NativeSelect.Option value="C">{$t`centre`}</NativeSelect.Option>
					<NativeSelect.Option value="R">{$t`right`}</NativeSelect.Option>
				</NativeSelect.Root>
			</FormField>
		{/if}

		{#if item}
			<fieldset class="grid grid-cols-3 gap-2 rounded-md border p-2">
				<legend class="px-1 text-xs font-medium text-muted-foreground">{$t`Repeat`}</legend>
				<FormField label={$t`Times`}><Input type="number" min="0" max="100" bind:value={item.repeat} /></FormField>
				<FormField label={$t`dx (mm)`}><Input type="number" step="0.1" bind:value={item.dx} /></FormField>
				<FormField label={$t`dy (mm)`}><Input type="number" step="0.1" bind:value={item.dy} /></FormField>
				<p class="col-span-3 text-xs text-muted-foreground">{$t`Prints the element several times, moved by dx/dy; # in a placeholder counts 1, 2, 3 …`}</p>
			</fieldset>
		{/if}

		<div class="flex gap-2">
			<Button size="sm" variant="outline" class="flex-1" onclick={() => editor.duplicate()}><Copy />{$t`Duplicate`}</Button>
			<Button size="sm" variant="outline" class="flex-1 text-destructive" onclick={() => editor.remove()}><Trash />{$t`Delete`}</Button>
		</div>
	{/if}
</div>
