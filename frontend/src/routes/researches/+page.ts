import { error, redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { ListResponse, ResearchCard } from '$lib/api/types';

export const load: PageLoad = async ({ fetch, url }) => {
	const meRes = await fetch('/api/v1/auth/me');
	if (!meRes.ok) {
		throw redirect(
			303,
			`/login?next=${encodeURIComponent(url.pathname + url.search)}`
		);
	}
	const meBody = (await meRes.json()) as { user?: unknown };
	if (!meBody?.user) {
		throw redirect(
			303,
			`/login?next=${encodeURIComponent(url.pathname + url.search)}`
		);
	}

	const listRes = await fetch('/api/v1/researches');
	if (!listRes.ok) {
		throw error(listRes.status >= 500 ? 502 : listRes.status, 'Unable to load researches');
	}
	const data = (await listRes.json()) as ListResponse<ResearchCard>;
	return { researches: data.items ?? [] };
};
