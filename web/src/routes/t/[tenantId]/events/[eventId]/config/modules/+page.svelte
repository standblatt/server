<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { Switch } from '$lib/components/ui/switch';
	import HelpButton from '$lib/help/HelpButton.svelte';
	import { has, moduleUI, modulesQuery, type ModuleKey } from '$lib/modules';
	import ModuleSlot from '$lib/modules/ModuleSlot.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { eventQuery, updateEvent } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const event = eventQuery(
		() => tenantId,
		() => Number(page.params.eventId)
	);
	const modules = modulesQuery();
	const list = $derived((modules.data ?? []).filter((m) => m.level === 'event').map((m) => ({ ...m, key: m.key as ModuleKey })));

	/** Why a module cannot be switched now (blocked) and what switching it does besides (note); "" when nothing applies. */
	function state(key: ModuleKey, requires: string[], needsRange: boolean): { blocked: string; note: string } {
		const ev = event.data!;
		/** Translated labels of modules, comma separated. */
		const names = (keys: string[]) => keys.map((k) => $t(moduleUI[k as ModuleKey].label)).join(', ');
		if (ev.archivedAt) return { blocked: $t`The event is closed.`, note: '' };
		if (has(ev, key)) {
			if (ev.modulesInUse.includes(key)) return { blocked: $t`In use: remove its data first.`, note: '' };
			const by = list.filter((m) => m.requires.includes(key) && has(ev, m.key)).map((m) => m.key);
			return { blocked: by.length ? $t`Needed by ${names(by)}.` : '', note: '' };
		}
		if (needsRange && !ev.ranges.length) return { blocked: $t`Needs a range with lanes.`, note: '' };
		const missing = requires.filter((r) => !has(ev, r as ModuleKey));
		return { blocked: '', note: missing.length ? $t`Switches on ${names(missing)} as well.` : '' };
	}

	/** Switches a module on or off. */
	async function toggle(key: ModuleKey, on: boolean) {
		const ev = event.data!;
		await updateEvent(queryClient, tenantId, ev, { modules: on ? [...ev.modules, key] : ev.modules.filter((k) => k !== key) });
		toast.success(on ? $t`${$t(moduleUI[key].label)} switched on` : $t`${$t(moduleUI[key].label)} switched off`);
	}
</script>

<PageHeader help="user-guide/events#modules" title={$t`Modules`} description={$t`Switch on what this event needs; everything else stays out of the way.`} />

{#if event.data}
	{@const ev = event.data}
	<ul class="grid max-w-5xl gap-3 md:grid-cols-2">
		{#each list as m (m.key)}
			{@const ui = moduleUI[m.key]}
			{@const on = has(ev, m.key)}
			{@const st = state(m.key, m.requires, m.needsRange)}
			<li class="flex flex-col gap-3 rounded-lg border bg-card p-4 {on ? '' : 'opacity-80'}" data-testid="module-{m.key}">
				<div class="flex items-start gap-3">
					<ui.icon class="mt-0.5 size-5 shrink-0 {on ? 'text-primary' : 'text-muted-foreground'}" />
					<div class="min-w-0 flex-1 text-sm">
						<div class="flex items-center gap-1 font-medium">{$t(ui.label)}<HelpButton page={ui.help} /></div>
						<p class="text-muted-foreground">{$t(ui.description)}</p>
						{#if st.blocked || st.note}<p class="mt-1 text-xs text-muted-foreground">{st.blocked || st.note}</p>{/if}
					</div>
					<Switch
						checked={on}
						disabled={!!st.blocked}
						aria-label={$t(ui.label)}
						onCheckedChange={(v) => toggle(m.key, v)}
					/>
				</div>
				{#if on && ui.settings && !ev.archivedAt}
					<div class="border-t pt-3"><ModuleSlot load={ui.settings} {tenantId} event={ev} /></div>
				{/if}
			</li>
		{/each}
	</ul>
{/if}
