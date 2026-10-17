import { createQuery } from '@tanstack/svelte-query';
import { api, call } from './api';

export const keys = {
	ranges: (t: number) => ['tenants', t, 'ranges'],
	events: (t: number) => ['tenants', t, 'events'],
	event: (t: number, e: number) => ['tenants', t, 'events', e],
	users: (t: number) => ['tenants', t, 'users'],
	clubs: (t: number) => ['tenants', t, 'clubs'],
	shooters: (t: number) => ['tenants', t, 'shooters'],
	shooter: (t: number, s: number) => ['tenants', t, 'shooters', s],
	programmes: (t: number, e: number) => ['tenants', t, 'events', e, 'programmes'],
	registrations: (t: number, e: number) => ['tenants', t, 'events', e, 'registrations'],
	lanes: (t: number, e: number) => ['tenants', t, 'events', e, 'lanes'],
	stats: (t: number, e: number) => ['tenants', t, 'events', e, 'stats'],
	rankings: (t: number, e: number) => ['tenants', t, 'events', e, 'rankings'],
	kioskTokens: (t: number, e: number) => ['tenants', t, 'events', e, 'kiosk-tokens'],
	helperCodes: (t: number, e: number) => ['tenants', t, 'events', e, 'helper-codes'],
	reports: (t: number, e: number) => ['tenants', t, 'events', e, 'reports'],
	booking: (t: number, e: number) => ['tenants', t, 'events', e, 'booking'],
	reservations: (t: number, e: number) => ['tenants', t, 'events', e, 'reservations']
};

/** Path parameters for tenant-scoped API calls. */
const path = (tenantId: number) => ({ params: { path: { tenantId } } });

/** Query of the tenant's ranges. */
export function rangesQuery(tenantId: () => number) {
	return createQuery(() => ({
		queryKey: keys.ranges(tenantId()),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/ranges', path(tenantId())))
	}));
}

/** Query of the tenant's events. */
export function eventsQuery(tenantId: () => number) {
	return createQuery(() => ({
		queryKey: keys.events(tenantId()),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events', path(tenantId())))
	}));
}

/** Query of the tenant's users. */
export function usersQuery(tenantId: () => number) {
	return createQuery(() => ({
		queryKey: keys.users(tenantId()),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/users', path(tenantId())))
	}));
}

/** Query of the tenant's clubs. */
export function clubsQuery(tenantId: () => number) {
	return createQuery(() => ({
		queryKey: keys.clubs(tenantId()),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/clubs', path(tenantId())))
	}));
}

/** Query of one event. */
export function eventQuery(tenantId: () => number, eventId: () => number) {
	return createQuery(() => ({
		queryKey: keys.event(tenantId(), eventId()),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}', { params: { path: { tenantId: tenantId(), eventId: eventId() } } }))
	}));
}

/** Query of the programmes of an event. */
export function programmesQuery(tenantId: () => number, eventId: () => number) {
	return createQuery(() => ({
		queryKey: keys.programmes(tenantId(), eventId()),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/programmes', { params: { path: { tenantId: tenantId(), eventId: eventId() } } }))
	}));
}

/** Query of the lane states of an event, per range. */
export function lanesQuery(tenantId: () => number, eventId: () => number) {
	return createQuery(() => ({
		queryKey: keys.lanes(tenantId(), eventId()),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/lanes', { params: { path: { tenantId: tenantId(), eventId: eventId() } } }))
	}));
}

/** Query of all registrations of an event. */
export function registrationsQuery(tenantId: () => number, eventId: () => number) {
	return createQuery(() => ({
		queryKey: keys.registrations(tenantId(), eventId()),
		queryFn: () =>
			call(api.GET('/api/tenants/{tenantId}/events/{eventId}/registrations', { params: { path: { tenantId: tenantId(), eventId: eventId() }, query: { all: true } } }))
	}));
}

/** Query of the booking requests of an event. */
export function bookingRequestsQuery(tenantId: () => number, eventId: () => number) {
	return createQuery(() => ({
		queryKey: [...keys.booking(tenantId(), eventId()), 'requests'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/booking-requests', { params: { path: { tenantId: tenantId(), eventId: eventId() } } }))
	}));
}
