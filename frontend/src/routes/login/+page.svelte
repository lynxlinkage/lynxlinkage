<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';

	onMount(() => {
		// Login is handled by Authelia. Redirect there, preserving the
		// intended destination so Authelia can redirect back after login.
		const next = page.url?.searchParams.get('next') ?? '/admin';
		window.location.href = `/admin`;
		// Full page load to /admin triggers Traefik → Authelia auth check,
		// which redirects to auth.lynxlinkage.com if not yet authenticated.
		void next;
	});
</script>

<svelte:head>
	<title>Sign in · Lynxlinkage</title>
	<meta name="robots" content="noindex,nofollow" />
</svelte:head>

<div class="splash">
	<p>Redirecting to login…</p>
</div>

<style>
	.splash {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100dvh;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
		color: #666;
		font-size: 14px;
	}
</style>
