import { marked } from 'marked';
import { browser } from '$app/environment';
import DOMPurify from 'dompurify';

// Configure marked options
marked.setOptions({
	breaks: true, // Convert line breaks to <br>
	gfm: true, // GitHub Flavored Markdown
});

/**
 * Renders markdown to sanitized HTML
 * @param markdown - The markdown string to render
 * @returns Sanitized HTML string
 */
export function renderMarkdown(markdown: string): string {
	if (!markdown) return '';
	
	// Only render in browser (DOMPurify requires browser environment)
	if (!browser) {
		// Return escaped text for SSR
		return markdown.replace(/</g, '&lt;').replace(/>/g, '&gt;');
	}
	
	try {
		// Parse markdown to HTML
		const html = marked.parse(markdown) as string;
		
		// Sanitize HTML to prevent XSS attacks
		const sanitized = DOMPurify.sanitize(html, {
			ALLOWED_TAGS: [
				'p', 'br', 'strong', 'em', 'u', 's', 'code', 'pre',
				'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
				'ul', 'ol', 'li', 'blockquote', 'hr',
				'a', 'table', 'thead', 'tbody', 'tr', 'th', 'td'
			],
			ALLOWED_ATTR: ['href', 'title', 'target', 'rel'],
			ALLOWED_URI_REGEXP: /^(?:(?:(?:f|ht)tps?|mailto|tel|callto|sms|cid|xmpp):|[^a-z]|[a-z+.\-]+(?:[^a-z+.\-:]|$))/i
		});
		
		return sanitized;
	} catch (error) {
		console.error('Error rendering markdown:', error);
		// Return escaped text if parsing fails
		return markdown.replace(/</g, '&lt;').replace(/>/g, '&gt;');
	}
}
