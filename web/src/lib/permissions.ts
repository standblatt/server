import { msg } from '@lingui/core/macro';
import type { MessageDescriptor } from '@lingui/core';
import type { Permission } from './api';
import type { ClubModuleKey } from './modules/types';

/** How the role matrix shows a permission; module hides it while that club module is off. */
type PermissionUI = { group: MessageDescriptor; label: MessageDescriptor; description: MessageDescriptor; module?: ClubModuleKey };

const events = msg`Events`;
const master = msg`Master data`;
const club = msg`Club`;

/** Every permission with its texts, in the order of the role matrix; svelte-check catches a missing one. */
export const permissionUI: Record<Permission, PermissionUI> = {
	'events.view': { group: events, label: msg`See events and rankings`, description: msg`Events, programmes, rankings and classifications.` },
	'events.operate': { group: events, label: msg`Run events`, description: msg`Registrations, results, awards, reservations, documents and printing.` },
	'events.configure': {
		group: events,
		label: msg`Configure events`,
		description: msg`Create, configure, copy, close and delete events with programmes, classifications, kiosk and helper links and online booking.`
	},
	'reports.view': { group: events, label: msg`Statistics and reports`, description: msg`Statistics, reports with amounts and ranking CSV files.` },
	'shooters.view': { group: master, label: msg`See shooters and clubs`, description: msg`The personal data of the master data.` },
	'shooters.edit': { group: master, label: msg`Edit shooters and clubs`, description: msg`Add and change shooters and clubs.` },
	'shooters.manage': { group: master, label: msg`Manage master data`, description: msg`Import shooters, rename, merge and delete clubs, delete shooters.` },
	'ranges.manage': { group: club, label: msg`Ranges and printers`, description: msg`Ranges, print agents and printer settings.` },
	'forms.manage': { group: club, label: msg`Print forms`, description: msg`Print forms and images.`, module: 'forms' },
	'club.manage': { group: club, label: msg`Users and club`, description: msg`Users, roles, invitations and club modules.` }
};

/** The permission keys in matrix order. */
export const permissionKeys = Object.keys(permissionUI) as Permission[];
