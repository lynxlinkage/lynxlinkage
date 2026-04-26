<script lang="ts">
	import { untrack } from 'svelte';
	import type { ContactKind } from '$lib/api/types';
	import { submitContact } from '$lib/api/client';

	interface Props {
		defaultKind?: ContactKind;
		title?: string;
		subtitle?: string;
	}

	let { defaultKind = 'general', title = 'Get in touch', subtitle }: Props = $props();

	let name = $state('');
	let email = $state('');
	let company = $state('');
	let message = $state('');
	// One-time initialization from the prop. `untrack` snapshots the prop
	// without subscribing to it, so subsequent prop updates don't clobber
	// the user's selection in the dropdown.
	let kind = $state<ContactKind>(untrack(() => defaultKind));
	let status = $state<'idle' | 'submitting' | 'ok' | 'error'>('idle');
	let errorMsg = $state('');

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (status === 'submitting') return;
		status = 'submitting';
		errorMsg = '';

		const result = await submitContact({
			name: name.trim(),
			email: email.trim(),
			company: company.trim() || undefined,
			message: message.trim(),
			kind
		});

		if (result.ok) {
			status = 'ok';
			name = email = company = message = '';
		} else {
			status = 'error';
			errorMsg = result.error ?? 'Something went wrong.';
		}
	}
</script>

<section class="contact">
	<header class="contact__head">
		<h2>{title}</h2>
		{#if subtitle}
			<p>{subtitle}</p>
		{/if}
	</header>

	<form class="contact__form" onsubmit={handleSubmit} novalidate>
		<div class="row">
			<label class="field">
				<span>Name</span>
				<input type="text" bind:value={name} required minlength="2" maxlength="120" autocomplete="name" />
			</label>
			<label class="field">
				<span>Email</span>
				<input type="email" bind:value={email} required maxlength="254" autocomplete="email" />
			</label>
		</div>

		<div class="row">
			<label class="field">
				<span>Company <em>(optional)</em></span>
				<input type="text" bind:value={company} maxlength="200" autocomplete="organization" />
			</label>
			<label class="field">
				<span>Topic</span>
				<select bind:value={kind}>
					<option value="general">General</option>
					<option value="partnership">Partnership</option>
					<option value="research">Research</option>
					<option value="hiring">Hiring</option>
				</select>
			</label>
		</div>

		<label class="field">
			<span>Message</span>
			<textarea bind:value={message} required minlength="10" maxlength="5000" rows="5"></textarea>
		</label>

		<div class="contact__footer">
			<button class="contact__submit" type="submit" disabled={status === 'submitting'}>
				{status === 'submitting' ? 'Sending…' : 'Send message'}
			</button>
			{#if status === 'ok'}
				<p class="contact__msg contact__msg--ok" role="status">
					Thanks &mdash; we&rsquo;ll be in touch soon.
				</p>
			{:else if status === 'error'}
				<p class="contact__msg contact__msg--err" role="alert">{errorMsg}</p>
			{/if}
		</div>
	</form>
</section>

<style>
	.contact {
		display: grid;
		gap: var(--space-6);
	}
	.contact__head h2 {
		margin-bottom: var(--space-2);
	}
	.contact__head p {
		margin: 0;
	}

	.contact__form {
		display: grid;
		gap: var(--space-4);
	}

	.row {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
		gap: var(--space-4);
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
		font-size: var(--text-sm);
		color: var(--text);
	}
	.field span {
		font-weight: 500;
	}
	.field em {
		color: var(--text-faint);
		font-style: normal;
	}

	input,
	select,
	textarea {
		font-family: inherit;
		font-size: var(--text-base);
		color: var(--text);
		background: var(--bg);
		border: 1px solid var(--border-strong);
		border-radius: var(--radius-sm);
		padding: 0.7rem 0.85rem;
		transition:
			border-color 140ms var(--ease-out),
			box-shadow 140ms var(--ease-out);
		width: 100%;
	}
	input:focus,
	select:focus,
	textarea:focus {
		outline: none;
		border-color: var(--accent);
		box-shadow: 0 0 0 3px var(--accent-soft);
	}
	textarea {
		resize: vertical;
		min-height: 120px;
	}

	.contact__footer {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
		align-items: center;
	}
	.contact__submit {
		font-family: inherit;
		font-size: var(--text-sm);
		font-weight: 500;
		padding: 0.75rem 1.4rem;
		background: var(--accent);
		color: var(--accent-ink);
		border: 1px solid var(--accent);
		border-radius: var(--radius);
		cursor: pointer;
		transition:
			background-color 140ms var(--ease-out),
			border-color 140ms var(--ease-out);
	}
	.contact__submit:hover:not(:disabled) {
		background: var(--accent-strong);
		border-color: var(--accent-strong);
	}
	.contact__submit:disabled {
		opacity: 0.65;
		cursor: not-allowed;
	}

	.contact__msg {
		margin: 0;
		font-size: var(--text-sm);
	}
	.contact__msg--ok {
		color: var(--success);
	}
	.contact__msg--err {
		color: var(--danger);
	}
</style>
