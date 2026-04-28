import type { JobPosting, ListResponse, Partner, ResearchCard } from './types';

/**
 * Server-side API helpers used by `+page.server.ts` load functions during
 * prerendering. These run on the build machine (or in `vite dev`) and call
 * the Go backend directly via BACKEND_URL.
 */

const BACKEND_URL = process.env.BACKEND_URL ?? 'http://localhost:8080';

type FetchLike = typeof fetch;

async function getJSON<T>(path: string, fetchFn: FetchLike): Promise<T> {
	const url = new URL(path, BACKEND_URL).toString();
	const res = await fetchFn(url);
	if (!res.ok) {
		throw new Error(`GET ${path} failed: ${res.status} ${res.statusText}`);
	}
	return (await res.json()) as T;
}

export async function loadResearches(
	fetchFn: FetchLike,
	opts: { limit?: number; tag?: string } = {}
): Promise<ResearchCard[]> {
	const params = new URLSearchParams();
	if (opts.limit) params.set('limit', String(opts.limit));
	if (opts.tag) params.set('tag', opts.tag);
	const qs = params.toString();
	const data = await getJSON<ListResponse<ResearchCard>>(
		`/api/v1/researches${qs ? `?${qs}` : ''}`,
		fetchFn
	);
	return data.items ?? [];
}

export async function loadJobs(fetchFn: FetchLike): Promise<JobPosting[]> {
	const data = await getJSON<ListResponse<JobPosting>>('/api/v1/jobs', fetchFn);
	return data.items ?? [];
}

export async function loadJob(fetchFn: FetchLike, id: string | number): Promise<JobPosting | null> {
	const url = new URL(`/api/v1/jobs/${id}`, BACKEND_URL).toString();
	const res = await fetchFn(url);
	if (res.status === 404) return null;
	if (!res.ok) {
		throw new Error(`GET /api/v1/jobs/${id} failed: ${res.status} ${res.statusText}`);
	}
	return (await res.json()) as JobPosting;
}

export async function loadPartners(fetchFn: FetchLike): Promise<Partner[]> {
	const data = await getJSON<ListResponse<Partner>>('/api/v1/partners', fetchFn);
	return data.items ?? [];
}
