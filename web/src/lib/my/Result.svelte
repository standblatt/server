<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import ArrowDownRight from '@lucide/svelte/icons/arrow-down-right';
	import ArrowUpRight from '@lucide/svelte/icons/arrow-up-right';
	import CloudOff from '@lucide/svelte/icons/cloud-off';
	import Medal from '@lucide/svelte/icons/medal';
	import Trophy from '@lucide/svelte/icons/trophy';
	import { Button } from '$lib/components/ui/button';
	import { disciplineName, fmt, fmtTotal, type Discipline, type TrainingSaved } from '$lib/my.svelte';

	let { result, discipline, onagain }: { result: TrainingSaved | 'offline'; discipline?: Discipline; onagain: () => void } = $props();
</script>

<main class="mx-auto flex min-h-svh max-w-md flex-col items-center justify-center gap-6 p-6 text-center">
	{#if result === 'offline'}
		<CloudOff class="size-12 text-muted-foreground" />
		<div>
			<h1 class="text-xl font-semibold">{$t`Saved on this device`}</h1>
			<p class="mt-1 text-sm text-muted-foreground">{$t`You are offline. The session is uploaded automatically as soon as you are back online.`}</p>
		</div>
	{:else}
		{@const s = result.session}
		{@const f = result.feedback}
		{#if f.personalBest}
			<div class="flex size-20 items-center justify-center rounded-full bg-amber-500/15"><Trophy class="size-10 text-amber-500" /></div>
		{:else}
			<div class="flex size-20 items-center justify-center rounded-full bg-primary/10"><Medal class="size-10 text-primary" /></div>
		{/if}
		<div>
			<p class="text-sm text-muted-foreground">{$t(disciplineName(discipline))}</p>
			<p class="text-6xl font-semibold tabular-nums" data-testid="result-total">{fmtTotal(s.total)}</p>
			<p class="text-sm text-muted-foreground tabular-nums">{$t`${s.shots} shots · Ø ${fmt(s.shots ? s.total / s.shots : 0, 2)} per shot`}</p>
		</div>
		<div class="grid w-full gap-2 text-sm" data-testid="feedback">
			{#if f.personalBest}
				<p class="rounded-lg bg-amber-500/15 px-4 py-3 font-semibold text-amber-700 dark:text-amber-300">{$t`New personal best!`}</p>
			{:else if f.of > 1}
				<p class="rounded-lg bg-muted px-4 py-3">{$t`Place ${f.rank} of your ${f.of} sessions with ${s.shots} shots`}</p>
			{:else}
				<p class="rounded-lg bg-muted px-4 py-3">{$t`Your first session in this format — the next one has something to beat.`}</p>
			{/if}
			{#if f.delta != null}
				<p class="flex items-center justify-center gap-1 rounded-lg px-4 py-3 {f.delta >= 0 ? 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-400' : 'bg-muted'}">
					{#if f.delta >= 0}<ArrowUpRight class="size-4" />{:else}<ArrowDownRight class="size-4" />{/if}
					{f.delta >= 0
						? $t`${fmt(f.delta, 2)} per shot better than your last 30 days`
						: $t`${fmt(-f.delta, 2)} per shot below your last 30 days`}
				</p>
			{/if}
		</div>
	{/if}
	<div class="grid w-full gap-2">
		<Button href="/my" size="lg">{$t`Done`}</Button>
		<div class="grid grid-cols-2 gap-2">
			<Button variant="outline" onclick={onagain}>{$t`Another session`}</Button>
			{#if result !== 'offline'}
				<Button variant="outline" href="/my/log/{result.session.id}">{$t`Details`}</Button>
			{:else}
				<Button variant="outline" href="/my/log">{$t`Log`}</Button>
			{/if}
		</div>
	</div>
</main>
