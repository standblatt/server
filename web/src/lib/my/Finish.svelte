<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import Angry from '@lucide/svelte/icons/angry';
	import Frown from '@lucide/svelte/icons/frown';
	import Laugh from '@lucide/svelte/icons/laugh';
	import Meh from '@lucide/svelte/icons/meh';
	import Smile from '@lucide/svelte/icons/smile';
	import Plus from '@lucide/svelte/icons/plus';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Sheet from '$lib/components/ui/sheet';
	import { Textarea } from '$lib/components/ui/textarea';
	import FormField from '$lib/FormField.svelte';
	import FormError from '$lib/FormError.svelte';
	import { myKey, type Weapon } from '$lib/my.svelte';
	import type { Draft } from './capture';

	let {
		open = $bindable(false),
		draft,
		saving,
		tags,
		weapons,
		onsave
	}: { open?: boolean; draft: Draft; saving: boolean; tags: string[]; weapons: Weapon[]; onsave: () => void } = $props();

	const queryClient = useQueryClient();
	let newWeapon = $state('');
	let weaponError = $state('');

	let newTag = $state('');
	const faces = [Angry, Frown, Meh, Smile, Laugh];
	const suggestions = $derived([...new Set([...draft.tags, ...tags])].slice(0, 12));

	/** Adds or removes a tag on the draft. */
	function toggle(tag: string) {
		draft.tags = draft.tags.includes(tag) ? draft.tags.filter((x) => x !== tag) : [...draft.tags, tag];
	}

	/** Creates a weapon with the typed name and selects it for the draft. */
	async function addWeapon(e: Event) {
		e.preventDefault();
		weaponError = '';
		try {
			const w = await call(api.POST('/api/my/weapons', { body: { name: newWeapon } }));
			await queryClient.invalidateQueries({ queryKey: myKey });
			draft.weaponId = w.id;
			newWeapon = '';
		} catch (err) {
			weaponError = (err as Error).message;
		}
	}

	/** Adds the typed tag to the draft unless it is empty or already there. */
	function addTag(e: Event) {
		e.preventDefault();
		const tag = newTag.trim();
		if (tag && !draft.tags.includes(tag)) draft.tags = [...draft.tags, tag];
		newTag = '';
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="bottom" class="mx-auto max-h-[90svh] max-w-xl overflow-y-auto rounded-t-2xl px-4 pb-[max(1rem,env(safe-area-inset-bottom))]">
		<Sheet.Header class="px-0">
			<Sheet.Title>{$t`Save session`}</Sheet.Title>
			<Sheet.Description>{$t`Everything below is optional — but it is what makes your analysis smart.`}</Sheet.Description>
		</Sheet.Header>
		<div class="grid gap-4">
			<div class="grid grid-cols-2 gap-2 rounded-lg bg-muted p-1 text-sm font-medium" role="radiogroup" aria-label={$t`Kind`}>
				{#each [['training', $t`Training`], ['competition', $t`Competition`]] as [k, label] (k)}
					<button
						type="button"
						role="radio"
						aria-checked={draft.kind === k}
						class="rounded-md py-2 {draft.kind === k ? 'bg-background shadow-sm' : 'text-muted-foreground'}"
						onclick={() => (draft.kind = k as Draft['kind'])}>{label}</button
					>
				{/each}
			</div>
			<FormField label={$t`Weapon`} hint={$t`Progress compares your results between weapons.`}>
				{#if weapons.length}
					<div class="flex flex-wrap gap-1.5" role="radiogroup" aria-label={$t`Weapon`}>
						{#each weapons as w (w.id)}
							<button
								type="button"
								role="radio"
								aria-checked={draft.weaponId === w.id}
								class="rounded-full border px-2.5 py-1 text-xs {draft.weaponId === w.id ? 'border-primary bg-primary text-primary-foreground' : 'hover:bg-muted'}"
								onclick={() => (draft.weaponId = draft.weaponId === w.id ? null : w.id)}>{w.name}</button
							>
						{/each}
					</div>
				{/if}
				<form class="flex gap-2" onsubmit={addWeapon}>
					<Input bind:value={newWeapon} placeholder={$t`New weapon, e.g. Steyr EVO 10`} maxlength={60} name="weapon" required />
					<Button type="submit" variant="outline" size="icon" aria-label={$t`Add weapon`}><Plus /></Button>
				</form>
				<FormError message={weaponError} />
			</FormField>
			<FormField label={$t`How did it feel?`}>
				<div class="flex justify-between gap-1" role="radiogroup" aria-label={$t`How did it feel?`}>
					{#each faces as Face, i (i)}
						<button
							type="button"
							role="radio"
							aria-checked={draft.rating === i + 1}
							aria-label={String(i + 1)}
							class="flex flex-1 justify-center rounded-lg border py-2.5 transition {draft.rating === i + 1 ? 'border-primary bg-primary/10 text-primary' : 'text-muted-foreground hover:text-foreground'}"
							onclick={() => (draft.rating = draft.rating === i + 1 ? null : i + 1)}><Face class="size-6" /></button
						>
					{/each}
				</div>
			</FormField>
			<FormField label={$t`Tags`} hint={$t`For example equipment, conditions or what you worked on. Progress compares your results by tag.`}>
				<div class="flex flex-wrap gap-1.5">
					{#each suggestions as tag (tag)}
						<button
							type="button"
							class="rounded-full border px-2.5 py-1 text-xs {draft.tags.includes(tag) ? 'border-primary bg-primary text-primary-foreground' : 'hover:bg-muted'}"
							aria-pressed={draft.tags.includes(tag)}
							onclick={() => toggle(tag)}>{tag}</button
						>
					{/each}
				</div>
				<form class="flex gap-2" onsubmit={addTag}>
					<Input bind:value={newTag} placeholder={$t`New tag`} maxlength={40} name="tag" />
					<Button type="submit" variant="outline" size="icon" aria-label={$t`Add tag`}><Plus /></Button>
				</form>
			</FormField>
			<div class="grid grid-cols-2 gap-3">
				<FormField label={$t`When`}><Input type="datetime-local" bind:value={draft.startedAt} name="startedAt" /></FormField>
				<FormField label={$t`Where`}><Input bind:value={draft.place} placeholder={$t`Range`} name="place" maxlength={200} /></FormField>
			</div>
			{#if draft.kind === 'competition'}
				<FormField label={$t`Competition`}><Input bind:value={draft.title} name="title" maxlength={200} /></FormField>
			{/if}
			<FormField label={$t`Notes`}><Textarea bind:value={draft.notes} rows={3} name="notes" placeholder={$t`What went well? What will you try next time?`} /></FormField>
			<Button size="lg" class="h-12 text-base" disabled={saving} onclick={onsave}>{$t`Save session`}</Button>
		</div>
	</Sheet.Content>
</Sheet.Root>
