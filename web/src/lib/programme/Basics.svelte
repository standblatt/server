<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import FormField from '$lib/FormField.svelte';
	import type { ProgrammeDraft } from './draft.svelte';

	let { draft, others }: { draft: ProgrammeDraft; others: { name: string; variant?: string }[] } = $props();
	const names = $derived([...new Set(others.map((o) => o.name))].sort());
	const siblings = $derived(others.filter((o) => o.name === draft.name.trim()));
</script>

<div class="grid gap-4 sm:grid-cols-2">
	<FormField label={$t`Name`}>
		<Input bind:value={draft.name} required list="programme-names" />
		<datalist id="programme-names">{#each names as n (n)}<option value={n}></option>{/each}</datalist>
	</FormField>
	<FormField label={$t`Variant`} hint={siblings.length ? '' : $t`For a variant, pick the name of an existing programme: a shooter registers for only one programme of the same name.`}>
		<Input bind:value={draft.variant} placeholder={$t`e.g. Stgw 90`} />
	</FormField>
	{#if siblings.length}
		<p class="rounded-md border border-sky-500/40 bg-sky-500/10 px-3 py-2 text-sm sm:col-span-2" data-testid="variant-hint">
			{$t`Variant of "${draft.name.trim()}": shooters register for only one of ${siblings.length + 1} programmes.`}
			{#if !draft.variant.trim() || siblings.some((o) => (o.variant ?? '') === draft.variant.trim())}<strong class="text-destructive">{$t`Give each variant its own variant name.`}</strong>{/if}
		</p>
	{/if}
	<FormField label={$t`Description`} class="sm:col-span-2"><Textarea rows={2} bind:value={draft.description} /></FormField>
	<FormField label={$t`Barcode`}><Input bind:value={draft.barcode} /></FormField>
</div>
