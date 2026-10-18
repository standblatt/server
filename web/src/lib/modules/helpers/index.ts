import { msg } from '@lingui/core/macro';
import Handshake from '@lucide/svelte/icons/handshake';
import type { ModuleUI } from '../types';

export default {
	label: msg`Helper links`,
	description: msg`Scorers, ammunition desk and results clerk work through a link without an account.`,
	help: 'user-guide/helpers',
	icon: Handshake,
	configPage: () => import('./HelperCodes.svelte')
} satisfies ModuleUI;
