<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import ArrowDownRight from '@lucide/svelte/icons/arrow-down-right';
	import ArrowUpRight from '@lucide/svelte/icons/arrow-up-right';
	import Flame from '@lucide/svelte/icons/flame';
	import Plus from '@lucide/svelte/icons/plus';
	import Trophy from '@lucide/svelte/icons/trophy';
	import { Button } from '$lib/components/ui/button';
	import SessionRow from '$lib/my/SessionRow.svelte';
	import { meQuery } from '$lib/me';
	import { disciplineName, fmt, fmtTotal, myQuery, projected, statsQuery } from '$lib/my.svelte';
	import TrendChart from '$lib/TrendChart.svelte';

	const my = myQuery();
	const me = meQuery();
	const data = $derived(my.data!);
	const disciplines = $derived(data.disciplines);
	const byId = $derived(new Map(disciplines.map((d) => [d.id, d])));
	const mainId = $derived(data.recent[0]?.disciplineId);
	const main = $derived(mainId ? byId.get(mainId) : undefined);
	const ordered = $derived(main ? [main, ...disciplines.filter((d) => d.id !== main.id)] : disciplines);
	const stats = statsQuery(
		() => mainId,
		() => !!my.data
	);
	const goal = $derived(data.goals.find((g) => g.disciplineId === mainId));
	const standard = $derived(main ? main.seriesCount * main.shotsPerSeries : 0);

	const hour = new Date().getHours();
	const greeting = $derived(hour < 11 ? $t`Good morning` : hour < 18 ? $t`Hello` : $t`Good evening`);
	const recent = $derived(stats.data?.stats.recent ?? null);
	const previous = $derived(stats.data?.stats.previous ?? null);
	const delta = $derived(recent != null && previous != null ? recent - previous : null);
</script>

<svelte:head><title>{$t`Today`} · Standblatt</title></svelte:head>

<section>
	<h1 class="text-2xl font-semibold tracking-tight">{greeting}, {data.profile.firstname || me.data?.user.name}</h1>
	<p class="text-sm text-muted-foreground">
		{#if data.week.sessions}
			{$plural(data.week.sessions, { one: '# training this week', other: '# trainings this week' })} · {$plural(data.week.shots, { one: '# shot', other: '# shots' })}
		{:else}
			{$t`No training yet this week. Ready when you are.`}
		{/if}
	</p>
</section>

<section class="rounded-xl border bg-card p-4 shadow-sm">
	<div class="mb-3 flex items-center justify-between gap-2">
		<h2 class="font-semibold">{$t`Start training`}</h2>
		{#if data.week.streak > 1}
			<span class="flex items-center gap-1 text-sm font-medium text-amber-600 dark:text-amber-400" title={$t`Weeks in a row with training`}>
				<Flame class="size-4" />{$plural(data.week.streak, { one: '# week in a row', other: '# weeks in a row' })}
			</span>
		{/if}
	</div>
	<div class="-mx-1 flex gap-2 overflow-x-auto px-1 pb-1 [scrollbar-width:none]">
		{#each ordered as d, i (d.id)}
			<a
				href="/my/train?d={d.id}"
				class="flex shrink-0 items-center gap-1.5 rounded-full border px-3 py-1.5 text-sm whitespace-nowrap {i === 0 && main ? 'border-primary bg-primary text-primary-foreground' : 'hover:bg-muted'}"
			>
				{#if i === 0 && main}<Plus class="size-4" />{/if}{$t(disciplineName(d))}
			</a>
		{/each}
	</div>
</section>

{#if data.competitions > 0}
	<a href="/my/log#competitions" class="flex items-center gap-3 rounded-xl border border-amber-500/40 bg-amber-500/10 p-4 text-sm hover:bg-amber-500/15">
		<Trophy class="size-5 shrink-0 text-amber-500" />
		<span class="flex-1">{$plural(data.competitions, { one: 'We found # competition result of yours. Add it to your log with one tap.', other: 'We found # competition results of yours. Add them to your log with one tap.' })}</span>
	</a>
{/if}

{#if main && stats.data && stats.data.stats.sessions > 0}
	{@const st = stats.data.stats}
	<a href="/my/progress?d={main.id}" class="block rounded-xl border bg-card p-4 shadow-sm hover:border-primary/40">
		<div class="flex flex-wrap items-end justify-between gap-3">
			<div>
				<p class="text-xs text-muted-foreground">{$t(disciplineName(main))} · {$t`last 30 days`}</p>
				<p class="text-3xl font-semibold tabular-nums">
					{fmt(recent ?? st.perShot, 2)}<span class="ml-1 text-sm font-normal text-muted-foreground">{$t`per shot`}</span>
				</p>
				<p class="text-sm text-muted-foreground tabular-nums">{$t`≙ ${fmtTotal(projected(recent ?? st.perShot, main))} of ${standard} shots`}</p>
			</div>
			{#if delta != null}
				<span class="flex items-center gap-1 rounded-full px-2 py-1 text-sm font-medium {delta >= 0 ? 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-400' : 'bg-red-500/10 text-red-700 dark:text-red-400'}">
					{#if delta >= 0}<ArrowUpRight class="size-4" />{:else}<ArrowDownRight class="size-4" />{/if}
					{delta >= 0 ? '+' : ''}{fmt(delta, 2)}
					<span class="font-normal opacity-80">{$t`vs. 30 days before`}</span>
				</span>
			{/if}
		</div>
		{#if st.points.length > 1}
			<div class="mt-3"><TrendChart points={st.points.slice(-20)} goal={goal ? goal.targetTotal / standard : null} label={$t`Trend`} /></div>
		{/if}
		{#if goal}
			{@const now = projected(recent ?? st.perShot, main)}
			<div class="mt-2">
				<div class="flex justify-between text-xs text-muted-foreground">
					<span>{$t`Goal ${fmtTotal(goal.targetTotal)}`}</span>
					<span>{now >= goal.targetTotal ? $t`reached` : $t`${fmt(goal.targetTotal - now, 1)} to go`}</span>
				</div>
				<div class="mt-1 h-2 overflow-hidden rounded-full bg-muted">
					<div class="h-full rounded-full bg-emerald-500" style:width="{Math.min(100, (now / goal.targetTotal) * 100)}%"></div>
				</div>
			</div>
		{/if}
	</a>
{/if}

<section class="overflow-hidden rounded-xl border bg-card shadow-sm">
	<div class="flex items-center justify-between border-b px-4 py-3">
		<h2 class="font-semibold">{$t`Recent sessions`}</h2>
		{#if data.recent.length}<a href="/my/log" class="text-sm text-muted-foreground hover:text-foreground">{$t`All`}</a>{/if}
	</div>
	{#each data.recent as s (s.id)}
		<SessionRow session={s} discipline={byId.get(s.disciplineId)} />
	{:else}
		<div class="grid gap-4 p-6 text-sm">
			<p class="text-base font-medium">{$t`Your log is empty — let's change that.`}</p>
			<ol class="grid gap-2 text-muted-foreground">
				<li><strong class="text-foreground">1.</strong> {$t`Tap the + after your next training.`}</li>
				<li><strong class="text-foreground">2.</strong> {$t`Type the series totals — or every shot, or tap the hits on the target.`}</li>
				<li><strong class="text-foreground">3.</strong> {$t`See right away how the session compares and where you improve.`}</li>
			</ol>
			<Button href="/my/train{ordered[0] ? `?d=${ordered[0].id}` : ''}" class="w-fit">{$t`Log my first training`}</Button>
		</div>
	{/each}
</section>
