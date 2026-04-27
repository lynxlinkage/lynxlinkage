<script lang="ts">
	import { site } from '$lib/site';

	const year = new Date().getFullYear();
</script>

<footer class="footer">
	<div class="container footer__inner">
		<div class="footer__brand">
			<div class="footer__name">{site.name}</div>
			<p class="footer__tag">{site.tagline}</p>
			<p class="footer__locations">{site.locations.join(' · ')}</p>
		</div>

		<nav aria-label="Footer">
			<div class="footer__col">
				<h4>Company</h4>
				<ul>
					{#each site.nav as item (item.href)}
						<li><a href={item.href}>{item.label}</a></li>
					{/each}
				</ul>
			</div>
		</nav>
	</div>

	<div class="container footer__bottom">
		<span>&copy; {year} {site.name}. All rights reserved.</span>
		<span>Trading involves risk. Past performance is not indicative of future results.</span>
	</div>
</footer>

<style>
	.footer {
		background: var(--surface);
		border-top: 1px solid var(--border);
		padding-block: var(--space-10) var(--space-6);
		color: var(--text-muted);
	}

	.footer__inner {
		display: grid;
		grid-template-columns: 1.4fr 1fr;
		gap: var(--space-8);
		margin-bottom: var(--space-8);
	}

	.footer__brand {
		max-width: 32ch;
	}
	.footer__name {
		font-weight: 700;
		font-size: var(--text-lg);
		color: var(--text);
		margin-bottom: var(--space-3);
		letter-spacing: -0.015em;
	}
	.footer__tag {
		margin-bottom: var(--space-3);
		color: var(--text-muted);
	}
	.footer__locations {
		margin: 0;
		font-size: var(--text-sm);
		color: var(--text-faint);
	}

	.footer__col h4 {
		font-size: var(--text-xs);
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: var(--text-faint);
		margin: 0 0 var(--space-3);
		font-weight: 600;
	}
	.footer__col ul {
		list-style: none;
		margin: 0;
		padding: 0;
		display: grid;
		gap: var(--space-2);
	}
	.footer__col a {
		color: var(--text);
		text-decoration: none;
		font-size: var(--text-sm);
	}
	.footer__col a:hover {
		color: var(--accent);
	}

	.footer__bottom {
		display: flex;
		justify-content: space-between;
		gap: var(--space-3);
		flex-wrap: wrap;
		font-size: var(--text-xs);
		color: var(--text-faint);
		padding-top: var(--space-5);
		border-top: 1px solid var(--border);
	}

	@media (max-width: 720px) {
		.footer__inner {
			grid-template-columns: 1fr;
		}
	}
</style>
