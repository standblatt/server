import { msg } from '@lingui/core/macro';
import CalendarClock from '@lucide/svelte/icons/calendar-clock';
import type { ModuleUI } from '../types';

export default {
	label: msg`Reservations`,
	description: msg`Lanes reserved for groups and clubs, shown on the dashboard; clubs can ask for lanes online.`,
	help: 'user-guide/reservations',
	icon: CalendarClock,
	operationPage: () => import('./Page.svelte'),
	badge: () => import('./Badge.svelte'),
	configPage: () => import('./booking/BookingSettings.svelte'),
	reports: [{ id: 'bookings', label: msg`Lane bookings` }]
} satisfies ModuleUI;
