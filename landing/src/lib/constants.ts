export const contact = {
	email: 'jezreel.veloso@gmail.com',
	github: 'https://github.com/woragis',
	instagram: 'https://instagram.com/y.jezreel.andrade',
	linkedin: 'https://linkedin.com/in/jezreel-andrade',
	phone: '+55 83 99691-2887',
	whatsapp: '+55 83 99691-2887',
	location: 'João Pessoa, Brazil'
} as const;

export const api = {
	baseURL: import.meta.env.PUBLIC_API_BASE_URL || 'http://localhost:8080/api',
	timeout: 10000
} as const;

