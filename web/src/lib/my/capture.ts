import type { Discipline, TrainingInput, TrainingSession } from '$lib/my.svelte';
import { localInput, uuid } from '$lib/my.svelte';

export type Mode = 'totals' | 'shots' | 'target';
export type Shot = { v: number; x?: number; y?: number };
/** A series being captured: only a total, or every shot (values) with optional positions. */
export type CaptureSeries = { shots: number; total: number | null; values: Shot[] | null };

export type Draft = {
	clientId: string;
	editId?: number;
	disciplineId: number;
	startedAt: string;
	kind: 'training' | 'competition';
	mode: Mode;
	series: CaptureSeries[];
	current: number;
	tags: string[];
	weaponId: number | null;
	rating: number | null;
	notes: string;
	place: string;
	title: string;
};

/** Rounds to one decimal. */
const round1 = (v: number) => Math.round(v * 10) / 10;

/** Empty draft for a discipline with its series layout, started now. */
export function newDraft(d: Discipline, weaponId: number | null = null): Draft {
	return {
		clientId: uuid(),
		disciplineId: d.id,
		startedAt: localInput(new Date()),
		kind: 'training',
		mode: 'totals',
		series: Array.from({ length: d.seriesCount }, () => ({ shots: d.shotsPerSeries, total: null, values: null })),
		current: 0,
		tags: [],
		weaponId,
		rating: null,
		notes: '',
		place: '',
		title: ''
	};
}

/** Draft for editing a saved session; the mode follows the stored data (positions, shots or totals). */
export function fromSession(s: TrainingSession): Draft {
	const series = s.series.map((x) => ({
		shots: x.shots,
		total: x.total,
		values: x.values?.length ? x.values.map((v, i) => ({ v, x: x.x?.[i], y: x.y?.[i] })) : null
	}));
	const positioned = series.some((x) => x.values?.some((v) => v.x !== undefined));
	return {
		clientId: uuid(),
		editId: s.id,
		disciplineId: s.disciplineId,
		startedAt: localInput(new Date(s.startedAt)),
		kind: s.kind,
		mode: positioned ? 'target' : series.some((x) => x.values) ? 'shots' : 'totals',
		series,
		current: series.length - 1,
		tags: [...s.tags],
		weaponId: s.weaponId ?? null,
		rating: s.rating ?? null,
		notes: s.notes,
		place: s.place,
		title: s.title
	};
}

/** The points of a series so far. */
export const seriesTotal = (s: CaptureSeries) => (s.values ? round1(s.values.reduce((a, b) => a + b.v, 0)) : (s.total ?? 0));

/** A series counts once it has a total or at least one shot. */
export const filled = (s: CaptureSeries) => (s.values ? s.values.length > 0 : s.total !== null);

/** Whether a series has its total or all its shots. */
export const complete = (s: CaptureSeries) => (s.values ? s.values.length >= s.shots : s.total !== null);

/** Points and shots of the filled series of a draft. */
export function totals(d: Draft) {
	let points = 0;
	let shots = 0;
	for (const s of d.series) {
		if (!filled(s)) continue;
		points += seriesTotal(s);
		shots += s.values ? s.values.length : s.shots;
	}
	return { points: round1(points), shots };
}

/** API input from a draft: only filled series, positions only when every shot has one, no client ID when editing. */
export function toInput(d: Draft): TrainingInput {
	const series = d.series.filter(filled).map((s) => {
		if (!s.values) return { shots: s.shots, total: s.total! };
		const positioned = s.values.every((v) => v.x !== undefined && v.y !== undefined);
		return {
			shots: s.values.length,
			total: seriesTotal(s),
			values: s.values.map((v) => v.v),
			...(positioned ? { x: s.values.map((v) => v.x!), y: s.values.map((v) => v.y!) } : {})
		};
	});
	return {
		disciplineId: d.disciplineId,
		kind: d.kind,
		startedAt: new Date(d.startedAt).toISOString(),
		title: d.title,
		place: d.place,
		notes: d.notes,
		tags: d.tags,
		...(d.weaponId ? { weaponId: d.weaponId } : {}),
		...(d.rating ? { rating: d.rating } : {}),
		clientId: d.editId ? undefined : d.clientId,
		series
	};
}

/**
 * Adds a key of the number pad to the text typed so far. A value is committed as soon as another digit
 * could not make it valid any more; with decimals the point is added automatically when the whole part is complete.
 */
export function typeKey(buf: string, key: string, max: number, decimal: boolean): { buf: string; commit?: number } {
	if (key === 'back') return { buf: buf.slice(0, -1) };
	if (key === '.') return decimal && !buf.includes('.') ? { buf: (buf || '0') + '.' } : { buf };
	const next = buf + key;
	if (Number(next) > max + 1e-9) return { buf };
	if (decimal && next.includes('.')) return { buf: '', commit: Number(next) };
	if (Number(next) * 10 > max + 1e-9) return decimal ? { buf: next + '.' } : { buf: '', commit: Number(next) };
	return { buf: next };
}

/** Typing a decimal shot value: two digits make 9.7, a leading 10 needs a third digit (10.4). */
export function typeShot(buf: string, key: string): { buf: string; commit?: number } {
	if (key === 'back') return { buf: buf.slice(0, -1) };
	const next = buf + key;
	if (next.length < (next.startsWith('10') ? 3 : 2)) return { buf: next };
	return { buf: '', commit: Number(next.slice(0, -1)) + Number(next.slice(-1)) / 10 };
}
