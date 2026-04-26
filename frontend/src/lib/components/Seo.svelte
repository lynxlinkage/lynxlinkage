<script lang="ts">
	import { site } from '$lib/site';

	interface Props {
		title?: string;
		description?: string;
		path?: string;
		image?: string;
	}

	let { title, description = site.description, path = '/', image }: Props = $props();

	const fullTitle = $derived(title ? `${title} | ${site.name}` : `${site.name} — ${site.tagline}`);
	const url = $derived(new URL(path, site.url).toString());
	// og.svg is shipped as a placeholder. For production OG previews on
	// Twitter/Facebook, export a 1200x630 PNG to /static/og.png and pass
	// `image="/og.png"` (or change the default below).
	const ogImage = $derived(image ?? new URL('/og.svg', site.url).toString());
</script>

<svelte:head>
	<title>{fullTitle}</title>
	<meta name="description" content={description} />
	<link rel="canonical" href={url} />

	<meta property="og:type" content="website" />
	<meta property="og:site_name" content={site.name} />
	<meta property="og:title" content={fullTitle} />
	<meta property="og:description" content={description} />
	<meta property="og:url" content={url} />
	<meta property="og:image" content={ogImage} />

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content={fullTitle} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content={ogImage} />
</svelte:head>
