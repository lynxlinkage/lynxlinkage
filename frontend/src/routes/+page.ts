import type { PageLoad } from './$types';
import type { JobPosting, ListResponse, Partner } from '$lib/api/types';

export const prerender = false;
export const ssr = false;

export const load: PageLoad = async ({ fetch }) => {
	const [jobs, partners] = await Promise.all([
		fetch('/api/v1/jobs')
			.then((r) => (r.ok ? (r.json() as Promise<ListResponse<JobPosting>>) : { items: [] }))
			.then((d) => (d.items ?? []).slice(0, 2))
			.catch(() => [] as JobPosting[]),
		fetch('/api/v1/partners')
			.then((r) => (r.ok ? (r.json() as Promise<ListResponse<Partner>>) : { items: [] }))
			.then((d) => (d.items ?? []).slice(0, 8))
			.catch(() => [] as Partner[])
	]);

	return { jobs, partners };
};
