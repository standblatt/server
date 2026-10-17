<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import FileText from '@lucide/svelte/icons/file-text';
	import Printer from '@lucide/svelte/icons/printer';
	import { toast } from 'svelte-sonner';
	import { Button } from './components/ui/button';
	import { pdfClick, pdfUrl, printDocument, type PrintTarget } from './print';

	let { tenantId, target }: { tenantId: number; target: PrintTarget } = $props();

	const text = $derived(
		{
			lane_sheet: [$t`Print lane sheet`, $t`Lane sheet sent to printer`],
			control_strip: [$t`Print control strip`, $t`Control strip sent to printer`],
			result_sheet: [$t`Print result sheet`, $t`Result sheet sent to printer`],
			wreath_card: [$t`Print wreath card`, $t`Wreath card sent to printer`],
			ranking: [$t`Print ranking`, $t`Ranking sent to printer`]
		}[target.document]
	);

	/** Sends the document to the printer (toast when queued) or opens the PDF. */
	async function print() {
		try {
			if ((await printDocument(tenantId, target)) === 'queued') toast.success(text[1]);
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
</script>

<span class="inline-flex items-center">
	<Button size="sm" variant="outline" class="rounded-r-none" onclick={print}><Printer />{text[0]}</Button>
	<Button size="sm" variant="outline" class="-ml-px rounded-l-none" href={pdfUrl(tenantId, target)} target="_blank" onclick={pdfClick} title={$t`Open as PDF`}><FileText />PDF</Button>
</span>
