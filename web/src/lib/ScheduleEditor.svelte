<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import X from '@lucide/svelte/icons/x';
	import type { ScheduleDay, ShootingTimes, WeeklyTimes } from './api';
	import { shootingDays, weekdayName } from './schedule';
	import { Button } from './components/ui/button';
	import { Checkbox } from './components/ui/checkbox';
	import { Input } from './components/ui/input';
	import { Switch } from './components/ui/switch';
	import FormField from './FormField.svelte';
	import { intlLocale } from './locale.svelte';

	let {
		mainDay,
		days = $bindable(),
		period = $bindable(),
		from = $bindable(),
		to = $bindable(),
		weekly = $bindable()
	}: {
		mainDay: string;
		days: ScheduleDay[];
		period: boolean;
		from: string;
		to: string;
		weekly: WeeklyTimes[];
	} = $props();

	const count = $derived(shootingDays(period ? { days, from, to, weekly } : { days }, mainDay).length);

	/** Appends a shooting day one day after the last one (or after the main day). */
	function addDay() {
		const d = new Date(days[days.length - 1]?.date || mainDay);
		d.setDate(d.getDate() + 1);
		days.push({ date: d.toISOString().slice(0, 10), times: [] });
	}

	/** Ticks or unticks a weekday of the period. */
	function toggleWeekday(weekday: number, on: boolean) {
		weekly = on ? [...weekly, { weekday, times: [] }].sort((a, b) => a.weekday - b.weekday) : weekly.filter((w) => w.weekday !== weekday);
	}
</script>

{#snippet timesEditor(times: ShootingTimes[], label: string)}
	{#each times as _, i (i)}
		<span class="flex items-center gap-1">
			<Input type="time" class="w-28" bind:value={times[i].from} required aria-label={$t`${label} from`} />
			<span>–</span>
			<Input type="time" class="w-28" bind:value={times[i].to} required aria-label={$t`${label} until`} />
			<Button variant="ghost" size="icon-sm" aria-label={$t`Remove time`} onclick={() => times.splice(i, 1)}><X /></Button>
		</span>
	{/each}
	<Button variant="ghost" size="sm" onclick={() => times.push({ from: '', to: '' })}><Plus />{$t`Time`}</Button>
{/snippet}

<FormField label={$t`Shooting days`} class="sm:col-span-2" hint={$t`Single days with optional shooting times, e.g. pre-shooting days or the official weekend. The main day is always a shooting day.`}>
	<div class="grid gap-2">
		{#each days as day, i (i)}
			<div class="flex flex-wrap items-center gap-2">
				<Input type="date" class="w-40" bind:value={day.date} required aria-label={$t`Shooting day`} />
				{@render timesEditor(day.times!, day.date)}
				<Button variant="ghost" size="icon-sm" aria-label={$t`Remove day`} onclick={() => days.splice(i, 1)}><Trash /></Button>
			</div>
		{/each}
		<div><Button variant="outline" size="sm" onclick={addDay}><Plus />{$t`Add day`}</Button></div>
	</div>
</FormField>

<label class="flex items-center justify-between gap-3 rounded-md border p-3 text-sm sm:col-span-2">
	<span>
		{$t`Period (season)`}
		<span class="block text-muted-foreground">{$t`For events that run over weeks or months, like the Feldschiessen.`}</span>
	</span>
	<Switch bind:checked={period} />
</label>
{#if period}
	<FormField label={$t`First day`}><Input type="date" bind:value={from} required /></FormField>
	<FormField label={$t`Last day`}><Input type="date" bind:value={to} min={from} required /></FormField>
	<FormField label={$t`Shooting times`} class="sm:col-span-2" hint={$t`Weekdays of the period that are shot on. None ticked: every day of the period.`}>
		<div class="grid gap-2 rounded-md border p-3" role="group" aria-label={$t`Shooting times`}>
			{#each [1, 2, 3, 4, 5, 6, 7] as wd (wd)}
				{@const w = weekly.find((x) => x.weekday === wd)}
				<div class="flex flex-wrap items-center gap-2 text-sm">
					<label class="flex w-32 items-center gap-2">
						<Checkbox checked={!!w} onCheckedChange={(v) => toggleWeekday(wd, v)} />
						{weekdayName(wd, intlLocale())}
					</label>
					{#if w}{@render timesEditor(w.times!, weekdayName(wd, intlLocale()))}{/if}
				</div>
			{/each}
		</div>
	</FormField>
{/if}
<p class="text-sm text-muted-foreground sm:col-span-2">{$plural(count, { one: '# shooting day', other: '# shooting days' })}</p>
