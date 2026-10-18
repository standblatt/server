<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { Event, Programme, ProgrammeDefinition } from './api';
	import { Button } from './components/ui/button';
	import FormError from './FormError.svelte';
	import { membershipOf, meQuery } from './me';
	import { enabled, enabledClub, type ProgrammeSection } from './modules';
	import ModuleSlot from './modules/ModuleSlot.svelte';
	import Basics from './programme/Basics.svelte';
	import type { ProgrammeDraft } from './programme/draft.svelte';
	import Prices from './programme/Prices.svelte';
	import Section from './programme/Section.svelte';
	import Series from './programme/Series.svelte';
	import Usage from './programme/Usage.svelte';

	let {
		tenantId,
		event,
		draft,
		programme,
		others = [],
		locked = false,
		submitLabel,
		onsubmit
	}: {
		tenantId: number;
		event: Event;
		draft: ProgrammeDraft;
		/** The saved programme, missing for a new one. */
		programme?: Programme;
		others?: { name: string; variant?: string }[];
		/** Shooters are registered: series are fixed, age categories and awards can only be added. */
		locked?: boolean;
		submitLabel: string;
		onsubmit: (d: ProgrammeDefinition) => Promise<void>;
	} = $props();

	let error = $state('');
	let busy = $state(false);
	const me = meQuery();
	const club = $derived(membershipOf(me.data, tenantId));
	/** Sections of the switched-on event and club modules that go after a core section. */
	const after = (anchor: ProgrammeSection['after']) =>
		[...enabled(event), ...enabledClub(club)].flatMap((m) => (m.programmeSections ?? []).filter((s) => s.after === anchor));

	/** Saves the definition through onsubmit. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		busy = true;
		error = '';
		try {
			await onsubmit(draft.definition());
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

{#snippet modules(anchor: ProgrammeSection['after'])}
	{#each after(anchor) as s, i (i)}<ModuleSlot load={s.load} props={{ tenantId, event, draft, locked, programme }} />{/each}
{/snippet}

<form class="space-y-4" onsubmit={submit}>
	<Section title={$t`Programme`}><Basics {draft} {others} /></Section>
	{@render modules('basics')}
	{#if locked}
		<p class="rounded-md border border-amber-500/40 bg-amber-500/10 px-3 py-2 text-sm" data-testid="locked-hint">
			{$t`Shooters are registered for this programme: the series stay as they are, and age categories and awards can be added but not removed.`}
		</p>
	{/if}
	<Section title={$t`Series (${draft.totalShots} shots, max. ${draft.totalPoints} points)`}><Series {draft} {locked} /></Section>
	{@render modules('series')}
	<Section
		title={$t`Prices (CHF)`}
		hint={draft.repeatable
			? $t`Only the checked age categories may shoot this programme. Further attempts cost the same unless a price is given; 0 means paid once.`
			: $t`Only the checked age categories may shoot this programme.`}
	>
		<Prices {draft} {locked} />
	</Section>
	{@render modules('prices')}
	<div class="sticky bottom-0 z-10 flex items-center gap-3 rounded-lg border bg-card/95 px-4 py-3 backdrop-blur">
		<FormError message={error} class="flex-1" />
		<Button type="submit" class="ml-auto" disabled={busy}>{submitLabel}</Button>
	</div>
</form>

{#if programme}
	<div class="mt-8 space-y-4">
		<Section title={$t`Use`} hint={$t`Saved at once, also while shooters are registered.`}>
			<Usage {tenantId} {event} {programme} {draft} />
		</Section>
		{@render modules('usage')}
	</div>
{/if}
