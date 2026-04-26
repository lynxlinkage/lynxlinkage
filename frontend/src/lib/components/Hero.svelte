<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		eyebrow?: string;
		title: string;
		subtitle?: string;
		cta?: Snippet;
	}

	let { eyebrow, title, subtitle, cta }: Props = $props();
</script>

<section class="hero">
	<div class="container hero__inner">
		{#if eyebrow}
			<span class="eyebrow">{eyebrow}</span>
		{/if}
		<h1 class="hero__title">{title}</h1>
		{#if subtitle}
			<p class="hero__subtitle">{subtitle}</p>
		{/if}
		{#if cta}
			<div class="hero__cta">{@render cta()}</div>
		{/if}
	</div>
	<div class="hero__pattern" aria-hidden="true"></div>
</section>

<style>
	.hero {
		position: relative;
		padding-block: clamp(var(--space-10), 9vw, var(--space-16));
		padding-top: calc(var(--header-height) + clamp(var(--space-8), 6vw, var(--space-12)));
		overflow: hidden;
	}

	.hero__inner {
		position: relative;
		z-index: 1;
		max-width: 920px;
	}

	.hero__title {
		max-width: 18ch;
		margin-bottom: var(--space-5);
	}

	.hero__subtitle {
		max-width: 56ch;
		font-size: clamp(var(--text-lg), 1.5vw + 0.6rem, var(--text-xl));
		color: var(--text-muted);
		margin-bottom: var(--space-7);
	}

	.hero__cta {
		display: flex;
		gap: var(--space-3);
		flex-wrap: wrap;
	}

	.hero__pattern {
		position: absolute;
		inset: 0;
		background:
			radial-gradient(circle at 1px 1px, rgba(11, 42, 107, 0.08) 1px, transparent 1.5px) 0 0 / 28px
				28px;
		mask-image: linear-gradient(to bottom, transparent 0%, black 25%, black 70%, transparent 100%);
		-webkit-mask-image: linear-gradient(
			to bottom,
			transparent 0%,
			black 25%,
			black 70%,
			transparent 100%
		);
		pointer-events: none;
		z-index: 0;
	}
</style>
