<script lang="ts">
	import { fmt } from './my.svelte';

	type Point = { perShot: number; moving: number; kind: string; at: string };
	let { points, goal, label }: { points: Point[]; goal?: number | null; label: string } = $props();

	const W = innerWidth < 640 ? 360 : 640;
	const H = W < 640 ? 200 : 220;
	const pad = { l: 40, r: 12, t: 12, b: 18 };

	const range = $derived.by(() => {
		const values = points.flatMap((p) => [p.perShot, p.moving]);
		if (goal) values.push(goal);
		let lo = Math.min(...values);
		let hi = Math.max(...values);
		if (hi - lo < 0.4) {
			const mid = (hi + lo) / 2;
			lo = mid - 0.2;
			hi = mid + 0.2;
		}
		const m = (hi - lo) * 0.1;
		return { lo: lo - m, hi: hi + m };
	});
	/** Horizontal SVG position of point i (centred when there is only one). */
	const x = (i: number) => pad.l + (points.length < 2 ? (W - pad.l - pad.r) / 2 : (i * (W - pad.l - pad.r)) / (points.length - 1));
	/** Vertical SVG position of a value within the chart range. */
	const y = (v: number) => pad.t + ((range.hi - v) * (H - pad.t - pad.b)) / (range.hi - range.lo);
	const ticks = $derived([0, 0.5, 1].map((f) => range.lo + f * (range.hi - range.lo)));
	const moving = $derived(points.map((p, i) => `${i ? 'L' : 'M'}${x(i)} ${y(p.moving)}`).join(''));
</script>

<svg viewBox="0 0 {W} {H}" class="w-full" role="img" aria-label={label}>
	{#each ticks as t (t)}
		<line x1={pad.l} x2={W - pad.r} y1={y(t)} y2={y(t)} class="stroke-border" stroke-width="1" />
		<text x={pad.l - 6} y={y(t)} text-anchor="end" dominant-baseline="central" class="fill-muted-foreground text-[11px]">{fmt(t, 2)}</text>
	{/each}
	{#if goal}
		<line x1={pad.l} x2={W - pad.r} y1={y(goal)} y2={y(goal)} stroke="#10b981" stroke-width="1.5" stroke-dasharray="6 4" />
	{/if}
	{#if points.length > 1}
		<path d={moving} fill="none" class="stroke-chart-3" stroke-width="2.5" stroke-linejoin="round" />
	{/if}
	{#each points as p, i (i)}
		<circle cx={x(i)} cy={y(p.perShot)} r={p.kind === 'competition' ? 5 : 3.5} class={p.kind === 'competition' ? 'fill-amber-500' : 'fill-foreground/70'}>
			<title>{fmt(p.perShot, 2)}</title>
		</circle>
	{/each}
</svg>
