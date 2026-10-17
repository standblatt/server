<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { untrack } from 'svelte';
	import type { ShooterInput } from './api';
	import { Button } from './components/ui/button';
	import { Input } from './components/ui/input';
	import FormError from './FormError.svelte';
	import FormField from './FormField.svelte';

	let {
		initial,
		submitLabel,
		onsubmit
	}: {
		initial?: ShooterInput;
		submitLabel: string;
		onsubmit: (input: ShooterInput) => Promise<void>;
	} = $props();

	const start: ShooterInput = untrack(() => initial ?? { firstname: '', lastname: '' });
	let licenseNumber = $state(start.licenseNumber ?? '');
	let firstname = $state(start.firstname);
	let lastname = $state(start.lastname);
	let birthday = $state(start.birthday ?? '');
	let postalcode = $state(start.postalcode ?? '');
	let city = $state(start.city ?? '');
	let error = $state('');
	let busy = $state(false);

	/** Passes the form fields to onsubmit and shows any error. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		busy = true;
		error = '';
		try {
			await onsubmit({ licenseNumber, firstname, lastname, birthday, postalcode, city });
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

<form class="grid gap-4 sm:grid-cols-2" onsubmit={submit}>
	<FormField label={$t`First name`}><Input bind:value={firstname} required /></FormField>
	<FormField label={$t`Last name`}><Input bind:value={lastname} required /></FormField>
	<FormField label={$t`License number`}><Input bind:value={licenseNumber} inputmode="numeric" /></FormField>
	<FormField label={$t`Birthday`}><Input type="date" bind:value={birthday} /></FormField>
	<FormField label={$t`Postal code`}><Input bind:value={postalcode} /></FormField>
	<FormField label={$t`City`}><Input bind:value={city} /></FormField>
	<FormError message={error} class="sm:col-span-2" />
	<div class="flex justify-end sm:col-span-2"><Button type="submit" disabled={busy}>{submitLabel}</Button></div>
</form>
