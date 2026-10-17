import type { Schemas } from './api';

type Queue = unknown[][];
declare global {
	interface Window {
		_paq?: Queue;
	}
}

let started = false;
let last = '';

/** Kiosk displays run for hours and would only inflate the numbers; helper and booking paths carry secrets. */
function tracked(path: string) {
	return !['/kiosk/', '/h/', '/b/'].some((p) => path.startsWith(p));
}

/**
 * Loads the Matomo tracker without cookies. Page views are reported by `trackPage`
 * with IDs replaced and without the query string, so no record IDs or tokens leave the app.
 */
export function startMatomo(m: Schemas['Matomo']) {
	if (started) return;
	started = true;
	const paq = (window._paq = window._paq || []);
	paq.push(['disableCookies']);
	paq.push(['enableLinkTracking']);
	paq.push(['setTrackerUrl', m.url + 'matomo.php']);
	paq.push(['setSiteId', String(m.siteId)]);
	const script = document.createElement('script');
	script.async = true;
	script.src = m.url + 'matomo.js';
	document.head.append(script);
	trackPage();
}

/** Reports the current path to Matomo (numeric segments as :id) once per path change. */
export function trackPage() {
	if (!started || !tracked(location.pathname) || location.pathname === last) return;
	last = location.pathname;
	const path = location.pathname.replace(/\/\d+(?=\/|$)/g, '/:id');
	window._paq!.push(['setCustomUrl', location.origin + path]);
	window._paq!.push(['setDocumentTitle', path]);
	window._paq!.push(['trackPageView']);
}
