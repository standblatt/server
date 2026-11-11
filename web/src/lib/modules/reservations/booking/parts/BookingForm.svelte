<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import type { BookingInfo, BookingRequestInput, BookingShooter } from '$lib/api';
	import { byDay, slotKey, slotTime } from '../booking';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import ShooterList from './ShooterList.svelte';

	let {
		info,
		initial,
		readonly = false,
		submitLabel,
		draftKey,
		onsubmit
	}: {
		info: BookingInfo;
		initial?: BookingRequestInput;
		readonly?: boolean;
		submitLabel: string;
		draftKey?: string;
		onsubmit: (input: BookingRequestInput) => Promise<void>;
	} = $props();

	type Draft = { club: string; contact: string; email: string; phone: string; note: string; lanes: Record<string, number[]>; shooters: BookingShooter[] };

	/** The form state from the request being changed, a saved draft or empty. */
	function start(): Draft {
		if (initial) {
			const lanes: Record<string, number[]> = {};
			for (const it of initial.items) lanes[slotKey(it.rangeId, it.start)] = it.wanted ?? [];
			return { club: initial.club, contact: initial.contact, email: initial.email, phone: initial.phone ?? '', note: initial.note ?? '', lanes, shooters: initial.shooters ?? [] };
		}
		try {
			const saved = draftKey && localStorage.getItem(draftKey);
			if (saved) {
				const draft: Draft = JSON.parse(saved);
				for (const [k, v] of Object.entries(draft.lanes)) if (!Array.isArray(v)) delete draft.lanes[k];
				return draft;
			}
		} catch {
			/* no storage: start empty */
		}
		return { club: '', contact: '', email: '', phone: '', note: '', lanes: {}, shooters: [] };
	}

	const form = $state(start());
	let error = $state('');
	let busy = $state(false);

	$effect(() => {
		if (!draftKey || initial) return;
		const json = JSON.stringify(form);
		try {
			localStorage.setItem(draftKey, json);
		} catch {
			/* drafts are a convenience */
		}
	});

	const days = $derived(byDay(info.slots, info.timezone));
	const multi = $derived(info.ranges.length > 1);
	const chosen = $derived(Object.values(form.lanes).filter((v) => v.length > 0).length);
	const total = $derived(Object.values(form.lanes).reduce((n, v) => n + v.length, 0));

	/** Lanes chosen in a slot on a range. */
	const lanesAt = (rangeId: number, start: string) => form.lanes[slotKey(rangeId, start)] ?? [];

	/** Free and requested lanes of a range in a slot. */
	const freeAt = (sl: BookingInfo['slots'][number], rangeId: number) => sl.free.find((f) => f.rangeId === rangeId) ?? { lanes: [] as number[], requested: [] as number[] };

	/** Chosen slots of a day, counted per range like the total. */
	const chosenIn = (slots: BookingInfo['slots']) => slots.reduce((n, sl) => n + info.ranges.filter((r) => lanesAt(r.id, sl.start).length > 0).length, 0);

	// Open at first: every day of a short event, else the days with a choice; then the club decides.
	// svelte-ignore state_referenced_locally
	const opened = new Set(days.filter(([, slots]) => days.length <= 2 || chosenIn(slots) > 0).map(([d]) => d));

	/** Adds a lane to the choice of a slot or takes it out. */
	function toggle(key: string, lane: number) {
		const mine = form.lanes[key] ?? [];
		form.lanes[key] = mine.includes(lane) ? mine.filter((l) => l !== lane) : [...mine, lane].sort((a, b) => a - b);
	}

	/** Colours of a lane button (traffic light): chosen, free, asked for by another club or taken. */
	const laneClass = (state: 'chosen' | 'free' | 'requested' | 'taken') =>
		({
			chosen: 'border-[3px] border-green-500 bg-brand-900 text-white hover:bg-brand-800 hover:text-white dark:bg-brand-100 dark:text-brand-950 dark:hover:bg-white dark:hover:text-brand-950',
			free: 'border-green-600 bg-green-500 text-white hover:bg-green-600 hover:text-white dark:border-green-500 dark:bg-green-600 dark:text-green-50 dark:hover:bg-green-500',
			requested: 'border-yellow-500 bg-yellow-400 text-yellow-950 hover:bg-yellow-500 hover:text-yellow-950 dark:bg-yellow-600 dark:text-yellow-50 dark:hover:bg-yellow-500',
			taken: 'border-red-600 bg-red-500 text-white opacity-55 disabled:opacity-55 dark:bg-red-700 dark:text-red-100 dark:opacity-50 dark:disabled:opacity-50'
		})[state];

	/** Builds the request from the form and hands it to the page. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		const items = info.slots.flatMap((sl) =>
			info.ranges
				.map((r) => {
					const wanted = lanesAt(r.id, sl.start);
					return { rangeId: r.id, start: sl.start, end: sl.end, lanes: wanted.length, wanted };
				})
				.filter((it) => it.lanes > 0)
		);
		if (items.length === 0) {
			error = $t`Choose at least one slot.`;
			return;
		}
		busy = true;
		try {
			const shooters = form.shooters.map(({ lastname, firstname, birthday, licenseNumber }) => ({ lastname, firstname, birthday, licenseNumber }));
			await onsubmit({ club: form.club, contact: form.contact, email: form.email, phone: form.phone, note: form.note, items, shooters });
			if (draftKey && !initial) {
				try {
					localStorage.removeItem(draftKey);
				} catch {
					/* nothing to clean */
				}
			}
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

<form class="grid gap-6" onsubmit={submit}>
	<fieldset class="grid gap-3" disabled={readonly}>
		<legend class="mb-1 text-base font-semibold">{$t`Lanes`}</legend>
		<p class="text-sm text-muted-foreground">
			{$t`Tap the lanes you want in each slot. The organizer confirms them or assigns others if they are taken by then.`}
			{#if info.maxLanes > 0}{$t`At most ${info.maxLanes} lanes per slot.`}{/if}
		</p>
		<div class="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
			<span class="flex items-center gap-1.5"><span class={['size-4 rounded border', laneClass('free')]}></span>{$t`free`}</span>
			<span class="flex items-center gap-1.5"><span class={['size-4 rounded border', laneClass('requested')]}></span>{$t`asked for by another club`}</span>
			<span class="flex items-center gap-1.5"><span class={['size-4 rounded border', laneClass('taken')]}></span>{$t`taken`}</span>
			<span class="flex items-center gap-1.5"><span class={['size-4 rounded border', laneClass('chosen')]}></span>{$t`your choice`}</span>
		</div>
		{#each days as [day, slots] (day)}
			{@const n = chosenIn(slots)}
			<details class="group overflow-hidden rounded-lg border bg-card" open={opened.has(day)}>
				<summary class="cursor-pointer px-4 py-2 text-sm font-semibold group-open:border-b">
					{day}
					{#if n > 0}<Badge variant="secondary" class="ml-2">{$plural(n, { one: '# slot chosen', other: '# slots chosen' })}</Badge>{/if}
				</summary>
				<div class="divide-y text-sm">
					{#each slots as sl (sl.start)}
						{@const time = slotTime(sl.start, info.timezone)}
						{@const picked = info.ranges.some((r) => lanesAt(r.id, sl.start).length > 0)}
						<div class={['grid gap-2 px-4 py-2 sm:grid-cols-[8rem_1fr]', picked && 'bg-primary/5']}>
							<span class={['pt-1 tabular-nums', picked ? 'font-semibold' : 'font-medium']}>{time}–{slotTime(sl.end, info.timezone)}</span>
							<div class="grid gap-2">
								{#each info.ranges as r (r.id)}
									{@const f = freeAt(sl, r.id)}
									{@const key = slotKey(r.id, sl.start)}
									{@const mine = lanesAt(r.id, sl.start)}
									{@const atMost = info.maxLanes > 0 && mine.length >= info.maxLanes}
									<div class="grid gap-1">
										{#if multi}<span class="text-xs text-muted-foreground">{r.name}</span>{/if}
										<div class="flex flex-wrap gap-1">
											{#each r.lanes as lane (lane)}
												{@const on = mine.includes(lane)}
												{@const state = on ? 'chosen' : !f.lanes.includes(lane) ? 'taken' : f.requested.includes(lane) ? 'requested' : 'free'}
												<Button
													variant="outline"
													size="sm"
													class={['h-8 w-9 px-0 tabular-nums disabled:opacity-100', laneClass(state)]}
													disabled={!on && (state === 'taken' || atMost)}
													aria-pressed={on}
													aria-label={$t`Lane ${lane} on ${r.name} at ${time}`}
													onclick={() => toggle(key, lane)}>{lane}</Button
												>
											{/each}
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			</details>
		{/each}
	</fieldset>

	<fieldset class="grid gap-4 sm:grid-cols-2" disabled={readonly}>
		<legend class="mb-1 text-base font-semibold">{$t`Club and contact`}</legend>
		<FormField label={$t`Club`} class="sm:col-span-2"><Input bind:value={form.club} required maxlength={200} autocomplete="organization" /></FormField>
		<FormField label={$t`Contact person`}><Input bind:value={form.contact} required maxlength={200} autocomplete="name" /></FormField>
		<FormField label={$t`Email`}><Input type="email" bind:value={form.email} required maxlength={200} autocomplete="email" /></FormField>
		<FormField label={$t`Phone`}><Input type="tel" bind:value={form.phone} maxlength={50} autocomplete="tel" /></FormField>
		<FormField label={$t`Note to the organizer`} class="sm:col-span-2"><Textarea bind:value={form.note} maxlength={2000} rows={3} /></FormField>
	</fieldset>

	<fieldset class="grid gap-2">
		<legend class="mb-1 text-base font-semibold">{$t`Shooters`}</legend>
		<p class="text-sm text-muted-foreground">{$t`Announce your shooters so they are found quickly at the registration desk. Programmes are chosen at the range.`}</p>
		<ShooterList bind:shooters={form.shooters} {readonly} />
	</fieldset>

	<FormError message={error} />
	{#if !readonly}
		<div class="sticky bottom-0 -mx-4 flex flex-wrap items-center gap-3 border-t bg-background/95 px-4 pt-3 pb-[max(0.75rem,env(safe-area-inset-bottom))] backdrop-blur">
			<span class="text-sm text-muted-foreground">{$plural(chosen, { one: '# slot chosen', other: '# slots chosen' })} · {$plural(total, { one: '# lane', other: '# lanes' })}</span>
			<Button type="submit" class="ml-auto" disabled={busy}>{submitLabel}</Button>
		</div>
	{/if}
</form>
