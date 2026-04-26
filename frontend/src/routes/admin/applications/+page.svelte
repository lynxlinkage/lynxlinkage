<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import {
		ApiError,
		adminGetApplication,
		adminListApplications,
		applicationFileUrl
	} from '$lib/api/client';
	import { auth } from '$lib/auth.svelte';
	import type { Application } from '$lib/api/types';

	let items = $state<Application[]>([]);
	let loading = $state(true);
	let listError = $state<string | null>(null);

	let selectedId = $state<number | null>(null);
	let selected = $state<Application | null>(null);
	let detailLoading = $state(false);
	let detailError = $state<string | null>(null);

	let filterJobId = $state<string>('');

	onMount(async () => {
		const user = await auth.load();
		if (!user) {
			void goto(`/login?next=${encodeURIComponent('/admin/applications')}`, { replaceState: true });
			return;
		}
		if (user.role !== 'hr') {
			listError = 'Your account is not authorised to view applications.';
			loading = false;
			return;
		}
		await refresh();
	});

	async function refresh() {
		loading = true;
		listError = null;
		try {
			const jobId = filterJobId.trim() ? Number(filterJobId.trim()) : undefined;
			items = await adminListApplications(jobId);
			if (selectedId != null && !items.some((x) => x.id === selectedId)) {
				selectedId = null;
				selected = null;
			}
		} catch (err) {
			if (err instanceof ApiError && err.status === 401) {
				void goto(`/login?next=${encodeURIComponent('/admin/applications')}`, {
					replaceState: true
				});
				return;
			}
			listError = err instanceof Error ? err.message : 'Failed to load applications';
		} finally {
			loading = false;
		}
	}

	async function openDetail(app: Application) {
		selectedId = app.id;
		detailError = null;
		detailLoading = true;
		try {
			selected = await adminGetApplication(app.id);
		} catch (err) {
			detailError = err instanceof Error ? err.message : 'Failed to load application';
			selected = null;
		} finally {
			detailLoading = false;
		}
	}

	function closeDetail() {
		selectedId = null;
		selected = null;
		detailError = null;
	}

	async function onLogout() {
		await auth.logout();
		void goto('/login', { replaceState: true });
	}

	function fmtDateTime(iso: string): string {
		const t = Date.parse(iso);
		if (Number.isNaN(t)) return iso;
		return new Date(t).toLocaleString();
	}

	function fmtBytes(n: number): string {
		if (n < 1024) return `${n} B`;
		if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
		return `${(n / (1024 * 1024)).toFixed(1)} MB`;
	}
</script>

<svelte:head>
	<title>Applications · Admin · Lynxlinkage</title>
	<meta name="robots" content="noindex,nofollow" />
</svelte:head>

<section class="admin">
	<div class="container">
		<header class="admin__header">
			<div>
				<p class="muted small">Admin</p>
				<h1>Applications</h1>
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
			<a class="admin__tab" href="/admin">Job postings</a>
			<a class="admin__tab admin__tab--active" href="/admin/applications">Applications</a>
		</nav>

		<div class="admin__toolbar">
			<form
				class="filter"
				onsubmit={(e) => {
					e.preventDefault();
					void refresh();
				}}
			>
				<label class="filter__field">
					<span class="muted small">Filter by job ID</span>
					<input
						type="number"
						min="1"
						placeholder="all jobs"
						bind:value={filterJobId}
					/>
				</label>
				<button type="submit" class="ghost">Apply filter</button>
				{#if filterJobId}
					<button
						type="button"
						class="ghost"
						onclick={() => {
							filterJobId = '';
							void refresh();
						}}>Clear</button
					>
				{/if}
			</form>
			<button type="button" class="ghost" onclick={refresh}>Refresh</button>
		</div>

		{#if listError}
			<p class="error">{listError}</p>
		{/if}

		<div class="layout" class:layout--with-detail={selectedId != null}>
			<div class="layout__list">
				{#if loading}
					<p class="muted">Loading…</p>
				{:else if items.length === 0}
					<p class="muted">No applications yet.</p>
				{:else}
					<div class="table-wrapper">
						<table class="table">
							<thead>
								<tr>
									<th>Submitted</th>
									<th>Candidate</th>
									<th>Email</th>
									<th>Role</th>
									<th class="num">Files</th>
								</tr>
							</thead>
							<tbody>
								{#each items as app (app.id)}
									<tr
										class:selected={selectedId === app.id}
										onclick={() => openDetail(app)}
									>
										<td title={app.createdAt}>{fmtDateTime(app.createdAt)}</td>
										<td><strong>{app.name}</strong></td>
										<td>{app.email}</td>
										<td>
											{#if app.jobTitle}
												<a
													href="/hiring/{app.jobId}"
													target="_blank"
													rel="noopener"
													onclick={(e) => e.stopPropagation()}
												>
													{app.jobTitle}
												</a>
											{:else}
												<span class="muted">job #{app.jobId}</span>
											{/if}
										</td>
										<td class="num">{app.files?.length ?? '—'}</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</div>

			{#if selectedId != null}
				<aside class="layout__detail" aria-label="Application detail">
					<div class="detail__head">
						<h2>Application</h2>
						<button type="button" class="ghost" onclick={closeDetail} aria-label="Close detail"
							>&times;</button
						>
					</div>

					{#if detailLoading}
						<p class="muted">Loading…</p>
					{:else if detailError}
						<p class="error">{detailError}</p>
					{:else if selected}
						<dl class="kv">
							<dt>Submitted</dt>
							<dd>{fmtDateTime(selected.createdAt)}</dd>

							<dt>Name</dt>
							<dd>{selected.name}</dd>

							<dt>Email</dt>
							<dd><a href="mailto:{selected.email}">{selected.email}</a></dd>

							<dt>Role</dt>
							<dd>
								{#if selected.jobTitle}
									<a href="/hiring/{selected.jobId}" target="_blank" rel="noopener"
										>{selected.jobTitle}</a
									>
								{:else}
									<span class="muted">job #{selected.jobId}</span>
								{/if}
							</dd>
						</dl>

						{#if selected.message}
							<h3>Message</h3>
							<p class="message">{selected.message}</p>
						{/if}

						<h3>Attachments ({selected.files?.length ?? 0})</h3>
						{#if selected.files && selected.files.length > 0}
							<ul class="files">
								{#each selected.files as f (f.id)}
									<li>
										<a
											href={applicationFileUrl(selected.id, f.id)}
											target="_blank"
											rel="noopener"
											download={f.originalName}
										>
											{f.originalName}
										</a>
										<span class="muted small">{f.contentType || '—'}</span>
										<span class="muted small">{fmtBytes(f.sizeBytes)}</span>
									</li>
								{/each}
							</ul>
						{:else}
							<p class="muted small">No files attached.</p>
						{/if}
					{/if}
				</aside>
			{/if}
		</div>
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
	.num {
		text-align: right;
		font-variant-numeric: tabular-nums;
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
		flex-wrap: wrap;
		gap: var(--space-3);
		align-items: flex-end;
		margin-bottom: var(--space-4);
	}
	.filter {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-2);
		align-items: flex-end;
	}
	.filter__field {
		display: grid;
		gap: 0.3rem;
	}
	.filter__field input {
		padding: 0.5rem 0.65rem;
		font: inherit;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
		color: var(--text);
		min-width: 140px;
	}

	.layout {
		display: grid;
		gap: var(--space-4);
		grid-template-columns: 1fr;
	}
	.layout--with-detail {
		grid-template-columns: 1fr;
	}
	@media (min-width: 1024px) {
		.layout--with-detail {
			grid-template-columns: minmax(0, 2fr) minmax(320px, 1fr);
		}
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
	.table th.num {
		text-align: right;
	}
	.table tbody tr {
		cursor: pointer;
	}
	.table tbody tr:hover {
		background: var(--surface-muted, #f7f8fb);
	}
	.table tbody tr.selected {
		background: var(--accent-soft);
	}
	.table tbody tr:last-child td {
		border-bottom: none;
	}
	.table a {
		color: var(--accent);
		text-decoration: none;
		font-weight: 500;
	}
	.table a:hover {
		text-decoration: underline;
	}

	.layout__detail {
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: var(--space-5);
		align-self: flex-start;
		position: sticky;
		top: var(--space-4);
	}
	.detail__head {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--space-3);
	}
	.detail__head h2 {
		margin: 0;
		font-size: var(--text-xl);
	}

	.kv {
		display: grid;
		grid-template-columns: max-content 1fr;
		column-gap: var(--space-3);
		row-gap: 0.4rem;
		margin: 0 0 var(--space-4);
		font-size: var(--text-sm);
	}
	.kv dt {
		color: var(--text-muted);
	}
	.kv dd {
		margin: 0;
		word-break: break-word;
	}
	.layout__detail h3 {
		font-size: var(--text-base);
		margin: var(--space-3) 0 var(--space-2);
	}
	.message {
		white-space: pre-wrap;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		padding: 0.7rem 0.9rem;
		font-size: var(--text-sm);
		line-height: 1.6;
	}
	.files {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	.files li {
		display: grid;
		grid-template-columns: 1fr auto auto;
		gap: var(--space-3);
		align-items: baseline;
		padding: 0.55rem 0.7rem;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
		font-size: var(--text-sm);
	}
	.files li a {
		color: var(--accent);
		text-decoration: none;
		font-weight: 500;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	.files li a:hover {
		text-decoration: underline;
	}

	.ghost {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0.55rem 0.9rem;
		font-size: var(--text-sm);
		font-weight: 500;
		border-radius: var(--radius-sm);
		cursor: pointer;
		border: 1px solid var(--border);
		background: transparent;
		color: var(--text);
		transition:
			background-color 120ms var(--ease-out),
			border-color 120ms var(--ease-out),
			color 120ms var(--ease-out);
	}
	.ghost:hover:not(:disabled) {
		border-color: var(--accent);
		color: var(--accent);
	}

	.error {
		margin: 0 0 var(--space-3);
		padding: 0.7rem 0.9rem;
		background: var(--danger-soft, #fdecec);
		color: var(--danger, #b42318);
		border-radius: var(--radius-sm);
		font-size: var(--text-sm);
	}
</style>
