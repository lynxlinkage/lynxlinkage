<script lang="ts">
	import Button from '$lib/components/Button.svelte';
	import Hero from '$lib/components/Hero.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { site } from '$lib/site';

	const planes = [
		{ code: 'STS', label: 'Strategy sessions' },
		{ code: 'TD', label: 'Venues & orders' },
		{ code: 'MD', label: 'Public feeds & tape' },
		{ code: 'SYM', label: 'Tick / step' },
		{ code: 'Paper', label: 'Paper trading' },
		{ code: 'API+UI', label: 'Control plane' }
	];

	const venues = ['Paper', 'Gate', 'Binance spot', 'Binance USD-M', 'Bybit'];

	const plans = [
		{
			name: 'Self-host',
			body: 'Free. Run the node yourself from source.'
		},
		{
			name: 'Standard',
			body: 'Hosted instance with reservation — not self-serve checkout.'
		},
		{
			name: 'Custom',
			body: 'VPC or distributed setups, plus strategy help when you need it.'
		}
	];
</script>

<Seo
	title="Product"
	description="MFTIK is LynxLinkage's self-hosted trading desk (TradeOps): one node, one Owner, Python strategies, and the ops loop to deploy, watch, and recover."
	path="/product"
/>

<Hero
	eyebrow="Product · MFTIK"
	title="The trading desk we run — and now publish."
	subtitle="MFTIK is self-hosted TradeOps software: one node, one Owner, strategies in Python, and the process control to deploy, watch, stop, ack failures, and pull logs."
>
	{#snippet cta()}
		<Button
			href="https://mftik.com"
			variant="primary"
			size="lg"
			target="_blank"
			rel="noopener noreferrer">mftik.com</Button
		>
		<Button
			href="https://docs.mftik.com"
			variant="secondary"
			size="lg"
			target="_blank"
			rel="noopener noreferrer">Docs</Button
		>
		<Button href="/#contact" variant="ghost" size="lg">Contact</Button>
	{/snippet}
</Hero>

<section class="section">
	<div class="container-narrow">
		<h2>Not only develop — TradeOps included</h2>
		<p>
			{site.name} is the team: statistical arbitrage and market making in crypto, owning the stack end
			to end. <strong>MFTIK</strong> is the desk software we built so a lean team can develop and operate
			strategies on a single self-hosted node.
		</p>
		<p>
			You host a node. You write strategies in Python. The node keeps processes up &mdash; deploy,
			watch, stop, acknowledge a failure, pull logs. It is a trading desk, not a backtester.
			Single-tenant: one Owner per node.
		</p>
	</div>
</section>

<section class="section section--surface">
	<div class="container">
		<div class="section-heading">
			<span class="eyebrow">Planes</span>
			<h2>STS, TD, MD, and the rest of the desk.</h2>
			<p>
				Planes talk through Redis; they do not import each other. Full architecture lives in the
				docs.
			</p>
		</div>
		<ul class="planes">
			{#each planes as plane (plane.code)}
				<li>
					<span class="planes__code">{plane.code}</span>
					<span class="planes__label">{plane.label}</span>
				</li>
			{/each}
		</ul>
	</div>
</section>

<section class="section">
	<div class="container-narrow">
		<h2>Venues</h2>
		<p class="lede">Already in the product:</p>
		<ul class="venues">
			{#each venues as venue (venue)}
				<li>{venue}</li>
			{/each}
		</ul>
	</div>
</section>

<section class="section section--surface">
	<div class="container">
		<div class="section-heading">
			<span class="eyebrow">Plans</span>
			<h2>Self-host, Standard, or Custom.</h2>
			<p>
				Pricing and reservation details live on
				<a href="https://mftik.com" target="_blank" rel="noopener noreferrer">mftik.com</a>.
			</p>
		</div>
		<ul class="plans">
			{#each plans as plan (plan.name)}
				<li>
					<h3>{plan.name}</h3>
					<p>{plan.body}</p>
				</li>
			{/each}
		</ul>
	</div>
</section>

<section class="section">
	<div class="container-narrow">
		<h2>Go further</h2>
		<p>
			Product site, docs, and source are public. For a demo or hosted reservation, use the company
			contact form &mdash; pick the MFTIK topic.
		</p>
		<div class="cta-row">
			<Button href="https://mftik.com" variant="primary" target="_blank" rel="noopener noreferrer"
				>mftik.com</Button
			>
			<Button
				href="https://docs.mftik.com"
				variant="secondary"
				target="_blank"
				rel="noopener noreferrer">docs.mftik.com</Button
			>
			<Button
				href="https://github.com/lynxlinkage/mftik"
				variant="secondary"
				target="_blank"
				rel="noopener noreferrer">GitHub</Button
			>
			<Button href="/#contact" variant="ghost">Contact</Button>
		</div>
	</div>
</section>

<style>
	.planes {
		list-style: none;
		padding: 0;
		margin: 0;
		display: grid;
		gap: var(--space-4);
		grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
	}
	.planes li {
		padding: var(--space-5);
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		display: grid;
		gap: var(--space-2);
	}
	.planes__code {
		font-family: var(--font-mono, ui-monospace, monospace);
		font-size: var(--text-sm);
		font-weight: 600;
		letter-spacing: 0.04em;
		color: var(--fg);
	}
	.planes__label {
		font-size: var(--text-sm);
		color: var(--muted);
	}

	.lede {
		color: var(--muted);
		margin-bottom: var(--space-4);
	}

	.venues {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
	}
	.venues li {
		padding: 0.45rem 0.85rem;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--surface);
		font-size: var(--text-sm);
		color: var(--fg);
	}

	.plans {
		list-style: none;
		padding: 0;
		margin: 0;
		display: grid;
		gap: var(--space-5);
		grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
	}
	.plans li {
		padding: var(--space-5);
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius);
	}
	.plans h3 {
		font-size: var(--text-lg);
		margin-bottom: var(--space-2);
	}
	.plans p {
		margin: 0;
		font-size: var(--text-sm);
		color: var(--muted);
	}

	.cta-row {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
		margin-top: var(--space-6);
	}
</style>
