<script lang="ts">
	import { page } from '$app/state';
	import Button from '$lib/components/Button.svelte';
</script>

<section class="section">
	<div class="container-narrow err">
		<span class="eyebrow">Error {page.status}</span>
		<h1>{page.status === 404 ? 'Page not found' : 'Something went wrong'}</h1>
		<p>
			{page.error?.message ?? 'The page you are looking for is no longer available.'}
		</p>
		<div class="cta">
			<Button href="/" variant="primary">Back to home</Button>
		</div>
	</div>
</section>

<style>
	.err {
		text-align: center;
		padding-top: var(--space-12);
	}
	.cta {
		margin-top: var(--space-6);
		display: inline-flex;
	}
</style>
