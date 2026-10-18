import { msg } from '@lingui/core/macro';
import FileText from '@lucide/svelte/icons/file-text';
import type { ModuleUI } from '../types';

export default {
	label: msg`Lane sheet`,
	description: msg`Printed at registration: shooter, lane, programmes and barcodes.`,
	help: 'user-guide/registrations',
	icon: FileText
} satisfies ModuleUI;
