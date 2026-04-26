<script lang="ts">
	import type { PageData } from './$types';
	import Hero from '$lib/components/Hero.svelte';
	import Button from '$lib/components/Button.svelte';
	import ResearchCard from '$lib/components/ResearchCard.svelte';
	import JobCard from '$lib/components/JobCard.svelte';
	import LogoWall from '$lib/components/LogoWall.svelte';
	import ContactForm from '$lib/components/ContactForm.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { site } from '$lib/site';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
</script>

<Seo path="/" />

<Hero
	eyebrow="Trading firm"
	title="Statistical arbitrage in crypto markets."
	subtitle="We are a agile team of researchers and engineers building systematic stat-arb strategies for digital-asset markets &mdash; cross-venue, basis, funding, and cointegrated pairs. We publish what we learn, hire technically deep people, and partner with the venues we depend on."
>
	{#snippet cta()}
		<Button href="/researches" variant="primary" size="lg">Read our research</Button>
		<Button href="/hiring" variant="secondary" size="lg">See open roles</Button>
	{/snippet}
</Hero>

<section class="section section--surface">
	<div class="container">
		<div class="section-heading">
			<span class="eyebrow">About us</span>
			<h2>A agile team going deep on crypto statistical arbitrage.</h2>
			<p>
				{site.name} is a research-led trading firm. We model price discovery in fragmented crypto
				markets &mdash; across venues, between spot and derivatives, and across cointegrated
				baskets &mdash; and we ship fast because the edges decay fast.
			</p>
		</div>
		<Button href="/about" variant="ghost">More about us &rarr;</Button>
	</div>
</section>

<section class="section">
	<div class="container">
		<div class="section-heading">
			<span class="eyebrow">Public researches</span>
			<h2>Selected pieces from our research desk.</h2>
			<p>We publish what we can &mdash; sometimes on Medium, sometimes on our own site.</p>
		</div>
		{#if data.researches.length}
			<div class="grid">
				{#each data.researches as card (card.id)}
					<ResearchCard {card} />
				{/each}
			</div>
		{:else}
			<p class="muted">Research notes will appear here soon.</p>
		{/if}
		<div class="section__cta">
			<Button href="/researches" variant="ghost">Browse all research &rarr;</Button>
		</div>
	</div>
</section>

<section class="section section--surface">
	<div class="container">
		<div class="section-heading">
			<span class="eyebrow">Hiring</span>
			<h2>Build the systems that trade markets.</h2>
			<p>We hire deeply technical people who care about the craft.</p>
		</div>
		{#if data.jobs.length}
			<div class="grid grid--two">
				{#each data.jobs as job (job.id)}
					<JobCard {job} />
				{/each}
			</div>
		{:else}
			<p class="muted">No open roles right now &mdash; check back soon.</p>
		{/if}
		<div class="section__cta">
			<Button href="/hiring" variant="ghost">All open roles &rarr;</Button>
		</div>
	</div>
</section>

<section class="section">
	<div class="container">
		<div class="section-heading">
			<span class="eyebrow">Partners</span>
			<h2>We work with leading venues and infrastructure providers.</h2>
		</div>
		{#if data.partners.length}
			<LogoWall partners={data.partners} compact />
		{:else}
			<p class="muted">Coming soon.</p>
		{/if}
		<div class="section__cta">
			<Button href="/partners" variant="ghost">Meet our partners &rarr;</Button>
		</div>
	</div>
</section>

<section class="section section--surface" id="contact">
	<div class="container-narrow">
		<ContactForm
			subtitle="Whether it's a research collaboration, a partnership opportunity, or a question about the firm &mdash; drop us a line."
		/>
	</div>
</section>

<style>
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: var(--space-5);
		margin-bottom: var(--space-6);
	}
	.grid--two {
		grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
	}

	.section__cta {
		display: flex;
		justify-content: flex-start;
		margin-top: var(--space-5);
	}
</style>
