import type { BookingShooter } from './api';
import { intlLocale } from './locale.svelte';

/** Formats the time of an ISO timestamp as HH:MM in the event's time zone. */
export const slotTime = (iso: string, timeZone: string) =>
	new Date(iso).toLocaleTimeString(intlLocale(), { hour: '2-digit', minute: '2-digit', timeZone });

/** Formats the day of an ISO timestamp with weekday in the event's time zone. */
export const slotDay = (iso: string, timeZone: string) =>
	new Date(iso).toLocaleDateString(intlLocale(), { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric', timeZone });

/** Formats a deadline (the moment after its last day) as the last day it still allows. */
export const deadlineDay = (iso: string, timeZone: string) =>
	new Date(new Date(iso).getTime() - 1).toLocaleDateString(intlLocale(), { day: 'numeric', month: 'long', year: 'numeric', timeZone });

/** Groups items with a start timestamp by their day in the event's time zone, in order. */
export function byDay<T extends { start: string }>(items: T[], timeZone: string): [string, T[]][] {
	const days = new Map<string, T[]>();
	for (const it of [...items].sort((a, b) => a.start.localeCompare(b.start))) {
		const day = slotDay(it.start, timeZone);
		days.set(day, [...(days.get(day) ?? []), it]);
	}
	return [...days.entries()];
}

/** Key of a slot on a range, for the chosen lanes. */
export const slotKey = (rangeId: number, start: string) => `${rangeId}|${new Date(start).toISOString()}`;

/** Turns DD.MM.YYYY, D.M.YYYY or YYYY-MM-DD into YYYY-MM-DD; empty when it is none of them. */
export function isoDate(text: string): string {
	const s = text.trim();
	if (/^\d{4}-\d{2}-\d{2}$/.test(s)) return s;
	const m = /^(\d{1,2})\.(\d{1,2})\.(\d{4})$/.exec(s);
	return m ? `${m[3]}-${m[2].padStart(2, '0')}-${m[1].padStart(2, '0')}` : '';
}

/**
 * Reads a pasted list (from Excel: tab separated, or ; or ,) with one shooter per line in the order
 * last name, first name, birthday, licence. A header line or lines without a date are skipped.
 */
export function parseShooters(text: string): BookingShooter[] {
	const out: BookingShooter[] = [];
	for (const line of text.split(/\r?\n/)) {
		const cells = line.split(/\t|;|,/).map((c) => c.trim());
		const birthday = isoDate(cells[2] ?? '');
		if (!cells[0] || !cells[1] || !birthday) continue;
		out.push({ lastname: cells[0], firstname: cells[1], birthday, licenseNumber: cells[3] ?? '' });
	}
	return out;
}
