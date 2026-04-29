import type { PageLoad } from './$types';
import type { Partner, ListResponse } from '$lib/api/types';

export const prerender = false;
export const ssr = false;

export const load: PageLoad = async ({ fetch }) => {
	try {
		const res = await fetch('/api/v1/partners');
		if (!res.ok) return { partners: [] as Partner[] };
		const data = (await res.json()) as ListResponse<Partner>;
		return { partners: data.items ?? [] };
	} catch {
		return { partners: [] as Partner[] };
	}
};
