import { msg } from '@lingui/core/macro';
import type { MessageDescriptor } from '@lingui/core';

/** Display name of a tenant role; translate with `$t(...)`. */
export const roleLabel = (role: string): MessageDescriptor => (role === 'admin' ? msg`Admin` : msg`Staff`);

/** Display name of a helper link role; translate with `$t(...)`. */
export const helperRoleLabel = (role: string): MessageDescriptor =>
	role === 'scorer' ? msg`Scorer` : role === 'ammunition' ? msg`Ammunition desk` : msg`Results clerk`;

/** Display name of an age category; translate with `$t(...)`. */
export function ageCategoryLabel(c: string): MessageDescriptor {
	switch (c) {
		case 'elite':
			return msg`Elite`;
		case 'senior':
			return msg`Senior`;
		case 'veteran':
			return msg`Veteran`;
		case 'senior_veteran':
			return msg`Senior veteran`;
		default:
			return { id: c, message: c.toUpperCase() };
	}
}

/** Display name of a registration status. */
export function registrationStatus(s: string): MessageDescriptor {
	return s === 'completed' ? msg`completed` : s === 'shooting' ? msg`shooting` : msg`registered`;
}

/** Display name of a print job status. */
export function jobStatus(s: string): MessageDescriptor {
	switch (s) {
		case 'queued':
			return msg`queued`;
		case 'printing':
			return msg`printing`;
		case 'done':
			return msg`done`;
		case 'failed':
			return msg`failed`;
		default:
			return { id: s, message: s };
	}
}

/** "live" or "offline" for a realtime connection. */
export const liveLabel = (connected?: boolean): MessageDescriptor => (connected ? msg`live` : msg`offline`);
