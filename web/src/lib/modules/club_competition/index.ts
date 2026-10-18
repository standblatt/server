import { msg } from '@lingui/core/macro';
import Building from '@lucide/svelte/icons/building-2';
import type { ModuleUI } from '../types';

export default {
	label: msg`Club competition`,
	description: msg`Clubs are ranked by the results of their shooters.`,
	help: 'user-guide/statistics-reports',
	icon: Building,
	settings: () => import('./Settings.svelte'),
	statisticsSections: () => import('./StatisticsSection.svelte'),
	reports: [{ id: 'club_competition', label: msg`Club competition` }]
} satisfies ModuleUI;
