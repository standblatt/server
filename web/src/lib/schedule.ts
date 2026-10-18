import type { Schedule } from './api';

/** ISO weekday of a YYYY-MM-DD date: 1 = Monday … 7 = Sunday. */
const isoWeekday = (iso: string) => ((new Date(`${iso}T00:00:00Z`).getUTCDay() + 6) % 7) + 1;

/** The shooting days of a schedule plus the main day, sorted; mirrors domain.Schedule.Expand on the server. */
export function shootingDays(s: Schedule, mainDay: string): string[] {
	const days = new Set<string>([mainDay, ...(s.days ?? []).map((d) => d.date).filter((d) => d)]);
	if (s.from && s.to && s.from <= s.to) {
		const weekdays = new Set((s.weekly ?? []).map((w) => w.weekday));
		for (let d = new Date(`${s.from}T00:00:00Z`); d.toISOString().slice(0, 10) <= s.to && days.size < 1000; d.setUTCDate(d.getUTCDate() + 1)) {
			const iso = d.toISOString().slice(0, 10);
			if (weekdays.size === 0 || weekdays.has(isoWeekday(iso))) days.add(iso);
		}
	}
	return [...days].sort();
}

/** Localized weekday name of an ISO weekday (1 = Monday). */
export function weekdayName(weekday: number, locale: string): string {
	return new Date(Date.UTC(2024, 0, weekday)).toLocaleDateString(locale, { weekday: 'long', timeZone: 'UTC' });
}
