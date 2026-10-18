<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { membershipOf, meQuery } from '$lib/me';
	import { hasClub } from '$lib/modules';

	let { children } = $props();
	const me = meQuery();
	const tenantId = $derived(Number(page.params.tenantId));
	const club = $derived(membershipOf(me.data, tenantId));
</script>

{#if club}
	{#if hasClub(club, 'forms')}
		{@render children()}
	{:else}
		<p class="text-sm text-muted-foreground">
			{$t`This module is off.`} <a class="underline" href="/t/{tenantId}/modules">{$t`Club modules`}</a>
		</p>
	{/if}
{/if}
