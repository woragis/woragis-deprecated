import { QueryClient } from '@tanstack/svelte-query';

export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			retry: 1,
			refetchOnWindowFocus: false
		}
	}
});
