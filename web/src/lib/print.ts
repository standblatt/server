import { api, call } from './api';

export type PrintTarget =
	| { document: 'lane_sheet' | 'control_strip' | 'result_sheet'; registrationId: number }
	| { document: 'wreath_card'; entryId: number }
	| { document: 'ranking'; programmeId: number };

/** API URL of the PDF for a print target. */
export function pdfUrl(tenantId: number, t: PrintTarget): string {
	switch (t.document) {
		case 'wreath_card':
			return `/api/tenants/${tenantId}/entries/${t.entryId}/wreath-card`;
		case 'ranking':
			return `/api/tenants/${tenantId}/programmes/${t.programmeId}/ranking`;
		default:
			return `/api/tenants/${tenantId}/registrations/${t.registrationId}/documents/${t.document}`;
	}
}

/** printDocument sends a document to the range's printer; without a printer or agent it opens the PDF instead. */
export async function printDocument(tenantId: number, t: PrintTarget): Promise<'queued' | 'opened'> {
	const r = await call(api.POST('/api/tenants/{tenantId}/print', { params: { path: { tenantId } }, body: t }));
	if (r.queued) return 'queued';
	window.open(pdfUrl(tenantId, t), '_blank');
	return 'opened';
}
