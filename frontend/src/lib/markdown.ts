import { marked } from 'marked';

marked.setOptions({
	gfm: true,
	breaks: false
});

/**
 * Render a Markdown string to HTML.
 *
 * The seed YAML is authored by us, so the input is trusted; we render at
 * build time and inject the result with `{@html}` in prerendered pages.
 * Do **not** call this on user-submitted content without sanitising first.
 */
export function renderMarkdown(md: string): string {
	if (!md) return '';
	return marked.parse(md, { async: false }) as string;
}
