import type { MessageDescriptor } from '@lingui/core';
import type { Component } from 'svelte';
import type Users from '@lucide/svelte/icons/users';
import type { Event } from '$lib/api';

/** Key of an event module, as the API names it. */
export type ModuleKey = Event['modules'][number];

/** Props every component a module plugs in receives. */
export type EventProps = { tenantId: number; event: Event };

/** A lazily loaded component of a module. */
export type Lazy = () => Promise<{ default: Component<EventProps> }>;

/** The UI of a module: what it is called and where it plugs in. Every extension point is optional. */
export type ModuleUI = {
	label: MessageDescriptor;
	description: MessageDescriptor;
	/** Help page (and anchor) under /help. */
	help: string;
	icon: typeof Users;
	/** Own page under Configuration, at config/<key>. */
	configPage?: Lazy;
	/** Own page in the operation sidebar, at events/<id>/<key>. */
	operationPage?: Lazy;
	/** Badge next to the operation page in the sidebar. */
	badge?: Lazy;
	/** Small settings shown on the modules page under the module's switch. */
	settings?: Lazy;
	/** Sections on the rankings page, after the programme rankings. */
	rankingsSections?: Lazy;
	/** Sections on the statistics page. */
	statisticsSections?: Lazy;
	/** Reports listed on the statistics page. */
	reports?: { id: string; label: MessageDescriptor }[];
};
