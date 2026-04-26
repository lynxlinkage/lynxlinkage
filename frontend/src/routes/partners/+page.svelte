<script lang="ts">
	import type { PageData } from './$types';
	import type { Partner, PartnerTier } from '$lib/api/types';
	import Hero from '$lib/components/Hero.svelte';
	import LogoWall from '$lib/components/LogoWall.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { tierLabel } from '$lib/format';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const TIER_ORDER: PartnerTier[] = ['strategic', 'exchange', 'broker', 'tech'];

	const grouped = $derived.by(() => {
		const buckets: Partial<Record<PartnerTier, Partner[]>> = {};
		for (const p of data.partners) {
			let bucket = buckets[p.tier];
			if (!bucket) {
				bucket = [];
				buckets[p.tier] = bucket;
			}
			bucket.push(p);
		}
		return TIER_ORDER.filter((t) => buckets[t]?.length).map((t) => [t, buckets[t]!] as const);
	});
</script>

<Seo
	title="Partners"
	description="Exchanges, brokers, and infrastructure partners that lynxlinkage works with around the world."
	path="/partners"
/>

<Hero
	eyebrow="Partners"
	title="The venues and infrastructure behind our trading."
	subtitle="We work closely with leading exchanges, brokers, and technology providers across regions."
/>

<section class="section">
	<div class="container">
		{#if data.partners.length}
			{#each grouped as [tier, partners] (tier)}
				<div class="tier">
					<h2 class="tier__title">{tierLabel(tier)}</h2>
					<LogoWall {partners} />
				</div>
			{/each}
		{:else}
			<p class="muted">Partner directory coming soon.</p>
		{/if}
	</div>
</section>

<style>
	.tier {
		margin-bottom: var(--space-10);
	}
	.tier__title {
		font-size: var(--text-sm);
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: var(--text-faint);
		margin-bottom: var(--space-5);
		padding-bottom: var(--space-3);
		border-bottom: 1px solid var(--border);
		font-weight: 600;
	}
</style>
