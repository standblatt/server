<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import Trophy from '@lucide/svelte/icons/trophy';
	import { disciplineName, fmt, fmtDay, fmtTotal, myQuery, type Discipline, type TrainingSession } from '$lib/my.svelte';

	let { session, discipline }: { session: TrainingSession; discipline?: Discipline } = $props();

	const my = myQuery();
	const weapon = $derived(my.data?.weapons.find((w) => w.id === session.weaponId));
	const perShot = $derived(session.shots ? session.total / session.shots : 0);
	const series = $derived(session.series.map((s) => (s.shots ? s.total / s.shots : 0)));
	const lo = $derived(Math.min(...series) - 0.3);
	const hi = $derived(Math.max(...series));
</script>

<a href="/my/log/{session.id}" class="flex items-center gap-3 px-4 py-3 hover:bg-muted/60" data-testid="session-row">
	<div class="min-w-0 flex-1">
		<div class="flex items-center gap-2 text-xs text-muted-foreground">
			<span>{fmtDay(session.startedAt)}</span>
			<span>·</span>
			<span class="truncate">{$t(disciplineName(discipline))}{weapon ? ` · ${weapon.name}` : ''}</span>
			{#if session.kind === 'competition'}<Trophy class="size-3.5 shrink-0 text-amber-500" aria-label={$t`Competition`} />{/if}
			{#if session.rank}<span class="shrink-0 font-medium text-amber-600 dark:text-amber-400">{$t`Rank ${session.rank} of ${session.participants ?? 0}`}</span>{/if}
		</div>
		<div class="truncate text-sm">
			{session.title || session.place || (session.kind === 'competition' ? $t`Competition` : $t`Training`)}
			{#each session.tags as tag (tag)}<span class="ml-1.5 rounded bg-muted px-1.5 py-0.5 text-[11px] text-muted-foreground">{tag}</span>{/each}
		</div>
	</div>
	{#if series.length > 1}
		<div class="hidden h-6 items-end gap-px sm:flex" aria-hidden="true">
			{#each series as v, i (i)}
				<span class="w-1.5 rounded-t-sm bg-chart-3/70" style:height="{Math.max(10, ((v - lo) / (hi - lo || 1)) * 100)}%"></span>
			{/each}
		</div>
	{/if}
	<div class="text-right">
		<div class="text-lg leading-tight font-semibold tabular-nums">{fmtTotal(session.total)}</div>
		<div class="text-[11px] text-muted-foreground tabular-nums">{$t`${session.shots} shots · Ø ${fmt(perShot, 2)}`}</div>
	</div>
	<ChevronRight class="size-4 shrink-0 text-muted-foreground" />
</a>
