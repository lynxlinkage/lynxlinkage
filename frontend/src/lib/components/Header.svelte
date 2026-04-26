<script lang="ts">
	import { page } from '$app/state';
	import { site } from '$lib/site';

	let scrolled = $state(false);
	let mobileOpen = $state(false);

	function onScroll() {
		scrolled = window.scrollY > 8;
	}

	function toggleMobile() {
		mobileOpen = !mobileOpen;
	}

	function closeMobile() {
		mobileOpen = false;
	}

	$effect(() => {
		onScroll();
		window.addEventListener('scroll', onScroll, { passive: true });
		return () => window.removeEventListener('scroll', onScroll);
	});

	const currentPath = $derived(page.url?.pathname ?? '/');
</script>

<header class="header" class:header--scrolled={scrolled}>
	<div class="container header__inner">
		<a class="header__brand" href="/" onclick={closeMobile}>
			<svg width="28" height="28" viewBox="0 0 32 32" aria-hidden="true">
				<rect width="32" height="32" rx="6" fill="currentColor" />
				<path d="M9 22V10h2v10h6v2z" fill="#fff" />
				<path d="M19 10h3l3 12h-2l-2.5-8.5L18 22h-2z" fill="#fff" opacity="0.85" />
			</svg>
			<span class="header__brand-text">{site.name}</span>
		</a>

		<nav class="header__nav" aria-label="Primary">
			<ul>
				{#each site.nav as item (item.href)}
					<li>
						<a
							href={item.href}
							class:active={currentPath === item.href || currentPath.startsWith(`${item.href}/`)}
						>
							{item.label}
						</a>
					</li>
				{/each}
			</ul>
		</nav>

		<div class="header__cta">
			<a class="header__contact" href="/#contact">Contact</a>
		</div>

		<button
			class="header__burger"
			class:header__burger--open={mobileOpen}
			onclick={toggleMobile}
			aria-label="Toggle navigation menu"
			aria-expanded={mobileOpen}
			aria-controls="mobile-nav"
		>
			<span></span>
			<span></span>
			<span></span>
		</button>
	</div>

	{#if mobileOpen}
		<div class="header__drawer" id="mobile-nav">
			<nav aria-label="Primary mobile">
				<ul>
					{#each site.nav as item (item.href)}
						<li>
							<a href={item.href} onclick={closeMobile}>{item.label}</a>
						</li>
					{/each}
					<li>
						<a class="header__drawer-cta" href="/#contact" onclick={closeMobile}>Contact</a>
					</li>
				</ul>
			</nav>
		</div>
	{/if}
</header>

<style>
	.header {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 50;
		height: var(--header-height);
		background: rgba(255, 255, 255, 0.7);
		backdrop-filter: saturate(140%) blur(10px);
		-webkit-backdrop-filter: saturate(140%) blur(10px);
		border-bottom: 1px solid transparent;
		transition:
			background-color 160ms var(--ease-out),
			border-color 160ms var(--ease-out),
			box-shadow 160ms var(--ease-out);
	}
	.header--scrolled {
		background: rgba(255, 255, 255, 0.92);
		border-bottom-color: var(--border);
		box-shadow: var(--shadow-sm);
	}

	.header__inner {
		display: flex;
		align-items: center;
		gap: var(--space-5);
		height: 100%;
	}

	.header__brand {
		display: inline-flex;
		align-items: center;
		gap: 0.55rem;
		color: var(--accent);
		text-decoration: none;
		font-weight: 700;
		letter-spacing: -0.02em;
	}
	.header__brand:hover {
		color: var(--accent-strong);
	}
	.header__brand-text {
		font-size: 1.05rem;
		color: var(--text);
	}

	.header__nav {
		flex: 1;
		display: flex;
		justify-content: center;
	}
	.header__nav ul {
		display: flex;
		gap: var(--space-5);
		list-style: none;
		margin: 0;
		padding: 0;
	}
	.header__nav a {
		color: var(--text-muted);
		font-size: var(--text-sm);
		font-weight: 500;
		text-decoration: none;
		padding: 0.4rem 0;
		position: relative;
	}
	.header__nav a:hover,
	.header__nav a.active {
		color: var(--text);
	}
	.header__nav a.active::after {
		content: '';
		position: absolute;
		left: 0;
		right: 0;
		bottom: -2px;
		height: 2px;
		background: var(--accent);
		border-radius: 1px;
	}

	.header__cta {
		display: flex;
		align-items: center;
	}
	.header__contact {
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--accent);
		padding: 0.5rem 1rem;
		border: 1px solid var(--accent-soft);
		border-radius: var(--radius);
		background: var(--accent-soft);
		text-decoration: none;
		transition:
			background-color 140ms var(--ease-out),
			color 140ms var(--ease-out);
	}
	.header__contact:hover {
		background: var(--accent);
		color: var(--accent-ink);
	}

	.header__burger {
		display: none;
		flex-direction: column;
		justify-content: center;
		gap: 4px;
		width: 40px;
		height: 40px;
		padding: 0;
		background: transparent;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		cursor: pointer;
	}
	.header__burger span {
		display: block;
		width: 18px;
		height: 1.6px;
		background: var(--text);
		margin: 0 auto;
		transition: transform 160ms var(--ease-out);
	}
	.header__burger--open span:nth-child(1) {
		transform: translateY(5.6px) rotate(45deg);
	}
	.header__burger--open span:nth-child(2) {
		opacity: 0;
	}
	.header__burger--open span:nth-child(3) {
		transform: translateY(-5.6px) rotate(-45deg);
	}

	.header__drawer {
		position: absolute;
		top: var(--header-height);
		left: 0;
		right: 0;
		background: var(--bg);
		border-bottom: 1px solid var(--border);
		box-shadow: var(--shadow-md);
	}
	.header__drawer ul {
		list-style: none;
		margin: 0;
		padding: var(--space-3) clamp(var(--space-4), 4vw, var(--space-7));
		display: grid;
		gap: var(--space-1);
	}
	.header__drawer a {
		display: block;
		padding: 0.85rem 0;
		font-size: var(--text-base);
		color: var(--text);
		text-decoration: none;
		border-bottom: 1px solid var(--border);
	}
	.header__drawer a:hover {
		color: var(--accent);
	}
	.header__drawer-cta {
		color: var(--accent) !important;
		font-weight: 500;
	}

	@media (max-width: 860px) {
		.header__nav,
		.header__cta {
			display: none;
		}
		.header__burger {
			display: inline-flex;
			margin-left: auto;
		}
	}
</style>
