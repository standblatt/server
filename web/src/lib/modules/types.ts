import type { MessageDescriptor } from '@lingui/core';
import type { Component } from 'svelte';
import type Users from '@lucide/svelte/icons/users';
import type { Event, Programme } from '$lib/api';
import type { ProgrammeDraft } from '$lib/programme/draft.svelte';

/** Key of an event module, as the API names it. */
export type ModuleKey = Event['modules'][number];

/** Props every component a module plugs in receives. */
export type EventProps = { tenantId: number; event: Event };

/** A lazily loaded component of a module. */
export type Lazy<P extends Record<string, any> = EventProps> = () => Promise<{ default: Component<P> }>;

/** Props of a section of the programme page: the draft being edited and, for a saved programme, the programme. */
export type ProgrammeProps = EventProps & { draft: ProgrammeDraft; locked: boolean; programme?: Programme };

/** Where a section goes on the programme page: after a core section. */
export type ProgrammeSection = { after: 'basics' | 'series' | 'prices' | 'usage'; load: Lazy<ProgrammeProps> };

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
	/** Sections on the programme page; `usage` ones save at once and only show for saved programmes. */
	programmeSections?: ProgrammeSection[];
	/** Reports listed on the statistics page. */
	reports?: { id: string; label: MessageDescriptor }[];
};
