import { derived, type Readable } from 'svelte/store';
import { locale } from './store';
import { getTranslation, type TranslationKey } from './translations';

export const t: Readable<(key: TranslationKey) => string> = derived(
	locale,
	($locale) => (key: TranslationKey) => getTranslation($locale, key)
);

