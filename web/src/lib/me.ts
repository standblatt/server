import { createQuery, type QueryClient } from '@tanstack/svelte-query';
import { api, call, type Me, type Permission } from './api';
import { setLocale, type Locale } from './locale.svelte';

export const meKey = ['me'];

/** The user's membership in a club, undefined when the user is no member. */
export const membershipOf = (me: Me | undefined, tenantId: number) => me?.tenants.find((m) => m.tenantId === tenantId);

/** Whether the role of a membership has a permission. */
export const can = (membership: { permissions: Permission[] } | undefined, permission: Permission) => !!membership?.permissions.includes(permission);

/** Query of the logged-in user; applies the user's saved locale. Pass enabled false where no user is logged in,
 * e.g. on helper links, since a 401 sends the browser to the login. */
export function meQuery(enabled: () => boolean = () => true) {
	return createQuery(() => ({
		enabled: enabled(),
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
