<script lang="ts">
	import { Trans, t } from 'lingui-for-svelte/macro';
	import { untrack } from 'svelte';
	import type { EntryResult, SeriesResult } from './api';
	import { Badge } from './components/ui/badge';
	import { Button } from './components/ui/button';
	import { Input } from './components/ui/input';
	import PrintButton from './PrintButton.svelte';

	let {
		entry,
		tenantId,
		readonly = false,
		print = true,
		printForm = print,
		onsave,
		onaward
	}: {
		entry: EntryResult;
		tenantId: number;
		readonly?: boolean;
		/** Show the wreath card print button (not for helper links). */
		print?: boolean;
		/** Show the print form button when the programme has a form. */
		printForm?: boolean;
		onsave: (series: SeriesResult[]) => Promise<void>;
		onaward: (tierId: number | undefined) => Promise<void>;
	} = $props();

	type Row = { points: string; mouches: string; shots: string[] };
	/** Form rows (points, mouches and single shots as strings) from the saved results of an entry. */
	function rowsFrom(e: EntryResult): Row[] {
		return e.series.map((s, i) => ({
			points: e.results[i] ? String(e.results[i].points) : '',
			mouches: e.results[i]?.mouches != null ? String(e.results[i].mouches) : '',
			shots: s.trackShots ? Array.from({ length: s.shots }, (_, j) => (e.results[i]?.shots?.[j] != null ? String(e.results[i].shots![j]) : '')) : []
		}));
	}

	/** Whether series i has any single shot typed; its points are then the sum of the shots. */
	const hasShots = (i: number) => rows[i].shots.some((v) => v !== '');
	/** The points of series i: the sum of its shots when they are typed, else the typed total. */
	const points = (i: number) => (hasShots(i) ? Math.round(rows[i].shots.reduce((n, v) => n + (Number(v) || 0), 0) * 100) / 100 : Number(rows[i].points) || 0);

	let rows = $state<Row[]>(untrack(() => rowsFrom(entry)));
	let error = $state('');
	let busy = $state(false);
	let tier = $state(untrack(() => entry.award?.id ?? entry.eligibleAwards[0]?.id ?? 0));
	let form = $state<HTMLFormElement>();

	$effect(() => {
		const e = entry;
		untrack(() => {
			rows = rowsFrom(e);
			tier = e.award?.id ?? e.eligibleAwards[0]?.id ?? 0;
		});
	});

	const readOnly = $derived(entry.combinedSeries.length > 0);
	const total = $derived(rows.reduce((n, _, i) => n + points(i), 0));
	/** Whether the points of series i exceed the series maximum. */
	const tooHigh = (i: number) => points(i) > entry.series[i].maxPoints;

	/** Moves focus to the next score input (or the save button) once the typed value cannot grow without exceeding max. */
	function advance(el: HTMLInputElement, value: string, max: number) {
		if (!/^\d+$/.test(value) || Number(value) * 10 <= max) return;
		const inputs = [...(form?.querySelectorAll<HTMLInputElement>('input[data-score]') ?? [])];
		const next = inputs[inputs.indexOf(el) + 1];
		if (next) next.focus();
		else form?.querySelector<HTMLButtonElement>('button[type=submit]')?.focus();
	}

	/** Validates that every series has points and saves the results through onsave. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		const partial = rows.findIndex((r) => hasShots(rows.indexOf(r)) && r.shots.some((v) => v === ''));
		if (partial >= 0) {
			error = $t`Series ${partial + 1}: enter every shot or none`;
			return;
		}
		if (rows.some((r, i) => r.points === '' && !hasShots(i))) {
			error = $t`Enter all series`;
			return;
		}
		busy = true;
		try {
			await onsave(
				rows.map((r, i) => ({
					points: points(i),
					mouches: entry.series[i].trackMouches && r.mouches !== '' ? Number(r.mouches) : undefined,
					shots: hasShots(i) ? r.shots.map(Number) : undefined
				}))
			);
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}

	/** Saves the chosen award tier (none when 0) through onaward. */
	async function saveAward() {
		error = '';
		try {
			await onaward(tier || undefined);
		} catch (err) {
			error = (err as Error).message;
		}
	}
</script>

<section class="rounded-lg border bg-card" aria-label={entry.repeatable ? $t`${entry.title}, attempt ${entry.attempt}` : entry.title}>
	<div class="flex flex-wrap items-center gap-3 border-b px-4 py-2.5">
		<h3 class="font-semibold">{entry.title}{#if entry.repeatable}<span class="font-normal text-muted-foreground"> · {$t`Attempt ${entry.attempt}`}</span>{/if}</h3>
		<Badge variant={entry.status === 'completed' ? 'default' : 'secondary'} class={entry.status === 'completed' ? 'bg-emerald-600 text-white' : ''}>
			{entry.status === 'completed' ? $t`completed` : $t`pending`}
		</Badge>
		<span class="ml-auto text-lg tabular-nums"><Trans>Total <strong>{readOnly ? entry.total : total}</strong> / {entry.maxPoints}</Trans></span>
	</div>

	<div class="p-4">
		{#if readOnly}
			<p class="text-sm text-muted-foreground">{$t`Takes series ${entry.combinedSeries.join(', ')} of ${entry.combinedFrom ?? ''}.`}</p>
			{#if entry.results.length > 0}
				<p class="mt-1 tabular-nums">{entry.results.map((r, i) => `S${i + 1}: ${r.points}`).join(' · ')}</p>
			{:else}
				<p class="mt-1 text-sm text-muted-foreground">{$t`Waiting for the results of ${entry.combinedFrom ?? ''}.`}</p>
			{/if}
		{:else}
			<form bind:this={form} class="flex flex-wrap items-start gap-3" onsubmit={submit}>
				<fieldset class="contents" disabled={readonly}>
					{#each entry.series as s, i (i)}
						<div class="grid gap-1.5">
							<label class="grid gap-1.5 text-sm">
								<span class="font-medium">{s.name || `S${i + 1}`} <span class="font-normal text-muted-foreground">/{s.maxPoints}</span></span>
								{#if hasShots(i)}
									<Input class="h-12 w-20 text-center text-xl tabular-nums md:text-xl {tooHigh(i) ? 'border-destructive ring-3 ring-destructive/20' : ''}" aria-label={$t`Series ${i + 1}`} value={points(i)} readonly tabindex={-1} />
								{:else}
									<Input
										class="h-12 w-20 text-center text-xl tabular-nums md:text-xl {tooHigh(i) ? 'border-destructive ring-3 ring-destructive/20' : ''}"
										data-score={s.trackShots ? undefined : ''}
										inputmode="decimal"
										aria-label={$t`Series ${i + 1}`}
										bind:value={rows[i].points}
										oninput={(e) => advance(e.currentTarget, rows[i].points, s.maxPoints)}
									/>
								{/if}
							</label>
							{#if s.trackShots}
								<div class="flex max-w-[34rem] flex-wrap gap-1" role="group" aria-label={$t`Shots of series ${i + 1}`}>
									{#each rows[i].shots as _, j (j)}
										<Input
											class="h-10 w-12 px-1 text-center tabular-nums"
											data-score
											inputmode="decimal"
											aria-label={$t`Series ${i + 1}, shot ${j + 1}`}
											bind:value={rows[i].shots[j]}
											oninput={(e) => advance(e.currentTarget, rows[i].shots[j], s.maxPoints / s.shots)}
										/>
									{/each}
								</div>
							{/if}
							{#if s.trackMouches}
								<label class="grid gap-1 text-xs text-muted-foreground">
									{$t`Mouches`}
									<Input
										class="h-9 w-20 text-center tabular-nums"
										data-score
										inputmode="numeric"
										aria-label={$t`Mouches series ${i + 1}`}
										bind:value={rows[i].mouches}
										oninput={(e) => advance(e.currentTarget, rows[i].mouches, s.shots)}
									/>
								</label>
							{/if}
						</div>
					{/each}
					<Button type="submit" size="lg" class="mt-6 h-12 px-5" disabled={busy || rows.some((_, i) => tooHigh(i))}>{$t`Save`}</Button>
				</fieldset>
			</form>
		{/if}

		{#if entry.status === 'completed'}
			{#if printForm && entry.form}<div class="mt-3"><PrintButton {tenantId} target={{ document: 'form', entryId: entry.entryId }} /></div>{/if}
			<div class="mt-4 border-t pt-3">
				{#if entry.eligibleAwards.length === 0}
					<p class="text-sm text-muted-foreground">{$t`No award reached.`}</p>
				{:else}
					<fieldset class="flex flex-wrap items-center gap-x-4 gap-y-2">
						<legend class="mb-2 text-sm font-medium">{entry.award ? $t`Award: ${entry.award.name}` : $t`Award (choose one)`}</legend>
						{#each entry.eligibleAwards as a (a.id)}
							<label class="flex items-center gap-2 text-sm">
								<input type="radio" class="size-4 accent-primary" name="award-{entry.entryId}" value={a.id} bind:group={tier} disabled={readonly} onchange={saveAward} />
								{a.name}{a.value ? ` (CHF ${a.value.toFixed(2)})` : ''}
							</label>
						{/each}
						{#if !readonly && !entry.award}<Button size="sm" onclick={saveAward}>{$t`Save award`}</Button>{/if}
						{#if entry.award?.kind === 'card'}
							{#if print}<PrintButton {tenantId} target={{ document: 'wreath_card', entryId: entry.entryId }} />{/if}
						{/if}
					</fieldset>
				{/if}
			</div>
		{/if}
		{#if error}<p class="mt-3 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive" role="alert">{error}</p>{/if}
	</div>
</section>
