import { msg } from '@lingui/core/macro';
import Layers from '@lucide/svelte/icons/layers';
import type { ModuleUI } from '../types';

export default {
	label: msg`Classifications`,
	description: msg`Rankings over several programmes or attempts, with classes and tie-breaks.`,
	help: 'user-guide/classifications',
	icon: Layers,
	configPage: () => import('./Classifications.svelte'),
	rankingsSections: () => import('./RankingsSection.svelte'),
	reports: [{ id: 'classifications', label: msg`Classifications` }]
} satisfies ModuleUI;
