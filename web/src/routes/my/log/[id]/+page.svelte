<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import Angry from '@lucide/svelte/icons/angry';
	import Frown from '@lucide/svelte/icons/frown';
	import Laugh from '@lucide/svelte/icons/laugh';
	import Meh from '@lucide/svelte/icons/meh';
	import Smile from '@lucide/svelte/icons/smile';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Trash from '@lucide/svelte/icons/trash-2';
	import Trophy from '@lucide/svelte/icons/trophy';
	import { api, call } from '$lib/api';
	import Bars from '$lib/Bars.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ask } from '$lib/confirm.svelte';
	import GroupInfo from '$lib/my/GroupInfo.svelte';
	import { disciplineName, fmt, fmtDateTime, fmtTotal, kindLabel, myKey, myQuery } from '$lib/my.svelte';
	import TargetFace from '$lib/TargetFace.svelte';
	import { targets } from '$lib/targets';

	const queryClient = useQueryClient();
	const id = $derived(Number(page.params.id));
	const my = myQuery();
	const q = createQuery(() => ({
		queryKey: [...myKey, 'training', id],
		queryFn: () => call(api.GET('/api/my/trainings/{id}', { params: { path: { id } } }))
	}));
	const s = $derived(q.data?.session);
	const f = $derived(q.data?.feedback);
	const discipline = $derived(my.data?.disciplines.find((d) => d.id === s?.disciplineId));
	const weapon = $derived(my.data?.weapons.find((w) => w.id === s?.weaponId));
	const target = $derived(discipline?.target ? targets[discipline.target] : undefined);
	const positioned = $derived((s?.series ?? []).flatMap((x) => (x.x ?? []).map((px, i) => ({ x: px, y: x.y![i] }))));
	const group = $derived.by(() => {
		if (!positioned.length) return null;
		const x = positioned.reduce((a, p) => a + p.x, 0) / positioned.length;
		const y = positioned.reduce((a, p) => a + p.y, 0) / positioned.length;
		const radius = positioned.reduce((a, p) => a + Math.hypot(p.x - x, p.y - y), 0) / positioned.length;
		return { x, y, radius, shots: positioned.length };
	});
	const perSeries = $derived((s?.series ?? []).map((x, i) => ({ label: String(i + 1), value: x.shots ? x.total / x.shots : 0, text: fmtTotal(x.total) })));
	const best = $derived(perSeries.length > 1 ? perSeries.reduce((b, x, i) => (x.value > perSeries[b].value ? i : b), 0) : undefined);
	const faces = [Angry, Frown, Meh, Smile, Laugh];

	/** Deletes this session after confirmation and returns to the log. */
	async function remove() {
		if (!(await ask({ title: $t`Delete this session?`, action: $t`Delete`, destructive: true }))) return;
		await call(api.DELETE('/api/my/trainings/{id}', { params: { path: { id } } }));
		await goto('/my/log');
		queryClient.removeQueries({ queryKey: [...myKey, 'training', id] });
		await queryClient.invalidateQueries({ queryKey: myKey });
	}
</script>

<svelte:head><title>{$t`Session`} · Standblatt</title></svelte:head>

{#if s && f}
	<div class="flex items-center gap-2">
		<Button variant="ghost" size="icon-sm" href="/my/log" aria-label={$t`Back`}><ArrowLeft /></Button>
		<div class="min-w-0 flex-1">
			<p class="truncate text-xs text-muted-foreground">{$t(disciplineName(discipline))}{weapon ? ` · ${weapon.name}` : ''} · {fmtDateTime(s.startedAt)}</p>
			<h1 class="flex items-center gap-2 truncate text-lg font-semibold">
				{#if s.kind === 'competition'}<Trophy class="size-4 text-amber-500" />{/if}
				{s.title || s.place || $t(kindLabel(s.kind))}
				{#if s.rank}<span class="shrink-0 rounded-full bg-amber-500/15 px-2 py-0.5 text-xs font-medium text-amber-700 dark:text-amber-400">{$t`Rank ${s.rank} of ${s.participants ?? 0}`}</span>{/if}
			</h1>
		</div>
		<Button variant="outline" size="icon-sm" href="/my/train?edit={s.id}" aria-label={$t`Edit`}><Pencil /></Button>
		<Button variant="outline" size="icon-sm" onclick={remove} aria-label={$t`Delete`}><Trash /></Button>
	</div>

	<section class="grid grid-cols-3 gap-2 rounded-xl border bg-card p-4 text-center shadow-sm">
		<div>
			<div class="text-3xl font-semibold tabular-nums">{fmtTotal(s.total)}</div>
			<div class="text-xs text-muted-foreground">{$t`Total`}</div>
		</div>
		<div>
			<div class="text-3xl font-semibold tabular-nums">{fmt(s.shots ? s.total / s.shots : 0, 2)}</div>
			<div class="text-xs text-muted-foreground">{$t`per shot`}</div>
		</div>
		<div>
			<div class="text-3xl font-semibold tabular-nums">{f.rank ? `${f.rank}/${f.of}` : '–'}</div>
			<div class="text-xs text-muted-foreground">{$t`Rank among sessions with ${s.shots} shots`}</div>
		</div>
	</section>

	{#if perSeries.length > 1}
		<section class="rounded-xl border bg-card p-4 shadow-sm">
			<h2 class="mb-3 font-semibold">{$t`Series`}</h2>
			<Bars items={perSeries} from={Math.max(0, Math.min(...perSeries.map((x) => x.value)) - 0.5)} highlight={best} label={$t`Points per shot by series`} />
		</section>
	{/if}

	{#if s.series.some((x) => x.values?.length)}
		<section class="rounded-xl border bg-card p-4 shadow-sm">
			<h2 class="mb-3 font-semibold">{$t`Shots`}</h2>
			<div class="grid gap-2">
				{#each s.series as x, i (i)}
					{#if x.values?.length}
						<div class="flex items-center gap-2 text-sm">
							<span class="w-6 text-xs text-muted-foreground">{i + 1}</span>
							<div class="flex flex-1 flex-wrap gap-1">
								{#each x.values as v, j (j)}
									<span class="min-w-8 rounded px-1 py-0.5 text-center tabular-nums {v >= 10 ? 'bg-amber-500/20 font-semibold' : 'bg-muted'}">{discipline?.decimal ? fmt(v, 1) : v}</span>
								{/each}
							</div>
							<span class="font-semibold tabular-nums">{fmtTotal(x.total)}</span>
						</div>
					{/if}
				{/each}
			</div>
		</section>
	{/if}

	{#if target && group}
		<section class="grid gap-4 rounded-xl border bg-card p-4 shadow-sm sm:grid-cols-2">
			<TargetFace {target} shots={positioned} {group} zoom label={$t`Hits of this session`} />
			<GroupInfo {group} />
		</section>
	{/if}

	{#if s.rating || s.tags.length || s.notes || s.place}
		<section class="grid gap-3 rounded-xl border bg-card p-4 text-sm shadow-sm">
			{#if s.rating}{@const Face = faces[s.rating - 1]}<Face class="size-7 text-primary" />{/if}
			{#if s.tags.length}
				<div class="flex flex-wrap gap-1.5">{#each s.tags as tag (tag)}<span class="rounded-full bg-muted px-2.5 py-1 text-xs">{tag}</span>{/each}</div>
			{/if}
			{#if s.place}<p class="text-muted-foreground">{s.place}</p>{/if}
			{#if s.notes}<p class="whitespace-pre-line">{s.notes}</p>{/if}
		</section>
	{/if}
{:else if q.isError}
	<p class="text-sm text-muted-foreground">{$t`Session not found.`} <a class="underline" href="/my/log">{$t`Back to the log`}</a></p>
{/if}
