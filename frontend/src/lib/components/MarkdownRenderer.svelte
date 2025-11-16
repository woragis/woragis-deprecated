<script lang="ts">
	import { renderMarkdown } from '$lib/utils/markdown';

	export let content: string;
	export let className = 'prose prose-invert max-w-none';

	$: html = renderMarkdown(content || '');
</script>

<div class={`markdown-content ${className}`}>
	{@html html}
</div>

<style>
	@reference "tailwindcss";

	:global(.markdown-content) {
		@apply text-slate-200;
	}

	:global(.markdown-content p) {
		@apply mb-3 leading-relaxed;
	}

	:global(.markdown-content strong) {
		@apply font-semibold text-white;
	}

	:global(.markdown-content em) {
		@apply italic text-slate-300;
	}

	:global(.markdown-content code) {
		@apply font-mono;
	}

	:global(.markdown-content pre) {
		@apply bg-slate-900 rounded-lg p-3 overflow-x-auto my-3;
	}

	:global(.markdown-content pre code) {
		@apply bg-transparent p-0;
	}

	:global(.markdown-content a) {
		@apply text-blue-400 hover:text-blue-300 underline;
	}

	:global(.markdown-content ul) {
		@apply list-disc list-inside space-y-1 my-2;
	}

	:global(.markdown-content ol) {
		@apply list-decimal list-inside space-y-1 my-2;
	}

	:global(.markdown-content blockquote) {
		@apply border-l-4 border-slate-600 pl-4 my-3 italic text-slate-400;
	}

	:global(.markdown-content h1) {
		@apply text-2xl font-bold mt-6 mb-3 text-white;
	}

	:global(.markdown-content h2) {
		@apply text-xl font-semibold mt-5 mb-2 text-white;
	}

	:global(.markdown-content h3) {
		@apply text-lg font-semibold mt-4 mb-2 text-slate-100;
	}
</style>

