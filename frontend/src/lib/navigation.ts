export type NavItem = {
	href: string;
	label: string;
	icon?: string;
	match?: (pathname: string) => boolean;
	children?: NavItem[];
};

const matchStartsWith =
	(base: string) =>
	(pathname: string) =>
		pathname === base || pathname.startsWith(`${base}/`);

export const primaryNav: NavItem[] = [
	{ href: '/', label: 'Home', match: (pathname) => pathname === '/' },
	{
		href: '/personal',
		label: 'Personal',
		match: (pathname) => pathname.startsWith('/personal'),
		children: [
			{ href: '/personal/finances', label: 'Finances', match: matchStartsWith('/personal/finances') },
			{ href: '/personal/ideas', label: 'Ideas', match: matchStartsWith('/personal/ideas') }
		]
	},
	{ href: '/chats', label: 'Chats', match: matchStartsWith('/chats') },
	{ href: '/projects', label: 'Projects', match: matchStartsWith('/projects') },
	{ href: '/clients', label: 'Clients', match: matchStartsWith('/clients') },
	{ href: '/reports', label: 'Reports', match: matchStartsWith('/reports') },
	{ href: '/schedules', label: 'Schedules', match: matchStartsWith('/schedules') },
	{ href: '/whatsapp', label: 'WhatsApp', match: matchStartsWith('/whatsapp') },
	{ href: '/monitoring', label: 'Monitoring', match: matchStartsWith('/monitoring') },
	{
		href: '/landing',
		label: 'Landing Pages',
		match: (pathname) => pathname.startsWith('/landing'),
		children: [
			{ href: '/landing/case-studies', label: 'Case Studies', match: matchStartsWith('/landing/case-studies') },
			{ href: '/landing/certifications', label: 'Certifications', match: matchStartsWith('/landing/certifications') },
			{ href: '/landing/posts', label: 'Posts', match: matchStartsWith('/landing/posts') },
			{ href: '/landing/problem-solutions', label: 'Problem Solutions', match: matchStartsWith('/landing/problem-solutions') },
			{ href: '/landing/skills', label: 'Skills', match: matchStartsWith('/landing/skills') },
			{ href: '/landing/social-media-posts', label: 'Social Media Posts', match: matchStartsWith('/landing/social-media-posts') },
			{ href: '/landing/system-designs', label: 'System Designs', match: matchStartsWith('/landing/system-designs') },
			{ href: '/landing/technical-writings', label: 'Technical Writings', match: matchStartsWith('/landing/technical-writings') },
			{ href: '/landing/testimonials', label: 'Testimonials', match: matchStartsWith('/landing/testimonials') }
		]
	},
	{ href: '/resumes', label: 'Resumes', match: matchStartsWith('/resumes') }
];

export const authNav: NavItem[] = [
	{ href: '/auth/login', label: 'Sign in', match: matchStartsWith('/auth/login') },
	{ href: '/auth/register', label: 'Register', match: matchStartsWith('/auth/register') },
{ href: '/auth/confirm-email', label: 'Confirm email', match: matchStartsWith('/auth/confirm-email') },
	{ href: '/auth/forgot', label: 'Password reset', match: matchStartsWith('/auth/forgot') },
	{ href: '/auth/profile', label: 'Profile', match: matchStartsWith('/auth/profile') },
	{ href: '/auth/mfa', label: 'MFA settings', match: matchStartsWith('/auth/mfa') },
	{ href: '/auth/sessions', label: 'Active sessions', match: matchStartsWith('/auth/sessions') },
	{ href: '/auth/connections', label: 'OAuth connections', match: matchStartsWith('/auth/connections') },
	{ href: '/api-keys', label: 'API Keys', match: matchStartsWith('/api-keys') }
];

