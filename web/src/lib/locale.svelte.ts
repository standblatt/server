import { setupI18n } from '@lingui/core';
import { messages as de } from '$locales/de';
import { messages as en } from '$locales/en';
import { messages as fr } from '$locales/fr';
import { messages as it } from '$locales/it';

export const localeNames = { en: 'English', de: 'Deutsch', fr: 'Français', it: 'Italiano' } as const;
export type Locale = keyof typeof localeNames;

/** The Lingui instance of the app, set as context in the root layout. */
export const i18n = setupI18n({ locale: 'en', messages: { en, de, fr, it } });

const storeKey = 'standblatt-locale';
const state = $state<{ current: Locale }>({ current: 'en' });

/** Supported locale for a language tag (first two letters), else undefined. */
function supported(l: string | null | undefined): Locale | undefined {
	const s = l?.slice(0, 2);
	return s && s in localeNames ? (s as Locale) : undefined;
}

/** The locale shown now. */
export const currentLocale = (): Locale => state.current;

/** BCP 47 tag for Intl formatting, e.g. de-CH. */
export const intlLocale = () => `${state.current}-CH`;

/** The remembered locale, else the first supported browser language. */
export function preferredLocale(): Locale {
	let stored = null;
	try {
		stored = localStorage.getItem(storeKey);
	} catch {
		// not remembered
	}
	return supported(stored) ?? navigator.languages.map(supported).find((l) => l) ?? 'en';
}

/** Switches the UI to a locale; remember keeps it for the next visit. */
export function setLocale(l: string, remember = true) {
	const loc = supported(l) ?? 'en';
	if (remember) {
		try {
			localStorage.setItem(storeKey, loc);
		} catch {
			// not remembered
		}
	}
	i18n.activate(loc);
	state.current = loc;
	document.documentElement.lang = loc;
}
