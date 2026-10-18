<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import Search from '@lucide/svelte/icons/search';
	import type { LayoutElement, LibraryImage } from '../api';
	import { Input } from '../components/ui/input';

	let {
		samples,
		images = [],
		library,
		onadd
	}: { samples: Record<string, string>; images?: LibraryImage[]; library: string; onadd: (el: Partial<LayoutElement>) => void } = $props();

	type Entry = { key: string; label: string; example: string; el: Partial<LayoutElement>; field?: boolean };

	let q = $state('');

	/** An entry that prints one placeholder as text. */
	const field = (key: string, label: string): Entry => ({ key, label, example: samples[key] ?? '', el: { type: 'text', value: `{{${key}}}`, size: 10 }, field: true });

	const series = $derived(
		Object.keys(samples)
			.map((k) => /^s(\d+)$/.exec(k)?.[1])
			.filter((n): n is string => !!n)
			.sort((a, b) => Number(a) - Number(b))
	);

	const groups = $derived.by((): { label: string; entries: Entry[] }[] => [
		{
			label: $t`Elements`,
			entries: [
				{ key: 'text', label: $t`Own text`, example: '', el: { type: 'text', value: $t`Text`, size: 10 } },
				{ key: 'qr', label: $t`QR code`, example: '{{registration}}', el: { type: 'qr', value: '{{registration}}', w: 15 } },
				{ key: 'barcode', label: $t`Barcode`, example: '{{shooter.license}}', el: { type: 'barcode', value: '{{shooter.license}}', w: 40, h: 8 } },
				{ key: 'line', label: $t`Line`, example: '', el: { type: 'line', w: 50, h: 0 } }
			]
		},
		{
			label: $t`Shooter`,
			entries: [
				field('shooter.name', $t`Name`),
				field('shooter.lastname', $t`Last name`),
				field('shooter.firstname', $t`First name`),
				field('shooter.license', $t`License number`),
				field('shooter.birthday', $t`Birthday`),
				field('shooter.birthyear', $t`Year of birth`),
				field('shooter.postalcode', $t`Postal code`),
				field('shooter.city', $t`City`),
				field('category', $t`Age category`)
			]
		},
		{
			label: $t`Club and event`,
			entries: [
				field('club', $t`Club`),
				field('club.ssv', $t`SSV number`),
				field('event.name', $t`Event`),
				field('event.date', $t`Date`),
				field('event.organizer', $t`Organizer`),
				field('today', $t`Today`),
				field('programme', $t`Programme`),
				field('lane', $t`Lane`),
				field('registration', $t`Registration number`)
			]
		},
		{
			label: $t`Result`,
			entries: [
				field('total', $t`Total`),
				field('mouches', $t`Mouches`),
				field('award', $t`Award`),
				field('award.value', $t`Value of the award`),
				field('award.value.words', $t`Value in words`),
				...series.flatMap((n) => [field(`s${n}`, $t`Series ${n}`), field(`s${n}.mouches`, $t`Mouches series ${n}`)])
			]
		},
		{
			label: $t`Registration`,
			entries: [
				field('category.short', $t`Age category, short`),
				field('programmes', $t`Programmes`),
				field('passes', $t`Passes`),
				{ key: 'registration.code', label: $t`QR code of the registration`, example: '', el: { type: 'qr', value: '{{registration.code}}', w: 20 }, field: true }
			]
		},
		{
			label: $t`Rank`,
			entries: [field('rank', $t`Rank`), field('ranking', $t`Ranking`), field('class', $t`Class`)]
		},
		{
			label: $t`Single shots`,
			entries: series.map((n) => {
				const shots = Object.keys(samples).filter((k) => k.startsWith(`s${n}.shot`)).length;
				return { key: `s${n}.shot#`, label: $t`Shots of series ${n}`, example: $t`${shots} boxes`, el: { type: 'text', value: `{{s${n}.shot#}}`, w: 7, size: 10, align: 'C', repeat: shots, dx: 8 } };
			})
		},
		{
			label: $t`Images`,
			entries: images.map((i) => ({
				key: `image-${i.id}`,
				label: i.name,
				example: `${i.width} × ${i.height} px`,
				el: { type: 'image', image: i.id, w: 30, h: Math.round((300 * i.height) / i.width) / 10 }
			}))
		}
	]);

	const shown = $derived(
		groups
			.map((g) => ({ ...g, entries: g.entries.filter((e) => (!e.field || e.key in samples) && (!q || `${e.label} ${e.key}`.toLowerCase().includes(q.trim().toLowerCase()))) }))
			.filter((g) => g.entries.length || (g.label === $t`Images` && !q))
	);

	/** Puts the element of an entry into a drag, so the sheet can drop it where the pointer is. */
	function dragstart(e: DragEvent, entry: Entry) {
		e.dataTransfer?.setData('text/standblatt-element', JSON.stringify(entry.el));
		if (e.dataTransfer) e.dataTransfer.effectAllowed = 'copy';
	}
</script>

<div class="grid content-start gap-2">
	<div class="relative">
		<Search class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
		<Input class="h-8 pl-8" type="search" placeholder={$t`Search fields`} bind:value={q} />
	</div>
	{#each shown as g (g.label)}
		<details class="rounded-md border" open>
			<summary class="cursor-pointer px-2.5 py-1.5 text-xs font-semibold text-muted-foreground uppercase">{g.label}</summary>
			<div class="grid gap-0.5 p-1">
				{#if g.label === $t`Images` && !g.entries.length}
					<p class="px-2 py-1 text-xs text-muted-foreground">{$t`No images yet.`} <a class="underline" href={library}>{$t`Upload images`}</a></p>
				{/if}
				{#each g.entries as entry (entry.key)}
					<button
						type="button"
						draggable="true"
						class="flex cursor-grab items-baseline justify-between gap-2 rounded px-2 py-1 text-left text-sm hover:bg-muted active:cursor-grabbing"
						title={$t`Click to add, or drag onto the sheet`}
						ondragstart={(e) => dragstart(e, entry)}
						onclick={() => onadd(entry.el)}
					>
						<span class="truncate">{entry.label}</span>
						{#if entry.example}<span class="shrink-0 truncate text-xs text-muted-foreground">{entry.example}</span>{/if}
					</button>
				{/each}
			</div>
		</details>
	{:else}
		<p class="text-sm text-muted-foreground">{$t`No field matches.`}</p>
	{/each}
</div>
