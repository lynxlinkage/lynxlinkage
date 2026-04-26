<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';

	let { children } = $props();
</script>

<a class="skip-link" href="#main">Skip to content</a>

<Header />

<main id="main">
	{@render children()}
</main>

<Footer />

<style>
	main {
		min-height: calc(100dvh - var(--header-height));
		padding-top: var(--header-height);
	}

	.skip-link {
		position: absolute;
		top: -40px;
		left: var(--space-3);
		padding: 0.6rem 1rem;
		background: var(--accent);
		color: var(--accent-ink);
		border-radius: var(--radius-sm);
		font-size: var(--text-sm);
		text-decoration: none;
		z-index: 100;
	}
	.skip-link:focus {
		top: var(--space-2);
	}
</style>
