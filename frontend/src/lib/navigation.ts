export type NavItem = {
	href: string;
	label: string;
	icon?: string;
	match?: (pathname: string) => boolean;
};

const matchStartsWith =
	(base: string) =>
	(pathname: string) =>
		pathname === base || pathname.startsWith(`${base}/`);

export const primaryNav: NavItem[] = [
	{ href: '/', label: 'Home', match: (pathname) => pathname === '/' },
	{ href: '/finances', label: 'Finances', match: matchStartsWith('/finances') },
	{ href: '/chats', label: 'Chats', match: matchStartsWith('/chats') },
	{ href: '/ideas', label: 'Ideas', match: matchStartsWith('/ideas') },
	{ href: '/projects', label: 'Projects', match: matchStartsWith('/projects') },
	{ href: '/clients', label: 'Clients', match: matchStartsWith('/clients') },
	{ href: '/reports', label: 'Reports', match: matchStartsWith('/reports') },
	{ href: '/schedules', label: 'Schedules', match: matchStartsWith('/schedules') },
	{ href: '/whatsapp', label: 'WhatsApp', match: matchStartsWith('/whatsapp') },
	{ href: '/monitoring', label: 'Monitoring', match: matchStartsWith('/monitoring') }
];

export const authNav: NavItem[] = [
	{ href: '/auth/login', label: 'Sign in', match: matchStartsWith('/auth/login') },
	{ href: '/auth/register', label: 'Register', match: matchStartsWith('/auth/register') },
{ href: '/auth/confirm-email', label: 'Confirm email', match: matchStartsWith('/auth/confirm-email') },
	{ href: '/auth/forgot', label: 'Password reset', match: matchStartsWith('/auth/forgot') },
	{ href: '/auth/profile', label: 'Profile', match: matchStartsWith('/auth/profile') },
	{ href: '/auth/mfa', label: 'MFA settings', match: matchStartsWith('/auth/mfa') },
	{ href: '/auth/sessions', label: 'Active sessions', match: matchStartsWith('/auth/sessions') },
	{ href: '/auth/connections', label: 'OAuth connections', match: matchStartsWith('/auth/connections') }
];

