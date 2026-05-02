<script lang="ts">
	import { onMount } from 'svelte';
	import {
		ApiError,
		adminCreateJob,
		adminDeleteJob,
		adminListJobs,
		adminUpdateJob
	} from '$lib/api/client';
	import { auth } from '$lib/auth.svelte';
	import type { EmploymentType, JobPosting, JobUpsertPayload } from '$lib/api/types';

	type Mode = 'list' | 'create' | 'edit';

	let mode = $state<Mode>('list');
	let jobs = $state<JobPosting[]>([]);
	let loading = $state(true);
	let listError = $state<string | null>(null);

	let editingId = $state<number | null>(null);
	let formError = $state<string | null>(null);
	let saving = $state(false);

	const employmentTypes: { value: EmploymentType; label: string }[] = [
		{ value: 'full_time', label: 'Full-time' },
		{ value: 'part_time', label: 'Part-time' },
		{ value: 'contract', label: 'Contract' },
		{ value: 'internship', label: 'Internship' }
	];

	const blankForm: JobUpsertPayload = {
		title: '',
		team: '',
		location: 'Remote (Taiwan)',
		employmentType: 'full_time',
		descriptionMd: '',
		applyUrlOrEmail: '',
		isActive: true
	};
	let form = $state<JobUpsertPayload>({ ...blankForm });

	onMount(async () => {
		const user = await auth.load();
		if (!user) {
			// Full page reload so Traefik+Authelia can gate the request and
			// redirect to the Authelia login portal if needed.
			window.location.href = '/admin';
			return;
		}
		if (user.role !== 'hr') {
			listError = 'Your account is not authorised to manage job postings.';
			loading = false;
			return;
		}
		await refresh();
	});

	async function refresh() {
		loading = true;
		listError = null;
		try {
			jobs = await adminListJobs();
		} catch (err) {
			if (err instanceof ApiError && err.status === 401) {
				window.location.href = '/admin';
				return;
			}
			listError = err instanceof Error ? err.message : 'Failed to load jobs';
		} finally {
			loading = false;
		}
	}

	// postedAtDisplay is shown as read-only text when editing an existing job.
	let postedAtDisplay = $state<string>('');

	function startCreate() {
		form = { ...blankForm };
		postedAtDisplay = '';
		editingId = null;
		formError = null;
		mode = 'create';
	}

	function startEdit(job: JobPosting) {
		form = {
			title: job.title,
			team: job.team,
			location: job.location,
			employmentType: job.employmentType,
			descriptionMd: job.descriptionMd,
			applyUrlOrEmail: job.applyUrlOrEmail,
			isActive: job.isActive
		};
		postedAtDisplay = job.postedAt ? job.postedAt.slice(0, 10) : '';
		editingId = job.id;
		formError = null;
		mode = 'edit';
	}

	function cancelEdit() {
		formError = null;
		mode = 'list';
		editingId = null;
		postedAtDisplay = '';
	}

	async function onSave(e: SubmitEvent) {
		e.preventDefault();
		if (saving) return;
		formError = null;
		saving = true;
		try {
			const payload: JobUpsertPayload = {
				...form,
				title: form.title.trim(),
				team: form.team.trim(),
				location: form.location.trim(),
				applyUrlOrEmail: form.applyUrlOrEmail.trim(),
				descriptionMd: form.descriptionMd
			};
			if (mode === 'edit' && editingId != null) {
				await adminUpdateJob(editingId, payload);
			} else {
				await adminCreateJob(payload);
			}
			await refresh();
			mode = 'list';
			editingId = null;
		} catch (err) {
			formError = err instanceof Error ? err.message : 'Save failed';
		} finally {
			saving = false;
		}
	}

	async function toggleActive(job: JobPosting) {
		try {
			await adminUpdateJob(job.id, {
				title: job.title,
				team: job.team,
				location: job.location,
				employmentType: job.employmentType,
				descriptionMd: job.descriptionMd,
				applyUrlOrEmail: job.applyUrlOrEmail,
				isActive: !job.isActive
			});
			await refresh();
		} catch (err) {
			listError = err instanceof Error ? err.message : 'Toggle failed';
		}
	}

	async function deleteJob(job: JobPosting) {
		const ok = window.confirm(
			`Permanently delete "${job.title}"?\n\nThis cannot be undone. ` +
				`Use "Hide" instead if you only want to take it off the public site.`
		);
		if (!ok) return;
		try {
			await adminDeleteJob(job.id);
			await refresh();
		} catch (err) {
			listError = err instanceof Error ? err.message : 'Delete failed';
		}
	}

	async function onLogout() {
		await auth.logout(); // redirects browser to Authelia logout
	}

	function fmtDate(iso: string): string {
		try {
			return new Date(iso).toISOString().slice(0, 10);
		} catch {
			return iso;
		}
	}

	function fmtRelative(iso: string | undefined): string {
		if (!iso) return '—';
		const t = Date.parse(iso);
		if (Number.isNaN(t)) return '—';
		const diffMs = Date.now() - t;
		const sec = Math.round(diffMs / 1000);
		if (sec < 60) return 'just now';
		const min = Math.round(sec / 60);
		if (min < 60) return `${min} min ago`;
		const hr = Math.round(min / 60);
		if (hr < 24) return `${hr}h ago`;
		const day = Math.round(hr / 24);
		if (day < 30) return `${day}d ago`;
		return new Date(t).toISOString().slice(0, 10);
	}
</script>

<svelte:head>
	<title>Admin · Lynxlinkage</title>
	<meta name="robots" content="noindex,nofollow" />
</svelte:head>

<section class="admin">
	<div class="container">
		<header class="admin__header">
			<div>
				<p class="muted small">Admin</p>
				<h1>Job postings</h1>
			</div>
			<div class="admin__user">
				{#if auth.user}
					<span class="muted small">
						{auth.user.email} · <strong>{auth.user.role.toUpperCase()}</strong>
					</span>
				{/if}
				<button type="button" class="ghost" onclick={onLogout}>Sign out</button>
			</div>
		</header>

		<nav class="admin__tabs" aria-label="Admin sections">
			<a class="admin__tab admin__tab--active" href="/admin">Job postings</a>
			<a class="admin__tab" href="/admin/applications">Applications</a>
			<a class="admin__tab" href="/admin/statuses">Workflow</a>
		</nav>

		{#if mode === 'list'}
			<div class="admin__toolbar">
				<button type="button" class="primary" onclick={startCreate}>+ New job</button>
				<button type="button" class="ghost" onclick={refresh}>Refresh</button>
			</div>

			{#if listError}
				<p class="error">{listError}</p>
			{/if}

			{#if loading}
				<p class="muted">Loading…</p>
			{:else if jobs.length === 0}
				<p class="muted">No job postings yet. Click "New job" to create the first one.</p>
			{:else}
				<div class="table-wrapper">
					<table class="table">
						<thead>
							<tr>
								<th>Title</th>
								<th>Team</th>
								<th>Type</th>
								<th>Posted</th>
								<th>Updated</th>
								<th>Status</th>
								<th class="actions">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#each jobs as job (job.id)}
								<tr class:inactive={!job.isActive}>
									<td>
										<a href="/hiring/{job.id}" target="_blank" rel="noopener">{job.title}</a>
									</td>
									<td>{job.team || '—'}</td>
									<td>{job.employmentType.replace('_', '-')}</td>
									<td>{fmtDate(job.postedAt)}</td>
									<td title={job.updatedAt ?? ''}>{fmtRelative(job.updatedAt)}</td>
									<td>
										<span class="badge" class:badge--off={!job.isActive}>
											{job.isActive ? 'Active' : 'Hidden'}
										</span>
									</td>
									<td class="actions">
										<button type="button" class="ghost small" onclick={() => startEdit(job)}>
											Edit
										</button>
										<button type="button" class="ghost small" onclick={() => toggleActive(job)}>
											{job.isActive ? 'Hide' : 'Show'}
										</button>
										<button
											type="button"
											class="ghost small danger"
											onclick={() => deleteJob(job)}
											title="Permanently delete this posting"
										>
											Delete
										</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		{:else}
			<form class="job-form" onsubmit={onSave} novalidate>
				<div class="job-form__head">
					<h2>{mode === 'edit' ? 'Edit job' : 'New job'}</h2>
					<button type="button" class="ghost" onclick={cancelEdit}>Cancel</button>
				</div>

				<div class="grid">
					<label class="field">
						<span>Title</span>
						<input type="text" required maxlength="200" bind:value={form.title} />
					</label>

					<label class="field">
						<span>Team</span>
						<input type="text" maxlength="80" bind:value={form.team} placeholder="e.g. Engineering" />
					</label>

					<label class="field">
						<span>Location</span>
						<input type="text" maxlength="120" bind:value={form.location} />
					</label>

					<label class="field">
						<span>Employment type</span>
						<select bind:value={form.employmentType}>
							{#each employmentTypes as et (et.value)}
								<option value={et.value}>{et.label}</option>
							{/each}
						</select>
					</label>

				<label class="field">
					<span>Apply URL or email</span>
					<input
						type="text"
						required
						maxlength="500"
						bind:value={form.applyUrlOrEmail}
						placeholder="hiring@example.com or https://…"
					/>
				</label>

				{#if mode === 'edit' && postedAtDisplay}
					<div class="field">
						<span>Posted at</span>
						<p class="field__readonly">{postedAtDisplay}</p>
					</div>
				{/if}

					<label class="field field--toggle">
						<input type="checkbox" bind:checked={form.isActive} />
						<span>Active (visible on the public site)</span>
					</label>
				</div>

				<label class="field">
					<span>Description (Markdown)</span>
					<textarea rows="14" bind:value={form.descriptionMd}></textarea>
					<small class="muted">Headings, lists, bold/italic, links and code blocks are supported.</small>
				</label>

				{#if formError}
					<p class="error">{formError}</p>
				{/if}

				<div class="job-form__actions">
					<button type="submit" class="primary" disabled={saving}>
						{saving ? 'Saving…' : mode === 'edit' ? 'Save changes' : 'Create job'}
					</button>
					<button type="button" class="ghost" onclick={cancelEdit} disabled={saving}>
						Cancel
					</button>
				</div>
			</form>
		{/if}
	</div>
</section>

<style>
	.admin {
		padding: var(--space-7) 0 var(--space-9);
	}
	.admin__header {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		gap: var(--space-3);
		margin-bottom: var(--space-5);
		flex-wrap: wrap;
	}
	.admin__header h1 {
		margin: 0;
		font-size: var(--text-3xl);
		letter-spacing: -0.01em;
	}
	.admin__user {
		display: flex;
		align-items: center;
		gap: var(--space-3);
	}
	.muted {
		color: var(--text-muted);
	}
	.small {
		font-size: var(--text-sm);
	}

	.admin__tabs {
		display: flex;
		gap: var(--space-1);
		border-bottom: 1px solid var(--border);
		margin-bottom: var(--space-5);
	}
	.admin__tab {
		padding: 0.55rem 1rem;
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--text-muted);
		text-decoration: none;
		border-bottom: 2px solid transparent;
		margin-bottom: -1px;
	}
	.admin__tab:hover {
		color: var(--accent);
	}
	.admin__tab--active {
		color: var(--text);
		border-bottom-color: var(--accent);
	}

	.admin__toolbar {
		display: flex;
		gap: var(--space-2);
		margin-bottom: var(--space-4);
	}

	.table-wrapper {
		overflow-x: auto;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		background: var(--bg);
	}
	.table {
		width: 100%;
		border-collapse: collapse;
		font-size: var(--text-sm);
	}
	.table th,
	.table td {
		padding: 0.7rem 0.85rem;
		text-align: left;
		border-bottom: 1px solid var(--border);
		vertical-align: middle;
	}
	.table thead th {
		background: var(--surface-muted, #f7f8fb);
		font-weight: 600;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		font-size: 0.72rem;
	}
	.table tbody tr:last-child td {
		border-bottom: none;
	}
	.table tbody tr.inactive td {
		opacity: 0.55;
	}
	.table .actions {
		display: flex;
		gap: 0.4rem;
		justify-content: flex-end;
	}
	.table th.actions {
		text-align: right;
	}
	.table a {
		color: var(--accent);
		text-decoration: none;
		font-weight: 500;
	}
	.table a:hover {
		text-decoration: underline;
	}

	.badge {
		display: inline-block;
		padding: 0.15rem 0.55rem;
		font-size: 0.72rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		border-radius: 999px;
		background: var(--accent-soft);
		color: var(--accent-strong, var(--accent));
	}
	.badge--off {
		background: var(--surface-muted, #eef0f4);
		color: var(--text-muted);
	}

	.job-form {
		display: grid;
		gap: var(--space-4);
		background: var(--bg);
		padding: var(--space-5);
		border: 1px solid var(--border);
		border-radius: var(--radius);
	}
	.job-form__head {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	.job-form__head h2 {
		margin: 0;
		font-size: var(--text-xl);
	}
	.job-form__actions {
		display: flex;
		gap: var(--space-2);
	}

	.grid {
		display: grid;
		gap: var(--space-3);
		grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
	}

	.field {
		display: grid;
		gap: 0.3rem;
	}
	.field__readonly {
		margin: 0;
		padding: 0.6rem 0.75rem;
		font-size: var(--text-sm);
		color: var(--text-muted);
		background: var(--surface-muted, #f7f8fb);
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
	}

	.field--toggle {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	.field--toggle input {
		width: auto;
	}
	.field span {
		font-size: var(--text-sm);
		font-weight: 500;
	}
	.field small {
		font-size: 0.75rem;
	}
	input[type='text'],
	input[type='date'],
	select,
	textarea {
		width: 100%;
		padding: 0.6rem 0.75rem;
		font: inherit;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
		color: var(--text);
	}
	textarea {
		font-family: var(--font-mono, ui-monospace, SFMono-Regular, monospace);
		font-size: 0.85rem;
		line-height: 1.5;
		resize: vertical;
	}
	input:focus,
	select:focus,
	textarea:focus {
		outline: none;
		border-color: var(--accent);
		box-shadow: 0 0 0 3px var(--accent-soft);
	}

	.primary,
	.ghost {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0.55rem 0.9rem;
		font-size: var(--text-sm);
		font-weight: 500;
		border-radius: var(--radius-sm);
		cursor: pointer;
		border: 1px solid transparent;
		transition:
			background-color 120ms var(--ease-out),
			border-color 120ms var(--ease-out),
			color 120ms var(--ease-out);
	}
	.primary {
		background: var(--accent);
		color: var(--accent-ink);
	}
	.primary:hover:not(:disabled) {
		background: var(--accent-strong);
	}
	.ghost {
		background: transparent;
		color: var(--text);
		border-color: var(--border);
	}
	.ghost:hover:not(:disabled) {
		border-color: var(--accent);
		color: var(--accent);
	}
	.ghost.small {
		padding: 0.35rem 0.6rem;
		font-size: var(--text-sm);
	}
	.ghost.danger {
		color: var(--danger, #b42318);
	}
	.ghost.danger:hover:not(:disabled) {
		border-color: var(--danger, #b42318);
		color: var(--danger, #b42318);
		background: var(--danger-soft, #fdecec);
	}
	button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.error {
		margin: 0;
		padding: 0.7rem 0.9rem;
		background: var(--danger-soft, #fdecec);
		color: var(--danger, #b42318);
		border-radius: var(--radius-sm);
		font-size: var(--text-sm);
	}
</style>
