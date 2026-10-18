import { msg } from '@lingui/core/macro';
import Package from '@lucide/svelte/icons/package';
import type { ModuleUI } from '../types';

export default {
	label: msg`Ammunition desk`,
	description: msg`The organizer hands out the ammunition: an Ammo button on the dashboard and a helper role.`,
	help: 'user-guide/helpers#ammunition',
	icon: Package
} satisfies ModuleUI;
