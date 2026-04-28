import type { PageServerLoad } from './$types';
import { loadJobs } from '$lib/api/server';

export const prerender = true;

export const load: PageServerLoad = async ({ fetch }) => {
	const jobs = await loadJobs(fetch).catch(() => []);
	return { jobs };
};
