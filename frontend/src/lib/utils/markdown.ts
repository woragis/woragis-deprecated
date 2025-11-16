/**
 * Simple markdown renderer for basic formatting
 * Supports: **bold**, *italic*, \n for line breaks, and other common markdown features
 */

export interface MarkdownOptions {
	/** Whether to convert line breaks to <br> tags */
	lineBreaks?: boolean;
	/** Additional CSS class for the rendered content */
	className?: string;
}

/**
 * Renders basic markdown text to HTML string
 * Supports:
 * - **bold** and __bold__
 * - *italic* and _italic_
 * - `code`
 * - \n for line breaks
 */
export function renderMarkdown(text: string, options: MarkdownOptions = {}): string {
	if (!text) return '';

	const { lineBreaks = true } = options;

	// Escape HTML first to prevent XSS
	let html = escapeHtml(text);

	// Convert line breaks to <br>
	if (lineBreaks) {
		html = html.replace(/\n/g, '<br>');
	}

	// Code blocks (```code```)
	html = html.replace(/```([\s\S]*?)```/g, '<pre><code>$1</code></pre>');

	// Inline code (`code`)
	html = html.replace(/`([^`]+)`/g, '<code class="px-1 py-0.5 rounded bg-slate-800 text-primary font-mono text-sm">$1</code>');

	// Bold (**text** or __text__)
	html = html.replace(/\*\*(.+?)\*\*/g, '<strong class="font-semibold">$1</strong>');
	html = html.replace(/__(.+?)__/g, '<strong class="font-semibold">$1</strong>');

	// Italic (*text* or _text_) - but not if part of bold or code
	html = html.replace(/(?<!<[^>]*)(?<!\*)\*(?!\*)([^*\n]+?)(?<!\*)\*(?!\*)(?!<[^>]*>)/g, '<em>$1</em>');
	html = html.replace(/(?<!<[^>]*)(?<!_)_(?!_)([^_\n]+?)(?<!_)_(?!_)(?!<[^>]*>)/g, '<em>$1</em>');

	// Headers
	html = html.replace(/^### (.+)$/gm, '<h3 class="text-lg font-semibold mt-4 mb-2">$1</h3>');
	html = html.replace(/^## (.+)$/gm, '<h2 class="text-xl font-semibold mt-4 mb-2">$1</h2>');
	html = html.replace(/^# (.+)$/gm, '<h1 class="text-2xl font-semibold mt-4 mb-2">$1</h1>');

	// Links [text](url)
	html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" class="text-primary hover:underline" target="_blank" rel="noopener noreferrer">$1</a>');

	// Lists (simple support)
	html = html.replace(/^\* (.+)$/gm, '<li class="ml-4">$1</li>');
	html = html.replace(/^- (.+)$/gm, '<li class="ml-4">$1</li>');
	// Wrap consecutive list items in <ul>
	html = html.replace(/(<li[^>]*>.*?<\/li>(?:\s*<li[^>]*>.*?<\/li>)*)/gs, '<ul class="list-disc list-inside my-2 space-y-1">$1</ul>');

	// Blockquotes
	html = html.replace(/^> (.+)$/gm, '<blockquote class="border-l-4 border-slate-600 pl-4 my-2 italic text-slate-400">$1</blockquote>');

	return html;
}

/**
 * Escape HTML to prevent XSS attacks
 */
function escapeHtml(text: string): string {
	const map: Record<string, string> = {
		'&': '&amp;',
		'<': '&lt;',
		'>': '&gt;',
		'"': '&quot;',
		"'": '&#039;'
	};
	return text.replace(/[&<>"']/g, (m) => map[m]);
}

/**
 * Strips markdown formatting from text (useful for previews)
 */
export function stripMarkdown(text: string): string {
	if (!text) return '';

	// Remove code blocks
	text = text.replace(/```[\s\S]*?```/g, '');
	// Remove inline code
	text = text.replace(/`[^`]+`/g, '');
	// Remove bold/italic
	text = text.replace(/\*\*([^*]+)\*\*/g, '$1');
	text = text.replace(/__([^_]+)__/g, '$1');
	text = text.replace(/\*([^*]+)\*/g, '$1');
	text = text.replace(/_([^_]+)_/g, '$1');
	// Remove headers
	text = text.replace(/^#{1,6}\s+/gm, '');
	// Remove links
	text = text.replace(/\[([^\]]+)\]\([^)]+\)/g, '$1');
	// Remove list markers
	text = text.replace(/^[\*\-\+]\s+/gm, '');

	return text.trim();
}

