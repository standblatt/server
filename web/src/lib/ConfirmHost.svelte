<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';
	import { confirmState } from './confirm.svelte';

	const p = $derived(confirmState.pending);
</script>

<AlertDialog.Root open={!!p} onOpenChange={(o) => !o && confirmState.answer(false)}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>{p?.title}</AlertDialog.Title>
			{#if p?.description}<AlertDialog.Description>{p.description}</AlertDialog.Description>{/if}
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={() => confirmState.answer(false)}>{$t`Cancel`}</AlertDialog.Cancel>
			<AlertDialog.Action class={p?.destructive ? buttonVariants({ variant: 'destructive' }) : ''} onclick={() => confirmState.answer(true)}>
				{p?.action ?? 'OK'}
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
