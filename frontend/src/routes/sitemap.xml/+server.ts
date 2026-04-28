import { site } from '$lib/site';

export const prerender = true;

const PATHS = ['/', '/about', '/hiring', '/partners'] as const;

export function GET() {
	const today = new Date().toISOString().slice(0, 10);
	const urls = PATHS.map(
		(p) => `<url><loc>${new URL(p, site.url).toString()}</loc><lastmod>${today}</lastmod></url>`
	).join('');
	const body = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${urls}</urlset>`;
	return new Response(body, {
		headers: {
			'Content-Type': 'application/xml',
			'Cache-Control': 'public, max-age=3600'
		}
	});
}
