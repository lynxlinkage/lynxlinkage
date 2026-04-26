<script lang="ts">
	import type { PageData } from './$types';
	import Hero from '$lib/components/Hero.svelte';
	import ResearchCard from '$lib/components/ResearchCard.svelte';
	import Seo from '$lib/components/Seo.svelte';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const allTags = $derived(
		Array.from(new Set(data.researches.flatMap((r) => r.tags))).sort()
	);

	let activeTag = $state<string | null>(null);

	const filtered = $derived(
		activeTag
			? data.researches.filter((r) => r.tags.includes(activeTag as string))
			: data.researches
	);

	function selectTag(tag: string | null) {
		activeTag = tag;
	}
</script>

<Seo
	title="Public researches"
	description="Selected research notes from the lynxlinkage research desk on market microstructure, volatility, and infrastructure."
	path="/researches"
/>

<Hero
	eyebrow="Public researches"
	title="Notes from the research desk."
	subtitle="A curated stream of writing on market microstructure, volatility, and the infrastructure that connects them."
/>

<section class="section">
	<div class="container">
		{#if allTags.length}
			<div class="tags" role="toolbar" aria-label="Filter by tag">
				<button
					class="tag"
					class:tag--active={activeTag === null}
					type="button"
					onclick={() => selectTag(null)}
				>
					All
				</button>
				{#each allTags as tag (tag)}
					<button
						class="tag"
						class:tag--active={activeTag === tag}
						type="button"
						onclick={() => selectTag(tag)}
					>
						{tag}
					</button>
				{/each}
			</div>
		{/if}

		{#if filtered.length}
			<div class="grid">
				{#each filtered as card (card.id)}
					<ResearchCard {card} />
				{/each}
			</div>
		{:else}
			<p class="muted">No research matches that filter yet.</p>
		{/if}
	</div>
</section>

<style>
	.tags {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-2);
		margin-bottom: var(--space-6);
	}
	.tag {
		font-family: inherit;
		font-size: var(--text-xs);
		padding: 0.4rem 0.85rem;
		background: var(--bg);
		color: var(--text-muted);
		border: 1px solid var(--border);
		border-radius: 999px;
		cursor: pointer;
		transition:
			border-color 140ms var(--ease-out),
			color 140ms var(--ease-out),
			background-color 140ms var(--ease-out);
	}
	.tag:hover {
		border-color: var(--border-strong);
		color: var(--text);
	}
	.tag--active {
		background: var(--accent);
		color: var(--accent-ink);
		border-color: var(--accent);
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: var(--space-5);
	}
</style>
