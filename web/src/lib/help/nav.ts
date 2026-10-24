import { msg } from '@lingui/core/macro';
import type { MessageDescriptor } from '@lingui/core';
import { loadPage, type HelpPage } from './markdown';

/** The help sidebar: groups of page slugs, in reading order. */
export const nav: { label: MessageDescriptor; pages: string[] }[] = [
	{ label: msg`Start`, pages: ['', 'getting-started/first-event', 'getting-started/concepts'] },
	{ label: msg`My training`, pages: ['my-training', 'my-training/first-session', 'my-training/logging', 'my-training/progress', 'my-training/account'] },
	{
		label: msg`How-to guides`,
		pages: [
			'how-to/next-year',
			'how-to/annual-programme',
			'how-to/section-championship',
			'how-to/feldschiessen',
			'how-to/import-results',
			'how-to/kiosk',
			'how-to/connect-printer',
			'how-to/army-sheet',
			'how-to/certificates',
			'how-to/hand-out-awards',
			'how-to/two-factor'
		]
	},
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
		label: msg`Configuration`,
		pages: [
			'user-guide/events',
			'user-guide/online-booking',
			'user-guide/programmes',
			'user-guide/classifications',
			'user-guide/shooters-clubs',
			'user-guide/ranges',
			'user-guide/printing',
			'user-guide/print-forms',
			'user-guide/users',
			'user-guide/login-language'
		]
	},
	{
		label: msg`Background`,
		pages: ['explanation/modules', 'explanation/rankings', 'explanation/closing', 'explanation/printing', 'explanation/data']
	},
	{
		label: msg`Reference`,
		pages: [
			'reference/rules',
			'reference/modules',
			'reference/permissions',
			'reference/documents',
			'reference/shortcuts',
			'reference/links',
			'reference/form-placeholders',
			'reference/csv-import',
			'reference/programme-file',
			'reference/codes',
			'reference/glossary',
			'changelog'
		]
	}
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
