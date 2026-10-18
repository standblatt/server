import { createQuery } from '@tanstack/svelte-query';
import { api, call } from '$lib/api';
import type { Event } from '$lib/api';
import lanes from './lanes';
import reservations from './reservations';
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
import forms from './forms';
import type { ClubModuleKey, ClubModuleUI, ModuleKey, ModuleUI } from './types';

export type { ClubModuleKey, ClubModuleUI, EventProps, Lazy, ModuleKey, ModuleUI, ProgrammeProps, ProgrammeSection } from './types';

/** The UI of every event module; a module the API knows but this list lacks fails the type check. */
export const moduleUI: Record<ModuleKey, ModuleUI> = {
	lanes: lanes,
	reservations: reservations,
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

/** The UI of every club module. */
export const clubModuleUI: Record<ClubModuleKey, ClubModuleUI> = {
	forms
};

/** Whether a club module is on for a club (its membership in Me). */
export const hasClub = (club: { modules: string[] } | undefined, key: ClubModuleKey) => !!club?.modules.includes(key);

/** The club modules that are on, with their UI. */
export const enabledClub = (club: { modules: ClubModuleKey[] } | undefined) => (club?.modules ?? []).map((key) => ({ key, ...clubModuleUI[key] }));

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
