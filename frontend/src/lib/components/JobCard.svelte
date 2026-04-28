<script lang="ts">
	import type { JobPosting } from '$lib/api/types';
	import { employmentLabel } from '$lib/format';

	interface Props {
		job: JobPosting;
	}

	let { job }: Props = $props();
</script>

<a class="card" href={`/hiring/${job.id}`} aria-label={`View role: ${job.title}`}>
	<div class="card__meta">
		<span class="card__team">{job.team || 'Team'}</span>
		<span class="card__sep" aria-hidden="true">&middot;</span>
		<span>{job.location}</span>
		<span class="card__sep" aria-hidden="true">&middot;</span>
		<span>{employmentLabel(job.employmentType)}</span>
	</div>
	<h3 class="card__title">{job.title}</h3>
	<span class="card__cta">
		View role
		<span class="card__arrow" aria-hidden="true">&rarr;</span>
	</span>
</a>

<style>
	.card {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
		padding: var(--space-6);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg);
		text-decoration: none;
		color: inherit;
		transition:
			border-color 140ms var(--ease-out),
			box-shadow 140ms var(--ease-out),
			transform 140ms var(--ease-out);
	}
	.card:hover {
		border-color: var(--border-strong);
		box-shadow: var(--shadow-md);
		transform: translateY(-1px);
	}
	.card:focus-visible {
		outline: 2px solid var(--accent);
		outline-offset: 2px;
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
	.card__arrow {
		display: inline-block;
		transition: transform 160ms var(--ease-out);
	}
	.card:hover .card__arrow {
		transform: translateX(3px);
	}
</style>
