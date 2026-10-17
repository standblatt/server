<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import type { BookingRequest } from '../api';
	import { Badge } from '../components/ui/badge';

	let { status }: { status: BookingRequest['status'] } = $props();

	const label = $derived(
		status === 'confirmed' ? $t`confirmed` : status === 'declined' ? $t`declined` : status === 'cancelled' ? $t`cancelled` : $t`waiting for confirmation`
	);
	const tone: Record<string, string> = {
		pending: 'border-amber-500/50 bg-amber-500/10 text-amber-700 dark:text-amber-400',
		confirmed: 'border-emerald-500/50 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400',
		declined: 'border-destructive/50 text-destructive',
		cancelled: 'text-muted-foreground'
	};
</script>

<Badge variant="outline" class={tone[status]}>{label}</Badge>
