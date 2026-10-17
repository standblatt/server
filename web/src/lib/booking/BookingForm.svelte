<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import type { BookingInfo, BookingRequestInput, BookingShooter } from '../api';
	import { byDay, slotKey, slotTime } from '../booking';
	import { Button } from '../components/ui/button';
	import { Input } from '../components/ui/input';
	import { Textarea } from '../components/ui/textarea';
	import FormError from '../FormError.svelte';
	import FormField from '../FormField.svelte';
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

	/** Most lanes a club can ask for in a slot: what is free, capped by the organizer's limit. */
	const most = (free: number) => (info.maxLanes > 0 ? Math.min(free, info.maxLanes) : free);

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
			await onsubmit({ club: form.club, contact: form.contact, email: form.email, phone: form.phone, note: form.note, items, shooters: form.shooters });
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
			<section class="overflow-hidden rounded-lg border bg-card">
				<h3 class="border-b px-4 py-2 text-sm font-semibold">{day}</h3>
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
							<tr>
								<td class="w-32 px-4 py-1.5 font-medium tabular-nums">{slotTime(sl.start, info.timezone)}–{slotTime(sl.end, info.timezone)}</td>
								{#each info.ranges as r (r.id)}
									{@const free = sl.free.find((f) => f.rangeId === r.id)?.free ?? 0}
									{@const key = slotKey(r.id, sl.start)}
									<td class="px-2 py-1.5">
										{#if free === 0 && !form.lanes[key]}
											<span class="text-muted-foreground">{$t`full`}</span>
										{:else}
											<span class="flex items-center gap-2">
												<Input
													class="w-20"
													type="number"
													min="0"
													max={most(free)}
													placeholder="0"
													bind:value={form.lanes[key]}
													aria-label={$t`Lanes on ${r.name} at ${slotTime(sl.start, info.timezone)}`}
												/>
												<span class="text-xs whitespace-nowrap text-muted-foreground">{$t`${free} free`}</span>
											</span>
										{/if}
									</td>
								{/each}
							</tr>
						{/each}
					</tbody>
				</table>
			</section>
		{/each}
		<p class="text-sm text-muted-foreground">{$plural(chosen, { one: '# slot chosen', other: '# slots chosen' })}</p>
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
	{#if !readonly}<div><Button type="submit" disabled={busy}>{submitLabel}</Button></div>{/if}
</form>
