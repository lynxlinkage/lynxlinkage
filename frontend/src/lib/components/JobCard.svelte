<script lang="ts">
	import type { JobPosting } from '$lib/api/types';
	import { employmentLabel } from '$lib/format';

	interface Props {
		job: JobPosting;
	}

	let { job }: Props = $props();

	const applyHref = $derived(
		job.applyUrlOrEmail.includes('@')
			? `mailto:${job.applyUrlOrEmail}?subject=${encodeURIComponent(`Application: ${job.title}`)}`
			: job.applyUrlOrEmail
	);

	const isExternalLink = $derived(!job.applyUrlOrEmail.includes('@'));
</script>

<article class="card">
	<div class="card__meta">
		<span class="card__team">{job.team || 'Team'}</span>
		<span class="card__sep" aria-hidden="true">&middot;</span>
		<span>{job.location}</span>
		<span class="card__sep" aria-hidden="true">&middot;</span>
		<span>{employmentLabel(job.employmentType)}</span>
	</div>
	<h3 class="card__title">{job.title}</h3>
	<a
		class="card__cta"
		href={applyHref}
		rel={isExternalLink ? 'noopener noreferrer' : undefined}
		target={isExternalLink ? '_blank' : undefined}
	>
		Apply now
		<span aria-hidden="true">&rarr;</span>
	</a>
</article>

<style>
	.card {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
		padding: var(--space-6);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg);
		transition:
			border-color 140ms var(--ease-out),
			box-shadow 140ms var(--ease-out);
	}
	.card:hover {
		border-color: var(--border-strong);
		box-shadow: var(--shadow-md);
	}

	.card__meta {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
		align-items: center;
		font-size: var(--text-sm);
		color: var(--text-muted);
	}
	.card__team {
		font-weight: 500;
		color: var(--accent);
	}
	.card__sep {
		opacity: 0.5;
	}

	.card__title {
		font-size: var(--text-xl);
		margin: 0;
	}

	.card__cta {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		align-self: flex-start;
		font-weight: 500;
		font-size: var(--text-sm);
		color: var(--accent);
		margin-top: var(--space-2);
	}
	.card__cta:hover span {
		transform: translateX(3px);
	}
	.card__cta span {
		display: inline-block;
		transition: transform 160ms var(--ease-out);
	}
</style>
