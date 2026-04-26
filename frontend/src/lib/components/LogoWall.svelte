<script lang="ts">
	import type { Partner } from '$lib/api/types';

	interface Props {
		partners: Partner[];
		compact?: boolean;
	}

	let { partners, compact = false }: Props = $props();
</script>

<ul class="wall" class:wall--compact={compact} aria-label="Partner logos">
	{#each partners as partner (partner.id)}
		<li class="wall__item">
			{#if partner.websiteUrl}
				<a
					class="wall__link"
					href={partner.websiteUrl}
					rel="noopener noreferrer"
					target="_blank"
					title={partner.name}
				>
					<img class="wall__logo" src={partner.logoUrl} alt={partner.name} loading="lazy" />
				</a>
			{:else}
				<img
					class="wall__logo"
					src={partner.logoUrl}
					alt={partner.name}
					title={partner.name}
					loading="lazy"
				/>
			{/if}
		</li>
	{/each}
</ul>

<style>
	.wall {
		list-style: none;
		margin: 0;
		padding: 0;
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
		gap: var(--space-5);
		align-items: stretch;
	}
	.wall--compact {
		grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
		gap: var(--space-4);
	}

	.wall__item {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 88px;
		padding: var(--space-4);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg);
		transition: border-color 140ms var(--ease-out);
	}
	.wall__item:hover {
		border-color: var(--border-strong);
	}

	.wall__link {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
	}

	.wall__logo {
		max-height: 44px;
		width: auto;
		max-width: 80%;
		filter: grayscale(100%);
		opacity: 0.78;
		transition:
			filter 160ms var(--ease-out),
			opacity 160ms var(--ease-out);
	}
	.wall__item:hover .wall__logo {
		filter: grayscale(0%);
		opacity: 1;
	}
</style>
