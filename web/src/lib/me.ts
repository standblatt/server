import { createQuery, type QueryClient } from '@tanstack/svelte-query';
import { api, call } from './api';
import { setLocale, type Locale } from './locale.svelte';

export const meKey = ['me'];

/** Query of the logged-in user; applies the user's saved locale. */
export function meQuery() {
	return createQuery(() => ({
		queryKey: meKey,
		queryFn: async () => {
			const me = await call(api.GET('/api/me'));
			if (me.user.locale) setLocale(me.user.locale);
			return me;
		},
		retry: false
	}));
}

/** Stores the UI language on the account, applies it and refreshes the user. */
export async function setMyLocale(queryClient: QueryClient, locale: string) {
	await call(api.PUT('/api/me/locale', { body: { locale: locale as Locale } }));
	setLocale(locale);
	await queryClient.invalidateQueries({ queryKey: meKey });
}

/** What this server offers on the login page: sign-up and the second factor. */
export function authOptionsQuery() {
	return createQuery(() => ({
		queryKey: ['auth-options'],
		queryFn: () => call(api.GET('/api/auth/options')),
		staleTime: Infinity
	}));
}

/** The server version, shown at the bottom of the sidebars. */
export function versionQuery() {
	return createQuery(() => ({
		queryKey: ['version'],
		queryFn: () => call(api.GET('/api/version')),
		staleTime: Infinity
	}));
}
