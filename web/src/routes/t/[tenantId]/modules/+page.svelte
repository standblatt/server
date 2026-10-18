<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call } from '$lib/api';
	import { meKey, meQuery } from '$lib/me';
	import { clubModuleUI, hasClub, modulesQuery, type ClubModuleKey } from '$lib/modules';
	import ModuleCards from '$lib/modules/ModuleCards.svelte';
	import PageHeader from '$lib/PageHeader.svelte';

	const tenantId = $derived(Number(page.params.tenantId));
	const queryClient = useQueryClient();
	const me = meQuery();
	const club = $derived(me.data?.tenants.find((m) => m.tenantId === tenantId));
	const modules = modulesQuery();
	const list = $derived((modules.data ?? []).filter((m) => m.level === 'club').map((m) => m.key as ClubModuleKey));

	$effect(() => {
		if (club && club.role !== 'admin') goto(`/t/${tenantId}`, { replaceState: true });
	});

	/** Switches a club module on or off. */
	async function toggle(key: ClubModuleKey, on: boolean) {
		const current = club?.modules ?? [];
		await call(api.PUT('/api/tenants/{tenantId}/modules', { params: { path: { tenantId } }, body: { modules: on ? [...current, key] : current.filter((k) => k !== key) } }));
		await queryClient.invalidateQueries({ queryKey: meKey });
		toast.success(on ? $t`${$t(clubModuleUI[key].label)} switched on` : $t`${$t(clubModuleUI[key].label)} switched off`);
	}
</script>

<PageHeader help="user-guide/events#club-modules" title={$t`Club modules`} description={$t`What the club uses besides its events.`} />

{#if club?.role === 'admin'}
	<ModuleCards
		cards={list.map((key) => ({
			key,
			ui: clubModuleUI[key],
			on: hasClub(club, key),
			blocked: hasClub(club, key) && club.modulesInUse.includes(key) ? $t`In use: remove its data first.` : '',
			note: ''
		}))}
		ontoggle={(k, on) => toggle(k as ClubModuleKey, on)}
	/>
{/if}
