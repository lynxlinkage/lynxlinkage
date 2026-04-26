import type { PageServerLoad } from './$types';
import { loadJobs, loadPartners, loadResearches } from '$lib/api/server';

export const prerender = true;

export const load: PageServerLoad = async ({ fetch }) => {
	// Fetch in parallel; tolerate failures so the home page can still render
	// even when one section is empty.
	const [researches, jobs, partners] = await Promise.all([
		loadResearches(fetch, { limit: 3 }).catch(() => []),
		loadJobs(fetch).catch(() => []),
		loadPartners(fetch).catch(() => [])
	]);

	return {
		researches,
		jobs: jobs.slice(0, 2),
		partners: partners.slice(0, 8)
	};
};
