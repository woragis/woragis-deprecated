import { browser } from '$app/environment';

const FINGERPRINT_STORAGE_KEY = 'woragis_device_fingerprint';

const generateFingerprint = () => {
	if (!browser) {
		return '';
	}
	return globalThis.crypto?.randomUUID?.() ?? Math.random().toString(36).slice(2);
};

export const getDeviceFingerprint = () => {
	if (!browser) {
		return '';
	}

	try {
		let fingerprint = window.localStorage.getItem(FINGERPRINT_STORAGE_KEY);
		if (!fingerprint) {
			fingerprint = generateFingerprint();
			window.localStorage.setItem(FINGERPRINT_STORAGE_KEY, fingerprint);
		}
		return fingerprint;
	} catch (error) {
		console.warn('Unable to access localStorage for device fingerprint.', error);
		return generateFingerprint();
	}
};

export const getDeviceName = () => {
	if (!browser) {
		return 'Browser';
	}

	const platform = navigator.platform || 'Unknown platform';
	const vendor = navigator.vendor || 'Generic';
	return `${vendor} ${platform}`.trim();
};

export const getUserAgent = () => {
	if (!browser) {
		return '';
	}
	return navigator.userAgent ?? '';
};

