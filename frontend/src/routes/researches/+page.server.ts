import type { PageServerLoad } from './$types';
import { loadResearches } from '$lib/api/server';

export const prerender = true;

export const load: PageServerLoad = async ({ fetch }) => {
	const researches = await loadResearches(fetch).catch(() => []);
	return { researches };
};
