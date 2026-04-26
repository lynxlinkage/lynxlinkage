/**
 * Static site-wide configuration. Edit these values; they're read by the
 * header, footer, SEO meta tags and structured-data blocks.
 */
export const site = {
	name: 'lynxlinkage',
	tagline: 'Statistical arbitrage in crypto markets.',
	description:
		'lynxlinkage is an agile team of researchers and engineers building systematic statistical-arbitrage strategies for crypto markets. We publish what we learn, hire technically deep people, and partner with the venues and infrastructure we depend on.',
	contactEmail: 'hello@lynxlinkage.com',
	careersEmail: 'careers@lynxlinkage.com',
	pressEmail: 'press@lynxlinkage.com',
	url: 'https://lynxlinkage.com',
	github: 'https://github.com/lynxlinkage',
	founded: 2024,
	locations: ['Remote', 'Taiwan'],
	nav: [
		{ href: '/about', label: 'About us' },
		{ href: '/researches', label: 'Public researches' },
		{ href: '/hiring', label: 'Hiring' },
		{ href: '/partners', label: 'Partners' }
	]
} as const;

export type NavItem = (typeof site.nav)[number];
