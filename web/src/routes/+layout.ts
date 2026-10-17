import { preferredLocale, setLocale } from '$lib/locale.svelte';

export const ssr = false;
export const prerender = false;

/** Applies the preferred locale before the first render. */
export const load = () => setLocale(preferredLocale(), false);
