import type {
	ContactPayload,
	JobPosting,
	JobUpsertPayload,
	ListResponse,
	User
} from './types';

export interface ContactResult {
	ok: boolean;
	id?: number;
	error?: string;
}

/**
 * Submits a contact form to the backend. Runs in the browser only and uses
 * a relative URL so it works in dev (Vite proxies /api to the Go server)
 * and in production (Go serves both the static frontend and the API on
 * the same origin).
 */
export async function submitContact(payload: ContactPayload): Promise<ContactResult> {
	try {
		const res = await fetch('/api/v1/contact', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload)
		});
		if (res.status === 429) {
			return { ok: false, error: 'Too many submissions. Please try again in a minute.' };
		}
		if (!res.ok) {
			const body = await res.json().catch(() => ({}));
			return { ok: false, error: body?.error ?? `Submission failed (${res.status}).` };
		}
		const body = (await res.json()) as { id?: number };
		return { ok: true, id: body.id };
	} catch (err) {
		return { ok: false, error: err instanceof Error ? err.message : 'Network error' };
	}
}

/**
 * Thrown by API helpers below when the server returns a non-2xx status.
 * The caller can read .status to distinguish 401/403 from generic errors.
 */
export class ApiError extends Error {
	status: number;
	constructor(message: string, status: number) {
		super(message);
		this.status = status;
		this.name = 'ApiError';
	}
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
	const res = await fetch(path, {
		credentials: 'same-origin',
		...init,
		headers: {
			'Content-Type': 'application/json',
			Accept: 'application/json',
			...(init.headers ?? {})
		}
	});
	if (!res.ok) {
		const body = await res.json().catch(() => ({}));
		throw new ApiError(
			(body as { error?: string })?.error ?? `${res.status} ${res.statusText}`,
			res.status
		);
	}
	if (res.status === 204) return undefined as T;
	return (await res.json()) as T;
}

export async function login(email: string, password: string): Promise<User> {
	const body = await request<{ user: User }>('/api/v1/auth/login', {
		method: 'POST',
		body: JSON.stringify({ email, password })
	});
	return body.user;
}

export async function logout(): Promise<void> {
	await request<unknown>('/api/v1/auth/logout', { method: 'POST' });
}

export async function fetchMe(): Promise<User | null> {
	try {
		const body = await request<{ user: User }>('/api/v1/auth/me', { method: 'GET' });
		return body.user;
	} catch (err) {
		if (err instanceof ApiError && (err.status === 401 || err.status === 403)) return null;
		throw err;
	}
}

export async function adminListJobs(): Promise<JobPosting[]> {
	const body = await request<ListResponse<JobPosting>>('/api/v1/admin/jobs', { method: 'GET' });
	return body.items ?? [];
}

export async function adminCreateJob(payload: JobUpsertPayload): Promise<JobPosting> {
	return request<JobPosting>('/api/v1/admin/jobs', {
		method: 'POST',
		body: JSON.stringify(payload)
	});
}

export async function adminUpdateJob(
	id: number,
	payload: JobUpsertPayload
): Promise<JobPosting> {
	return request<JobPosting>(`/api/v1/admin/jobs/${id}`, {
		method: 'PUT',
		body: JSON.stringify(payload)
	});
}

export async function adminDeleteJob(id: number): Promise<void> {
	await request<void>(`/api/v1/admin/jobs/${id}`, { method: 'DELETE' });
}
