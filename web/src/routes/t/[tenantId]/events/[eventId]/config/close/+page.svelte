<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { Snippet } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import Archive from '@lucide/svelte/icons/archive';
	import ArchiveRestore from '@lucide/svelte/icons/archive-restore';
	import CopyPlus from '@lucide/svelte/icons/copy-plus';
	import LayoutTemplate from '@lucide/svelte/icons/layout-template';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import FormError from '$lib/FormError.svelte';
	import FormField from '$lib/FormField.svelte';
	import { ask } from '$lib/confirm.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { eventQuery, keys } from '$lib/tenant';

	const tenantId = $derived(Number(page.params.tenantId));
	const eventId = $derived(Number(page.params.eventId));
	const queryClient = useQueryClient();
	const event = eventQuery(
		() => tenantId,
		() => eventId
	);

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

	let copy = $state<{ name: string; startsOn: string; template?: boolean } | null>(null);
	let copyError = $state('');

	/** Event name for next year: years and a leading edition number count up. */
	const nextName = (name: string) =>
		name.replace(/\b(19|20)\d{2}\b/g, (y) => String(Number(y) + 1)).replace(/^(\d+)\./, (_, n) => `${Number(n) + 1}.`);
	/** Same date one year later (YYYY-MM-DD). */
	const nextYear = (date: string) => `${Number(date.slice(0, 4)) + 1}${date.slice(4)}`;

	/** Opens the dialog that saves the event as a template of the club, named without its year. */
	function openTemplate() {
		copyError = '';
		const name = (event.data?.name ?? '').replace(/\s*\b(19|20)\d{2}\b/g, '').trim();
		copy = { name: name || (event.data?.name ?? ''), startsOn: event.data?.startsOn ?? '', template: true };
	}

	/** Opens the copy dialog with next year's name and date. */
	function openCopy() {
		copyError = '';
		copy = { name: nextName(event.data?.name ?? ''), startsOn: nextYear(event.data?.startsOn ?? '') };
	}

	/** Copies the event for next year and opens the copy's configuration. */
	async function copyEvent(e: SubmitEvent) {
		e.preventDefault();
		if (!copy) return;
		try {
			const ev = await call(api.POST('/api/tenants/{tenantId}/events/{eventId}/copy', { params: { path: { tenantId, eventId } }, body: copy }));
			await queryClient.invalidateQueries({ queryKey: keys.events(tenantId) });
			toast.success(ev.template ? $t`Template "${ev.name}" saved` : $t`Event "${ev.name}" created`);
			copy = null;
			await goto(`/t/${tenantId}/events/${ev.id}/config`);
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

<PageHeader help="user-guide/events#close" title={$t`Copy & close`} description={$t`Copy the event for next year, close it when it is over or delete it.`} />

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
{#snippet templateButton()}<Button variant="outline" onclick={openTemplate}><LayoutTemplate />{$t`Save as template`}</Button>{/snippet}
{#snippet copyButton()}<Button variant="outline" onclick={openCopy}><CopyPlus />{$t`Copy`}</Button>{/snippet}
{#snippet deleteButton()}<Button variant="destructive" onclick={remove}><Trash />{$t`Delete`}</Button>{/snippet}

{#if event.data}
	<div class="max-w-3xl space-y-4">
		{#if !event.data.template}
			{@render card($t`Copy for next year`, $t`A new event with the same settings, modules and programmes, without registrations.`, copyButton)}
			{@render card($t`Save as template`, $t`New events of the club can start from it, with its modules, settings, ranges, programmes and classifications.`, templateButton)}
		{/if}
		{#if event.data.archivedAt}
			{@render card($t`Reopen this event`, $t`Registrations, results and settings can be changed again.`, unarchive)}
		{:else}
			{@render card($t`Close this event`, $t`When the event is over: the rankings are final and reach the shooters' personal training logs, everything stays visible, but results and registrations can no longer be changed.`, archiveButton)}
			{@render card($t`Delete this event`, $t`Only possible as long as nobody is registered.`, deleteButton, true)}
		{/if}
	</div>
{/if}

<Dialog.Root open={copy !== null} onOpenChange={(o) => !o && (copy = null)}>
	<Dialog.Content class="sm:max-w-md">
		{#if copy}
			<Dialog.Header>
				<Dialog.Title>{copy.template ? $t`Save as template` : $t`Copy event`}</Dialog.Title>
				<Dialog.Description>{$t`Range, organizer, language, settings, modules, all programmes with prices, awards, combinations and special prizes, and the classifications are copied. Registrations, results and kiosk links are not.`}</Dialog.Description>
			</Dialog.Header>
			<form class="grid gap-4" onsubmit={copyEvent}>
				<FormField label={$t`Name`}><Input name="name" bind:value={copy.name} required maxlength={200} /></FormField>
				{#if !copy.template}<FormField label={$t`Date`}><Input name="startsOn" type="date" bind:value={copy.startsOn} required /></FormField>{/if}
				<FormError message={copyError} />
				<Dialog.Footer><Button type="submit"><CopyPlus />{copy.template ? $t`Save template` : $t`Create copy`}</Button></Dialog.Footer>
			</form>
		{/if}
	</Dialog.Content>
</Dialog.Root>
