import { createQuery } from '@tanstack/svelte-query';
import { api, call } from '$lib/api';
import type { Event } from '$lib/api';
import reservations from './reservations';
import booking from './booking';
import lanesheet from './lane_sheet';
import controlstrip from './control_strip';
import ammo from './ammo';
import attempts from './attempts';
import awards from './awards';
import specialprizes from './special_prizes';
import clubcompetition from './club_competition';
import classifications from './classifications';
import kiosk from './kiosk';
import helpers from './helpers';
import type { ModuleKey, ModuleUI } from './types';

export type { EventProps, Lazy, ModuleKey, ModuleUI } from './types';

/** The UI of every event module; a module the API knows but this list lacks fails the type check. */
export const moduleUI: Record<ModuleKey, ModuleUI> = {
	reservations: reservations,
	booking: booking,
	lane_sheet: lanesheet,
	control_strip: controlstrip,
	ammo: ammo,
	attempts: attempts,
	awards: awards,
	special_prizes: specialprizes,
	club_competition: clubcompetition,
	classifications: classifications,
	kiosk: kiosk,
	helpers: helpers
};

/** Whether a module is on for an event. */
export const has = (event: Pick<Event, 'modules'> | undefined, key: ModuleKey) => !!event?.modules.includes(key);

/** The modules that are on for an event, in the order of the API, with their UI. */
export const enabled = (event: Pick<Event, 'modules'> | undefined) => (event?.modules ?? []).map((key) => ({ key, ...moduleUI[key] }));

/** Query of the modules the server knows, with what they require (never changes while the app runs). */
export function modulesQuery() {
	return createQuery(() => ({
		queryKey: ['modules'],
		queryFn: () => call(api.GET('/api/modules')),
		staleTime: Infinity
	}));
}
