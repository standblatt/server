<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import Minus from '@lucide/svelte/icons/minus';
	import Plus from '@lucide/svelte/icons/plus';
	import { untrack, type Snippet } from 'svelte';
	import type { ClubRef, LaneState, Programme, ProgrammeChoice, RangeLanes, Registration, Shooter } from './api';
	import { ageCategoryFor } from './categories';
	import { ageCategoryLabel } from './labels';
	import { Badge } from './components/ui/badge';
	import { Button } from './components/ui/button';
	import { Checkbox } from './components/ui/checkbox';
	import * as Dialog from './components/ui/dialog';
	import * as NativeSelect from './components/ui/native-select';
	import { ask } from './confirm.svelte';
	import { formatDate, programmeTitle, withRange } from './format';

	let {
		shooter,
		eventDate,
		programmes,
		lanes,
		initial,
		initialClubId,
		used,
		readonly = false,
		lanesOnly = false,
		submitLabel,
		onsubmit,
		footer
	}: {
		shooter: Pick<Shooter, 'firstname' | 'lastname' | 'licenseNumber' | 'birthday'> & { clubs?: ClubRef[] };
		eventDate: string;
		programmes: Programme[];
		lanes: RangeLanes[];
		initial?: Registration;
		initialClubId?: number;
		/** Attempts the shooter already has per programme in other registrations of the event. */
		used?: Map<number, number>;
		readonly?: boolean;
		/** Only the lane may change (ammunition desk): programmes stay as they are. */
		lanesOnly?: boolean;
		submitLabel: string;
		onsubmit: (r: { clubId?: number; rangeId?: number; lane?: number; programmes: ProgrammeChoice[] }) => Promise<void>;
		footer?: Snippet;
	} = $props();

	const category = $derived(initial?.ageCategory ?? ageCategoryFor(shooter.birthday, eventDate));
	const kept = untrack(() => new Set(initial?.entries.map((e) => e.programmeId) ?? []));

	let clubId = $state(untrack(() => initial?.clubId ?? initialClubId ?? shooter.clubs?.[0]?.id ?? 0));
	let lane = $state(untrack(() => initial?.lane ?? 0));
	let rangeId = $state(untrack(() => initial?.rangeId ?? (lanes.length === 1 ? lanes[0].rangeId : 0)));
	let selected = $state<Record<number, boolean>>(untrack(() => Object.fromEntries((initial?.entries ?? []).map((e) => [e.programmeId, true]))));
	let combine = $state<Record<number, boolean>>(
		untrack(() => Object.fromEntries((initial?.entries ?? []).filter((e) => e.combinedSeries.length > 0).map((e) => [e.programmeId, true])))
	);
	let counts = $state<Record<number, number>>(
		untrack(() => {
			const out: Record<number, number> = {};
			for (const e of initial?.entries ?? []) out[e.programmeId] = (out[e.programmeId] ?? 0) + 1;
			return out;
		})
	);
	let series = $state<Record<number, number[]>>(untrack(() => Object.fromEntries((initial?.entries ?? []).map((e) => [e.programmeId, e.combinedSeries]))));
	let error = $state('');
	let busy = $state(false);

	type Group = { name: string; items: Programme[] };
	const groups = $derived.by(() => {
		const out: Group[] = [];
		for (const p of programmes) {
			if (!p.active && !kept.has(p.id)) continue;
			const g = out.find((x) => x.name === p.name);
			if (g) g.items.push(p);
			else out.push({ name: p.name, items: [p] });
		}
		return out;
	});

	/** Price of a programme for the shooter's age category, undefined when not offered. */
	const price = (p: Programme) => p.prices.find((x) => x.ageCategory === category)?.price;
	/** Whether a shooter may shoot the programme more than once. */
	const repeatable = (p: Programme) => !!p.unlimitedAttempts || (p.attempts ?? 1) > 1;
	/** Attempts the shooter may still register for the programme. */
	const left = (p: Programme) => (p.unlimitedAttempts ? Infinity : Math.max(p.attempts ?? 1, 1)) - (used?.get(p.id) ?? 0);
	/** Attempts of the programme in this registration. */
	const count = (p: Programme) => Math.max(counts[p.id] ?? 1, 1);
	/** Price of n attempts after the shooter's earlier ones: further attempts cost the repeat price, if any. */
	function cost(p: Programme, n: number) {
		const first = price(p) ?? 0;
		const repeat = p.prices.find((x) => x.ageCategory === category)?.repeatPrice ?? first;
		let sum = 0;
		for (let i = 0; i < n; i++) sum += (used?.get(p.id) ?? 0) + i === 0 ? first : repeat;
		return sum;
	}
	/** Whether the programme can be chosen: offered for the category, attempts left, not lanes-only. */
	const allowed = (p: Programme) => !lanesOnly && price(p) !== undefined && left(p) > 0;
	/** Reason a programme cannot be chosen. */
	const unavailable = (p: Programme) => (left(p) > 0 ? $t`not for this age category` : repeatable(p) ? $t`no attempts left` : $t`already registered`);
	/** Programme by ID. */
	const byId = (id: number) => programmes.find((p) => p.id === id);
	const total = $derived(programmes.filter((p) => selected[p.id]).reduce((n, p) => n + cost(p, repeatable(p) ? count(p) : 1), 0));
	/** Range the chosen programmes are tied to, if any; the lane picker then shows only that range. */
	const requiredRange = $derived(programmes.find((p) => selected[p.id] && p.rangeId)?.rangeId);
	const shownLanes = $derived(requiredRange ? lanes.filter((r) => r.rangeId === requiredRange) : lanes);
	/** Name of a range from the lane states. */
	const rangeTitle = (id: number) => lanes.find((r) => r.rangeId === id)?.rangeName ?? '';
	const laneLabel = $derived(
		lane
			? withRange(
					lanes.map((r) => ({ id: r.rangeId, name: r.rangeName })),
					rangeId,
					$t`Lane ${lane}`
				)
			: $t`No lane: the shooter waits.`
	);
	$effect(() => {
		if (requiredRange && lane && rangeId !== requiredRange) lane = 0;
	});

	/** Selects one variant of a programme group and deselects the others. */
	function choose(group: Group, id: number) {
		for (const p of group.items) selected[p.id] = p.id === id;
	}

	/** Selects a lane (asking first when it is reserved) or clears it when clicked again. */
	async function pickLane(r: RangeLanes, l: LaneState) {
		if (lane === l.lane && rangeId === r.rangeId) {
			lane = 0;
			return;
		}
		if (l.state === 'reserved' && !(await ask({ title: $t`Lane ${l.lane} is reserved`, description: $t`Reserved for "${l.reservation ?? ''}". Use it anyway?`, action: $t`Use lane` }))) return;
		rangeId = r.rangeId;
		lane = l.lane;
	}

	/** Toggles series n as a combined series of the programme, up to max series. */
	function toggleSeries(id: number, n: number, max: number) {
		const list = series[id] ?? [];
		if (list.includes(n)) series[id] = list.filter((x) => x !== n);
		else if (list.length < max) series[id] = [...list, n].sort((a, b) => a - b);
	}

	/** CSS classes of a lane button for its state and the current selection. */
	function laneClass(r: RangeLanes, l: LaneState) {
		if (lane === l.lane && rangeId === r.rangeId) return 'border-amber-500 bg-amber-400 text-brand-950 ring-2 ring-amber-400/40';
		if (l.state === 'disabled') return 'lane-disabled';
		if (l.state === 'occupied' && l.registrationId !== initial?.id) return 'lane-occupied opacity-50';
		if (l.state === 'reserved') return 'lane-reserved hover:border-amber-400';
		return 'lane-free hover:border-emerald-400';
	}

	/** Selects combination targets of combined programmes and passes club, lane and programmes to onsubmit. */
	async function submit(e: SubmitEvent) {
		e.preventDefault();
		busy = true;
		error = '';
		try {
			const chosen = programmes.filter((p) => selected[p.id]);
			for (const p of chosen) {
				const target = p.combinableWithId ? byId(p.combinableWithId) : undefined;
				if (combine[p.id] && target && !selected[target.id]) selected[target.id] = true;
			}
			await onsubmit({
				clubId: clubId || undefined,
				rangeId: lane ? rangeId || undefined : undefined,
				lane: lane || undefined,
				programmes: programmes
					.filter((p) => selected[p.id])
					.map((p) => ({ programmeId: p.id, combinedSeries: combine[p.id] ? (series[p.id] ?? []) : [], count: repeatable(p) ? count(p) : undefined }))
			});
		} catch (err) {
			error = (err as Error).message;
		} finally {
			busy = false;
		}
	}
</script>

<form class="flex min-h-0 flex-1 flex-col" onsubmit={submit}>
	<fieldset class="contents" disabled={readonly}>
		<div class="grid min-h-0 flex-1 gap-6 overflow-y-auto px-6 py-4 md:grid-cols-[minmax(0,1fr)_minmax(0,1.2fr)]">
			<section class="space-y-5">
				<div class="flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
					<span>{$t`License ${shooter.licenseNumber || '–'}`}</span>
					<span>·</span>
					<span>{$t`born ${shooter.birthday ? formatDate(shooter.birthday) : '–'}`}</span>
					<Badge variant="secondary" data-testid="category">{$t(ageCategoryLabel(category))}</Badge>
				</div>

				{#if shooter.clubs}
					<label class="grid gap-1.5 text-sm">
						<span class="font-medium">{$t`Club`}</span>
						<NativeSelect.Root class="w-full" bind:value={clubId}>
							<NativeSelect.Option value={0}>{$t`none`}</NativeSelect.Option>
							{#each shooter.clubs as c (c.id)}<NativeSelect.Option value={c.id}>{c.name}</NativeSelect.Option>{/each}
						</NativeSelect.Root>
					</label>
				{/if}

				{#if lanes.length}
				<div>
					<div class="mb-2 flex items-baseline justify-between">
						<h3 class="text-sm font-medium">{$t`Lane`}</h3>
						<span class="text-xs text-muted-foreground">{laneLabel}</span>
					</div>
					{#each shownLanes as r (r.rangeId)}
						{#if lanes.length > 1}<h4 class="mt-2 mb-1 text-xs font-medium text-muted-foreground">{r.rangeName}</h4>{/if}
						<div class="grid grid-cols-[repeat(auto-fill,minmax(2.75rem,1fr))] gap-1.5" role="group" aria-label={r.rangeName}>
							{#each r.lanes as l (l.lane)}
								<button
									type="button"
									class="h-11 rounded-md border text-sm font-semibold transition-colors disabled:cursor-not-allowed {laneClass(r, l)}"
									disabled={l.state === 'disabled' || (l.state === 'occupied' && l.registrationId !== initial?.id)}
									title={l.state === 'occupied' ? l.shooterName : l.state === 'reserved' ? $t`reserved: ${l.reservation ?? ''}` : l.state === 'disabled' ? $t`out of service` : $t`free`}
									aria-pressed={lane === l.lane && rangeId === r.rangeId}
									onclick={() => pickLane(r, l)}>{l.lane}</button
								>
							{/each}
						</div>
					{/each}
					{#if requiredRange && lanes.length > 1}
						<p class="mt-2 text-xs text-muted-foreground">{$t`The chosen programmes are shot on ${rangeTitle(requiredRange ?? 0)}.`}</p>
					{/if}
				</div>
				{/if}
			</section>

			<section>
				<h3 class="mb-2 text-sm font-medium">{$t`Programmes`}</h3>
				<div class="divide-y rounded-md border">
					{#each groups as g (g.name)}
						<div class="p-3">
							{#if g.items.length > 1}
								<div class="mb-1 flex items-center justify-between">
									<span class="text-sm font-medium">{g.name}</span>
									<Button variant="ghost" size="xs" onclick={() => choose(g, 0)}>{$t`clear`}</Button>
								</div>
								<div class="grid gap-1" role="radiogroup" aria-label={g.name}>
									{#each g.items as p (p.id)}
										<label class="flex items-center gap-2.5 rounded-sm px-1 py-1 text-sm hover:bg-muted/60" class:opacity-50={!allowed(p)}>
											<input
												type="radio"
												class="size-4 accent-primary"
												name="group-{g.name}"
												checked={!!selected[p.id]}
												disabled={!allowed(p)}
												onchange={() => choose(g, p.id)}
											/>
											<span>{p.variant}</span>
											<span class="ml-auto text-muted-foreground tabular-nums">{allowed(p) ? `CHF ${price(p)?.toFixed(2)}` : unavailable(p)}</span>
										</label>
									{/each}
								</div>
							{:else}
								{@const p = g.items[0]}
								<label class="flex items-center gap-2.5 text-sm" class:opacity-50={!allowed(p)}>
									<Checkbox checked={!!selected[p.id]} onCheckedChange={(v) => (selected[p.id] = v)} disabled={!allowed(p)} />
									<span class="font-medium">{programmeTitle(p)}</span>
									<span class="ml-auto text-muted-foreground tabular-nums">{allowed(p) ? `CHF ${price(p)?.toFixed(2)}` : unavailable(p)}</span>
								</label>
							{/if}
							{#each g.items.filter((p) => selected[p.id] && repeatable(p)) as p (p.id)}
								<div class="mt-2 ml-6 flex items-center gap-2 text-sm">
									<span class="text-muted-foreground">{$t`Attempts`}</span>
									<Button variant="outline" size="icon-sm" aria-label={$t`Fewer attempts of ${programmeTitle(p)}`} disabled={count(p) <= 1} onclick={() => (counts[p.id] = count(p) - 1)}><Minus /></Button>
									<span class="w-6 text-center font-medium tabular-nums" data-testid="attempts">{count(p)}</span>
									<Button variant="outline" size="icon-sm" aria-label={$t`More attempts of ${programmeTitle(p)}`} disabled={count(p) >= left(p)} onclick={() => (counts[p.id] = count(p) + 1)}><Plus /></Button>
								</div>
							{/each}
							{#each g.items.filter((p) => selected[p.id] && p.combinableWithId) as p (p.id)}
								{@const target = byId(p.combinableWithId!)}
								{#if target}
									<div class="mt-2 ml-6 space-y-2 rounded-md bg-muted/50 p-2">
										<label class="flex items-center gap-2 text-sm">
											<Checkbox checked={!!combine[p.id]} onCheckedChange={(v) => (combine[p.id] = v)} />
											{$t`take ${p.series.length} series from ${programmeTitle(target)}`}
										</label>
										{#if combine[p.id]}
											<div class="flex flex-wrap gap-1.5" role="group" aria-label={$t`Series of ${programmeTitle(target)}`}>
												{#each target.series as _, i (i)}
													{@const picked = (series[p.id] ?? []).includes(i + 1)}
													<button
														type="button"
														class="size-8 rounded-md border text-sm font-medium {picked ? 'border-primary bg-primary text-primary-foreground' : 'bg-background hover:bg-muted'}"
														aria-pressed={picked}
														disabled={!picked && (series[p.id] ?? []).length >= p.series.length}
														onclick={() => toggleSeries(p.id, i + 1, p.series.length)}>{i + 1}</button
													>
												{/each}
											</div>
										{/if}
									</div>
								{/if}
							{/each}
						</div>
					{:else}
						<p class="p-3 text-sm text-muted-foreground">{$t`This event has no active programmes.`}</p>
					{/each}
				</div>
			</section>
		</div>

		{#if error}<p class="mx-6 mb-2 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive" role="alert">{error}</p>{/if}
		<Dialog.Footer class="items-center border-t px-6 py-3 sm:justify-between">
			<div class="flex flex-wrap items-center gap-2">{@render footer?.()}</div>
			<div class="flex items-center gap-4">
				<span class="text-base font-semibold tabular-nums">{$t`Total: CHF ${total.toFixed(2)}`}</span>
				<Button type="submit" disabled={busy}>{submitLabel}</Button>
			</div>
		</Dialog.Footer>
	</fieldset>
</form>
