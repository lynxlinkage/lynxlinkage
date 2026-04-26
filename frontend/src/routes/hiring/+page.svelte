<script lang="ts">
	import type { PageData } from './$types';
	import Hero from '$lib/components/Hero.svelte';
	import JobCard from '$lib/components/JobCard.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { site } from '$lib/site';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const grouped = $derived.by(() => {
		const buckets: Record<string, typeof data.jobs> = {};
		for (const j of data.jobs) {
			const key = j.team || 'Other';
			(buckets[key] ??= []).push(j);
		}
		return Object.entries(buckets).sort(([a], [b]) => a.localeCompare(b));
	});
</script>

<Seo
	title="Hiring"
	description="Open roles at lynxlinkage in research, engineering, and operations."
	path="/hiring"
/>

<Hero
	eyebrow="Hiring"
	title="Build the systems that trade markets."
	subtitle="We hire deeply technical people who care about the craft. Researchers and engineers work side-by-side in small teams with full ownership."
/>

<section class="section">
	<div class="container">
		{#if data.jobs.length}
			{#each grouped as [team, jobs] (team)}
				<div class="team">
					<h2 class="team__title">{team}</h2>
					<div class="grid">
						{#each jobs as job (job.id)}
							<JobCard {job} />
						{/each}
					</div>
				</div>
			{/each}
		{:else}
			<div class="empty">
				<h2>No open roles right now.</h2>
				<p>
					We&rsquo;re always interested in talking to exceptional researchers and engineers. Drop us
					a line at
					<a href={`mailto:${site.careersEmail}`}>{site.careersEmail}</a>.
				</p>
			</div>
		{/if}
	</div>
</section>

<style>
	.team {
		margin-bottom: var(--space-10);
	}
	.team__title {
		font-size: var(--text-xl);
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: var(--text-faint);
		margin-bottom: var(--space-5);
		padding-bottom: var(--space-3);
		border-bottom: 1px solid var(--border);
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
		gap: var(--space-4);
	}

	.empty {
		max-width: 56ch;
	}
</style>
