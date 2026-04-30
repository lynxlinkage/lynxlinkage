<script lang="ts">
	import Hero from '$lib/components/Hero.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { site } from '$lib/site';

	const focusAreas = [
		{
			title: 'Cross-venue statistical arbitrage',
			body: 'The same instrument trades on dozens of venues with different microstructures, fee schedules, and latency profiles. We model the spread as a stationary process, test for mean-reversion with rigorous statistical criteria, and size positions from the full residual distribution — not a point estimate.'
		},
		{
			title: 'Basis & funding arbitrage',
			body: 'Spot, perpetuals, and dated futures are linked by hard arbitrage constraints. When funding pressure or basis dislocations break those constraints, there is a clear causal mechanism behind the trade — not just a correlation. We model the reversion speed and hold time probabilistically.'
		},
		{
			title: 'Cointegrated pairs & baskets',
			body: 'Related tokens move in statistically measurable clusters. We apply cointegration tests and error-correction models to isolate pairs and baskets where mean-reversion has a structural cause, a measurable half-life, and a residual distribution that is stable enough to trade.'
		},
		{
			title: 'Market making',
			body: 'We provide liquidity across crypto spot and derivatives venues. Our quoting models are built on arrival-rate distributions, adverse-selection metrics, and real-time inventory risk — not fixed spreads. Execution is a first-class research problem, not an afterthought.'
		}
	];

	const principles = [
		{
			title: 'Statistics, not stories',
			body: 'Every edge is a statistical hypothesis with a null. We write the null before we run the test, report the full distribution of outcomes across regimes, and only go live when the evidence clears a high bar. A great Sharpe on one window is not evidence.'
		},
		{
			title: 'Causality over correlation',
			body: 'A correlated pair is a starting point, never a trade. We require a structural reason — funding mechanics, index pressure, cross-venue constraints — before committing research time. If we cannot write a causal story for why the spread should mean-revert, we do not build the model.'
		},
		{
			title: 'Theory precedes data',
			body: 'Each strategy begins as a written conjecture about how a market must behave given its microstructure. The backtest is the disconfirmation step, not the research itself. If the theory breaks in the data, the theory was wrong, and we discard it without ceremony.'
		},
		{
			title: 'Own the stack',
			body: 'From tick-data capture to risk management to execution — we build what we need. Off-the-shelf tools optimised for equities miss the microstructure details that determine whether a crypto strategy actually works at scale.'
		},
		{
			title: 'Lean by design',
			body: 'Researchers and engineers are the same people. There are no hand-offs, no approval layers, no strategy decks that go nowhere. Everyone owns a strategy end-to-end — from conjecture to live risk. The team stays small on purpose because that is how full ownership works.'
		}
	];

	const stats = [
		{ label: 'Founded', value: String(site.founded) },
		{ label: 'Markets traded', value: 'Crypto only' },
		{ label: 'Strategy style', value: 'Stat arb & MM' },
		{ label: 'Operating hours', value: '24 / 7' }
	];
</script>

<Seo
	title="About us"
	description="LynxLinkage is a lean team of researchers and engineers building statistical arbitrage and market-making strategies in crypto, grounded in probability theory and causal inference."
	path="/about"
/>

<Hero
	eyebrow="About LynxLinkage"
	title="Statistical arbitrage and market making, grounded in probability and causality."
	subtitle="We are a lean team of researchers and engineers who go all the way down: causal models, owned infrastructure, and strategies that only trade when the statistics say to."
/>

<section class="section">
	<div class="container-narrow">
		<h2>What we focus on</h2>
		<p>
			{site.name} is built around one conviction: <strong>crypto's structural inefficiencies are
				still large enough to reward deep research — but not large enough to forgive sloppy
				thinking.</strong>
			We work in two disciplines where probability theory and causal inference determine whether you
			make money or give it back: statistical arbitrage and market making. Every strategy starts as
			a falsifiable claim about how a market behaves, and every claim must survive statistical
			scrutiny before it sees real capital.
		</p>
		<p>
			We are deliberately lean. Every person on the team owns research, code, and live risk
			simultaneously. There are no hand-offs, no approval layers, and no organisational distance
			between the person who writes the model and the person who answers for the PnL &mdash;
			because they are the same person.
		</p>
	</div>
</section>

<section class="section section--surface">
	<div class="container">
		<div class="section-heading">
			<span class="eyebrow">Where we trade</span>
			<h2>Statistical arbitrage and market making across crypto.</h2>
			<p>
				Each strategy family has its own data requirements, its own causal model, and its own
				execution profile. We allocate research effort based on capacity, regime, and edge — and we
				only run strategies that have a theory we can defend.
			</p>
		</div>
		<ul class="focus">
			{#each focusAreas as area (area.title)}
				<li>
					<h3>{area.title}</h3>
					<p>{area.body}</p>
				</li>
			{/each}
		</ul>
	</div>
</section>

<section class="section">
	<div class="container-narrow">
		<h2>How we work</h2>
		<p class="principles-intro">
			Five principles that are non-negotiable regardless of market conditions, team size, or how
			good the last backtest looked.
		</p>
		<ul class="principles">
			{#each principles as p (p.title)}
				<li>
					<h3>{p.title}</h3>
					<p>{p.body}</p>
				</li>
			{/each}
		</ul>
	</div>
</section>

<section class="section section--surface">
	<div class="container">
		<dl class="stats">
			{#each stats as s (s.label)}
				<div class="stats__item">
					<dt>{s.label}</dt>
					<dd>{s.value}</dd>
				</div>
			{/each}
		</dl>
	</div>
</section>

<section class="section">
	<div class="container-narrow">
		<h2>Where we are</h2>
		<p>
			The team is fully remote and based in Taiwan. We organise around overlapping working hours
			rather than a single office, and we hire wherever the right person is &mdash; provided they
			are comfortable working on Taiwan time and are happy with a flat team where your ideas go
			straight into production without three levels of sign-off.
		</p>
	</div>
</section>

<style>
	.focus {
		list-style: none;
		padding: 0;
		margin: 0;
		display: grid;
		gap: var(--space-5);
		grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
	}
	.focus li {
		padding: var(--space-5);
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius);
	}
	.focus h3 {
		font-size: var(--text-lg);
		margin-bottom: var(--space-2);
	}
	.focus p {
		margin: 0;
		font-size: var(--text-sm);
		color: var(--muted);
	}

	.principles-intro {
		color: var(--muted);
		margin-bottom: var(--space-6);
	}

	.principles {
		list-style: none;
		padding: 0;
		margin: 0;
		display: grid;
		gap: var(--space-5);
		grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
	}
	.principles li {
		padding: var(--space-5);
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius);
	}
	.principles h3 {
		font-size: var(--text-lg);
		margin-bottom: var(--space-2);
	}
	.principles p {
		margin: 0;
		font-size: var(--text-sm);
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: var(--space-5);
		margin: 0;
	}
	.stats__item {
		padding: var(--space-5);
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius);
	}
	.stats__item dt {
		font-size: var(--text-xs);
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: var(--muted);
		margin-bottom: var(--space-2);
	}
	.stats__item dd {
		margin: 0;
		font-size: var(--text-xl);
		font-weight: 600;
		color: var(--fg);
	}
</style>
