<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { useQueryClient } from '@tanstack/svelte-query';
	import ImageUp from '@lucide/svelte/icons/image-up';
	import Trash from '@lucide/svelte/icons/trash-2';
	import { toast } from 'svelte-sonner';
	import { api, call, type LibraryImage } from '../api';
	import { Button } from '../components/ui/button';
	import { ask } from '../confirm.svelte';
	import { imagesKey, imagesQuery, imageUrl, uploadImage } from './images';

	let { tenantId }: { tenantId: number } = $props();
	const queryClient = useQueryClient();
	const images = imagesQuery(() => tenantId);
	let input = $state<HTMLInputElement>();

	/** Adds the chosen file to the library, named after the file. */
	async function upload(e: Event) {
		const f = (e.currentTarget as HTMLInputElement).files?.[0];
		if (!f) return;
		try {
			const img = await uploadImage(tenantId, f.name.replace(/\.[^.]+$/, ''), f);
			await queryClient.invalidateQueries({ queryKey: imagesKey(tenantId) });
			toast.success($t`Image "${img.name}" added`);
		} catch (err) {
			toast.error((err as Error).message);
		} finally {
			if (input) input.value = '';
		}
	}

	/** Deletes an image after confirmation. */
	async function remove(img: LibraryImage) {
		if (!(await ask({ title: $t`Delete image "${img.name}"?`, description: $t`Only possible while no print form shows it.`, action: $t`Delete`, destructive: true }))) return;
		await call(api.DELETE('/api/tenants/{tenantId}/images/{imageId}', { params: { path: { tenantId, imageId: img.id } } }));
		await queryClient.invalidateQueries({ queryKey: imagesKey(tenantId) });
		toast.success($t`Image deleted`);
	}
</script>

<section class="overflow-hidden rounded-lg border bg-card" aria-label={$t`Images`}>
	<div class="flex flex-wrap items-center gap-2 border-b px-4 py-2.5">
		<div class="mr-auto">
			<h2 class="text-sm font-semibold">{$t`Images`}</h2>
			<p class="text-xs text-muted-foreground">{$t`Logos and other images for print forms (PNG or JPEG, at most 2 MB).`}</p>
		</div>
		<Button size="sm" variant="outline" onclick={() => input?.click()}><ImageUp />{$t`Upload image`}</Button>
		<input bind:this={input} class="hidden" type="file" accept="image/png,image/jpeg" onchange={upload} aria-label={$t`Upload image`} />
	</div>
	<ul class="grid gap-3 p-4 sm:grid-cols-2 lg:grid-cols-4">
		{#each images.data ?? [] as img (img.id)}
			<li class="flex items-center gap-3 rounded-md border p-2">
				<img src={imageUrl(tenantId, img.id)} alt="" class="size-12 shrink-0 rounded border bg-white object-contain" />
				<div class="min-w-0 flex-1 text-sm">
					<div class="truncate font-medium">{img.name}</div>
					<div class="text-xs text-muted-foreground">{img.width} × {img.height} px</div>
				</div>
				<Button variant="ghost" size="icon-sm" title={$t`Delete`} aria-label={$t`Delete ${img.name}`} onclick={() => remove(img)}><Trash /></Button>
			</li>
		{:else}
			<li class="text-sm text-muted-foreground">{$t`No images yet.`}</li>
		{/each}
	</ul>
</section>
