<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { untrack } from 'svelte';
	import type { AgeCategory, ProgrammeDefinition } from './api';
	import { ageCategories } from './categories';
	import { ageCategoryLabel } from './labels';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { Button } from './components/ui/button';
	import { Checkbox } from './components/ui/checkbox';
	import { Input } from './components/ui/input';
	import * as NativeSelect from './components/ui/native-select';
	import { Switch } from './components/ui/switch';
	import * as Table from './components/ui/table';
	import { Textarea } from './components/ui/textarea';
	import FormError from './FormError.svelte';
	import FormField from './FormField.svelte';

	let {
		initial,
		others = [],
		locked = false,
		submitLabel,
		onsubmit
	}: {
		initial?: ProgrammeDefinition;
		others?: { name: string; variant?: string }[];
		/** Shooters are registered: series are fixed, age categories and awards can only be added. */
		locked?: boolean;
		submitLabel: string;
		onsubmit: (d: ProgrammeDefinition) => Promise<void>;
	} = $props();

	type Award = { name: string; kind: 'natural' | 'card' | 'cash'; value: number; limits: Record<AgeCategory, string> };

	const start: ProgrammeDefinition = untrack(
		() =>
			initial ?? {
				name: '',
				variant: '',
				showOnRankings: true,
				series: [{ shots: 10, maxPoints: 100 }],
				prices: [],
				awards: []
			}
	);

	/** Award limits with an empty value for every age category. */
	function emptyLimits(): Record<AgeCategory, string> {
		return Object.fromEntries(ageCategories.map((c) => [c, ''])) as Record<AgeCategory, string>;
	}

	let name = $state(start.name);
	let variant = $state(start.variant ?? '');
	const names = $derived([...new Set(others.map((o) => o.name))].sort());
	const siblings = $derived(others.filter((o) => o.name === name.trim()));
	let description = $state(start.description ?? '');
	let barcode = $state(start.barcode ?? '');
	let countsForClubCompetition = $state(start.countsForClubCompetition ?? false);
	let showOnRankings = $state(start.showOnRankings ?? true);
	let series = $state(start.series.map((s) => ({ shots: s.shots, maxPoints: s.maxPoints, name: s.name ?? '', trackMouches: s.trackMouches ?? false, trackShots: s.trackShots ?? false })));
	let prices = $state(
		Object.fromEntries(
			ageCategories.map((c) => {
				const p = start.prices.find((x) => x.ageCategory === c);
				return [c, { enabled: !!p, price: p?.price ?? 0 }];
			})
		) as Record<AgeCategory, { enabled: boolean; price: number }>
	);
	let awards = $state<Award[]>(
		start.awards.map((a) => {
			const limits = emptyLimits();
			for (const l of a.limits) limits[l.ageCategory] = String(l.points);
			return { name: a.name, kind: a.kind, value: a.value, limits };
		})
	);
	let error = $state('');
	let busy = $state(false);
	const fixedPrices = untrack(() => new Set(start.prices.map((p) => p.ageCategory)));
	const fixedAwards = untrack(() => start.awards.length);

	const totalShots = $derived(series.reduce((n, s) => n + (Number(s.shots) || 0), 0));
	const totalPoints = $derived(series.reduce((n, s) => n + (Number(s.maxPoints) || 0), 0));

	/** Appends a series that copies the shots, maximum and mouches setting of the last one. */
	function addSeries() {
		const last = series[series.length - 1];
		series.push({ shots: last?.shots ?? 10, maxPoints: last?.maxPoints ?? 100, name: '', trackMouches: last?.trackMouches ?? false, trackShots: last?.trackShots ?? false });
	}

	/** Appends an empty award tier. */
	function addAward() {
		awards.push({ name: '', kind: 'card', value: 0, limits: emptyLimits() });
	}

	/** Builds the ProgrammeDefinition from the form and passes it to onsubmit. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		busy = true;
		error = '';
		try {
			await onsubmit({
				name,
				variant,
				description,
				barcode,
				countsForClubCompetition,
				showOnRankings,
				series: series.map((s) => ({ shots: Number(s.shots), maxPoints: Number(s.maxPoints), name: s.name, trackMouches: s.trackMouches, trackShots: s.trackShots })),
				prices: ageCategories.filter((c) => prices[c].enabled).map((c) => ({ ageCategory: c, price: Number(prices[c].price) || 0 })),
				awards: awards.map((a) => ({
					name: a.name,
					kind: a.kind,
					value: Number(a.value) || 0,
					limits: ageCategories.filter((c) => a.limits[c] !== '').map((c) => ({ ageCategory: c, points: Number(a.limits[c]) }))
				}))
			});
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

{#snippet section(title: string, hint?: string)}
	<div class="mb-3">
		<h3 class="text-sm font-semibold">{title}</h3>
		{#if hint}<p class="text-xs text-muted-foreground">{hint}</p>{/if}
	</div>
{/snippet}

<form class="flex min-h-0 flex-1 flex-col" onsubmit={submit}>
	<div class="flex-1 space-y-8 overflow-y-auto px-6 py-5">
		<section class="grid gap-4 sm:grid-cols-2">
			<FormField label={$t`Name`}>
				<Input bind:value={name} required list="programme-names" />
				<datalist id="programme-names">{#each names as n (n)}<option value={n}></option>{/each}</datalist>
			</FormField>
			<FormField label={$t`Variant`} hint={siblings.length ? '' : $t`For a variant, pick the name of an existing programme: a shooter registers for only one programme of the same name.`}>
				<Input bind:value={variant} placeholder={$t`e.g. Stgw 90`} />
			</FormField>
			{#if siblings.length}
				<p class="rounded-md border border-sky-500/40 bg-sky-500/10 px-3 py-2 text-sm sm:col-span-2" data-testid="variant-hint">
					{$t`Variant of "${name.trim()}": shooters register for only one of ${siblings.length + 1} programmes.`}
					{#if !variant.trim() || siblings.some((o) => (o.variant ?? '') === variant.trim())}<strong class="text-destructive">{$t`Give each variant its own variant name.`}</strong>{/if}
				</p>
			{/if}
			<FormField label={$t`Description`} class="sm:col-span-2"><Textarea rows={2} bind:value={description} /></FormField>
			<FormField label={$t`Barcode`}><Input bind:value={barcode} /></FormField>
			<div class="grid content-end gap-2">
				<label class="flex items-center justify-between gap-3 text-sm">{$t`Counts for the club competition`} <Switch bind:checked={countsForClubCompetition} /></label>
				<label class="flex items-center justify-between gap-3 text-sm">{$t`Show on kiosk displays`} <Switch bind:checked={showOnRankings} /></label>
			</div>
		</section>

		{#if locked}
			<p class="rounded-md border border-amber-500/40 bg-amber-500/10 px-3 py-2 text-sm" data-testid="locked-hint">
				{$t`Shooters are registered for this programme: the series stay as they are, and age categories and awards can be added but not removed.`}
			</p>
		{/if}

		<section>
			<h3 class="mb-3 text-sm font-semibold">{$t`Series (${totalShots} shots, max. ${totalPoints} points)`}</h3>
			<div class="overflow-hidden rounded-md border">
				<Table.Root>
					<Table.Header>
						<Table.Row><Table.Head class="hidden w-10 sm:table-cell">#</Table.Head><Table.Head>{$t`Shots`}</Table.Head><Table.Head class="whitespace-normal">{$t`Max. points`}</Table.Head><Table.Head>{$t`Name`}</Table.Head><Table.Head><span class="sm:hidden" title={$t`Mouches`}>{$t({ message: "M", comment: "Column header: mouches" })}</span><span class="hidden sm:inline">{$t`Mouches`}</span></Table.Head><Table.Head><span class="sm:hidden" title={$t`Single shots`}>{$t({ message: "SH", comment: "Column header: single shots" })}</span><span class="hidden sm:inline">{$t`Single shots`}</span></Table.Head><Table.Head></Table.Head></Table.Row>
					</Table.Header>
					<Table.Body>
						{#each series as s, i (i)}
							<Table.Row>
								<Table.Cell class="hidden text-muted-foreground sm:table-cell">{i + 1}</Table.Cell>
								<Table.Cell><Input class="h-8 w-14 sm:w-20" type="number" min="1" max="100" bind:value={s.shots} required disabled={locked} /></Table.Cell>
								<Table.Cell><Input class="h-8 w-16 sm:w-24" type="number" min="1" bind:value={s.maxPoints} required disabled={locked} /></Table.Cell>
								<Table.Cell class="w-full"><Input class="h-8 min-w-16" bind:value={s.name} /></Table.Cell>
								<Table.Cell><Checkbox bind:checked={s.trackMouches} disabled={locked} aria-label={$t`Track mouches in series ${i + 1}`} /></Table.Cell>
								<Table.Cell><Checkbox bind:checked={s.trackShots} disabled={locked} aria-label={$t`Record single shots in series ${i + 1}`} /></Table.Cell>
								<Table.Cell class="text-right">
									<Button variant="ghost" size="icon-sm" title={$t`Remove series`} aria-label={$t`Remove series ${i + 1}`} onclick={() => series.splice(i, 1)} disabled={locked || series.length === 1}><Trash /></Button>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
			<Button variant="outline" size="sm" class="mt-2" onclick={addSeries} disabled={locked}><Plus />{$t`Add series`}</Button>
		</section>

		<section>
			{@render section($t`Prices (CHF)`, $t`Only the checked age categories may shoot this programme.`)}
			<div class="grid grid-cols-2 gap-3 *:min-w-0 sm:grid-cols-3 lg:grid-cols-5">
				{#each ageCategories as c (c)}
					<div class="grid min-w-0 gap-1.5 rounded-md border p-2.5 {prices[c].enabled ? '' : 'bg-muted/40'}">
						<label class="flex items-center gap-2 text-sm font-medium"><Checkbox bind:checked={prices[c].enabled} disabled={locked && fixedPrices.has(c)} /><span class="min-w-0 break-words hyphens-auto">{$t(ageCategoryLabel(c))}</span></label>
						<Input class="h-8" type="number" min="0" step="0.05" bind:value={prices[c].price} disabled={!prices[c].enabled} aria-label={$t`Price ${$t(ageCategoryLabel(c))}`} />
					</div>
				{/each}
			</div>
		</section>

		<section>
			{@render section($t`Awards`, $t`Minimum total per age category; leave empty where the award is not available.`)}
			<div class="space-y-3">
				{#each awards as a, i (i)}
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
							<Button variant="ghost" size="icon" title={$t`Remove award`} aria-label={$t`Remove award ${i + 1}`} onclick={() => awards.splice(i, 1)} disabled={locked && i < fixedAwards}><Trash /></Button>
						</div>
						<div class="grid grid-cols-2 gap-3 *:min-w-0 sm:grid-cols-3 lg:grid-cols-5">
							{#each ageCategories as c (c)}
								<FormField label={$t(ageCategoryLabel(c))}><Input class="h-8" type="number" min="0" max={totalPoints} step="0.1" bind:value={a.limits[c]} /></FormField>
							{/each}
						</div>
					</fieldset>
				{/each}
			</div>
			<Button variant="outline" size="sm" class="mt-2" onclick={addAward}><Plus />{$t`Add award`}</Button>
		</section>
	</div>

	<div class="flex items-center gap-3 border-t px-6 py-3">
		<FormError message={error} class="flex-1" />
		<Button type="submit" class="ml-auto" disabled={busy}>{submitLabel}</Button>
	</div>
</form>
