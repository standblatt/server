<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import ArrowDown from '@lucide/svelte/icons/arrow-down';
	import ArrowUp from '@lucide/svelte/icons/arrow-up';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call, type AgeCategory, type Classification, type ClassificationDefinition } from '$lib/api';
	import { ageCategories } from '$lib/categories';
	import { ageCategoryLabel } from '$lib/labels';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import * as Sheet from '$lib/components/ui/sheet';
	import { Switch } from '$lib/components/ui/switch';
	import { ask } from '$lib/confirm.svelte';
	import { programmeTitle } from '$lib/format';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { keys, programmesQuery } from '$lib/tenant';
	import { has, type EventProps } from '..';

	let { tenantId, event }: EventProps = $props();
	const eventId = $derived(event.id);
	const queryClient = useQueryClient();
	const list = createQuery(() => ({
		queryKey: keys.classifications(tenantId, eventId),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/classifications', { params: { path: { tenantId, eventId } } }))
	}));
	const programmes = programmesQuery(
		() => tenantId,
		() => eventId
	);

	type TieBreak = ClassificationDefinition['tieBreak'][number];
	const criteria = $derived<[TieBreak, string][]>([
		['series', $t`Series from the last one backwards`],
		['shots', $t`More of the highest single shots`],
		['mouches', $t`More mouches`],
		['age', $t`Age (SSV rule)`]
	]);

	let editing = $state<Classification | 'new' | null>(null);
	let draft = $state<ClassificationDefinition>(blank());
	let error = $state('');
	let busy = $state(false);

	/** An empty classification with the SSV tie-breaks. */
	function blank(): ClassificationDefinition {
		return { name: '', parts: [], best: 0, complete: false, tieBreak: ['series', 'mouches', 'age'], classes: [], showOnKiosk: true };
	}

	/** Opens the editor with a copy of the classification, or an empty one. */
	function open(c: Classification | 'new') {
		error = '';
		editing = c;
		draft = c === 'new' ? blank() : structuredClone($state.snapshot(c));
	}

	/** Title of a programme of the event. */
	const titleOf = (id: number) => {
		const p = programmes.data?.find((x) => x.id === id);
		return p ? programmeTitle(p) : `#${id}`;
	};

	/** Moves a tie-break up or down in the order. */
	function moveTie(i: number, by: number) {
		const list = [...draft.tieBreak];
		[list[i], list[i + by]] = [list[i + by], list[i]];
		draft.tieBreak = list;
	}

	/** Turns a tie-break on (at the end) or off. */
	function toggleTie(c: TieBreak, on: boolean) {
		draft.tieBreak = on ? [...draft.tieBreak, c] : draft.tieBreak.filter((x) => x !== c);
	}

	/** Puts an age category into a class, taking it out of the others. */
	function toggleCategory(i: number, cat: AgeCategory, on: boolean) {
		draft.classes = draft.classes.map((cl, j) => ({
			...cl,
			categories: j === i ? (on ? [...cl.categories, cat] : cl.categories.filter((x) => x !== cat)) : cl.categories.filter((x) => !on || x !== cat)
		}));
	}

	/** Creates or saves the classification. */
	async function save(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		busy = true;
		const body = { ...draft, best: Number(draft.best) || 0, parts: draft.parts.map((p) => ({ ...p, count: Number(p.count) || 1 })) };
		try {
			if (editing === 'new') await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/classifications', { params: { path: { tenantId, eventId } }, body }));
			else if (editing) await call(api.PUT('/api/tenants/{tenantId}/events/{eventId}/classifications/{classificationId}', { params: { path: { tenantId, eventId, classificationId: editing.id } }, body }));
			await queryClient.invalidateQueries({ queryKey: keys.classifications(tenantId, eventId) });
			await queryClient.invalidateQueries({ queryKey: keys.classificationRankings(tenantId, eventId) });
			toast.success($t`Classification "${draft.name}" saved`);
			editing = null;
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}

	/** Deletes a classification after confirmation. */
	async function remove(c: Classification) {
		if (!(await ask({ title: $t`Delete classification?`, description: $t`"${c.name}" is deleted; the results stay.`, action: $t`Delete`, destructive: true }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/events/{eventId}/classifications/{classificationId}', { params: { path: { tenantId, eventId, classificationId: c.id } } }));
		await queryClient.invalidateQueries({ queryKey: keys.classifications(tenantId, eventId) });
		await queryClient.invalidateQueries({ queryKey: keys.classificationRankings(tenantId, eventId) });
	}
</script>

<section class="rounded-lg border bg-card">
	<div class="flex flex-wrap items-center gap-3 border-b px-4 py-3">
		<div class="min-w-48 flex-1">
			<h2 class="text-sm font-semibold">{$t`Classifications`}</h2>
			<p class="text-xs text-muted-foreground">{$t`Rankings over several programmes or attempts, e.g. a section or club championship.`}</p>
		</div>
		<Button size="sm" onclick={() => open('new')}><Plus />{$t`New classification`}</Button>
	</div>
	<ul class="divide-y">
		{#each list.data ?? [] as c (c.id)}
			<li class="flex items-center gap-3 px-4 py-2.5 text-sm">
				<div class="min-w-0 flex-1">
					<div class="font-medium">{c.name}</div>
					<div class="truncate text-xs text-muted-foreground">{c.parts.map((p) => (p.count > 1 ? `${titleOf(p.programmeId)} ×${p.count}` : titleOf(p.programmeId))).join(' + ')}</div>
				</div>
				<Button variant="ghost" size="icon-sm" aria-label={$t`Edit ${c.name}`} onclick={() => open(c)}><Pencil /></Button>
				<Button variant="ghost" size="icon-sm" aria-label={$t`Delete ${c.name}`} onclick={() => remove(c)}><Trash /></Button>
			</li>
		{:else}
			<li class="px-4 py-6 text-center text-sm text-muted-foreground">{$t`No classifications yet. Programme rankings need none.`}</li>
		{/each}
	</ul>
</section>

{#if editing}
	<Sheet.Root open onOpenChange={(o) => !o && (editing = null)}>
		<Sheet.Content class="gap-0 p-0 data-[side=right]:w-full data-[side=right]:sm:max-w-2xl">
			<Sheet.Header class="border-b px-6 py-4">
				<Sheet.Title>{editing === 'new' ? $t`New classification` : $t`Edit classification`}</Sheet.Title>
				<Sheet.Description>{$t`Which programmes count, how ties are decided and who is ranked together.`}</Sheet.Description>
			</Sheet.Header>
			<form class="flex min-h-0 flex-1 flex-col" onsubmit={save}>
				<div class="flex-1 space-y-6 overflow-y-auto px-6 py-5">
					<FormField label={$t`Name`}><Input bind:value={draft.name} required maxlength={100} /></FormField>

					<section class="space-y-2">
						<h3 class="text-sm font-semibold">{$t`Programmes`}</h3>
						{#each draft.parts as part, i (i)}
							<div class="flex items-end gap-2">
								<FormField label={$t`Programme ${i + 1}`} class="min-w-0 flex-1">
									<NativeSelect.Root class="w-full" bind:value={part.programmeId}>
										{#each programmes.data ?? [] as p (p.id)}<NativeSelect.Option value={p.id}>{programmeTitle(p)}</NativeSelect.Option>{/each}
									</NativeSelect.Root>
								</FormField>
								<FormField label={$t`Counting attempts`} class="w-32"><Input type="number" min="1" max="100" bind:value={part.count} /></FormField>
								<Button variant="ghost" size="icon" aria-label={$t`Remove programme ${i + 1}`} onclick={() => draft.parts.splice(i, 1)}><Trash /></Button>
							</div>
						{/each}
						<Button
							variant="outline"
							size="sm"
							disabled={!programmes.data?.length}
							onclick={() => draft.parts.push({ programmeId: (programmes.data ?? []).find((p) => !draft.parts.some((x) => x.programmeId === p.id))?.id ?? programmes.data![0].id, count: 1 })}
							><Plus />{$t`Add programme`}</Button
						>
						<p class="text-xs text-muted-foreground">{$t`Each programme counts the given number of its best (or first or last, as set on the programme) attempts.`}</p>
						{#if draft.parts.length > 1}
							<FormField label={$t`Count only the best programmes`} hint={$t`0 counts all of them.`} class="w-56">
								<Input type="number" min="0" max={draft.parts.length} bind:value={draft.best} />
							</FormField>
						{/if}
						<label class="flex items-center justify-between gap-3 text-sm">{$t`Rank only shooters with all counting results`} <Switch bind:checked={draft.complete} /></label>
					</section>

					<section class="space-y-2">
						<h3 class="text-sm font-semibold">{$t`Tie-break`}</h3>
						<p class="text-xs text-muted-foreground">{$t`Equal totals are decided by these criteria in order.`}</p>
						<ol class="divide-y rounded-md border">
							{#each draft.tieBreak as c, i (c)}
								<li class="flex items-center gap-2 px-3 py-1.5 text-sm">
									<span class="w-5 text-muted-foreground tabular-nums">{i + 1}.</span>
									<span class="flex-1">{criteria.find(([k]) => k === c)?.[1]}</span>
									<Button variant="ghost" size="icon-sm" aria-label={$t`Move up`} disabled={i === 0} onclick={() => moveTie(i, -1)}><ArrowUp /></Button>
									<Button variant="ghost" size="icon-sm" aria-label={$t`Move down`} disabled={i === draft.tieBreak.length - 1} onclick={() => moveTie(i, 1)}><ArrowDown /></Button>
								</li>
							{/each}
						</ol>
						<div class="flex flex-wrap gap-x-4 gap-y-1">
							{#each criteria as [k, label] (k)}
								<label class="flex items-center gap-2 text-sm"><Checkbox checked={draft.tieBreak.includes(k)} onCheckedChange={(v) => toggleTie(k, v)} />{label}</label>
							{/each}
						</div>
					</section>

					<section class="space-y-2">
						<h3 class="text-sm font-semibold">{$t`Classes`}</h3>
						<p class="text-xs text-muted-foreground">{$t`Without classes everybody is in one ranking. Shooters whose category is in no class are not listed.`}</p>
						{#each draft.classes as cl, i (i)}
							<fieldset class="space-y-2 rounded-md border p-3">
								<div class="flex items-end gap-2">
									<FormField label={$t`Class ${i + 1}`} class="flex-1"><Input bind:value={cl.name} required maxlength={100} /></FormField>
									<Button variant="ghost" size="icon" aria-label={$t`Remove class ${i + 1}`} onclick={() => draft.classes.splice(i, 1)}><Trash /></Button>
								</div>
								<div class="grid grid-cols-2 gap-1 sm:grid-cols-3">
									{#each ageCategories as cat (cat)}
										<label class="flex items-center gap-2 text-sm">
											<Checkbox checked={cl.categories.includes(cat)} onCheckedChange={(v) => toggleCategory(i, cat, v)} aria-label={$t`${$t(ageCategoryLabel(cat))} in class ${i + 1}`} />{$t(ageCategoryLabel(cat))}
										</label>
									{/each}
								</div>
							</fieldset>
						{/each}
						<Button variant="outline" size="sm" onclick={() => draft.classes.push({ name: '', categories: [] })}><Plus />{$t`Add class`}</Button>
					</section>

					{#if has(event, 'kiosk')}<label class="flex items-center justify-between gap-3 text-sm">{$t`Show on kiosk displays`} <Switch bind:checked={draft.showOnKiosk} /></label>{/if}
				</div>
				<div class="flex items-center gap-3 border-t px-6 py-3">
					<FormError message={error} class="flex-1" />
					<Button type="submit" class="ml-auto" disabled={busy}>{editing === 'new' ? $t`Create` : $t`Save`}</Button>
				</div>
			</form>
		</Sheet.Content>
	</Sheet.Root>
{/if}
