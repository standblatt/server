<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { api, call, type Shooter, type ShooterInput } from './api';
	import * as Dialog from './components/ui/dialog';
	import { Input } from './components/ui/input';
	import FormField from './FormField.svelte';
	import ShooterForm from './ShooterForm.svelte';
	import { clubsQuery, keys } from './tenant';

	let {
		tenantId,
		open = $bindable(false),
		submitLabel,
		oncreated
	}: { tenantId: number; open?: boolean; submitLabel?: string; oncreated: (s: Shooter) => void } = $props();

	const queryClient = useQueryClient();
	const clubs = clubsQuery(() => tenantId);
	let club = $state('');

	/** Creates the shooter with the optional club, refreshes clubs and shooters and closes the dialog. */
	async function create(input: ShooterInput) {
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
		<FormField label={$t`Club (existing or new)`}>
			<Input list="club-names" bind:value={club} placeholder={$t`Club (existing or new)`} />
		</FormField>
		<datalist id="club-names">{#each clubs.data ?? [] as c (c.id)}<option value={c.name}></option>{/each}</datalist>
		{#if open}<ShooterForm submitLabel={submitLabel ?? $t`Create`} onsubmit={create} />{/if}
	</Dialog.Content>
</Dialog.Root>
