<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import Download from '@lucide/svelte/icons/download';
	import Trophy from '@lucide/svelte/icons/trophy';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { formatDate } from '$lib/format';
	import SessionRow from '$lib/my/SessionRow.svelte';
	import { disciplineName, fmtTotal, myKey, myQuery, type TrainingSession } from '$lib/my.svelte';
	import { intlLocale } from '$lib/locale.svelte';
	import PageHeader from '$lib/PageHeader.svelte';

	const queryClient = useQueryClient();
	const my = myQuery();
	const disciplines = $derived(my.data?.disciplines ?? []);
	const byId = $derived(new Map(disciplines.map((d) => [d.id, d])));

	let discipline = $state(0);
	let kind = $state<'' | 'training' | 'competition'>('');
	let limit = $state(50);

	const sessions = createQuery(() => ({
		queryKey: [...myKey, 'trainings', discipline, kind, limit],
		queryFn: () => call(api.GET('/api/my/trainings', { params: { query: { discipline, kind, limit } } })),
		placeholderData: (prev: TrainingSession[] | undefined) => prev
	}));
	const competitions = createQuery(() => ({
		queryKey: [...myKey, 'competitions'],
		queryFn: () => call(api.GET('/api/my/competitions'))
	}));

	const groups = $derived.by(() => {
		const out: { month: string; items: TrainingSession[] }[] = [];
		const f = new Intl.DateTimeFormat(intlLocale(), { month: 'long', year: 'numeric' });
		for (const s of sessions.data ?? []) {
			const month = f.format(new Date(s.startedAt));
			if (out.at(-1)?.month !== month) out.push({ month, items: [] });
			out.at(-1)!.items.push(s);
		}
		return out;
	});

	const preferred = $derived(my.data?.recent[0]?.disciplineId ?? disciplines[0]?.id ?? 0);
	let picks = $state<Record<number, number>>({});
	let importing = $state(false);

	/** Imports one competition entry into the log with the chosen discipline. */
	async function importOne(entryId: number) {
		await call(api.POST('/api/my/competitions/{id}', { params: { path: { id: entryId } }, body: { disciplineId: picks[entryId] ?? preferred } }));
	}

	/** Imports the given competition entries one after the other and refreshes the log. */
	async function add(ids: number[]) {
		importing = true;
		try {
			for (const id of ids) await importOne(id);
			toast.success($plural(ids.length, { one: 'Competition added to your log', other: '# competitions added to your log' }));
			await queryClient.invalidateQueries({ queryKey: myKey });
		} catch (err) {
			toast.error((err as Error).message);
		} finally {
			importing = false;
		}
	}
</script>

<svelte:head><title>{$t`Log`} · Standblatt</title></svelte:head>

<PageHeader help="my-training/logging" title={$t`Log`} description={$t`Every training and competition you logged.`}>
	{#snippet actions()}
		<Button variant="outline" size="sm" href="/api/my/export" download><Download />{$t`CSV`}</Button>
	{/snippet}
</PageHeader>

{#if competitions.data?.length}
	<section id="competitions" class="overflow-hidden rounded-xl border border-amber-500/40 bg-card shadow-sm">
		<div class="flex flex-wrap items-center gap-3 border-b bg-amber-500/10 px-4 py-3">
			<Trophy class="size-5 text-amber-500" />
			<div class="min-w-0 flex-1">
				<h2 class="font-semibold">{$t`Your competition results`}</h2>
				<p class="text-xs text-muted-foreground">{$t`Found by your license number and birthday in closed competitions run with Standblatt.`}</p>
			</div>
			<Button size="sm" disabled={importing} onclick={() => add(competitions.data!.map((c) => c.entryId))}>{$t`Add all`}</Button>
		</div>
		<ul class="divide-y">
			{#each competitions.data as c (c.entryId)}
				<li class="flex flex-wrap items-center gap-3 px-4 py-3 text-sm">
					<div class="min-w-0 flex-1">
						<div class="font-medium">{c.event}</div>
						<div class="text-xs text-muted-foreground">{formatDate(c.date)} · {c.club} · {c.programme}{#if c.rank} · {$t`Rank ${c.rank} of ${c.participants}`}{/if}</div>
					</div>
					<span class="text-lg font-semibold tabular-nums">{fmtTotal(c.total)}</span>
					<select
						class="rounded-md border bg-background px-2 py-1 text-xs"
						value={picks[c.entryId] ?? preferred}
						onchange={(e) => (picks[c.entryId] = Number(e.currentTarget.value))}
						aria-label={$t`Discipline`}
					>
						{#each disciplines as d (d.id)}<option value={d.id}>{$t(disciplineName(d))}</option>{/each}
					</select>
					<Button size="sm" variant="outline" disabled={importing} onclick={() => add([c.entryId])}>{$t`Add`}</Button>
				</li>
			{/each}
		</ul>
	</section>
{/if}

<div class="flex flex-wrap items-center gap-2">
	<select class="rounded-md border bg-background px-2 py-1.5 text-sm" bind:value={discipline} aria-label={$t`Discipline`}>
		<option value={0}>{$t`All disciplines`}</option>
		{#each disciplines as d (d.id)}<option value={d.id}>{$t(disciplineName(d))}</option>{/each}
	</select>
	<div class="flex rounded-lg bg-muted p-0.5 text-sm" role="radiogroup" aria-label={$t`Kind`}>
		{#each [['', $t`All`], ['training', $t`Training`], ['competition', $t`Competition`]] as [k, label] (k)}
			<button type="button" role="radio" aria-checked={kind === k} class="rounded-md px-3 py-1 {kind === k ? 'bg-background shadow-sm' : 'text-muted-foreground'}" onclick={() => (kind = k as typeof kind)}>{label}</button>
		{/each}
	</div>
</div>

{#each groups as g (g.month)}
	<section class="overflow-hidden rounded-xl border bg-card shadow-sm">
		<h2 class="border-b bg-muted/40 px-4 py-2 text-xs font-semibold tracking-wide text-muted-foreground uppercase">{g.month}</h2>
		<div class="divide-y">
			{#each g.items as s (s.id)}<SessionRow session={s} discipline={byId.get(s.disciplineId)} />{/each}
		</div>
	</section>
{:else}
	{#if sessions.isSuccess}
		<p class="rounded-xl border border-dashed p-8 text-center text-sm text-muted-foreground">{$t`No sessions yet.`} <a class="underline" href="/my/train">{$t`Log one now`}</a></p>
	{/if}
{/each}

{#if (sessions.data?.length ?? 0) >= limit}
	<Button variant="outline" class="mx-auto" onclick={() => (limit += 50)}>{$t`Show more`}</Button>
{/if}
