<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAnchorAttributes, HTMLButtonAttributes } from 'svelte/elements';

	type Variant = 'primary' | 'secondary' | 'ghost';
	type Size = 'md' | 'lg';

	type CommonProps = {
		variant?: Variant;
		size?: Size;
		children: Snippet;
	};

	type AsButton = CommonProps & { href?: undefined } & Omit<HTMLButtonAttributes, 'children'>;
	type AsLink = CommonProps & { href: string } & Omit<HTMLAnchorAttributes, 'children'>;

	let { variant = 'primary', size = 'md', children, ...rest }: AsButton | AsLink = $props();

	const classes = $derived(['btn', `btn--${variant}`, `btn--${size}`].join(' '));
</script>

{#if 'href' in rest && rest.href}
	<a class={classes} {...rest}>
		{@render children()}
	</a>
{:else}
	<button class={classes} {...rest as HTMLButtonAttributes}>
		{@render children()}
	</button>
{/if}

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		font-weight: 500;
		font-family: inherit;
		text-decoration: none;
		cursor: pointer;
		border: 1px solid transparent;
		border-radius: var(--radius);
		transition:
			background-color 140ms var(--ease-out),
			border-color 140ms var(--ease-out),
			color 140ms var(--ease-out),
			box-shadow 140ms var(--ease-out),
			transform 140ms var(--ease-out);
		white-space: nowrap;
		user-select: none;
	}
	.btn:disabled {
		opacity: 0.55;
		cursor: not-allowed;
	}
	.btn:focus-visible {
		outline: 2px solid var(--accent);
		outline-offset: 2px;
	}

	.btn--md {
		font-size: var(--text-sm);
		padding: 0.6rem 1.1rem;
		line-height: 1.2;
	}
	.btn--lg {
		font-size: var(--text-base);
		padding: 0.85rem 1.5rem;
		line-height: 1.2;
	}

	.btn--primary {
		background: var(--accent);
		color: var(--accent-ink);
		border-color: var(--accent);
		box-shadow: var(--shadow-sm);
	}
	.btn--primary:hover {
		background: var(--accent-strong);
		border-color: var(--accent-strong);
		box-shadow: var(--shadow-md);
	}

	.btn--secondary {
		background: var(--bg);
		color: var(--text);
		border-color: var(--border-strong);
	}
	.btn--secondary:hover {
		background: var(--surface);
		border-color: var(--accent);
		color: var(--accent);
	}

	.btn--ghost {
		background: transparent;
		color: var(--text);
		border-color: transparent;
	}
	.btn--ghost:hover {
		background: var(--surface);
		color: var(--accent);
	}
</style>
