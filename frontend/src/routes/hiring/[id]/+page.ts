import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { fetchJob } from '$lib/api/client';
import { renderMarkdown } from '$lib/markdown';

export const load: PageLoad = async ({ params }) => {
	const job = await fetchJob(params.id).catch(() => null);
	if (!job) {
		throw error(404, 'Job posting not found');
	}
	return {
		job,
		descriptionHtml: renderMarkdown(job.descriptionMd)
	};
};
