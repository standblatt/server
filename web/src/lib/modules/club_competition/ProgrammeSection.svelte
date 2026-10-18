<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { Switch } from '$lib/components/ui/switch';
	import type { ProgrammeProps } from '..';

	let { draft }: ProgrammeProps = $props();
</script>

<label class="flex items-center justify-between gap-3 rounded-lg border bg-card p-4 text-sm">{$t`Counts for the club competition`} <Switch bind:checked={draft.countsForClubCompetition} /></label>
