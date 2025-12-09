import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';

export type Locale = 'en' | 'pt';

const LOCALE_STORAGE_KEY = 'woragis-locale';

// Get initial locale from localStorage or default to 'en'
function getInitialLocale(): Locale {
	if (typeof window !== 'undefined') {
		const stored = localStorage.getItem(LOCALE_STORAGE_KEY);
		if (stored === 'en' || stored === 'pt') {
			return stored;
		}
	}
	return 'en';
}

export const locale: Writable<Locale> = writable<Locale>(getInitialLocale());

// Subscribe to locale changes and persist to localStorage
if (typeof window !== 'undefined') {
	locale.subscribe((value) => {
		localStorage.setItem(LOCALE_STORAGE_KEY, value);
	});
}

