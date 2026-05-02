<script lang="ts">
	import { onMount } from 'svelte';
	import {
		ApiError,
		adminGetApplication,
		adminListApplications,
		adminListJobs,
		adminListStatuses,
		adminUpdateApplicationStatus,
		applicationFileUrl
	} from '$lib/api/client';
	import { auth } from '$lib/auth.svelte';
	import type {
		Application,
		ApplicationSort,
		ApplicationStatus,
		JobPosting
	} from '$lib/api/types';

	let items = $state<Application[]>([]);
	let loading = $state(true);
	let listError = $state<string | null>(null);

	let statuses = $state<ApplicationStatus[]>([]);
	let jobs = $state<JobPosting[]>([]);

	let filterStatusId = $state<number | ''>('');
	let filterJobId = $state<number | ''>('');
	let sortOrder = $state<ApplicationSort>('newest');

	// `expandedId` — the row currently open inline. `expanded` holds the
	// hydrated detail (loaded lazily on click so the list stays fast).
	let expandedId = $state<number | null>(null);
	let expanded = $state<Application | null>(null);
	let detailLoading = $state(false);
	let detailError = $state<string | null>(null);

	let pendingStatusId = $state<number | ''>('');
	let pendingNote = $state('');
	let savingStatus = $state(false);
	let statusUpdateError = $state<string | null>(null);

	onMount(async () => {
		const user = await auth.load();
		if (!user) {
			window.location.href = '/admin/applications';
			return;
		}
		if (user.role !== 'hr') {
			listError = 'Your account is not authorised to view applications.';
			loading = false;
			return;
		}
		const [statusRes, jobsRes] = await Promise.allSettled([
			adminListStatuses(),
			adminListJobs()
		]);
		if (statusRes.status === 'fulfilled') statuses = statusRes.value;
		if (jobsRes.status === 'fulfilled') jobs = jobsRes.value;
		await refresh();
	});

	async function refresh() {
		loading = true;
		listError = null;
		try {
			items = await adminListApplications({
				jobId: filterJobId === '' ? undefined : Number(filterJobId),
				statusId: filterStatusId === '' ? undefined : Number(filterStatusId),
				sort: sortOrder
			});
			if (expandedId != null && !items.some((x) => x.id === expandedId)) {
				collapse();
			}
		} catch (err) {
			if (err instanceof ApiError && err.status === 401) {
				window.location.href = '/admin/applications';
				return;
			}
			listError = err instanceof Error ? err.message : 'Failed to load applications';
		} finally {
			loading = false;
		}
	}

	function clearFilters() {
		filterStatusId = '';
		filterJobId = '';
		sortOrder = 'newest';
		void refresh();
	}

	const hasFilters = $derived(filterStatusId !== '' || filterJobId !== '' || sortOrder !== 'newest');

	async function toggleExpand(app: Application) {
		// Same row clicked twice — collapse.
		if (expandedId === app.id) {
			collapse();
			return;
		}
		expandedId = app.id;
		expanded = null;
		detailError = null;
		statusUpdateError = null;
		detailLoading = true;
		try {
			const full = await adminGetApplication(app.id);
			// Race guard: user may have collapsed or moved to another
			// row while we were fetching.
			if (expandedId !== app.id) return;
			expanded = full;
			pendingStatusId = full.statusId ?? '';
			pendingNote = '';
		} catch (err) {
			if (expandedId !== app.id) return;
			detailError = err instanceof Error ? err.message : 'Failed to load application';
		} finally {
			if (expandedId === app.id) detailLoading = false;
		}
	}

	function collapse() {
		expandedId = null;
		expanded = null;
		detailError = null;
		statusUpdateError = null;
		pendingStatusId = '';
		pendingNote = '';
	}

	async function saveStatus() {
		if (!expanded || pendingStatusId === '' || savingStatus) return;
		savingStatus = true;
		statusUpdateError = null;
		try {
			const updated = await adminUpdateApplicationStatus(
				expanded.id,
				Number(pendingStatusId),
				pendingNote.trim() || undefined
			);
			expanded = updated;
			pendingNote = '';
			pendingStatusId = updated.statusId ?? '';
			items = items.map((a) =>
				a.id === updated.id
					? {
							...a,
							status: updated.status,
							statusId: updated.statusId,
							statusUpdatedAt: updated.statusUpdatedAt
						}
					: a
			);
		} catch (err) {
			statusUpdateError = err instanceof Error ? err.message : 'Failed to update status';
		} finally {
			savingStatus = false;
		}
	}

	async function onLogout() {
		await auth.logout(); // redirects browser to Authelia logout
	}

	function fmtDateTime(iso: string | undefined): string {
		if (!iso) return '—';
		const t = Date.parse(iso);
		if (Number.isNaN(t)) return iso;
		return new Date(t).toLocaleString();
	}

	function fmtBytes(n: number): string {
		if (n < 1024) return `${n} B`;
		if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
		return `${(n / (1024 * 1024)).toFixed(1)} MB`;
	}

	function badgeStyle(status: ApplicationStatus | undefined): string {
		const c = status?.color || '#64748b';
		return `--badge-color: ${c};`;
	}

	const dirtyStatus = $derived(
		expanded != null && pendingStatusId !== '' && Number(pendingStatusId) !== (expanded.statusId ?? -1)
	);

	function onRowKey(e: KeyboardEvent, app: Application) {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			void toggleExpand(app);
		}
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
			<a class="admin__tab" href="/admin/statuses">Workflow</a>
		</nav>

		<div class="admin__toolbar">
			<form
				class="filters"
				onsubmit={(e) => {
					e.preventDefault();
					void refresh();
				}}
			>
				<label class="filter">
					<span class="muted small">Status</span>
					<select bind:value={filterStatusId} onchange={() => void refresh()}>
						<option value="">All statuses</option>
						{#each statuses as s (s.id)}
							<option value={s.id}>{s.name}</option>
						{/each}
					</select>
				</label>
				<label class="filter">
					<span class="muted small">Role</span>
					<select bind:value={filterJobId} onchange={() => void refresh()}>
						<option value="">All roles</option>
						{#each jobs as j (j.id)}
							<option value={j.id}>{j.title}{j.isActive ? '' : ' (hidden)'}</option>
						{/each}
					</select>
				</label>
				<label class="filter">
					<span class="muted small">Submitted</span>
					<select bind:value={sortOrder} onchange={() => void refresh()}>
						<option value="newest">Newest first</option>
						<option value="oldest">Oldest first</option>
					</select>
				</label>
				{#if hasFilters}
					<button type="button" class="ghost" onclick={clearFilters}>Clear</button>
				{/if}
			</form>
			<button type="button" class="ghost" onclick={refresh}>Refresh</button>
		</div>

		{#if listError}
			<p class="error">{listError}</p>
		{/if}

		{#if loading}
			<p class="muted">Loading…</p>
		{:else if items.length === 0}
			<p class="muted">No applications match these filters.</p>
		{:else}
			<div class="table-wrapper">
				<table class="table">
					<thead>
						<tr>
							<th class="caret-col" aria-hidden="true"></th>
							<th>Submitted</th>
							<th>Candidate</th>
							<th>Email</th>
							<th>Role</th>
							<th>Status</th>
						</tr>
					</thead>
					<tbody>
						{#each items as app (app.id)}
							{@const isOpen = expandedId === app.id}
							<tr
								class="row"
								class:row--open={isOpen}
								tabindex="0"
								role="button"
								aria-expanded={isOpen}
								aria-controls={`detail-${app.id}`}
								onclick={() => toggleExpand(app)}
								onkeydown={(e) => onRowKey(e, app)}
							>
								<td class="caret-col">
									<span class="caret" class:caret--open={isOpen} aria-hidden="true">▸</span>
								</td>
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
								<td>
									{#if app.status}
										<span class="badge" style={badgeStyle(app.status)}>{app.status.name}</span>
									{:else}
										<span class="muted small">—</span>
									{/if}
								</td>
							</tr>

							{#if isOpen}
								<tr class="expansion">
									<td class="expansion__cell" colspan="6" id={`detail-${app.id}`}>
										{#if detailLoading}
											<p class="muted">Loading…</p>
										{:else if detailError}
											<p class="error">{detailError}</p>
										{:else if expanded}
											<div class="detail">
												<header class="detail__head">
													<div class="detail__title">
														<h2>{expanded.name}</h2>
														<span class="muted small">{expanded.email}</span>
														{#if expanded.status}
															<span class="badge" style={badgeStyle(expanded.status)}
																>{expanded.status.name}</span
															>
														{/if}
													</div>
													<button type="button" class="ghost" onclick={collapse}>Collapse</button>
												</header>

												<section class="status-changer">
													<div class="status-changer__row">
														<label class="field">
															<span class="muted small">Move to</span>
															<select bind:value={pendingStatusId} disabled={savingStatus}>
																{#each statuses as s (s.id)}
																	<option value={s.id}>{s.name}</option>
																{/each}
															</select>
														</label>
														<label class="field field--grow">
															<span class="muted small">Note (optional)</span>
															<input
																type="text"
																placeholder="e.g. phone screen passed"
																maxlength="500"
																bind:value={pendingNote}
																disabled={savingStatus}
															/>
														</label>
														<button
															type="button"
															class="primary"
															disabled={!dirtyStatus || savingStatus}
															onclick={saveStatus}
														>
															{savingStatus ? 'Saving…' : 'Save'}
														</button>
													</div>
													{#if statusUpdateError}
														<p class="error small">{statusUpdateError}</p>
													{/if}
												</section>

												<div class="detail__grid">
													<div class="detail__col">
														<h3>Candidate</h3>
														<dl class="kv">
															<dt>Submitted</dt>
															<dd>{fmtDateTime(expanded.createdAt)}</dd>

															<dt>Role</dt>
															<dd>
																{#if expanded.jobTitle}
																	<a
																		href="/hiring/{expanded.jobId}"
																		target="_blank"
																		rel="noopener">{expanded.jobTitle}</a
																	>
																{:else}
																	<span class="muted">job #{expanded.jobId}</span>
																{/if}
															</dd>

															<dt>Status</dt>
															<dd>
																{#if expanded.status}
																	<span class="badge" style={badgeStyle(expanded.status)}
																		>{expanded.status.name}</span
																	>
																	{#if expanded.statusUpdatedAt}
																		<span class="muted small">
																			· updated {fmtDateTime(expanded.statusUpdatedAt)}</span
																		>
																	{/if}
																{:else}
																	<span class="muted">—</span>
																{/if}
															</dd>
														</dl>

														{#if expanded.message}
															<h3>Message</h3>
															<p class="message">{expanded.message}</p>
														{/if}
													</div>

													<div class="detail__col">
														<h3>Attachments ({expanded.files?.length ?? 0})</h3>
														{#if expanded.files && expanded.files.length > 0}
															<ul class="files">
																{#each expanded.files as f (f.id)}
																	<li>
																		<a
																			href={applicationFileUrl(expanded.id, f.id)}
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

														<h3>History</h3>
														{#if expanded.history && expanded.history.length > 0}
															<ol class="history">
																{#each expanded.history as ev (ev.id)}
																	<li>
																		<div class="history__line">
																			{#if ev.fromStatusName}
																				<span class="muted small">{ev.fromStatusName}</span>
																				<span aria-hidden="true">→</span>
																			{/if}
																			<strong>{ev.toStatusName}</strong>
																			<span class="muted small"
																				>· {fmtDateTime(ev.createdAt)}</span
																			>
																		</div>
																		<div class="history__meta">
																			{#if ev.actorEmail}
																				<span class="muted small">by {ev.actorEmail}</span>
																			{:else}
																				<span class="muted small">automatic</span>
																			{/if}
																			{#if ev.note}
																				<span class="history__note">— {ev.note}</span>
																			{/if}
																		</div>
																	</li>
																{/each}
															</ol>
														{:else}
															<p class="muted small">No history yet.</p>
														{/if}
													</div>
												</div>
											</div>
										{/if}
									</td>
								</tr>
							{/if}
						{/each}
					</tbody>
				</table>
			</div>
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
		flex-wrap: wrap;
		gap: var(--space-3);
		align-items: flex-end;
		justify-content: space-between;
		margin-bottom: var(--space-4);
	}
	.filters {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
		align-items: flex-end;
	}
	.filter {
		display: grid;
		gap: 0.25rem;
		min-width: 160px;
	}
	.filter select {
		padding: 0.5rem 0.65rem;
		font: inherit;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
		color: var(--text);
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
	.caret-col {
		width: 1.5rem;
		padding-right: 0;
	}
	.caret {
		display: inline-block;
		color: var(--text-muted);
		transition: transform 120ms var(--ease-out);
	}
	.caret--open {
		transform: rotate(90deg);
		color: var(--accent);
	}

	.row {
		cursor: pointer;
	}
	.row:hover {
		background: var(--surface-muted, #f7f8fb);
	}
	.row:focus-visible {
		outline: 2px solid var(--accent);
		outline-offset: -2px;
	}
	.row--open {
		background: var(--accent-soft);
	}
	.row--open + .expansion .expansion__cell {
		border-top: none;
	}
	.row--open td {
		border-bottom-color: transparent;
	}
	.row a {
		color: var(--accent);
		text-decoration: none;
		font-weight: 500;
	}
	.row a:hover {
		text-decoration: underline;
	}

	.expansion .expansion__cell {
		padding: 0;
		background: var(--surface, #fafbfd);
		border-bottom: 1px solid var(--border);
	}

	.detail {
		padding: var(--space-4) var(--space-5);
		display: grid;
		gap: var(--space-4);
	}
	.detail__head {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: var(--space-3);
		flex-wrap: wrap;
	}
	.detail__title {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: var(--space-2);
	}
	.detail__title h2 {
		margin: 0;
		font-size: var(--text-xl);
	}
	.detail h3 {
		font-size: var(--text-sm);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		color: var(--text-muted);
		margin: 0 0 var(--space-2);
	}

	.status-changer {
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		padding: var(--space-3);
	}
	.status-changer__row {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-3);
		align-items: flex-end;
	}
	.field {
		display: grid;
		gap: 0.25rem;
		min-width: 180px;
	}
	.field--grow {
		flex: 1 1 280px;
		min-width: 240px;
	}
	.status-changer select,
	.status-changer input {
		padding: 0.55rem 0.7rem;
		font: inherit;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
		color: var(--text);
	}
	.status-changer select:focus,
	.status-changer input:focus {
		outline: none;
		border-color: var(--accent);
		box-shadow: 0 0 0 3px var(--accent-soft);
	}

	.detail__grid {
		display: grid;
		gap: var(--space-5);
		grid-template-columns: 1fr;
	}
	@media (min-width: 900px) {
		.detail__grid {
			grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
		}
	}
	.detail__col {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
		min-width: 0;
	}

	.kv {
		display: grid;
		grid-template-columns: max-content 1fr;
		column-gap: var(--space-3);
		row-gap: 0.4rem;
		margin: 0;
		font-size: var(--text-sm);
	}
	.kv dt {
		color: var(--text-muted);
	}
	.kv dd {
		margin: 0;
		word-break: break-word;
	}

	.badge {
		display: inline-block;
		padding: 0.15rem 0.55rem;
		font-size: 0.72rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		border-radius: 999px;
		background: color-mix(in srgb, var(--badge-color, #64748b) 14%, transparent);
		color: var(--badge-color, #64748b);
		border: 1px solid color-mix(in srgb, var(--badge-color, #64748b) 35%, transparent);
	}

	.message {
		white-space: pre-wrap;
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		padding: 0.7rem 0.9rem;
		font-size: var(--text-sm);
		line-height: 1.6;
		margin: 0;
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

	.history {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		font-size: var(--text-sm);
	}
	.history li {
		padding: 0.55rem 0.7rem;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
	}
	.history__line {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
		align-items: baseline;
	}
	.history__meta {
		margin-top: 0.2rem;
	}
	.history__note {
		font-size: var(--text-sm);
		color: var(--text);
	}

	.primary,
	.ghost {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0.55rem 0.95rem;
		font-size: var(--text-sm);
		font-weight: 500;
		border-radius: var(--radius-sm);
		cursor: pointer;
		border: 1px solid transparent;
		transition:
			background-color 120ms var(--ease-out),
			border-color 120ms var(--ease-out),
			color 120ms var(--ease-out);
		white-space: nowrap;
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
	button:disabled {
		opacity: 0.55;
		cursor: not-allowed;
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
