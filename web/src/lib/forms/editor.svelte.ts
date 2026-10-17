import type { Layout, LayoutElement, LayoutInput } from '../api';

/** An element of the form being edited, with an ID that lives only in the editor. */
export type Item = LayoutElement & { id: string };

type Snapshot = { name: string; width: number; height: number; offsetX: number; offsetY: number; items: Item[] };

const historyMax = 100;

/** A number from a field that may still hold the string an input typed into it. */
export const num = (v: unknown) => Number(v) || 0;

/** Rounds millimetres to a tenth, the precision the editor stores. */
export const tenth = (v: number) => Math.round(v * 10) / 10;

/** A fresh editor ID. */
const newId = () => Math.random().toString(36).slice(2, 10);

/** The editor state of a print form: its fields, elements, selection and undo history. */
export class FormEditor {
	name = $state('');
	width = $state(148);
	height = $state(210);
	offsetX = $state(0);
	offsetY = $state(0);
	items = $state<Item[]>([]);
	selected = $state<string[]>([]);
	#past = $state<string[]>([]);
	#future = $state<string[]>([]);
	#clean = $state('');

	/** Takes a stored form as the clean starting point. */
	load(l: Layout) {
		this.name = l.name;
		this.width = l.width;
		this.height = l.height;
		this.offsetX = l.offsetX ?? 0;
		this.offsetY = l.offsetY ?? 0;
		this.items = l.elements.map((e) => ({ ...e, id: newId() }));
		this.selected = [];
		this.#past = [];
		this.#future = [];
		this.markClean();
	}

	/** The state as a string, for the history and to tell unsaved changes. */
	#snapshot(): string {
		const s: Snapshot = { name: this.name, width: num(this.width), height: num(this.height), offsetX: num(this.offsetX), offsetY: num(this.offsetY), items: $state.snapshot(this.items) };
		return JSON.stringify(s);
	}

	/** Puts a snapshot back. */
	#restore(json: string) {
		const s = JSON.parse(json) as Snapshot;
		Object.assign(this, { name: s.name, width: s.width, height: s.height, offsetX: s.offsetX, offsetY: s.offsetY, items: s.items });
		this.selected = this.selected.filter((id) => s.items.some((i) => i.id === id));
	}

	/** Records the state before a change, so it can be undone. */
	remember() {
		this.#past = [...this.#past.slice(-(historyMax - 1)), this.#snapshot()];
		this.#future = [];
	}

	/** The state before an edit in a field, to record it once the edit ends (see commit). */
	begin(): string {
		return this.#snapshot();
	}

	/** Records the state from begin as one undo step if it changed since. */
	commit(before: string) {
		if (before === this.#snapshot()) return;
		this.#past = [...this.#past.slice(-(historyMax - 1)), before];
		this.#future = [];
	}

	/** Undoes the last change. */
	undo() {
		const prev = this.#past.at(-1);
		if (!prev) return;
		this.#future = [...this.#future, this.#snapshot()];
		this.#past = this.#past.slice(0, -1);
		this.#restore(prev);
	}

	/** Redoes the last undone change. */
	redo() {
		const next = this.#future.at(-1);
		if (!next) return;
		this.#past = [...this.#past, this.#snapshot()];
		this.#future = this.#future.slice(0, -1);
		this.#restore(next);
	}

	/** Whether there is a change to undo. */
	get canUndo() {
		return this.#past.length > 0;
	}

	/** Whether there is an undone change to redo. */
	get canRedo() {
		return this.#future.length > 0;
	}

	/** Marks the current state as saved. */
	markClean() {
		this.#clean = this.#snapshot();
	}

	/** Whether there are changes since loading or saving. */
	get dirty() {
		return this.#snapshot() !== this.#clean;
	}

	/** The selected elements, in drawing order. */
	get selection(): Item[] {
		return this.items.filter((i) => this.selected.includes(i.id));
	}

	/** Adds an element with defaults for what it leaves out and selects it. */
	add(el: Partial<LayoutElement>): Item {
		this.remember();
		const item: Item = { type: 'text', x: 10, y: 10, w: 0, h: 0, size: 10, bold: false, align: 'L', value: '', repeat: 0, dx: 0, dy: 0, ...el, id: newId() };
		this.items.push(item);
		this.selected = [item.id];
		return item;
	}

	/** Copies the selected elements a little below and to the right, and selects the copies. */
	duplicate() {
		if (!this.selected.length) return;
		this.remember();
		const copies = this.selection.map((i) => ({ ...$state.snapshot(i), id: newId(), x: tenth(num(i.x) + 2), y: tenth(num(i.y) + 2) }));
		this.items.push(...copies);
		this.selected = copies.map((c) => c.id);
	}

	/** Removes the selected elements. */
	remove() {
		if (!this.selected.length) return;
		this.remember();
		this.items = this.items.filter((i) => !this.selected.includes(i.id));
		this.selected = [];
	}

	/** Moves the selected elements to the end of the list, so they print on top. */
	toFront() {
		this.remember();
		this.items = [...this.items.filter((i) => !this.selected.includes(i.id)), ...this.selection];
	}

	/** The form as the API stores it: numbers as numbers, no editor IDs. */
	input(): LayoutInput {
		return {
			name: this.name.trim(),
			width: num(this.width),
			height: num(this.height),
			offsetX: num(this.offsetX),
			offsetY: num(this.offsetY),
			elements: $state.snapshot(this.items).map(({ id: _, ...e }) => ({
				...e,
				x: num(e.x),
				y: num(e.y),
				w: num(e.w),
				h: num(e.h),
				size: num(e.size),
				repeat: num(e.repeat),
				dx: num(e.dx),
				dy: num(e.dy)
			}))
		};
	}
}

/** The text an element shows on the sheet: placeholders as examples (# counted as n), or as names when asked. */
export function fill(value: string | undefined, samples: Record<string, string>, n: number, names = false): string {
	return (value ?? '').replace(/\{\{\s*([A-Za-z0-9_.#]+)\s*\}\}/g, (_, k: string) => {
		const key = k.replace(/#/g, String(n));
		return names ? `{{${k}}}` : (samples[key] ?? `{{${key}}}`);
	});
}
