import type { Helper } from './api';

/** What the event desk (dashboard) lets the viewer do; staff get everything, helper links a subset per role. */
export type DeskRights = {
	search: boolean;
	register: boolean;
	assign: boolean;
	results: boolean;
	finish: boolean;
	ammo: boolean;
	print: boolean;
	/** Print the results onto the print form of their programme. */
	forms: boolean;
	waiting: boolean;
	completed: boolean;
	/** Mark the awards of a registration as handed out. */
	awards: boolean;
	/** Take an ammunition or awards handout back (club members only). */
	revoke: boolean;
};

export const staffRights: DeskRights = { search: true, register: true, assign: true, results: true, finish: true, ammo: true, print: true, forms: true, waiting: true, completed: true, awards: true, revoke: true };

const none: DeskRights = { search: false, register: false, assign: false, results: false, finish: false, ammo: false, print: false, forms: false, waiting: false, completed: false, awards: false, revoke: false };

/** Desk rights of a helper link by role: scorer, ammunition desk or results desk. */
export function helperRights(role: Helper['role']): DeskRights {
	switch (role) {
		case 'scorer':
			return { ...none, results: true, finish: true, forms: true };
		case 'ammunition':
			return { ...none, search: true, assign: true, ammo: true, print: true, waiting: true };
		default:
			return { ...none, search: true, results: true, completed: true, awards: true, forms: true };
	}
}
