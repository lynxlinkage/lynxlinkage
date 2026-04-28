import type {
	Application,
	ApplicationListFilter,
	ApplicationStatus,
	ApplicationStatusUpsertPayload,
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

export interface SubmitApplicationInput {
	name: string;
	email: string;
	message?: string;
	files: File[];
}

export interface SubmitApplicationResult {
	ok: boolean;
	id?: number;
	files?: number;
	error?: string;
}

/**
 * Submits a candidate application as multipart/form-data. We don't go
 * through `request()` because that helper forces JSON content-type and
 * we need the browser to set its own multipart boundary.
 */
export async function submitApplication(
	jobId: number,
	input: SubmitApplicationInput
): Promise<SubmitApplicationResult> {
	const fd = new FormData();
	fd.set('name', input.name);
	fd.set('email', input.email);
	if (input.message) fd.set('message', input.message);
	for (const f of input.files) fd.append('files', f, f.name);

	try {
		const res = await fetch(`/api/v1/jobs/${jobId}/applications`, {
			method: 'POST',
			body: fd
		});
		if (res.status === 429) {
			return { ok: false, error: 'Too many submissions. Please try again in a minute.' };
		}
		if (!res.ok) {
			const body = await res.json().catch(() => ({}));
			return {
				ok: false,
				error: (body as { error?: string })?.error ?? `Submission failed (${res.status}).`
			};
		}
		const body = (await res.json()) as { id?: number; files?: number };
		return { ok: true, id: body.id, files: body.files };
	} catch (err) {
		return { ok: false, error: err instanceof Error ? err.message : 'Network error' };
	}
}

export async function adminListApplications(
	filter: ApplicationListFilter = {}
): Promise<Application[]> {
	const qs = new URLSearchParams();
	if (filter.jobId) qs.set('jobId', String(filter.jobId));
	if (filter.statusId) qs.set('statusId', String(filter.statusId));
	if (filter.sort) qs.set('sort', filter.sort);
	if (filter.limit) qs.set('limit', String(filter.limit));
	const path = `/api/v1/admin/applications${qs.size ? `?${qs.toString()}` : ''}`;
	const body = await request<ListResponse<Application>>(path, { method: 'GET' });
	return body.items ?? [];
}

export async function adminGetApplication(id: number): Promise<Application> {
	return request<Application>(`/api/v1/admin/applications/${id}`, { method: 'GET' });
}

export async function adminUpdateApplicationStatus(
	applicationId: number,
	statusId: number,
	note?: string
): Promise<Application> {
	return request<Application>(`/api/v1/admin/applications/${applicationId}/status`, {
		method: 'PUT',
		body: JSON.stringify({ statusId, note: note ?? '' })
	});
}

/**
 * Returns the URL to the download endpoint for an application file.
 * The browser will follow same-origin cookies automatically when the
 * URL is hit via an <a href> click or window.open.
 */
export function applicationFileUrl(applicationId: number, fileId: number): string {
	return `/api/v1/admin/applications/${applicationId}/files/${fileId}`;
}

export async function adminListStatuses(): Promise<ApplicationStatus[]> {
	const body = await request<ListResponse<ApplicationStatus>>('/api/v1/admin/application-statuses', {
		method: 'GET'
	});
	return body.items ?? [];
}

export async function adminCreateStatus(
	payload: ApplicationStatusUpsertPayload
): Promise<ApplicationStatus> {
	return request<ApplicationStatus>('/api/v1/admin/application-statuses', {
		method: 'POST',
		body: JSON.stringify(payload)
	});
}

export async function adminUpdateStatus(
	id: number,
	payload: ApplicationStatusUpsertPayload
): Promise<ApplicationStatus> {
	return request<ApplicationStatus>(`/api/v1/admin/application-statuses/${id}`, {
		method: 'PUT',
		body: JSON.stringify(payload)
	});
}

export async function adminDeleteStatus(id: number): Promise<void> {
	await request<void>(`/api/v1/admin/application-statuses/${id}`, { method: 'DELETE' });
}
