<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { untrack } from 'svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { programmeTitle } from '$lib/format';
	import Section from '$lib/programme/Section.svelte';
	import { keys } from '$lib/tenant';
	import type { ProgrammeProps } from '..';

	let { tenantId, event, programme }: ProgrammeProps = $props();
	const queryClient = useQueryClient();
	const start = untrack(() => programme?.specialPrizes);
	let ranks = $state(start?.ranks ? start.ranks : 3);
	let minCategory = $state(start?.ranks ? start.minCategory : 10);
	let minClub = $state(start?.ranks ? start.minClub : 4);
	let error = $state('');

	/** Saves the special prize rules of the programme. */
	async function save(e: SubmitEvent) {
		e.preventDefault();
		if (!programme) return;
		error = '';
		try {
			await call(
				api.PUT('/api/tenants/{tenantId}/programmes/{programmeId}/special-prizes', {
					params: { path: { tenantId, programmeId: programme.id } },
					body: { ranks: Number(ranks), minCategory: Number(minCategory), minClub: Number(minClub) }
				})
			);
			await queryClient.invalidateQueries({ queryKey: keys.programmes(tenantId, event.id) });
			await queryClient.invalidateQueries({ queryKey: ['tenants', tenantId, 'programme', programme.id] });
			toast.success(Number(ranks) ? $t`Special prizes of ${programmeTitle(programme)} saved` : $t`Special prizes of ${programmeTitle(programme)} turned off`);
		} catch (err) {
			error = (err as Error).message;
		}
	}
</script>

{#if programme}
	<Section title={$t`Special prizes`} hint={$t`The best ranks of every age category get a special prize. Shooters whose club has too few participants in this programme are skipped. Can be changed while shooters are registered.`}>
		<form class="grid items-end gap-4 sm:grid-cols-3" onsubmit={save} aria-label={$t`Special prizes`}>
			<FormField label={$t`Ranks per category`} hint={$t`0 turns special prizes off.`}><Input type="number" min="0" max="10" bind:value={ranks} required /></FormField>
			<FormField label={$t`Minimum participants per category`}><Input type="number" min="0" bind:value={minCategory} required /></FormField>
			<FormField label={$t`Minimum participants of the club in this programme`}><Input type="number" min="0" bind:value={minClub} required /></FormField>
			<FormError message={error} class="sm:col-span-2" />
			<Button type="submit" variant="outline" class="sm:col-start-3">{$t`Save special prizes`}</Button>
		</form>
	</Section>
{/if}
