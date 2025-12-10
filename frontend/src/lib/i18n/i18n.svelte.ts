import { derived, type Readable } from 'svelte/store';
import { locale } from './store';
import { getTranslation, type TranslationKey } from './translations';

export const t: Readable<(key: TranslationKey) => string> = derived(
	locale,
	($locale) => (key: TranslationKey) => getTranslation($locale, key)
);

// Helper function to use translation in Svelte 5 runes mode
export function useTranslation() {
	// Track the current locale value
	let currentLocale = $state<'en' | 'pt'>('en');
	
	// Subscribe to locale changes
	$effect(() => {
		const unsubscribe = locale.subscribe((l) => {
			currentLocale = l;
		});
		return unsubscribe;
	});
	
	// Return a derived translation function that reacts to locale changes
	const translationFn = $derived.by(() => (key: TranslationKey) => getTranslation(currentLocale, key));
	return translationFn;
}
