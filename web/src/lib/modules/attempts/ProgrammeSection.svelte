<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import { Switch } from '$lib/components/ui/switch';
	import FormField from '$lib/FormField.svelte';
	import Section from '$lib/programme/Section.svelte';
	import type { ProgrammeProps } from '..';

	let { draft }: ProgrammeProps = $props();
</script>

<Section title={$t`Attempts`} hint={$t`How often a shooter may shoot the programme at this event, e.g. further attempts at a shooting festival or an annual programme.`}>
	<div class="grid items-end gap-4 sm:grid-cols-3">
		<FormField label={$t`Attempts per shooter`}>
			<Input type="number" min="1" max="1000" bind:value={draft.attempts} disabled={draft.unlimited} />
		</FormField>
		<label class="flex h-9 items-center gap-2 text-sm"><Checkbox bind:checked={draft.unlimited} />{$t`Unlimited`}</label>
		{#if draft.repeatable}
			<FormField label={$t`Counting attempt`}>
				<NativeSelect.Root class="w-full" bind:value={draft.attemptRule}>
					<NativeSelect.Option value="best">{$t`The best`}</NativeSelect.Option>
					<NativeSelect.Option value="first">{$t`The first`}</NativeSelect.Option>
					<NativeSelect.Option value="last">{$t`The last`}</NativeSelect.Option>
				</NativeSelect.Root>
			</FormField>
			<label class="flex items-center justify-between gap-3 text-sm sm:col-span-3">{$t`Only the counting attempt receives an award`} <Switch bind:checked={draft.awardCountingOnly} /></label>
		{/if}
	</div>
</Section>
