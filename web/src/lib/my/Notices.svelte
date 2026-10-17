<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { toast } from 'svelte-sonner';
	import CloudOff from '@lucide/svelte/icons/cloud-off';
	import Mail from '@lucide/svelte/icons/mail';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { outbox, type MyOverview } from '$lib/my.svelte';

	let { my }: { my: MyOverview } = $props();

	let sent = $state(false);

	/** Sends the email verification link again. */
	async function resend() {
		await call(api.POST('/api/my/email/verify'));
		sent = true;
		toast.success($t`We sent you a new link.`);
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
