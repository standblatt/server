import { msg } from '@lingui/core/macro';
import Repeat from '@lucide/svelte/icons/repeat';
import type { ModuleUI } from '../types';

export default {
	label: msg`Attempts`,
	description: msg`Programmes can be shot several times, e.g. further attempts or an annual programme.`,
	help: 'user-guide/programmes#attempts',
	icon: Repeat
} satisfies ModuleUI;
