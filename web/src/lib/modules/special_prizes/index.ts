import { msg } from '@lingui/core/macro';
import Trophy from '@lucide/svelte/icons/trophy';
import type { ModuleUI } from '../types';

export default {
	label: msg`Special prizes`,
	description: msg`The best ranks of every age category get a special prize.`,
	help: 'user-guide/programmes#special-prizes',
	icon: Trophy,
	reports: [{ id: 'special_prizes', label: msg`Special prizes` }]
} satisfies ModuleUI;
