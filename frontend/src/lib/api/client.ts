import type { ContactPayload } from './types';

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
