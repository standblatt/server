<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import ClipboardPaste from '@lucide/svelte/icons/clipboard-paste';
	import Plus from '@lucide/svelte/icons/plus';
	import X from '@lucide/svelte/icons/x';
	import type { BookingShooter } from '../api';
	import { parseShooters } from '../booking';
	import { Button } from '../components/ui/button';
	import * as Dialog from '../components/ui/dialog';
	import { Input } from '../components/ui/input';
	import { Textarea } from '../components/ui/textarea';

	let { shooters = $bindable(), readonly = false }: { shooters: BookingShooter[]; readonly?: boolean } = $props();

	let pasting = $state(false);
	let pasted = $state('');
	const preview = $derived(parseShooters(pasted));

	/** Adds an empty row. */
	function add() {
		shooters.push({ firstname: '', lastname: '', birthday: '', licenseNumber: '' });
	}

	/** Appends the shooters read from the pasted text and closes the dialog. */
	function takePasted() {
		shooters.push(...preview);
		pasted = '';
		pasting = false;
	}
</script>

<div class="grid gap-2">
	{#if shooters.length > 0}
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead class="text-left text-xs text-muted-foreground">
					<tr>
						<th class="w-8 py-1 font-medium">#</th>
						<th class="py-1 font-medium">{$t`Last name`}</th>
						<th class="py-1 font-medium">{$t`First name`}</th>
						<th class="py-1 font-medium">{$t`Birthday`}</th>
						<th class="py-1 font-medium">{$t`Licence (optional)`}</th>
						<th class="w-10"></th>
					</tr>
				</thead>
				<tbody>
					{#each shooters as s, i (i)}
						<tr>
							<td class="py-1 text-muted-foreground tabular-nums">{i + 1}</td>
							<td class="py-1 pr-2"><Input bind:value={s.lastname} required disabled={readonly} aria-label={$t`Last name of shooter ${i + 1}`} /></td>
							<td class="py-1 pr-2"><Input bind:value={s.firstname} required disabled={readonly} aria-label={$t`First name of shooter ${i + 1}`} /></td>
							<td class="py-1 pr-2"><Input type="date" bind:value={s.birthday} required disabled={readonly} aria-label={$t`Birthday of shooter ${i + 1}`} /></td>
							<td class="py-1 pr-2"><Input bind:value={s.licenseNumber} maxlength={20} disabled={readonly} aria-label={$t`Licence of shooter ${i + 1}`} /></td>
							<td class="py-1">
								{#if !readonly}<Button variant="ghost" size="icon-sm" aria-label={$t`Remove shooter ${i + 1}`} onclick={() => shooters.splice(i, 1)}><X /></Button>{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else}
		<p class="text-sm text-muted-foreground">{$t`No shooters yet. You can add them now or later until the deadline.`}</p>
	{/if}
	{#if !readonly}
		<div class="flex flex-wrap gap-2">
			<Button variant="outline" onclick={add}><Plus />{$t`Add shooter`}</Button>
			<Button variant="outline" onclick={() => (pasting = true)}><ClipboardPaste />{$t`Paste list`}</Button>
		</div>
	{/if}
</div>

<Dialog.Root bind:open={pasting}>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>{$t`Paste list`}</Dialog.Title>
			<Dialog.Description>{$t`Copy the rows from Excel or a text file: last name, first name, birthday (e.g. 01.05.1990), licence. One shooter per line.`}</Dialog.Description>
		</Dialog.Header>
		<Textarea bind:value={pasted} rows={8} aria-label={$t`List of shooters`} />
		<p class="text-sm text-muted-foreground">{$plural(preview.length, { one: '# shooter recognised', other: '# shooters recognised' })}</p>
		<Dialog.Footer><Button disabled={preview.length === 0} onclick={takePasted}>{$t`Add to the list`}</Button></Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
