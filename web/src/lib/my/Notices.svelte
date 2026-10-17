<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { toast } from 'svelte-sonner';
	import CloudOff from '@lucide/svelte/icons/cloud-off';
	import Mail from '@lucide/svelte/icons/mail';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import TriangleAlert from '@lucide/svelte/icons/triangle-alert';
	import { ask } from '$lib/confirm.svelte';
	import { discardFailed, outbox, type MyOverview } from '$lib/my.svelte';

	let { my }: { my: MyOverview } = $props();

	let sent = $state(false);

	/** Sends the email verification link again. */
	async function resend() {
		await call(api.POST('/api/my/email/verify'));
		sent = true;
		toast.success($t`We sent you a new link.`);
	}

	/** Drops the sessions whose upload the server refused, after asking. */
	async function discard() {
		if (await ask({ title: $t`Discard these sessions?`, description: $t`They are removed from this device and cannot be uploaded anymore.`, action: $t`Discard`, destructive: true })) discardFailed(my.userId);
	}
</script>

{#if my.email && !my.emailVerified}
	<div class="flex flex-wrap items-center gap-3 rounded-lg border bg-card px-4 py-3 text-sm" role="status">
		<Mail class="size-5 shrink-0 text-muted-foreground" />
		<p class="min-w-0 flex-1 [overflow-wrap:anywhere]">{$t`Please confirm your email address ${my.email} with the link we sent you.`}</p>
		<Button size="sm" variant="ghost" disabled={sent} onclick={resend}>{$t`Send again`}</Button>
	</div>
{/if}
{#if outbox.count}
	<div class="flex items-center gap-3 rounded-lg border bg-card px-4 py-3 text-sm" role="status">
		<CloudOff class="size-5 shrink-0 text-muted-foreground" />
		<p>{$plural(outbox.count, { one: '# session is saved on this device and will be uploaded when you are online.', other: '# sessions are saved on this device and will be uploaded when you are online.' })}</p>
	</div>
{/if}
{#if outbox.failed.length}
	<div class="flex flex-wrap items-center gap-3 rounded-lg border border-destructive/50 bg-card px-4 py-3 text-sm" role="alert" data-testid="outbox-failed">
		<TriangleAlert class="size-5 shrink-0 text-destructive" />
		<div class="min-w-0 flex-1 [overflow-wrap:anywhere]">
			<p>{$plural(outbox.failed.length, { one: '# session saved on this device could not be uploaded:', other: '# sessions saved on this device could not be uploaded:' })}</p>
			<ul class="list-disc pl-5 text-muted-foreground">{#each [...new Set(outbox.failed)] as e (e)}<li>{e}</li>{/each}</ul>
		</div>
		<Button size="sm" variant="ghost" onclick={discard}>{$t`Discard`}</Button>
	</div>
{/if}
