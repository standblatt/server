<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { Snippet } from 'svelte';
	import { goto } from '$app/navigation';
	import * as NativeSelect from '$lib/components/ui/native-select';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Archive from '@lucide/svelte/icons/archive';
	import ArchiveRestore from '@lucide/svelte/icons/archive-restore';
	import CopyPlus from '@lucide/svelte/icons/copy-plus';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call, type EventInput } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { ask } from '$lib/confirm.svelte';
	import EventForm from '$lib/EventForm.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import ProgrammeList from '$lib/ProgrammeList.svelte';
	import HelperCodes from '$lib/HelperCodes.svelte';
	import BookingSettings from '$lib/BookingSettings.svelte';
	import KioskTokens from '$lib/KioskTokens.svelte';
	import Classifications from '$lib/Classifications.svelte';
	import { meQuery } from '$lib/me';
	import { clubsQuery, eventQuery, keys, rangesQuery } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const queryClient = useQueryClient();
	const me = meQuery();
	const ranges = rangesQuery(() => tenantId);
	const clubs = clubsQuery(() => tenantId);
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);
	const isAdmin = $derived(me.data?.tenants.find((t) => t.tenantId === tenantId)?.role === 'admin');
	const tabs = $derived([
		['programmes', $t`Programmes`],
		['event', $t`Event`],
		['classifications', $t`Classifications`],
		['kiosk', $t`Kiosk displays`],
		['helpers', $t`Helpers`],
		...(event.data?.ranges.length ? [['booking', $t`Online booking`]] : []),
		['more', $t`Copy & close`]
	]);
	const setupHelp: Record<string, string> = {
		programmes: 'user-guide/programmes',
		event: 'user-guide/events',
		classifications: 'user-guide/classifications',
		kiosk: 'user-guide/rankings-kiosk',
		helpers: 'user-guide/helpers',
		booking: 'user-guide/online-booking',
		more: 'user-guide/events'
	};
	const tab = $derived(isAdmin && tabs.some(([k]) => k === page.url.searchParams.get('tab')) ? page.url.searchParams.get('tab') : 'programmes');

	/** Saves the event settings. */
	async function save(input: EventInput) {
		await call(api.PUT('/api/tenants/{tenantId}/events/{eventId}', { params: { path: { tenantId, eventId } }, body: input }));
		await queryClient.invalidateQueries({ queryKey: keys.events(tenantId) });
		await queryClient.invalidateQueries({ queryKey: keys.lanes(tenantId, eventId) });
		await queryClient.invalidateQueries({ queryKey: keys.event(tenantId, eventId) });
		toast.success($t`Event saved`);
	}

	/** Closes (after confirmation) or reopens the event; closing releases the rankings to the shooters' logs. */
	async function setArchived(archived: boolean) {
		if (
			archived &&
			!(await ask({
				title: $t`Close event?`,
				description: $t`"${event.data?.name ?? ''}" is final: the rankings are released and the results appear in the shooters' personal training logs with their rank. Nothing can be changed until the event is reopened.`,
				action: $t`Close event`
			}))
		)
			return;
		await call(api.PUT('/api/tenants/{tenantId}/events/{eventId}/archived', { params: { path: { tenantId, eventId } }, body: { archived } }));
		await queryClient.invalidateQueries({ queryKey: keys.events(tenantId) });
		await queryClient.invalidateQueries({ queryKey: keys.event(tenantId, eventId) });
		toast.success(archived ? $t`Event closed` : $t`Event reopened`);
	}

	let copy = $state<{ name: string; startsOn: string } | null>(null);
	let copyError = $state('');

	/** Event name for next year: years and a leading edition number count up. */
	const nextName = (name: string) =>
		name.replace(/\b(19|20)\d{2}\b/g, (y) => String(Number(y) + 1)).replace(/^(\d+)\./, (_, n) => `${Number(n) + 1}.`);
	/** Same date one year later (YYYY-MM-DD). */
	const nextYear = (date: string) => `${Number(date.slice(0, 4)) + 1}${date.slice(4)}`;

	/** Opens the copy dialog with next year's name and date. */
	function openCopy() {
		copyError = '';
		copy = { name: nextName(event.data?.name ?? ''), startsOn: nextYear(event.data?.startsOn ?? '') };
	}

	/** Copies the event for next year and opens the copy's setup. */
	async function copyEvent(e: SubmitEvent) {
		e.preventDefault();
		if (!copy) return;
		try {
			const ev = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/copy', { params: { path: { tenantId, eventId } }, body: copy }));
			await queryClient.invalidateQueries({ queryKey: keys.events(tenantId) });
			toast.success($t`Event "${ev.name}" created`);
			copy = null;
			await goto(`/t/${tenantId}/events/${ev.id}/setup`);
		} catch (err) {
			copyError = (err as Error).message;
		}
	}

	/** Deletes the event after confirmation and returns to the events list. */
	async function remove() {
		if (!(await ask({ title: $t`Delete event?`, description: $t`"${event.data?.name ?? ''}" is deleted. Events with registrations cannot be deleted.`, action: $t`Delete`, destructive: true }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/events/{eventId}', { params: { path: { tenantId, eventId } } }));
		await goto(`/t/${tenantId}`);
		queryClient.removeQueries({ queryKey: keys.event(tenantId, eventId) });
		await queryClient.invalidateQueries({ queryKey: keys.events(tenantId) });
		toast.success($t`Event deleted`);
	}
</script>

<PageHeader help={setupHelp[tab ?? 'programmes']} title={$t`Setup`} description={$t`Programmes, event settings and kiosk displays.`} />

{#snippet card(title: string, text: string, action: Snippet, danger = false)}
	<section class="flex flex-wrap items-center gap-3 rounded-lg border bg-card p-4 {danger ? 'border-destructive/30' : ''}">
		<div class="min-w-48 flex-1 text-sm">
			<div class="font-medium">{title}</div>
			<div class="text-muted-foreground">{text}</div>
		</div>
		{@render action()}
	</section>
{/snippet}

{#snippet unarchive()}<Button variant="outline" onclick={() => setArchived(false)}><ArchiveRestore />{$t`Reopen`}</Button>{/snippet}
{#snippet archiveButton()}<Button variant="outline" onclick={() => setArchived(true)}><Archive />{$t`Close event`}</Button>{/snippet}
{#snippet copyButton()}<Button variant="outline" onclick={openCopy}><CopyPlus />{$t`Copy`}</Button>{/snippet}
{#snippet deleteButton()}<Button variant="destructive" onclick={remove}><Trash />{$t`Delete`}</Button>{/snippet}

{#if event.data}
	{@const ev = event.data}
	{#if isAdmin}
		<NativeSelect.Root class="mb-4 w-full sm:hidden" value={tab} onchange={(e) => goto(`?tab=${e.currentTarget.value}`, { replaceState: true, noScroll: true })} aria-label={$t`Setup`}>
			{#each tabs as [key, label] (key)}<NativeSelect.Option value={key}>{label}</NativeSelect.Option>{/each}
		</NativeSelect.Root>
		<div class="mb-4 hidden gap-1 overflow-x-auto border-b sm:flex" role="tablist" aria-label={$t`Setup`}>
			{#each tabs as [key, label] (key)}
				<a
					href="?tab={key}"
					role="tab"
					aria-selected={tab === key}
					data-sveltekit-replacestate
					data-sveltekit-noscroll
					class="-mb-px border-b-2 px-3 py-2 text-sm font-medium whitespace-nowrap {tab === key
						? 'border-primary text-foreground'
						: 'border-transparent text-muted-foreground hover:text-foreground'}">{label}</a
				>
			{/each}
		</div>
	{/if}
	<div class="space-y-4 {tab === 'programmes' ? '' : 'max-w-3xl'}">
		{#if tab === 'programmes'}
			<section class="overflow-hidden rounded-lg border bg-card"><ProgrammeList {tenantId} {eventId} isAdmin={isAdmin && !ev.archivedAt} /></section>
		{:else if tab === 'event' && ranges.data}
			{#if ev.archivedAt}
				{@render card($t`Reopen this event`, $t`Registrations, results and settings can be changed again.`, unarchive)}
			{:else}
				<section class="rounded-lg border bg-card p-4">
					{#key ev.id}<EventForm ranges={ranges.data} clubs={clubs.data ?? []} initial={{ ...ev, rangeIds: ev.ranges.map((r) => r.id) }} submitLabel={$t`Save`} onsubmit={save} />{/key}
				</section>
			{/if}
		{:else if tab === 'classifications'}
			<Classifications {tenantId} {eventId} />
		{:else if tab === 'kiosk'}
			<section class="overflow-hidden rounded-lg border bg-card"><KioskTokens {tenantId} {eventId} /></section>
		{:else if tab === 'helpers'}
			<section class="overflow-hidden rounded-lg border bg-card"><HelperCodes {tenantId} {eventId} /></section>
		{:else if tab === 'booking'}
			<BookingSettings {tenantId} event={ev} />
		{:else if tab === 'more'}
			{@render card($t`Copy for next year`, $t`A new event with the same settings and programmes, without registrations.`, copyButton)}
			{#if ev.archivedAt}
				{@render card($t`Reopen this event`, $t`Registrations, results and settings can be changed again.`, unarchive)}
			{:else}
				{@render card($t`Close this event`, $t`When the event is over: the rankings are final and reach the shooters' personal training logs, everything stays visible, but results and registrations can no longer be changed.`, archiveButton)}
				{@render card($t`Delete this event`, $t`Only possible as long as nobody is registered.`, deleteButton, true)}
			{/if}
		{/if}
	</div>
{/if}

<Dialog.Root open={copy !== null} onOpenChange={(o) => !o && (copy = null)}>
	<Dialog.Content class="sm:max-w-md">
		{#if copy}
			<Dialog.Header>
				<Dialog.Title>{$t`Copy event`}</Dialog.Title>
				<Dialog.Description>{$t`Range, organizer, language, settings, all programmes with prices, awards, combinations and special prizes, and the classifications are copied. Registrations, results and kiosk links are not.`}</Dialog.Description>
			</Dialog.Header>
			<form class="grid gap-4" onsubmit={copyEvent}>
				<FormField label={$t`Name`}><Input name="name" bind:value={copy.name} required maxlength={200} /></FormField>
				<FormField label={$t`Date`}><Input name="startsOn" type="date" bind:value={copy.startsOn} required /></FormField>
				<FormError message={copyError} />
				<Dialog.Footer><Button type="submit"><CopyPlus />{$t`Create copy`}</Button></Dialog.Footer>
			</form>
		{/if}
	</Dialog.Content>
</Dialog.Root>
