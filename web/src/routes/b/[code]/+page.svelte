<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import CircleCheck from '@lucide/svelte/icons/circle-check';
	import { api, call, type BookingRequestInput } from '$lib/api';
	import { deadlineDay } from '$lib/modules/reservations/booking/booking';
	import BookingForm from '$lib/modules/reservations/booking/parts/BookingForm.svelte';
	import PublicShell from '$lib/PublicShell.svelte';
	import { formatDays } from '$lib/format';

	const code = $derived(page.params.code ?? '');
	const info = createQuery(() => ({
		queryKey: ['booking', code],
		queryFn: () => call(api.GET('/api/booking/{code}', { params: { path: { code } } })),
		retry: false
	}));
	let sentTo = $state('');

	/** Sends the request; the page then shows where the confirmation went. */
	async function send(input: BookingRequestInput) {
		await call(api.POST('/api/booking/{code}/requests', { params: { path: { code } }, body: input }));
		sentTo = input.email;
		window.scrollTo({ top: 0 });
	}
</script>

<svelte:head><title>{info.data ? `${info.data.event} – ${$t`Book lanes`}` : $t`Book lanes`}</title></svelte:head>

<PublicShell>
	{#if info.isError}
		<p class="rounded-lg border border-dashed p-6 text-center text-muted-foreground">{$t`This booking link is not valid. Ask the organizer for a new one.`}</p>
	{:else if info.data}
		{@const ev = info.data}
		<header class="grid gap-1">
			<p class="text-sm text-muted-foreground">{ev.organizer}</p>
			<h1 class="text-2xl font-semibold">{ev.event}</h1>
			<p class="text-muted-foreground">{formatDays(ev.days)}</p>
		</header>
		{#if sentTo}
			<div class="flex gap-3 rounded-lg border border-emerald-500/40 bg-emerald-500/10 p-4" role="status">
				<CircleCheck class="mt-0.5 size-5 shrink-0 text-emerald-600" />
				<div class="grid gap-1 text-sm">
					<p class="font-medium">{$t`Request received`}</p>
					<p>{$t`We sent a confirmation to ${sentTo} with a link to check, change or cancel your request. The lanes are reserved once the organizer confirms them.`}</p>
				</div>
			</div>
		{:else if !ev.open}
			<p class="rounded-lg border border-dashed p-6 text-center text-muted-foreground">{$t`The deadline for online booking has passed. Contact the organizer.`}</p>
		{:else}
			<p class="text-sm">
				{$t`Ask for lanes until ${deadlineDay(ev.deadline, ev.timezone)}. The organizer confirms your request by mail.`}
				<a href="/help/user-guide/online-booking#club" target="_blank" rel="noopener" class="underline">{$t`How does this work?`}</a>
			</p>
			<BookingForm info={ev} submitLabel={$t`Send request`} draftKey={`standblatt-booking-${code}`} onsubmit={send} />
		{/if}
	{/if}
</PublicShell>
