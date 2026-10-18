<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { createQuery } from '@tanstack/svelte-query';
	import Award from '@lucide/svelte/icons/award';
	import { api, call, type Event } from '../api';
	import { Button } from '../components/ui/button';
	import * as Dialog from '../components/ui/dialog';
	import { Input } from '../components/ui/input';
	import * as NativeSelect from '../components/ui/native-select';
	import FormField from '../FormField.svelte';
	import { meQuery } from '../me';
	import { hasClub } from '../modules';
	import PrintButton from '../PrintButton.svelte';
	import type { CertificatesTarget } from '../print';

	let {
		tenantId,
		event,
		title,
		programmeId,
		classificationId,
		classes = []
	}: { tenantId: number; event: Event; title: string; programmeId?: number; classificationId?: number; classes?: string[] } = $props();
	const me = meQuery();
	const forms = $derived(hasClub(me.data?.tenants.find((m) => m.tenantId === tenantId), 'forms'));
	const layouts = createQuery(() => ({
		queryKey: ['tenants', tenantId, 'layouts'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/layouts', { params: { path: { tenantId } } })),
		enabled: forms
	}));
	const rankForms = $derived((layouts.data ?? []).filter((l) => l.source === 'rank'));
	let open = $state(false);
	let layoutId = $state(0);
	let cls = $state('');
	let from = $state(1);
	let to = $state(3);
	const target = $derived<CertificatesTarget>({
		document: 'certificate',
		eventId: event.id,
		layoutId: layoutId || rankForms[0]?.id || 0,
		programmeId,
		classificationId,
		class: cls || undefined,
		from: Number(from) || 1,
		to: Number(to) || 1
	});
</script>

{#if forms && rankForms.length}
	<Button size="sm" variant="outline" onclick={() => (open = true)} aria-label={$t`Certificates of ${title}`}><Award />{$t`Certificates`}</Button>
	<Dialog.Root bind:open>
		<Dialog.Content class="sm:max-w-md">
			<Dialog.Header>
				<Dialog.Title>{$t`Certificates: ${title}`}</Dialog.Title>
				<Dialog.Description>{$t`One page per ranked shooter, printed onto a print form with the rank as data source.`}</Dialog.Description>
			</Dialog.Header>
			<div class="grid gap-4">
				<FormField label={$t`Print form`}>
					<NativeSelect.Root class="w-full" value={target.layoutId} onchange={(e) => (layoutId = Number(e.currentTarget.value))}>
						{#each rankForms as l (l.id)}<NativeSelect.Option value={l.id}>{l.name}</NativeSelect.Option>{/each}
					</NativeSelect.Root>
				</FormField>
				{#if classes.length > 1}
					<FormField label={$t`Class`}>
						<NativeSelect.Root class="w-full" bind:value={cls}>
							<NativeSelect.Option value="">{$t`All classes`}</NativeSelect.Option>
							{#each classes as c (c)}<NativeSelect.Option value={c}>{c}</NativeSelect.Option>{/each}
						</NativeSelect.Root>
					</FormField>
				{/if}
				<div class="grid grid-cols-2 gap-3">
					<FormField label={$t`From rank`}><Input type="number" min="1" bind:value={from} /></FormField>
					<FormField label={$t`To rank`}><Input type="number" min="1" bind:value={to} /></FormField>
				</div>
				<div class="flex justify-end"><PrintButton {tenantId} {target} /></div>
			</div>
		</Dialog.Content>
	</Dialog.Root>
{/if}
