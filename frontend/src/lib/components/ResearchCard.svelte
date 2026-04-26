<script lang="ts">
	import type { ResearchCard } from '$lib/api/types';
	import { formatDate, sourceLabel } from '$lib/format';

	interface Props {
		card: ResearchCard;
	}

	let { card }: Props = $props();
</script>

<article class="card">
	<header class="card__head">
		<time class="card__date" datetime={card.publishedAt}>{formatDate(card.publishedAt)}</time>
		{#if card.tags.length > 0}
			<ul class="card__tags" aria-label="Tags">
				{#each card.tags as tag (tag)}
					<li class="card__tag">{tag}</li>
				{/each}
			</ul>
		{/if}
	</header>

	<h3 class="card__title">{card.title}</h3>
	<p class="card__summary">{card.summary}</p>

	<a class="card__link" href={card.externalUrl} rel="noopener noreferrer" target="_blank">
		{sourceLabel(card.source)}
		<span class="card__arrow" aria-hidden="true">&rarr;</span>
	</a>
</article>

<style>
	.card {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: var(--space-6);
		transition:
			border-color 140ms var(--ease-out),
			box-shadow 140ms var(--ease-out);
		height: 100%;
	}
	.card:hover {
		border-color: var(--border-strong);
		box-shadow: var(--shadow-md);
	}

	.card__head {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-2) var(--space-3);
		align-items: center;
		justify-content: space-between;
		font-size: var(--text-xs);
		color: var(--text-faint);
	}

	.card__date {
		letter-spacing: 0.04em;
		text-transform: uppercase;
	}

	.card__tags {
		list-style: none;
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-2);
		margin: 0;
		padding: 0;
	}

	.card__tag {
		font-size: var(--text-xs);
		padding: 0.15rem 0.55rem;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 999px;
		color: var(--text-muted);
	}

	.card__title {
		font-size: var(--text-xl);
		margin: var(--space-2) 0 0;
	}

	.card__summary {
		flex: 1;
		color: var(--text-muted);
		margin: 0;
	}

	.card__link {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		font-weight: 500;
		font-size: var(--text-sm);
		color: var(--accent);
	}
	.card__link:hover .card__arrow {
		transform: translateX(3px);
	}
	.card__arrow {
		display: inline-block;
		transition: transform 160ms var(--ease-out);
	}
</style>
