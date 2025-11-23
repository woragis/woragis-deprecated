<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import Nav from '$lib/components/Nav.svelte';
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import { queryClient } from '$lib/query-client';
	
	let { children } = $props();
</script>

<svelte:head>
	<title>Woragis - Software Development Agency</title>
	<meta name="description" content="Woragis is a software development agency that specializes in building custom software solutions for businesses.">
	<meta name="keywords" content="software development, software agency, software solutions, software development agency, software development company, software development services, software development company in Brazil, software development company in São Paulo, software development company in Rio de Janeiro, software development company in Belo Horizonte, software development company in Curitiba, software development company in Porto Alegre, software development company in Brasília, software development company in Fortaleza, software development company in Recife, software development company in Salvador, software development company in João Pessoa, software development company in Maceió, software development company in Aracaju, software development company in João Pessoa, software development company in Maceió, software development company in Aracaju, software development company in João Pessoa, software development company in Maceió, software development company in Aracaju">
	<meta name="author" content="Woragis">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="robots" content="index, follow">
	<meta name="googlebot" content="index, follow">
	<meta name="bingbot" content="index, follow">
	<meta name="yandexbot" content="index, follow">
	<link rel="icon" href={favicon} />
</svelte:head>

<QueryClientProvider client={queryClient}>
	<Nav />
	<div class="pt-16">
		{@render children()}
	</div>
</QueryClientProvider>
