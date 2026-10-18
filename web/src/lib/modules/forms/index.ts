import { msg } from '@lingui/core/macro';
import ScrollText from '@lucide/svelte/icons/scroll-text';
import type { ClubModuleUI } from '../types';

export default {
	label: msg`Print forms`,
	description: msg`Results printed onto pre-printed sheets, such as the army's lane sheets.`,
	help: 'user-guide/print-forms',
	icon: ScrollText,
	clubPage: { path: 'forms', label: msg`Print forms`, permission: 'forms.manage' },
	programmeSections: [{ after: 'usage', load: () => import('./ProgrammeSection.svelte') }]
} satisfies ClubModuleUI;
