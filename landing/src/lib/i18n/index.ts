import { writable, get, derived } from 'svelte/store';
import { browser } from '$app/environment';
import type { Translations } from './translations';

export type Language = 'en' | 'pt-BR' | 'fr' | 'es' | 'de' | 'ru' | 'ja' | 'ko' | 'zh-CN' | 'el' | 'la';

// Load translations
import enTranslations from './locales/en.json';
import ptBRTranslations from './locales/pt-BR.json';
import frTranslations from './locales/fr.json';
import esTranslations from './locales/es.json';
import deTranslations from './locales/de.json';
import ruTranslations from './locales/ru.json';
import jaTranslations from './locales/ja.json';
import koTranslations from './locales/ko.json';
import zhCNTranslations from './locales/zh-CN.json';
import elTranslations from './locales/el.json';
import laTranslations from './locales/la.json';

const translations: Record<Language, Translations> = {
	'en': enTranslations,
	'pt-BR': ptBRTranslations,
	'fr': frTranslations,
	'es': esTranslations,
	'de': deTranslations,
	'ru': ruTranslations,
	'ja': jaTranslations,
	'ko': koTranslations,
	'zh-CN': zhCNTranslations,
	'el': elTranslations,
	'la': laTranslations
};

// Get initial language from localStorage or browser, default to English
function getInitialLanguage(): Language {
	if (!browser) return 'en';
	
	const stored = localStorage.getItem('language') as Language | null;
	if (stored && translations[stored]) {
		return stored;
	}
	
	// Try to detect from browser
	const browserLang = navigator.language.toLowerCase();
	if (browserLang.startsWith('pt')) return 'pt-BR';
	if (browserLang.startsWith('fr')) return 'fr';
	if (browserLang.startsWith('es')) return 'es';
	if (browserLang.startsWith('de')) return 'de';
	if (browserLang.startsWith('ru')) return 'ru';
	if (browserLang.startsWith('ja')) return 'ja';
	if (browserLang.startsWith('ko')) return 'ko';
	if (browserLang.startsWith('zh')) return 'zh-CN';
	if (browserLang.startsWith('el')) return 'el';
	if (browserLang.startsWith('la')) return 'la';
	
	return 'en';
}

// Create language store
export const language = writable<Language>(getInitialLanguage());

// Subscribe to language changes and save to localStorage
if (browser) {
	language.subscribe((lang) => {
		localStorage.setItem('language', lang);
	});
}

// Create a derived store for translations that updates when language changes
export const translationsStore = derived(language, ($lang) => {
	return (key: string, params?: Record<string, string | number>): string => {
		return getTranslation($lang, key, params);
	};
});

// Translation function (reactive - use in components)
// Usage in template: {$t('hero.title')} where $t = $translationsStore
// Or use: $: text = $translationsStore('hero.title')
export function t(key: string, params?: Record<string, string | number>): string {
	const currentLang = get(language);
	return getTranslation(currentLang, key, params);
}

// Internal translation getter
function getTranslation(lang: Language, key: string, params?: Record<string, string | number>): string {
	const keys = key.split('.');
	let value: any = translations[lang];
	
	for (const k of keys) {
		value = value?.[k];
		if (value === undefined) {
			// Fallback to English
			value = translations['en'];
			for (const k2 of keys) {
				value = value?.[k2];
			}
			break;
		}
	}
	
	if (typeof value !== 'string') {
		return key; // Return key if translation not found
	}
	
	// Replace parameters
	if (params) {
		return value.replace(/\{\{(\w+)\}\}/g, (match, param) => {
			return params[param]?.toString() || match;
		});
	}
	
	return value;
}

// Get current translations object
export function getTranslations(): Translations {
	const currentLang = get(language);
	return translations[currentLang];
}

// Language metadata
export const languages: Array<{ code: Language; name: string; flag: string }> = [
	{ code: 'en', name: 'English', flag: '🇺🇸' },
	{ code: 'pt-BR', name: 'Português (BR)', flag: '🇧🇷' },
	{ code: 'fr', name: 'Français', flag: '🇫🇷' },
	{ code: 'es', name: 'Español', flag: '🇪🇸' },
	{ code: 'de', name: 'Deutsch', flag: '🇩🇪' },
	{ code: 'ru', name: 'Русский', flag: '🇷🇺' },
	{ code: 'ja', name: '日本語', flag: '🇯🇵' },
	{ code: 'ko', name: '한국어', flag: '🇰🇷' },
	{ code: 'zh-CN', name: '中文 (简体)', flag: '🇨🇳' },
	{ code: 'el', name: 'Ελληνικά', flag: '🇬🇷' },
	{ code: 'la', name: 'Latina', flag: '🏛️' }
];

