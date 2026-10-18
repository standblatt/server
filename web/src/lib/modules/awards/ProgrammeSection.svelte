<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { ageCategories } from '$lib/categories';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import FormField from '$lib/FormField.svelte';
	import { ageCategoryLabel } from '$lib/labels';
	import { emptyLimits } from '$lib/programme/draft.svelte';
	import Section from '$lib/programme/Section.svelte';
	import type { ProgrammeProps } from '..';

	let { draft, locked }: ProgrammeProps = $props();
</script>

<Section title={$t`Awards`} hint={$t`Minimum total per age category; leave empty where the award is not available.`}>
	<div class="space-y-3">
		{#each draft.awards as a, i (i)}
			<fieldset class="rounded-md border p-4">
				<legend class="px-1 text-sm font-medium">{$t`Award ${i + 1}`}</legend>
				<div class="mb-4 flex flex-wrap items-end gap-3">
					<FormField label={$t`Name`} class="min-w-40 flex-1"><Input bind:value={a.name} required /></FormField>
					<FormField label={$t`Kind`}>
						<NativeSelect.Root bind:value={a.kind}>
							<NativeSelect.Option value="card">{$t`Wreath card`}</NativeSelect.Option>
							<NativeSelect.Option value="natural">{$t`Prize in kind`}</NativeSelect.Option>
							<NativeSelect.Option value="cash">{$t`Cash prize`}</NativeSelect.Option>
						</NativeSelect.Root>
					</FormField>
					<FormField label={$t`Value (CHF)`}><Input class="w-28" type="number" min="0" step="0.05" bind:value={a.value} /></FormField>
					<Button variant="ghost" size="icon" title={$t`Remove award`} aria-label={$t`Remove award ${i + 1}`} onclick={() => draft.awards.splice(i, 1)} disabled={locked && i < draft.fixedAwards}><Trash /></Button>
				</div>
				<div class="grid grid-cols-2 gap-3 *:min-w-0 sm:grid-cols-3 lg:grid-cols-5">
					{#each ageCategories as c (c)}
						<FormField label={$t(ageCategoryLabel(c))}><Input class="h-8" type="number" min="0" max={draft.totalPoints} step="0.1" bind:value={a.limits[c]} /></FormField>
					{/each}
				</div>
			</fieldset>
		{/each}
	</div>
	<Button variant="outline" size="sm" class="mt-2" onclick={() => draft.awards.push({ name: '', kind: 'card', value: 0, limits: emptyLimits() })}><Plus />{$t`Add award`}</Button>
</Section>
