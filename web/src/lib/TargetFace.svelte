<script lang="ts">
	import { ringRadius, type Target } from './targets';

	type Shot = { x: number; y: number; current?: boolean };
	let {
		target,
		shots = [],
		group,
		zoom = false,
		ontap,
		label,
		class: className = ''
	}: {
		target: Target;
		shots?: Shot[];
		group?: { x: number; y: number; radius: number } | null;
		zoom?: boolean;
		ontap?: (x: number, y: number) => void;
		label: string;
		class?: string;
	} = $props();

	const paper = '#f3efe4';
	const ink = '#14181c';
	const outer = $derived(ringRadius(target, 1));
	const spread = $derived(Math.max(0, ...shots.map((s) => Math.hypot(s.x, s.y) + target.bullet * 2), group ? Math.hypot(group.x, group.y) + group.radius : 0));
	const view = $derived((zoom ? Math.min(outer, Math.max(target.black + target.ring, spread)) : outer) * 1.04);
	const line = $derived(view / 260);
	const rings = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
	let svg = $state<SVGSVGElement>();

	/** Converts a pointer position into target coordinates (y up, rounded to 0.01) and passes them to ontap. */
	function tap(e: PointerEvent) {
		if (!ontap || !svg) return;
		const m = svg.getScreenCTM();
		if (!m) return;
		const p = new DOMPoint(e.clientX, e.clientY).matrixTransform(m.inverse());
		ontap(Math.round(p.x * 100) / 100, Math.round(-p.y * 100) / 100);
	}
</script>

<svg
	bind:this={svg}
	viewBox="{-view} {-view} {2 * view} {2 * view}"
	class="aspect-square w-full touch-none select-none {ontap ? 'cursor-crosshair' : ''} {className}"
	role={ontap ? 'button' : 'img'}
	aria-label={label}
	onpointerdown={tap}
>
	<circle r={outer} fill={paper} stroke={ink} stroke-width={line} />
	<circle r={target.black} fill={ink} />
	{#each rings as n (n)}
		{@const r = ringRadius(target, n)}
		<circle r={r} fill={n === 10 && r < target.black ? '#e9e4d6' : 'none'} stroke={r < target.black ? paper : ink} stroke-width={line} />
		{#if n < 9 && (!zoom || r <= view)}
			<text
				y={-(r - target.ring / 2)}
				fill={r - target.ring < target.black ? paper : ink}
				font-size={target.ring * 0.42}
				text-anchor="middle"
				dominant-baseline="central"
				opacity="0.75">{n}</text
			>
		{/if}
	{/each}
	{#if group}
		<circle cx={group.x} cy={-group.y} r={Math.max(group.radius, line * 4)} fill="#f59e0b" fill-opacity="0.18" stroke="#f59e0b" stroke-width={line * 1.5} stroke-dasharray="{line * 6} {line * 4}" />
		<path
			d="M{group.x - line * 10} {-group.y}h{line * 20}M{group.x} {-group.y - line * 10}v{line * 20}"
			stroke="#f59e0b"
			stroke-width={line * 2.5}
		/>
	{/if}
	{#each shots as s, i (i)}
		<circle
			cx={s.x}
			cy={-s.y}
			r={target.bullet}
			fill={s.current ? '#f59e0b' : '#3b82f6'}
			fill-opacity={s.current ? 0.95 : 0.75}
			stroke="#fff"
			stroke-width={line}
		/>
	{/each}
</svg>
