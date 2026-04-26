<script lang="ts">
	import type { PageData } from './$types';
	import Button from '$lib/components/Button.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { employmentLabel, formatDate } from '$lib/format';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
	const job = $derived(data.job);

	const applyHref = $derived(
		job.applyUrlOrEmail.includes('@')
			? `mailto:${job.applyUrlOrEmail}?subject=${encodeURIComponent(`Application: ${job.title}`)}`
			: job.applyUrlOrEmail
	);
	const isExternalLink = $derived(!job.applyUrlOrEmail.includes('@'));
</script>

<Seo
	title={job.title}
	description={`Open role at lynxlinkage: ${job.title} on the ${job.team} team, ${job.location}.`}
	path={`/hiring/${job.id}`}
/>

<article class="job">
	<header class="job__header">
		<div class="container-narrow">
			<a class="job__back" href="/hiring">
				<span aria-hidden="true">&larr;</span> All open roles
			</a>
			<div class="job__meta">
				<span class="job__team">{job.team || 'Team'}</span>
				<span class="job__sep" aria-hidden="true">&middot;</span>
				<span>{job.location}</span>
				<span class="job__sep" aria-hidden="true">&middot;</span>
				<span>{employmentLabel(job.employmentType)}</span>
				<span class="job__sep" aria-hidden="true">&middot;</span>
				<span>Posted {formatDate(job.postedAt)}</span>
			</div>
			<h1 class="job__title">{job.title}</h1>
			<div class="job__cta">
				<Button
					href={applyHref}
					variant="primary"
					size="lg"
					rel={isExternalLink ? 'noopener noreferrer' : undefined}
					target={isExternalLink ? '_blank' : undefined}
				>
					Apply now
				</Button>
				<Button href="/hiring" variant="ghost">Other roles</Button>
			</div>
		</div>
	</header>

	<section class="section">
		<div class="container-narrow">
			<!-- eslint-disable-next-line svelte/no-at-html-tags -->
			<div class="prose">{@html data.descriptionHtml}</div>

			<div class="job__apply-footer">
				<h2>Sounds like you?</h2>
				<p>
					Send a CV (or a link to your work) to
					<a href={`mailto:${job.applyUrlOrEmail.includes('@') ? job.applyUrlOrEmail : 'careers@lynxlinkage.com'}?subject=${encodeURIComponent(`Application: ${job.title}`)}`}>
						{job.applyUrlOrEmail.includes('@') ? job.applyUrlOrEmail : 'careers@lynxlinkage.com'}
					</a>. A short note about why this role catches your eye goes a long way.
				</p>
				<Button
					href={applyHref}
					variant="primary"
					rel={isExternalLink ? 'noopener noreferrer' : undefined}
					target={isExternalLink ? '_blank' : undefined}
				>
					Apply for {job.title}
				</Button>
			</div>
		</div>
	</section>
</article>

<style>
	.job__header {
		padding-block: var(--space-9) var(--space-7);
		background: var(--surface);
		border-bottom: 1px solid var(--border);
	}
	.job__back {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		color: var(--text-muted);
		font-size: var(--text-sm);
		text-decoration: none;
		margin-bottom: var(--space-5);
	}
	.job__back:hover {
		color: var(--accent);
	}
	.job__meta {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
		align-items: center;
		font-size: var(--text-sm);
		color: var(--text-muted);
		margin-bottom: var(--space-3);
	}
	.job__team {
		font-weight: 500;
		color: var(--accent);
	}
	.job__sep {
		opacity: 0.5;
	}
	.job__title {
		font-size: var(--text-3xl);
		margin: 0 0 var(--space-5);
		letter-spacing: -0.02em;
	}
	.job__cta {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
	}

	.prose :global(h1),
	.prose :global(h2),
	.prose :global(h3) {
		margin-top: var(--space-7);
		margin-bottom: var(--space-3);
	}
	.prose :global(h2) {
		font-size: var(--text-xl);
	}
	.prose :global(h3) {
		font-size: var(--text-lg);
	}
	.prose :global(p) {
		margin-bottom: var(--space-4);
		line-height: 1.7;
	}
	.prose :global(ul),
	.prose :global(ol) {
		margin: 0 0 var(--space-4) var(--space-5);
		padding: 0;
	}
	.prose :global(li) {
		margin-bottom: var(--space-2);
		line-height: 1.7;
	}
	.prose :global(strong) {
		color: var(--text);
	}

	.job__apply-footer {
		margin-top: var(--space-8);
		padding: var(--space-6);
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius);
	}
	.job__apply-footer h2 {
		font-size: var(--text-xl);
		margin-bottom: var(--space-3);
	}
	.job__apply-footer p {
		margin-bottom: var(--space-4);
		color: var(--text-muted);
	}
</style>
