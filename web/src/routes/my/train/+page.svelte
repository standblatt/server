<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import Check from '@lucide/svelte/icons/check';
	import Delete from '@lucide/svelte/icons/delete';
	import Minus from '@lucide/svelte/icons/minus';
	import Plus from '@lucide/svelte/icons/plus';
	import Undo from '@lucide/svelte/icons/undo-2';
	import X from '@lucide/svelte/icons/x';
	import ZoomIn from '@lucide/svelte/icons/zoom-in';
	import ZoomOut from '@lucide/svelte/icons/zoom-out';
	import { api, call, ApiError } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { ask } from '$lib/confirm.svelte';
	import { complete, filled, fromSession, newDraft, seriesTotal, toInput, totals, typeKey, typeShot, type Draft, type Mode } from '$lib/my/capture';
	import Finish from '$lib/my/Finish.svelte';
	import Result from '$lib/my/Result.svelte';
	import { disciplineName, fmt, fmtTotal, flushOutbox, isOffline, maxShot, myKey, myQuery, queueOffline, type Discipline, type TrainingSaved } from '$lib/my.svelte';
	import TargetFace from '$lib/TargetFace.svelte';
	import { score, targets } from '$lib/targets';

	const draftKey = 'standblatt-draft';
	const queryClient = useQueryClient();
	const my = myQuery();
	const disciplines = $derived(my.data?.disciplines ?? []);

	let draft = $state<Draft>();
	let buf = $state('');
	let zoom = $state(true);
	let finishing = $state(false);
	let saving = $state(false);
	let result = $state<TrainingSaved | 'offline'>();

	const discipline = $derived(disciplines.find((d) => d.id === draft?.disciplineId));
	const target = $derived(discipline?.target ? targets[discipline.target] : undefined);
	const cur = $derived(draft?.series[draft.current]);
	const sum = $derived(draft ? totals(draft) : { points: 0, shots: 0 });
	const done = $derived(draft?.series.filter(filled).length ?? 0);
	const max = $derived(discipline ? maxShot(discipline) : 10);
	const decimal = $derived(!!discipline?.decimal);
	const shotButtons = $derived(!decimal && max <= 10 ? Array.from({ length: max + 1 }, (_, i) => max - i) : null);

	/** New draft for a discipline, with the last weapon and place used. */
	const start = (d: Discipline) => ({
		...newDraft(d, my.data?.recent.find((s) => s.disciplineId === d.id)?.weaponId ?? null),
		place: my.data?.recent.find((s) => s.kind === 'training')?.place ?? ''
	});

	onMount(async () => {
		const edit = Number(page.url.searchParams.get('edit'));
		if (edit) {
			const s = await call(api.GET('/api/my/trainings/{id}', { params: { path: { id: edit } } }));
			draft = fromSession(s.session);
			return;
		}
		try {
			const stored = JSON.parse(localStorage.getItem(draftKey) ?? localStorage.getItem('targetlog-draft') ?? 'null') as Draft | null;
			if (stored && !stored.editId && stored.series.some(filled)) {
				draft = stored;
				toast.info($t`Your unfinished session is back.`);
				return;
			}
		} catch {
			// start fresh
		}
	});

	$effect(() => {
		if (draft || !my.data || page.url.searchParams.has('edit')) return;
		const wanted = Number(page.url.searchParams.get('d'));
		const d = disciplines.find((x) => x.id === wanted) ?? disciplines.find((x) => x.id === my.data!.recent[0]?.disciplineId) ?? disciplines[0];
		if (d) draft = start(d);
	});

	$effect(() => {
		if (!draft || draft.editId || result) return;
		const json = JSON.stringify(draft);
		try {
			localStorage.setItem(draftKey, json);
		} catch {
			// not kept
		}
	});

	/** Removes the stored draft. */
	function clearDraft() {
		try {
			localStorage.removeItem(draftKey);
		} catch {
			// nothing stored
		}
	}

	/** Switches the discipline; keeps entered series, otherwise starts a fresh draft. */
	function changeDiscipline(id: number) {
		const d = disciplines.find((x) => x.id === id);
		if (!d || !draft) return;
		if (draft.series.some(filled)) {
			draft.disciplineId = id;
		} else {
			draft = { ...start(d), clientId: draft.clientId, editId: draft.editId };
		}
		if (!d.target && draft.mode === 'target') draft.mode = 'shots';
		buf = '';
	}

	/** Switches the input mode of the current series. */
	function setMode(m: Mode) {
		if (!draft || !cur) return;
		draft.mode = m;
		buf = '';
		if (filled(cur)) return;
		cur.values = m === 'totals' ? null : [];
	}

	/** Makes series i current and picks the mode matching its values. */
	function select(i: number) {
		if (!draft) return;
		draft.current = i;
		buf = '';
		const s = draft.series[i];
		if (s.values) draft.mode = s.values.some((v) => v.x !== undefined) && target ? 'target' : 'shots';
		else if (s.total !== null) draft.mode = 'totals';
		else s.values = draft.mode === 'totals' ? null : [];
	}

	/** Moves to the next unfilled series, or to finishing after the last one. */
	function advance() {
		if (!draft) return;
		const next = draft.series.findIndex((s, i) => i > draft!.current && !filled(s));
		if (next >= 0) select(next);
		else if (draft.current === draft.series.length - 1 && draft.series.every(complete)) finishing = true;
	}

	/** Appends a series and selects it. */
	function addSeries() {
		if (!draft || !discipline) return;
		draft.series.push({ shots: cur?.shots ?? discipline.shotsPerSeries, total: null, values: draft.mode === 'totals' ? null : [] });
		select(draft.series.length - 1);
	}

	/** Removes the current series (at least one stays). */
	function removeSeries() {
		if (!draft || draft.series.length < 2) return;
		draft.series.splice(draft.current, 1);
		select(Math.min(draft.current, draft.series.length - 1));
	}

	/** Changes the shot count of the current series by delta (1 to 100). */
	function shotsPerSeries(delta: number) {
		if (!cur) return;
		cur.shots = Math.min(100, Math.max(1, cur.shots + delta));
		if (cur.total !== null && cur.total > cur.shots * max) cur.total = null;
	}

	/** Handles a keypad key while entering a series total. */
	function totalKey(key: string) {
		if (!cur) return;
		if (key === 'ok') {
			if (buf) {
				cur.total = Number(buf.replace(/\.$/, ''));
				buf = '';
			}
			if (cur.total !== null) advance();
			return;
		}
		if (key === 'back' && !buf && cur.total !== null) {
			buf = String(cur.total);
			cur.total = null;
		}
		const r = typeKey(buf, key, cur.shots * max, decimal);
		buf = r.buf;
		if (r.commit !== undefined) {
			cur.total = r.commit;
			advance();
		} else if (buf) {
			cur.total = null;
		}
	}

	/** Adds a shot to the current series, moving on when the series is full. */
	function addShot(v: number, x?: number, y?: number) {
		if (!cur) return;
		if (!cur.values) cur.values = [];
		if (cur.values.length >= cur.shots) {
			advance();
			if (!draft || draft.series[draft.current] === cur) return;
			return addShot(v, x, y);
		}
		cur.values.push({ v: Math.round(v * 10) / 10, x, y });
		cur.total = null;
		if (cur.values.length >= cur.shots) advance();
	}

	/** Handles a keypad key while typing single shot values. */
	function shotKey(key: string) {
		if (!cur) return;
		if (key === 'back') {
			if (buf) buf = buf.slice(0, -1);
			else undoShot();
			return;
		}
		const r = typeShot(buf, key);
		buf = r.buf;
		if (r.commit !== undefined) addShot(r.commit);
	}

	/** Removes the last shot, going back to the previous series if needed. */
	function undoShot() {
		if (!draft || !cur) return;
		if (cur.values?.length) cur.values.pop();
		else if (draft.current > 0) {
			select(draft.current - 1);
			draft.series[draft.current].values?.pop();
		}
	}

	/** Scores a tap on the target face and adds it as a shot. */
	function tapTarget(x: number, y: number) {
		if (!target) return;
		addShot(score(target, x, y, decimal), x, y);
	}

	/** Maps the hardware keyboard onto the keypad. */
	function onkeydown(e: KeyboardEvent) {
		if (finishing || result || !draft || (e.target as HTMLElement).closest('input, textarea, select')) return;
		const key = e.key === 'Backspace' ? 'back' : e.key === 'Enter' ? 'ok' : e.key === ',' ? '.' : e.key;
		if (!/^([0-9.]|back|ok)$/.test(key)) return;
		e.preventDefault();
		if (draft.mode === 'totals') totalKey(key);
		else if (decimal || !shotButtons) shotKey(key === 'ok' ? '' : key);
	}

	/** Leaves the capture, asking before discarding unsaved values. */
	async function close() {
		if (draft?.series.some(filled) && !draft.editId) {
			if (!(await ask({ title: $t`Discard this session?`, description: $t`The values you entered are lost.`, action: $t`Discard`, destructive: true }))) return;
		}
		clearDraft();
		await goto(draft?.editId ? `/my/log/${draft.editId}` : '/my');
	}

	/** Saves the session (create or update) and shows the result. */
	async function save() {
		if (!draft) return;
		const input = toInput(draft);
		if (!input.series.length) {
			toast.error($t`Enter at least one series.`);
			return;
		}
		saving = true;
		try {
			result = draft.editId
				? await call(api.PUT('/api/my/trainings/{id}', { params: { path: { id: draft.editId } }, body: input }))
				: await call(api.POST('/api/my/trainings', { body: input }));
			clearDraft();
			finishing = false;
			await queryClient.invalidateQueries({ queryKey: myKey });
			if (my.data) flushOutbox(queryClient, my.data.userId);
		} catch (err) {
			if (isOffline(err) && !draft.editId && my.data) {
				queueOffline(input, my.data.userId);
				clearDraft();
				finishing = false;
				result = 'offline';
			} else {
				toast.error((err as ApiError).message);
			}
		} finally {
			saving = false;
		}
	}

	/** Starts another session of the same discipline. */
	function again() {
		const d = discipline;
		result = undefined;
		draft = d ? start(d) : undefined;
	}
</script>

<svelte:head><title>{$t`Training`} · Standblatt</title></svelte:head>
<svelte:window {onkeydown} />

{#if result}
	<Result {result} discipline={discipline} onagain={again} />
{:else if draft && discipline && cur}
	<div class="mx-auto flex h-svh max-w-xl flex-col bg-background">
		<header class="flex items-center gap-2 border-b px-2 py-2">
			<Button variant="ghost" size="icon" onclick={close} aria-label={$t`Close`}><X /></Button>
			<select
				class="min-w-0 flex-1 truncate rounded-md border bg-background px-2 py-1.5 text-sm font-medium"
				value={draft.disciplineId}
				onchange={(e) => changeDiscipline(Number(e.currentTarget.value))}
				aria-label={$t`Discipline`}
			>
				{#each disciplines as d (d.id)}<option value={d.id}>{$t(disciplineName(d))}</option>{/each}
			</select>
			<Button size="sm" disabled={!done} onclick={() => (finishing = true)}><Check />{$t`Finish`}</Button>
		</header>

		<div class="flex items-end justify-between px-4 pt-3">
			<div>
				<div class="text-4xl font-semibold tabular-nums" data-testid="session-total">{fmtTotal(sum.points)}</div>
				<div class="text-xs text-muted-foreground tabular-nums">
					{$plural(sum.shots, { one: '# shot', other: '# shots' })}{#if sum.shots} · Ø {fmt(sum.points / sum.shots, 2)}{/if}
				</div>
			</div>
			<div class="flex rounded-lg bg-muted p-0.5 text-xs font-medium" role="tablist" aria-label={$t`Input`}>
				{#each [['totals', $t`Series`], ['shots', $t`Shots`], ...(target ? [['target', $t`Target`]] : [])] as [m, label] (m)}
					<button
						type="button"
						role="tab"
						aria-selected={draft.mode === m}
						class="rounded-md px-2.5 py-1.5 {draft.mode === m ? 'bg-background shadow-sm' : 'text-muted-foreground'}"
						onclick={() => setMode(m as Mode)}>{label}</button
					>
				{/each}
			</div>
		</div>

		<div class="flex items-center gap-1.5 overflow-x-auto px-4 py-3 [scrollbar-width:none]" role="tablist" aria-label={$t`Series`}>
			{#each draft.series as s, i (i)}
				<button
					type="button"
					role="tab"
					aria-selected={i === draft.current}
					aria-label={$t`Series ${i + 1}`}
					class="flex min-w-12 shrink-0 flex-col items-center rounded-lg border px-2 py-1 {i === draft.current
						? 'border-primary ring-2 ring-primary/30'
						: complete(s)
							? 'bg-muted/60'
							: ''}"
					onclick={() => select(i)}
				>
					<span class="text-[10px] text-muted-foreground">{i + 1}</span>
					<span class="text-sm font-semibold tabular-nums">{filled(s) ? fmtTotal(seriesTotal(s)) : '–'}</span>
				</button>
			{/each}
			<button type="button" class="flex size-10 shrink-0 items-center justify-center rounded-lg border border-dashed text-muted-foreground" onclick={addSeries} aria-label={$t`Add series`}><Plus class="size-4" /></button>
		</div>

		<div class="flex items-center justify-between px-4 text-xs text-muted-foreground">
			<span class="flex items-center gap-1">
				{$t`Series ${draft.current + 1}`} ·
				<button type="button" class="rounded p-1 hover:bg-muted" onclick={() => shotsPerSeries(-1)} aria-label={$t`Fewer shots`}><Minus class="size-3" /></button>
				<span class="tabular-nums">{$plural(cur.shots, { one: '# shot', other: '# shots' })}</span>
				<button type="button" class="rounded p-1 hover:bg-muted" onclick={() => shotsPerSeries(1)} aria-label={$t`More shots`}><Plus class="size-3" /></button>
			</span>
			{#if draft.series.length > 1}
				<button type="button" class="rounded px-1.5 py-1 hover:bg-muted" onclick={removeSeries}>{$t`Remove series`}</button>
			{/if}
		</div>

		<div class="flex min-h-0 flex-1 flex-col overflow-y-auto px-4 pt-2 pb-[max(1rem,env(safe-area-inset-bottom))]">
			<div class="mt-auto grid gap-3">
			{#if draft.mode === 'totals'}
				<div class="flex flex-col items-center justify-center py-2">
					<div class="text-6xl font-semibold tabular-nums" data-testid="series-input">{#if buf || cur.total !== null}{buf || fmtTotal(cur.total!)}{:else}<span class="text-muted-foreground/30">0</span>{/if}</div>
					<div class="mt-1 text-xs text-muted-foreground">{$t`max. ${fmtTotal(cur.shots * max)}`}</div>
				</div>
				<div class="grid grid-cols-3 gap-2" data-testid="keypad">
					{#each ['1', '2', '3', '4', '5', '6', '7', '8', '9', decimal ? '.' : '', '0', 'back'] as k, i (i)}
						{#if k}
							<button type="button" class="h-[clamp(2.75rem,7svh,3.5rem)] rounded-xl bg-muted text-2xl font-medium active:bg-muted-foreground/20" onclick={() => totalKey(k)} aria-label={k === 'back' ? $t`Delete` : k}>
								{#if k === 'back'}<Delete class="mx-auto size-6" />{:else}{k}{/if}
							</button>
						{:else}<span></span>{/if}
					{/each}
				</div>
				<Button size="lg" class="h-12 text-base" onclick={() => totalKey('ok')}>{$t`Next`}</Button>
			{:else}
				<div class="flex flex-wrap justify-center gap-1" aria-label={$t`Shots of this series`}>
					{#each Array.from({ length: cur.shots }) as _, i (i)}
						{@const v = cur.values?.[i]}
						<span class="flex h-8 min-w-8 items-center justify-center rounded-md border px-1 text-sm font-semibold tabular-nums {v ? 'bg-muted' : i === (cur.values?.length ?? 0) ? 'border-primary' : 'text-muted-foreground'}">
							{v ? (decimal ? fmt(v.v, 1) : v.v) : i === (cur.values?.length ?? 0) && buf ? buf : ''}
						</span>
					{/each}
				</div>
				{#if draft.mode === 'target' && target}
					<div class="relative mx-auto w-full max-w-[min(100%,44svh)]">
						<TargetFace
							{target}
							{zoom}
							label={$t`Tap where the shot hit`}
							ontap={tapTarget}
							shots={[
								...draft.series.slice(0, draft.current).flatMap((s) => (s.values ?? []).filter((v) => v.x !== undefined).map((v) => ({ x: v.x!, y: v.y! }))),
								...(cur.values ?? []).filter((v) => v.x !== undefined).map((v) => ({ x: v.x!, y: v.y!, current: true }))
							]}
						/>
						<div class="absolute top-0 right-0 flex gap-1">
							<Button variant="secondary" size="icon-sm" onclick={() => (zoom = !zoom)} aria-label={zoom ? $t`Zoom out` : $t`Zoom in`}>
								{#if zoom}<ZoomOut />{:else}<ZoomIn />{/if}
							</Button>
						</div>
					</div>
					<div class="flex justify-center gap-2">
						<Button variant="outline" onclick={undoShot}><Undo />{$t`Undo`}</Button>
						{#if !complete(cur) && (cur.values?.length ?? 0) > 0}
							<Button variant="outline" onclick={advance}>{$t`Next series`}</Button>
						{/if}
					</div>
				{:else if shotButtons}
					<div class="grid grid-cols-4 gap-2" data-testid="keypad">
						{#each shotButtons as v (v)}
							<button type="button" class="h-[clamp(2.75rem,7svh,3.5rem)] rounded-xl text-2xl font-semibold active:opacity-70 {v >= max - 1 ? 'bg-primary text-primary-foreground' : 'bg-muted'}" onclick={() => addShot(v)}>{v}</button>
						{/each}
						<button type="button" class="h-[clamp(2.75rem,7svh,3.5rem)] rounded-xl bg-muted active:opacity-70" onclick={undoShot} aria-label={$t`Undo`}><Undo class="mx-auto size-6" /></button>
					</div>
				{:else}
					<p class="text-center text-xs text-muted-foreground">{$t`Type 97 for 9.7 and 104 for 10.4.`}</p>
					<div class="grid grid-cols-3 gap-2" data-testid="keypad">
						{#each ['1', '2', '3', '4', '5', '6', '7', '8', '9', '', '0', 'back'] as k, i (i)}
							{#if k}
								<button type="button" class="h-[clamp(2.75rem,7svh,3.5rem)] rounded-xl bg-muted text-2xl font-medium active:bg-muted-foreground/20" onclick={() => shotKey(k)} aria-label={k === 'back' ? $t`Delete` : k}>
									{#if k === 'back'}<Delete class="mx-auto size-6" />{:else}{k}{/if}
								</button>
							{:else}<span></span>{/if}
						{/each}
					</div>
				{/if}
			{/if}
			</div>
		</div>
	</div>
	<Finish bind:open={finishing} {draft} {saving} tags={my.data?.tags ?? []} weapons={my.data?.weapons ?? []} onsave={save} />
{/if}
