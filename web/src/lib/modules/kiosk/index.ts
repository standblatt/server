import { msg } from '@lingui/core/macro';
import Monitor from '@lucide/svelte/icons/monitor';
import type { ModuleUI } from '../types';

export default {
	label: msg`Kiosk displays`,
	description: msg`Live rankings on screens at the range through a link.`,
	help: 'user-guide/rankings-kiosk',
	icon: Monitor,
	configPage: () => import('./ConfigPage.svelte'),
	programmeSections: [{ after: 'basics', load: () => import('./ProgrammeSection.svelte') }]
} satisfies ModuleUI;
