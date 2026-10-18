import { msg } from '@lingui/core/macro';
import Ticket from '@lucide/svelte/icons/ticket';
import type { ModuleUI } from '../types';

export default {
	label: msg`Control strip`,
	description: msg`Printed at registration: the passes the shooter may shoot.`,
	help: 'user-guide/registrations',
	icon: Ticket
} satisfies ModuleUI;
