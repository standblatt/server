<script lang="ts">
	import { pdfClick } from '$lib/print';
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import FileSpreadsheet from '@lucide/svelte/icons/file-spreadsheet';
	import FileText from '@lucide/svelte/icons/file-text';
	import { ageCategoryLabel } from '$lib/labels';
	import { formatDate, programmeTitle } from '$lib/format';
	import { intlLocale } from '$lib/locale.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import PageHeader from '$lib/PageHeader.svelte';
	import PanelHead from '$lib/PanelHead.svelte';
	import ModuleSlot from '$lib/modules/ModuleSlot.svelte';
	import { enabled, has } from '$lib/modules';
	import { eventQuery, keys, programmesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const path = $derived({ params: { path: { tenantId, eventId } } });
	const stats = createQuery(() => ({
		queryKey: [...keys.reports(tenantId, eventId), 'statistics'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/statistics', path))
	}));
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const programmes = programmesQuery(
		() => tenantId,
		() => eventId
	);

	const reports = $derived([
		{ id: 'financial', label: $t`Financial report` },
		{ id: 'revenue', label: $t`Revenue by programme and age category` },
		{ id: 'club_participants', label: $t`Club participants` },
		...enabled(event.data).flatMap((m) => (m.reports ?? []).map((r) => ({ id: r.id, label: $t(r.label) })))
	]);
	/** Download URL of a report in the given format. */
	const reportUrl = (id: string, format: string) => `/api/tenants/${tenantId}/events/${eventId}/reports/${id}?format=${format}`;
	/** Formats an amount in Swiss francs. */
	const chf = (v: number) => `CHF ${v.toLocaleString(intlLocale(), { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
	/** Key figures as [value, label] pairs. */
	const figures = (s: NonNullable<typeof stats.data>) => [
		[String(s.registered + s.shooting + s.completed), $t`registrations`],
		[`${s.completionRate}%`, $t`completed`],
		[String(s.entries), $t`programmes sold`],
		[chf(s.revenue), $t`revenue`],
		...(has(event.data, 'awards')
			? [
					[String(s.awards.reduce((n, a) => n + a.count, 0)), $t`awards`],
					[chf(s.awardsValue), $t`award value`]
				]
			: [])
	];
</script>

<PageHeader help="user-guide/statistics-reports" title={$t`Statistics`} description={$t`Key figures, club competition and exports of this event.`} />

<div class="space-y-4">
	{#if stats.data}
		{@const s = stats.data}
		<section class="grid grid-cols-2 divide-x divide-y overflow-hidden rounded-lg border bg-card sm:grid-cols-3 sm:divide-y-0 xl:grid-cols-6" aria-label={$t`Key figures`}>
			{#each figures(s) as [value, label] (label)}
				<div class="px-4 py-3"><div class="text-xl font-semibold tabular-nums">{value}</div><div class="text-xs text-muted-foreground">{label}</div></div>
			{/each}
		</section>
	{/if}

	<div class="grid gap-4 *:min-w-0 xl:grid-cols-2">
		{#if stats.data}
			{@const s = stats.data}
			<section class="overflow-hidden rounded-lg border bg-card">
				<PanelHead title={$t`Programmes`} />
				<Table.Root>
					<Table.Header><Table.Row><Table.Head>{$t`Programme`}</Table.Head><Table.Head class="text-right">{$t({ message: "Entries", comment: "Number of programmes shot" })}</Table.Head><Table.Head class="hidden sm:table-cell text-right">{$t`Completed`}</Table.Head><Table.Head class="text-right">{$t`Revenue`}</Table.Head></Table.Row></Table.Header>
					<Table.Body>
						{#each s.programmes as p (p.programmeId)}
							<Table.Row><Table.Cell class="font-medium whitespace-normal">{p.title}</Table.Cell><Table.Cell class="text-right tabular-nums">{p.entries}</Table.Cell><Table.Cell class="hidden sm:table-cell text-right tabular-nums">{p.completed}</Table.Cell><Table.Cell class="text-right tabular-nums">{chf(p.revenue)}</Table.Cell></Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</section>

			<section class="overflow-hidden rounded-lg border bg-card">
				<PanelHead title={$t`Age categories`} />
				<Table.Root>
					<Table.Body>
						{#each s.categories as c (c.category)}
							<Table.Row><Table.Cell>{$t(ageCategoryLabel(c.category))}</Table.Cell><Table.Cell class="text-right tabular-nums">{c.count}</Table.Cell></Table.Row>
						{:else}
							<Table.Row><Table.Cell class="py-6 text-center text-muted-foreground">{$t`No registrations yet.`}</Table.Cell></Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</section>

			{#if s.days.length > 1}
				<section class="overflow-hidden rounded-lg border bg-card">
					<PanelHead title={$t`Shooting days`} />
					<Table.Root>
						<Table.Header>
							<Table.Row><Table.Head>{$t`Day`}</Table.Head><Table.Head class="text-right">{$t`Registrations`}</Table.Head><Table.Head class="text-right">{$t`Revenue`}</Table.Head></Table.Row>
						</Table.Header>
						<Table.Body>
							{#each s.days as d (d.date)}
								<Table.Row><Table.Cell>{formatDate(d.date)}</Table.Cell><Table.Cell class="text-right tabular-nums">{d.registrations}</Table.Cell><Table.Cell class="text-right tabular-nums">{chf(d.revenue)}</Table.Cell></Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</section>
			{/if}

		{/if}

		{#if event.data}
			{#each enabled(event.data) as m (m.key)}
				{#if m.statisticsSections}<ModuleSlot load={m.statisticsSections} {tenantId} event={event.data} />{/if}
			{/each}
		{/if}

		<section class="overflow-hidden rounded-lg border bg-card" aria-label={$t`Exports`}>
			<PanelHead title={$t`Exports`} hint={$t`PDF for printing, CSV for Excel (semicolon separated).`} />
			<ul class="divide-y">
				{#each reports as r (r.id)}
					<li class="flex items-center gap-2 px-4 py-2 text-sm">
						<span class="mr-auto">{r.label}</span>
						<Button size="sm" variant="outline" href={reportUrl(r.id, 'pdf')} target="_blank" onclick={pdfClick}><FileText />PDF</Button>
						<Button size="sm" variant="outline" href={reportUrl(r.id, 'csv')}><FileSpreadsheet />CSV</Button>
					</li>
				{/each}
				{#each programmes.data ?? [] as p (p.id)}
					{#if p.showOnRankings}
						<li class="flex items-center gap-2 px-4 py-2 text-sm">
							<span class="mr-auto">{$t`Ranking ${programmeTitle(p)}`}</span>
							<Button size="sm" variant="outline" href="/api/tenants/{tenantId}/programmes/{p.id}/ranking" target="_blank" onclick={pdfClick}><FileText />PDF</Button>
							<Button size="sm" variant="outline" href="/api/tenants/{tenantId}/programmes/{p.id}/ranking.csv"><FileSpreadsheet />CSV</Button>
						</li>
					{/if}
				{/each}
			</ul>
		</section>
	</div>
</div>
