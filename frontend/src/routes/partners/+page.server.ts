import type { PageServerLoad } from './$types';
import { loadPartners } from '$lib/api/server';

export const prerender = true;

export const load: PageServerLoad = async ({ fetch }) => {
	const partners = await loadPartners(fetch).catch(() => []);
	return { partners };
};
