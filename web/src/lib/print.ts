import { msg } from '@lingui/core/macro';
import { toast } from 'svelte-sonner';
import { api, call, fetchFile, helperCode } from './api';
import { i18n } from './locale.svelte';

/** Certificates: one page per ranked shooter of a programme or classification (class) in ranks from to to. */
export type CertificatesTarget = {
	document: 'certificate';
	eventId: number;
	layoutId: number;
	programmeId?: number;
	classificationId?: number;
	class?: string;
	from: number;
	to: number;
};

export type PrintTarget =
	| { document: 'lane_sheet' | 'control_strip' | 'result_sheet'; registrationId: number }
	| { document: 'wreath_card' | 'form'; entryId: number }
	| { document: 'form'; registrationId: number; layoutId: number }
	| CertificatesTarget;

/** API URL of the PDF for a print target. */
export function pdfUrl(tenantId: number, t: PrintTarget): string {
	if (t.document === 'certificate') {
		const q = new URLSearchParams({ layout: String(t.layoutId), from: String(t.from), to: String(t.to) });
		if (t.programmeId) q.set('programme', String(t.programmeId));
		if (t.classificationId) q.set('classification', String(t.classificationId));
		if (t.class) q.set('class', t.class);
		return `/api/tenants/${tenantId}/events/${t.eventId}/certificates?${q}`;
	}
	if ('layoutId' in t) return `/api/tenants/${tenantId}/registrations/${t.registrationId}/forms/${t.layoutId}`;
	if ('entryId' in t) return `/api/tenants/${tenantId}/entries/${t.entryId}/${t.document === 'wreath_card' ? 'wreath-card' : 'form'}`;
	return `/api/tenants/${tenantId}/registrations/${t.registrationId}/documents/${t.document}`;
}

/** printDocument sends a document to the range's printer; without a printer or agent it opens the PDF instead. */
export async function printDocument(tenantId: number, t: PrintTarget): Promise<'queued' | 'opened'> {
	const r = await call(api.POST('/api/tenants/{tenantId}/print', { params: { path: { tenantId } }, body: t }));
	if (r.queued) return 'queued';
	await openPdf(pdfUrl(tenantId, t));
	return 'opened';
}

/** True on touch screens, where a PDF goes to the share sheet instead of a new tab. */
const touch = () => matchMedia('(pointer: coarse)').matches;

const pdfReady = msg`The PDF is ready.`;
const shareLabel = msg`Share`;

/** Opens a PDF: a new tab on a desktop, the share sheet (or a download) on a phone; helper links fetch it with their code. */
async function openPdf(url: string) {
	if (!touch() && !helperCode()) {
		window.open(url, '_blank');
		return;
	}
	const file = await fetchFile(url);
	if (!touch()) {
		window.open(URL.createObjectURL(file), '_blank');
		return;
	}
	if (!navigator.canShare?.({ files: [file] })) {
		const a = document.createElement('a');
		a.href = URL.createObjectURL(file);
		a.download = file.name;
		a.click();
		setTimeout(() => URL.revokeObjectURL(a.href), 60_000);
		return;
	}
	try {
		await navigator.share({ files: [file] });
	} catch (err) {
		if ((err as Error).name !== 'NotAllowedError') return;
		toast(i18n._(pdfReady), { action: { label: i18n._(shareLabel), onClick: () => navigator.share({ files: [file] }).catch(() => {}) } });
	}
}

/** Click handler for PDF links: a desktop follows the link, a phone or helper link goes through openPdf. */
export function pdfClick(e: MouseEvent) {
	if (!touch() && !helperCode()) return;
	e.preventDefault();
	openPdf((e.currentTarget as HTMLAnchorElement).href);
}
