<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { toast } from 'svelte-sonner';
	import { api, call, type BookingConfig, type BookingSettings, type Event } from '../api';
	import { Button } from '../components/ui/button';
	import { Input } from '../components/ui/input';
	import FormError from '../FormError.svelte';
	import FormField from '../FormField.svelte';
	import { intlLocale } from '../locale.svelte';

	let { tenantId, event, settings, onsaved }: { tenantId: number; event: Event; settings: BookingSettings; onsaved: (c: BookingConfig) => void } = $props();

	/** The form state from the stored settings, with an empty pair for every day of the event. */
	function formFrom(s: BookingSettings) {
		const days: Record<string, { from: string; to: string }> = {};
		for (const d of event.days) days[d] = { from: s.days?.[d]?.from ?? '', to: s.days?.[d]?.to ?? '' };
		return {
			from: s.from ?? '',
			to: s.to ?? '',
			breakFrom: s.breakFrom ?? '',
			breakTo: s.breakTo ?? '',
			deadline: s.deadline ?? '',
			slotMinutes: (s.slotMinutes || null) as number | null,
			maxLanes: (s.maxLanes || null) as number | null,
			days
		};
	}

	// The parent re-creates this form with {#key} whenever the stored settings change.
	// svelte-ignore state_referenced_locally
	const form = $state(formFrom(settings));
	let error = $state('');

	/** Saves the adjusted settings; empty fields keep the defaults. */
	async function save(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		const days: Record<string, { from: string; to: string }> = {};
		for (const [d, h] of Object.entries(form.days)) if (h.from && h.to) days[d] = h;
		const body: BookingSettings = {
			from: form.from,
			to: form.to,
			breakFrom: form.breakFrom,
			breakTo: form.breakTo,
			deadline: form.deadline,
			slotMinutes: Number(form.slotMinutes) || 0,
			maxLanes: Number(form.maxLanes) || 0,
			days
		};
		try {
			onsaved(await call(api.PUT('/api/tenants/{tenantId}/events/{eventId}/booking', { params: { path: { tenantId, eventId: event.id } }, body })));
			toast.success($t`Booking settings saved`);
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** A day of the event as a short label. */
	const dayLabel = (d: string) => new Date(`${d}T12:00:00`).toLocaleDateString(intlLocale(), { weekday: 'short', day: 'numeric', month: 'short' });
</script>

<details class="rounded-lg border bg-card p-4">
	<summary class="cursor-pointer text-sm font-medium">{$t`Adjust`}</summary>
	<form class="mt-4 grid gap-4 sm:grid-cols-2" onsubmit={save}>
		<FormField label={$t`Shooting from`} hint={$t`Empty: the shooting times of the day, else 08:00`}><Input type="time" bind:value={form.from} /></FormField>
		<FormField label={$t`Shooting until`} hint={$t`Empty: the shooting times of the day, else 17:00`}><Input type="time" bind:value={form.to} /></FormField>
		<FormField label={$t`Break from`} hint={$t`Empty: no break`}><Input type="time" bind:value={form.breakFrom} /></FormField>
		<FormField label={$t`Break until`}><Input type="time" bind:value={form.breakTo} /></FormField>
		<FormField label={$t`Slot length (minutes)`}><Input type="number" min="5" max="480" bind:value={form.slotMinutes} placeholder="30" /></FormField>
		<FormField label={$t`Lanes per slot and club`} hint={$t`Empty: no limit`}><Input type="number" min="1" bind:value={form.maxLanes} /></FormField>
		<FormField label={$t`Deadline`} hint={$t`Empty: 14 days before the first day`}><Input type="date" bind:value={form.deadline} /></FormField>
		<div class="grid gap-2 sm:col-span-2">
			<div class="text-sm font-medium">{$t`Days with other hours`}</div>
			{#each event.days as d (d)}
				<div class="flex flex-wrap items-center gap-2 text-sm">
					<span class="w-32">{dayLabel(d)}</span>
					<Input class="w-32" type="time" bind:value={form.days[d].from} aria-label={$t`${dayLabel(d)} from`} />
					<span>–</span>
					<Input class="w-32" type="time" bind:value={form.days[d].to} aria-label={$t`${dayLabel(d)} until`} />
				</div>
			{/each}
		</div>
		<FormError message={error} class="sm:col-span-2" />
		<div class="sm:col-span-2"><Button type="submit">{$t`Save`}</Button></div>
	</form>
</details>
