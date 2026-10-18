<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { goto } from '$app/navigation';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Building from '@lucide/svelte/icons/building-2';
	import Dumbbell from '@lucide/svelte/icons/dumbbell';
	import LogOut from '@lucide/svelte/icons/log-out';
	import UserCog from '@lucide/svelte/icons/user-cog';
	import { roleLabel } from '$lib/labels';
	import { meQuery } from '$lib/me';
	import { logout } from '$lib/my.svelte';
	import Wordmark from '$lib/Wordmark.svelte';

	const me = meQuery();
	const queryClient = useQueryClient();

	$effect(() => {
		// Without a club there is only the training log, so go there straight away.
		if (me.data?.tenants.length === 0) goto('/my', { replaceState: true });
	});
</script>

<main class="flex min-h-svh items-center justify-center bg-sidebar p-6">
	<div class="w-full max-w-3xl">
		<div class="mb-6 flex items-center gap-3 text-white">
			<Wordmark class="text-3xl" />
		</div>
		{#if me.isPending}
			<p class="text-sm text-sidebar-foreground">{$t`Loading …`}</p>
		{:else if me.data && me.data.tenants.length > 0}
			<h1 class="text-2xl font-semibold text-white">{$t`Where to?`}</h1>
			<p class="mb-6 text-sm text-sidebar-foreground">{$t`Logged in as ${me.data.user.name}`}</p>
			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each me.data.tenants as m (m.tenantId)}
					<a class="flex flex-col gap-3 rounded-lg border bg-card p-5 shadow-xl transition hover:-translate-y-0.5 hover:ring-2 hover:ring-sidebar-primary" href="/t/{m.tenantId}">
						<Building class="size-8 text-primary" />
						<span class="font-semibold">{m.tenantName}</span>
						<span class="text-xs text-muted-foreground">{m.roleName || $t(roleLabel(m.roleKey))}</span>
					</a>
				{/each}
				<a class="flex flex-col gap-3 rounded-lg border bg-card p-5 shadow-xl transition hover:-translate-y-0.5 hover:ring-2 hover:ring-sidebar-primary" href="/my">
					<Dumbbell class="size-8 text-primary" />
					<span class="font-semibold">{$t`My training`}</span>
					<span class="text-xs text-muted-foreground">{$t`Your personal training log`}</span>
				</a>
			</div>
			<div class="mt-6 flex flex-wrap gap-4 text-sm text-sidebar-foreground">
				<a href="/account" class="inline-flex items-center gap-1.5 hover:text-white"><UserCog class="size-4" />{$t`Account`}</a>
				<button type="button" class="inline-flex items-center gap-1.5 hover:text-white" onclick={() => logout(queryClient)}><LogOut class="size-4" />{$t`Log out`}</button>
			</div>
		{/if}
	</div>
</main>
