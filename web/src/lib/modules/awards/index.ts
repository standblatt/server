import { msg } from '@lingui/core/macro';
import Award from '@lucide/svelte/icons/award';
import type { ModuleUI } from '../types';

export default {
	label: msg`Awards`,
	description: msg`Wreath cards, prizes in kind and cash prizes by points per age category.`,
	help: 'user-guide/programmes',
	icon: Award,
	settings: () => import('./Settings.svelte'),
	statisticsSections: () => import('./StatisticsSection.svelte'),
	reports: [
		{ id: 'awards', label: msg`Awards` },
		{ id: 'awards_open', label: msg`Awards to hand out` }
	]
} satisfies ModuleUI;
