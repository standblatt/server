import { msg } from '@lingui/core/macro';
import Globe from '@lucide/svelte/icons/globe';
import type { ModuleUI } from '../types';

export default {
	label: msg`Online booking`,
	description: msg`Clubs ask for lanes through a public link; confirmed requests become reservations.`,
	help: 'user-guide/online-booking',
	icon: Globe,
	configPage: () => import('./BookingSettings.svelte'),
	reports: [{ id: 'bookings', label: msg`Lane bookings` }]
} satisfies ModuleUI;
