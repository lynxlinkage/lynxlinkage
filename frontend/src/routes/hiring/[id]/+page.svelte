<script lang="ts">
	import type { PageData } from './$types';
	import Button from '$lib/components/Button.svelte';
	import Seo from '$lib/components/Seo.svelte';
	import { employmentLabel, formatDate } from '$lib/format';
	import { submitApplication } from '$lib/api/client';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
	const job = $derived(data.job);

	const MAX_FILES = 3;
	const MAX_FILE_BYTES = 10 * 1024 * 1024;

	// Only PDF and common image formats are accepted.
	const ALLOWED_MIME = new Set([
		'application/pdf',
		'image/jpeg',
		'image/png',
		'image/gif',
		'image/webp'
	]);
	const ALLOWED_EXT = /\.(pdf|jpe?g|png|gif|webp)$/i;

	function isAllowedFile(f: File): boolean {
		// MIME type from the browser is reliable for these well-known types;
		// fall back to extension when the browser reports an empty string.
		return f.type ? ALLOWED_MIME.has(f.type) : ALLOWED_EXT.test(f.name);
	}

	let name = $state('');
	let email = $state('');
	let message = $state('');
	let files = $state<File[]>([]);
	let formError = $state<string | null>(null);
	let submitting = $state(false);
	let submitted = $state(false);

	let fileInputEl = $state<HTMLInputElement | null>(null);

	function fmtBytes(n: number): string {
		if (n < 1024) return `${n} B`;
		if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
		return `${(n / (1024 * 1024)).toFixed(1)} MB`;
	}

	function onPickFiles(e: Event) {
		const input = e.target as HTMLInputElement;
		const picked = input.files ? Array.from(input.files) : [];
		const next = [...files];
		for (const f of picked) {
			if (!isAllowedFile(f)) {
				formError = `"${f.name}" is not allowed. Only PDF and image files (JPEG, PNG, GIF, WebP) are accepted.`;
				input.value = '';
				return;
			}
			if (next.some((x) => x.name === f.name && x.size === f.size)) continue;
			next.push(f);
		}
		const truncated = next.slice(0, MAX_FILES);
		const oversized = truncated.find((f) => f.size > MAX_FILE_BYTES);
		if (oversized) {
			formError = `"${oversized.name}" is larger than 10 MB. Please upload a smaller file.`;
			input.value = '';
			return;
		}
		if (next.length > MAX_FILES) {
			formError = `You can attach at most ${MAX_FILES} files.`;
		} else {
			formError = null;
		}
		files = truncated;
		input.value = '';
	}

	function removeFile(idx: number) {
		files = files.filter((_, i) => i !== idx);
		formError = null;
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();
		if (submitting) return;
		formError = null;

		if (!name.trim()) {
			formError = 'Please enter your name.';
			return;
		}
		if (!email.trim() || !email.includes('@')) {
			formError = 'Please enter a valid email.';
			return;
		}
		if (files.length > MAX_FILES) {
			formError = `You can attach at most ${MAX_FILES} files.`;
			return;
		}
		for (const f of files) {
			if (!isAllowedFile(f)) {
				formError = `"${f.name}" is not a PDF or image file.`;
				return;
			}
			if (f.size > MAX_FILE_BYTES) {
				formError = `"${f.name}" is larger than 10 MB.`;
				return;
			}
		}

		submitting = true;
		const result = await submitApplication(job.id, {
			name: name.trim(),
			email: email.trim(),
			message: message.trim(),
			files
		});
		submitting = false;
		if (!result.ok) {
			formError = result.error ?? 'Submission failed. Please try again.';
			return;
		}
		submitted = true;
	}
</script>

<Seo
	title={job.title}
	description={`Open role at LynxLinkage: ${job.title} on the ${job.team} team, ${job.location}.`}
	path={`/hiring/${job.id}`}
/>

<article class="job">
	<header class="job__header">
		<div class="container-narrow">
			<a class="job__back" href="/hiring">
				<span aria-hidden="true">&larr;</span> All open roles
			</a>
			<div class="job__meta">
				<span class="job__team">{job.team || 'Team'}</span>
				<span class="job__sep" aria-hidden="true">&middot;</span>
				<span>{job.location}</span>
				<span class="job__sep" aria-hidden="true">&middot;</span>
				<span>{employmentLabel(job.employmentType)}</span>
				<span class="job__sep" aria-hidden="true">&middot;</span>
				<span>Posted {formatDate(job.postedAt)}</span>
			</div>
			<h1 class="job__title">{job.title}</h1>
			<div class="job__cta">
				<Button href="#apply" variant="primary" size="lg">Apply now</Button>
				<Button href="/hiring" variant="ghost">Other roles</Button>
			</div>
		</div>
	</header>

	<section class="section">
		<div class="container-narrow">
			<!-- eslint-disable-next-line svelte/no-at-html-tags -->
			<div class="prose">{@html data.descriptionHtml}</div>

			<div id="apply" class="apply">
				{#if submitted}
					<div class="apply__success" role="status" aria-live="polite">
						<h2>Thanks — we got it.</h2>
						<p>
							Your application for <strong>{job.title}</strong> is in our queue. Our team reviews
							every submission; expect to hear back within a couple of weeks.
						</p>
						<Button href="/hiring" variant="ghost">Browse other roles</Button>
					</div>
				{:else}
					<h2>Apply for {job.title}</h2>
					<p class="apply__lede">
						Tell us a bit about yourself and attach your CV. A short note about why this role
						catches your eye goes a long way.
					</p>

					<form class="apply__form" novalidate onsubmit={handleSubmit}>
						<div class="apply__row">
							<label class="apply__field">
								<span>Full name</span>
								<input
									type="text"
									name="name"
									autocomplete="name"
									required
									maxlength="200"
									bind:value={name}
									disabled={submitting}
								/>
							</label>
							<label class="apply__field">
								<span>Email</span>
								<input
									type="email"
									name="email"
									autocomplete="email"
									required
									maxlength="320"
									bind:value={email}
									disabled={submitting}
								/>
							</label>
						</div>

						<label class="apply__field">
							<span>Anything you'd like us to know <em>(optional)</em></span>
							<textarea
								name="message"
								rows="5"
								maxlength="4096"
								placeholder="Why this role, what you're working on, links to writing or code…"
								bind:value={message}
								disabled={submitting}
							></textarea>
						</label>

						<div class="apply__field">
							<span>Attachments</span>
							<p class="apply__hint">
								Up to {MAX_FILES} files, max 10 MB each. PDF and images only (JPEG, PNG, GIF, WebP).
								CV, cover letter, samples — anything that helps us read you.
							</p>
							<input
								bind:this={fileInputEl}
								type="file"
								name="files"
								multiple
								accept="application/pdf,image/jpeg,image/png,image/gif,image/webp"
								onchange={onPickFiles}
								disabled={submitting || files.length >= MAX_FILES}
							/>
							{#if files.length > 0}
								<ul class="apply__files">
									{#each files as file, idx (`${file.name}-${file.size}-${idx}`)}
										<li>
											<span class="apply__file-name">{file.name}</span>
											<span class="apply__file-size">{fmtBytes(file.size)}</span>
											<button
												type="button"
												class="apply__file-remove"
												aria-label="Remove {file.name}"
												onclick={() => removeFile(idx)}
												disabled={submitting}
											>
												Remove
											</button>
										</li>
									{/each}
								</ul>
							{/if}
						</div>

						{#if formError}
							<div class="apply__error" role="alert">{formError}</div>
						{/if}

						<div class="apply__actions">
							<Button type="submit" variant="primary" disabled={submitting}>
								{submitting ? 'Sending…' : `Submit application`}
							</Button>
							<span class="apply__privacy">
								We'll only use the info you share to evaluate this application.
							</span>
						</div>
					</form>
				{/if}
			</div>
		</div>
	</section>
</article>

<style>
	.job__header {
		padding-block: var(--space-9) var(--space-7);
		background: var(--surface);
		border-bottom: 1px solid var(--border);
	}
	.job__back {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		color: var(--text-muted);
		font-size: var(--text-sm);
		text-decoration: none;
		margin-bottom: var(--space-5);
	}
	.job__back:hover {
		color: var(--accent);
	}
	.job__meta {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
		align-items: center;
		font-size: var(--text-sm);
		color: var(--text-muted);
		margin-bottom: var(--space-3);
	}
	.job__team {
		font-weight: 500;
		color: var(--accent);
	}
	.job__sep {
		opacity: 0.5;
	}
	.job__title {
		font-size: var(--text-3xl);
		margin: 0 0 var(--space-5);
		letter-spacing: -0.02em;
	}
	.job__cta {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
	}

	.prose :global(h1),
	.prose :global(h2),
	.prose :global(h3) {
		margin-top: var(--space-7);
		margin-bottom: var(--space-3);
	}
	.prose :global(h2) {
		font-size: var(--text-xl);
	}
	.prose :global(h3) {
		font-size: var(--text-lg);
	}
	.prose :global(p) {
		margin-bottom: var(--space-4);
		line-height: 1.7;
	}
	.prose :global(ul),
	.prose :global(ol) {
		margin: 0 0 var(--space-4) var(--space-5);
		padding: 0;
	}
	.prose :global(li) {
		margin-bottom: var(--space-2);
		line-height: 1.7;
	}
	.prose :global(strong) {
		color: var(--text);
	}

	.apply {
		margin-top: var(--space-8);
		padding: var(--space-6);
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		scroll-margin-top: 80px;
	}
	.apply h2 {
		font-size: var(--text-xl);
		margin: 0 0 var(--space-2);
	}
	.apply__lede {
		color: var(--text-muted);
		margin-bottom: var(--space-5);
	}

	.apply__form {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}
	.apply__row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--space-4);
	}
	@media (max-width: 640px) {
		.apply__row {
			grid-template-columns: 1fr;
		}
	}

	.apply__field {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}
	.apply__field > span {
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--text);
	}
	.apply__field em {
		font-style: normal;
		color: var(--text-muted);
		font-weight: 400;
	}
	.apply__hint {
		font-size: var(--text-sm);
		color: var(--text-muted);
		margin: 0;
	}

	.apply__field input[type='text'],
	.apply__field input[type='email'],
	.apply__field textarea {
		font: inherit;
		padding: 0.7rem 0.8rem;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm, 6px);
		background: var(--bg);
		color: var(--text);
		width: 100%;
	}
	.apply__field input[type='text']:focus,
	.apply__field input[type='email']:focus,
	.apply__field textarea:focus {
		outline: 2px solid var(--accent);
		outline-offset: 1px;
		border-color: var(--accent);
	}
	.apply__field input[type='file'] {
		font: inherit;
		font-size: var(--text-sm);
	}
	.apply__field textarea {
		resize: vertical;
		min-height: 100px;
	}

	.apply__files {
		list-style: none;
		padding: 0;
		margin: 0.5rem 0 0;
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}
	.apply__files li {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: 0.5rem 0.7rem;
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius-sm, 6px);
		font-size: var(--text-sm);
	}
	.apply__file-name {
		flex: 1 1 auto;
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	.apply__file-size {
		color: var(--text-muted);
		font-variant-numeric: tabular-nums;
	}
	.apply__file-remove {
		font: inherit;
		font-size: var(--text-sm);
		background: transparent;
		border: none;
		color: var(--text-muted);
		cursor: pointer;
		padding: 0.1rem 0.3rem;
		border-radius: 4px;
	}
	.apply__file-remove:hover {
		color: var(--danger);
	}

	.apply__error {
		padding: 0.7rem 0.9rem;
		background: var(--danger-soft);
		border: 1px solid color-mix(in oklch, var(--danger) 45%, var(--color-gray-800));
		color: var(--color-rose-200);
		border-radius: var(--radius-sm, 6px);
		font-size: var(--text-sm);
	}

	.apply__actions {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
		align-items: center;
		margin-top: var(--space-2);
	}
	.apply__privacy {
		font-size: var(--text-sm);
		color: var(--text-muted);
	}

	.apply__success h2 {
		margin-bottom: var(--space-3);
	}
	.apply__success p {
		color: var(--text-muted);
		margin-bottom: var(--space-4);
	}
</style>
