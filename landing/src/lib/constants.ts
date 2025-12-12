import { PUBLIC_API_BASE_URL, PUBLIC_API_KEY } from '$env/static/public';

export const contact = {
	email: 'jezreel.veloso@gmail.com',
	github: 'https://github.com/woragis',
	instagram: 'https://instagram.com/y.jezreel.andrade',
	linkedin: 'https://linkedin.com/in/jezreel-andrade',
	phone: '+55 83 99691-2887',
	whatsapp: '+55 83 99691-2887',
	location: 'João Pessoa, Brazil'
} as const;

// User ID for resume downloads (from backend)
export const userId = '6ad0d828-f605-45fc-a545-3441e17a015c';

export const api = {
	baseURL: PUBLIC_API_BASE_URL || 'http://localhost:8080/api',
	timeout: 10000,
	apiKey: PUBLIC_API_KEY || null
} as const;

export { skills } from './constants/skills';
export { interests } from './constants/interests';
export { caseStudies, systemDesigns, problemSolutions } from './constants/technical';

