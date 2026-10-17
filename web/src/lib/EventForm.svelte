<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { untrack } from 'svelte';
	import type { Club, EventInput, Range } from './api';
	import { Button } from './components/ui/button';
	import { Checkbox } from './components/ui/checkbox';
	import { Input } from './components/ui/input';
	import * as NativeSelect from './components/ui/native-select';
	import { Switch } from './components/ui/switch';
	import FormError from './FormError.svelte';
	import FormField from './FormField.svelte';
	import { currentLocale, localeNames, type Locale } from './locale.svelte';
	import ScheduleEditor from './ScheduleEditor.svelte';

	let {
		ranges,
		clubs,
		initial,
		submitLabel,
		onsubmit
	}: {
		ranges: Range[];
		clubs: Club[];
		initial?: EventInput;
		submitLabel: string;
		onsubmit: (input: EventInput) => Promise<void>;
	} = $props();

	const start: EventInput = untrack(
		() =>
			initial ?? {
				name: '',
				startsOn: new Date().toISOString().slice(0, 10),
				rangeIds: ranges.length === 1 ? [ranges[0].id] : [],
				schedule: {},
				organizer: '',
				locale: currentLocale(),
				settings: { showLaneStatus: false, printAwardValue: true, ammoIssue: false }
			}
	);
	let name = $state(start.name);
	let startsOn = $state(start.startsOn);
	let rangeIds = $state<number[]>(untrack(() => [...start.rangeIds]));
	const sched = untrack(() => $state.snapshot(start.schedule ?? {}));
	let days = $state((sched.days ?? []).filter((d) => d.date !== start.startsOn || d.times?.length).map((d) => ({ ...d, times: d.times ?? [] })));
	let period = $state(!!sched.from);
	let from = $state(sched.from ?? '');
	let to = $state(sched.to ?? '');
	let weekly = $state((sched.weekly ?? []).map((w) => ({ ...w, times: w.times ?? [] })));
	/** Adds or removes a range from the selected ranges. */
	function toggleRange(id: number, on: boolean) {
		rangeIds = on ? [...rangeIds, id] : rangeIds.filter((r) => r !== id);
	}
	let organizer = $state(start.organizer ?? '');
	let locale = $state<Locale>(start.locale ?? currentLocale());
	let showLaneStatus = $state(start.settings?.showLaneStatus ?? false);
	let ammoIssue = $state(start.settings?.ammoIssue ?? false);
	let printAwardValue = $state(start.settings?.printAwardValue ?? true);
	let excludedClubId = $state(start.settings?.excludedClubId ?? 0);
	let error = $state('');
	let busy = $state(false);

	/** Collects the form into an EventInput and passes it to onsubmit. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		busy = true;
		error = '';
		try {
			await onsubmit({
				name,
				startsOn,
				rangeIds: ranges.filter((r) => rangeIds.includes(r.id)).map((r) => r.id),
				schedule: period ? { days: $state.snapshot(days), from, to, weekly: $state.snapshot(weekly) } : { days: $state.snapshot(days) },
				organizer,
				timezone: start.timezone,
				locale,
				settings: { showLaneStatus, printAwardValue, ammoIssue, excludedClubId: excludedClubId || undefined }
			});
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

{#if ranges.length === 0}
	<p class="text-sm text-muted-foreground">{$t`Create a range first.`}</p>
{:else}
	<form class="grid gap-4 sm:grid-cols-2" onsubmit={submit}>
		<FormField label={$t`Name`} class="sm:col-span-2"><Input bind:value={name} required /></FormField>
		<FormField label={$t`Date`} hint={$t`The main day. Age categories follow the SSV rule: competition year minus year of birth, so only the year of this date matters.`}><Input type="date" bind:value={startsOn} required /></FormField>
		<FormField label={$t`Organizer`}><Input bind:value={organizer} /></FormField>
		<FormField label={$t`Ranges`} class="sm:col-span-2" hint={ranges.length > 1 ? $t`Tick every range the event is shot on.` : ''}>
			<div class="grid gap-1.5 rounded-md border p-3 sm:grid-cols-2" role="group" aria-label={$t`Ranges`}>
				{#each ranges as s (s.id)}
					<label class="flex items-center gap-2 text-sm">
						<Checkbox checked={rangeIds.includes(s.id)} onCheckedChange={(v) => toggleRange(s.id, v)} />
						{$t`${s.name} (${s.laneCount} lanes)`}
					</label>
				{/each}
			</div>
		</FormField>
		<ScheduleEditor mainDay={startsOn} bind:days bind:period bind:from bind:to bind:weekly />
		<FormField label={$t`Language of documents and kiosk`}>
			<NativeSelect.Root class="w-full" bind:value={locale}>
				{#each Object.entries(localeNames) as [value, label] (value)}<NativeSelect.Option {value}>{label}</NativeSelect.Option>{/each}
			</NativeSelect.Root>
		</FormField>
		<FormField label={$t`Club excluded from the club competition`} class="sm:col-span-2">
			<NativeSelect.Root class="w-full" bind:value={excludedClubId}>
				<NativeSelect.Option value={0}>{$t`none`}</NativeSelect.Option>
				{#each clubs as c (c.id)}<NativeSelect.Option value={c.id}>{c.name}</NativeSelect.Option>{/each}
			</NativeSelect.Root>
		</FormField>
		<label class="flex items-center justify-between gap-3 rounded-md border p-3 text-sm sm:col-span-2">
			{$t`Show lane status on the rankings`}
			<Switch bind:checked={showLaneStatus} />
		</label>
		<label class="flex items-center justify-between gap-3 rounded-md border p-3 text-sm sm:col-span-2">
			{$t`Print the award value on wreath cards`}
			<Switch bind:checked={printAwardValue} />
		</label>
		<label class="flex items-center justify-between gap-3 rounded-md border p-3 text-sm sm:col-span-2">
			{$t`The organizer hands out the ammunition (Ammo button on the dashboard)`}
			<Switch bind:checked={ammoIssue} />
		</label>
		<FormError message={error} class="sm:col-span-2" />
		<div class="flex justify-end sm:col-span-2"><Button type="submit" disabled={busy}>{submitLabel}</Button></div>
	</form>
{/if}
