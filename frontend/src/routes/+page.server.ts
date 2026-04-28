import type { PageServerLoad } from './$types';
import { loadJobs, loadPartners } from '$lib/api/server';

export const prerender = true;

export const load: PageServerLoad = async ({ fetch }) => {
	// Fetch in parallel; tolerate failures so the home page can still render
	// even when one section is empty.
	const [jobs, partners] = await Promise.all([
		loadJobs(fetch).catch(() => []),
		loadPartners(fetch).catch(() => [])
	]);

	return {
		jobs: jobs.slice(0, 2),
		partners: partners.slice(0, 8)
	};
};
