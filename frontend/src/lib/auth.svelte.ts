import { fetchMe, logout as apiLogout } from '$lib/api/client';
import type { User } from '$lib/api/types';

/**
 * App-wide authenticated session, exposed as a Svelte 5 rune state.
 *
 * Usage:
 *
 *   import { auth, ensureAuthLoaded } from '$lib/auth.svelte';
 *   // in <script>: const user = $derived(auth.user);
 *
 * The store is intentionally not persisted to localStorage; the source of
 * truth is the HttpOnly session cookie set by the backend, and we re-hydrate
 * by calling /api/v1/auth/me on mount.
 */
class AuthStore {
	user = $state<User | null>(null);
	loading = $state(false);
	loaded = $state(false);
	error = $state<string | null>(null);

	async load(force = false): Promise<User | null> {
		if (this.loaded && !force) return this.user;
		this.loading = true;
		this.error = null;
		try {
			this.user = await fetchMe();
			return this.user;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to load session';
			this.user = null;
			return null;
		} finally {
			this.loading = false;
			this.loaded = true;
		}
	}

	async logout(): Promise<void> {
		this.user = null;
		this.loaded = true;
		await apiLogout(); // redirects browser to Authelia logout
	}
}

export const auth = new AuthStore();

/** Convenience helper for routes that need the user before rendering. */
export async function ensureAuthLoaded(): Promise<User | null> {
	return auth.load();
}
