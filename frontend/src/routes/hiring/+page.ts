import type { PageLoad } from './$types';
import type { JobPosting, ListResponse } from '$lib/api/types';

export const prerender = false;
export const ssr = false;

export const load: PageLoad = async ({ fetch }) => {
	try {
		const res = await fetch('/api/v1/jobs');
		if (!res.ok) return { jobs: [] as JobPosting[] };
		const data = (await res.json()) as ListResponse<JobPosting>;
		return { jobs: data.items ?? [] };
	} catch {
		return { jobs: [] as JobPosting[] };
	}
};
