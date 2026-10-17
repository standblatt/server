import type { AgeCategory } from './api';

export const ageCategories: AgeCategory[] = ['u10', 'u13', 'u15', 'u17', 'u19', 'u21', 'elite', 'senior', 'veteran', 'senior_veteran'];

export const ageCategoryShort: Record<AgeCategory, string> = {
	u10: 'U10',
	u13: 'U13',
	u15: 'U15',
	u17: 'U17',
	u19: 'U19',
	u21: 'U21',
	elite: 'E',
	senior: 'S',
	veteran: 'V',
	senior_veteran: 'SV'
};

/** Age category from the birth year and the event year (SSV year rule, display only; the server decides). */
export function ageCategoryFor(birthday: string, eventDate: string): AgeCategory {
	const born = Number(birthday.slice(0, 4));
	const year = Number(eventDate.slice(0, 4));
	if (!born || !year) return 'elite';
	const age = year - born;
	if (age <= 9) return 'u10';
	if (age <= 12) return 'u13';
	if (age <= 14) return 'u15';
	if (age <= 16) return 'u17';
	if (age <= 18) return 'u19';
	if (age <= 20) return 'u21';
	if (age <= 45) return 'elite';
	if (age <= 59) return 'senior';
	if (age <= 69) return 'veteran';
	return 'senior_veteran';
}
