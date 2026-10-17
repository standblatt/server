<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { LayoutElement } from './api';

	let {
		width,
		height,
		elements = $bindable(),
		selected = $bindable(-1),
		background,
		samples,
		onadd
	}: {
		width: number;
		height: number;
		elements: LayoutElement[];
		selected?: number;
		background?: string;
		samples: Record<string, string>;
		onadd: (x: number, y: number) => void;
	} = $props();

	let box = $state(0);
	let sheet = $state<HTMLDivElement>();
	const scale = $derived(box / width);
	const mmPerPoint = 0.3528;

	/** Rounds to half millimetres, the grid of the editor. */
	const snap = (v: number) => Math.round(v * 2) / 2;
	/** A number field of an element, 0 when empty (inputs hand back strings while typing). */
	const n = (v: unknown) => Number(v) || 0;

	/** The text an element shows: placeholders become their examples (# the count n); unknown ones stay visible. */
	function fill(value: string | undefined, count: number) {
		return (value ?? '').replace(/\{\{\s*([A-Za-z0-9_.#]+)\s*\}\}/g, (_, k: string) => {
			const key = k.replace(/#/g, String(count));
			return samples[key] ?? `{{${key}}}`;
		});
	}

	/** Screen pixels per millimetre for pointer maths; differs from scale where the page is zoomed (CSS zoom). */
	const pointerScale = () => (sheet?.getBoundingClientRect().width ?? 0) / width;

	type Drag = { i: number; resize: boolean; px: number; py: number; x: number; y: number; w: number; h: number; scale: number };
	let drag: Drag | null = null;

	/** Starts moving element i, or resizing it when the handle was grabbed. */
	function down(e: PointerEvent, i: number, resize = false) {
		e.stopPropagation();
		e.preventDefault();
		selected = i;
		sheet?.focus();
		const el = elements[i];
		drag = { i, resize, px: e.clientX, py: e.clientY, x: n(el.x), y: n(el.y), w: n(el.w), h: n(el.h), scale: pointerScale() };
		(e.currentTarget as Element).setPointerCapture(e.pointerId);
	}

	/** Moves or resizes the dragged element along the pointer, on the half millimetre grid. */
	function move(e: PointerEvent) {
		if (!drag || !drag.scale) return;
		const dx = (e.clientX - drag.px) / drag.scale;
		const dy = (e.clientY - drag.py) / drag.scale;
		const el = elements[drag.i];
		if (!drag.resize) {
			el.x = snap(drag.x + dx);
			el.y = snap(drag.y + dy);
		} else if (el.type === 'line') {
			el.w = snap(drag.w + dx);
			el.h = snap(drag.h + dy);
		} else {
			el.w = Math.max(1, snap(drag.w + dx));
			el.h = Math.max(1, snap(drag.h + dy));
		}
	}

	/** Ends a drag. */
	const up = () => (drag = null);

	/** Arrow keys move the selected element by half a millimetre, with Shift by five. */
	function key(e: KeyboardEvent) {
		const el = elements[selected];
		const step = e.shiftKey ? 5 : 0.5;
		const d = ({ ArrowLeft: [-step, 0], ArrowRight: [step, 0], ArrowUp: [0, -step], ArrowDown: [0, step] } as Record<string, number[]>)[e.key];
		if (!el || !d) return;
		e.preventDefault();
		el.x = snap(n(el.x) + d[0]);
		el.y = snap(n(el.y) + d[1]);
	}

	/** A double click on the empty sheet adds a text element at that spot. */
	function dblclick(e: MouseEvent) {
		const ps = pointerScale();
		if (e.target !== e.currentTarget || !ps) return;
		const r = (e.currentTarget as HTMLElement).getBoundingClientRect();
		onadd(snap((e.clientX - r.left) / ps), snap((e.clientY - r.top) / ps));
	}

	/** The copies of an element: one, or Repeat of them moved by dx/dy. */
	const copies = (el: LayoutElement) => Array.from({ length: Math.max(n(el.repeat), 1) }, (_, k) => ({ k, x: n(el.x) + k * n(el.dx), y: n(el.y) + k * n(el.dy) }));
	/** Width and height of an element in millimetres as printed: text without a width runs to the right edge. */
	function size(el: LayoutElement, x: number) {
		const font = n(el.size) || 10;
		if (el.type === 'qr') return { w: n(el.w), h: n(el.w) };
		if (el.type === 'text') return { w: n(el.w) || width - x, h: n(el.h) || font * 0.45 };
		return { w: n(el.w), h: n(el.h) };
	}
</script>

<!-- The sheet takes focus for the arrow keys; the table next to it is the accessible way to edit. -->
<!-- svelte-ignore a11y_no_noninteractive_tabindex, a11y_no_noninteractive_element_interactions -->
<div
	bind:this={sheet}
	bind:clientWidth={box}
	class="relative w-full touch-none overflow-hidden rounded-lg border bg-white outline-none select-none focus-visible:ring-2 focus-visible:ring-ring"
	style:aspect-ratio="{width} / {height}"
	tabindex="0"
	role="application"
	aria-label={$t`Sheet`}
	onkeydown={key}
	onpointerdown={() => (selected = -1)}
	ondblclick={dblclick}
>
	{#if background}<img src={background} alt="" class="pointer-events-none absolute inset-0 h-full w-full" />{/if}
	<svg class="pointer-events-none absolute inset-0 h-full w-full" viewBox="0 0 {width} {height}" aria-hidden="true">
		{#each elements as el, i (i)}
			{#if el.type === 'line'}
				{#each copies(el) as c (c.k)}
					<line x1={c.x} y1={c.y} x2={c.x + n(el.w)} y2={c.y + n(el.h)} stroke={selected === i ? '#0284c7' : 'black'} stroke-width={selected === i ? 0.6 : 0.3} />
				{/each}
			{/if}
		{/each}
	</svg>
	{#each elements as el, i (i)}
		{#each copies(el) as c (c.k)}
			{@const s = size(el, c.x)}
			{#if el.type === 'line'}
				<div
					class="absolute cursor-move"
					style:left="{(Math.min(c.x, c.x + n(el.w)) - 1) * scale}px"
					style:top="{(Math.min(c.y, c.y + n(el.h)) - 1) * scale}px"
					style:width="{(Math.abs(n(el.w)) + 2) * scale}px"
					style:height="{(Math.abs(n(el.h)) + 2) * scale}px"
					role="button"
					tabindex="-1"
					aria-label={$t`Element ${i + 1}`}
					onpointerdown={(e) => down(e, i)}
					onpointermove={move}
					onpointerup={up}
				></div>
			{:else}
				<div
					class="absolute flex cursor-move items-center overflow-hidden leading-none whitespace-nowrap text-black {selected === i ? 'outline-2 outline-sky-600' : 'outline-1 outline-sky-400/50 outline-dashed hover:outline-sky-600'} {el.align === 'C' ? 'justify-center' : el.align === 'R' ? 'justify-end' : 'justify-start'}"
					style:left="{c.x * scale}px"
					style:top="{c.y * scale}px"
					style:width="{s.w * scale}px"
					style:height="{s.h * scale}px"
					style:font-size="{(n(el.size) || 10) * mmPerPoint * scale}px"
					style:font-weight={el.bold ? 700 : 400}
					style:font-family="Helvetica, Arial, sans-serif"
					style:background={el.type === 'barcode' ? 'repeating-linear-gradient(90deg, #000 0 1px, transparent 1px 3px)' : el.type === 'qr' ? 'repeating-conic-gradient(#000 0 25%, transparent 0 50%) 0 0 / 20% 20%' : undefined}
					role="button"
					tabindex="-1"
					aria-label={$t`Element ${i + 1}`}
					data-element={i}
					onpointerdown={(e) => down(e, i)}
					onpointermove={move}
					onpointerup={up}
				>
					{#if el.type === 'text'}{fill(el.value, c.k + 1)}{/if}
					{#if selected === i && c.k === 0}
						<span
							class="absolute right-0 bottom-0 size-2.5 cursor-nwse-resize bg-sky-600"
							role="button"
							tabindex="-1"
							aria-label={$t`Resize element ${i + 1}`}
							onpointerdown={(e) => down(e, i, true)}
							onpointermove={move}
							onpointerup={up}
						></span>
					{/if}
				</div>
			{/if}
		{/each}
	{/each}
</div>
