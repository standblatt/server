<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import ArrowDownRight from '@lucide/svelte/icons/arrow-down-right';
	import ArrowUpRight from '@lucide/svelte/icons/arrow-up-right';
	import Target from '@lucide/svelte/icons/target';
	import { api, call } from '$lib/api';
	import Bars from '$lib/Bars.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { formatDate } from '$lib/format';
	import FormField from '$lib/FormField.svelte';
	import GroupInfo from '$lib/my/GroupInfo.svelte';
	import { disciplineName, fmt, fmtDay, fmtTotal, myKey, myQuery, projected, statsQuery } from '$lib/my.svelte';
	import TargetFace from '$lib/TargetFace.svelte';
	import { targets } from '$lib/targets';
	import TrendChart from '$lib/TrendChart.svelte';

	const queryClient = useQueryClient();
	const my = myQuery();
	const data = $derived(my.data!);
	const used = $derived(new Set(data.recent.map((s) => s.disciplineId)));
	const fallback = $derived(data.recent[0]?.disciplineId ?? data.disciplines[0]?.id);
	const id = $derived(Number(page.url.searchParams.get('d')) || fallback);
	const discipline = $derived(data.disciplines.find((d) => d.id === id));
	const weaponId = $derived(Number(page.url.searchParams.get('w')) || 0);
	/** Name of the weapon with the given ID, empty if unknown. */
	const weaponName = (wid: number) => data.weapons.find((w) => w.id === wid)?.name ?? '';
	const stats = statsQuery(
		() => id,
		() => true,
		() => weaponId
	);
	const st = $derived(stats.data?.stats);
	const standard = $derived(discipline ? discipline.seriesCount * discipline.shotsPerSeries : 0);
	const goal = $derived(stats.data?.goal);
	const target = $derived(discipline?.target ? targets[discipline.target] : undefined);
	const recent = $derived(st?.recent ?? st?.perShot ?? 0);
	const delta = $derived(st?.recent != null && st?.previous != null ? st.recent - st.previous : null);

	const seriesItems = $derived((st?.series ?? []).map((x) => ({ label: String(x.number), value: x.perShot, text: fmt(x.perShot, 2) })));
	const strongest = $derived(seriesItems.length > 1 ? seriesItems.reduce((b, x, i) => (x.value > seriesItems[b].value ? i : b), 0) : undefined);
	const weakest = $derived(seriesItems.length > 1 ? seriesItems.reduce((b, x, i) => (x.value < seriesItems[b].value ? i : b), 0) : undefined);
	const shotsLogged = $derived((st?.distribution ?? []).reduce((a, b) => a + b, 0));
	const distribution = $derived(
		(st?.distribution ?? [])
			.map((n, v) => ({ label: String(v), value: n, text: shotsLogged ? `${Math.round((n / shotsLogged) * 100)}%` : '' }))
			.reverse()
			.filter((x, i, all) => x.value > 0 || all.slice(i).some((y) => y.value > 0))
	);

	let goalTotal = $state('');
	let goalDue = $state('');
	let editingGoal = $state(false);

	/** Opens the goal editor with the current goal. */
	function editGoal() {
		goalTotal = goal ? String(goal.targetTotal) : '';
		goalDue = goal?.dueOn ?? '';
		editingGoal = true;
	}

	/** Saves the goal of this discipline. */
	async function saveGoal(e: SubmitEvent) {
		e.preventDefault();
		try {
			await call(api.PUT('/api/my/goals/{id}', { params: { path: { id } }, body: { targetTotal: Number(goalTotal), dueOn: goalDue } }));
			editingGoal = false;
			await queryClient.invalidateQueries({ queryKey: myKey });
		} catch (err) {
			toast.error((err as Error).message);
		}
	}

	/** Deletes the goal of this discipline. */
	async function removeGoal() {
		await call(api.DELETE('/api/my/goals/{id}', { params: { path: { id } } }));
		editingGoal = false;
		await queryClient.invalidateQueries({ queryKey: myKey });
	}

	/** Shows the progress of a discipline and weapon (0 = all weapons). */
	function show(d: number, w: number) {
		goto(`/my/progress?d=${d}${w ? `&w=${w}` : ''}`, { replaceState: true, noScroll: true });
	}
</script>

<svelte:head><title>{$t`Progress`} · Standblatt</title></svelte:head>

<div class="flex flex-wrap items-center gap-3">
	<h1 class="text-lg font-semibold tracking-tight">{$t`Progress`}</h1>
	<select
		class="ml-auto rounded-md border bg-background px-2 py-1.5 text-sm"
		value={id}
		onchange={(e) => show(Number(e.currentTarget.value), weaponId)}
		aria-label={$t`Discipline`}
	>
		{#each [...data.disciplines].sort((a, b) => Number(used.has(b.id)) - Number(used.has(a.id))) as d (d.id)}
			<option value={d.id}>{$t(disciplineName(d))}</option>
		{/each}
	</select>
	{#if data.weapons.length}
		<select class="rounded-md border bg-background px-2 py-1.5 text-sm" value={weaponId} onchange={(e) => show(id, Number(e.currentTarget.value))} aria-label={$t`Weapon`}>
			<option value={0}>{$t`All weapons`}</option>
			{#each data.weapons as w (w.id)}
				<option value={w.id}>{w.name}</option>
			{/each}
		</select>
	{/if}
</div>

{#if st && st.sessions === 0}
	<section class="rounded-xl border border-dashed p-8 text-center text-sm text-muted-foreground">
		<Target class="mx-auto mb-3 size-8" />
		<p>{$t`No sessions in this discipline yet. Log a few and your progress shows up here.`}</p>
		<Button class="mt-4" href="/my/train?d={id}">{$t`Start training`}</Button>
	</section>
{:else if st && discipline}
	<section class="grid grid-cols-2 gap-2 sm:grid-cols-4">
		<div class="rounded-xl border bg-card p-3 shadow-sm">
			<div class="text-xs text-muted-foreground">{$t`Last 30 days`}</div>
			<div class="text-2xl font-semibold tabular-nums">{fmt(recent, 2)}</div>
			{#if delta != null}
				<div class="flex items-center gap-0.5 text-xs {delta >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}">
					{#if delta >= 0}<ArrowUpRight class="size-3" />{:else}<ArrowDownRight class="size-3" />{/if}{delta >= 0 ? '+' : ''}{fmt(delta, 2)}
				</div>
			{:else}
				<div class="text-xs text-muted-foreground">{$t`per shot`}</div>
			{/if}
		</div>
		<div class="rounded-xl border bg-card p-3 shadow-sm">
			<div class="text-xs text-muted-foreground">{$t`Personal best`}</div>
			{#if st.best}
				<a href="/my/log/{st.best.sessionId}" class="block text-2xl font-semibold tabular-nums hover:underline">{fmtTotal(st.best.total)}</a>
				<div class="text-xs text-muted-foreground">{fmtDay(st.best.at)}</div>
			{:else}
				<div class="text-2xl font-semibold">–</div>
				<div class="text-xs text-muted-foreground">{$t`${standard} shots`}</div>
			{/if}
		</div>
		<div class="rounded-xl border bg-card p-3 shadow-sm">
			<div class="text-xs text-muted-foreground">{$t`Sessions`}</div>
			<div class="text-2xl font-semibold tabular-nums">{st.sessions}</div>
			<div class="text-xs text-muted-foreground tabular-nums">{$t`${st.shots} shots`}</div>
		</div>
		<div class="rounded-xl border bg-card p-3 shadow-sm">
			<div class="text-xs text-muted-foreground">{$t`Competition vs. training`}</div>
			{#if st.competition != null && st.training != null}
				<div class="text-2xl font-semibold tabular-nums">{st.competition - st.training >= 0 ? '+' : ''}{fmt(st.competition - st.training, 2)}</div>
				<div class="text-xs text-muted-foreground tabular-nums">{fmt(st.competition, 2)} / {fmt(st.training, 2)}</div>
			{:else}
				<div class="text-2xl font-semibold">–</div>
				<div class="text-xs text-muted-foreground">{$t`needs both`}</div>
			{/if}
		</div>
	</section>

	{#if st.weapons.length > 1}
		<section class="rounded-xl border bg-card p-4 shadow-sm" data-testid="weapons">
			<h2 class="mb-1 font-semibold">{$t`Weapons compared`}</h2>
			<p class="mb-3 text-xs text-muted-foreground">{$t`Points per shot of all sessions with each weapon. Tap a weapon to see only its sessions.`}</p>
			<div class="grid grid-cols-[minmax(0,1fr)_auto_auto_auto] items-center gap-x-4 text-sm">
				<span></span>
				<span class="text-right text-xs text-muted-foreground">{$t`Average`}</span>
				<span class="text-right text-xs text-muted-foreground">{$t`30 days`}</span>
				<span class="text-right text-xs text-muted-foreground">{$t`Best`}</span>
				{#each st.weapons as w, i (w.weaponId)}
					<button
						type="button"
						class="col-span-4 grid grid-cols-subgrid items-center rounded-md py-2 text-left hover:bg-muted/60 {w.weaponId === weaponId ? 'bg-muted' : ''}"
						onclick={() => show(id, w.weaponId === weaponId ? 0 : w.weaponId)}
						aria-pressed={w.weaponId === weaponId}
					>
						<span class="min-w-0">
							<span class="block truncate font-medium">{weaponName(w.weaponId)}</span>
							<span class="block text-xs text-muted-foreground">{$t`${w.sessions} sessions`}</span>
						</span>
						<span class="text-right font-semibold tabular-nums {i === 0 ? 'text-emerald-600 dark:text-emerald-400' : ''}">{fmt(w.perShot, 2)}</span>
						<span class="text-right tabular-nums">{w.recent != null ? fmt(w.recent, 2) : '–'}</span>
						<span class="text-right tabular-nums">{w.best != null ? fmtTotal(w.best) : '–'}</span>
					</button>
				{/each}
			</div>
		</section>
	{/if}

	<section class="rounded-xl border bg-card p-4 shadow-sm">
		<div class="mb-2 flex flex-wrap items-center justify-between gap-2">
			<h2 class="font-semibold">{$t`Trend`}</h2>
			<div class="flex items-center gap-3 text-xs text-muted-foreground">
				<span class="flex items-center gap-1"><span class="size-2 rounded-full bg-foreground/70"></span>{$t`Training`}</span>
				<span class="flex items-center gap-1"><span class="size-2.5 rounded-full bg-amber-500"></span>{$t`Competition`}</span>
				<span class="flex items-center gap-1"><span class="h-0.5 w-3 bg-chart-3"></span>{$t`Average of 5`}</span>
				{#if goal}<span class="flex items-center gap-1"><span class="h-0.5 w-3 border-t-2 border-dashed border-emerald-500"></span>{$t`Goal`}</span>{/if}
			</div>
		</div>
		<TrendChart points={st.points} goal={goal ? goal.targetTotal / standard : null} label={$t`Points per shot of every session`} />
		<p class="mt-1 text-xs text-muted-foreground">{$t`Points per shot of every session, oldest on the left.`}</p>
	</section>

	<section class="rounded-xl border bg-card p-4 shadow-sm">
		<div class="flex items-center justify-between gap-2">
			<h2 class="font-semibold">{$t`Goal`}</h2>
			{#if !editingGoal}<Button variant="ghost" size="sm" onclick={editGoal}>{goal ? $t`Change` : $t`Set a goal`}</Button>{/if}
		</div>
		{#if editingGoal}
			<form class="mt-3 grid gap-3 sm:grid-cols-[1fr_1fr_auto] sm:items-end" onsubmit={saveGoal}>
				<FormField label={$t`Total over ${standard} shots`}>
					<Input type="number" step="0.1" min="1" max={standard * (discipline.decimal ? 10.9 : discipline.maxShot)} bind:value={goalTotal} required name="goal" />
				</FormField>
				<FormField label={$t`By (optional)`}><Input type="date" bind:value={goalDue} name="due" /></FormField>
				<div class="flex gap-2">
					<Button type="submit">{$t`Save`}</Button>
					{#if goal}<Button variant="ghost" onclick={removeGoal}>{$t`Remove`}</Button>{/if}
				</div>
			</form>
		{:else if goal}
			{@const now = projected(recent, discipline)}
			<p class="mt-1 text-sm text-muted-foreground">
				{$t`${fmtTotal(goal.targetTotal)} over ${standard} shots`}{#if goal.dueOn} · {$t`by ${formatDate(goal.dueOn)}`}{/if}
			</p>
			<div class="mt-2 h-3 overflow-hidden rounded-full bg-muted">
				<div class="h-full rounded-full bg-emerald-500" style:width="{Math.min(100, (now / goal.targetTotal) * 100)}%"></div>
			</div>
			<p class="mt-1 text-sm">
				{now >= goal.targetTotal
					? $t`Your last 30 days average ${fmtTotal(now)} — goal reached!`
					: $t`Your last 30 days average ${fmtTotal(now)}. ${fmt(Math.max(0.01, goal.targetTotal / standard - recent), 2)} more per shot to go.`}
			</p>
		{:else}
			<p class="mt-1 text-sm text-muted-foreground">{$t`A goal gives every session a meaning. Set the total you want to reach.`}</p>
		{/if}
	</section>

	<div class="grid gap-4 md:grid-cols-2 *:min-w-0">
		{#if seriesItems.length > 1}
			<section class="rounded-xl border bg-card p-4 shadow-sm">
				<h2 class="mb-1 font-semibold">{$t`Series curve`}</h2>
				<p class="mb-3 text-xs text-muted-foreground">
					{$t`Points per shot by series. Series ${(strongest ?? 0) + 1} is your strongest, series ${(weakest ?? 0) + 1} your weakest.`}
				</p>
				<Bars items={seriesItems} from={Math.max(0, Math.min(...seriesItems.map((x) => x.value)) - 0.3)} highlight={strongest} label={$t`Series curve`} />
			</section>
		{/if}
		{#if shotsLogged}
			<section class="rounded-xl border bg-card p-4 shadow-sm">
				<h2 class="mb-1 font-semibold">{$t`Shot values`}</h2>
				<p class="mb-3 text-xs text-muted-foreground">{$t`Share of every value over ${shotsLogged} logged shots.`}</p>
				<Bars items={distribution} highlight={0} label={$t`Shot values`} />
			</section>
		{/if}
	</div>

	{#if st.group && target}
		<section class="grid gap-4 rounded-xl border bg-card p-4 shadow-sm sm:grid-cols-2">
			<TargetFace {target} group={st.group} zoom label={$t`Centre and size of your group`} />
			<GroupInfo group={st.group} />
		</section>
	{:else if target}
		<p class="rounded-xl border border-dashed p-4 text-sm text-muted-foreground">
			{$t`Tip: tap your hits on the target when you log a session and Standblatt shows where your group sits and how big it is.`}
		</p>
	{/if}

	{#if st.tags.length}
		<section class="rounded-xl border bg-card p-4 shadow-sm">
			<h2 class="mb-1 font-semibold">{$t`What your tags tell you`}</h2>
			<p class="mb-3 text-xs text-muted-foreground">{$t`Points per shot of sessions with a tag compared with your average.`}</p>
			<ul class="divide-y text-sm">
				{#each st.tags as tag (tag.tag)}
					<li class="flex items-center gap-3 py-2">
						<span class="rounded-full bg-muted px-2.5 py-0.5 text-xs">{tag.tag}</span>
						<span class="text-xs text-muted-foreground">{$t`${tag.sessions} sessions`}</span>
						<span class="ml-auto tabular-nums">{fmt(tag.perShot, 2)}</span>
						<span class="w-14 text-right font-semibold tabular-nums {tag.delta >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}">
							{tag.delta >= 0 ? '+' : ''}{fmt(tag.delta, 2)}
						</span>
					</li>
				{/each}
			</ul>
		</section>
	{/if}
{/if}
