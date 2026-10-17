<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import ListTree from '@lucide/svelte/icons/list-tree';
	import { Button } from '$lib/components/ui/button';
	import HelpArticle from '$lib/help/HelpArticle.svelte';
	import HelpNav from '$lib/help/HelpNav.svelte';
	import { currentLocale, localeNames, setLocale, type Locale } from '$lib/locale.svelte';
	import PublicShell from '$lib/PublicShell.svelte';

	const wanted = page.url.searchParams.get('lang');
	if (wanted && wanted in localeNames) setLocale(wanted as Locale, false);

	const slug = $derived((page.params.path ?? '').replace(/\/$/, ''));
	const anchor = $derived(decodeURIComponent(page.url.hash.slice(1)));
	const lang = $derived(currentLocale());
	let menu = $state(false);

	/** Closes the contents on phones once a page is chosen. */
	const closeOnLink = (e: MouseEvent) => (e.target as HTMLElement).closest('a') && (menu = false);
</script>

<PublicShell wide>
	<div class="grid gap-6 md:grid-cols-[15rem_minmax(0,1fr)] md:gap-8">
		<aside class="md:sticky md:top-6 md:max-h-[calc(100svh-3rem)] md:self-start md:overflow-y-auto">
			<Button variant="outline" size="sm" class="md:hidden" aria-expanded={menu} onclick={() => (menu = !menu)}><ListTree />{$t`Contents`}</Button>
			<!-- svelte-ignore a11y_click_events_have_key_events, a11y_no_static_element_interactions -->
			<div class="{menu ? 'mt-3' : 'hidden'} md:mt-0 md:block" onclick={closeOnLink}><HelpNav {lang} current={slug} /></div>
		</aside>
		<div class="min-w-0">
			{#key `${lang}/${slug}`}<HelpArticle {slug} {lang} {anchor} />{/key}
		</div>
	</div>
</PublicShell>
