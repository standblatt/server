/** Formats an ISO date as dd.mm.yyyy. */
export function formatDate(iso: string): string {
	const [y, m, d] = iso.split('-');
	return `${d}.${m}.${y}`;
}

/** First and last of the sorted days: "08.01.2027" or "20.12.2026–10.01.2027". */
export function formatDays(days: string[], fallback = ''): string {
	if (days.length === 0) return fallback ? formatDate(fallback) : '';
	const first = formatDate(days[0]);
	const last = formatDate(days[days.length - 1]);
	return first === last ? first : `${first}–${last}`;
}

/** Names of the ranges, comma separated. */
export const rangeNames = (ranges: { name: string }[]) => ranges.map((r) => r.name).join(', ');

/** Programme name with the variant in parentheses: "Championship (Sport)". */
export const programmeTitle = (p: { name: string; variant?: string }) => (p.variant ? `${p.name} (${p.variant})` : p.name);

/** Pads a number to two digits. */
export const pad2 = (n: number) => String(n).padStart(2, '0');

/** Local date and time for an <input type="datetime-local">. */
export function localInput(d: Date): string {
	return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())}T${pad2(d.getHours())}:${pad2(d.getMinutes())}`;
}

/** Prefixes a lane text with the range name when the event has several ranges: "300m · Lane 4". */
export function withRange(ranges: { id: number; name: string }[], rangeId: number | undefined, laneText: string): string {
	if (ranges.length < 2) return laneText;
	const name = ranges.find((r) => r.id === rangeId)?.name ?? '';
	return name ? `${name} · ${laneText}` : laneText;
}

/** Writes lane numbers like the server: runs of three or more as 1-4, others separated by commas. */
export function formatLanes(lanes: number[]): string {
	const parts: string[] = [];
	for (let i = 0; i < lanes.length; ) {
		let j = i;
		while (j + 1 < lanes.length && lanes[j + 1] === lanes[j] + 1) j++;
		parts.push(j - i >= 2 ? `${lanes[i]}-${lanes[j]}` : lanes.slice(i, j + 1).join(', '));
		i = j + 1;
	}
	return parts.join(', ');
}

/** Reads lane lists like "1, 5, 7-10" into sorted, unique lane numbers; NaN marks an invalid part. */
export function parseLanes(s: string): number[] {
	const out = new Set<number>();
	for (const part of s.split(',').map((x) => x.trim()).filter(Boolean)) {
		const [a, b = a] = part.split('-').map((x) => Number(x.trim()));
		if (!Number.isInteger(a) || !Number.isInteger(b) || b < a) out.add(NaN);
		else for (let n = a; n <= b; n++) out.add(n);
	}
	return [...out].sort((x, y) => x - y);
}

/** Last page index (0-based) of a list of total items in pages of size, at least 0. */
export const lastPage = (total: number, size: number) => Math.max(0, Math.ceil(total / size) - 1);

/** The items of one page; a page past the end shows the last one. */
export const pageOf = <T>(list: T[], page: number, size: number) => {
	const p = Math.min(page, lastPage(list.length, size));
	return list.slice(p * size, (p + 1) * size);
};
