import { createQuery } from '@tanstack/svelte-query';
import { api, call, fetchFile, type LayoutInput } from '../api';

/** An image as stored in an exported form file: its old ID, name and the PNG or JPEG as base64. */
export type ExportedImage = { id: number; name: string; data: string };

/** Query key of a club's images. */
export const imagesKey = (tenantId: number) => ['tenants', tenantId, 'images'];

/** Query of a club's image library. */
export function imagesQuery(tenantId: () => number) {
	return createQuery(() => ({
		queryKey: imagesKey(tenantId()),
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/images', { params: { path: { tenantId: tenantId() } } }))
	}));
}

/** URL of an image of the library. */
export const imageUrl = (tenantId: number, id: number) => `/api/tenants/${tenantId}/images/${id}`;

/** Adds a PNG or JPEG to the library. */
export async function uploadImage(tenantId: number, name: string, file: Blob) {
	const form = new FormData();
	form.append('name', name);
	form.append('file', file);
	return call(api.POST('/api/tenants/{tenantId}/images', { params: { path: { tenantId } }, body: {} as never, bodySerializer: () => form }));
}

/** The images a form shows, as base64 for an export file. */
export async function exportImages(tenantId: number, form: LayoutInput, names: Map<number, string>): Promise<ExportedImage[]> {
	const ids = [...new Set(form.elements.filter((e) => e.type === 'image' && e.image).map((e) => e.image!))];
	return Promise.all(
		ids.map(async (id) => {
			const bytes = new Uint8Array(await (await fetchFile(imageUrl(tenantId, id))).arrayBuffer());
			let bin = '';
			for (const b of bytes) bin += String.fromCharCode(b);
			return { id, name: names.get(id) ?? `image-${id}`, data: btoa(bin) };
		})
	);
}

/** Uploads the images of an imported form file and points its image elements at them; an image whose name the club
 * already has is reused. */
export async function importImages(tenantId: number, form: LayoutInput, images: ExportedImage[], existing: Map<string, number>) {
	const ids = new Map<number, number>();
	for (const img of images) {
		let id = existing.get(img.name);
		if (id === undefined) {
			const bytes = Uint8Array.from(atob(img.data), (c) => c.charCodeAt(0));
			id = (await uploadImage(tenantId, img.name, new Blob([bytes]))).id;
		}
		ids.set(img.id, id);
	}
	form.elements = form.elements.map((e) => (e.type === 'image' && e.image ? { ...e, image: ids.get(e.image) ?? e.image } : e));
}
