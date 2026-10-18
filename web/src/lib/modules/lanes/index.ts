import { msg } from '@lingui/core/macro';
import Grid from '@lucide/svelte/icons/grid-3x3';
import type { ModuleUI } from '../types';

export default {
	label: msg`Lane assignment`,
	description: msg`Shooters are registered on lanes: lane grid on the dashboard, lane status on kiosk displays and scorer links.`,
	help: 'user-guide/events#modules',
	icon: Grid
} satisfies ModuleUI;
