<script lang="ts">
	import { untrack } from 'svelte';
	import { plural, t } from 'lingui-for-svelte/macro';
	import { jobStatus } from '$lib/labels';
	import { intlLocale } from '$lib/locale.svelte';
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Plus from '@lucide/svelte/icons/plus';
	import Server from '@lucide/svelte/icons/server';
	import { toast } from 'svelte-sonner';
	import { api, call, type PrintAgent, type PrintSetting, type Printer } from '$lib/api';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import * as Table from '$lib/components/ui/table';
	import { ask } from '$lib/confirm.svelte';
	import FormError from '$lib/FormError.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { rangesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const rangeId = $derived(Number(page.params.rangeId));
	const p = $derived({ params: { path: { tenantId, rangeId } } });
	const queryClient = useQueryClient();
	const key = $derived(['tenants', tenantId, 'ranges', rangeId, 'printing']);

	const agents = createQuery(() => ({ queryKey: [...key, 'agents'], queryFn: () => call(api.GET('/api/tenants/{tenantId}/ranges/{rangeId}/print-agents', p)) }));
	const settings = createQuery(() => ({ queryKey: [...key, 'settings'], queryFn: () => call(api.GET('/api/tenants/{tenantId}/ranges/{rangeId}/print-settings', p)) }));
	const jobs = createQuery(() => ({
		queryKey: [...key, 'jobs'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/ranges/{rangeId}/print-jobs', p)),
		refetchInterval: 10000
	}));

	const documents: { id: PrintSetting['document']; label: string }[] = $derived([
		{ id: 'lane_sheet', label: $t`Lane sheet (A5)` },
		{ id: 'control_strip', label: $t`Control strip (label 90 × 29 mm)` },
		{ id: 'result_sheet', label: $t`Result sheet (A5)` },
		{ id: 'wreath_card', label: $t`Wreath card (A6, pre-printed)` },
		{ id: 'ranking', label: $t`Ranking (A4)` }
	]);
	const routedOptions = ['InputSlot', 'PageSize', 'MediaType'];

	const printers = $derived.by(() => {
		const all = new Map<string, Printer>();
		for (const a of agents.data ?? []) for (const pr of a.printers) all.set(pr.name, pr);
		return [...all.values()];
	});
	/** Whether the agent was seen within the last minute. */
	const online = (a: PrintAgent) => !!a.lastSeenAt && Date.now() - new Date(a.lastSeenAt).getTime() < 60000;

	type Draft = { printer: string; options: Record<string, string>; extra: string };
	let drafts = $state<Record<string, Draft>>({});
	$effect(() => {
		const list = settings.data;
		if (!list) return;
		const next: Record<string, Draft> = {};
		for (const d of untrack(() => documents)) {
			const s = list.find((x) => x.document === d.id);
			const options: Record<string, string> = {};
			const extra: string[] = [];
			for (const [k, v] of Object.entries(s?.options ?? {})) {
				if (routedOptions.includes(k)) options[k] = v;
				else extra.push(v ? `${k}=${v}` : k);
			}
			next[d.id] = { printer: s?.printer ?? '', options, extra: extra.join(', ') };
		}
		drafts = next;
	});

	const ranges = rangesQuery(() => tenantId);
	const range = $derived(ranges.data?.find((r) => r.id === rangeId));
	let error = $state('');
	let agentName = $state('');
	let secret = $state('');
	let addingAgent = $state(false);

	/** Saves the printer and options (incl. extra key=value pairs) of a document. */
	async function save(doc: PrintSetting['document']) {
		error = '';
		const d = drafts[doc];
		const options: Record<string, string> = {};
		for (const [k, v] of Object.entries(d.options)) if (v) options[k] = v;
		for (const part of d.extra.split(',').map((x) => x.trim()).filter(Boolean)) {
			const [k, ...v] = part.split('=');
			options[k.trim()] = v.join('=').trim();
		}
		try {
			await call(api.PUT('/api/tenants/{tenantId}/ranges/{rangeId}/print-settings', { ...p, body: { document: doc, printer: d.printer, options } }));
			toast.success($t`Saved.`);
			await queryClient.invalidateQueries({ queryKey: [...key, 'settings'] });
		} catch (err) {
			error = (err as Error).message;
		}
	}

	/** Creates a print agent and shows its secret once. */
	async function createAgent(e: SubmitEvent) {
		e.preventDefault();
		const r = await call(api.POST('/api/tenants/{tenantId}/ranges/{rangeId}/print-agents', { ...p, body: { name: agentName } }));
		secret = r.secret;
		addingAgent = true;
		agentName = '';
		await queryClient.invalidateQueries({ queryKey: [...key, 'agents'] });
	}

	/** Revokes a print agent after confirmation. */
	async function revoke(a: PrintAgent) {
		if (!(await ask({ title: $t`Revoke print agent "${a.name}"?`, description: $t`The agent can no longer fetch jobs.`, action: $t`Revoke`, destructive: true }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/print-agents/{agentId}', { params: { path: { tenantId, agentId: a.id } } }));
		await queryClient.invalidateQueries({ queryKey: [...key, 'agents'] });
	}
</script>

<PageHeader help="user-guide/printing" title={range ? $t`Printing · ${range.name}` : $t`Printing`} description={$t`Without an online agent, documents open as PDF in the browser.`} />

<div class="grid gap-4 *:min-w-0 2xl:grid-cols-[28rem_minmax(0,1fr)]">
	<section class="h-fit overflow-hidden rounded-lg border bg-card">
		<div class="border-b px-4 py-3">
			<h2 class="text-sm font-semibold">{$t`Print agents`}</h2>
			<p class="text-xs text-muted-foreground">
				{$t`The print agent runs on a computer next to the printers (e.g. a Raspberry Pi with CUPS). It connects to Standblatt, reports the printers and prints the jobs of this range.`}
			</p>
		</div>
		<ul class="divide-y">
			{#each agents.data ?? [] as a (a.id)}
				<li class="flex flex-wrap items-center gap-2 px-4 py-2.5 text-sm">
					<Server class="size-4 text-muted-foreground" />
					<span class="font-medium">{a.name}</span>
					<Badge variant="outline" class={online(a) ? 'border-emerald-500/50 text-emerald-700 dark:text-emerald-400' : ''}>
						{online(a) ? $t`online` : a.lastSeenAt ? $t`last seen ${new Date(a.lastSeenAt).toLocaleString(intlLocale())}` : $t`never connected`}
					</Badge>
					<span class="text-xs text-muted-foreground">{$plural(a.printers.length, { one: '# printer', other: '# printers' })}</span>
					<Button variant="ghost" size="sm" class="ml-auto text-destructive" onclick={() => revoke(a)}>{$t`revoke`}</Button>
				</li>
			{:else}
				<li class="px-4 py-3 text-sm text-muted-foreground">{$t`No print agent yet.`}</li>
			{/each}
		</ul>
		<form class="flex gap-2 border-t p-4" onsubmit={createAgent}>
			<Input placeholder={$t`Agent name, e.g. Pi hall 1`} bind:value={agentName} required aria-label={$t`Agent name`} />
			<Button type="submit" variant="outline"><Plus />{$t`Add print agent`}</Button>
		</form>
	</section>

	<div class="space-y-4">
		<section class="overflow-hidden rounded-lg border bg-card">
			<div class="border-b px-4 py-3">
				<h2 class="text-sm font-semibold">{$t`Documents`}</h2>
				{#if printers.length === 0}<p class="text-xs text-muted-foreground">{$t`The printers appear here once an agent has connected.`}</p>{/if}
			</div>
			<Table.Root>
				<Table.Header><Table.Row><Table.Head>{$t`Document`}</Table.Head><Table.Head>{$t`Printer`}</Table.Head><Table.Head>{$t`Options`}</Table.Head><Table.Head></Table.Head></Table.Row></Table.Header>
				<Table.Body>
					{#each documents as d (d.id)}
						{#if drafts[d.id]}
							{@const pr = printers.find((x) => x.name === drafts[d.id].printer)}
							<Table.Row>
								<Table.Cell class="font-medium">{d.label}</Table.Cell>
								<Table.Cell>
									<NativeSelect.Root size="sm" bind:value={drafts[d.id].printer} aria-label={$t`Printer for ${d.label}`}>
										<NativeSelect.Option value="">{$t`browser (PDF)`}</NativeSelect.Option>
										{#each printers as x (x.name)}<NativeSelect.Option value={x.name}>{x.name}</NativeSelect.Option>{/each}
										{#if drafts[d.id].printer && !pr}<NativeSelect.Option value={drafts[d.id].printer}>{$t`${drafts[d.id].printer} (offline)`}</NativeSelect.Option>{/if}
									</NativeSelect.Root>
								</Table.Cell>
								<Table.Cell>
									<div class="flex flex-wrap items-end gap-2">
										{#each (pr?.options ?? []).filter((o) => routedOptions.includes(o.name)) as o (o.name)}
											<label class="grid gap-1 text-xs text-muted-foreground">
												{o.label || o.name}
												<NativeSelect.Root size="sm" bind:value={drafts[d.id].options[o.name]}>
													<NativeSelect.Option value="">{$t`default (${o.default ?? ''})`}</NativeSelect.Option>
													{#each o.values as v (v)}<NativeSelect.Option value={v}>{v}</NativeSelect.Option>{/each}
												</NativeSelect.Root>
											</label>
										{/each}
										<Input class="h-8 w-56" placeholder={$t`more lp options, e.g. fit-to-page`} bind:value={drafts[d.id].extra} aria-label={$t`More options for ${d.label}`} />
									</div>
								</Table.Cell>
								<Table.Cell class="text-right"><Button size="sm" variant="outline" onclick={() => save(d.id)}>{$t`Save`}</Button></Table.Cell>
							</Table.Row>
						{/if}
					{/each}
				</Table.Body>
			</Table.Root>
			<FormError message={error} class="m-4" />
		</section>

		<section class="overflow-hidden rounded-lg border bg-card">
			<div class="border-b px-4 py-3"><h2 class="text-sm font-semibold">{$t`Recent jobs`}</h2></div>
			<Table.Root>
				<Table.Header><Table.Row><Table.Head>{$t`Time`}</Table.Head><Table.Head>{$t`Document`}</Table.Head><Table.Head>{$t`Printer`}</Table.Head><Table.Head>{$t`Status`}</Table.Head></Table.Row></Table.Header>
				<Table.Body>
					{#each jobs.data ?? [] as j (j.id)}
						<Table.Row>
							<Table.Cell class="tabular-nums">{new Date(j.createdAt).toLocaleTimeString(intlLocale())}</Table.Cell>
							<Table.Cell>{j.title}</Table.Cell>
							<Table.Cell>{j.printer}</Table.Cell>
							<Table.Cell>
								<Badge
									variant="outline"
									class={j.status === 'failed' ? 'border-destructive/50 text-destructive' : j.status === 'done' ? 'border-emerald-500/50 text-emerald-700 dark:text-emerald-400' : ''}>{$t(jobStatus(j.status))}</Badge
								>
								{#if j.error}<span class="ml-2 text-xs text-destructive">{j.error}</span>{/if}
							</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row><Table.Cell colspan={4} class="py-6 text-center text-muted-foreground">{$t`No print jobs yet.`}</Table.Cell></Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</section>
	</div>
</div>

<Dialog.Root bind:open={addingAgent}>
	<Dialog.Content class="sm:max-w-xl">
		<Dialog.Header>
			<Dialog.Title>{$t`Print agent added`}</Dialog.Title>
			<Dialog.Description>{$t`Start the agent with these settings. The token is shown only now.`}</Dialog.Description>
		</Dialog.Header>
		<pre class="overflow-x-auto rounded-md bg-muted p-3 text-xs" role="status">STANDBLATT_URL={location.origin}
STANDBLATT_AGENT_TOKEN=<span data-testid="agent-token">{secret}</span></pre>
		<Dialog.Footer>
			<Button
				variant="outline"
				onclick={() => {
					navigator.clipboard.writeText(`STANDBLATT_URL=${location.origin}\nSTANDBLATT_AGENT_TOKEN=${secret}`);
					toast.success($t`Copied`);
				}}>{$t`Copy`}</Button
			>
			<Button onclick={() => (addingAgent = false)}>{$t`Done`}</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
