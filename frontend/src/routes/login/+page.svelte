<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { auth } from '$lib/auth.svelte';
	import { onMount } from 'svelte';

	let email = $state('');
	let password = $state('');
	let submitting = $state(false);
	let errorMsg = $state<string | null>(null);

	const next = $derived(page.url?.searchParams.get('next') ?? '/admin');

	onMount(async () => {
		const u = await auth.load();
		if (u) {
			void goto(next, { replaceState: true });
		}
	});

	async function onSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (submitting) return;
		errorMsg = null;
		submitting = true;
		try {
			await auth.login(email.trim(), password);
			void goto(next, { replaceState: true });
		} catch (err) {
			errorMsg = err instanceof Error ? err.message : 'Login failed';
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:head>
	<title>Sign in · Lynxlinkage admin</title>
	<meta name="robots" content="noindex,nofollow" />
</svelte:head>

<section class="login">
	<div class="login__card container">
		<header>
			<h1>Sign in</h1>
			<p class="muted">Restricted to authorised users (HR &amp; admin team).</p>
		</header>

		<form onsubmit={onSubmit} novalidate>
			<label class="field">
				<span>Email</span>
				<input
					type="email"
					autocomplete="username"
					required
					bind:value={email}
					disabled={submitting}
				/>
			</label>

			<label class="field">
				<span>Password</span>
				<input
					type="password"
					autocomplete="current-password"
					required
					minlength="1"
					bind:value={password}
					disabled={submitting}
				/>
			</label>

			{#if errorMsg}
				<p class="error" role="alert">{errorMsg}</p>
			{/if}

			<button type="submit" class="primary" disabled={submitting}>
				{submitting ? 'Signing in…' : 'Sign in'}
			</button>
		</form>

		<p class="muted small">
			Contact your administrator if you don't have credentials yet.
		</p>
	</div>
</section>

<style>
	.login {
		min-height: calc(100dvh - var(--header-height));
		display: grid;
		place-items: center;
		padding: var(--space-7) var(--space-3);
		background: var(--surface-muted, #f7f8fb);
	}
	.login__card {
		width: min(420px, 100%);
		background: var(--bg);
		padding: var(--space-6);
		border-radius: var(--radius-lg);
		box-shadow: var(--shadow-md);
		display: grid;
		gap: var(--space-4);
	}
	header h1 {
		margin: 0 0 0.25rem;
		font-size: var(--text-3xl);
		letter-spacing: -0.01em;
	}
	.muted {
		color: var(--text-muted);
		margin: 0;
	}
	.small {
		font-size: var(--text-sm);
	}
	form {
		display: grid;
		gap: var(--space-3);
	}
	.field {
		display: grid;
		gap: 0.35rem;
	}
	.field span {
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--text);
	}
	input {
		width: 100%;
		padding: 0.7rem 0.85rem;
		font-size: var(--text-base);
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
		transition: border-color 120ms var(--ease-out);
	}
	input:focus {
		outline: none;
		border-color: var(--accent);
		box-shadow: 0 0 0 3px var(--accent-soft);
	}
	.primary {
		display: inline-flex;
		justify-content: center;
		align-items: center;
		padding: 0.75rem 1rem;
		border: none;
		border-radius: var(--radius-sm);
		background: var(--accent);
		color: var(--accent-ink);
		font-size: var(--text-base);
		font-weight: 600;
		cursor: pointer;
		transition: background-color 140ms var(--ease-out);
	}
	.primary:hover:not(:disabled) {
		background: var(--accent-strong);
	}
	.primary:disabled {
		opacity: 0.7;
		cursor: progress;
	}
	.error {
		margin: 0;
		padding: 0.65rem 0.85rem;
		background: var(--danger-soft, #fdecec);
		color: var(--danger, #b42318);
		border-radius: var(--radius-sm);
		font-size: var(--text-sm);
	}
</style>
