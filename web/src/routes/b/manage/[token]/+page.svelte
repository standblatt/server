<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import createClient from 'openapi-fetch';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import Ban from '@lucide/svelte/icons/ban';
	import { toast } from 'svelte-sonner';
	import { call, type BookingRequestInput } from '$lib/api';
	import type { paths } from '$lib/api.gen';
	import { deadlineDay } from '$lib/modules/reservations/booking/booking';
	import BookingForm from '$lib/modules/reservations/booking/parts/BookingForm.svelte';
	import PublicShell from '$lib/PublicShell.svelte';
	import BookingStatus from '$lib/modules/reservations/booking/parts/BookingStatus.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ask } from '$lib/confirm.svelte';
	import { currentLocale } from '$lib/locale.svelte';
	import { formatDays } from '$lib/format';

	const token = page.params.token ?? '';
	const client = createClient<paths>({ baseUrl: '', headers: { Authorization: `Bearer ${token}` } });
	client.use({
		/** Sends the UI language so error messages come translated. */
		onRequest({ request }) {
			request.headers.set('Accept-Language', currentLocale());
		}
	});
	const queryClient = useQueryClient();
	const key = ['booking-manage', token];
	const data = createQuery(() => ({
		queryKey: key,
		queryFn: () => call(client.GET('/api/booking-manage')),
		retry: false
	}));
	const editable = $derived(!!data.data && data.data.info.open && (data.data.request.status === 'pending' || data.data.request.status === 'confirmed'));
	let form = $state(0);

	/** Saves the changed request; a confirmed one with other slots needs a new confirmation. */
	async function save(input: BookingRequestInput) {
		const r = data.data!.request;
		const before = JSON.stringify(r.items.map((i) => [i.rangeId, new Date(i.start).toISOString(), i.lanes]).sort());
		const after = JSON.stringify(input.items.map((i) => [i.rangeId, new Date(i.start).toISOString(), i.lanes]).sort());
		if (
			r.status === 'confirmed' &&
			before !== after &&
			!(await ask({
				title: $t`Change confirmed slots?`,
				description: $t`Your reserved lanes are released and the organizer has to confirm the request again.`,
				action: $t`Change`
			}))
		)
			return;
		const updated = await call(client.PUT('/api/booking-manage', { body: input }));
		toast.success(updated.status === 'pending' && r.status === 'confirmed' ? $t`Changed; waiting for a new confirmation` : $t`Request saved`);
		await queryClient.invalidateQueries({ queryKey: key });
		form++;
	}

	/** Cancels the request after confirmation. */
	async function cancel() {
		if (!(await ask({ title: $t`Cancel request?`, description: $t`Reserved lanes are released. This cannot be undone.`, action: $t`Cancel request`, destructive: true }))) return;
		await call(client.DELETE('/api/booking-manage'));
		toast.success($t`Request cancelled`);
		await queryClient.invalidateQueries({ queryKey: key });
	}
</script>

<svelte:head><title>{$t`Your lane request`}</title></svelte:head>

<PublicShell>
	{#if data.isError}
		<p class="rounded-lg border border-dashed p-6 text-center text-muted-foreground">{$t`This link is not valid.`}</p>
	{:else if data.data}
		{@const { request: r, info } = data.data}
		<header class="grid gap-1">
			<p class="text-sm text-muted-foreground">{info.organizer}</p>
			<h1 class="text-2xl font-semibold">{info.event}</h1>
			<p class="text-muted-foreground">{formatDays(info.days)}</p>
		</header>
		<section class="grid gap-2 rounded-lg border bg-card p-4 text-sm">
			<div class="flex flex-wrap items-center gap-3">
				<span class="font-medium">{$t`Your request`}</span>
				<BookingStatus status={r.status} />
				{#if editable}<Button variant="outline" size="sm" class="ml-auto" onclick={cancel}><Ban />{$t`Cancel request`}</Button>{/if}
			</div>
			{#if r.status === 'declined' && r.declineReason}<p>{r.declineReason}</p>{/if}
			{#if editable}
				<p class="text-muted-foreground">{$t`You can change or cancel the request until ${deadlineDay(info.deadline, info.timezone)}.`}</p>
			{:else if r.status === 'pending' || r.status === 'confirmed'}
				<p class="text-muted-foreground">{$t`The deadline has passed. Contact the organizer for changes.`}</p>
			{/if}
		</section>
		{#key form}
			<BookingForm {info} initial={{ ...r, items: r.items, shooters: r.shooters }} readonly={!editable} submitLabel={$t`Save changes`} onsubmit={save} />
		{/key}
	{/if}
</PublicShell>
