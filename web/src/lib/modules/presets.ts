import { msg } from '@lingui/core/macro';
import type { MessageDescriptor } from '@lingui/core';
import type { ModuleKey } from './types';

/** A starting point for a new event: only the modules it switches on, everything can be changed later. */
export type Preset = { id: string; label: MessageDescriptor; description: MessageDescriptor; modules: ModuleKey[] };

/** The presets offered when creating an event, the default first. */
export const presets: Preset[] = [
	{
		id: 'festival',
		label: msg`Shooting festival`,
		description: msg`Desk with lanes, online booking, awards, special prizes, club competition, kiosk and helpers.`,
		modules: ['lanes', 'reservations', 'lane_sheet', 'control_strip', 'awards', 'special_prizes', 'club_competition', 'kiosk', 'helpers']
	},
	{
		id: 'annual',
		label: msg`Annual programme or training`,
		description: msg`Results over the season with several attempts and classifications, no desk.`,
		modules: ['attempts', 'classifications']
	},
	{
		id: 'field',
		label: msg`Feldschiessen`,
		description: msg`Lanes, awards and club competition, kiosk displays and helper links.`,
		modules: ['lanes', 'awards', 'club_competition', 'kiosk', 'helpers']
	},
	{
		id: 'championship',
		label: msg`Club championship`,
		description: msg`Several attempts ranked in classifications, shown on kiosk displays.`,
		modules: ['attempts', 'classifications', 'kiosk']
	},
	{ id: 'custom', label: msg`Custom selection`, description: msg`Start with nothing and switch on modules as needed.`, modules: [] }
];
