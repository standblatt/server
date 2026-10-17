<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import * as ContextMenu from '../components/ui/context-menu';
	import { fill, num, tenth, type FormEditor, type Item } from './editor.svelte';

	let {
		editor,
		samples,
		names = false,
		background,
		zoom = 1,
		onzoom
	}: { editor: FormEditor; samples: Record<string, string>; names?: boolean; background?: string; zoom?: number; onzoom: (z: number) => void } = $props();

	const scale = $derived(4 * zoom);
	const mmPerPoint = 0.3528;
	const ruler = 22;
	const snapPx = 8;
	const W = $derived(Math.max(num(editor.width), 1));
	const H = $derived(Math.max(num(editor.height), 1));

	type Box = { x: number; y: number; w: number; h: number };
	type Kind = 'move' | 'size' | 'end' | 'spread';
	type Drag = { kind: Kind; id: string; start: { x: number; y: number }; orig: Map<string, Box & { dx: number; dy: number }>; boxes: Map<string, Box>; moved: boolean; solo: boolean };

	let sheet = $state<HTMLElement | null>(null);
	let guides = $state<{ axis: 'x' | 'y'; at: number }[]>([]);
	let band = $state<Box | null>(null);
	let bandStart = { x: 0, y: 0 };
	let bandBase: string[] = [];
	let drag: Drag | null = null;

	/** Screen pixels per millimetre; differs from scale under the CSS zoom of the content pane. */
	const screen = () => (sheet?.getBoundingClientRect().width ?? 0) / W;
	/** The millimetre position of a pointer on the sheet. */
	function at(e: { clientX: number; clientY: number }) {
		const r = sheet!.getBoundingClientRect();
		const s = screen();
		return { x: (e.clientX - r.left) / s, y: (e.clientY - r.top) / s };
	}
	/** Keeps v between lo and hi. */
	const clamp = (v: number, lo: number, hi: number) => Math.min(Math.max(v, lo), Math.max(lo, hi));
	/** The box of an element in millimetres as drawn on the sheet. */
	function box(it: Item): Box {
		if (it.type === 'line') {
			const x2 = num(it.x) + num(it.w);
			const y2 = num(it.y) + num(it.h);
			return { x: Math.min(num(it.x), x2), y: Math.min(num(it.y), y2), w: Math.abs(num(it.w)), h: Math.abs(num(it.h)) };
		}
		const node = sheet?.querySelector<HTMLElement>(`[data-id="${it.id}"]`);
		return { x: num(it.x), y: num(it.y), w: node ? node.offsetWidth / scale : num(it.w), h: node ? node.offsetHeight / scale : num(it.h) };
	}
	/** The copies an element prints: one, or Repeat of them moved by dx/dy. */
	const copies = (it: Item) => Array.from({ length: Math.max(num(it.repeat), 1) }, (_, k) => ({ k, x: num(it.x) + k * num(it.dx), y: num(it.y) + k * num(it.dy) }));
	/** Takes the focus away from inputs, so the arrow keys move elements. */
	const blur = () => (document.activeElement as HTMLElement | null)?.blur?.();

	/** Starts moving the selection (with Shift or Ctrl adding to it), or resizing, stretching or spreading one element. */
	function grab(e: PointerEvent, it: Item, kind: Kind = 'move') {
		if (e.button !== 0) return;
		e.stopPropagation();
		e.preventDefault();
		blur();
		const solo = kind === 'move' && !(e.shiftKey || e.ctrlKey || e.metaKey) && editor.selected.length > 1;
		if (kind === 'move') {
			if (e.shiftKey || e.ctrlKey || e.metaKey) {
				editor.selected = editor.selected.includes(it.id) ? editor.selected.filter((i) => i !== it.id) : [...editor.selected, it.id];
				if (!editor.selected.includes(it.id)) return;
			} else if (!editor.selected.includes(it.id)) {
				editor.selected = [it.id];
			}
		} else {
			editor.selected = [it.id];
		}
		const ids = kind === 'move' ? editor.selected : [it.id];
		const orig = new Map<string, Box & { dx: number; dy: number }>();
		const boxes = new Map<string, Box>();
		for (const i of editor.items) {
			if (!ids.includes(i.id)) continue;
			const b = box(i);
			boxes.set(i.id, b);
			orig.set(i.id, { x: num(i.x), y: num(i.y), w: i.type === 'line' ? num(i.w) : num(i.w) || b.w, h: i.type === 'line' ? num(i.h) : num(i.h) || b.h, dx: num(i.dx), dy: num(i.dy) });
		}
		drag = { kind, id: it.id, start: at(e), orig, boxes, moved: false, solo };
	}

	/** Snaps a moved box to the edges and centres of the other elements and of the sheet, else to whole millimetres. */
	function snap(b: Box, dx: number, dy: number, moving: Map<string, Box>) {
		const tolerance = snapPx / screen();
		const xs = [0, W / 2, W];
		const ys = [0, H / 2, H];
		for (const i of editor.items) {
			if (moving.has(i.id)) continue;
			const o = box(i);
			xs.push(o.x, o.x + o.w / 2, o.x + o.w);
			ys.push(o.y, o.y + o.h / 2, o.y + o.h);
		}
		/** The nearest line to any edge within the tolerance. */
		const nearest = (edges: number[], lines: number[]) => {
			let best: { d: number; line: number } | null = null;
			for (const e of edges) for (const l of lines) if (Math.abs(l - e) <= tolerance && (!best || Math.abs(l - e) < Math.abs(best.d))) best = { d: l - e, line: l };
			return best;
		};
		const x = b.x + dx;
		const y = b.y + dy;
		const sx = nearest([x, x + b.w / 2, x + b.w], xs);
		const sy = nearest([y, y + b.h / 2, y + b.h], ys);
		const found: { axis: 'x' | 'y'; at: number }[] = [];
		if (sx) found.push({ axis: 'x', at: sx.line });
		if (sy) found.push({ axis: 'y', at: sy.line });
		return { dx: sx ? dx + sx.d : Math.round(x) - b.x, dy: sy ? dy + sy.d : Math.round(y) - b.y, guides: found };
	}

	/** Follows the pointer: the rubber band, or the drag started by grab. */
	function move(e: PointerEvent) {
		if (band) {
			const p = at(e);
			band = { x: Math.min(p.x, bandStart.x), y: Math.min(p.y, bandStart.y), w: Math.abs(p.x - bandStart.x), h: Math.abs(p.y - bandStart.y) };
			const b = band;
			const hit = editor.items.filter((i) => {
				const o = box(i);
				return o.x < b.x + b.w && o.x + Math.max(o.w, 0.5) > b.x && o.y < b.y + b.h && o.y + Math.max(o.h, 0.5) > b.y;
			});
			editor.selected = [...new Set([...bandBase, ...hit.map((i) => i.id)])];
			return;
		}
		if (!drag) return;
		const p = at(e);
		let dx = p.x - drag.start.x;
		let dy = p.y - drag.start.y;
		if (!drag.moved) {
			if (Math.hypot(dx, dy) * screen() < 3) return;
			editor.remember();
			drag.moved = true;
		}
		const it = editor.items.find((i) => i.id === drag!.id);
		const o = drag.orig.get(drag.id);
		if (!it || !o) return;
		if (drag.kind === 'size') {
			it.w = Math.max(1, tenth(o.w + dx));
			it.h = Math.max(1, tenth(o.h + dy));
			return;
		}
		if (drag.kind === 'end') {
			it.w = tenth(o.w + dx);
			it.h = tenth(o.h + dy);
			return;
		}
		if (drag.kind === 'spread') {
			const n = Math.max(num(it.repeat) - 1, 1);
			it.dx = tenth(o.dx + dx / n);
			it.dy = tenth(o.dy + dy / n);
			return;
		}
		const s = e.altKey ? null : snap(drag.boxes.get(drag.id)!, dx, dy, drag.boxes);
		if (s) ({ dx, dy } = s);
		guides = s?.guides ?? [];
		const all = [...drag.boxes.values()];
		dx = clamp(dx, -Math.min(...all.map((b) => b.x)), W - Math.max(...all.map((b) => b.x + b.w)));
		dy = clamp(dy, -Math.min(...all.map((b) => b.y)), H - Math.max(...all.map((b) => b.y + b.h)));
		for (const i of editor.items) {
			const io = drag.orig.get(i.id);
			if (!io) continue;
			i.x = tenth(io.x + dx);
			i.y = tenth(io.y + dy);
		}
	}

	/** Ends a drag or the rubber band; a click without moving on a selected element selects only it. */
	function up() {
		if (drag?.solo && !drag.moved) editor.selected = [drag.id];
		drag = null;
		band = null;
		guides = [];
	}

	/** Starts a rubber band on the empty sheet; without Shift or Ctrl it clears the selection first. */
	function startBand(e: PointerEvent) {
		if (e.button !== 0) return;
		blur();
		const additive = e.shiftKey || e.ctrlKey || e.metaKey;
		bandBase = additive ? editor.selected : [];
		if (!additive) editor.selected = [];
		bandStart = at(e);
		band = { ...bandStart, w: 0, h: 0 };
	}

	/** Lets a field from the palette be dropped onto the sheet. */
	function dragover(e: DragEvent) {
		if (!e.dataTransfer?.types.includes('text/standblatt-element')) return;
		e.preventDefault();
		e.dataTransfer.dropEffect = 'copy';
	}

	/** Adds the dropped field where the pointer is, on whole millimetres. */
	function drop(e: DragEvent) {
		const raw = e.dataTransfer?.getData('text/standblatt-element');
		if (!raw) return;
		e.preventDefault();
		const p = at(e);
		editor.add({ ...JSON.parse(raw), x: Math.round(p.x), y: Math.round(p.y) });
	}

	/** Keyboard shortcuts of the editor, unless the focus is in an input. */
	function key(e: KeyboardEvent) {
		const tag = (e.target as HTMLElement | null)?.tagName?.toLowerCase() ?? '';
		const typing = tag === 'input' || tag === 'textarea' || tag === 'select';
		const mod = e.ctrlKey || e.metaKey;
		const k = e.key.toLowerCase();
		if (typing) return;
		if (mod && k === 'z') {
			e.preventDefault();
			if (e.shiftKey) editor.redo();
			else editor.undo();
			return;
		}
		if (mod && k === 'y') {
			e.preventDefault();
			editor.redo();
			return;
		}
		if (mod && k === 'a') {
			e.preventDefault();
			editor.selected = editor.items.map((i) => i.id);
			return;
		}
		if (e.key === 'Escape') {
			editor.selected = [];
			return;
		}
		if (!editor.selected.length) return;
		if (e.key === 'Delete' || e.key === 'Backspace') {
			e.preventDefault();
			editor.remove();
			return;
		}
		if (mod && k === 'd') {
			e.preventDefault();
			editor.duplicate();
			return;
		}
		const step = e.shiftKey ? 5 : 0.5;
		const d = ({ ArrowLeft: [-step, 0], ArrowRight: [step, 0], ArrowUp: [0, -step], ArrowDown: [0, step] } as Record<string, number[]>)[e.key];
		if (!d) return;
		e.preventDefault();
		editor.remember();
		for (const i of editor.selection) {
			i.x = tenth(num(i.x) + d[0]);
			i.y = tenth(num(i.y) + d[1]);
		}
	}

	/** Ctrl and the mouse wheel zoom the sheet. */
	function wheel(e: WheelEvent) {
		if (!e.ctrlKey && !e.metaKey) return;
		e.preventDefault();
		onzoom(zoom + (e.deltaY < 0 ? 0.1 : -0.1));
	}

	/** Lines up or spreads the selected elements, or centres them on the sheet. */
	function align(how: 'left' | 'hcenter' | 'right' | 'top' | 'vcenter' | 'bottom' | 'dist-h' | 'dist-v' | 'page-h' | 'page-v') {
		const sel = editor.selection.map((i) => ({ i, b: box(i) }));
		if (!sel.length) return;
		editor.remember();
		const L = Math.min(...sel.map((s) => s.b.x));
		const R = Math.max(...sel.map((s) => s.b.x + s.b.w));
		const T = Math.min(...sel.map((s) => s.b.y));
		const B = Math.max(...sel.map((s) => s.b.y + s.b.h));
		/** Moves an element so its box starts at x/y. */
		const to = (s: { i: Item; b: Box }, x: number, y: number) => {
			s.i.x = tenth(num(s.i.x) + x - s.b.x);
			s.i.y = tenth(num(s.i.y) + y - s.b.y);
		};
		/** Spaces the centres evenly between the outermost two along one axis. */
		const spread = (pos: 'x' | 'y', size: 'w' | 'h') => {
			const sorted = [...sel].sort((a, b) => a.b[pos] + a.b[size] / 2 - (b.b[pos] + b.b[size] / 2));
			const first = sorted[0].b[pos] + sorted[0].b[size] / 2;
			const last = sorted.at(-1)!.b[pos] + sorted.at(-1)!.b[size] / 2;
			sorted.forEach((s, n) => {
				const c = first + ((last - first) * n) / (sorted.length - 1) - s.b[size] / 2;
				to(s, pos === 'x' ? c : s.b.x, pos === 'y' ? c : s.b.y);
			});
		};
		for (const s of sel) {
			if (how === 'left') to(s, L, s.b.y);
			if (how === 'right') to(s, R - s.b.w, s.b.y);
			if (how === 'hcenter') to(s, (L + R) / 2 - s.b.w / 2, s.b.y);
			if (how === 'top') to(s, s.b.x, T);
			if (how === 'bottom') to(s, s.b.x, B - s.b.h);
			if (how === 'vcenter') to(s, s.b.x, (T + B) / 2 - s.b.h / 2);
			if (how === 'page-h') to(s, s.b.x + W / 2 - (L + R) / 2, s.b.y);
			if (how === 'page-v') to(s, s.b.x, s.b.y + H / 2 - (T + B) / 2);
		}
		if (how === 'dist-h' && sel.length > 2) spread('x', 'w');
		if (how === 'dist-v' && sel.length > 2) spread('y', 'h');
	}

	/** A right click on an element selects it unless it already is. */
	function context(it: Item) {
		if (!editor.selected.includes(it.id)) editor.selected = [it.id];
	}

	const n = $derived(editor.selected.length);
	const single = $derived(n === 1 ? editor.selection[0] : undefined);
	const ticksX = $derived(Array.from({ length: Math.floor(W) + 1 }, (_, i) => i));
	const ticksY = $derived(Array.from({ length: Math.floor(H) + 1 }, (_, i) => i));
</script>

<svelte:window onpointermove={move} onpointerup={up} onkeydown={key} />

<div class="overflow-auto rounded-lg border bg-muted/40 p-2" data-swipe onwheel={wheel}>
	<div class="relative" style:width="{ruler + W * scale}px" style:height="{ruler + H * scale}px">
		<svg class="absolute top-0 text-muted-foreground" style:left="{ruler}px" width={W * scale} height={ruler} aria-hidden="true">
			{#each ticksX as i (i)}
				<line x1={i * scale} x2={i * scale} y1={ruler} y2={ruler - (i % 10 === 0 ? 10 : i % 5 === 0 ? 6 : scale >= 4 ? 3 : 0)} stroke="currentColor" stroke-width="1" />
				{#if i % 10 === 0}<text x={i * scale + 2} y="10" font-size="9" fill="currentColor">{i}</text>{/if}
			{/each}
		</svg>
		<svg class="absolute left-0 text-muted-foreground" style:top="{ruler}px" width={ruler} height={H * scale} aria-hidden="true">
			{#each ticksY as i (i)}
				<line y1={i * scale} y2={i * scale} x1={ruler} x2={ruler - (i % 10 === 0 ? 10 : i % 5 === 0 ? 6 : scale >= 4 ? 3 : 0)} stroke="currentColor" stroke-width="1" />
				{#if i % 10 === 0 && i > 0}<text x="1" y={i * scale - 2} font-size="9" fill="currentColor">{i}</text>{/if}
			{/each}
		</svg>
		<ContextMenu.Root>
			<ContextMenu.Trigger>
				{#snippet child({ props })}
					<!-- The sheet is a pointer surface; the properties panel is the accessible way to edit. -->
					<!-- svelte-ignore a11y_no_static_element_interactions -->
					<div
						{...props}
						bind:this={sheet}
						class="absolute touch-none overflow-hidden bg-white shadow select-none"
						style:left="{ruler}px"
						style:top="{ruler}px"
						style:width="{W * scale}px"
						style:height="{H * scale}px"
						role="application"
						aria-label={$t`Sheet`}
						onpointerdown={(e) => {
							(props.onpointerdown as ((e: PointerEvent) => void) | undefined)?.(e);
							startBand(e);
						}}
						ondragover={dragover}
						ondrop={drop}
					>
						{#if background}<img src={background} alt="" class="pointer-events-none absolute inset-0 h-full w-full" />{/if}
						<svg class="pointer-events-none absolute inset-0 h-full w-full" viewBox="0 0 {W} {H}" preserveAspectRatio="none" aria-hidden="true">
							{#each editor.items as it (it.id)}
								{#if it.type === 'line'}
									{#each copies(it) as c (c.k)}
										<line x1={c.x} y1={c.y} x2={c.x + num(it.w)} y2={c.y + num(it.h)} stroke={editor.selected.includes(it.id) ? '#0284c7' : 'black'} stroke-width={editor.selected.includes(it.id) ? 0.5 : 0.25} />
									{/each}
								{/if}
							{/each}
							{#each guides as g, i (i)}
								{#if g.axis === 'x'}<line x1={g.at} x2={g.at} y1="0" y2={H} stroke="#db2777" stroke-width="0.2" stroke-dasharray="1 0.6" />
								{:else}<line y1={g.at} y2={g.at} x1="0" x2={W} stroke="#db2777" stroke-width="0.2" stroke-dasharray="1 0.6" />{/if}
							{/each}
						</svg>
						{#each editor.items as it (it.id)}
							{@const sel = editor.selected.includes(it.id)}
							{#each copies(it) as c (c.k)}
								{#if it.type === 'line'}
									<div
										class="absolute cursor-move"
										style:left="{(Math.min(c.x, c.x + num(it.w)) - 1) * scale}px"
										style:top="{(Math.min(c.y, c.y + num(it.h)) - 1) * scale}px"
										style:width="{(Math.abs(num(it.w)) + 2) * scale}px"
										style:height="{(Math.abs(num(it.h)) + 2) * scale}px"
										data-id={c.k === 0 ? it.id : undefined}
										role="presentation"
										onpointerdown={(e) => grab(e, it)}
										oncontextmenu={() => context(it)}
									></div>
									{#if single === it && c.k === 0}
										<span class="absolute size-2.5 -translate-1/2 cursor-crosshair rounded-full bg-sky-600" style:left="{(c.x + num(it.w)) * scale}px" style:top="{(c.y + num(it.h)) * scale}px" role="presentation" onpointerdown={(e) => grab(e, it, 'end')}></span>
									{/if}
								{:else}
									<div
										class="absolute flex items-center leading-none whitespace-nowrap text-black {c.k === 0 ? 'cursor-move' : 'opacity-70'} {sel ? 'outline-2 outline-sky-600' : 'outline-1 outline-sky-500/40 outline-dashed hover:outline-sky-600'} {it.align === 'C' ? 'justify-center' : it.align === 'R' ? 'justify-end' : 'justify-start'}"
										style:left="{c.x * scale}px"
										style:top="{c.y * scale}px"
										style:width={it.type === 'qr' || num(it.w) ? `${num(it.w) * scale}px` : undefined}
										style:height="{(it.type === 'qr' ? num(it.w) : num(it.h) || (num(it.size) || 10) * 0.45) * scale}px"
										style:font-size="{(num(it.size) || 10) * mmPerPoint * scale}px"
										style:font-weight={it.bold ? 700 : 400}
										style:font-family="Helvetica, Arial, sans-serif"
										style:background={it.type === 'barcode' ? 'repeating-linear-gradient(90deg, #000 0 1px, transparent 1px 3px)' : it.type === 'qr' ? 'repeating-conic-gradient(#000 0 25%, transparent 0 50%) 0 0 / 25% 25%' : undefined}
										data-id={c.k === 0 ? it.id : undefined}
										data-copy={c.k}
										role="presentation"
										onpointerdown={(e) => grab(e, it)}
										oncontextmenu={() => context(it)}
									>
										{#if it.type === 'text'}{fill(it.value, samples, c.k + 1, names)}{/if}
										{#if single === it && c.k === 0}
											<span class="absolute -right-1 -bottom-1 size-2.5 cursor-nwse-resize bg-sky-600" role="presentation" onpointerdown={(e) => grab(e, it, 'size')}></span>
										{/if}
										{#if single === it && c.k > 0 && c.k === num(it.repeat) - 1}
											<span class="absolute -top-1 -right-1 size-2.5 cursor-ew-resize rounded-full bg-pink-600" title={$t`Drag to change the spacing`} role="presentation" onpointerdown={(e) => grab(e, it, 'spread')}></span>
										{/if}
									</div>
								{/if}
							{/each}
						{/each}
						{#if band}
							<div class="pointer-events-none absolute border border-sky-600 bg-sky-500/10" style:left="{band.x * scale}px" style:top="{band.y * scale}px" style:width="{band.w * scale}px" style:height="{band.h * scale}px"></div>
						{/if}
					</div>
				{/snippet}
			</ContextMenu.Trigger>
			<ContextMenu.Content class="w-60">
				<ContextMenu.Item disabled={!n} onSelect={() => editor.duplicate()}>{$t`Duplicate`}<ContextMenu.Shortcut>Ctrl+D</ContextMenu.Shortcut></ContextMenu.Item>
				<ContextMenu.Item disabled={!n} onSelect={() => editor.toFront()}>{$t`Bring to front`}</ContextMenu.Item>
				<ContextMenu.Separator />
				<ContextMenu.Sub>
					<ContextMenu.SubTrigger disabled={n < 2}>{$t`Align`}</ContextMenu.SubTrigger>
					<ContextMenu.SubContent class="w-56">
						<ContextMenu.Item onSelect={() => align('left')}>{$t`Align left edges`}</ContextMenu.Item>
						<ContextMenu.Item onSelect={() => align('hcenter')}>{$t`Align horizontal centres`}</ContextMenu.Item>
						<ContextMenu.Item onSelect={() => align('right')}>{$t`Align right edges`}</ContextMenu.Item>
						<ContextMenu.Separator />
						<ContextMenu.Item onSelect={() => align('top')}>{$t`Align top edges`}</ContextMenu.Item>
						<ContextMenu.Item onSelect={() => align('vcenter')}>{$t`Align vertical centres`}</ContextMenu.Item>
						<ContextMenu.Item onSelect={() => align('bottom')}>{$t`Align bottom edges`}</ContextMenu.Item>
						<ContextMenu.Separator />
						<ContextMenu.Item disabled={n < 3} onSelect={() => align('dist-h')}>{$t`Distribute horizontally`}</ContextMenu.Item>
						<ContextMenu.Item disabled={n < 3} onSelect={() => align('dist-v')}>{$t`Distribute vertically`}</ContextMenu.Item>
					</ContextMenu.SubContent>
				</ContextMenu.Sub>
				<ContextMenu.Item disabled={!n} onSelect={() => align('page-h')}>{$t`Centre on the sheet horizontally`}</ContextMenu.Item>
				<ContextMenu.Item disabled={!n} onSelect={() => align('page-v')}>{$t`Centre on the sheet vertically`}</ContextMenu.Item>
				<ContextMenu.Separator />
				<ContextMenu.Item disabled={!n} variant="destructive" onSelect={() => editor.remove()}>{$t`Delete`}<ContextMenu.Shortcut>Del</ContextMenu.Shortcut></ContextMenu.Item>
			</ContextMenu.Content>
		</ContextMenu.Root>
	</div>
</div>
