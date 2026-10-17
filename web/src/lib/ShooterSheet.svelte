<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Trash from '@lucide/svelte/icons/trash-2';
	import X from '@lucide/svelte/icons/x';
	import { toast } from 'svelte-sonner';
	import { api, call, type ShooterInput } from './api';
	import { Badge } from './components/ui/badge';
	import { Button } from './components/ui/button';
	import * as Sheet from './components/ui/sheet';
	import { ask } from './confirm.svelte';
	import FormError from './FormError.svelte';
	import { formatDate } from './format';
	import ShooterForm from './ShooterForm.svelte';
	import ClubPicker from './ClubPicker.svelte';
	import { keys } from './tenant';

	let { tenantId, shooterId, isAdmin, onclose }: { tenantId: number; shooterId: number; isAdmin: boolean; onclose: () => void } = $props();

	const queryClient = useQueryClient();
	const shooter = createQuery(() => ({
		queryKey: keys.shooter(tenantId, shooterId),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/shooters/{shooterId}', { params: { path: { tenantId, shooterId } } }))
	}));
	const p = $derived({ params: { path: { tenantId, shooterId } } });
	let clubName = $state('');
	let clubReady = $state(true);
	let clubError = $state('');

	/** Invalidates the shooters and clubs of the tenant. */
	async function refresh() {
		await queryClient.invalidateQueries({ queryKey: keys.shooters(tenantId) });
		await queryClient.invalidateQueries({ queryKey: keys.clubs(tenantId) });
	}

	/** Saves the shooter's master data. */
	async function save(input: ShooterInput) {
		await call(api.PUT('/api/tenants/{tenantId}/shooters/{shooterId}', { ...p, body: input }));
		await refresh();
		toast.success($t`Shooter saved`);
	}

	/** Adds the shooter to the club with the typed name. */
	async function addClub(e: SubmitEvent) {
		e.preventDefault();
		clubError = '';
		if (!clubReady) {
			clubError = $t`Pick an existing club or confirm the new one.`;
			return;
		}
		try {
			await call(api.POST('/api/tenants/{tenantId}/shooters/{shooterId}/clubs', { ...p, body: { clubName: clubName.trim() } }));
			clubName = '';
			await refresh();
		} catch (err) {
			clubError = (err as Error).message;
		}
	}

	/** Removes the shooter from a club. */
	async function removeClub(clubId: number) {
		await call(api.DELETE('/api/tenants/{tenantId}/shooters/{shooterId}/clubs/{clubId}', { params: { path: { tenantId, shooterId, clubId } } }));
		await refresh();
	}

	/** Asks for confirmation, then deletes the shooter and closes the sheet. */
	async function remove() {
		const name = `${shooter.data?.firstname} ${shooter.data?.lastname}`;
		if (!(await ask({ title: $t`Delete shooter?`, description: $t`${name} is deleted. Shooters with registrations cannot be deleted.`, action: $t`Delete`, destructive: true }))) return;
		try {
			await call(api.DELETE('/api/tenants/{tenantId}/shooters/{shooterId}', p));
			await refresh();
			toast.success($t`${name} deleted`);
			onclose();
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
</script>

<Sheet.Root open onOpenChange={(o) => !o && onclose()}>
	<Sheet.Content class="gap-0 p-0 data-[side=right]:w-full data-[side=right]:sm:max-w-xl">
		{#if shooter.data}
			{@const sh = shooter.data}
			<Sheet.Header class="border-b px-6 py-4">
				<Sheet.Title class="text-lg">{sh.lastname} {sh.firstname}</Sheet.Title>
				<Sheet.Description>{$t`License ${sh.licenseNumber || '–'} · born ${sh.birthday ? formatDate(sh.birthday) : '–'}`}</Sheet.Description>
			</Sheet.Header>
			<div class="flex-1 space-y-8 overflow-y-auto px-6 py-5">
				<section>
					<h3 class="mb-3 text-sm font-semibold">{$t`Clubs`}</h3>
					<div class="mb-3 flex flex-wrap gap-2">
						{#each sh.clubs as c (c.id)}
							<Badge variant="secondary" class="gap-1 py-1 pr-1">
								{c.name}
								<button class="rounded-sm p-0.5 hover:bg-background" title={$t`Remove ${c.name}`} aria-label={$t`Remove ${c.name}`} onclick={() => removeClub(c.id)}><X class="size-3" /></button>
							</Badge>
						{:else}
							<span class="text-sm text-muted-foreground">{$t`No club assigned.`}</span>
						{/each}
					</div>
					<form class="grid grid-cols-[1fr_auto] items-start gap-2 *:min-w-0" onsubmit={addClub}>
						<div><ClubPicker {tenantId} id="sheet-club" bind:value={clubName} bind:ready={clubReady} required /></div>
						<Button type="submit" variant="outline">{$t`Add club`}</Button>
					</form>
					<FormError message={clubError} class="mt-2" />
				</section>
				<section>
					<h3 class="mb-3 text-sm font-semibold">{$t`Master data`}</h3>
					{#key sh.id}<ShooterForm initial={sh} submitLabel={$t`Save`} onsubmit={save} />{/key}
				</section>
			</div>
			{#if isAdmin}
				<div class="flex border-t px-6 py-3">
					<Button variant="destructive" onclick={remove}><Trash />{$t`Delete shooter`}</Button>
				</div>
			{/if}
		{:else if shooter.error}
			<p class="p-6 text-sm text-destructive">{shooter.error.message}</p>
		{/if}
	</Sheet.Content>
</Sheet.Root>
