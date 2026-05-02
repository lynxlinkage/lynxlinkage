<script lang="ts">
	import { onMount } from 'svelte';
	import {
		ApiError,
		adminCreateStatus,
		adminDeleteStatus,
		adminListStatuses,
		adminUpdateStatus
	} from '$lib/api/client';
	import { auth } from '$lib/auth.svelte';
	import type {
		ApplicationStatus,
		ApplicationStatusKind,
		ApplicationStatusUpsertPayload
	} from '$lib/api/types';

	type Mode = 'list' | 'create' | 'edit';

	let mode = $state<Mode>('list');
	let items = $state<ApplicationStatus[]>([]);
	let loading = $state(true);
	let listError = $state<string | null>(null);

	let editingId = $state<number | null>(null);
	let formError = $state<string | null>(null);
	let saving = $state(false);

	const kindOptions: { value: ApplicationStatusKind; label: string; hint: string }[] = [
		{ value: 'open', label: 'Open', hint: 'In-flight pipeline state' },
		{ value: 'accept', label: 'Accept', hint: 'Terminal — candidate hired' },
		{ value: 'reject', label: 'Reject', hint: 'Terminal — candidate declined' }
	];

	const blank: ApplicationStatusUpsertPayload = {
		slug: '',
		name: '',
		kind: 'open',
		color: '#0ea5e9',
		displayOrder: 0,
		isDefault: false
	};
	let form = $state<ApplicationStatusUpsertPayload>({ ...blank });

	onMount(async () => {
		const user = await auth.load();
		if (!user) {
			window.location.href = '/admin/statuses';
			return;
		}
		if (user.role !== 'hr') {
			listError = 'Your account is not authorised to manage the workflow.';
			loading = false;
			return;
		}
		await refresh();
	});

	async function refresh() {
		loading = true;
		listError = null;
		try {
			items = await adminListStatuses();
		} catch (err) {
			if (err instanceof ApiError && err.status === 401) {
				window.location.href = '/admin/statuses';
				return;
			}
			listError = err instanceof Error ? err.message : 'Failed to load statuses';
		} finally {
			loading = false;
		}
	}

	function startCreate() {
		const nextOrder =
			items.length > 0
				? Math.max(...items.map((i) => i.displayOrder)) + 10
				: 0;
		form = { ...blank, displayOrder: nextOrder };
		editingId = null;
		formError = null;
		mode = 'create';
	}

	function startEdit(status: ApplicationStatus) {
		form = {
			slug: status.slug,
			name: status.name,
			kind: status.kind,
			color: status.color,
			displayOrder: status.displayOrder,
			isDefault: status.isDefault
		};
		editingId = status.id;
		formError = null;
		mode = 'edit';
	}

	function cancel() {
		editingId = null;
		formError = null;
		mode = 'list';
	}

	async function onSave(e: SubmitEvent) {
		e.preventDefault();
		if (saving) return;
		formError = null;
		saving = true;
		try {
			const payload: ApplicationStatusUpsertPayload = {
				slug: form.slug?.trim() || undefined,
				name: form.name.trim(),
				kind: form.kind,
				color: form.color?.trim() || undefined,
				displayOrder: form.displayOrder,
				isDefault: form.isDefault
			};
			if (mode === 'edit' && editingId != null) {
				await adminUpdateStatus(editingId, payload);
			} else {
				await adminCreateStatus(payload);
			}
			await refresh();
			cancel();
		} catch (err) {
			formError = err instanceof Error ? err.message : 'Save failed';
		} finally {
			saving = false;
		}
	}

	async function onDelete(status: ApplicationStatus) {
		const ok = window.confirm(
			`Delete "${status.name}"?\n\nThis only works if no applications still reference it.`
		);
		if (!ok) return;
		try {
			await adminDeleteStatus(status.id);
			await refresh();
		} catch (err) {
			listError = err instanceof Error ? err.message : 'Delete failed';
		}
	}

	async function onLogout() {
		await auth.logout(); // redirects browser to Authelia logout
	}

	function badgeStyle(color: string): string {
		const c = color || '#64748b';
		return `--badge-color: ${c};`;
	}
</script>

<svelte:head>
	<title>Workflow · Admin · Lynxlinkage</title>
	<meta name="robots" content="noindex,nofollow" />
</svelte:head>

<section class="admin">
	<div class="container">
		<header class="admin__header">
			<div>
				<p class="muted small">Admin</p>
				<h1>Hiring workflow</h1>
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
			<a class="admin__tab" href="/admin/applications">Applications</a>
			<a class="admin__tab admin__tab--active" href="/admin/statuses">Workflow</a>
		</nav>

		<p class="muted">
			HR-defined pipeline. New applications start in the status marked
			<strong>Default</strong>; ordering controls the dropdowns HR sees on the application
			detail page. Two terminal kinds (<em>accept</em> / <em>reject</em>) are tagged so the
			rest of the system can tell in-flight from closed.
		</p>

		{#if mode === 'list'}
			<div class="admin__toolbar">
				<button type="button" class="primary" onclick={startCreate}>+ New status</button>
				<button type="button" class="ghost" onclick={refresh}>Refresh</button>
			</div>

			{#if listError}
				<p class="error">{listError}</p>
			{/if}

			{#if loading}
				<p class="muted">Loading…</p>
			{:else if items.length === 0}
				<p class="muted">No statuses yet — create one to begin.</p>
			{:else}
				<div class="table-wrapper">
					<table class="table">
						<thead>
							<tr>
								<th>Order</th>
								<th>Name</th>
								<th>Slug</th>
								<th>Kind</th>
								<th>Default</th>
								<th class="actions">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#each items as st (st.id)}
								<tr>
									<td>{st.displayOrder}</td>
									<td>
										<span class="badge" style={badgeStyle(st.color)}>{st.name}</span>
									</td>
									<td><code>{st.slug}</code></td>
									<td>{st.kind}</td>
									<td>{st.isDefault ? 'Yes' : '—'}</td>
									<td class="actions">
										<button type="button" class="ghost small" onclick={() => startEdit(st)}>
											Edit
										</button>
										<button
											type="button"
											class="ghost small danger"
											onclick={() => onDelete(st)}
											title="Delete this status (only if unused)"
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
			<form class="status-form" onsubmit={onSave} novalidate>
				<div class="status-form__head">
					<h2>{mode === 'edit' ? 'Edit status' : 'New status'}</h2>
					<button type="button" class="ghost" onclick={cancel}>Cancel</button>
				</div>

				<div class="grid">
					<label class="field">
						<span>Name</span>
						<input type="text" required maxlength="80" bind:value={form.name} />
					</label>

					<label class="field">
						<span>Slug <em>(optional, derived from name)</em></span>
						<input type="text" maxlength="80" bind:value={form.slug} placeholder="auto" />
					</label>

					<label class="field">
						<span>Kind</span>
						<select bind:value={form.kind}>
							{#each kindOptions as opt (opt.value)}
								<option value={opt.value}>{opt.label} — {opt.hint}</option>
							{/each}
						</select>
					</label>

					<label class="field">
						<span>Display order</span>
						<input type="number" step="1" bind:value={form.displayOrder} />
					</label>

					<label class="field">
						<span>Badge colour</span>
						<input type="color" bind:value={form.color} />
					</label>

					<label class="field field--toggle">
						<input type="checkbox" bind:checked={form.isDefault} />
						<span>Default for new submissions</span>
					</label>
				</div>

				{#if formError}
					<p class="error">{formError}</p>
				{/if}

				<div class="status-form__actions">
					<button type="submit" class="primary" disabled={saving}>
						{saving ? 'Saving…' : mode === 'edit' ? 'Save changes' : 'Create status'}
					</button>
					<button type="button" class="ghost" onclick={cancel} disabled={saving}>
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
		margin: var(--space-4) 0 var(--space-3);
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
	.table .actions {
		display: flex;
		gap: 0.4rem;
		justify-content: flex-end;
	}
	.table th.actions {
		text-align: right;
	}
	.table code {
		font-family: var(--font-mono, ui-monospace, SFMono-Regular, monospace);
		font-size: 0.85em;
		color: var(--text-muted);
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

	.status-form {
		display: grid;
		gap: var(--space-4);
		background: var(--bg);
		padding: var(--space-5);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		margin-top: var(--space-4);
	}
	.status-form__head {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	.status-form__head h2 {
		margin: 0;
		font-size: var(--text-xl);
	}
	.status-form__actions {
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
	.field em {
		font-style: normal;
		color: var(--text-muted);
		font-weight: 400;
	}
	.field input[type='text'],
	.field input[type='number'],
	.field select {
		width: 100%;
		padding: 0.6rem 0.75rem;
		font: inherit;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
		color: var(--text);
	}
	.field input[type='color'] {
		width: 60px;
		height: 36px;
		padding: 0;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		background: var(--bg);
	}
	.field input:focus,
	.field select:focus {
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
		opacity: 0.55;
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
