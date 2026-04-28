import type { EntryGenerator, PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';
import { loadJob, loadJobs } from '$lib/api/server';
import { renderMarkdown } from '$lib/markdown';

export const prerender = true;

export const entries: EntryGenerator = async () => {
	const jobs = await loadJobs(fetch).catch(() => []);
	return jobs.map((j) => ({ id: String(j.id) }));
};

export const load: PageServerLoad = async ({ fetch, params }) => {
	const job = await loadJob(fetch, params.id).catch(() => null);
	if (!job) {
		throw error(404, 'Job posting not found');
	}
	return {
		job,
		descriptionHtml: renderMarkdown(job.descriptionMd)
	};
};
