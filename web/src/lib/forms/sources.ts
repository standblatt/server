import { msg } from '@lingui/core/macro';
import type { MessageDescriptor } from '@lingui/core';
import type { Source } from './editor.svelte';

/** The data sources of print forms with their labels; translate with `$t(...)`. */
export const sourceLabels: Record<Source, MessageDescriptor> = {
	results: msg`Results of a programme`,
	registration: msg`A registration`,
	rank: msg`A rank`
};

/** The data sources in the order the UI offers them. */
export const sourceKeys = Object.keys(sourceLabels) as Source[];
