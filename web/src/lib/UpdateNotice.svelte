<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { untrack } from 'svelte';
	import { beforeNavigate } from '$app/navigation';
	import { page, updated } from '$app/state';
	import { toast } from 'svelte-sonner';

	$effect(() => {
		if (!updated.current) return;
		if (page.url.pathname.startsWith('/kiosk/')) {
			location.reload();
			return;
		}
		untrack(() =>
			toast.warning($t`A new version of Standblatt is available.`, {
				position: 'top-right',
				duration: Infinity,
				action: { label: $t`Reload`, onClick: () => location.reload() }
			})
		);
	});

	beforeNavigate(({ willUnload, to }) => {
		if (updated.current && !willUnload && to) location.href = to.url.href;
	});
</script>
