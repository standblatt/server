import { redirect } from '@sveltejs/kit';

const targets: Record<string, string> = {
	programmes: 'config/programmes',
	event: 'config',
	classifications: 'config/classifications',
	kiosk: 'config/kiosk',
	helpers: 'config/helpers',
	booking: 'config/booking',
	more: 'config/close'
};

/** The old setup page with tabs: sends its links to the pages of the configuration. */
export function load({ params, url }) {
	redirect(307, `/t/${params.tenantId}/events/${params.eventId}/${targets[url.searchParams.get('tab') ?? ''] ?? 'config/programmes'}`);
}
