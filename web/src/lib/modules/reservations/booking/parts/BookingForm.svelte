<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import Minus from '@lucide/svelte/icons/minus';
	import Plus from '@lucide/svelte/icons/plus';
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

	type Draft = { club: string; contact: string; email: string; phone: string; note: string; lanes: Record<string, number>; shooters: BookingShooter[] };

	/** The form state from the request being changed, a saved draft or empty. */
	function start(): Draft {
		if (initial) {
			const lanes: Record<string, number> = {};
			for (const it of initial.items) lanes[slotKey(it.rangeId, it.start)] = it.lanes;
			return { club: initial.club, contact: initial.contact, email: initial.email, phone: initial.phone ?? '', note: initial.note ?? '', lanes, shooters: initial.shooters ?? [] };
		}
		try {
			const saved = draftKey && localStorage.getItem(draftKey);
			if (saved) return JSON.parse(saved);
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
	const chosen = $derived(Object.values(form.lanes).reduce((n, v) => n + (Number(v) > 0 ? 1 : 0), 0));
	const total = $derived(Object.values(form.lanes).reduce((n, v) => n + (Number(v) || 0), 0));

	/** Most lanes a club can ask for in a slot: what is free, capped by the organizer's limit. */
	const most = (free: number) => (info.maxLanes > 0 ? Math.min(free, info.maxLanes) : free);

	/** Lanes chosen in a slot on a range. */
	const lanesAt = (rangeId: number, start: string) => Number(form.lanes[slotKey(rangeId, start)]) || 0;

	/** Free lanes of a range in a slot. */
	const freeAt = (sl: BookingInfo['slots'][number], rangeId: number) => sl.free.find((f) => f.rangeId === rangeId)?.free ?? 0;

	/** Chosen slots of a day, counted per range like the total. */
	const chosenIn = (slots: BookingInfo['slots']) => slots.reduce((n, sl) => n + info.ranges.filter((r) => lanesAt(r.id, sl.start) > 0).length, 0);

	// Open at first: every day of a short event, else the days with a choice; then the club decides.
	// svelte-ignore state_referenced_locally
	const opened = new Set(days.filter(([, slots]) => days.length <= 2 || chosenIn(slots) > 0).map(([d]) => d));

	/** Changes the lanes of a slot by one, between 0 and the most allowed. */
	function step(key: string, by: number, max: number) {
		form.lanes[key] = Math.max(0, Math.min(max, (Number(form.lanes[key]) || 0) + by));
	}

	/** Builds the request from the form and hands it to the page. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		const items = info.slots.flatMap((sl) =>
			info.ranges
				.map((r) => ({ rangeId: r.id, start: sl.start, end: sl.end, lanes: Number(form.lanes[slotKey(r.id, sl.start)] ?? 0) }))
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
			{$t`Enter how many lanes you need in each slot. The organizer assigns the lane numbers when confirming.`}
			{#if info.maxLanes > 0}{$t`At most ${info.maxLanes} lanes per slot.`}{/if}
		</p>
		{#each days as [day, slots] (day)}
			{@const n = chosenIn(slots)}
			<details class="group overflow-hidden rounded-lg border bg-card" open={opened.has(day)}>
				<summary class="cursor-pointer px-4 py-2 text-sm font-semibold group-open:border-b">
					{day}
					{#if n > 0}<Badge variant="secondary" class="ml-2">{$plural(n, { one: '# slot chosen', other: '# slots chosen' })}</Badge>{/if}
				</summary>
				<table class="w-full text-sm">
					{#if multi}
						<thead class="text-left text-xs text-muted-foreground">
							<tr>
								<th class="px-4 py-1.5 font-medium">{$t`Time`}</th>
								{#each info.ranges as r (r.id)}<th class="px-2 py-1.5 font-medium">{r.name}</th>{/each}
							</tr>
						</thead>
					{/if}
					<tbody class="divide-y">
						{#each slots as sl (sl.start)}
							{@const time = slotTime(sl.start, info.timezone)}
							{@const picked = info.ranges.some((r) => lanesAt(r.id, sl.start) > 0)}
							{@const full = info.ranges.every((r) => freeAt(sl, r.id) === 0)}
							<tr class={[picked && 'bg-primary/5', full && !picked && 'text-muted-foreground']}>
								<td class={['w-32 px-4 py-1.5 tabular-nums', picked ? 'font-semibold' : 'font-medium']}>{time}–{slotTime(sl.end, info.timezone)}</td>
								{#each info.ranges as r (r.id)}
									{@const free = freeAt(sl, r.id)}
									{@const key = slotKey(r.id, sl.start)}
									{@const max = most(free)}
									<td class="px-2 py-1.5">
										{#if free === 0 && !form.lanes[key]}
											<span class="text-muted-foreground">{$t`full`}</span>
										{:else}
											<span class="flex items-center gap-1">
												<Button
													variant="outline"
													size="icon-sm"
													disabled={lanesAt(r.id, sl.start) === 0}
													aria-label={$t`Fewer lanes on ${r.name} at ${time}`}
													onclick={() => step(key, -1, max)}><Minus /></Button
												>
												<Input
													class="w-16 text-center"
													type="number"
													min="0"
													{max}
													placeholder="0"
													bind:value={form.lanes[key]}
													aria-label={$t`Lanes on ${r.name} at ${time}`}
												/>
												<Button
													variant="outline"
													size="icon-sm"
													disabled={lanesAt(r.id, sl.start) >= max}
													aria-label={$t`More lanes on ${r.name} at ${time}`}
													onclick={() => step(key, 1, max)}><Plus /></Button
												>
												<span class={['ml-1 text-xs whitespace-nowrap', free <= r.laneCount / 4 ? 'text-amber-700 dark:text-amber-400' : 'text-muted-foreground']}>{$t`${free} free`}</span>
											</span>
										{/if}
									</td>
								{/each}
							</tr>
						{/each}
					</tbody>
				</table>
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
