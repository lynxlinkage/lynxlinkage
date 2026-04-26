<script lang="ts">
	import Hero from '$lib/components/Hero.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { site } from '$lib/site';

	const focusAreas = [
		{
			title: 'Cross-venue mispricings',
			body: 'The same instrument trades on dozens of venues with different microstructures, fee schedules, and risk policies. We model the gap and trade it.'
		},
		{
			title: 'Basis & funding',
			body: 'Spot, perpetuals, and dated futures are tied together by funding and basis. When the rope goes slack, we know where to pull.'
		},
		{
			title: 'Cointegration & pairs',
			body: 'Long-tail tokens move in clusters. We build cointegrated baskets and pair trades that mean-revert with measurable half-lives.'
		},
		{
			title: 'Microstructure-aware execution',
			body: 'Order-book imbalance, taker/maker dynamics, and latency budgets decide whether a signal becomes PnL or slippage. We treat execution as a first-class research problem.'
		}
	];

	const principles = [
		{
			title: 'Research first, code second',
			body: 'Every trading line of code starts as a question. We treat backtests as evidence, not marketing material, and we are quick to throw out our own ideas.'
		},
		{
			title: 'Own the stack',
			body: 'From market-data capture to risk to execution \u2014 we build it ourselves when off-the-shelf tools fall short of the bar that crypto microstructure demands.'
		},
		{
			title: 'Theory-driven, not data-driven',
			body: 'Each strategy starts as a statistical conjecture about how a market should behave. The data is where we look for disconfirming evidence, and the live market is the final referee.'
		},
		{
			title: 'Small team, full ownership',
			body: 'Researchers and engineers ship together. There is no handoff, no committee, and no headcount we have to feed \u2014 the team stays small on purpose.'
		},
		{
			title: 'Share what we learn',
			body: 'Crypto research benefits from being in public. We publish on Medium and our own platform when we can \u2014 it keeps us honest and the conversations interesting.'
		}
	];

	const stats = [
		{ label: 'Founded', value: String(site.founded) },
		{ label: 'Markets traded', value: 'Crypto only' },
		{ label: 'Strategy style', value: 'Statistical arbitrage' },
		{ label: 'Operating hours', value: '24 / 7' }
	];
</script>

<Seo
	title="About us"
	description="lynxlinkage is an agile team of researchers and engineers focused on statistical arbitrage in crypto markets."
	path="/about"
/>

<Hero
	eyebrow="About us"
	title="An agile team going deep on statistical arbitrage in crypto market."
	subtitle="We trade systematic strategies across digital-asset markets &mdash; cross-venue, basis, funding, and cointegrated pairs. No retail fluff, no AI buzzwords. Careful research, robust infrastructure, and tight feedback loops between the two."
/>

<section class="section">
	<div class="container-narrow">
		<h2>What we focus on</h2>
		<p>
			{site.name} is built around one idea: <strong
				>statistical arbitrage in crypto markets is still a research problem, not a solved one.</strong
			>
			Crypto suits this style of trading. Markets run 24 / 7, public order books are deep, and the
			structural inefficiencies are still strong enough to reward serious research &mdash; provided
			you do the work to capture clean data, model microstructure honestly, and execute with
			discipline.
		</p>
		<p>
			We are deliberately small. Every person on the team writes code, looks at data, and owns
			strategies end-to-end. We do not run customer money, we do not chase narrative trades, and we
			do not over-fit slide decks &mdash; we trade our own book and let the equity curve do the
			talking.
		</p>
	</div>
</section>

<section class="section section--surface">
	<div class="container">
		<div class="section-heading">
			<span class="eyebrow">Where we play</span>
			<h2>Four families of statistical-arbitrage strategies.</h2>
			<p>
				Each strategy family has its own data, its own models, and its own execution profile. We
				rotate research effort between them based on capacity, regime, and edge.
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
			are happy to work on Taiwan time.
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
