import { msg } from '@lingui/core/macro';
import type { MessageDescriptor } from '@lingui/core';
import { loadPage, type HelpPage } from './markdown';

/** The help sidebar: groups of page slugs, in reading order. */
export const nav: { label: MessageDescriptor; pages: string[] }[] = [
	{ label: msg`Start`, pages: ['', 'getting-started/concepts', 'getting-started/first-event', 'changelog'] },
	{ label: msg`My training`, pages: ['my-training', 'my-training/logging', 'my-training/progress', 'my-training/account'] },
	{
		label: msg`At the event`,
		pages: [
			'user-guide/dashboard',
			'user-guide/registrations',
			'user-guide/results',
			'user-guide/results-entry',
			'user-guide/rankings-kiosk',
			'user-guide/reservations',
			'user-guide/statistics-reports',
			'user-guide/helpers'
		]
	},
	{
		label: msg`Setup`,
		pages: [
			'user-guide/events',
			'user-guide/online-booking',
			'user-guide/programmes',
			'user-guide/shooters-clubs',
			'user-guide/ranges',
			'user-guide/printing',
			'user-guide/print-forms',
			'user-guide/users',
			'user-guide/login-language'
		]
	},
	{ label: msg`Reference`, pages: ['reference/rules', 'reference/csv-import', 'reference/programme-file', 'reference/codes'] }
];

const cache = new Map<string, Promise<HelpPage[]>>();

/** Every page of the sidebar in a language, loaded once, for titles and the search. */
export function allPages(lang: string): Promise<HelpPage[]> {
	let all = cache.get(lang);
	if (!all) {
		all = Promise.all(nav.flatMap((g) => g.pages).map((slug) => loadPage(lang, slug))).then((pages) => pages.filter((p) => !!p));
		cache.set(lang, all);
	}
	return all;
}
