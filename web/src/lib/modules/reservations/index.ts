import { msg } from '@lingui/core/macro';
import CalendarClock from '@lucide/svelte/icons/calendar-clock';
import type { ModuleUI } from '../types';

export default {
	label: msg`Reservations`,
	description: msg`Lanes reserved for groups and clubs, shown on the dashboard.`,
	help: 'user-guide/reservations',
	icon: CalendarClock,
	operationPage: () => import('./Page.svelte'),
	badge: () => import('./Badge.svelte')
} satisfies ModuleUI;
