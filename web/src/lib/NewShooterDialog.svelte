<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call, type Shooter, type ShooterInput } from './api';
	import ClubPicker from './ClubPicker.svelte';
	import * as Dialog from './components/ui/dialog';
	import FormField from './FormField.svelte';
	import ShooterForm from './ShooterForm.svelte';
	import { keys } from './tenant';

	let {
		tenantId,
		open = $bindable(false),
		submitLabel,
		oncreated
	}: { tenantId: number; open?: boolean; submitLabel?: string; oncreated: (s: Shooter) => void } = $props();

	const queryClient = useQueryClient();
	let club = $state('');
	let clubReady = $state(true);

	/** Creates the shooter with the optional club, refreshes clubs and shooters and closes the dialog. */
	async function create(input: ShooterInput) {
		if (!clubReady) throw new Error($t`Pick an existing club or confirm the new one.`);
		const sh = await call(
			api.POST('/api/tenants/{tenantId}/shooters', {
				params: { path: { tenantId } },
				body: { ...input, clubNames: club.trim() ? [club.trim()] : [] }
			})
		);
		await queryClient.invalidateQueries({ queryKey: keys.clubs(tenantId) });
		await queryClient.invalidateQueries({ queryKey: keys.shooters(tenantId) });
		toast.success($t`${sh.firstname} ${sh.lastname} created`);
		club = '';
		open = false;
		oncreated(sh);
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-xl">
		<Dialog.Header>
			<Dialog.Title>{$t`New shooter`}</Dialog.Title>
			<Dialog.Description>{$t`Master data of this club. License number, birthday and address are optional.`}</Dialog.Description>
		</Dialog.Header>
		<FormField label={$t`Club`}>
			<ClubPicker {tenantId} id="new-shooter-club" bind:value={club} bind:ready={clubReady} />
		</FormField>
		{#if open}<ShooterForm submitLabel={submitLabel ?? $t`Create`} onsubmit={create} />{/if}
	</Dialog.Content>
</Dialog.Root>
